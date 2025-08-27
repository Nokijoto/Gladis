package commands

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gladis/internal/ai"

	"github.com/bwmarrin/discordgo"
)

// Pomocnicze: zwięzłe ID wskaźnika
func dbgID(v any) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%p", v)
}

// Buduje łańcuch debug z kontekstem
func debugCtx(where string, aiManager *ai.Manager, page int, extra string) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("where=%s ", where))
	sb.WriteString(fmt.Sprintf("ai=%s ", dbgID(aiManager)))
	if aiManager != nil {
		sb.WriteString(fmt.Sprintf("db=%s ", dbgID(aiManager.DB)))
	}
	if page > 0 {
		sb.WriteString(fmt.Sprintf("page=%d ", page))
	}
	if extra != "" {
		sb.WriteString(extra)
	}
	sb.WriteString(fmt.Sprintf(" ts=%s", time.Now().Format(time.RFC3339)))
	return sb.String()
}

// Czy pokazać debug użytkownikowi
func debugToUser() bool {
	return os.Getenv("DEBUG_DB") == "1"
}

var SetModelCommand = func() *discordgo.ApplicationCommand {
	min := 1.0
	return &discordgo.ApplicationCommand{
		Name:        "setmodel",
		Description: "Displays a list of available models to choose from.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "page",
				Description: "Page number to display",
				Required:    false,
				MinValue:    &min,
			},
		},
	}
}()

// HandleSetModel obsługuje komendę i pierwsze renderowanie
func HandleSetModel(s *discordgo.Session, i *discordgo.InteractionCreate, aiManager *ai.Manager) {
	fmt.Printf("[SetModel] %s\n", debugCtx("slash", aiManager, 0, ""))

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	page := 1
	if opt, ok := optionMap["page"]; ok {
		page = int(opt.IntValue())
		if page < 1 {
			page = 1
		}
	}

	sendModelSelect(s, i, page, aiManager, false)
}

// sendModelSelect renderuje dropdown i przyciski paginacji.
// update true oznacza edycję istniejącej wiadomości.
func sendModelSelect(s *discordgo.Session, i *discordgo.InteractionCreate, page int, aiManager *ai.Manager, update bool) {
	fmt.Printf("[sendModelSelect] %s\n", debugCtx("render", aiManager, page, fmt.Sprintf("update=%t", update)))

	// 1. Szybkie ACK żeby nie złapać 10062
	if update {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
	} else {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}

	// 2. Walidacje po ACK. Błędy pokażemy przez edit.
	if aiManager == nil {
		editErrorEmbed(s, i, "❌ Error", "AI Manager nie jest zainicjalizowany.", nil)
		return
	}
	if aiManager.DB == nil {
		editErrorEmbed(s, i, "❌ Error", "Połączenie z bazą danych nie jest dostępne.", nil)
		return
	}

	const pageSize = int64(10)

	// 3. Pobierz stronę modeli z DB
	rows, err := aiManager.DB.GetModels(int64(page), pageSize)
	if err != nil {
		if debugToUser() {
			editErrorEmbed(s, i, "❌ Error", "Nie udało się pobrać listy modeli z bazy danych.",
				[]*discordgo.MessageEmbedField{{Name: "Debug", Value: fmt.Sprintf("%v", err), Inline: false}})
		} else {
			editErrorEmbed(s, i, "❌ Error", "Nie udało się pobrać listy modeli z bazy danych.", nil)
		}
		return
	}
	if len(rows) == 0 {
		editErrorEmbed(s, i, "❌ No Models Found", fmt.Sprintf("Brak modeli dla strony %d.", page), nil)
		return
	}

	// 4. Stabilne sortowanie w obrębie strony
	sort.Slice(rows, func(a, b int) bool {
		if rows[a].Name == rows[b].Name {
			return rows[a].ID.Hex() < rows[b].ID.Hex()
		}
		return rows[a].Name < rows[b].Name
	})

	// 5. Zbuduj komponenty UI
	current := aiManager.GetCurrentModel()

	var selectOptions []discordgo.SelectMenuOption
	for _, m := range rows {
		idHex := m.ID.Hex()
		isDefault := m.Name == current.Name && m.Provider == string(current.Provider)
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label:   fmt.Sprintf("%s (%s)", m.Name, m.Provider),
			Value:   idHex, // unikalne ID
			Default: isDefault,
		})
	}

	selectMenu := discordgo.SelectMenu{
		CustomID:    "model_select",
		Placeholder: "Choose a model...",
		Options:     selectOptions,
	}
	rowSelect := discordgo.ActionsRow{Components: []discordgo.MessageComponent{selectMenu}}

	var buttons []discordgo.MessageComponent
	if page > 1 {
		buttons = append(buttons, discordgo.Button{
			Label:    "Previous Page",
			CustomID: fmt.Sprintf("page-%d", page-1),
			Style:    discordgo.PrimaryButton,
		})
	}
	if len(rows) == int(pageSize) {
		buttons = append(buttons, discordgo.Button{
			Label:    "Next Page",
			CustomID: fmt.Sprintf("page-%d", page+1),
			Style:    discordgo.PrimaryButton,
		})
	}
	components := []discordgo.MessageComponent{rowSelect}
	if len(buttons) > 0 {
		components = append(components, discordgo.ActionsRow{Components: buttons})
	}

	responseContent := fmt.Sprintf(
		"Current model: `%s` Provider: `%s`\nWybierz nowy model z listy poniżej.",
		current.Name, current.Provider,
	)

	embed := &discordgo.MessageEmbed{
		Title:       "📚 Available Models",
		Description: fmt.Sprintf("Strona %d. Wybierz model z rozwijanej listy.", page),
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Current Model**",
				Value:  fmt.Sprintf("`%s` Provider: `%s`", current.Name, current.Provider),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Model selection",
		},
	}
	if debugToUser() {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  "Debug",
			Value: debugCtx("ok", aiManager, page, fmt.Sprintf("pageSize=%d listLen=%d", pageSize, len(rows))),
		})
	}

	// 6. Edytujemy odpowiedź po ACK
	edit := &discordgo.WebhookEdit{
		Content:    strPtr(responseContent),
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &components,
	}
	if _, err := s.InteractionResponseEdit(i.Interaction, edit); err != nil {
		fmt.Printf("[sendModelSelect] edit error: %v ctx=%s\n", err, debugCtx("edit-err", aiManager, page, ""))
	}
}

