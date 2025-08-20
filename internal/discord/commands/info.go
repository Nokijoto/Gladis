package commands

import (
	"fmt"
	"gladis/internal/ai"
	"time" // Add import for time package

	"github.com/bwmarrin/discordgo"
)

var InfoCommand = &discordgo.ApplicationCommand{
	Name:        "info",
	Description: "Displays information about the bot, including the current model.",
}

// HandleInfo is the handler for the /info command.
// It requires the AI manager, system prompt, and context length to display information.
func HandleInfo(s *discordgo.Session, i *discordgo.InteractionCreate, aiManager *ai.Manager, systemPrompt *string, contextLength *int) {
	prompt := *systemPrompt
	if prompt == "" {
		prompt = "no system prompt set"
	}

	currentModel := aiManager.GetCurrentModel() // Ensure current model is up-to-date

	// Create a more visually appealing embed
	embed := &discordgo.MessageEmbed{
		Title:       "🤖 Bot Information",
		Description: "Here is some information about the current bot configuration.",
		Color:       0x3498db, // Set a friendly blue color
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://example.com/bot-icon.png", // Add an icon or bot image if available
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**🔧 Current Model**",
				Value:  fmt.Sprintf("`%s` (Provider: `%s`)", currentModel.Name, currentModel.Provider),
				Inline: false,
			},
			{
				Name:   "**💬 System Prompt**",
				Value:  fmt.Sprintf("`%s`", prompt),
				Inline: false,
			},
			{
				Name:   "**📄 Context Length**",
				Value:  fmt.Sprintf("`%d` messages", *contextLength),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Powered by Gladis | Keep learning!",
			IconURL: "https://example.com/footer-icon.png", // Footer icon, add one if available
		},
		Timestamp: time.Now().Format(time.RFC3339), // Use time.Now().Format() to get the current timestamp
	}

	// Send the embed response to the user
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
