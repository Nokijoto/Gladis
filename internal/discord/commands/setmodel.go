package commands

import (
	"fmt"
	"sort"
	"strconv" // Added import
	"strings" // Added import

	"gladis/internal/ai"
	"gladis/internal/models"

	"github.com/bwmarrin/discordgo"
)

var SetModelCommand = &discordgo.ApplicationCommand{
	Name:        "setmodel",
	Description: "Displays a list of available models to choose from.",
}

// HandleSetModel handles the /setmodel command and implements pagination.
func HandleSetModel(s *discordgo.Session, i *discordgo.InteractionCreate, aiManager *ai.Manager) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Determine which page to display
	page := 1
	if opt, ok := optionMap["page"]; ok {
		page = int(opt.IntValue()) // Parse page number from interaction
	}

	sendModelSelect(s, i, page, aiManager)
}

// sendModelSelect displays the model selection dropdown with pagination.
func sendModelSelect(s *discordgo.Session, i *discordgo.InteractionCreate, page int, aiManager *ai.Manager) {
	// Get all models and paginate
	allModels := models.GetAllModels()

	// Calculate the start and end indices for the page
	startIndex := (page - 1) * 25
	endIndex := startIndex + 25
	if endIndex > len(allModels) {
		endIndex = len(allModels)
	}

	// Slice models to only show the selected page
	modelsList := allModels[startIndex:endIndex]

	// Handle empty model list
	if len(modelsList) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "❌ No Models Found",
			Description: fmt.Sprintf("No models found for page %d.", page),
			Color:       0xFF0000,
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		return
	}

	// Sort models by name for consistent order
	sort.Slice(modelsList, func(k, l int) bool {
		return modelsList[k].Name < modelsList[l].Name
	})

	// Create select options
	var selectOptions []discordgo.SelectMenuOption
	currentModel := aiManager.GetCurrentModel()
	for _, model := range modelsList {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label:   fmt.Sprintf("%s (%s)", model.Name, model.Provider),
			Value:   model.Name, // Use model.Name as the value
			Default: model.Name == currentModel.Name,
		})
	}

	// Create select menu
	selectMenu := discordgo.SelectMenu{
		CustomID:    "model_select",
		Placeholder: "Choose a model...",
		Options:     selectOptions,
	}

	// Wrap selectMenu in ActionsRow component (required by Discord)
	actionsRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// Prepare response content
	responseContent := fmt.Sprintf("Current model: `%s` (Provider: `%s`)\nPlease select a new model from the dropdown.", currentModel.Name, currentModel.Provider)

	// Prepare pagination buttons
	var components []discordgo.MessageComponent
	if page > 1 {
		previousButton := discordgo.Button{
			Label:    "Previous Page",
			CustomID: fmt.Sprintf("page-%d", page-1), // Set page number for previous
			Style:    discordgo.PrimaryButton,
		}
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{previousButton},
		})
	}

	if endIndex < len(allModels) {
		nextButton := discordgo.Button{
			Label:    "Next Page",
			CustomID: fmt.Sprintf("page-%d", page+1), // Set page number for next
			Style:    discordgo.PrimaryButton,
		}
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{nextButton},
		})
	}

	// Create Embed for the model selection response
	embed := &discordgo.MessageEmbed{
		Title:       "📚 Available Models",
		Description: "Select a model from the dropdown below:",
		Color:       0x00FF00, // Green for success
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Current Model**",
				Value:  fmt.Sprintf("`%s` (Provider: `%s`)", currentModel.Name, currentModel.Provider),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Model selection",
		},
	}

	// Respond to the interaction with the Embed, select menu, and pagination buttons
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    responseContent,
			Flags:      discordgo.MessageFlagsEphemeral,
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: append([]discordgo.MessageComponent{actionsRow}, components...),
		},
	})

	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}

// HandleInteraction handles message component interactions (e.g., button clicks, select menus).
func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, aiManager *ai.Manager) {
	if i.Type == discordgo.InteractionMessageComponent { // Corrected check for message component interaction
		customID := i.MessageComponentData().CustomID

		if strings.HasPrefix(customID, "page-") {
			// Handle pagination button click
			pageStr := customID[len("page-"):]
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				fmt.Println("Error parsing page number from custom ID:", err)
				// Respond with an error message to the user
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Invalid page number.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}
			sendModelSelect(s, i, page, aiManager)
		} else if customID == "model_select" {
			// Handle model selection from the dropdown
			if len(i.MessageComponentData().Values) > 0 {
				selectedModelName := i.MessageComponentData().Values[0]
				currentModel := aiManager.GetCurrentModel()

				if selectedModelName != currentModel.Name {
					allModels := models.GetAllModels()
					var newModel *models.ModelInfo // Corrected to models.ModelInfo
					for _, m := range allModels {
						if m.Name == selectedModelName {
							newModel = &m
							break
						}
					}

					if newModel != nil {
						aiManager.SetModel(newModel.Name) // Corrected to pass model name string
						responseContent := fmt.Sprintf("Model updated to: `%s` (Provider: `%s`)", newModel.Name, newModel.Provider)
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: responseContent,
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Error: Selected model not found.",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})
					}
				} else {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You have already selected this model.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}
			}
		}
	}
}
