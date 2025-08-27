package ai

import (
	"context"
	"fmt"
	"sync"

	"gladis/internal/config"
	"gladis/internal/database"
	"gladis/internal/models"
)

type Manager struct {
	config       *config.Config
	geminiClient *GeminiClient
	openRouter   *OpenRouterClient
	DB           *database.MongoDB
	currentModel models.ModelInfo
	mu           sync.RWMutex
}

func NewManager(cfg *config.Config, db *database.MongoDB) (*Manager, error) {
	var geminiClient *GeminiClient
	var err error

	if cfg.GeminiAPIKey != "" {
		geminiClient, err = NewGeminiClient(cfg.GeminiAPIKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create Gemini client: %w", err)
		}
	}

	var openRouterClient *OpenRouterClient
	if cfg.OpenRouterAPIKey != "" {
		openRouterClient = NewOpenRouterClient(cfg.OpenRouterAPIKey)
	}

	if geminiClient == nil && openRouterClient == nil {
		return nil, fmt.Errorf("at least one API key must be provided")
	}

	return &Manager{
		config:       cfg,
		geminiClient: geminiClient,
		openRouter:   openRouterClient,
		DB:           db, // KLUCZOWE, wcześniej brakowało
		currentModel: models.ModelInfo{Name: "gemini-2.5-flash-lite", Provider: models.ProviderGemini},
	}, nil
}

func (m *Manager) SetModelByID(idHex string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.DB == nil {
		return fmt.Errorf("database handle is nil")
	}
	doc, err := m.DB.GetModelByID(idHex)
	if err != nil {
		return fmt.Errorf("failed to get model by id: %w", err)
	}
	if doc == nil {
		return fmt.Errorf("model not found")
	}

	selected := models.ModelInfo{
		Name:     doc.Name,
		Provider: models.AIProvider(doc.Provider),
	}

	switch selected.Provider {
	case models.ProviderGemini:
		if m.geminiClient == nil {
			return fmt.Errorf("Gemini API key not configured")
		}
	case models.ProviderOpenRouter:
		if m.openRouter == nil {
			return fmt.Errorf("OpenRouter API key not configured")
		}
	}

	m.currentModel = selected
	return nil
}

func (m *Manager) GetCurrentModel() models.ModelInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentModel
}

func (m *Manager) SetCurrentModel(model models.ModelInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentModel = model
}

func (m *Manager) GenerateContent(ctx context.Context, prompt string, images [][]byte) (string, error) {
	m.mu.RLock()
	currentModel := m.currentModel
	m.mu.RUnlock()

	switch currentModel.Provider {
	case models.ProviderGemini:
		if m.geminiClient == nil {
			return "", fmt.Errorf("Gemini client not initialized")
		}
		return m.geminiClient.GenerateContent(ctx, prompt, images, currentModel.Name)
	case models.ProviderOpenRouter:
		if m.openRouter == nil {
			return "", fmt.Errorf("OpenRouter client not initialized")
		}
		return m.openRouter.GenerateContent(ctx, prompt, images, currentModel.Name)
	default:
		return "", fmt.Errorf("unknown provider: %s", currentModel.Provider)
	}
}

func (m *Manager) Close() error {
	if m.geminiClient != nil {
		return m.geminiClient.Close()
	}
	return nil
}
