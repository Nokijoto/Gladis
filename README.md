# Gladis Discord Bot

Gladis is a versatile Discord bot built with Go, designed to integrate various AI models and provide interactive functionalities within your Discord server. It supports multiple AI backends (currently Gemini and OpenRouter) and offers commands for managing AI interactions, setting context, and configuring system prompts.

## Features

- **AI Integration**: Connects to Gemini and OpenRouter APIs for advanced AI capabilities.
- **Configurable AI Models**: Easily switch between different AI models.
- **Context Management**: Control the number of previous messages sent to the AI for better conversational flow.
- **Custom System Prompts**: Set a custom system prompt to guide the AI's behavior.
- **Logging**: Integrates with a logging system, with optional webhook support for external log monitoring.
- **Discord Commands**:
    - `/help`: Displays a list of available commands.
    - `/ping`: Tests the bot's responsiveness.
    - `/setcontext`: Sets the number of previous messages to send to the bot (0-10).
    - `/setmodel`: Interactively changes the AI model being used.
    - `/setsystemprompt`: Sets a custom message that is always sent to the bot.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.21 or higher)
- Discord Bot Token
- API Key for Gemini or OpenRouter (or both)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Nokijoto/Gladis.git
   cd Gladis
   ```

2. **Set up environment variables:**
   Create a `.env` file in the root directory of the project and add the following:

   ```
   DISCORD_TOKEN=YOUR_DISCORD_BOT_TOKEN
   GEMINI_API_KEY=YOUR_GEMINI_API_KEY # Optional, if using Gemini
   OPENROUTER_API_KEY=YOUR_OPENROUTER_API_KEY # Optional, if using OpenRouter
   ENVIRONMENT=development # or production
   WEBHOOK_URL=YOUR_WEBHOOK_URL # Optional, for logging to a webhook
   ```
   Replace `YOUR_DISCORD_BOT_TOKEN`, `YOUR_GEMINI_API_KEY`, `YOUR_OPENROUTER_API_KEY`, and `YOUR_WEBHOOK_URL` with your actual keys and URLs.

3. **Run the bot:**
   ```bash
   go run cmd/bot/main.go
   ```
   The bot should now be running and connected to your Discord server.

## Testing

The project includes comprehensive unit tests for core functionality:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/config
```

## Development

### Code Quality

The project uses standard Go tools for code quality:

```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Build the project
go build ./...
```

## Project Structure

```
.
├── cmd/
│   └── bot/
│       └── main.go           # Main entry point for the Discord bot
├── internal/
│   ├── ai/
│   │   ├── gemini.go         # Gemini AI integration
│   │   ├── manager.go        # Manages AI model interactions
│   │   └── openrouter.go     # OpenRouter AI integration
│   ├── config/
│   │   └── config.go         # Configuration loading and management
│   ├── discord/
│   │   ├── bot.go            # Discord bot core logic
│   │   └── commands/         # Discord command definitions and handlers
│   │       ├── commands.go
│   │       ├── help.go
│   │       ├── info.go
│   │       ├── ping.go
│   │       ├── setcontext.go
│   │       ├── setmodel.go
│   │       └── setsystemprompt.go
│   ├── logger/
│   │   └── logger.go         # Logging utility
│   └── models/
│       └── models.go         # Data models
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
└── README.md                 # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
