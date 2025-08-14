package ai

import (
	"context"
	"fmt"
	"sync"

	"gladis/internal/config"
	"gladis/internal/models"
)

type Manager struct {
	config       *config.Config
	geminiClient *GeminiClient
	openRouter   *OpenRouterClient
	currentModel models.ModelInfo
	mu           sync.RWMutex
}

func NewManager(cfg *config.Config) (*Manager, error) {
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
		return nil, fmt.Errorf("at least one API key (Gemini or OpenRouter) must be provided")
	}

	return &Manager{
		config:       cfg,
		geminiClient: geminiClient,
		openRouter:   openRouterClient,
		currentModel: models.ModelInfo{Name: "gemini-2.5-flash-lite", Provider: models.ProviderGemini},
	}, nil
}

func (m *Manager) SetModel(modelName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	availableModels := models.GetAvailableModels()
	var selectedModel models.ModelInfo
	found := false
	for _, available := range availableModels {
		if available.Name == modelName {
			selectedModel = available
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("model %s is not available", modelName)
	}

	switch selectedModel.Provider {
	case models.ProviderGemini:
		if m.geminiClient == nil {
			return fmt.Errorf("Gemini API key not configured")
		}
	case models.ProviderOpenRouter:
		if m.openRouter == nil {
			return fmt.Errorf("OpenRouter API key not configured")
		}
	}

	m.currentModel = selectedModel
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
