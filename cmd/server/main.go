package main

import (
	"context"
	"os"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/actor"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/batch"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/billing"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/bots"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/companies"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/destinations"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/external"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/files"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/messages"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/notifications"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/projects"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/provision"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/users"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/middleware"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/repositories"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/storage"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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
	dbURL := getEnv("DATABASE_URL", "postgres://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable")
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

	// Initialize Redis connection
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.Fatal("Failed to parse Redis URL: ", err)
	}
	redisClient := redis.NewClient(opt)
	defer redisClient.Close()

	// Test Redis connection
	ctx := context.Background()
	if err = redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis: ", err)
	}
	logger.Info("Successfully connected to Redis!")

	// Initialize repositories and services
	companyRepo := repositories.NewCompanyRepository(db)
	companyService := services.NewCompanyService(companyRepo)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)

	teamsBotRepo := repositories.NewTeamsBotRepository(db)
	teamsBotService := services.NewteamsBotService(teamsBotRepo)
	installationRepo := repositories.NewBotInstallationRepository(db)
	messagesService := services.NewMessagesService(teamsBotRepo, installationRepo, getEnv("TEAMS_TENANT_ID", ""))

	destinationRepo := repositories.NewDestinationRepository(db)
	destinationService := services.NewDestinationService(destinationRepo)

	// Provision service (company + project + destination)
	provisionService := services.NewProvisionService(companyRepo, projectRepo, destinationRepo)

	notificationRepo := repositories.NewNotificationRepository(db)
	notificationDestRepo := repositories.NewNotificationDestinationRepository(db)
	broadcaster := services.NewBroadcastService(teamsBotRepo, installationRepo, destinationRepo)

	// Initialize Redis Queue
	redisQueue := actor.NewRedisQueue(redisClient)
	notificationService := services.NewNotificationService(notificationRepo, destinationRepo, notificationDestRepo, broadcaster, redisQueue)

	// Initialize Actor Pool for async notification sending
	actorDB := actor.NewActorDB(notificationDestRepo, installationRepo, teamsBotRepo, notificationRepo)
	actorPool := actor.NewActorPool(redisClient, actorDB, 10, 10*time.Second) // Max 10 actors, 10s poll interval
	actorPool.Start(ctx)
	defer actorPool.Stop()
	logger.Info("Actor Pool started successfully")

	// Start Redis Queue consumer (high -> normal -> low)
	consumer := actor.NewQueueConsumer(redisQueue, actorPool)
	consumer.Start(ctx)
	defer consumer.Stop()
	logger.Info("Queue Consumer started successfully")

	// External service
	externalService := services.NewExternalService(projectRepo, destinationRepo, notificationService)

	// Billing service
	billingPlanRepo := repositories.NewBillingPlanRepository(db)
	usageRecordRepo := repositories.NewUsageRecordRepository(db)
	companyBillingRepo := repositories.NewCompanyBillingRepository(db)
	billingService := services.NewBillingService(usageRecordRepo, billingPlanRepo, companyRepo, projectRepo, userRepo)
	_ = companyBillingRepo // Will be used when implementing company billing features

	// File service
	localStorage := storage.NewLocalStorage("uploads") // Assuming "uploads" directory
	fileRepo := repositories.NewFileRepository(db)
	fileService := services.NewFileService(fileRepo, localStorage)

	// Batch service
	batchRepo := repositories.NewBatchRepository(db)
	batchService := services.NewBatchService(batchRepo, notificationService, projectRepo, destinationRepo)

	// Queue system replaced by Actor Pool V2

	// Create handlers
	companyHandler := companies.NewHandler(companyService)
	userHandler := users.NewHandler(userService)
	projectHandler := projects.NewHandler(projectService)
	botHandler := bots.NewHandler(teamsBotService)
	messagesHandler := messages.NewHandler(messagesService)
	destinationHandler := destinations.NewHandler(destinationService)
	notificationHandler := notifications.NewHandler(notificationService)
	provisionHandler := provision.NewHandler(provisionService)
	externalHandler := external.NewHandler(externalService)
	billingHandler := billing.NewHandler(billingService)
	fileHandler := files.NewHandler(fileService)
	batchHandler := batch.NewHandler(batchService)
	// Queue API handler removed - replaced by Actor Pool V2

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
	server.RegisterHandlers(companyHandler, userHandler, projectHandler, botHandler, destinationHandler, notificationHandler, messagesHandler, provisionHandler, externalHandler, billingHandler, fileHandler, batchHandler)

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
