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

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// OpenAIRequest represents the request body for OpenAI API
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// OpenAIMessage represents a message in the conversation
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response from OpenAI API
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// SetAPIKey sets the API key
func (p *OpenAIProvider) SetAPIKey(key string) {
	p.apiKey = key
}

// SetModel sets the model
func (p *OpenAIProvider) SetModel(model string) {
	p.model = model
}

// AvailableModels returns available models
func (p *OpenAIProvider) AvailableModels() []string {
	return []string{
		"gpt-4",
		"gpt-4-turbo",
		"gpt-4o",
		"gpt-3.5-turbo",
	}
}

// Complete sends a completion request to OpenAI
func (p *OpenAIProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if p.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	// Build request body
	model := p.model
	if req.Model != "" {
		model = req.Model
	}
	if model == "" {
		model = "gpt-4"
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2000
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	payload := OpenAIRequest{
		Model:       model,
		Messages:    []OpenAIMessage{},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	// Add system message if provided
	if req.SystemPrompt != "" {
		payload.Messages = append(payload.Messages, OpenAIMessage{
			Role:    "system",
			Content: req.SystemPrompt,
		})
	}

	// Add user message
	payload.Messages = append(payload.Messages, OpenAIMessage{
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

	// Parse response (even for error status codes)
	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResponse, err)
	}

	// Check for API error in response body
	if openaiResp.Error != nil {
		return nil, NewAPIError(p.Name(), resp.StatusCode, openaiResp.Error.Message)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, NewAPIError(p.Name(), resp.StatusCode, string(body))
	}

	// Extract content from response
	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("%w: no choices in response", ErrInvalidResponse)
	}

	choice := openaiResp.Choices[0]
	return &CompletionResponse{
		Content:      choice.Message.Content,
		FinishReason: choice.FinishReason,
		Model:        openaiResp.Model,
		Usage: TokenUsage{
			PromptTokens:     openaiResp.Usage.PromptTokens,
			CompletionTokens: openaiResp.Usage.CompletionTokens,
			TotalTokens:      openaiResp.Usage.TotalTokens,
		},
	}, nil
}
