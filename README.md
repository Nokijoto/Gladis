# Gladis Discord Bot

A versatile Discord bot that integrates with multiple AI providers including Google Gemini and OpenRouter. The bot can process text and images, and supports model switching through Discord commands.

## Features

- **Multi-Provider Support**: Works with Google Gemini and OpenRouter
- **Image Processing**: Can analyze and respond to images sent in messages
- **Model Switching**: Change AI models on-the-fly using Discord commands
- **Slash Commands**: Modern Discord slash command interface
- **Mention-Based**: Responds when mentioned in messages
- **Docker Support**: Easy deployment with Docker

## Supported Models

### Google Gemini
- gemini-2.5-pro
- gemini-2.5-flash
- gemini-2.5-flash-lite
- gemini-1.5-pro
- gemini-1.5-flash

### OpenRouter Models
- anthropic/claude-3.5-sonnet
- anthropic/claude-3.5-haiku
- openai/gpt-4o
- openai/gpt-4o-mini
- google/gemini-2.5-pro-exp-03-25
- google/gemini-2.0-flash-001
- mistralai/mistral-large
- meta-llama/llama-3.1-405b-instruct

## Setup

### Environment Variables

Create a `.env` file with the following variables:

```bash
DISCORD_TOKEN=your_discord_bot_token
GEMINI_API_KEY=your_google_gemini_api_key
OPENROUTER_API_KEY=your_openrouter_api_key
ENVIRONMENT=production
```

### Running Locally

```bash
go mod tidy
go run cmd/bot/main.go
```

### Docker Deployment

```bash
docker build -t gladis-bot .
docker run -d --env-file .env gladis-bot
```

## Usage

1. Invite the bot to your Discord server
2. Mention the bot in any channel: `@Gladis your question here`
3. Attach images for image analysis
4. Use `/setmodel` to change the AI model

## Commands

- `/ping` - Test bot responsiveness
- `/help` - Display available commands
- `/setmodel` - Change the active AI model

## Development

The codebase is organized into modular packages:

- `internal/ai/` - AI provider integrations
- `internal/discord/` - Discord bot functionality
- `internal/config/` - Configuration management
- `internal/models/` - Model definitions and utilities
- `cmd/bot/` - Main application entry point
