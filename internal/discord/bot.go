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
)

type Bot struct {
	session        *discordgo.Session
	aiManager      *ai.Manager
	commandHandler *CommandHandler
}

func NewBot(token string, aiManager *ai.Manager) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
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

	b.commandHandler = NewCommandHandler(b.aiManager.GetCurrentModel())
}

func (b *Bot) readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Logged in as: %v#%v\n", event.User.Username, event.User.Discriminator)
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

	images, err := b.downloadImages(m.Attachments)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error processing images: %v", err))
		return
	}

	prompt := strings.ReplaceAll(m.Content, fmt.Sprintf("<@%s>", s.State.User.ID), "")
	prompt = strings.TrimSpace(prompt)

	response, err := b.aiManager.GenerateContent(ctx, prompt, images)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error generating response: %v", err))
		return
	}

	if len(response) > 2000 {
		response = response[:1997] + "..."
	}

	s.ChannelMessageSend(m.ChannelID, response)
}

func (b *Bot) downloadImages(attachments []*discordgo.MessageAttachment) ([][]byte, error) {
	var images [][]byte

	for _, attachment := range attachments {
		if !strings.HasPrefix(attachment.ContentType, "image/") {
			continue
		}

		resp, err := http.Get(attachment.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to download image: %w", err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read image data: %w", err)
		}

		images = append(images, data)
	}

	return images, nil
}

func (b *Bot) interactionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		b.handleApplicationCommand(s, i)
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
