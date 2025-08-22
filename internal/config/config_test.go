package config

import (
	"os"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		discordToken     string
		geminiAPIKey     string
		openRouterAPIKey string
		expectError      bool
	}{
		{
			name:         "Valid config with Discord and Gemini",
			discordToken: "test-discord-token",
			geminiAPIKey: "test-gemini-key",
			expectError:  false,
		},
		{
			name:             "Valid config with Discord and OpenRouter",
			discordToken:     "test-discord-token",
			openRouterAPIKey: "test-openrouter-key",
			expectError:      false,
		},
		{
			name:             "Valid config with all keys",
			discordToken:     "test-discord-token",
			geminiAPIKey:     "test-gemini-key",
			openRouterAPIKey: "test-openrouter-key",
			expectError:      false,
		},
		{
			name:        "Missing Discord token",
			geminiAPIKey: "test-gemini-key",
			expectError: true,
		},
		{
			name:         "Missing AI API keys",
			discordToken: "test-discord-token",
			expectError:  true,
		},
		{
			name:        "Empty config",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				DiscordToken:     tt.discordToken,
				GeminiAPIKey:     tt.geminiAPIKey,
				OpenRouterAPIKey: tt.openRouterAPIKey,
			}

			err := cfg.Validate()
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestConfigLoad(t *testing.T) {
	// Save original env vars
	originalDiscordToken := os.Getenv("DISCORD_TOKEN")
	originalGeminiKey := os.Getenv("GEMINI_API_KEY")
	
	// Clean up after test
	defer func() {
		os.Setenv("DISCORD_TOKEN", originalDiscordToken)
		os.Setenv("GEMINI_API_KEY", originalGeminiKey)
	}()

	// Set test env vars
	os.Setenv("DISCORD_TOKEN", "test-token")
	os.Setenv("GEMINI_API_KEY", "test-key")

	cfg := Load()

	if cfg.DiscordToken != "test-token" {
		t.Errorf("Expected DiscordToken to be 'test-token', got '%s'", cfg.DiscordToken)
	}
	if cfg.GeminiAPIKey != "test-key" {
		t.Errorf("Expected GeminiAPIKey to be 'test-key', got '%s'", cfg.GeminiAPIKey)
	}
}