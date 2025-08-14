package discord

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gladis/internal/models"
)

const modelsPerPage = 10

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Responds with PONG!",
	},
	{
		Name:        "help",
		Description: "Displays this help message.",
	},
	{
		Name:        "setmodel",
		Description: "Displays a list of available models to choose from, with pagination.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "page",
				Description: "The page number to display (e.g., 1 for the first 25 models).",
				Required:    false,
			},
		},
	},
	{
		Name:        "info",
		Description: "Displays information about the bot, including the current model.",
	},
	{
		Name:        "setsystemprompt",
		Description: "Sets a system prompt to be sent as the first message in every interaction.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The system prompt to set.",
				Required:    true,
			},
		},
	},
	{
		Name:        "setcontext",
		Description: "Sets the number of previous messages to include as context. Max 10.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "count",
				Description: "Number of messages to include from 0 to 10. Set to 0 to disable.",
				Required:    true,
				// MinValue i MaxValue pomijamy celowo
			},
		},
	},
}

type CommandHandler struct {
	currentModel  models.ModelInfo
	systemPrompt  string
	contextLength int
	aiManager     AIModelManager // Add AIModelManager interface
}

// AIModelManager interface to decouple CommandHandler from concrete AI manager
type AIModelManager interface {
	SetModel(modelName string) error
	GetCurrentModel() models.ModelInfo
}

func NewCommandHandler(aiManager AIModelManager) *CommandHandler {
	return &CommandHandler{
		aiManager:     aiManager,
		currentModel:  aiManager.GetCurrentModel(),
		systemPrompt:  "",
		contextLength: 0,
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func (h *CommandHandler) HandlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "PONG!",
		},
	})
}

func (h *CommandHandler) HandleSetContext(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var newContextLength int
	if opt, ok := optionMap["count"]; ok {
		newContextLength = int(opt.IntValue())
	}

	newContextLength = clampInt(newContextLength, 0, 10)
	h.contextLength = newContextLength

	resp := "Context messages disabled."
	if newContextLength > 0 {
		resp = fmt.Sprintf("Context length set to: `%d` messages.", newContextLength)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: resp,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *CommandHandler) HandleInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	prompt := h.systemPrompt
	if prompt == "" {
		prompt = "no system prompt set"
	}

	h.currentModel = h.aiManager.GetCurrentModel() // Ensure current model is up-to-date

	embed := &discordgo.MessageEmbed{
		Title:       "Bot Information",
		Description: "Here is some information about the bot.",
		Color:       0x0099FF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Current Model",
				Value:  fmt.Sprintf("`%s` (Provider: `%s`)", h.currentModel.Name, h.currentModel.Provider),
				Inline: true,
			},
			{
				Name:   "System Prompt",
				Value:  fmt.Sprintf("`%s`", prompt),
				Inline: false,
			},
			{
				Name:   "Context Length",
				Value:  fmt.Sprintf("`%d` messages", h.contextLength),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Powered by Gladis",
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func (h *CommandHandler) HandleSetSystemPrompt(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var newPrompt string
	if opt, ok := optionMap["prompt"]; ok {
		newPrompt = opt.StringValue()
	}

	h.systemPrompt = newPrompt

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("System prompt set to: `%s`", newPrompt),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *CommandHandler) HandleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	helpMessage := "Available commands:\n"
	for _, cmd := range Commands {
		helpMessage += fmt.Sprintf("`/%s` - %s\n", cmd.Name, cmd.Description)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMessage,
		},
	})
}

func (h *CommandHandler) HandleSetModel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	page := 0 // Default to page 0 (first page)
	if opt, ok := optionMap["page"]; ok {
		// Discord command options are 1-indexed for users, so convert to 0-indexed
		page = int(opt.IntValue()) - 1
		if page < 0 {
			page = 0
		}
	}
	h.sendModelSelect(s, i, page)
}

func (h *CommandHandler) sendModelSelect(s *discordgo.Session, i *discordgo.InteractionCreate, page int) {
	modelsList := models.GetAvailableModels()
	if len(modelsList) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "No Models Found",
			Description: "No models available for content generation.",
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

	// Sort models by name for consistent pagination
	sort.Slice(modelsList, func(k, l int) bool {
		return modelsList[k].Name < modelsList[l].Name
	})

	totalPages := (len(modelsList) + modelsPerPage - 1) / modelsPerPage
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}
	if totalPages == 0 { // Handle case where there are no models
		page = 0
	}

	start := page * modelsPerPage
	end := start + modelsPerPage
	if end > len(modelsList) {
		end = len(modelsList)
	}

	pagedModels := modelsList[start:end]

	var selectOptions []discordgo.SelectMenuOption
	for _, model := range pagedModels {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label:   fmt.Sprintf("%s (%s)", model.Name, model.Provider),
			Value:   model.Name, // Use model.Name as the value
			Default: model.Name == h.aiManager.GetCurrentModel().Name,
		})
	}

	selectMenu := discordgo.SelectMenu{
		CustomID:    fmt.Sprintf("model_select_%d", page), // Include page in CustomID
		Placeholder: fmt.Sprintf("Choose a model (Page %d/%d)...", page+1, totalPages),
		Options:     selectOptions,
	}

	var components []discordgo.MessageComponent
	components = append(components, selectMenu)

	var buttons []discordgo.MessageComponent
	if page > 0 {
		buttons = append(buttons, discordgo.Button{
			Label:    "Previous",
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("model_prev_%d", page-1),
		})
	}
	if page < totalPages-1 {
		buttons = append(buttons, discordgo.Button{
			Label:    "Next",
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("model_next_%d", page+1),
		})
	}

	if len(buttons) > 0 {
		components = append(components, discordgo.ActionsRow{Components: buttons})
	}

	responseContent := fmt.Sprintf("Current model: `%s` (Provider: `%s`)\nPlease select a new model from the dropdown:", h.aiManager.GetCurrentModel().Name, h.aiManager.GetCurrentModel().Provider)

	var interactionType discordgo.InteractionResponseType
	if i.Type == discordgo.InteractionApplicationCommand {
		interactionType = discordgo.InteractionResponseChannelMessageWithSource
	} else {
		interactionType = discordgo.InteractionResponseUpdateMessage
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: interactionType,
		Data: &discordgo.InteractionResponseData{
			Content:    responseContent,
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: components,
		},
	})
}

func (h *CommandHandler) HandleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	if strings.HasPrefix(data.CustomID, "model_select_") {
		if len(data.Values) == 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "No model selected.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		chosenModelName := data.Values[0]
		err := h.aiManager.SetModel(chosenModelName)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content:    fmt.Sprintf("Failed to set model: %s", err.Error()),
					Flags:      discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{},
				},
			})
			return
		}
		h.currentModel = h.aiManager.GetCurrentModel() // Update local state

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    fmt.Sprintf("Model set to: `%s` (Provider: `%s`)", h.currentModel.Name, h.currentModel.Provider),
				Flags:      discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{},
			},
		})
	} else if strings.HasPrefix(data.CustomID, "model_prev_") {
		pageStr := strings.TrimPrefix(data.CustomID, "model_prev_")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			fmt.Printf("Error parsing page number from CustomID: %v\n", err)
			return
		}
		h.sendModelSelect(s, i, page)
	} else if strings.HasPrefix(data.CustomID, "model_next_") {
		pageStr := strings.TrimPrefix(data.CustomID, "model_next_")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			fmt.Printf("Error parsing page number from CustomID: %v\n", err)
			return
		}
		h.sendModelSelect(s, i, page)
	}
}
