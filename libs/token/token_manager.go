package token

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/actor"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenType represents the type of token
type TokenType string

const (
	ConnectorToken TokenType = "connector"
	GraphToken     TokenType = "graph"
)

// TokenInfo represents cached token information
type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Type      TokenType `json:"type"`
}

// BotConfig represents bot configuration
type BotConfig struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	TenantID  string `json:"tenant_id"`
}

// TokenManager defines the interface for token management
type TokenManager interface {
	GetConnectorToken(ctx context.Context, botID, tenantID string) (string, error)
	GetGraphToken(ctx context.Context, botID, tenantID string) (string, error)
	InvalidateToken(ctx context.Context, botID, tenantID string, tokenType TokenType) error
	RefreshToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error)
	GetBotConfig(ctx context.Context, botID string) (*BotConfig, error)
}

// redisTokenManager implements TokenManager using Redis cache
type redisTokenManager struct {
	redis  *redis.Client
	config map[string]*BotConfig
}

// NewTokenManager creates a new token manager
func NewTokenManager(redis *redis.Client) TokenManager {
	return &redisTokenManager{
		redis:  redis,
		config: make(map[string]*BotConfig),
	}
}

// GetConnectorToken gets a connector token for Teams Bot Framework
func (tm *redisTokenManager) GetConnectorToken(ctx context.Context, botID, tenantID string) (string, error) {
	return tm.getToken(ctx, botID, tenantID, ConnectorToken)
}

// GetGraphToken gets a Microsoft Graph token
func (tm *redisTokenManager) GetGraphToken(ctx context.Context, botID, tenantID string) (string, error) {
	return tm.getToken(ctx, botID, tenantID, GraphToken)
}

// getToken gets a token from cache or fetches a new one
func (tm *redisTokenManager) getToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error) {
	cacheKey := tm.getCacheKey(botID, tenantID, tokenType)
	
	// Try to get from cache first
	if token, err := tm.getCachedToken(ctx, cacheKey); err == nil && tm.isTokenValid(token) {
		return token.Token, nil
	}
	
	// Fetch and cache new token
	return tm.fetchAndCacheToken(ctx, botID, tenantID, tokenType)
}

// getCachedToken retrieves a token from Redis cache
func (tm *redisTokenManager) getCachedToken(ctx context.Context, cacheKey string) (*TokenInfo, error) {
	val, err := tm.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	
	var token TokenInfo
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return nil, err
	}
	
	return &token, nil
}

// isTokenValid checks if a token is still valid
func (tm *redisTokenManager) isTokenValid(token *TokenInfo) bool {
	// Add 5 minute buffer to avoid edge cases
	return time.Now().Add(5 * time.Minute).Before(token.ExpiresAt)
}

// fetchAndCacheToken fetches a new token and caches it
func (tm *redisTokenManager) fetchAndCacheToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error) {
	// Get bot configuration
	config, err := tm.GetBotConfig(ctx, botID)
	if err != nil {
		return "", fmt.Errorf("failed to get bot config: %w", err)
	}
	
	// Fetch token from Microsoft
	tokenInfo, err := tm.fetchTokenFromMicrosoft(ctx, config, tokenType)
	if err != nil {
		return "", fmt.Errorf("failed to fetch token: %w", err)
	}
	
	// Cache the token
	cacheKey := tm.getCacheKey(botID, tenantID, tokenType)
	if err := tm.cacheToken(ctx, cacheKey, tokenInfo); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to cache token: %v\n", err)
	}
	
	return tokenInfo.Token, nil
}

// fetchTokenFromMicrosoft fetches a token from Microsoft OAuth2 endpoint
func (tm *redisTokenManager) fetchTokenFromMicrosoft(ctx context.Context, config *BotConfig, tokenType TokenType) (*TokenInfo, error) {
	var scope, tokenURL string
	
	switch tokenType {
	case ConnectorToken:
		scope = "https://api.botframework.com/.default"
		tokenURL = "https://login.microsoftonline.com/" + config.TenantID + "/oauth2/v2.0/token"
	case GraphToken:
		scope = "https://graph.microsoft.com/.default"
		tokenURL = "https://login.microsoftonline.com/" + config.TenantID + "/oauth2/v2.0/token"
	default:
		return nil, fmt.Errorf("unsupported token type: %s", tokenType)
	}
	
	// Prepare form data
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", config.AppID)
	form.Set("client_secret", config.AppSecret)
	form.Set("scope", scope)
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed: %s: %s", resp.Status, string(body))
	}
	
	// Parse response
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	
	return &TokenInfo{
		Token:     tokenResponse.AccessToken,
		ExpiresAt: expiresAt,
		Type:      tokenType,
	}, nil
}

// cacheToken stores a token in Redis cache
func (tm *redisTokenManager) cacheToken(ctx context.Context, cacheKey string, token *TokenInfo) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}
	
	// Cache for the duration of the token minus 5 minutes buffer
	ttl := time.Until(token.ExpiresAt) - 5*time.Minute
	if ttl <= 0 {
		ttl = time.Hour // Default to 1 hour if calculation fails
	}
	
	return tm.redis.Set(ctx, cacheKey, string(data), ttl).Err()
}

// getCacheKey generates a cache key for a token
func (tm *redisTokenManager) getCacheKey(botID, tenantID string, tokenType TokenType) string {
	return fmt.Sprintf("token:%s:%s:%s", botID, tenantID, tokenType)
}

