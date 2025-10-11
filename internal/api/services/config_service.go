package services

import (
	"context"
	"os"
	"strconv"
	"strings"
)

// ConfigService provides system configuration management
type ConfigService interface {
	GetConfig(ctx context.Context) (*ConfigResponse, error)
	ValidateConfig(ctx context.Context) (*ConfigValidationResponse, error)
}

// ConfigResponse represents the configuration response
type ConfigResponse struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Teams    TeamsConfig    `json:"teams"`
	Actor    ActorConfig    `json:"actor"`
	Logger   LoggerConfig   `json:"logger"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string `json:"port"`
	Environment  string `json:"environment"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
	RateLimit    int    `json:"rate_limit"`
	Burst        int    `json:"burst"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime string `json:"conn_max_lifetime"`
	ConnMaxIdleTime string `json:"conn_max_idle_time"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
	MaxRetries   int    `json:"max_retries"`
	Timeout      string `json:"timeout"`
}

// TeamsConfig represents Teams configuration
type TeamsConfig struct {
	AppID    string `json:"app_id"`
	TenantID string `json:"tenant_id"`
	BaseURL  string `json:"base_url"`
	Timeout  string `json:"timeout"`
}

// ActorConfig represents actor configuration
type ActorConfig struct {
	MaxActors      int                  `json:"max_actors"`
	PollInterval   string               `json:"poll_interval"`
	RetryDelay     string               `json:"retry_delay"`
	MaxRetries     int                  `json:"max_retries"`
	CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker"`
}

// CircuitBreakerConfig represents circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold int    `json:"failure_threshold"`
	RecoveryTimeout  string `json:"recovery_timeout"`
	HalfOpenMaxCalls int    `json:"half_open_max_calls"`
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// ConfigValidationResponse represents configuration validation response
type ConfigValidationResponse struct {
	Valid    bool            `json:"valid"`
	Errors   []ConfigError   `json:"errors"`
	Warnings []ConfigWarning `json:"warnings"`
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ConfigWarning represents a configuration warning
type ConfigWarning struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

// configService implements ConfigService
type configService struct{}

// NewConfigService creates a new config service
func NewConfigService() ConfigService {
	return &configService{}
}

// GetConfig returns current system configuration
func (s *configService) GetConfig(ctx context.Context) (*ConfigResponse, error) {
	return &ConfigResponse{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			ReadTimeout:  getEnv("READ_TIMEOUT", "30s"),
			WriteTimeout: getEnv("WRITE_TIMEOUT", "30s"),
			RateLimit:    getEnvInt("RATE_LIMIT", 100),
			Burst:        getEnvInt("BURST", 200),
		},
		Database: DatabaseConfig{
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnv("DB_CONN_MAX_LIFETIME", "5m"),
			ConnMaxIdleTime: getEnv("DB_CONN_MAX_IDLE_TIME", "1m"),
		},
		Redis: RedisConfig{
			PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 2),
			MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
			Timeout:      getEnv("REDIS_TIMEOUT", "5s"),
		},
		Teams: TeamsConfig{
			AppID:    getEnv("TEAMS_BOT_APP_ID", ""),
			TenantID: getEnv("TEAMS_TENANT_ID", ""),
			BaseURL:  getEnv("TEAMS_BASE_URL", "https://smba.trafficmanager.net/apac"),
			Timeout:  getEnv("TEAMS_TIMEOUT", "30s"),
		},
		Actor: ActorConfig{
			MaxActors:    getEnvInt("ACTOR_MAX_ACTORS", 10),
			PollInterval: getEnv("ACTOR_POLL_INTERVAL", "1s"),
			RetryDelay:   getEnv("ACTOR_RETRY_DELAY", "5s"),
			MaxRetries:   getEnvInt("ACTOR_MAX_RETRIES", 3),
			CircuitBreaker: CircuitBreakerConfig{
				FailureThreshold: getEnvInt("CIRCUIT_BREAKER_FAILURE_THRESHOLD", 5),
				RecoveryTimeout:  getEnv("CIRCUIT_BREAKER_RECOVERY_TIMEOUT", "30s"),
				HalfOpenMaxCalls: getEnvInt("CIRCUIT_BREAKER_HALF_OPEN_MAX_CALLS", 3),
			},
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}, nil
}

// ValidateConfig validates system configuration
func (s *configService) ValidateConfig(ctx context.Context) (*ConfigValidationResponse, error) {
	var errors []ConfigError
	var warnings []ConfigWarning

	// Required environment variables
	requiredVars := []string{
		"DATABASE_URL",
		"REDIS_URL",
		"TEAMS_BOT_APP_ID",
		"TEAMS_TENANT_ID",
		"TEAMS_BOT_APP_PASSWORD",
	}

	for _, varName := range requiredVars {
		if value := os.Getenv(varName); value == "" {
			errors = append(errors, ConfigError{
				Field:   varName,
				Message: "missing required environment variable",
				Code:    "MISSING_REQUIRED",
			})
		}
	}

	// Validate specific configurations
	if port := getEnv("PORT", "8080"); port != "" {
		if _, err := strconv.Atoi(port); err != nil {
			errors = append(errors, ConfigError{
				Field:   "PORT",
				Message: "invalid port number",
				Code:    "INVALID_PORT",
			})
		}
	}

	// Check for potential issues
	if env := getEnv("ENVIRONMENT", "development"); env == "production" {
		if logLevel := getEnv("LOG_LEVEL", "info"); logLevel == "debug" {
			warnings = append(warnings, ConfigWarning{
				Field:      "LOG_LEVEL",
				Message:    "debug logging in production",
				Suggestion: "consider using 'info' or 'warn' level in production",
			})
		}
	}

	// Check Redis configuration
	if redisURL := getEnv("REDIS_URL", ""); redisURL != "" {
		if !strings.HasPrefix(redisURL, "redis://") && !strings.HasPrefix(redisURL, "rediss://") {
			warnings = append(warnings, ConfigWarning{
				Field:      "REDIS_URL",
				Message:    "Redis URL should use redis:// or rediss:// protocol",
				Suggestion: "use redis:// for unencrypted or rediss:// for encrypted connections",
			})
		}
	}

	// Check database configuration
	if dbURL := getEnv("DATABASE_URL", ""); dbURL != "" {
		if !strings.HasPrefix(dbURL, "postgres://") && !strings.HasPrefix(dbURL, "postgresql://") {
			warnings = append(warnings, ConfigWarning{
				Field:      "DATABASE_URL",
				Message:    "Database URL should use postgres:// or postgresql:// protocol",
				Suggestion: "use postgresql:// for PostgreSQL connections",
			})
		}
	}

	return &ConfigValidationResponse{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// Helper functions
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
