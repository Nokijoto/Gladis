package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

// LogMessage represents the structure of a log entry.
type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Context   string    `json:"context,omitempty"` // e.g., file, function, or custom context
	Error     string    `json:"error,omitempty"`   // For error logs
}

var (
	webhookURL string
	logLevel   LogLevel = LevelInfo // Default log level
)

// InitLogger initializes the logger with configuration.
// It reads the webhook URL and log level from environment variables.
func InitLogger() {
	webhookURL = os.Getenv("WEBHOOK_URL")
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr != "" {
		logLevel = LogLevel(strings.ToUpper(levelStr))
	}
	if webhookURL == "" {
		log.Println("WEBHOOK_URL not set. Logs will not be sent to webhook.")
	}
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

	// Send to webhook if URL is set and level is appropriate
	if webhookURL != "" && (level == LevelError || level == LevelFatal || level == LevelWarn) {
		go sendToWebhook(logJSON)
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

// sendToWebhook sends a JSON log message to the configured webhook URL.
func sendToWebhook(logData []byte) {
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(logData))
	if err != nil {
		log.Printf("Failed to create webhook request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send log to webhook: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("Webhook returned non-OK status: %s", resp.Status)
	}
}
