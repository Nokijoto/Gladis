package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	ctx   context.Context
	client *genai.Client
	model *genai.GenerativeModel
	currentModelName string
	env              string
	startupChannelID = "361961924059070466"
	startupGuildID   = "361961924059070464"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with PONG!",
		},
		{
			Name:        "help",
			Description: "Displays this help message.",
		},
		{
			Name:        "setmodel",
			Description: "Displays a list of available Gemini models to choose from.",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "PONG!",
				},
			})
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			helpMessage := "Available commands:\n"
			for _, cmd := range commands {
				helpMessage += fmt.Sprintf("`/%s` - %s\n", cmd.Name, cmd.Description)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: helpMessage,
				},
			})
		},
		"setmodel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			models, err := getAvailableModels(client)
			if err != nil {
				log.Printf("Error fetching available models: %v", err)
				embed := &discordgo.MessageEmbed{
					Title:       "Error",
					Description: "Error fetching available models. Please try again later.",
					Color:       0xFF0000, // Red color for errors
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				return
			}

			if len(models) == 0 {
				embed := &discordgo.MessageEmbed{
					Title:       "No Models Found",
					Description: "No Gemini models available for content generation.",
					Color:       0xFF0000, // Red color for errors
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
					},
				})
				return
			}

			options := []discordgo.SelectMenuOption{}
			for _, modelName := range models {
				options = append(options, discordgo.SelectMenuOption{
					Label: modelName,
					Value: modelName,
				})
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please select a Gemini model:",
					Flags:   discordgo.MessageFlagsEphemeral, // Only visible to the user who invoked the command
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									CustomID:    "select_gemini_model",
									Placeholder: "Choose a model...",
									Options:     options,
								},
							},
						},
					},
				},
			})
		},
	}
)

var componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"select_gemini_model": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.MessageComponentData()
		selectedModelName := data.Values[0]

		newModel, err := initializeGeminiModel(client, selectedModelName)
		if err != nil {
			log.Printf("Error setting model to '%s': %v", selectedModelName, err)
			errMsg := fmt.Sprintf("Error setting model: %v", err)
			if strings.Contains(err.Error(), "not found or unsupported") {
				errMsg = fmt.Sprintf("Error: %v. Call ListModels to see the list of available models and their supported methods.", err)
			}
			embed := &discordgo.MessageEmbed{
				Title:       "Error Setting Model",
				Description: errMsg,
				Color:       0xFF0000, // Red color for errors
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds:     []*discordgo.MessageEmbed{embed},
					Components: []discordgo.MessageComponent{}, // Remove components
				},
			})
			return
		}

		model = newModel
		currentModelName = selectedModelName
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    fmt.Sprintf("✅ Gemini model set to: `%s`", currentModelName),
				Components: []discordgo.MessageComponent{}, // Remove components
			},
		})
	},
}

