package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/bots"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/companies"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/destinations"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/external"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/messages"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/notifications"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/projects"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/provision"
	"github.com/evencycu/TeamsNotifyGoV2/internal/api/handlers/users"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	server *http.Server
}

// Config represents server configuration
type Config struct {
	Port         string
	Environment  string
	RateLimit    rate.Limit
	Burst        int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewServer creates a new API server
func NewServer(cfg Config) *Server {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting
	limiter := rate.NewLimiter(cfg.RateLimit, cfg.Burst)
	router.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// API v1 routes will be registered by handlers

	return &Server{
		router: router,
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// RegisterHandlers registers all API handlers
func (s *Server) RegisterHandlers(handlers ...Handler) {
	v1 := s.router.Group("/api/v1")
	for _, handler := range handlers {
		handler.RegisterRoutes(v1)
	}
}

// Handler interface for registering routes
type Handler interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

// Start starts the server
func (s *Server) Start() error {
	// Start server in goroutine
	go func() {
		log.Printf("Starting server on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
	return nil
}

// GetRouter returns the router for testing
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// RegisterRoutes registers all API routes
func (s *Server) RegisterRoutes(
	companyHandler *companies.Handler,
	userHandler *users.Handler,
	projectHandler *projects.Handler,
	botHandler *bots.Handler,
	destinationHandler *destinations.Handler,
	notificationHandler *notifications.Handler,
	messagesHandler *messages.Handler,
	provisionHandler *provision.Handler,
	externalHandler *external.Handler,
) {
	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Company routes
		v1.POST("/companies", companyHandler.CreateCompany)
		v1.GET("/companies", companyHandler.ListCompanies)
		v1.GET("/companies/:id", companyHandler.GetCompany)
		v1.PUT("/companies/:id", companyHandler.UpdateCompany)
		v1.DELETE("/companies/:id", companyHandler.DeleteCompany)
		v1.PATCH("/companies/:id/status", companyHandler.UpdateCompanyStatus)
		v1.PATCH("/companies/:id/billing", companyHandler.UpdateBillingStatus)

		// User routes
		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users", userHandler.ListUsers)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
		v1.PATCH("/users/:id/password", userHandler.ChangePassword)
		v1.PATCH("/users/:id/api-key", userHandler.GenerateAPIKey)
		v1.DELETE("/users/:id/api-key", userHandler.RevokeAPIKey)
		v1.GET("/users/company/:companyId", userHandler.GetUsersByCompany)
		v1.GET("/users/role/:role", userHandler.GetUsersByRole)

		// Project routes
		v1.POST("/projects", projectHandler.CreateProject)
		v1.GET("/projects", projectHandler.ListProjects)
		v1.GET("/projects/:id", projectHandler.GetProject)
		v1.PUT("/projects/:id", projectHandler.UpdateProject)
		v1.DELETE("/projects/:id", projectHandler.DeleteProject)
		v1.PATCH("/projects/:id/limits", projectHandler.UpdateLimits)
		v1.GET("/projects/company/:companyId", projectHandler.GetProjectsByCompany)
		v1.GET("/projects/key/:keyName", projectHandler.GetProjectByKeyName)

		// Bot routes
		v1.POST("/bots/platform", botHandler.CreateTeamsBotService)
		v1.GET("/bots/platform", botHandler.ListteamsBotServices)
		v1.GET("/bots/platform/:id", botHandler.GetteamsBotService)
		v1.PUT("/bots/platform/:id", botHandler.UpdateteamsBotService)
		v1.DELETE("/bots/platform/:id", botHandler.DeleteteamsBotService)
		v1.PATCH("/bots/platform/:id/status", botHandler.UpdateteamsBotServiceStatus)
		v1.PATCH("/bots/platform/:id/capabilities", botHandler.UpdateteamsBotServiceCapabilities)
		v1.POST("/bots/platform/:id/test", botHandler.TestteamsBotServiceConnection)

		// third-party routes removed
		// v1.GET("/bots/company/:companyId", botHandler.GetBotsByCompany) // removed with third-party
		v1.GET("/bots/status/:status", botHandler.GetBotsByStatus)

		// Destination routes
		v1.POST("/destinations", destinationHandler.CreateDestination)
		v1.GET("/destinations", destinationHandler.ListDestinations)
		v1.GET("/destinations/:id", destinationHandler.GetDestination)
		v1.PUT("/destinations/:id", destinationHandler.UpdateDestination)
		v1.DELETE("/destinations/:id", destinationHandler.DeleteDestination)
		v1.PATCH("/destinations/:id/targets", destinationHandler.UpdateTargets)
		v1.POST("/destinations/:id/validate", destinationHandler.ValidateTargets)
		v1.GET("/destinations/project/:projectId", destinationHandler.GetDestinationsByProject)
		v1.GET("/destinations/bot/:botId", destinationHandler.GetDestinationsByBot)
		v1.GET("/destinations/search", destinationHandler.SearchDestinations)

		// Notification routes
		v1.POST("/notifications", notificationHandler.SendNotification)
		v1.GET("/notifications", notificationHandler.ListNotifications)
		v1.GET("/notifications/:id", notificationHandler.GetNotification)
		v1.POST("/notifications/:id/retry", notificationHandler.RetryNotification)
		v1.DELETE("/notifications/:id", notificationHandler.CancelNotification)
		v1.GET("/notifications/project/:projectId", notificationHandler.GetNotificationsByProject)
		v1.GET("/notifications/sender/:senderId", notificationHandler.GetNotificationsBySender)
		v1.GET("/notifications/status/:status", notificationHandler.GetNotificationsByStatus)
		v1.GET("/notifications/date-range", notificationHandler.GetNotificationsByDateRange)

		// Bot Framework messages webhook
		v1.POST("/messages", messagesHandler.Handle)
		v1.POST("/messages/proactive/test", messagesHandler.ProactiveTest)

		// Provision routes
		provisionHandler.RegisterRoutes(v1)

		// External API routes
		externalHandler.RegisterRoutes(v1)
	}
}
