package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenRouterClient struct {
	client *http.Client
	apiKey string
}

type OpenRouterMessage struct {
	Role    string                   `json:"role"`
	Content []OpenRouterContentBlock `json:"content"`
}

type OpenRouterContentBlock struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL *struct {
		URL string `json:"url"`
	} `json:"image_url,omitempty"`
}

type OpenRouterRequest struct {
	Model    string              `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		client: &http.Client{
			Timeout: 60 * time.Second, // Add 60 second timeout
		},
		apiKey: apiKey,
	}
}

func (oc *OpenRouterClient) GenerateContent(ctx context.Context, prompt string, images [][]byte, modelName string) (string, error) {
	message := OpenRouterMessage{
		Role: "user",
		Content: []OpenRouterContentBlock{
			{Type: "text", Text: prompt},
		},
	}

	for _, img := range images {
		imgBase64 := base64.StdEncoding.EncodeToString(img)
		imgURL := fmt.Sprintf("data:image/jpeg;base64,%s", imgBase64)
		message.Content = append(message.Content, OpenRouterContentBlock{
			Type: "image_url",
			ImageURL: &struct {
				URL string `json:"url"`
			}{URL: imgURL},
		})
	}

	request := OpenRouterRequest{
		Model:    modelName,
		Messages: []OpenRouterMessage{message},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+oc.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/Nokijoto/Gladis")
	req.Header.Set("X-Title", "Gladis Discord Bot")

	resp, err := oc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if we received any choices in the response
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return response.Choices[0].Message.Content, nil
}