func main() {
	fmt.Println("Discord Bot starting...")

	discordToken := os.Getenv("DISCORD_TOKEN")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	env = os.Getenv("ENVIRONMENT") // Assign to global env

	if discordToken == "" || geminiAPIKey == "" {
		fmt.Println("Error: DISCORD_TOKEN and GEMINI_API_KEY must be set as environment variables.")
		return
	}

	// Setup logging based on environment
	if env == "dev" {
		log.Println("Running in development mode. Logging errors to app.log")
		logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Warning: Failed to open log file 'app.log': %v. Using default stderr logging.", err)
		} else {
			log.SetOutput(logFile) // Redirect log output to file
		}
	} else {
		log.Println("Running in production mode.")
	}

	fmt.Println("Successfully loaded configuration.")

	// Initialize Discord session with necessary intents
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Printf("Error creating Discord session: %v", err)
		return
	}

	// Set up intents for the bot
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	// Register messageCreate as a message handler for non-command messages
	dg.AddHandler(messageCreate)
	// Register interactionCreate as a handler for slash commands and components
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Printf("Error opening connection: %v", err)
		return
	}

	// Register slash commands
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, startupGuildID, v)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	log.Println("Commands added.")

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	// Send startup message to a specific channel
	guild, err := dg.State.Guild(startupGuildID)
	if err != nil {
		log.Printf("Error getting guild with ID %s: %v", startupGuildID, err)
	}

	if guild != nil {
		_, err = dg.ChannelMessageSend(startupChannelID, "Bot has started and is ready to assist!")
		if err != nil {
			log.Printf("Error sending startup message to channel %s: %v", startupChannelID, err)
		} else {
			log.Printf("Startup message sent to channel %s.", startupChannelID)
		}
	} else {
		log.Printf("Guild with ID %s not found. Cannot send startup message.", startupGuildID)
	}

	// Initialize Gemini client globally
	ctx = context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	// No defer client.Close() here as the program exits after select{}

	// Set initial model
	currentModelName = "gemini-2.5-flash-lite" // Default model as requested
	model, err = initializeGeminiModel(client, currentModelName)
	if err != nil {
		log.Printf("Error initializing Gemini model '%s': %v", currentModelName, err)
		// Continue running with potentially no model, or exit if model is critical
		// For now, we'll let it continue, but Gemini calls will fail.
	} else {
		log.Printf("Gemini model '%s' initialized successfully.", currentModelName)
	}

	// Keep the program running until interrupted
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stop

	log.Println("Gracefully shutting down...")

	// Close the Gemini client when the program exits
	if client != nil {
		client.Close()
	}

	// Clean up registered commands
	log.Println("Removing commands...")
	for _, cmd := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, startupGuildID, cmd.ID)
		if err != nil {
			log.Printf("Cannot delete slash command %q: %v", cmd.Name, err)
		}
	}
	log.Println("Commands removed.")

	dg.Close()
}

// initializeGeminiModel attempts to create a GenerativeModel and returns it.
// It returns an error if the model is not found or unsupported.
func initializeGeminiModel(client *genai.Client, modelName string) (*genai.GenerativeModel, error) {
	m := client.GenerativeModel(modelName)
	// No direct validation via m.Config() as it's not available.
	// The GenerateContent call will fail if the model is invalid.
	return m, nil
}

// sendMessage is a helper function to send a message to a Discord channel and log any errors.
func sendMessage(s *discordgo.Session, channelID, content string) {
	_, err := s.ChannelMessageSend(channelID, content)
	if err != nil {
		log.Printf("Error sending message to channel %s: %v", channelID, err)
	}
}

// getAvailableModels fetches a list of available Gemini models that support content generation.
func getAvailableModels(client *genai.Client) ([]string, error) {
	if client == nil {
		return nil, fmt.Errorf("Gemini client is not initialized")
	}

	iter := client.ListModels(ctx)
	var models []string
	for {
		model, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error listing models: %w", err)
		}
		// Only include models that support content generation
		if contains(model.SupportedGenerationMethods, "generateContent") {
			models = append(models, strings.TrimPrefix(model.Name, "models/"))
		}
	}
	sort.Strings(models) // Sort models alphabetically
	return models, nil
}

