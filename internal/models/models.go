package models
const (
	MAX_IMAGE_SIZE_MB      = 2
	MAX_IMAGES_PER_MESSAGE = 5
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

func GetModelsPage1() []ModelInfo {
	return []ModelInfo{
		{Name: "gemini-2.5-pro", Provider: ProviderGemini},
		{Name: "gemini-2.5-flash", Provider: ProviderGemini},
		{Name: "gemini-2.5-flash-lite", Provider: ProviderGemini},
		{Name: "google/gemini-2.5-pro-exp-03-25", Provider: ProviderOpenRouter},
		{Name: "google/gemini-2.0-flash-exp:free", Provider: ProviderOpenRouter},
		{Name: "mistralai/mistral-small-3.2-24b-instruct:free", Provider: ProviderOpenRouter},
		{Name: "meta-llama/llama-3.3-70b-instruct:free", Provider: ProviderOpenRouter},
		{Name: "google/gemma-3-12b-it:free", Provider: ProviderOpenRouter},
		{Name: "google/gemma-3-27b-it:free", Provider: ProviderOpenRouter},
		{Name: "meta-llama/llama-3.1-405b-instruct:free", Provider: ProviderOpenRouter},
	}
}

func GetModelsPage2() []ModelInfo {
	return []ModelInfo{
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
	}
}

func GetModelsPage3() []ModelInfo {
	return []ModelInfo{
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
	}
}

func GetModelsPage4() []ModelInfo {
	return []ModelInfo{
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
	}
}

func GetModelsPage5() []ModelInfo {
	return []ModelInfo{
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
	}
}

func GetModelsPage6() []ModelInfo {
	return []ModelInfo{
		{Name: "mistralai/mistral-7b-instruct:free", Provider: ProviderOpenRouter},
		{Name: "google/gemma-3n-e2b-it:free", Provider: ProviderOpenRouter},
		{Name: "google/gemma-3n-e4b-it:free", Provider: ProviderOpenRouter},
		{Name: "qwen/qwen2.5-vl-32b-instruct:free", Provider: ProviderOpenRouter},
		{Name: "deepseek/deepseek-r1-distill-llama-70b:free", Provider: ProviderOpenRouter},
		{Name: "google/gemma-2-9b-it:free", Provider: ProviderOpenRouter},
	}
}

// GetAllModels returns a concatenated list of all models from all pages.
func GetAllModels() []ModelInfo {
	allModels := []ModelInfo{}
	allModels = append(allModels, GetModelsPage1()...)
	allModels = append(allModels, GetModelsPage2()...)
	allModels = append(allModels, GetModelsPage3()...)
	allModels = append(allModels, GetModelsPage4()...)
	allModels = append(allModels, GetModelsPage5()...)
	allModels = append(allModels, GetModelsPage6()...)
	return allModels
}

func GetAvailableModels(page int, modelsPerPage int) ([]ModelInfo, int) {
	allModels := GetAllModels()
	totalModels := len(allModels)
	totalPages := (totalModels + modelsPerPage - 1) / modelsPerPage

	// Ensure page is within valid bounds
	if page < 0 {
		page = 0
	}
	if page >= totalPages && totalPages > 0 {
		page = totalPages - 1
	}
	if totalPages == 0 { // Safety check for no models
		return []ModelInfo{}, 0
	}

	startIndex := page * modelsPerPage
	endIndex := startIndex + modelsPerPage

	if startIndex >= totalModels {
		return []ModelInfo{}, totalPages
	}
	if endIndex > totalModels {
		endIndex = totalModels
	}

	return allModels[startIndex:endIndex], totalPages
}
