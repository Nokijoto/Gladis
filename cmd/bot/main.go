package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gladis/internal/ai"
	"gladis/internal/config"
	"gladis/internal/discord"
)

func main() {
	cfg := config.Load()

	aiManager, err := ai.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create AI manager: %v", err)
	}
	defer aiManager.Close()

	bot, err := discord.NewBot(cfg.DiscordToken, aiManager)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	defer bot.Stop()

	if err := bot.RegisterCommands(); err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Shutting down...")
}
