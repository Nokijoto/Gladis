package commands

import (
	"fmt"
	"os"
	"strings"

	"gladis/internal/database"
	"log" // Dodano import log

	"github.com/bwmarrin/discordgo"
)

// AvailableProvidersCommand defines the /availableproviders command.
var AvailableProvidersCommand = &discordgo.ApplicationCommand{
	Name:        "availableproviders",
	Description: "Wyświetla dostępnych dostawców modeli AI i status ich kluczy API.",
}

// AvailableProvidersHandler handles the /availableproviders command.
func AvailableProvidersHandler(s *discordgo.Session, i *discordgo.InteractionCreate, mongoDB *database.MongoDB) {
	providers, err := mongoDB.GetUniqueProviders()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Wystąpił błąd podczas pobierania dostawców: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	log.Printf("Pobrani dostawcy z bazy danych: %v", providers) // Dodano logowanie

	if len(providers) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Nie znaleziono żadnych dostawców modeli w bazie danych.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var responseBuilder strings.Builder
	responseBuilder.WriteString("Dostępni dostawcy modeli AI:\n")

	for _, provider := range providers {
		envVarName := strings.ToUpper(provider) + "_API_KEY"
		apiKey := os.Getenv(envVarName)

		if apiKey != "" {
			responseBuilder.WriteString(fmt.Sprintf("- **%s**: Dostępny (klucz API znaleziony w .env)\n", provider))
		} else {
			responseBuilder.WriteString(fmt.Sprintf("- **%s**: Brak klucza API (ustaw zmienną środowiskową `%s`)\n", provider, envVarName))
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseBuilder.String(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
