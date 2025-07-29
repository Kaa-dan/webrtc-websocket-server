package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log  *logrus.Logger
	once sync.Once
	mu   sync.RWMutex
)

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string // "json" or "text"
	Output     io.Writer
	TimeFormat string
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Format:     "json",
		Output:     os.Stdout,
		TimeFormat: time.RFC3339,
	}
}

// Init initializes the logger with given configuration
// This function is thread-safe and can be called multiple times
func Init(level string) error {
	return InitWithConfig(&Config{
		Level:      level,
		Format:     "json",
		Output:     os.Stdout,
		TimeFormat: time.RFC3339,
	})
}

// InitWithConfig initializes the logger with detailed configuration
func InitWithConfig(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	var initErr error
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()

		log = logrus.New()

		// Set output
		if config.Output != nil {
			log.SetOutput(config.Output)
		} else {
			log.SetOutput(os.Stdout)
		}

		// Set formatter
		switch config.Format {
		case "text":
			log.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: config.TimeFormat,
			})
		default: // json
			log.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: config.TimeFormat,
			})
		}

		// Set level
		logLevel, err := logrus.ParseLevel(config.Level)
		if err != nil {
			logLevel = logrus.InfoLevel
			initErr = fmt.Errorf("invalid log level '%s', defaulting to info: %w", config.Level, err)
		}
		log.SetLevel(logLevel)

		// Add hooks for production (optional)
		// log.AddHook(&SomeProductionHook{})
	})

	return initErr
}

// GetLogger returns the logger instance
// If not initialized, it will initialize with default config
func GetLogger() *logrus.Logger {
	mu.RLock()
	if log != nil {
		mu.RUnlock()
		return log
	}
	mu.RUnlock()

	// Initialize with default config if not already done
	if err := Init("info"); err != nil {
		// Fallback to a basic logger if initialization fails
		fallbackLogger := logrus.New()
		fallbackLogger.SetLevel(logrus.InfoLevel)
		fallbackLogger.SetFormatter(&logrus.JSONFormatter{})
		fallbackLogger.Warn("Failed to initialize logger properly, using fallback")
		return fallbackLogger
	}

	mu.RLock()
	defer mu.RUnlock()
	return log
}

// IsInitialized checks if logger has been initialized
func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return log != nil
}

// SetLevel dynamically changes log level
func SetLevel(level string) error {
	logger := GetLogger()
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level '%s': %w", level, err)
	}
	logger.SetLevel(logLevel)
	return nil
}

// GetLevel returns current log level
func GetLevel() string {
	return GetLogger().GetLevel().String()
}

// Structured logging methods
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}

// Standard logging methods
func Trace(args ...interface{}) {
	GetLogger().Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	GetLogger().Tracef(format, args...)
}

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

// Convenience methods for structured logging
func InfoWithFields(msg string, fields logrus.Fields) {
	GetLogger().WithFields(fields).Info(msg)
}

func ErrorWithFields(msg string, fields logrus.Fields) {
	GetLogger().WithFields(fields).Error(msg)
}

func WarnWithFields(msg string, fields logrus.Fields) {
	GetLogger().WithFields(fields).Warn(msg)
}

func DebugWithFields(msg string, fields logrus.Fields) {
	GetLogger().WithFields(fields).Debug(msg)
}
