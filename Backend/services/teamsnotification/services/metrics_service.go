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
	IncrementNotificationSent(ctx context.Context) error
	IncrementNotificationFailed(ctx context.Context) error
	IncrementNotificationPending(ctx context.Context) error
	IncrementTeamsAPICall(ctx context.Context) error
	IncrementTeamsAPIError(ctx context.Context) error
	UpdateResponseTime(ctx context.Context, responseTimeMs float64) error
	UpdateQueueProcessingRate(ctx context.Context, rate float64) error
	UpdateActiveProjects(ctx context.Context, count int64) error
	UpdateActiveDestinations(ctx context.Context, count int64) error
}

// MetricsResponse represents the metrics response
type MetricsResponse struct {
	Timestamp  time.Time         `json:"timestamp"`
	System     SystemMetrics     `json:"system"`
	TokenCache TokenCacheMetrics `json:"token_cache"`
	WorkerPool WorkerPoolMetrics `json:"worker_pool"`
	Database   DatabaseMetrics   `json:"database"`
	Redis      RedisMetrics      `json:"redis"`
	Business   BusinessMetrics   `json:"business"`
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

// WorkerPoolMetrics represents worker pool metrics
type WorkerPoolMetrics struct {
	ActiveWorkers  int64 `json:"active_workers"`
	MaxWorkers     int64 `json:"max_workers"`
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

// BusinessMetrics represents business-level metrics
type BusinessMetrics struct {
	NotificationsSent    int64   `json:"notifications_sent"`
	NotificationsFailed  int64   `json:"notifications_failed"`
	NotificationsPending int64   `json:"notifications_pending"`
	SuccessRate          float64 `json:"success_rate"`
	AverageResponseTime  float64 `json:"average_response_time_ms"`
	TeamsAPICalls        int64   `json:"teams_api_calls"`
	TeamsAPIErrors       int64   `json:"teams_api_errors"`
	QueueProcessingRate  float64 `json:"queue_processing_rate_per_minute"`
	ActiveProjects       int64   `json:"active_projects"`
	ActiveDestinations   int64   `json:"active_destinations"`
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

	// Worker pool metrics (placeholder - would need actual worker pool implementation)
	workerPoolMetrics := WorkerPoolMetrics{
		ActiveWorkers:  0,  // TODO: Get from actual worker pool
		MaxWorkers:     10, // TODO: Get from actual worker pool config
		QueueSize:      0,  // TODO: Get from actual worker pool
		ProcessedCount: 0,  // TODO: Get from actual worker pool
		ErrorCount:     0,  // TODO: Get from actual worker pool
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

	// Business metrics
	businessMetrics, err := s.getBusinessMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Create response
	response := &MetricsResponse{
		Timestamp:  now,
		System:     systemMetrics,
		TokenCache: tokenCacheMetrics,
		WorkerPool: workerPoolMetrics,
		Database:   databaseMetrics,
		Redis:      redisMetrics,
		Business:   businessMetrics,
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

// getBusinessMetrics collects business-level metrics
func (s *metricsService) getBusinessMetrics(ctx context.Context) (BusinessMetrics, error) {
	// Get metrics from Redis cache
	notificationsSent, _ := s.redisClient.Get(ctx, "metrics:notifications_sent").Int64()
	notificationsFailed, _ := s.redisClient.Get(ctx, "metrics:notifications_failed").Int64()
	notificationsPending, _ := s.redisClient.Get(ctx, "metrics:notifications_pending").Int64()
	teamsAPICalls, _ := s.redisClient.Get(ctx, "metrics:teams_api_calls").Int64()
	teamsAPIErrors, _ := s.redisClient.Get(ctx, "metrics:teams_api_errors").Int64()
	activeProjects, _ := s.redisClient.Get(ctx, "metrics:active_projects").Int64()
	activeDestinations, _ := s.redisClient.Get(ctx, "metrics:active_destinations").Int64()
	
	// Calculate success rate
	var successRate float64
	totalNotifications := notificationsSent + notificationsFailed
	if totalNotifications > 0 {
		successRate = float64(notificationsSent) / float64(totalNotifications) * 100
	}
	
	// Get average response time from Redis
	avgResponseTime, _ := s.redisClient.Get(ctx, "metrics:avg_response_time_ms").Float64()
	
	// Get queue processing rate (notifications per minute)
	queueProcessingRate, _ := s.redisClient.Get(ctx, "metrics:queue_processing_rate").Float64()
	
	return BusinessMetrics{
		NotificationsSent:     notificationsSent,
		NotificationsFailed:   notificationsFailed,
		NotificationsPending:  notificationsPending,
		SuccessRate:           successRate,
		AverageResponseTime:   avgResponseTime,
		TeamsAPICalls:         teamsAPICalls,
		TeamsAPIErrors:        teamsAPIErrors,
		QueueProcessingRate:   queueProcessingRate,
		ActiveProjects:        activeProjects,
		ActiveDestinations:    activeDestinations,
	}, nil
}

// IncrementNotificationSent increments the notification sent counter
func (s *metricsService) IncrementNotificationSent(ctx context.Context) error {
	return s.redisClient.Incr(ctx, "metrics:notifications_sent").Err()
}

// IncrementNotificationFailed increments the notification failed counter
func (s *metricsService) IncrementNotificationFailed(ctx context.Context) error {
	return s.redisClient.Incr(ctx, "metrics:notifications_failed").Err()
}

// IncrementNotificationPending increments the notification pending counter
func (s *metricsService) IncrementNotificationPending(ctx context.Context) error {
	return s.redisClient.Incr(ctx, "metrics:notifications_pending").Err()
}

// IncrementTeamsAPICall increments the Teams API call counter
func (s *metricsService) IncrementTeamsAPICall(ctx context.Context) error {
	return s.redisClient.Incr(ctx, "metrics:teams_api_calls").Err()
}

// IncrementTeamsAPIError increments the Teams API error counter
func (s *metricsService) IncrementTeamsAPIError(ctx context.Context) error {
	return s.redisClient.Incr(ctx, "metrics:teams_api_errors").Err()
}

// UpdateResponseTime updates the average response time
func (s *metricsService) UpdateResponseTime(ctx context.Context, responseTimeMs float64) error {
	// Update running average
	pipe := s.redisClient.Pipeline()
	pipe.Get(ctx, "metrics:avg_response_time_ms")
	pipe.Get(ctx, "metrics:response_time_count")
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	
	oldAvg, _ := cmds[0].(*redis.StringCmd).Float64()
	count, _ := cmds[1].(*redis.StringCmd).Int64()
	
	newCount := count + 1
	newAvg := (oldAvg*float64(count) + responseTimeMs) / float64(newCount)
	
	pipe = s.redisClient.Pipeline()
	pipe.Set(ctx, "metrics:avg_response_time_ms", newAvg, 0)
	pipe.Set(ctx, "metrics:response_time_count", newCount, 0)
	_, err = pipe.Exec(ctx)
	return err
}

// UpdateQueueProcessingRate updates the queue processing rate
func (s *metricsService) UpdateQueueProcessingRate(ctx context.Context, rate float64) error {
	return s.redisClient.Set(ctx, "metrics:queue_processing_rate", rate, 0).Err()
}

// UpdateActiveProjects updates the active projects count
func (s *metricsService) UpdateActiveProjects(ctx context.Context, count int64) error {
	return s.redisClient.Set(ctx, "metrics:active_projects", count, 0).Err()
}

// UpdateActiveDestinations updates the active destinations count
func (s *metricsService) UpdateActiveDestinations(ctx context.Context, count int64) error {
	return s.redisClient.Set(ctx, "metrics:active_destinations", count, 0).Err()
}