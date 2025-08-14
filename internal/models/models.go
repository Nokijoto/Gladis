package models

import (
	"strings"
)

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

func GetAvailableModels() []string {
	return []string{
		"gemini-2.5-pro",
		"gemini-2.5-flash",
		"gemini-2.5-flash-lite",
		"google/gemini-2.5-pro-exp-03-25",
		"google/gemini-2.0-flash-exp:free",
		"qwen/qwen3-coder:free",
		"tngtech/deepseek-r1t2-chimera:free",
		"deepseek/deepseek-r1-0528:free",
		"tngtech/deepseek-r1t-chimera:free",
		"microsoft/mai-ds-r1:free",
		"deepseek/deepseek-chat-v3-0324:free",
		"deepseek/deepseek-r1:free",
		"openai/gpt-oss-20b:free",
		"z-ai/glm-4.5-air:free",
		"mistralai/mistral-small-3.2-24b-instruct:free",
		"moonshotai/kimi-dev-72b:free",
		"deepseek/deepseek-r1-0528-qwen3-8b:free",
		"qwen/qwen3-235b-a22b:free",
		"moonshotai/kimi-vl-a3b-thinking:free",
		"nvidia/llama-3.1-nemotron-ultra-253b-v1:free",
		"nousresearch/deephermes-3-llama-3-8b-preview:free",
		"meta-llama/llama-3.2-3b-instruct:free",
		"meta-llama/llama-3.2-11b-vision-instruct:free",
		"mistralai/mistral-nemo:free",
		"mistralai/mistral-small-3.1-24b-instruct:free",
		"agentica-org/deepcoder-14b-preview:free",
		"google/gemma-3-12b-it:free",
		"google/gemma-3-27b-it:free",
		"meta-llama/llama-3.3-70b-instruct:free",
		"meta-llama/llama-3.1-405b-instruct:free",
		"deepseek/deepseek-r1-distill-qwen-14b:free",
		"qwen/qwen3-4b:free",
		"qwen/qwen3-30b-a3b:free",
		"qwen/qwen3-8b:free",
		"qwen/qwen3-14b:free",
		"moonshotai/kimi-k2:free",
		"cognitivecomputations/dolphin-mistral-24b-venice-edition:free",
		"tencent/hunyuan-a13b-instruct:free",
		"sarvamai/sarvam-m:free",
		"mistralai/devstral-small-2505:free",
		"shisa-ai/shisa-v2-llama3.3-70b:free",
		"arliai/qwq-32b-arliai-rpr-v1:free",
		"featherless/qwerky-72b:free",
		"google/gemma-3-4b-it:free",
		"rekaai/reka-flash-3:free",
		"qwen/qwq-32b:free",
		"cognitivecomputations/dolphin3.0-r1-mistral-24b:free",
		"cognitivecomputations/dolphin3.0-mistral-24b:free",
		"qwen/qwen2.5-vl-72b-instruct:free",
		"mistralai/mistral-small-24b-instruct-2501:free",
		"qwen/qwen-2.5-coder-32b-instruct:free",
		"qwen/qwen-2.5-72b-instruct:free",
		"mistralai/mistral-7b-instruct:free",
		"google/gemma-3n-e2b-it:free",
		"google/gemma-3n-e4b-it:free",
		"qwen/qwen2.5-vl-32b-instruct:free",
		"deepseek/deepseek-r1-distill-llama-70b:free",
		"google/gemma-2-9b-it:free",
	}
}

func ParseModelString(modelStr string) (provider AIProvider, modelName string) {
	if strings.HasPrefix(modelStr, "openrouter:") {
		return ProviderOpenRouter, strings.TrimPrefix(modelStr, "openrouter:")
	}
	return ProviderGemini, modelStr
}
