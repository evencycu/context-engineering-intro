package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTokenManagerAdapter_GetConnectorToken(t *testing.T) {
	// Create mock token manager
	mockTM := &MockTokenManager{}

	// Create adapter
	adapter := &TokenManagerAdapter{tm: mockTM}

	// Mock the underlying token manager
	mockTM.On("GetConnectorToken", mock.Anything, "test_bot_id", "test_tenant_id").Return("test_token", nil)

	// Test
	ctx := context.Background()
	token, err := adapter.GetConnectorToken(ctx, "test_bot_id", "test_tenant_id")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "test_token", token)
	mockTM.AssertExpectations(t)
}

func TestTokenManagerAdapter_GetGraphToken(t *testing.T) {
	// Create mock token manager
	mockTM := &MockTokenManager{}

	// Create adapter
	adapter := &TokenManagerAdapter{tm: mockTM}

	// Mock the underlying token manager
	mockTM.On("GetGraphToken", mock.Anything, "test_bot_id", "test_tenant_id").Return("test_graph_token", nil)

	// Test
	ctx := context.Background()
	token, err := adapter.GetGraphToken(ctx, "test_bot_id", "test_tenant_id")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "test_graph_token", token)
	mockTM.AssertExpectations(t)
}

func TestTokenManagerAdapter_InvalidateToken(t *testing.T) {
	// Create mock token manager
	mockTM := &MockTokenManager{}

	// Create adapter
	adapter := &TokenManagerAdapter{tm: mockTM}

	// Mock the underlying token manager
	mockTM.On("InvalidateToken", mock.Anything, "test_bot_id", "test_tenant_id", TokenTypeConnector).Return(nil)

	// Test
	ctx := context.Background()
	err := adapter.InvalidateToken(ctx, "test_bot_id", "test_tenant_id", TokenTypeConnector)

	// Assertions
	assert.NoError(t, err)
	mockTM.AssertExpectations(t)
}

func TestTokenManagerAdapter_RefreshToken(t *testing.T) {
	// Create mock token manager
	mockTM := &MockTokenManager{}

	// Create adapter
	adapter := &TokenManagerAdapter{tm: mockTM}

	// Mock the underlying token manager
	mockTM.On("RefreshToken", mock.Anything, "test_bot_id", "test_tenant_id", TokenTypeConnector).Return("refreshed_token", nil)

	// Test
	ctx := context.Background()
	token, err := adapter.RefreshToken(ctx, "test_bot_id", "test_tenant_id", TokenTypeConnector)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "refreshed_token", token)
	mockTM.AssertExpectations(t)
}

func TestTokenManagerAdapter_GetBotConfig(t *testing.T) {
	// Create mock token manager
	mockTM := &MockTokenManager{}

	// Create adapter
	adapter := &TokenManagerAdapter{tm: mockTM}

	// Mock the underlying token manager
	expectedConfig := &BotConfig{
		BotID:       "test_bot_id",
		AppPassword: "test_password",
		TenantID:    "test_tenant_id",
		CompanyID:   "test_company_id",
	}
	mockTM.On("GetBotConfig", mock.Anything, "test_bot_id").Return(expectedConfig, nil)

	// Test
	ctx := context.Background()
	config, err := adapter.GetBotConfig(ctx, "test_bot_id")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
	mockTM.AssertExpectations(t)
}

// MockTokenManager is a mock implementation of TokenManager
type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GetConnectorToken(ctx context.Context, botID, tenantID string) (string, error) {
	args := m.Called(ctx, botID, tenantID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) GetGraphToken(ctx context.Context, botID, tenantID string) (string, error) {
	args := m.Called(ctx, botID, tenantID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) InvalidateToken(ctx context.Context, botID, tenantID string, tokenType TokenType) error {
	args := m.Called(ctx, botID, tenantID, tokenType)
	return args.Error(0)
}

func (m *MockTokenManager) RefreshToken(ctx context.Context, botID, tenantID string, tokenType TokenType) (string, error) {
	args := m.Called(ctx, botID, tenantID, tokenType)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) GetBotConfig(ctx context.Context, botID string) (*BotConfig, error) {
	args := m.Called(ctx, botID)
	return args.Get(0).(*BotConfig), args.Error(1)
}
