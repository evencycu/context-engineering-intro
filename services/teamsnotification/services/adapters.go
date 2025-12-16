package services

import (
	"context"
	"fmt"

	"github.com/evencycu/TeamsNotifyGoV3/libs/token"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/actor"
	"github.com/google/uuid"
)

// TokenManagerAdapter adapts libs/token.TokenManager to actor.TokenManager interface
type TokenManagerAdapter struct {
	tm token.TokenManager
}

// NewTokenManagerAdapter creates a new adapter
func NewTokenManagerAdapter(tm token.TokenManager) *TokenManagerAdapter {
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
	// Convert interface{} to token.TokenType
	var t token.TokenType
	switch v := tokenType.(type) {
	case string:
		t = token.TokenType(v)
	case token.TokenType:
		t = v
	default:
		return fmt.Errorf("invalid token type: %T", tokenType)
	}
	return a.tm.InvalidateToken(ctx, botID, tenantID, t)
}

// RefreshToken implements actor.TokenManager interface
func (a *TokenManagerAdapter) RefreshToken(ctx context.Context, botID, tenantID string, tokenType interface{}) (string, error) {
	// Convert interface{} to token.TokenType
	var t token.TokenType
	switch v := tokenType.(type) {
	case string:
		t = token.TokenType(v)
	case token.TokenType:
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
	bs BillingService
}

// NewBillingServiceAdapter creates a new adapter
func NewBillingServiceAdapter(bs BillingService) *BillingServiceAdapter {
	return &BillingServiceAdapter{bs: bs}
}

// RecordUsage implements actor.BillingService interface
func (a *BillingServiceAdapter) RecordUsage(ctx context.Context, req *actor.RecordUsageRequest) error {
	// Convert actor.RecordUsageRequest to services.RecordUsageRequest
	projectID := uuid.Nil
	if req.ProjectID != nil {
		projectID = *req.ProjectID
	}
	
	return a.bs.RecordUsage(ctx, &RecordUsageRequest{
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