// proste helpery

func editErrorEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, title, desc string, extra []*discordgo.MessageEmbedField) {
	emb := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       0xFF0000,
	}
	if extra != nil {
		emb.Fields = append(emb.Fields, extra...)
	}
	emptyComponents := []discordgo.MessageComponent{}
	edit := &discordgo.WebhookEdit{
		Content:    strPtr(""),
		Embeds:     &[]*discordgo.MessageEmbed{emb},
		Components: &emptyComponents,
	}
	_, _ = s.InteractionResponseEdit(i.Interaction, edit)
}

func strPtr(s string) *string { return &s }

// HandleInteraction obsługuje kliknięcia przycisków i wybór w dropdownie
func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, aiManager *ai.Manager) {
	fmt.Printf("[Interaction] %s\n", debugCtx("component", aiManager, 0, fmt.Sprintf("customID=%q", i.MessageComponentData().CustomID)))

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID

	// Paginacja
	if strings.HasPrefix(customID, "page-") {
		pageStr := customID[len("page-"):]
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			// szybki ACK i mały błąd
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
			editErrorEmbed(s, i, "❌ Error", "Invalid page number.", nil)
			return
		}
		sendModelSelect(s, i, page, aiManager, true)
		return
	}

	// Wybór modelu po ID z Mongo
	if customID == "model_select" {
		values := i.MessageComponentData().Values
		if len(values) == 0 {
			return
		}
		idHex := values[0]

		// ACK komponentu
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		if aiManager == nil {
			editErrorEmbed(s, i, "❌ Error", "AI Manager nie jest zainicjalizowany.", nil)
			return
		}

		if err := aiManager.SetModelByID(idHex); err != nil {
			editErrorEmbed(s, i, "❌ Error", fmt.Sprintf("Error: %v", err), nil)
			return
		}

		// Po ustawieniu modelu odświeżamy widok na pierwszą stronę
		sendModelSelect(s, i, 1, aiManager, true)
		return
	}
}