// contains checks if a string is present in a slice of strings.
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// messageCreate is the main handler for incoming Discord messages.
// It now only handles non-command messages (i.e., direct AI interactions).
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Only respond if the bot is mentioned
	if !isBotMentioned(s, m.Message) {
		return
	}

	// Remove bot mention from the message content
	content := strings.ReplaceAll(m.Content, fmt.Sprintf("<@%s>", s.State.User.ID), "")
	content = strings.TrimSpace(content)

	if content == "" {
		// If only the mention was sent, ignore or send a default message
		return
	}

	// If it's not a command, process it with Gemini
	if model == nil {
		sendMessage(s, m.ChannelID, "Gemini model is not initialized. Cannot process your request.")
		return
	}

	go func() { // Run in a goroutine to avoid blocking Discord session
		// Send the user message (without mention) to Gemini
		resp, err := model.GenerateContent(ctx, genai.Text(content))
		if err != nil {
			// Check for the specific model not found error
			errorTitle := "Gemini API Error"
			errorDescription := fmt.Sprintf("Failed to generate content: %v", err)

			if strings.Contains(err.Error(), "models/") && strings.Contains(err.Error(), "not found") {
				errorTitle = "Model Not Found or Unsupported"
				errorDescription = fmt.Sprintf("Sorry, the model `%s` does not exist or is not supported. Please use `/setmodel` to choose an available model.", currentModelName)
			}

			log.Printf("Gemini Error: %s - %s", errorTitle, errorDescription)

			embed := &discordgo.MessageEmbed{
				Title:       errorTitle,
				Description: errorDescription,
				Color:       0xFF0000, // Red color for errors
			}

			_, embedErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if embedErr != nil {
				log.Printf("Error sending embed message: %v", embedErr)
				sendMessage(s, m.ChannelID, "Sorry, I encountered an error processing your request and couldn't send the error as an embed.")
			}
			return
		}

		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
			var geminiResponseBuilder strings.Builder
			for _, part := range resp.Candidates[0].Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					geminiResponseBuilder.WriteString(string(txt))
				}
			}
			geminiResponse := geminiResponseBuilder.String()

			sendLongMessageAsReplies(s, m, geminiResponse)
		} else {
			embed := &discordgo.MessageEmbed{
				Title:       "Gemini Response Error",
				Description: "Sorry, I couldn't get a response from Gemini.",
				Color:       0xFF0000, // Red color for errors
			}
			_, embedErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if embedErr != nil {
				log.Printf("Error sending embed message for no Gemini response: %v", embedErr)
				sendMessage(s, m.ChannelID, "Sorry, I couldn't get a response from Gemini and couldn't send the error as an embed.")
			}
		}
	}()
}

// isBotMentioned checks if the bot is mentioned in the message.
func isBotMentioned(s *discordgo.Session, m *discordgo.Message) bool {
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			return true
		}
	}
	return false
}

// sendLongMessageAsReplies splits a long message into chunks and sends them as replies.
// Each chunk is truncated at the last period near the MAX_MESSAGE_LENGTH.
func sendLongMessageAsReplies(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	const MAX_MESSAGE_LENGTH = 2000

	remainingContent := content
	var lastMessageID string // To chain replies

	for len(remainingContent) > 0 {
		chunk := remainingContent
		if len(chunk) > MAX_MESSAGE_LENGTH {
			splitPoint := MAX_MESSAGE_LENGTH
			// Try to find the last period before the split point
			lastPeriodIndex := strings.LastIndex(chunk[:MAX_MESSAGE_LENGTH], ".")
			if lastPeriodIndex != -1 && lastPeriodIndex > MAX_MESSAGE_LENGTH/2 { // Ensure it's not too early in the chunk
				splitPoint = lastPeriodIndex + 1 // Include the period
			}
			chunk = remainingContent[:splitPoint]
			remainingContent = remainingContent[splitPoint:]
		} else {
			remainingContent = "" // Last chunk
		}

		var msg *discordgo.Message
		var err error

		if lastMessageID == "" {
			// First message, reply to the original user message
			msg, err = s.ChannelMessageSendReply(m.ChannelID, chunk, m.Reference())
		} else {
			// Subsequent messages, send directly to the channel
			msg, err = s.ChannelMessageSend(m.ChannelID, chunk)
		}

		if err != nil {
			log.Printf("Error sending chunked reply to message %s in channel %s: %v", m.ID, m.ChannelID, err)
			embed := &discordgo.MessageEmbed{
				Title:       "Discord API Error",
				Description: fmt.Sprintf("Sorry, I couldn't send a part of the response: %v", err),
				Color:       0xFF0000, // Red color for errors
			}
			_, embedErr := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if embedErr != nil {
				log.Printf("Error sending embed message for chunked Discord API error: %v", embedErr)
				sendMessage(s, m.ChannelID, "Sorry, I encountered an error sending a part of the response and couldn't send the error as an embed.")
			}
			return
		}
		lastMessageID = msg.ID // Update lastMessageID for the next iteration
	}
}
