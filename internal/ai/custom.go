package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CustomProvider implements a generic OpenAI-compatible provider
type CustomProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// CustomRequest represents a generic OpenAI-compatible request
type CustomRequest struct {
	Model       string          `json:"model"`
	Messages    []CustomMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// CustomMessage represents a message
type CustomMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CustomResponse represents a generic OpenAI-compatible response
type CustomResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      CustomMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// Name returns the provider name
func (p *CustomProvider) Name() string {
	return "custom"
}

// SetAPIKey sets the API key
func (p *CustomProvider) SetAPIKey(key string) {
	p.apiKey = key
}

// SetModel sets the model
func (p *CustomProvider) SetModel(model string) {
	p.model = model
}

// AvailableModels returns available models
// For custom providers, this is empty as we don't know what models are available
func (p *CustomProvider) AvailableModels() []string {
	return []string{p.model}
}

// Complete sends a completion request to the custom endpoint
func (p *CustomProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if p.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	if p.baseURL == "" {
		return nil, NewValidationError("baseURL", "custom provider requires baseURL")
	}

	// Build request body
	model := p.model
	if req.Model != "" {
		model = req.Model
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2000
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	payload := CustomRequest{
		Model:       model,
		Messages:    []CustomMessage{},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	// Add system message if provided
	if req.SystemPrompt != "" {
		payload.Messages = append(payload.Messages, CustomMessage{
			Role:    "system",
			Content: req.SystemPrompt,
		})
	}

	// Add user message
	payload.Messages = append(payload.Messages, CustomMessage{
		Role:    "user",
		Content: req.UserPrompt,
	})

	// Serialize request
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/chat/completions", p.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Make request
	client := p.httpClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.Canceled {
			return nil, ErrContextCanceled
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, ErrAPITimeout
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var customResp CustomResponse
	if err := json.Unmarshal(body, &customResp); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResponse, err)
	}

	// Check for API error in response body
	if customResp.Error != nil {
		return nil, NewAPIError(p.Name(), resp.StatusCode, customResp.Error.Message)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, NewAPIError(p.Name(), resp.StatusCode, string(body))
	}

	// Extract content from response
	if len(customResp.Choices) == 0 {
		return nil, fmt.Errorf("%w: no choices in response", ErrInvalidResponse)
	}

	choice := customResp.Choices[0]
	modelUsed := customResp.Model
	if modelUsed == "" {
		modelUsed = model
	}

	return &CompletionResponse{
		Content:      choice.Message.Content,
		FinishReason: choice.FinishReason,
		Model:        modelUsed,
		Usage: TokenUsage{
			PromptTokens:     customResp.Usage.PromptTokens,
			CompletionTokens: customResp.Usage.CompletionTokens,
			TotalTokens:      customResp.Usage.TotalTokens,
		},
	}, nil
}
