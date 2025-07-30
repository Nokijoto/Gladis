package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Discord Bot starting...")

	discordToken := os.Getenv("DISCORD_TOKEN")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	if discordToken == "" || geminiAPIKey == "" {
		fmt.Println("Error: DISCORD_TOKEN and GEMINI_API_KEY must be set as environment variables.")
		return // Or handle error more gracefully
	}

	fmt.Println("Successfully loaded configuration.")

	// TODO: Initialize Discord session
	// TODO: Initialize Gemini client
	// TODO: Implement message handling and AI interaction

	// Placeholder for bot logic
	fmt.Println("Bot is running (placeholder). Press Ctrl+C to stop.")

	// Keep the program running until interrupted
	select {}
}
