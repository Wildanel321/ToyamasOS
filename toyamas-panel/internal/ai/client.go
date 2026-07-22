package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type OllamaClient struct {
	BaseURL    string
	Model      string
	httpClient *http.Client
}

func NewOllamaClient() *OllamaClient {
	baseURL := os.Getenv("OLLAMA_HOST")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "qwen:latest"
	}

	return &OllamaClient{
		BaseURL: baseURL,
		Model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *OllamaClient) IsAvailable() bool {
	resp, err := c.httpClient.Get(c.BaseURL + "/api/tags")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()
	return true
}

type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Stream bool   `json:"stream"`
}

type OllamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *OllamaClient) Generate(systemPrompt, userPrompt string) (string, error) {
	if !c.IsAvailable() {
		return "", fmt.Errorf("ollama API is unreachable at %s. Ensure Ollama container is running", c.BaseURL)
	}

	reqPayload := OllamaGenerateRequest{
		Model:  c.Model,
		Prompt: userPrompt,
		System: systemPrompt,
		Stream: false,
	}

	data, err := json.Marshal(reqPayload)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %s", resp.Status)
	}

	var res OllamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Response, nil
}
