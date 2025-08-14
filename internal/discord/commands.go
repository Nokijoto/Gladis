package discord

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	"gladis/internal/models"
)

var Commands = []*discordgo.ApplicationCommand{
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
		Description: "Displays a list of available models to choose from.",
	},
}

type CommandHandler struct {
	currentModel string
}

func NewCommandHandler(currentModel string) *CommandHandler {
	return &CommandHandler{currentModel: currentModel}
}

func (h *CommandHandler) HandlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "PONG!",
		},
	})
}

func (h *CommandHandler) HandleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	helpMessage := "Available commands:\n"
	for _, cmd := range Commands {
		helpMessage += fmt.Sprintf("`/%s` - %s\n", cmd.Name, cmd.Description)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMessage,
		},
	})
}

func (h *CommandHandler) HandleSetModel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modelsList := models.GetAvailableModels()
	if len(modelsList) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "No Models Found",
			Description: "No models available for content generation.",
			Color:       0xFF0000,
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		return
	}

	sort.Strings(modelsList)

	var components []discordgo.MessageComponent
	for _, model := range modelsList {
		buttonStyle := discordgo.PrimaryButton
		if model == h.currentModel {
			buttonStyle = discordgo.SuccessButton
		}
		components = append(components, discordgo.Button{
			Label:    model,
			Style:    buttonStyle,
			CustomID: "select_model_" + model,
		})
	}

	var actionRows []discordgo.MessageComponent
	for i := 0; i < len(components); i += 5 {
		end := i + 5
		if end > len(components) {
			end = len(components)
		}
		actionRows = append(actionRows, discordgo.ActionsRow{
			Components: components[i:end],
		})
	}

	responseContent := fmt.Sprintf("Current model: `%s`\nPlease select a new model:", h.currentModel)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    responseContent,
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: actionRows,
		},
	})
}
