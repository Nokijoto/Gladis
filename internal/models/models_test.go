package models

import (
	"testing"
)

func TestGetAllModels(t *testing.T) {
	models := GetAllModels()

	if len(models) == 0 {
		t.Error("Expected at least one model, got none")
	}

	// Check that we have some expected models
	expectedModels := []string{
		"gemini-2.5-pro",
		"gemini-2.5-flash",
		"gemini-2.5-flash-lite",
	}

	modelMap := make(map[string]bool)
	for _, model := range models {
		modelMap[model.Name] = true
	}

	for _, expected := range expectedModels {
		if !modelMap[expected] {
			t.Errorf("Expected model '%s' not found in models list", expected)
		}
	}
}

func TestGetModelsPage(t *testing.T) {
	// Test first page
	page0 := GetModelsPage(0)
	if len(page0) == 0 {
		t.Error("Expected models on page 0, got none")
	}
	if len(page0) > MODELS_PER_PAGE {
		t.Errorf("Expected at most %d models per page, got %d", MODELS_PER_PAGE, len(page0))
	}

	// Test empty page (beyond available models)
	emptyPage := GetModelsPage(999)
	if len(emptyPage) != 0 {
		t.Errorf("Expected empty page for page 999, got %d models", len(emptyPage))
	}
}

func TestGetTotalPages(t *testing.T) {
	totalPages := GetTotalPages()
	totalModels := len(GetAllModels())
	expectedPages := (totalModels + MODELS_PER_PAGE - 1) / MODELS_PER_PAGE

	if totalPages != expectedPages {
		t.Errorf("Expected %d pages, got %d", expectedPages, totalPages)
	}
}

func TestModelInfo(t *testing.T) {
	model := ModelInfo{
		Name:     "test-model",
		Provider: ProviderGemini,
	}

	if model.Name != "test-model" {
		t.Errorf("Expected Name 'test-model', got '%s'", model.Name)
	}
	if model.Provider != ProviderGemini {
		t.Errorf("Expected Provider '%s', got '%s'", ProviderGemini, model.Provider)
	}
}

func TestBackwardCompatibility(t *testing.T) {
	// Test that legacy functions still work
	page1 := GetModelsPage1()
	newPage1 := GetModelsPage(0)

	if len(page1) != len(newPage1) {
		t.Errorf("Legacy GetModelsPage1() returned %d models, but GetModelsPage(0) returned %d", len(page1), len(newPage1))
	}

	// Check that models are the same
	for i, model := range page1 {
		if i >= len(newPage1) || model.Name != newPage1[i].Name {
			t.Errorf("Legacy function compatibility broken for model %d", i)
		}
	}
}
