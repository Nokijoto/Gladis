package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	ctx   context.Context
	model *genai.GenerativeModel
)

func main() {
	fmt.Println("Discord Bot starting...")

	discordToken := os.Getenv("DISCORD_TOKEN")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	if discordToken == "" || geminiAPIKey == "" {
		fmt.Println("Error: DISCORD_TOKEN and GEMINI_API_KEY must be set as environment variables.")
		return
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

	// Initialize Gemini client globally
	ctx = context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	// No defer client.Close() here as the program exits after select{}

	// Get the generative model
	model = client.GenerativeModel("gemini-pro") // Or another suitable model

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
	} else { // Process other messages with Gemini
		go func() { // Run in a goroutine to avoid blocking Discord session
			resp, err := model.GenerateContent(ctx, genai.Text(m.Content))
			if err != nil {
				log.Printf("Failed to generate content: %v", err)
				s.ChannelMessageSend(m.ChannelID, "Sorry, I encountered an error processing your request.")
				return
			}

			if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
				geminiResponse := ""
				for _, part := range resp.Candidates[0].Content.Parts {
					geminiResponse += fmt.Sprintf("%v", part)
				}
				s.ChannelMessageSend(m.ChannelID, geminiResponse)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't get a response from Gemini.")
			}
		}()
	}
}
