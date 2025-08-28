package logger

import (
	"encoding/json"
	"fmt" // Keep fmt for console logging
	"log" // Keep log for internal logger errors
	"os"
	"strings"
	"time"

	"gladis/internal/database" // Import database package

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LogLevel defines the severity of a log message.
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// LogMessage represents the structure of a log entry for MongoDB.
type LogMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Level     LogLevel           `bson:"level" json:"level"`
	Message   string             `bson:"message" json:"message"`
	Context   string             `bson:"context,omitempty" json:"context,omitempty"` // e.g., file, function, or custom context
	Error     string             `bson:"error,omitempty" json:"error,omitempty"`     // For error logs
}

var (
	mongoDB  *database.MongoDB
	logLevel LogLevel = LevelInfo // Default log level
)

// InitLogger initializes the logger with a MongoDB instance and configuration.
func InitLogger(db *database.MongoDB) {
	mongoDB = db
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr != "" {
		logLevel = LogLevel(strings.ToUpper(levelStr))
	}
	log.Println("Logger initialized. Logs will be sent to MongoDB.")
}

// SetLevel sets the minimum log level to be logged.
func SetLevel(level LogLevel) {
	logLevel = level
}

// Log sends a log message with the specified level and context.
func Log(level LogLevel, message string, context string) {
	if level < logLevel {
		return
	}

	logEntry := LogMessage{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Context:   context,
	}

	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Failed to marshal log message: %v", err)
		return
	}

	// Log to console
	fmt.Printf("%s\n", logJSON)

	// Send to MongoDB asynchronously
	if mongoDB != nil {
		go func(entry LogMessage) {
			entry.ID = primitive.NewObjectID() // Generate a new ObjectID for each log entry
			if err := mongoDB.InsertLogEntry(entry); err != nil {
				log.Printf("Failed to insert log entry into MongoDB: %v", err)
			}
		}(logEntry)
	}
}

// Info logs an informational message.
func Info(message string, context string) {
	Log(LevelInfo, message, context)
}

// Debug logs a debug message.
func Debug(message string, context string) {
	Log(LevelDebug, message, context)
}

// Warn logs a warning message.
func Warn(message string, context string) {
	Log(LevelWarn, message, context)
}

// Error logs an error message.
func Error(message string, context string, err error) {
	errorMessage := message
	if err != nil {
		errorMessage = fmt.Sprintf("%s: %v", message, err)
	}
	Log(LevelError, errorMessage, context)
}

// Fatal logs a fatal error message and exits the application.
func Fatal(message string, context string, err error) {
	errorMessage := message
	if err != nil {
		errorMessage = fmt.Sprintf("%s: %v", message, err)
	}
	Log(LevelFatal, errorMessage, context)
	os.Exit(1)
}

// TryCatch executes a function and logs any errors using the custom logger.
// It takes a function that returns an error and a context string.
// It returns the error if one occurred, otherwise nil.
func TryCatch(fn func() error, context string) error {
	if err := fn(); err != nil {
		Error("An error occurred", context, err) // Use the logger's Error function
		return err
	}
	return nil
}

// getLogColor returns a color code based on the log level.
// This function is no longer needed as logs are sent to MongoDB, not Discord webhooks.
/*
func getLogColor(level LogLevel) int {
	switch level {
	case LevelDebug:
		return 0x7289DA // Blue for debug
	case LevelInfo:
		return 0x00FF00 // Green for info
	case LevelWarn:
		return 0xFFFF00 // Yellow for warnings
	case LevelError:
		return 0xFF0000 // Red for errors
	case LevelFatal:
		return 0xFF0000 // Red for fatal errors
	default:
		return 0xFFFFFF // White for unknown levels
	}
}
*/
