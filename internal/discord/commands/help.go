package commands

import (
	"github.com/bwmarrin/discordgo"
)

var HelpCommand = &discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Displays this help message.",
}

func HandleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Help message with formatted Embed
	helpMessage := &discordgo.MessageEmbed{
		Title:       "🤖 Bot Help",
		Description: "Here are the available commands for this bot:",
		Color:       0x3498db, // Friendly blue color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "`/help` - Help",
				Value:  "Displays this help message.",
				Inline: false,
			},
			{
				Name:   "`/ping` - PONG",
				Value:  "Tests the bot's responsiveness.",
				Inline: false,
			},
			{
				Name:   "`/setcontext` - Set Context",
				Value:  "Sets the number of previous messages to send to the bot (0-10).",
				Inline: false,
			},
			{
				Name:   "`/setmodel` - Set Model",
				Value:  "Interactively change the AI model being used.",
				Inline: false,
			},
			{
				Name:   "`/setsystemprompt` - Set System Prompt",
				Value:  "Sets a custom message that is always sent to the bot.",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Powered by Gladis",
		},
	}

	// Send the formatted embed to the user
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{helpMessage},
		},
	})
}
