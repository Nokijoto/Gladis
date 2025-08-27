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
