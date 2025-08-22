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

	"github.com/bwmarrin/discordgo"
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
	} else {
		log.Println("WEBHOOK_URL loaded successfully")
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
		go sendToWebhookAsEmbed(logEntry)
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

// sendToWebhookAsEmbed sends a log entry as an embed to the configured webhook URL.
func sendToWebhookAsEmbed(logEntry LogMessage) {
	// Create the embed to send to the webhook
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🔴 %s Log", logEntry.Level),
		Description: logEntry.Message,
		Color:       getLogColor(logEntry.Level),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Context**",
				Value:  logEntry.Context,
				Inline: false,
			},
			{
				Name:   "**Timestamp**",
				Value:  logEntry.Timestamp.Format(time.RFC3339),
				Inline: true,
			},
			{
				Name:   "**Error (if any)**",
				Value:  logEntry.Error,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Log sent at %s", logEntry.Timestamp.Format(time.RFC3339)),
		},
	}

	// Send the embed to the webhook
	reqBody := map[string]interface{}{
		"embeds": []interface{}{embed},
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Failed to marshal embed message: %v", err)
		return
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(reqJSON))
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

	// Log the response code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("Webhook returned non-OK status: %s", resp.Status)
	} else {
		log.Printf("Webhook response status: %s", resp.Status)
	}
}

// getLogColor returns a color code based on the log level.
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
