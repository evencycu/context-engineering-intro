package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger represents application logger
type Logger struct {
	*logrus.Logger
}

// Config represents logger configuration
type Config struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// NewLogger creates a new logger instance
func NewLogger(config Config) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set formatter
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// Set output
	switch config.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "file":
		if config.FilePath != "" {
			// TODO: Implement file output with rotation
			logger.SetOutput(os.Stdout)
		} else {
			logger.SetOutput(os.Stdout)
		}
	default:
		logger.SetOutput(os.Stdout)
	}

	return &Logger{Logger: logger}
}

// WithFields creates a new logger with fields
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithField creates a new logger with a single field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithError creates a new logger with an error field
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithContext creates a new logger with context fields
func (l *Logger) WithContext(ctx interface{}) *logrus.Entry {
	return l.Logger.WithField("context", ctx)
}

// DefaultLogger creates a default logger
func DefaultLogger() *Logger {
	return NewLogger(Config{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})
}
