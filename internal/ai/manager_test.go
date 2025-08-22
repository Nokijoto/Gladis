package ai

import (
	"testing"

	"gladis/internal/config"
	"gladis/internal/models"
)

func TestManagerCreation(t *testing.T) {
	// Test with valid Gemini config
	cfg := &config.Config{
		GeminiAPIKey: "test-gemini-key",
	}

	// Note: This will fail because we can't create real clients without valid API keys
	// But we can test the validation logic
	_, err := NewManager(cfg)
	if err == nil {
		t.Log("Manager creation succeeded (this may fail in CI without real API keys)")
	} else {
		t.Logf("Manager creation failed as expected without real API keys: %v", err)
	}

	// Test with no API keys - should definitely fail
	emptyCfg := &config.Config{}
	_, err = NewManager(emptyCfg)
	if err == nil {
		t.Error("Expected error when creating manager with no API keys")
	}
}

func TestManagerModelOperations(t *testing.T) {
	// Create a mock manager for testing
	manager := &Manager{
		currentModel: models.ModelInfo{
			Name:     "test-model",
			Provider: models.ProviderGemini,
		},
	}

	// Test GetCurrentModel
	current := manager.GetCurrentModel()
	if current.Name != "test-model" {
		t.Errorf("Expected current model name 'test-model', got '%s'", current.Name)
	}

	// Test SetCurrentModel
	newModel := models.ModelInfo{
		Name:     "new-test-model",
		Provider: models.ProviderOpenRouter,
	}
	manager.SetCurrentModel(newModel)

	updated := manager.GetCurrentModel()
	if updated.Name != "new-test-model" {
		t.Errorf("Expected updated model name 'new-test-model', got '%s'", updated.Name)
	}
	if updated.Provider != models.ProviderOpenRouter {
		t.Errorf("Expected updated provider '%s', got '%s'", models.ProviderOpenRouter, updated.Provider)
	}
}

func TestSetModelValidation(t *testing.T) {
	// Create a mock manager
	manager := &Manager{}

	// Test setting an invalid model
	err := manager.SetModel("invalid-model-name")
	if err == nil {
		t.Error("Expected error when setting invalid model")
	}

	// Test setting a valid model name
	validModels := models.GetAllModels()
	if len(validModels) > 0 {
		err = manager.SetModel(validModels[0].Name)
		// This might fail because the clients aren't initialized, but the model should be found
		if err != nil && err.Error() != "Gemini API key not configured" && err.Error() != "OpenRouter API key not configured" {
			t.Errorf("Unexpected error when setting valid model: %v", err)
		}
	}
}

func TestManagerClose(t *testing.T) {
	// Test closing manager with no clients
	manager := &Manager{}
	err := manager.Close()
	if err != nil {
		t.Errorf("Expected no error when closing empty manager, got %v", err)
	}
}
