package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var SetSystemPromptCommand = &discordgo.ApplicationCommand{
	Name:        "setsystemprompt",
	Description: "Sets a system prompt to be sent as the first message in every interaction.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "prompt",
			Description: "The system prompt to set.",
			Required:    true,
		},
	},
}

// HandleSetSystemPrompt is the handler for the /setsystemprompt command.
// It requires a pointer to the system prompt string to update it.
func HandleSetSystemPrompt(s *discordgo.Session, i *discordgo.InteractionCreate, systemPrompt *string) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var newPrompt string
	if opt, ok := optionMap["prompt"]; ok {
		newPrompt = opt.StringValue()
	}

	*systemPrompt = newPrompt

	// Create Embed response
	embed := &discordgo.MessageEmbed{
		Title:       "✅ System Prompt Set",
		Description: fmt.Sprintf("The system prompt has been successfully updated to the following:"),
		Color:       0x00FF00, // Green color for success
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**New System Prompt**",
				Value:  fmt.Sprintf("`%s`", newPrompt),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Bot powered by Gladis",
		},
	}

	// Send the response with Embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
