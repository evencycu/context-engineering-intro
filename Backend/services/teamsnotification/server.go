package teamsnotification

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
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/evencycu/TeamsNotifyGoV3/configs"
	"github.com/evencycu/TeamsNotifyGoV3/libs/middleware"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	server *http.Server
}

// NewServer creates a new API server
func NewServer(serverCfg configs.ServerConfig) *Server {
	// Set Gin mode
	if serverCfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware in performance-optimized order
	// 1. Recovery first (safety)
	router.Use(gin.Recovery())

	// 2. Request ID for tracing (lightweight)
	router.Use(requestid.New())

	// 3. Rate limiting (early rejection for performance)
	limiter := rate.NewLimiter(rate.Limit(serverCfg.RateLimit), serverCfg.Burst)
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

	// 4. CORS (after rate limiting)
	// Create two CORS middlewares: one for Bot Framework (allow all), one for others (restricted)
	botFrameworkCORS := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Bot Framework doesn't use credentials
		MaxAge:           12 * time.Hour,
	})
	defaultCORS := cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:8082", "https://yourdomain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	// Conditional CORS: use botFrameworkCORS for /api/v1/messages, defaultCORS for others
	router.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/api/v1/messages" {
			botFrameworkCORS(c)
			return
		}
		defaultCORS(c)
	})

	// 5. Business metrics logger (for business events)
	businessMetricsMiddleware := middleware.NewBusinessMetricsMiddleware(logrus.New())
	router.Use(businessMetricsMiddleware.BusinessMetricsLogger())

	// 6. Logger last (most expensive, after other checks)
	router.Use(gin.Logger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})
	// Explicitly support HEAD for Docker HEALTHCHECK (wget --spider issues HEAD)
	router.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// System endpoints are now handled by the system handler under /api/v1

	// API v1 routes will be registered by handlers

	return &Server{
		router: router,
		server: &http.Server{
			Addr:         ":" + serverCfg.Port,
			Handler:      router,
			ReadTimeout:  serverCfg.ReadTimeout,
			WriteTimeout: serverCfg.WriteTimeout,
		},
	}
}

// RegisterHandlers registers all API handlers (external/public APIs)
func (s *Server) RegisterHandlers(handlers ...Handler) {
	v1 := s.router.Group("/api/v1")
	for _, handler := range handlers {
		handler.RegisterRoutes(v1)
	}
}

// RegisterInternalHandlers registers internal API handlers
func (s *Server) RegisterInternalHandlers(handlers ...Handler) {
	internalV1 := s.router.Group("/internal/v1")
	for _, handler := range handlers {
		handler.RegisterRoutes(internalV1)
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
