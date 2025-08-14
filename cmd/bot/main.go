package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gladis/internal/ai"
	"gladis/internal/config"
	"gladis/internal/discord"
	"gladis/internal/logger"
)

func main() {
	logger.InitLogger() // Initialize the logger

	cfg := config.Load()

	aiManager, err := ai.NewManager(cfg)
	if err != nil {
		logger.Fatal("Failed to create AI manager", "main", err)
	}
	defer aiManager.Close()

	bot, err := discord.NewBot(cfg.DiscordToken, aiManager)
	if err != nil {
		logger.Fatal("Failed to create bot", "main", err)
	}

	if err := bot.Start(); err != nil {
		logger.Fatal("Failed to start bot", "main", err)
	}
	defer bot.Stop()

	if err := bot.RegisterCommands(); err != nil {
		logger.Fatal("Failed to register commands", "main", err)
	}

	// Update the command handler with the initial system prompt from config, if any
	// This assumes you might want to load a default system prompt from config
	// For now, it's initialized as empty in NewCommandHandler, but this is where
	// you'd pass it if it came from config.
	// bot.commandHandler.systemPrompt = cfg.DefaultSystemPrompt // Example if you add it to config

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Shutting down...")
}
