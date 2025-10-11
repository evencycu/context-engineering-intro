package services

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// MetricsService provides system metrics collection
type MetricsService interface {
	GetMetrics(ctx context.Context) (*MetricsResponse, error)
}

// MetricsResponse represents the metrics response
type MetricsResponse struct {
	Timestamp  time.Time         `json:"timestamp"`
	System     SystemMetrics     `json:"system"`
	TokenCache TokenCacheMetrics `json:"token_cache"`
	ActorPool  ActorPoolMetrics  `json:"actor_pool"`
	Database   DatabaseMetrics   `json:"database"`
	Redis      RedisMetrics      `json:"redis"`
}

// SystemMetrics represents system-level metrics
type SystemMetrics struct {
	Uptime      string  `json:"uptime"`
	MemoryUsage uint64  `json:"memory_usage"`
	CPUUsage    float64 `json:"cpu_usage"`
	Goroutines  int     `json:"goroutines"`
}

// TokenCacheMetrics represents token cache metrics
type TokenCacheMetrics struct {
	CacheHits    int64 `json:"cache_hits"`
	CacheMisses  int64 `json:"cache_misses"`
	ActiveTokens int64 `json:"active_tokens"`
	RefreshCount int64 `json:"refresh_count"`
	ErrorCount   int64 `json:"error_count"`
}

// ActorPoolMetrics represents actor pool metrics
type ActorPoolMetrics struct {
	ActiveActors   int64 `json:"active_actors"`
	MaxActors      int64 `json:"max_actors"`
	QueueSize      int64 `json:"queue_size"`
	ProcessedCount int64 `json:"processed_count"`
	ErrorCount     int64 `json:"error_count"`
}

// DatabaseMetrics represents database metrics
type DatabaseMetrics struct {
	OpenConnections  int64   `json:"open_connections"`
	IdleConnections  int64   `json:"idle_connections"`
	MaxConnections   int64   `json:"max_connections"`
	QueryCount       int64   `json:"query_count"`
	QueryTimeAvg     float64 `json:"query_time_avg_ms"`
	SlowQueries      int64   `json:"slow_queries"`
	ConnectionErrors int64   `json:"connection_errors"`
	LastQueryTime    string  `json:"last_query_time"`
}

// RedisMetrics represents Redis metrics
type RedisMetrics struct {
	Connected        bool    `json:"connected"`
	MemoryUsage      string  `json:"memory_usage"`
	KeyCount         int64   `json:"key_count"`
	OperationsPerSec float64 `json:"operations_per_second"`
	HitRate          float64 `json:"hit_rate"`
	MissRate         float64 `json:"miss_rate"`
	EvictedKeys      int64   `json:"evicted_keys"`
	ExpiredKeys      int64   `json:"expired_keys"`
	ConnectedClients int64   `json:"connected_clients"`
}

// metricsService implements MetricsService
type metricsService struct {
	redisClient *redis.Client
	startTime   time.Time
	// Cache for performance optimization
	cache      *MetricsResponse
	cacheMutex sync.RWMutex
	lastUpdate time.Time
	cacheTTL   time.Duration
}

// NewMetricsService creates a new metrics service
func NewMetricsService(redisClient *redis.Client) MetricsService {
	return &metricsService{
		redisClient: redisClient,
		startTime:   time.Now(),
		cacheTTL:    5 * time.Second, // Cache metrics for 5 seconds
	}
}

