package services

import (
	"context"
	"time"
)

// MockRedisClient 模擬 Redis 客戶端，用於測試
type MockRedisClient struct {
	data map[string]interface{}
}

// NewMockRedisClient 創建新的 Mock Redis 客戶端
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]interface{}),
	}
}

// Set 模擬 Redis SET 操作
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.data[key] = value
	return nil
}

// Get 模擬 Redis GET 操作
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if val, exists := m.data[key]; exists {
		if str, ok := val.(string); ok {
			return str, nil
		}
	}
	return "", nil
}

// Incr 模擬 Redis INCR 操作
func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	if val, exists := m.data[key]; exists {
		if count, ok := val.(int64); ok {
			m.data[key] = count + 1
			return count + 1, nil
		}
	}
	m.data[key] = int64(1)
	return 1, nil
}

// Ping 模擬 Redis PING 操作
func (m *MockRedisClient) Ping(ctx context.Context) error {
	return nil
}
