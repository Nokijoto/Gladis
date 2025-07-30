package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
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

	// Initialize Discord session
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Register messageCreate as a message handler
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// Keep the program running until interrupted
	select {}
}

// messageCreate will be called (due to AddHandler) every time a new
// message is received with the prefix "!".
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "!ping"
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// TODO: Implement Gemini AI interaction here
}
