package main

import (
	"github.com/evencycu/TeamsNotifyGoV3/configs"
	"github.com/evencycu/TeamsNotifyGoV3/pkg/app" // Import the new app package
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load default environment variables from configs/.env if present.
	// Existing OS environment variables take precedence over values in the file.
	_ = godotenv.Load("configs/.env")

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg := configs.LoadConfig()
	if err := cfg.ValidateConfig(); err != nil {
		logger.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize application
	application, err := app.NewApplication(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Close()

	// Start server
	if err := application.Run(); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
