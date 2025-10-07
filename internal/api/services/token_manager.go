package services

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

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenType represents the type of token
type TokenType string

const (
	TokenTypeConnector TokenType = "connector"
	TokenTypeGraph     TokenType = "graph"
)

// TokenInfo represents cached token information
type TokenInfo struct {
	AccessToken   string    `json:"access_token"`
	TokenType     string    `json:"token_type"`
	ExpiresIn     int       `json:"expires_in"`
	ExpiresAt     int64     `json:"expires_at"`
	Scope         string    `json:"scope"`
	BotID         string    `json:"bot_id"`
	TenantID      string    `json:"tenant_id"`
	TokenTypeEnum TokenType `json:"token_type_enum"`
}

// BotConfig represents bot configuration
type BotConfig struct {
	BotID       string `json:"bot_id"`
	AppPassword string `json:"app_password"`
	TenantID    string `json:"tenant_id"`
	CompanyID   string `json:"company_id"`
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
	redis        *redis.Client
	botConfigs   map[string]*BotConfig
	cacheTTL     time.Duration
	refreshAhead time.Duration
	maxRetries   int
}

// NewTokenManager creates a new token manager
func NewTokenManager(redis *redis.Client) TokenManager {
	return &redisTokenManager{
		redis:        redis,
		botConfigs:   make(map[string]*BotConfig),
		cacheTTL:     50 * time.Minute, // 50 minutes cache
		refreshAhead: 5 * time.Minute,  // Refresh 5 minutes before expiry
		maxRetries:   3,
	}
}

// GetConnectorToken gets a Bot Framework connector token
func (tm *redisTokenManager) GetConnectorToken(ctx context.Context, botID, tenantID string) (string, error) {
	return tm.getToken(ctx, botID, tenantID, TokenTypeConnector)
}

// GetGraphToken gets a Microsoft Graph API token
func (tm *redisTokenManager) GetGraphToken(ctx context.Context, botID, tenantID string) (string, error) {
	return tm.getToken(ctx, botID, tenantID, TokenTypeGraph)
}

// getToken retrieves a token from cache or fetches a new one
func (tm *redisTokenManager) getToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error) {
	cacheKey := tm.getCacheKey(botID, tenantID, tokenType)

	// Try to get from cache first
	cached, err := tm.getCachedToken(ctx, cacheKey)
	if err == nil && cached != nil && tm.isTokenValid(cached) {
		return cached.AccessToken, nil
	}

	// If cache miss or token invalid, fetch new token
	return tm.fetchAndCacheToken(ctx, botID, tenantID, tokenType)
}

// getCachedToken retrieves token from Redis cache
func (tm *redisTokenManager) getCachedToken(ctx context.Context, cacheKey string) (*TokenInfo, error) {
	val, err := tm.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var tokenInfo TokenInfo
	if err := json.Unmarshal([]byte(val), &tokenInfo); err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

// isTokenValid checks if token is still valid
func (tm *redisTokenManager) isTokenValid(token *TokenInfo) bool {
	if token == nil {
		return false
	}

	// Check if token expires within refreshAhead time
	now := time.Now().Unix()
	return token.ExpiresAt > now+int64(tm.refreshAhead.Seconds())
}

// fetchAndCacheToken fetches a new token and caches it
func (tm *redisTokenManager) fetchAndCacheToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error) {
	// Use distributed lock to prevent concurrent token refresh
	lockKey := fmt.Sprintf("teams:token:refresh:%s:%s:%s", botID, tenantID, tokenType)
	lockValue := uuid.New().String()

	// Try to acquire lock with 30 second timeout
	acquired, err := tm.redis.SetNX(ctx, lockKey, lockValue, 30*time.Second).Result()
	if err != nil {
		return "", fmt.Errorf("failed to acquire refresh lock: %w", err)
	}

	if !acquired {
		// Another process is refreshing, wait and retry
		time.Sleep(1 * time.Second)
		return tm.getToken(ctx, botID, tenantID, tokenType)
	}

	// Ensure lock is released
	defer func() {
		// Use Lua script to safely release lock
		script := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
		tm.redis.Eval(ctx, script, []string{lockKey}, lockValue)
	}()

	// Get bot configuration
	config, err := tm.GetBotConfig(ctx, botID)
	if err != nil {
		return "", fmt.Errorf("failed to get bot config: %w", err)
	}

	// Fetch new token
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

	return tokenInfo.AccessToken, nil
}