// InvalidateToken removes a token from cache
func (tm *redisTokenManager) InvalidateToken(ctx context.Context, botID, tenantID string, tokenType TokenType) error {
	cacheKey := tm.getCacheKey(botID, tenantID, tokenType)
	return tm.redis.Del(ctx, cacheKey).Err()
}

// RefreshToken forces a token refresh
func (tm *redisTokenManager) RefreshToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error) {
	// Invalidate existing token
	if err := tm.InvalidateToken(ctx, botID, tenantID, tokenType); err != nil {
		// Log error but don't fail
		fmt.Printf("Warning: failed to invalidate token: %v\n", err)
	}
	
	// Fetch new token
	return tm.getToken(ctx, botID, tenantID, tokenType)
}

// GetBotConfig gets bot configuration
func (tm *redisTokenManager) GetBotConfig(ctx context.Context, botID string) (*BotConfig, error) {
	// Check cache first
	if config, exists := tm.config[botID]; exists {
		return config, nil
	}
	
	// Load from database (this would need to be implemented)
	config, err := tm.loadBotConfigFromDB(ctx, botID)
	if err != nil {
		// Fallback to environment variables
		config = tm.getFallbackConfig(botID)
	}
	
	// Cache the config
	tm.config[botID] = config
	return config, nil
}

// loadBotConfigFromDB loads bot configuration from database
func (tm *redisTokenManager) loadBotConfigFromDB(ctx context.Context, botID string) (*BotConfig, error) {
	// This would need to be implemented to load from database
	// For now, return fallback config
	return tm.getFallbackConfig(botID), nil
}

// getFallbackConfig gets configuration from environment variables
func (tm *redisTokenManager) getFallbackConfig(botID string) *BotConfig {
	return &BotConfig{
		AppID:     os.Getenv("TEAMS_BOT_APP_ID"),
		AppSecret: os.Getenv("TEAMS_BOT_APP_PASSWORD"),
		TenantID:  os.Getenv("TEAMS_TENANT_ID"),
	}
}

// TokenManagerAdapter adapts token.TokenManager to actor.TokenManager interface
type TokenManagerAdapter struct {
	tm TokenManager
}

// NewTokenManagerAdapter creates a new adapter
func NewTokenManagerAdapter(tm TokenManager) *TokenManagerAdapter {
	return &TokenManagerAdapter{tm: tm}
}

// GetConnectorToken implements actor.TokenManager interface
func (a *TokenManagerAdapter) GetConnectorToken(ctx context.Context, botID, tenantID string) (string, error) {
	return a.tm.GetConnectorToken(ctx, botID, tenantID)
}

// GetGraphToken implements actor.TokenManager interface
func (a *TokenManagerAdapter) GetGraphToken(ctx context.Context, botID, tenantID string) (string, error) {
	return a.tm.GetGraphToken(ctx, botID, tenantID)
}

// InvalidateToken implements actor.TokenManager interface
func (a *TokenManagerAdapter) InvalidateToken(ctx context.Context, botID, tenantID string, tokenType interface{}) error {
	// Convert interface{} to TokenType
	var t TokenType
	switch v := tokenType.(type) {
	case string:
		t = TokenType(v)
	case TokenType:
		t = v
	default:
		return fmt.Errorf("invalid token type: %T", tokenType)
	}
	return a.tm.InvalidateToken(ctx, botID, tenantID, t)
}

// RefreshToken implements actor.TokenManager interface
func (a *TokenManagerAdapter) RefreshToken(ctx context.Context, botID, tenantID string, tokenType interface{}) (string, error) {
	// Convert interface{} to TokenType
	var t TokenType
	switch v := tokenType.(type) {
	case string:
		t = TokenType(v)
	case TokenType:
		t = v
	default:
		return "", fmt.Errorf("invalid token type: %T", tokenType)
	}
	return a.tm.RefreshToken(ctx, botID, tenantID, t)
}

// GetBotConfig implements actor.TokenManager interface
func (a *TokenManagerAdapter) GetBotConfig(ctx context.Context, botID string) (interface{}, error) {
	return a.tm.GetBotConfig(ctx, botID)
}

// BillingServiceAdapter adapts services.BillingService to actor.BillingService interface
type BillingServiceAdapter struct {
	bs services.BillingService
}

// NewBillingServiceAdapter creates a new adapter
func NewBillingServiceAdapter(bs services.BillingService) *BillingServiceAdapter {
	return &BillingServiceAdapter{bs: bs}
}

// RecordUsage implements actor.BillingService interface
func (a *BillingServiceAdapter) RecordUsage(ctx context.Context, req *actor.RecordUsageRequest) error {
	// Convert actor.RecordUsageRequest to services.RecordUsageRequest
	projectID := uuid.Nil
	if req.ProjectID != nil {
		projectID = *req.ProjectID
	}
	
	return a.bs.RecordUsage(ctx, &services.RecordUsageRequest{
		ProjectID:   projectID,
		UserID:      req.UserID,
		RecordType:  req.RecordType,
		Quantity:    req.Quantity,
		UnitCost:    req.UnitPrice,
		TotalCost:   req.TotalCost,
		Metadata:    map[string]any{
			"notification_id": req.NotificationID,
			"bot_id": req.BotID,
			"bot_type": req.BotType,
			"billing_period": req.BillingPeriod,
		},
	})
}
