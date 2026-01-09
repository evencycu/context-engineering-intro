package actor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TeamsSendPayload represents the data needed to send a message to Teams
type TeamsSendPayload struct {
	ConversationID string
	ServiceURL     string
	Message        string
	Attachments    []map[string]any
	RecipientID    string
	FromID         string
	BotAppID       string
	BotAppPassword string
	TenantID       string
}

// TeamsSendResult contains the outcome of a Teams message send attempt
type TeamsSendResult struct {
	Success        bool
	TeamsMessageID string
	RetryAfter     int // From Retry-After header, if any
	StatusCode     int
	Error          error
}

// TokenManager interface for token management
type TokenManager interface {
	GetConnectorToken(ctx context.Context, botID, tenantID string) (string, error)
	GetGraphToken(ctx context.Context, botID, tenantID string) (string, error)
	InvalidateToken(ctx context.Context, botID, tenantID string, tokenType interface{}) error
	RefreshToken(ctx context.Context, botID, tenantID string, tokenType interface{}) (string, error)
	GetBotConfig(ctx context.Context, botID string) (interface{}, error)
}

// TeamsConnector defines the interface for interacting with Microsoft Teams API
type TeamsConnector interface {
	Send(ctx context.Context, payload *TeamsSendPayload) (*TeamsSendResult, error)
}

// defaultTeamsConnector implements TeamsConnector interface
type defaultTeamsConnector struct {
	httpClient     *http.Client
	tokenManager   TokenManager
	circuitBreaker CircuitBreaker
}

// NewTeamsConnector creates a new Teams connector
func NewTeamsConnector(tokenManager TokenManager, circuitBreaker CircuitBreaker) TeamsConnector {
	return &defaultTeamsConnector{
		httpClient:     &http.Client{Timeout: 30 * time.Second},
		tokenManager:   tokenManager,
		circuitBreaker: circuitBreaker,
	}
}

// Send sends a notification to Microsoft Teams via Bot Framework Connector API
func (c *defaultTeamsConnector) Send(ctx context.Context, payload *TeamsSendPayload) (*TeamsSendResult, error) {
	// 1. Check Circuit Breaker
	if c.circuitBreaker != nil {
		open, err := c.circuitBreaker.IsOpen(ctx)
		if err == nil && open {
			return &TeamsSendResult{
				Success: false,
				Error:   fmt.Errorf("circuit breaker is open"),
			}, nil
		}
	}

	// 2. Get Access Token
	token, err := c.getAccessToken(ctx, payload)
	if err != nil {
		return &TeamsSendResult{
			Success:    false,
			StatusCode: http.StatusUnauthorized,
			Error:      fmt.Errorf("failed to get access token: %w", err),
		}, nil
	}

	// 3. Construct Bot Framework Payload
	teamsPayload := map[string]any{
		"type": "message",
		"text": payload.Message,
		"from": map[string]any{
			"id":   payload.RecipientID, // In Bot Framework 'from' is the bot
			"name": "Notification Bot",
		},
		"recipient": map[string]any{
			"id": payload.FromID, // 'recipient' is the user/conversation
		},
		"conversation": map[string]any{
			"id": payload.ConversationID,
		},
		"channelId":  "msteams",
		"serviceUrl": payload.ServiceURL,
	}

	if len(payload.Attachments) > 0 {
		teamsPayload["attachments"] = payload.Attachments
	}

	jsonBody, err := json.Marshal(teamsPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal teams payload: %w", err)
	}

	// 4. Execute HTTP Request
	apiURL := fmt.Sprintf("%s/v3/conversations/%s/activities", 
		strings.TrimRight(payload.ServiceURL, "/"), 
		url.PathEscape(payload.ConversationID))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.circuitBreaker != nil {
			_ = c.circuitBreaker.RecordFailure(ctx, 0)
		}
		return &TeamsSendResult{Success: false, Error: err}, nil
	}
	defer resp.Body.Close()

	// 5. Process Response
	result := &TeamsSendResult{
		StatusCode: resp.StatusCode,
		Success:    resp.StatusCode >= 200 && resp.StatusCode < 300,
	}

	if !result.Success {
		// Handle failure (e.g. 429 Rate Limit)
		if resp.StatusCode == http.StatusTooManyRequests {
			if raStr := resp.Header.Get("Retry-After"); raStr != "" {
				var ra int
				fmt.Sscanf(raStr, "%d", &ra)
				result.RetryAfter = ra
			}
		}

		body, _ := io.ReadAll(resp.Body)
		result.Error = fmt.Errorf("teams api error (status %d): %s", resp.StatusCode, string(body))

		if c.circuitBreaker != nil {
			_ = c.circuitBreaker.RecordFailure(ctx, result.RetryAfter)
		}
	} else {
		// Handle success
		var respData map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&respData); err == nil {
			if id, ok := respData["id"].(string); ok {
				result.TeamsMessageID = id
			}
		}

		if c.circuitBreaker != nil {
			_ = c.circuitBreaker.RecordSuccess(ctx)
		}
	}

	return result, nil
}

func (c *defaultTeamsConnector) getAccessToken(ctx context.Context, payload *TeamsSendPayload) (string, error) {
	tenantID := payload.TenantID
	if tenantID == "" {
		tenantID = "botframework.com"
	}

	// Try Token Manager first
	if c.tokenManager != nil {
		token, err := c.tokenManager.GetConnectorToken(ctx, payload.BotAppID, tenantID)
		if err == nil && token != "" {
			return token, nil
		}
		log.Printf("TeamsConnector: TokenManager failed, falling back to direct fetch: %v", err)
	}

	// Fallback to direct fetch
	return c.fetchTokenDirect(ctx, payload, tenantID)
}

func (c *defaultTeamsConnector) fetchTokenDirect(ctx context.Context, payload *TeamsSendPayload, tenantID string) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", payload.BotAppID)
	form.Set("client_secret", payload.BotAppPassword)
	form.Set("scope", "https://api.botframework.com/.default")

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", url.PathEscape(tenantID))
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed (%d): %s", resp.StatusCode, string(body))
	}

	var tr struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	return tr.AccessToken, nil
}