package app

import (
	"context"
	"time"

	configs "github.com/evencycu/TeamsNotifyGoV3/configs"
	"github.com/evencycu/TeamsNotifyGoV3/libs/middleware"
	"github.com/evencycu/TeamsNotifyGoV3/libs/storage"
	"github.com/evencycu/TeamsNotifyGoV3/libs/token"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/actor"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/billing"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/bots"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/companies"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/destinations"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/directory"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/files"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/messages"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/monitoring"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/notifications"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/notify"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/projects"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/provision"
	qhandler "github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/queue"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/system"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/templates"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/handlers/users"
	repositories "github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
	services "github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Application holds the application's dependencies
type Application struct {
	Config    *configs.Config
	Logger    *logrus.Logger
	DB        *sqlx.DB
	Redis     *redis.Client
	Server    *teamsnotification.Server
	ActorPool *actor.NotificationProcessor
}

// NewApplication initializes and wires all application dependencies
func NewApplication(cfg *configs.Config, logger *logrus.Logger) (*Application, error) {
	// Initialize database connection
	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	// Test database connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	logger.Info("Successfully connected to database!")

	// Initialize Redis connection
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		db.Close()
		return nil, err
	}
	opt.Password = cfg.Redis.Password
	opt.DB = cfg.Redis.DB
	opt.PoolSize = cfg.Redis.PoolSize
	opt.MinIdleConns = cfg.Redis.MinIdleConns
	opt.MaxRetries = cfg.Redis.MaxRetries
	opt.DialTimeout = cfg.Redis.Timeout
	opt.ReadTimeout = cfg.Redis.Timeout
	opt.WriteTimeout = cfg.Redis.Timeout

	redisClient := redis.NewClient(opt)
	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = redisClient.Ping(ctx).Err(); err != nil {
		db.Close()
		redisClient.Close()
		return nil, err
	}
	logger.Info("Successfully connected to Redis!")

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.JWTSecret)
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	errorHandler := middleware.NewErrorHandler(logger)

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
	messagesService := services.NewMessagesService(teamsBotRepo, installationRepo, cfg.Teams.TenantID)

	destinationRepo := repositories.NewDestinationRepository(db)
	destinationService := services.NewDestinationService(destinationRepo)

	provisionService := services.NewProvisionService(companyRepo, projectRepo, destinationRepo)

	notificationRepo := repositories.NewNotificationRepository(db)
	notificationDestRepo := repositories.NewNotificationDestinationRepository(db)

	templateRepo := repositories.NewTemplateRepository(db)
	templateService := services.NewTemplateService(templateRepo)

	metricsService := services.NewMetricsService(redisClient)
	broadcaster := services.NewBroadcastService(teamsBotRepo, installationRepo, destinationRepo, metricsService)
	tokenManager := token.NewTokenManager(redisClient)
	logger.Info("Token Manager initialized successfully")

	// Directory & Graph
	directoryRepo := repositories.NewDirectoryRepository(db)
	graphService := services.NewGraphService(tokenManager, cfg.Teams.TenantID, cfg.Teams.AppID)
	directoryService := services.NewDirectoryService(directoryRepo, graphService)

	billingPlanRepo := repositories.NewBillingPlanRepository(db)
	usageRecordRepo := repositories.NewUsageRecordRepository(db)
	projectBillingRepo := repositories.NewProjectBillingRepository(db)
	billingService := services.NewBillingService(usageRecordRepo, billingPlanRepo, projectRepo, userRepo)
	_ = projectBillingRepo

	redisQueue := actor.NewRedisQueue(redisClient)
	notificationService := services.NewNotificationService(notificationRepo, destinationRepo, notificationDestRepo, broadcaster, metricsService)

	// Initialize notification processor (unified ActorPool + QueueConsumer + EnqueueWorker)
	notificationProcessor := actor.NewNotificationProcessor(
		redisClient,
		redisQueue,
		actor.NewActorDB(notificationDestRepo, installationRepo, teamsBotRepo, notificationRepo, projectRepo),
		services.NewTokenManagerAdapter(tokenManager),
		services.NewBillingServiceAdapter(billingService),
		cfg.Actor.MaxActors,
		notificationDestRepo,
	)
	notificationProcessor.Start(context.Background()) // Start in background
	logger.Info("Notification Processor started successfully")

	externalService := services.NewNotifyService(projectRepo, destinationRepo, notificationService)

	localStorage := storage.NewLocalStorage("uploads")
	fileRepo := repositories.NewFileRepository(db)
	fileService := services.NewFileService(fileRepo, localStorage)

	configService := services.NewConfigService()
	monitoringService := services.NewMonitoringService(
		metricsService,
		projectRepo,
		notificationRepo,
		destinationRepo,
	)

	// Create handlers
	companyHandler := companies.NewHandler(companyService)
	userHandler := users.NewHandler(userService)
	projectHandler := projects.NewHandler(projectService)
	botHandler := bots.NewHandler(teamsBotService)
	messagesHandler := messages.NewHandler(messagesService)
	destinationHandler := destinations.NewHandler(destinationService)
	notificationHandler := notifications.NewHandler(notificationService)
	provisionHandler := provision.NewHandler(provisionService)
	externalHandler := notify.NewHandler(externalService)
	billingHandler := billing.NewHandler(billingService)
	fileHandler := files.NewHandler(fileService)
	queueHandler := qhandler.NewHandler(redisClient)
	systemHandler := system.NewHandler(metricsService, configService)
	monitoringHandler := monitoring.NewMonitoringHandler(monitoringService)
	templateHandler := templates.NewHandler(templateService)
	directoryHandler := directory.NewHandler(directoryService)

	// Create server
	server := teamsnotification.NewServer(cfg.Server)

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
			c.JSON(200, gin.H{"message": "Login endpoint - TODO"})
		})
		auth.POST("/refresh", authMiddleware.RefreshToken)
	}

	// Register external/public API routes (api/v1)
	server.RegisterHandlers(externalHandler, provisionHandler, messagesHandler)

	// Register internal API routes (api/internal/v1)
	server.RegisterInternalHandlers(
		companyHandler,
		userHandler,
		projectHandler,
		botHandler,
		destinationHandler,
		notificationHandler,
		systemHandler,
		billingHandler,
		fileHandler,
		queueHandler,
		monitoringHandler,
		templateHandler,
		directoryHandler,
	)

	return &Application{
		Config:    cfg,
		Logger:    logger,
		DB:        db,
		Redis:     redisClient,
		Server:    server,
		ActorPool: notificationProcessor,
	}, nil
}

// Run starts the application server
func (app *Application) Run() error {
	app.Logger.Info("Starting server on port " + app.Config.Server.Port)
	return app.Server.Start()
}

// Close gracefully closes application resources
func (app *Application) Close() {
	if app.DB != nil {
		app.DB.Close()
	}
	if app.Redis != nil {
		app.Redis.Close()
	}
	if app.ActorPool != nil {
		app.ActorPool.Stop()
	}
}