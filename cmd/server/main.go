package main

import (
	"os"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/bots"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/companies"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/destinations"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/messages"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/notifications"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/projects"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/users"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/middleware"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/repositories"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg := loadConfig()

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	errorHandler := middleware.NewErrorHandler(logger)

	// Initialize database connection
	dbURL := getEnv("DATABASE_URL", "postgres://teamsnotify:teamsnotify123@localhost:5432/teamsnotify?sslmode=disable")
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	// Test database connection
	if err = db.Ping(); err != nil {
		logger.Fatal("Failed to ping database: ", err)
	}
	logger.Info("Successfully connected to database!")

	// Initialize repositories and services
	companyRepo := repositories.NewCompanyRepository(db)
	companyService := services.NewCompanyService(companyRepo)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)

	platformBotRepo := repositories.NewPlatformBotRepository(db)
	platformBotService := services.NewPlatformBotService(platformBotRepo)
	installationRepo := repositories.NewBotInstallationRepository(db)
	messagesService := services.NewMessagesService(platformBotRepo, installationRepo, getEnv("TEAMS_TENANT_ID", ""))
	thirdPartyBotRepo := repositories.NewThirdPartyBotRepository(db)
	thirdPartyBotService := services.NewThirdPartyBotService(thirdPartyBotRepo)

	destinationRepo := repositories.NewDestinationRepository(db)
	destinationService := services.NewDestinationService(destinationRepo)

	notificationRepo := repositories.NewNotificationRepository(db)
	broadcaster := services.NewBroadcastService(platformBotRepo, thirdPartyBotRepo, installationRepo, destinationRepo)
	notificationService := services.NewNotificationService(notificationRepo, broadcaster)

	// Create handlers
	companyHandler := companies.NewHandler(companyService)
	userHandler := users.NewHandler(userService)
	projectHandler := projects.NewHandler(projectService)
	botHandler := bots.NewHandler(platformBotService, thirdPartyBotService)
	messagesHandler := messages.NewHandler(messagesService)
	destinationHandler := destinations.NewHandler(destinationService)
	notificationHandler := notifications.NewHandler(notificationService)

	// Create server
	server := api.NewServer(api.Config{
		Port:         cfg.Port,
		Environment:  "development",
		RateLimit:    rate.Limit(100), // 100 requests per second
		Burst:        200,             // Allow bursts up to 200 requests
		ReadTimeout:  time.Duration(30) * time.Second,
		WriteTimeout: time.Duration(30) * time.Second,
	})

	// Get router and add middleware
	router := server.GetRouter()

	// Add global middleware
	router.Use(loggingMiddleware.StartTimeMiddleware())
	router.Use(loggingMiddleware.RequestLogger())
	router.Use(loggingMiddleware.RequestBodyLogger())
	router.Use(loggingMiddleware.ResponseLogger())
	router.Use(loggingMiddleware.ErrorLogger())
	router.Use(loggingMiddleware.SecurityHeaders())
	router.Use(loggingMiddleware.RateLimitLogger())
	router.Use(errorHandler.ErrorHandlerMiddleware())
	router.Use(errorHandler.PanicRecovery())

	// Add 404 and 405 handlers
	router.NoRoute(errorHandler.NotFoundHandler())
	router.NoMethod(errorHandler.MethodNotAllowedHandler())

	// Add auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			// TODO: Implement login
			c.JSON(200, gin.H{"message": "Login endpoint - TODO"})
		})
		auth.POST("/refresh", authMiddleware.RefreshToken)
	}

	// Register API routes
	server.RegisterRoutes(companyHandler, userHandler, projectHandler, botHandler, destinationHandler, notificationHandler, messagesHandler)

	// Start server
	logger.Info("Starting server on port " + cfg.Port)
	if err := server.Start(); err != nil {
		logger.Fatal("Failed to start server: ", err)
	}
}

// Config represents server configuration
type Config struct {
	Port      string
	JWTSecret string
	RateLimit rate.Limit
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		RateLimit: rate.Limit(100), // 100 requests per second
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
