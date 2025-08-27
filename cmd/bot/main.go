package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gladis/internal/ai"
	"gladis/internal/config"
	"gladis/internal/database"
	"gladis/internal/discord"
	"gladis/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	// Inicjalizacja loggera na samym początku
	logger.InitLogger()

	// Ładowanie .env
	if err := godotenv.Load(); err != nil {
		logger.Warn("Error loading .env file", err.Error())
	} else {
		logger.Info("Environment variables loaded successfully", "main")
	}

	// Sprawdzenie WEBHOOK_URL
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		logger.Warn("WEBHOOK_URL is not set. Logs will not be sent to webhook.", "main")
	} else {
		fmt.Println("WEBHOOK_URL loaded successfully:", webhookURL)
	}

	// Konfiguracja aplikacji
	cfg := config.Load()

	// Upewnij się, że w .env masz MONGO_DB i że config.Load mapuje to na cfg.MongoDBName
	// np. MONGO_URI=..., MONGO_DB=gladis
	if cfg.MongoURI == "" || cfg.MongoDBName == "" {
		logger.Fatal("Mongo configuration missing. Set MONGO_URI and MONGO_DB in .env", "main", nil)
	}

	// Połączenie z Mongo z wybraną bazą
	db, err := database.InitMongoDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", "main", err)
	}
	defer db.Close()

	// Usunięto insert modeli przy starcie zgodnie z Twoją prośbą

	// Inicjalizacja AI managera
	aiManager, err := ai.NewManager(cfg, db)
	if err != nil {
		logger.Fatal("Failed to create AI manager", "main", err)
	}
	defer aiManager.Close()

	// Tworzenie i start bota
	bot, err := discord.NewBot(cfg.DiscordToken, aiManager, db, cfg.GiphyAPIKey)
	if err != nil {
		logger.Fatal("Failed to create bot", "main", err)
	}

	if err := bot.Start(); err != nil {
		logger.Fatal("Failed to start bot", "main", err)
	}
	defer bot.Stop()

	// Rejestracja komend
	if err := bot.RegisterCommands(); err != nil {
		logger.Fatal("Failed to register commands", "main", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Shutting down...")
}
