package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var SetContextCommand = &discordgo.ApplicationCommand{
	Name:        "setcontext",
	Description: "Sets the number of previous messages to include as context. Max 10.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "count",
			Description: "Number of messages to include from 0 to 10. Set to 0 to disable.",
			Required:    true,
		},
	},
}

func HandleSetContext(s *discordgo.Session, i *discordgo.InteractionCreate, contextLength *int) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var newContextLength int
	if opt, ok := optionMap["count"]; ok {
		newContextLength = int(opt.IntValue())
	}

	// Clamp context length between 0 and 10
	newContextLength = clampInt(newContextLength, 0, 10)
	*contextLength = newContextLength

	// Create Embed response based on the new context length
	var embed *discordgo.MessageEmbed
	if newContextLength > 0 {
		// Fetch recent messages in the channel to use as context
		channelID := i.ChannelID
		messages, err := s.ChannelMessages(channelID, newContextLength, "", "", "")
		if err != nil {
			embed = &discordgo.MessageEmbed{
				Title:       "❌ Error",
				Description: "Failed to fetch messages for context.",
				Color:       0xFF0000,
			}
		} else {
			// Prepare context (show messages)
			var contextMessages []string
			for _, message := range messages {
				// Check if the message is from the user or the bot
				if message.Author.ID == i.Member.User.ID || message.Author.ID == s.State.User.ID {
					// Append the message content, including the bot's message
					contextMessages = append(contextMessages, fmt.Sprintf("**%s**: %s", message.Author.Username, message.Content))
				}
			}

			// Join messages into a single string
			contextContent := strings.Join(contextMessages, "\n")

			// Update embed with context messages
			embed = &discordgo.MessageEmbed{
				Title:       "✅ Context Set",
				Description: fmt.Sprintf("Context length set to: `%d` messages.", newContextLength),
				Color:       0x00FF00, // Green for success
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "📜 Last Messages",
						Value:  contextContent,
						Inline: false,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Context updated successfully!",
				},
			}
		}
	} else {
		// When context length is set to 0, disable context
		embed = &discordgo.MessageEmbed{
			Title:       "❌ Context Disabled",
			Description: "Context messages have been disabled.",
			Color:       0xFF0000,
		}
	}

	// Send the embed response to the user
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

// clampInt is a helper function to clamp an integer value within a range.
func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
