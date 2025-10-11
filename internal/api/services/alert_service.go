package services

import (
	"context"
	"time"
)

// AlertService provides alerting capabilities
type AlertService interface {
	CheckAlerts(ctx context.Context, metrics *MetricsResponse) ([]Alert, error)
	GetAlertHistory(ctx context.Context) ([]Alert, error)
}

// Alert represents an alert condition
type Alert struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Severity   string                 `json:"severity"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Resolved   bool                   `json:"resolved"`
	ResolvedAt *time.Time             `json:"resolved_at,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// AlertThresholds defines alert thresholds
type AlertThresholds struct {
	MemoryUsageMB    uint64  `json:"memory_usage_mb"`
	CPUUsagePercent  float64 `json:"cpu_usage_percent"`
	ResponseTimeMS   float64 `json:"response_time_ms"`
	ErrorRatePercent float64 `json:"error_rate_percent"`
	RedisHitRate     float64 `json:"redis_hit_rate"`
	DatabaseConnMax  int64   `json:"database_conn_max"`
}

// alertService implements AlertService
type alertService struct {
	thresholds AlertThresholds
	alerts     []Alert
}

// NewAlertService creates a new alert service
func NewAlertService() AlertService {
	return &alertService{
		thresholds: AlertThresholds{
			MemoryUsageMB:    500,  // Alert if memory usage > 500MB
			CPUUsagePercent:  80.0, // Alert if CPU usage > 80%
			ResponseTimeMS:   1000, // Alert if response time > 1000ms
			ErrorRatePercent: 5.0,  // Alert if error rate > 5%
			RedisHitRate:     0.8,  // Alert if Redis hit rate < 80%
			DatabaseConnMax:  80,   // Alert if DB connections > 80
		},
		alerts: make([]Alert, 0),
	}
}

// CheckAlerts checks metrics against thresholds and generates alerts
func (s *alertService) CheckAlerts(ctx context.Context, metrics *MetricsResponse) ([]Alert, error) {
	var newAlerts []Alert
	now := time.Now()

	// Check memory usage
	if metrics.System.MemoryUsage > s.thresholds.MemoryUsageMB {
		alert := Alert{
			ID:        "memory_high",
			Type:      "system",
			Severity:  "warning",
			Message:   "High memory usage detected",
			Timestamp: now,
			Metadata: map[string]interface{}{
				"memory_usage_mb": metrics.System.MemoryUsage,
				"threshold_mb":    s.thresholds.MemoryUsageMB,
			},
		}
		newAlerts = append(newAlerts, alert)
	}

	// Check CPU usage
	if metrics.System.CPUUsage > s.thresholds.CPUUsagePercent {
		alert := Alert{
			ID:        "cpu_high",
			Type:      "system",
			Severity:  "critical",
			Message:   "High CPU usage detected",
			Timestamp: now,
			Metadata: map[string]interface{}{
				"cpu_usage_percent": metrics.System.CPUUsage,
				"threshold_percent": s.thresholds.CPUUsagePercent,
			},
		}
		newAlerts = append(newAlerts, alert)
	}

	// Check Redis hit rate
	if metrics.Redis.HitRate < s.thresholds.RedisHitRate {
		alert := Alert{
			ID:        "redis_hit_rate_low",
			Type:      "redis",
			Severity:  "warning",
			Message:   "Low Redis hit rate detected",
			Timestamp: now,
			Metadata: map[string]interface{}{
				"hit_rate":       metrics.Redis.HitRate,
				"threshold_rate": s.thresholds.RedisHitRate,
			},
		}
		newAlerts = append(newAlerts, alert)
	}

	// Check database connections
	if metrics.Database.OpenConnections > s.thresholds.DatabaseConnMax {
		alert := Alert{
			ID:        "db_connections_high",
			Type:      "database",
			Severity:  "warning",
			Message:   "High database connection count detected",
			Timestamp: now,
			Metadata: map[string]interface{}{
				"open_connections": metrics.Database.OpenConnections,
				"threshold_conn":   s.thresholds.DatabaseConnMax,
			},
		}
		newAlerts = append(newAlerts, alert)
	}

	// Check Redis connection
	if !metrics.Redis.Connected {
		alert := Alert{
			ID:        "redis_disconnected",
			Type:      "redis",
			Severity:  "critical",
			Message:   "Redis connection lost",
			Timestamp: now,
			Metadata: map[string]interface{}{
				"connected": metrics.Redis.Connected,
			},
		}
		newAlerts = append(newAlerts, alert)
	}

	// Store new alerts
	s.alerts = append(s.alerts, newAlerts...)

	return newAlerts, nil
}

// GetAlertHistory returns alert history
func (s *alertService) GetAlertHistory(ctx context.Context) ([]Alert, error) {
	return s.alerts, nil
}
