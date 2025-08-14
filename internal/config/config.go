package config

import "os"

type Config struct {
	DiscordToken     string
	GeminiAPIKey     string
	OpenRouterAPIKey string
	Environment      string
	StartupChannelID string
	StartupGuildID   string
}

func Load() *Config {
	return &Config{
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		Environment:      os.Getenv("ENVIRONMENT"),
		StartupChannelID: "361961924059070466",
		StartupGuildID:   "361961924059070464",
	}
}