// GetMetrics collects and returns system metrics with caching
func (s *metricsService) GetMetrics(ctx context.Context) (*MetricsResponse, error) {
	// Check cache first
	s.cacheMutex.RLock()
	if s.cache != nil && time.Since(s.lastUpdate) < s.cacheTTL {
		cached := *s.cache
		s.cacheMutex.RUnlock()
		return &cached, nil
	}
	s.cacheMutex.RUnlock()

	// Cache miss or expired, collect fresh metrics
	now := time.Now()

	// System metrics
	systemMetrics, err := s.getSystemMetrics()
	if err != nil {
		return nil, err
	}

	// Token cache metrics (placeholder - would need actual token cache implementation)
	tokenCacheMetrics := TokenCacheMetrics{
		CacheHits:    0, // TODO: Get from actual token cache
		CacheMisses:  0, // TODO: Get from actual token cache
		ActiveTokens: 0, // TODO: Get from actual token cache
		RefreshCount: 0, // TODO: Get from actual token cache
		ErrorCount:   0, // TODO: Get from actual token cache
	}

	// Actor pool metrics (placeholder - would need actual actor pool implementation)
	actorPoolMetrics := ActorPoolMetrics{
		ActiveActors:   0,  // TODO: Get from actual actor pool
		MaxActors:      10, // TODO: Get from actual actor pool config
		QueueSize:      0,  // TODO: Get from actual actor pool
		ProcessedCount: 0,  // TODO: Get from actual actor pool
		ErrorCount:     0,  // TODO: Get from actual actor pool
	}

	// Database metrics (enhanced with more details)
	databaseMetrics := DatabaseMetrics{
		OpenConnections:  0,       // TODO: Get from actual database pool
		IdleConnections:  0,       // TODO: Get from actual database pool
		MaxConnections:   100,     // TODO: Get from actual database pool config
		QueryCount:       0,       // TODO: Get from actual database stats
		QueryTimeAvg:     0.0,     // TODO: Calculate average query time
		SlowQueries:      0,       // TODO: Count slow queries (>1s)
		ConnectionErrors: 0,       // TODO: Track connection errors
		LastQueryTime:    "never", // TODO: Track last query timestamp
	}

	// Redis metrics
	redisMetrics, err := s.getRedisMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Create response
	response := &MetricsResponse{
		Timestamp:  now,
		System:     systemMetrics,
		TokenCache: tokenCacheMetrics,
		ActorPool:  actorPoolMetrics,
		Database:   databaseMetrics,
		Redis:      redisMetrics,
	}

	// Update cache
	s.cacheMutex.Lock()
	s.cache = response
	s.lastUpdate = now
	s.cacheMutex.Unlock()

	return response, nil
}

// getSystemMetrics collects system-level metrics
func (s *metricsService) getSystemMetrics() (SystemMetrics, error) {
	// Calculate uptime
	uptime := time.Since(s.startTime)

	// Get memory stats from runtime
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Convert bytes to MB for better readability
	memoryUsageMB := memStats.Alloc / 1024 / 1024

	return SystemMetrics{
		Uptime:      uptime.String(),
		MemoryUsage: uint64(memoryUsageMB),
		CPUUsage:    0, // TODO: Get from system monitoring (requires gopsutil)
		Goroutines:  runtime.NumGoroutine(),
	}, nil
}

// getRedisMetrics collects Redis metrics
func (s *metricsService) getRedisMetrics(ctx context.Context) (RedisMetrics, error) {
	// Test Redis connection
	_, err := s.redisClient.Ping(ctx).Result()
	connected := err == nil

	// Initialize default values
	keyCount := int64(0)
	memoryUsage := "unknown"
	hitRate := 0.0
	missRate := 0.0
	evictedKeys := int64(0)
	expiredKeys := int64(0)
	connectedClients := int64(0)

	if connected {
		// Get basic Redis metrics
		keyCount, _ = s.redisClient.DBSize(ctx).Result()
		memoryUsage = "available"

		// Get Redis INFO for more detailed metrics
		_, err = s.redisClient.Info(ctx, "stats").Result()
		if err == nil {
			// Parse basic stats from Redis INFO
			// This is a simplified version - in production you'd want to parse the full INFO response
			connectedClients = 1 // Default for single connection
		}

		// Calculate hit/miss rates (simplified)
		hitRate = 0.95  // TODO: Calculate from actual Redis stats
		missRate = 0.05 // TODO: Calculate from actual Redis stats
	}

	return RedisMetrics{
		Connected:        connected,
		MemoryUsage:      memoryUsage,
		KeyCount:         keyCount,
		OperationsPerSec: 0, // TODO: Calculate from actual Redis stats
		HitRate:          hitRate,
		MissRate:         missRate,
		EvictedKeys:      evictedKeys,
		ExpiredKeys:      expiredKeys,
		ConnectedClients: connectedClients,
	}, nil
}
