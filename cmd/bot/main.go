package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gladis/internal/ai"
	"gladis/internal/config" // Changed import path
	"gladis/internal/discord"
	"gladis/internal/logger"

	"github.com/joho/godotenv" // Import godotenv
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// Log a warning but continue running the app
		logger.Warn("Error loading .env file", err.Error())
	} else {
		// Log successful loading of environment variables
		logger.Info("Environment variables loaded successfully", "main")
	}

	// Initialize the logger
	logger.InitLogger()

	// Check if the WEBHOOK_URL is properly loaded
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		logger.Warn("WEBHOOK_URL is not set. Logs will not be sent to webhook.", "main")
	} else {
		// Optional: Log the webhook URL for verification purposes (only in debug mode)
		fmt.Println("WEBHOOK_URL loaded successfully:", webhookURL)
	}

	cfg := config.Load()

	// Initialize AI manager
	aiManager, err := ai.NewManager(cfg)
	if err != nil {
		logger.Fatal("Failed to create AI manager", "main", err)
	}
	defer aiManager.Close()

	// Create and start the bot
	bot, err := discord.NewBot(cfg.DiscordToken, aiManager)
	if err != nil {
		logger.Fatal("Failed to create bot", "main", err)
	}

	// Start the bot
	if err := bot.Start(); err != nil {
		logger.Fatal("Failed to start bot", "main", err)
	}
	defer bot.Stop()

	// Register commands
	if err := bot.RegisterCommands(); err != nil {
		logger.Fatal("Failed to register commands", "main", err)
	}

	// Optional: Set the system prompt for the bot if necessary
	// bot.commandHandler.systemPrompt = cfg.DefaultSystemPrompt // Example if you want to add this feature

	// Log that the bot is running
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Graceful shutdown
	fmt.Println("Shutting down...")
}
