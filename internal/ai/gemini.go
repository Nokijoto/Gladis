package ai

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"strings"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &GeminiClient{client: client}, nil
}

func (gc *GeminiClient) GenerateContent(ctx context.Context, prompt string, images [][]byte, modelName string) (string, error) {
	model := gc.client.GenerativeModel(modelName)

	var parts []genai.Part
	parts = append(parts, genai.Text(prompt))

	for _, img := range images {
		parts = append(parts, genai.ImageData("jpeg", img))
	}

	resp, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("no response generated")
	}

	var response strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			response.WriteString(string(text))
		}
	}

	return response.String(), nil
}

func (gc *GeminiClient) Close() error {
	return gc.client.Close()
}
