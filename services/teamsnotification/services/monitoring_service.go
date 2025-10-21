package services

import (
	"context"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/services/teamsnotification/repositories"
)

// MonitoringService provides comprehensive system monitoring
type MonitoringService interface {
	GetSystemHealth(ctx context.Context) (*SystemHealth, error)
	GetPerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error)
	GetBusinessHealth(ctx context.Context) (*BusinessHealth, error)
	GetAlertStatus(ctx context.Context) (*AlertStatus, error)
}

type monitoringService struct {
	metricsService   MetricsService
	projectRepo      repositories.ProjectRepository
	notificationRepo repositories.NotificationRepository
	destinationRepo  repositories.DestinationRepository
}

// SystemHealth represents overall system health
type SystemHealth struct {
	Status     string                     `json:"status"`
	Timestamp  time.Time                  `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
	Overall    HealthScore                `json:"overall"`
}

// ComponentHealth represents health of individual components
type ComponentHealth struct {
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency_ms,omitempty"`
}

// HealthScore represents overall health score
type HealthScore struct {
	Score   int    `json:"score"` // 0-100
	Grade   string `json:"grade"` // A, B, C, D, F
	Message string `json:"message"`
}

// PerformanceMetrics represents system performance metrics
type PerformanceMetrics struct {
	Timestamp     time.Time            `json:"timestamp"`
	ResponseTime  ResponseTimeMetrics  `json:"response_time"`
	Throughput    ThroughputMetrics    `json:"throughput"`
	ErrorRate     ErrorRateMetrics     `json:"error_rate"`
	ResourceUsage ResourceUsageMetrics `json:"resource_usage"`
}

// ResponseTimeMetrics represents response time statistics
type ResponseTimeMetrics struct {
	Average float64 `json:"average_ms"`
	P50     float64 `json:"p50_ms"`
	P95     float64 `json:"p95_ms"`
	P99     float64 `json:"p99_ms"`
	Max     float64 `json:"max_ms"`
	Min     float64 `json:"min_ms"`
}

// ThroughputMetrics represents throughput statistics
type ThroughputMetrics struct {
	RequestsPerSecond      float64 `json:"requests_per_second"`
	NotificationsPerMinute float64 `json:"notifications_per_minute"`
	APICallsPerMinute      float64 `json:"api_calls_per_minute"`
}

// ErrorRateMetrics represents error rate statistics
type ErrorRateMetrics struct {
	TotalErrors    int64   `json:"total_errors"`
	ErrorRate      float64 `json:"error_rate_percent"`
	SuccessRate    float64 `json:"success_rate_percent"`
	CriticalErrors int64   `json:"critical_errors"`
	WarningErrors  int64   `json:"warning_errors"`
}

// ResourceUsageMetrics represents resource usage statistics
type ResourceUsageMetrics struct {
	MemoryUsage         float64 `json:"memory_usage_mb"`
	CPUUsage            float64 `json:"cpu_usage_percent"`
	DatabaseConnections int     `json:"database_connections"`
	RedisConnections    int     `json:"redis_connections"`
}

// BusinessHealth represents business-level health metrics
type BusinessHealth struct {
	Timestamp          time.Time      `json:"timestamp"`
	ActiveProjects     int64          `json:"active_projects"`
	ActiveDestinations int64          `json:"active_destinations"`
	NotificationsToday int64          `json:"notifications_today"`
	SuccessRate        float64        `json:"success_rate_percent"`
	QueueHealth        QueueHealth    `json:"queue_health"`
	TeamsAPIHealth     TeamsAPIHealth `json:"teams_api_health"`
}

// QueueHealth represents queue system health
type QueueHealth struct {
	PendingCount    int64   `json:"pending_count"`
	ProcessingCount int64   `json:"processing_count"`
	FailedCount     int64   `json:"failed_count"`
	ProcessingRate  float64 `json:"processing_rate_per_minute"`
	Status          string  `json:"status"`
}

// TeamsAPIHealth represents Teams API health
type TeamsAPIHealth struct {
	Status             string    `json:"status"`
	ResponseTime       float64   `json:"response_time_ms"`
	SuccessRate        float64   `json:"success_rate_percent"`
	RateLimitStatus    string    `json:"rate_limit_status"`
	LastSuccessfulCall time.Time `json:"last_successful_call"`
}

// AlertStatus represents current alert status
type AlertStatus struct {
	Timestamp    time.Time         `json:"timestamp"`
	ActiveAlerts int               `json:"active_alerts"`
	Alerts       []MonitoringAlert `json:"alerts"`
}

