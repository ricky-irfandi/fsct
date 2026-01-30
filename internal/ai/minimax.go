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

// MiniMaxProvider implements the Provider interface for MiniMax
type MiniMaxProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

// MiniMaxRequest represents the request body for MiniMax API
type MiniMaxRequest struct {
	Model       string           `json:"model"`
	Messages    []MiniMaxMessage `json:"messages"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	Temperature float64          `json:"temperature,omitempty"`
}

// MiniMaxMessage represents a message in the conversation
type MiniMaxMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MiniMaxResponse represents the response from MiniMax API
type MiniMaxResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Index        int            `json:"index"`
		Message      MiniMaxMessage `json:"message"`
		FinishReason string         `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	BaseResp struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

// Name returns the provider name
func (p *MiniMaxProvider) Name() string {
	return "minimax"
}

// SetAPIKey sets the API key
func (p *MiniMaxProvider) SetAPIKey(key string) {
	p.apiKey = key
}

// SetModel sets the model
func (p *MiniMaxProvider) SetModel(model string) {
	p.model = model
}

// AvailableModels returns available models
func (p *MiniMaxProvider) AvailableModels() []string {
	return []string{
		"abab6.5s-chat",
		"abab6.5-chat",
		"abab5.5s-chat",
		"abab5.5-chat",
	}
}

// Complete sends a completion request to MiniMax
func (p *MiniMaxProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if p.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	// Build request body
	model := p.model
	if req.Model != "" {
		model = req.Model
	}
	if model == "" {
		model = "abab6.5s-chat"
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2000
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	payload := MiniMaxRequest{
		Model:       model,
		Messages:    []MiniMaxMessage{},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	// Add system message if provided
	if req.SystemPrompt != "" {
		payload.Messages = append(payload.Messages, MiniMaxMessage{
			Role:    "system",
			Content: req.SystemPrompt,
		})
	}

	// Add user message
	payload.Messages = append(payload.Messages, MiniMaxMessage{
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

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, NewAPIError(p.Name(), resp.StatusCode, string(body))
	}

	// Parse response
	var minimaxResp MiniMaxResponse
	if err := json.Unmarshal(body, &minimaxResp); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidResponse, err)
	}

	// Check MiniMax-specific error
	if minimaxResp.BaseResp.StatusCode != 0 {
		return nil, NewAPIError(p.Name(), minimaxResp.BaseResp.StatusCode, minimaxResp.BaseResp.StatusMsg)
	}

	// Extract content from response
	if len(minimaxResp.Choices) == 0 {
		return nil, fmt.Errorf("%w: no choices in response", ErrInvalidResponse)
	}

	choice := minimaxResp.Choices[0]
	return &CompletionResponse{
		Content:      choice.Message.Content,
		FinishReason: choice.FinishReason,
		Model:        model,
		Usage: TokenUsage{
			PromptTokens:     minimaxResp.Usage.PromptTokens,
			CompletionTokens: minimaxResp.Usage.CompletionTokens,
			TotalTokens:      minimaxResp.Usage.TotalTokens,
		},
	}, nil
}
