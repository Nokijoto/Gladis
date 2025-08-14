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
	currentModel string
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
		currentModel: "gemini-1.5-flash",
	}, nil
}

func (m *Manager) SetModel(model string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	availableModels := models.GetAvailableModels()
	found := false
	for _, available := range availableModels {
		if available == model {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("model %s is not available", model)
	}

	provider, _ := models.ParseModelString(model)

	switch provider {
	case models.ProviderGemini:
		if m.geminiClient == nil {
			return fmt.Errorf("Gemini API key not configured")
		}
	case models.ProviderOpenRouter:
		if m.openRouter == nil {
			return fmt.Errorf("OpenRouter API key not configured")
		}
	}

	m.currentModel = model
	return nil
}

func (m *Manager) GetCurrentModel() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentModel
}

func (m *Manager) GenerateContent(ctx context.Context, prompt string, images [][]byte) (string, error) {
	m.mu.RLock()
	currentModel := m.currentModel
	m.mu.RUnlock()

	provider, modelName := models.ParseModelString(currentModel)

	switch provider {
	case models.ProviderGemini:
		if m.geminiClient == nil {
			return "", fmt.Errorf("Gemini client not initialized")
		}
		return m.geminiClient.GenerateContent(ctx, prompt, images, modelName)
	case models.ProviderOpenRouter:
		if m.openRouter == nil {
			return "", fmt.Errorf("OpenRouter client not initialized")
		}
		return m.openRouter.GenerateContent(ctx, prompt, images, modelName)
	default:
		return "", fmt.Errorf("unknown provider: %s", provider)
	}
}

func (m *Manager) Close() error {
	if m.geminiClient != nil {
		return m.geminiClient.Close()
	}
	return nil
}