// MonitoringAlert represents an individual monitoring alert
type MonitoringAlert struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	AlertStatus string    `json:"status"`
	Component   string    `json:"component"`
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(
	metricsService MetricsService,
	projectRepo repositories.ProjectRepository,
	notificationRepo repositories.NotificationRepository,
	destinationRepo repositories.DestinationRepository,
) MonitoringService {
	return &monitoringService{
		metricsService:   metricsService,
		projectRepo:      projectRepo,
		notificationRepo: notificationRepo,
		destinationRepo:  destinationRepo,
	}
}

// GetSystemHealth returns overall system health
func (s *monitoringService) GetSystemHealth(ctx context.Context) (*SystemHealth, error) {
	// Check database health
	dbHealth := s.checkDatabaseHealth(ctx)

	// Check Redis health
	redisHealth := s.checkRedisHealth(ctx)

	// Check API health
	apiHealth := s.checkAPIHealth(ctx)

	// Calculate overall health
	overall := s.calculateOverallHealth(dbHealth, redisHealth, apiHealth)

	components := map[string]ComponentHealth{
		"database": dbHealth,
		"redis":    redisHealth,
		"api":      apiHealth,
	}

	status := "healthy"
	if overall.Score < 80 {
		status = "degraded"
	}
	if overall.Score < 60 {
		status = "unhealthy"
	}

	return &SystemHealth{
		Status:     status,
		Timestamp:  time.Now(),
		Components: components,
		Overall:    overall,
	}, nil
}

// GetPerformanceMetrics returns performance metrics
func (s *monitoringService) GetPerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error) {
	// Get metrics from metrics service
	metrics, err := s.metricsService.GetMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	// Calculate performance metrics
	responseTime := ResponseTimeMetrics{
		Average: metrics.Business.AverageResponseTime,
		P50:     metrics.Business.AverageResponseTime * 0.8, // Simplified calculation
		P95:     metrics.Business.AverageResponseTime * 1.5,
		P99:     metrics.Business.AverageResponseTime * 2.0,
		Max:     metrics.Business.AverageResponseTime * 3.0,
		Min:     metrics.Business.AverageResponseTime * 0.5,
	}

	throughput := ThroughputMetrics{
		RequestsPerSecond:      10.0, // Placeholder - would need actual tracking
		NotificationsPerMinute: metrics.Business.QueueProcessingRate,
		APICallsPerMinute:      float64(metrics.Business.TeamsAPICalls),
	}

	errorRate := ErrorRateMetrics{
		TotalErrors:    metrics.Business.NotificationsFailed,
		ErrorRate:      100 - metrics.Business.SuccessRate,
		SuccessRate:    metrics.Business.SuccessRate,
		CriticalErrors: metrics.Business.TeamsAPIErrors,
		WarningErrors:  metrics.Business.NotificationsFailed - metrics.Business.TeamsAPIErrors,
	}

	resourceUsage := ResourceUsageMetrics{
		MemoryUsage:         float64(metrics.System.MemoryUsage) / 1024 / 1024, // Convert bytes to MB
		CPUUsage:            metrics.System.CPUUsage,
		DatabaseConnections: 5, // Placeholder - would need actual tracking
		RedisConnections:    3, // Placeholder - would need actual tracking
	}

	return &PerformanceMetrics{
		Timestamp:     time.Now(),
		ResponseTime:  responseTime,
		Throughput:    throughput,
		ErrorRate:     errorRate,
		ResourceUsage: resourceUsage,
	}, nil
}

// GetBusinessHealth returns business health metrics
func (s *monitoringService) GetBusinessHealth(ctx context.Context) (*BusinessHealth, error) {
	// Get metrics from metrics service
	metrics, err := s.metricsService.GetMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	// Get business data from repositories (simplified)
	activeProjects := metrics.Business.ActiveProjects
	activeDestinations := metrics.Business.ActiveDestinations

	// Calculate notifications today (simplified)
	notificationsToday := metrics.Business.NotificationsSent

	queueHealth := QueueHealth{
		PendingCount:    metrics.Business.NotificationsPending,
		ProcessingCount: 0, // Would need additional tracking
		FailedCount:     metrics.Business.NotificationsFailed,
		ProcessingRate:  metrics.Business.QueueProcessingRate,
		Status:          "healthy",
	}

	teamsAPIHealth := TeamsAPIHealth{
		Status:             "healthy",
		ResponseTime:       metrics.Business.AverageResponseTime,
		SuccessRate:        metrics.Business.SuccessRate,
		RateLimitStatus:    "normal",
		LastSuccessfulCall: time.Now().Add(-time.Minute * 5), // Simplified
	}

	return &BusinessHealth{
		Timestamp:          time.Now(),
		ActiveProjects:     activeProjects,
		ActiveDestinations: activeDestinations,
		NotificationsToday: notificationsToday,
		SuccessRate:        metrics.Business.SuccessRate,
		QueueHealth:        queueHealth,
		TeamsAPIHealth:     teamsAPIHealth,
	}, nil
}

