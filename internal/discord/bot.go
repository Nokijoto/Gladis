package discord

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gladis/internal/ai"
	"gladis/internal/logger" // Import the logger package
)

type Bot struct {
	session        *discordgo.Session
	aiManager      *ai.Manager
	commandHandler *CommandHandler
}

func NewBot(token string, aiManager *ai.Manager) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		// Use logger.Error for non-fatal errors during initialization
		logger.Error("Failed to create Discord session", "NewBot", err)
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	bot := &Bot{
		session:   session,
		aiManager: aiManager,
	}

	bot.setupHandlers()
	return bot, nil
}

func (b *Bot) setupHandlers() {
	b.session.AddHandler(b.readyHandler)
	b.session.AddHandler(b.messageCreateHandler)
	b.session.AddHandler(b.interactionCreateHandler)

	// Initialize CommandHandler with the AI manager
	b.commandHandler = NewCommandHandler(b.aiManager)
}

func (b *Bot) readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	// Use the custom Log function for informational messages
	logger.Log(logger.LevelInfo, fmt.Sprintf("Logged in as: %v#%v", event.User.Username, event.User.Discriminator), "readyHandler")
}

func (b *Bot) messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.Contains(m.Content, fmt.Sprintf("<@%s>", s.State.User.ID)) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use logger.TryCatch for error handling
	var images [][]byte // Declare images here to be accessible after TryCatch
	err := logger.TryCatch(func() error {
		var err error // Declare err locally within the function
		images, err = b.downloadImages(m.Attachments)
		if err != nil {
			// Log the error and send a user-friendly message
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error processing images: %v", err))
			return err // Return the error so TryCatch can log it
		}
		return nil
	}, "messageCreateHandler")

	// If TryCatch returned an error, we've already sent a message and logged it.
	// We should not proceed if images failed to download.
	if err != nil {
		return
	}

	prompt := strings.ReplaceAll(m.Content, fmt.Sprintf("<@%s>", s.State.User.ID), "")
	prompt = strings.TrimSpace(prompt)

	// Fetch recent messages for context if contextLength is greater than 0
	if b.commandHandler.contextLength > 0 {
		messages, err := s.ChannelMessages(m.ChannelID, b.commandHandler.contextLength, m.ID, "", "")
		if err != nil {
			logger.Error("Failed to fetch channel messages for context", "messageCreateHandler", err)
			// Continue without context if there's an error fetching messages
		} else {
			// Build context string from fetched messages, oldest first
			contextMessages := ""
			for i := len(messages) - 1; i >= 0; i-- { // Iterate in reverse to get chronological order
				msg := messages[i]
				// Exclude bot's own messages and the current message
				if msg.Author.ID == s.State.User.ID || msg.ID == m.ID {
					continue
				}
				contextMessages += fmt.Sprintf("%s: %s\n", msg.Author.Username, msg.Content)
			}
			if contextMessages != "" {
				prompt = contextMessages + "\n" + prompt
			}
		}
	}

	// Prepend system prompt if it exists
	if b.commandHandler.systemPrompt != "" {
		prompt = b.commandHandler.systemPrompt + "\n" + prompt
	}

	// Use logger.TryCatch for AI content generation
	var response string // Declare response here to be accessible after TryCatch
	err = logger.TryCatch(func() error {
		var err error // Declare err locally within the function
		response, err = b.aiManager.GenerateContent(ctx, prompt, images)
		if err != nil {
			// Log the error and send a user-friendly message
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error generating response: %v", err))
			return err // Return the error so TryCatch can log it
		}
		return nil
	}, "messageCreateHandler")

	// If TryCatch returned an error, we've already sent a message and logged it.
	// We should not proceed if AI generation failed.
	if err != nil {
		return
	}

	if len(response) > 2000 {
		response = response[:1997] + "..."
	}

	s.ChannelMessageSend(m.ChannelID, response)
}

func (b *Bot) downloadImages(attachments []*discordgo.MessageAttachment) ([][]byte, error) {
	var images [][]byte
	var downloadErr error // Declare downloadErr here

	for _, attachment := range attachments {
		if !strings.HasPrefix(attachment.ContentType, "image/") {
			continue
		}

		// Use logger.TryCatch for http.Get and io.ReadAll
		// Wrap the entire loop body in TryCatch to handle potential errors in Get or ReadAll
		err := logger.TryCatch(func() error {
			resp, err := http.Get(attachment.URL)
			if err != nil {
				// Log the error and return it
				return fmt.Errorf("failed to download image: %w", err)
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				// Log the error and return it
				return fmt.Errorf("failed to read image data: %w", err)
			}
			images = append(images, data)
			return nil
		}, "downloadImages")

		// If TryCatch returned an error, we should stop processing this attachment
		// and return the error.
		if err != nil {
			downloadErr = err // Store the error to return it after the loop
			// We can break here if we want to stop processing further attachments on error,
			// or continue if we want to try processing other attachments.
			// For now, let's break to stop processing on the first error.
			break
		}
	}
	// Return the accumulated error if any occurred during the loop.
	return images, downloadErr
}

func (b *Bot) interactionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		b.handleApplicationCommand(s, i)
	case discordgo.InteractionMessageComponent:
		b.handleMessageComponent(s, i)
	}
}

func (b *Bot) handleMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	// The CustomID for model selection now includes the page number, e.g., "model_select_0"
	if strings.HasPrefix(data.CustomID, "model_select_") {
		if len(data.Values) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No model selected.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		chosenModelName := data.Values[0]
		err := b.aiManager.SetModel(chosenModelName)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content:    fmt.Sprintf("Failed to set model: %s", err.Error()),
					Flags:      discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{},
				},
			})
			return
		}
		// Update the command handler's model after successful setting in aiManager
		b.commandHandler.currentModel = b.aiManager.GetCurrentModel()

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    fmt.Sprintf("Model changed to: `%s` (Provider: `%s`)", b.commandHandler.currentModel.Name, b.commandHandler.currentModel.Provider),
				Flags:      discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{},
			},
		})
	} else if strings.HasPrefix(data.CustomID, "model_prev_") || strings.HasPrefix(data.CustomID, "model_next_") {
		b.commandHandler.HandleComponent(s, i)
	}
}

func (b *Bot) handleApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		b.commandHandler.HandlePing(s, i)
	case "help":
		b.commandHandler.HandleHelp(s, i)
	case "setmodel":
		b.commandHandler.HandleSetModel(s, i)
	case "info":
		b.commandHandler.HandleInfo(s, i)
	case "setsystemprompt":
		b.commandHandler.HandleSetSystemPrompt(s, i)
	case "setcontext":
		b.commandHandler.HandleSetContext(s, i)
	}
}

func (b *Bot) Start() error {
	return b.session.Open()
}

func (b *Bot) Stop() error {
	return b.session.Close()
}

func (b *Bot) RegisterCommands() error {
	commands := Commands
	_, err := b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", commands)
	return err
}

func (b *Bot) UpdateCommandHandlerModel() {
	b.commandHandler.currentModel = b.aiManager.GetCurrentModel()
}

func (b *Bot) GetSystemPrompt() string {
	return b.commandHandler.systemPrompt
}