// fetchTokenFromMicrosoft fetches token from Microsoft OAuth endpoint
func (tm *redisTokenManager) fetchTokenFromMicrosoft(ctx context.Context, config *BotConfig, tokenType TokenType) (*TokenInfo, error) {
	var scope string
	switch tokenType {
	case TokenTypeConnector:
		scope = "https://api.botframework.com/.default"
	case TokenTypeGraph:
		scope = "https://graph.microsoft.com/.default"
	default:
		return nil, fmt.Errorf("unsupported token type: %s", tokenType)
	}

	// Prepare token request
	data := url.Values{}
	data.Set("client_id", config.BotID)
	data.Set("client_secret", config.AppPassword)
	data.Set("scope", scope)
	data.Set("grant_type", "client_credentials")

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", config.TenantID)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed: %d %s", resp.StatusCode, string(body))
	}

	// Parse token response
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	// Create token info
	now := time.Now()
	tokenInfo := &TokenInfo{
		AccessToken:   tokenResp.AccessToken,
		TokenType:     tokenResp.TokenType,
		ExpiresIn:     tokenResp.ExpiresIn,
		ExpiresAt:     now.Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix(),
		Scope:         scope,
		BotID:         config.BotID,
		TenantID:      config.TenantID,
		TokenTypeEnum: tokenType,
	}

	return tokenInfo, nil
}

// cacheToken stores token in Redis cache
func (tm *redisTokenManager) cacheToken(ctx context.Context, cacheKey string, token *TokenInfo) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	return tm.redis.Set(ctx, cacheKey, data, tm.cacheTTL).Err()
}

// getCacheKey generates cache key for token
func (tm *redisTokenManager) getCacheKey(botID, tenantID string, tokenType TokenType) string {
	return fmt.Sprintf("teams:token:%s:%s:%s", botID, tenantID, tokenType)
}

// InvalidateToken invalidates a cached token
func (tm *redisTokenManager) InvalidateToken(ctx context.Context, botID, tenantID string, tokenType TokenType) error {
	cacheKey := tm.getCacheKey(botID, tenantID, tokenType)
	return tm.redis.Del(ctx, cacheKey).Err()
}

// RefreshToken forces a token refresh
func (tm *redisTokenManager) RefreshToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error) {
	// Invalidate existing token first
	if err := tm.InvalidateToken(ctx, botID, tenantID, tokenType); err != nil {
		return "", fmt.Errorf("failed to invalidate token: %w", err)
	}

	// Fetch new token
	return tm.fetchAndCacheToken(ctx, botID, tenantID, tokenType)
}

// GetBotConfig retrieves bot configuration
func (tm *redisTokenManager) GetBotConfig(ctx context.Context, botID string) (*BotConfig, error) {
	// Check if config is already loaded
	if config, exists := tm.botConfigs[botID]; exists {
		return config, nil
	}

	// Try to load from database first
	config, err := tm.loadBotConfigFromDB(ctx, botID)
	if err != nil {
		// Fallback to environment variables
		config = tm.getFallbackConfig(botID)
		if config.AppPassword == "" || config.TenantID == "" {
			return nil, fmt.Errorf("bot configuration not found for bot ID: %s", botID)
		}
	}

	// Cache the config
	tm.botConfigs[botID] = config
	return config, nil
}

// loadBotConfigFromDB loads bot configuration from database
func (tm *redisTokenManager) loadBotConfigFromDB(ctx context.Context, botID string) (*BotConfig, error) {
	// TODO: Implement database loading
	// This would query the teams_bots table for the bot configuration
	return nil, fmt.Errorf("database loading not implemented")
}

// getFallbackConfig gets configuration from environment variables
func (tm *redisTokenManager) getFallbackConfig(botID string) *BotConfig {
	return &BotConfig{
		BotID:       botID,
		AppPassword: getTokenEnv("TEAMS_BOT_APP_PASSWORD", ""),
		TenantID:    getTokenEnv("TEAMS_TENANT_ID", ""),
		CompanyID:   "",
	}
}

// getTokenEnv gets environment variable with default value for token manager
func getTokenEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// TokenManagerAdapter adapts services.TokenManager to actor.TokenManager interface
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
	if tt, ok := tokenType.(TokenType); ok {
		return a.tm.InvalidateToken(ctx, botID, tenantID, tt)
	}
	// If it's a string, convert to TokenType
	if str, ok := tokenType.(string); ok {
		return a.tm.InvalidateToken(ctx, botID, tenantID, TokenType(str))
	}
	return fmt.Errorf("invalid token type: %T", tokenType)
}

// RefreshToken implements actor.TokenManager interface
func (a *TokenManagerAdapter) RefreshToken(ctx context.Context, botID, tenantID string, tokenType interface{}) (string, error) {
	// Convert interface{} to TokenType
	if tt, ok := tokenType.(TokenType); ok {
		return a.tm.RefreshToken(ctx, botID, tenantID, tt)
	}
	// If it's a string, convert to TokenType
	if str, ok := tokenType.(string); ok {
		return a.tm.RefreshToken(ctx, botID, tenantID, TokenType(str))
	}
	return "", fmt.Errorf("invalid token type: %T", tokenType)
}

// GetBotConfig implements actor.TokenManager interface
func (a *TokenManagerAdapter) GetBotConfig(ctx context.Context, botID string) (interface{}, error) {
	return a.tm.GetBotConfig(ctx, botID)
}
