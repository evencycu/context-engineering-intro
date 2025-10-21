package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents application configuration
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Teams    TeamsConfig
	Server   ServerConfig
	Actor    ActorConfig
	Auth     AuthConfig
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	URL          string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	Timeout      time.Duration
}

// TeamsConfig represents Teams Bot configuration
type TeamsConfig struct {
	AppID       string
	AppPassword string
	TenantID    string
	BaseURL     string
	Timeout     time.Duration
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	RateLimit    int
	Burst        int
}

// ActorConfig represents Actor Pool configuration
type ActorConfig struct {
	MaxActors      int
	PollInterval   time.Duration
	RetryDelay     time.Duration
	MaxRetries     int
	CircuitBreaker CircuitBreakerConfig
}

// CircuitBreakerConfig represents circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold int
	RecoveryTimeout  time.Duration
	HalfOpenMaxCalls int
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	JWTSecret   string
	TokenExpiry time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "1h"),
			ConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", "30m"),
		},
		Redis: RedisConfig{
			URL:          getEnv("REDIS_URL", "redis://localhost:6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			Timeout:      getEnvAsDuration("REDIS_TIMEOUT", "5s"),
		},
		Teams: TeamsConfig{
			AppID:       getEnv("TEAMS_BOT_APP_ID", ""),
			AppPassword: getEnv("TEAMS_BOT_APP_PASSWORD", ""),
			TenantID:    getEnv("TEAMS_TENANT_ID", ""),
			BaseURL:     getEnv("TEAMS_BASE_URL", "https://smba.trafficmanager.net/apis"),
			Timeout:     getEnvAsDuration("TEAMS_TIMEOUT", "30s"),
		},
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", "30s"),
			WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", "30s"),
			RateLimit:    getEnvAsInt("RATE_LIMIT_GLOBAL", 100),
			Burst:        getEnvAsInt("RATE_LIMIT_BURST", 200),
		},
		Actor: ActorConfig{
			MaxActors:    getEnvAsInt("ACTOR_POOL_SIZE", 10),
			PollInterval: getEnvAsDuration("QUEUE_CHECK_INTERVAL", "1s"),
			RetryDelay:   getEnvAsDuration("RETRY_DELAY", "5s"),
			MaxRetries:   getEnvAsInt("MAX_RETRIES", 3),
			CircuitBreaker: CircuitBreakerConfig{
				FailureThreshold: getEnvAsInt("CIRCUIT_BREAKER_THRESHOLD", 5),
				RecoveryTimeout:  getEnvAsDuration("CIRCUIT_BREAKER_RECOVERY_TIMEOUT", "30s"),
				HalfOpenMaxCalls: getEnvAsInt("CIRCUIT_BREAKER_HALF_OPEN_MAX_CALLS", 3),
			},
		},
		Auth: AuthConfig{
			JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
			TokenExpiry: getEnvAsDuration("TOKEN_EXPIRY", "24h"),
		},
	}
}

// ValidateConfig validates the configuration
func (c *Config) ValidateConfig() error {
	if c.Teams.AppID == "" {
		return fmt.Errorf("TEAMS_BOT_APP_ID is required")
	}
	if c.Teams.AppPassword == "" {
		return fmt.Errorf("TEAMS_BOT_APP_PASSWORD is required")
	}
	if c.Teams.TenantID == "" {
		return fmt.Errorf("TEAMS_TENANT_ID is required")
	}
	return nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets environment variable as duration
func getEnvAsDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Hour // fallback
}
