package commands

import (
	"github.com/bwmarrin/discordgo"
)

var PingCommand = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Responds with PONG!",
}

func HandlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Create an Embed for the response
	embed := &discordgo.MessageEmbed{
		Title:       "🏓 Ping Pong",
		Description: "Bot is alive and kicking! Here's your ping response.",
		Color:       0x00FF00, // Green color to represent a successful ping
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🏓 Ping Response",
				Value:  "PONG!",
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Bot powered by Gladis",
		},
	}

	// Respond with the embed message
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
