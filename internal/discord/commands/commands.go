package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Commands is a slice of all application commands.
// It will be populated by aggregating commands from individual files.
var Commands = []*discordgo.ApplicationCommand{
	PingCommand,
	HelpCommand,
	SetModelCommand,
	InfoCommand,
	SetSystemPromptCommand,
	SetContextCommand,
	AvailableProvidersCommand,
}
