// Package config provides configuration loading and validation for the Gladis Discord bot.
// It loads configuration from environment variables and validates that required
// values are present before the bot starts.
package config

import (
	"fmt"
	"os"
)

type Config struct {
	DiscordToken     string
	GeminiAPIKey     string
	OpenRouterAPIKey string
	Environment      string
	StartupChannelID string
	StartupGuildID   string
	WebhookURL       string // Added WebhookURL
	GiphyAPIKey      string // Added GiphyAPIKey
}

func Load() *Config {
	return &Config{
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		Environment:      os.Getenv("ENVIRONMENT"),
		WebhookURL:       os.Getenv("WEBHOOK_URL"),   // Added retrieval for WEBHOOK_URL
		GiphyAPIKey:      os.Getenv("GIPHY_API_KEY"), // Added retrieval for GIPHY_API_KEY
		StartupChannelID: os.Getenv("STARTUP_CHANNEL_ID"),
		StartupGuildID:   os.Getenv("STARTUP_GUILD_ID"),
	}
}

// Validate checks if the configuration is valid for the bot to run
func (c *Config) Validate() error {
	if c.DiscordToken == "" {
		return fmt.Errorf("DISCORD_TOKEN is required")
	}

	if c.GeminiAPIKey == "" && c.OpenRouterAPIKey == "" {
		return fmt.Errorf("at least one AI API key (GEMINI_API_KEY or OPENROUTER_API_KEY) is required")
	}

	return nil
}