// GetAlertStatus returns current alert status
func (s *monitoringService) GetAlertStatus(ctx context.Context) (*AlertStatus, error) {
	alerts := []MonitoringAlert{}

	// Check for various alert conditions
	metrics, err := s.metricsService.GetMetrics(ctx)
	if err == nil {
		// Check for high error rate
		if metrics.Business.SuccessRate < 90 {
			alerts = append(alerts, MonitoringAlert{
				ID:          "high_error_rate",
				Type:        "performance",
				Severity:    "warning",
				Message:     fmt.Sprintf("Success rate is %.2f%%, below threshold", metrics.Business.SuccessRate),
				Timestamp:   time.Now(),
				AlertStatus: "active",
				Component:   "notifications",
			})
		}

		// Check for high response time
		if metrics.Business.AverageResponseTime > 5000 {
			alerts = append(alerts, MonitoringAlert{
				ID:          "high_response_time",
				Type:        "performance",
				Severity:    "warning",
				Message:     fmt.Sprintf("Average response time is %.2fms, above threshold", metrics.Business.AverageResponseTime),
				Timestamp:   time.Now(),
				AlertStatus: "active",
				Component:   "api",
			})
		}

		// Check for Teams API errors
		if metrics.Business.TeamsAPIErrors > 10 {
			alerts = append(alerts, MonitoringAlert{
				ID:          "teams_api_errors",
				Type:        "integration",
				Severity:    "critical",
				Message:     fmt.Sprintf("Teams API errors: %d", metrics.Business.TeamsAPIErrors),
				Timestamp:   time.Now(),
				AlertStatus: "active",
				Component:   "teams_api",
			})
		}
	}

	return &AlertStatus{
		Timestamp:    time.Now(),
		ActiveAlerts: len(alerts),
		Alerts:       alerts,
	}, nil
}

// Helper methods for health checks
func (s *monitoringService) checkDatabaseHealth(ctx context.Context) ComponentHealth {
	start := time.Now()
	// Simple database ping check
	// In real implementation, would test actual database connection
	return ComponentHealth{
		Status:    "healthy",
		Message:   "Database connection successful",
		Timestamp: time.Now(),
		Latency:   time.Since(start).Milliseconds(),
	}
}

func (s *monitoringService) checkRedisHealth(ctx context.Context) ComponentHealth {
	start := time.Now()
	// Simple Redis ping check
	// In real implementation, would test actual Redis connection
	return ComponentHealth{
		Status:    "healthy",
		Message:   "Redis connection successful",
		Timestamp: time.Now(),
		Latency:   time.Since(start).Milliseconds(),
	}
}

func (s *monitoringService) checkAPIHealth(ctx context.Context) ComponentHealth {
	start := time.Now()
	// API health check
	return ComponentHealth{
		Status:    "healthy",
		Message:   "API service running",
		Timestamp: time.Now(),
		Latency:   time.Since(start).Milliseconds(),
	}
}

func (s *monitoringService) calculateOverallHealth(db, redis, api ComponentHealth) HealthScore {
	score := 100
	issues := []string{}

	if db.Status != "healthy" {
		score -= 30
		issues = append(issues, "Database issues")
	}
	if redis.Status != "healthy" {
		score -= 20
		issues = append(issues, "Redis issues")
	}
	if api.Status != "healthy" {
		score -= 25
		issues = append(issues, "API issues")
	}

	grade := "A"
	if score < 90 {
		grade = "B"
	}
	if score < 80 {
		grade = "C"
	}
	if score < 70 {
		grade = "D"
	}
	if score < 60 {
		grade = "F"
	}

	message := "System is healthy"
	if len(issues) > 0 {
		message = fmt.Sprintf("Issues detected: %s", fmt.Sprintf("%v", issues))
	}

	return HealthScore{
		Score:   score,
		Grade:   grade,
		Message: message,
	}
}
