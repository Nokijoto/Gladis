// Package models defines the data structures for AI models and providers
// supported by the Gladis Discord bot. It provides functionality for
// managing and paginating through available AI models.
package models

const (
	MAX_IMAGE_SIZE_MB      = 2
	MAX_IMAGES_PER_MESSAGE = 5
	MODELS_PER_PAGE        = 10 // Define the pagination size as a constant
)

type AIProvider string

const (
	ProviderGemini     AIProvider = "gemini"
	ProviderOpenRouter AIProvider = "openrouter"
)

type ModelInfo struct {
	Name     string
	Provider AIProvider
}

// allModels contains all available models
var allModels = []ModelInfo{
	// Gemini models
	{Name: "gemini-2.5-pro", Provider: ProviderGemini},
	{Name: "gemini-2.5-flash", Provider: ProviderGemini},
	{Name: "gemini-2.5-flash-lite", Provider: ProviderGemini},

	// OpenRouter models
	{Name: "google/gemini-2.5-pro-exp-03-25", Provider: ProviderOpenRouter},
	{Name: "google/gemini-2.0-flash-exp:free", Provider: ProviderOpenRouter},
	{Name: "mistralai/mistral-small-3.2-24b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "meta-llama/llama-3.3-70b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "google/gemma-3-12b-it:free", Provider: ProviderOpenRouter},
	{Name: "google/gemma-3-27b-it:free", Provider: ProviderOpenRouter},
	{Name: "meta-llama/llama-3.1-405b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen3-8b:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen3-14b:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen3-30b-a3b:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen3-coder:free", Provider: ProviderOpenRouter},
	{Name: "mistralai/mistral-nemo:free", Provider: ProviderOpenRouter},
	{Name: "mistralai/mistral-small-3.1-24b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "agentica-org/deepcoder-14b-preview:free", Provider: ProviderOpenRouter},
	{Name: "moonshotai/kimi-dev-72b:free", Provider: ProviderOpenRouter},
	{Name: "deepseek/deepseek-r1-0528:free", Provider: ProviderOpenRouter},
	{Name: "tngtech/deepseek-r1t2-chimera:free", Provider: ProviderOpenRouter},
	{Name: "tngtech/deepseek-r1t-chimera:free", Provider: ProviderOpenRouter},
	{Name: "microsoft/mai-ds-r1:free", Provider: ProviderOpenRouter},
	{Name: "deepseek/deepseek-chat-v3-0324:free", Provider: ProviderOpenRouter},
	{Name: "deepseek/deepseek-r1:free", Provider: ProviderOpenRouter},
	{Name: "openai/gpt-oss-20b:free", Provider: ProviderOpenRouter},
	{Name: "z-ai/glm-4.5-air:free", Provider: ProviderOpenRouter},
	{Name: "moonshotai/kimi-vl-a3b-thinking:free", Provider: ProviderOpenRouter},
	{Name: "nvidia/llama-3.1-nemotron-ultra-253b-v1:free", Provider: ProviderOpenRouter},
	{Name: "nousresearch/deephermes-3-llama-3-8b-preview:free", Provider: ProviderOpenRouter},
	{Name: "meta-llama/llama-3.2-3b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "meta-llama/llama-3.2-11b-vision-instruct:free", Provider: ProviderOpenRouter},
	{Name: "deepseek/deepseek-r1-distill-qwen-14b:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen3-4b:free", Provider: ProviderOpenRouter},
	{Name: "moonshotai/kimi-k2:free", Provider: ProviderOpenRouter},
	{Name: "cognitivecomputations/dolphin-mistral-24b-venice-edition:free", Provider: ProviderOpenRouter},
	{Name: "tencent/hunyuan-a13b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "sarvamai/sarvam-m:free", Provider: ProviderOpenRouter},
	{Name: "mistralai/devstral-small-2505:free", Provider: ProviderOpenRouter},
	{Name: "shisa-ai/shisa-v2-llama3.3-70b:free", Provider: ProviderOpenRouter},
	{Name: "arliai/qwq-32b-arliai-rpr-v1:free", Provider: ProviderOpenRouter},
	{Name: "featherless/qwerky-72b:free", Provider: ProviderOpenRouter},
	{Name: "google/gemma-3-4b-it:free", Provider: ProviderOpenRouter},
	{Name: "rekaai/reka-flash-3:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwq-32b:free", Provider: ProviderOpenRouter},
	{Name: "cognitivecomputations/dolphin3.0-r1-mistral-24b:free", Provider: ProviderOpenRouter},
	{Name: "cognitivecomputations/dolphin3.0-mistral-24b:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen2.5-vl-72b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "mistralai/mistral-small-24b-instruct-2501:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen-2.5-coder-32b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen-2.5-72b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "mistralai/mistral-7b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "google/gemma-3n-e2b-it:free", Provider: ProviderOpenRouter},
	{Name: "google/gemma-3n-e4b-it:free", Provider: ProviderOpenRouter},
	{Name: "qwen/qwen2.5-vl-32b-instruct:free", Provider: ProviderOpenRouter},
	{Name: "deepseek/deepseek-r1-distill-llama-70b:free", Provider: ProviderOpenRouter},
	{Name: "google/gemma-2-9b-it:free", Provider: ProviderOpenRouter},
}

// GetAllModels returns all available models
func GetAllModels() []ModelInfo {
	return allModels
}

// GetModelsPage returns models for a specific page
func GetModelsPage(page int) []ModelInfo {
	startIndex := page * MODELS_PER_PAGE
	endIndex := startIndex + MODELS_PER_PAGE

	if startIndex >= len(allModels) {
		return []ModelInfo{}
	}

	if endIndex > len(allModels) {
		endIndex = len(allModels)
	}

	return allModels[startIndex:endIndex]
}

// GetTotalPages returns the total number of pages
func GetTotalPages() int {
	return (len(allModels) + MODELS_PER_PAGE - 1) / MODELS_PER_PAGE
}

// Legacy functions for backward compatibility - these can be removed in a future version
func GetModelsPage1() []ModelInfo { return GetModelsPage(0) }
func GetModelsPage2() []ModelInfo { return GetModelsPage(1) }
func GetModelsPage3() []ModelInfo { return GetModelsPage(2) }
func GetModelsPage4() []ModelInfo { return GetModelsPage(3) }
func GetModelsPage5() []ModelInfo { return GetModelsPage(4) }
func GetModelsPage6() []ModelInfo { return GetModelsPage(5) }
