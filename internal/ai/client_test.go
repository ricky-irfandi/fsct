package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	provider := &MockProvider{}
	client := NewClient(provider, nil)

	if client == nil {
		t.Fatal("expected client to be created")
	}

	if client.provider != provider {
		t.Error("expected provider to be set")
	}

	if client.config == nil {
		t.Error("expected default config to be set")
	}
}

func TestClientComplete(t *testing.T) {
	mockProvider := &MockProvider{
		Response: &CompletionResponse{
			Content:      "Test response",
			FinishReason: "stop",
			Model:        "test-model",
			Usage: TokenUsage{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		},
	}

	client := NewClient(mockProvider, &ClientConfig{
		Timeout:    5 * time.Second,
		MaxRetries: 0,
	})

	req := &CompletionRequest{
		SystemPrompt: "You are a test assistant",
		UserPrompt:   "Hello",
		MaxTokens:    100,
		Temperature:  0.5,
	}

	resp, err := client.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Content != "Test response" {
		t.Errorf("expected 'Test response', got '%s'", resp.Content)
	}

	if resp.Model != "test-model" {
		t.Errorf("expected 'test-model', got '%s'", resp.Model)
	}

	if resp.Usage.TotalTokens != 15 {
		t.Errorf("expected 15 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestClientCompleteWithRetry(t *testing.T) {
	// Use a retryable error (APIError with 500 status)
	mockProvider := &MockProvider{
		FailCount: 2,
		Error:     NewAPIError("mock", 500, "server error"),
		Response: &CompletionResponse{
			Content: "Success after retries",
		},
	}

	client := NewClient(mockProvider, &ClientConfig{
		Timeout:       5 * time.Second,
		MaxRetries:    3,
		RetryDelay:    10 * time.Millisecond,
		RetryBackoff:  1.0, // No backoff for faster tests
		MaxRetryDelay: 100 * time.Millisecond,
	})

	req := &CompletionRequest{
		UserPrompt: "Test",
	}

	resp, err := client.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Content != "Success after retries" {
		t.Errorf("expected 'Success after retries', got '%s'", resp.Content)
	}

	if mockProvider.CallCount != 3 {
		t.Errorf("expected 3 calls (1 original + 2 retries), got %d", mockProvider.CallCount)
	}
}

func TestClientCompleteTimeout(t *testing.T) {
	mockProvider := &MockProvider{
		Delay: 100 * time.Millisecond,
		Response: &CompletionResponse{
			Content: "Delayed response",
		},
	}

	client := NewClient(mockProvider, &ClientConfig{
		Timeout:    50 * time.Millisecond,
		MaxRetries: 0,
	})

	req := &CompletionRequest{
		UserPrompt: "Test",
	}

	// Use client timeout (not context)
	_, err := client.Complete(context.Background(), req)
	// This test verifies the timeout is applied - the delay exceeds timeout
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestOpenAIProvider(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected /chat/completions, got %s", r.URL.Path)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-api-key" {
			t.Errorf("expected 'Bearer test-api-key', got '%s'", authHeader)
		}

		// Parse request body
		var req OpenAIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Model != "gpt-4" {
			t.Errorf("expected model 'gpt-4', got '%s'", req.Model)
		}

		// Send response
		resp := OpenAIResponse{
			ID:     "test-id",
			Model:  "gpt-4",
			Object: "chat.completion",
			Choices: []struct {
				Index        int           `json:"index"`
				Message      OpenAIMessage `json:"message"`
				FinishReason string        `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: OpenAIMessage{
						Role:    "assistant",
						Content: "This is a test response",
					},
					FinishReason: "stop",
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     10,
				CompletionTokens: 5,
				TotalTokens:      15,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create provider
	provider := &OpenAIProvider{
		baseURL: server.URL,
		model:   "gpt-4",
	}
	provider.SetAPIKey("test-api-key")

	// Make request
	req := &CompletionRequest{
		SystemPrompt: "You are a helpful assistant",
		UserPrompt:   "Hello",
		MaxTokens:    100,
		Temperature:  0.7,
	}

	resp, err := provider.Complete(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Content != "This is a test response" {
		t.Errorf("expected 'This is a test response', got '%s'", resp.Content)
	}

	if resp.Model != "gpt-4" {
		t.Errorf("expected 'gpt-4', got '%s'", resp.Model)
	}
}

func TestConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg.Provider != "minimax" {
		t.Errorf("expected default provider 'minimax', got '%s'", cfg.Provider)
	}

	if cfg.APIKeyEnv != "AI_API_KEY" {
		t.Errorf("expected default API key env 'AI_API_KEY', got '%s'", cfg.APIKeyEnv)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid with api key",
			config: &Config{
				Enabled:  true,
				Provider: "openai",
				APIKey:   "test-key",
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: &Config{
				Enabled:  true,
				Provider: "openai",
			},
			wantErr: true,
		},
		{
			name: "offline mode",
			config: &Config{
				Enabled:  true,
				Offline:  true,
				Provider: "openai",
			},
			wantErr: false,
		},
		{
			name: "custom without baseurl",
			config: &Config{
				Enabled:  true,
				Provider: "custom",
				APIKey:   "test-key",
			},
			wantErr: true,
		},
		{
			name: "invalid provider",
			config: &Config{
				Enabled:  true,
				Provider: "invalid",
				APIKey:   "test-key",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// MockProvider is a mock implementation for testing
type MockProvider struct {
	Response  *CompletionResponse
	Error     error
	FailCount int
	CallCount int
	Delay     time.Duration
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) SetAPIKey(key string) {}

func (m *MockProvider) SetModel(model string) {}

func (m *MockProvider) AvailableModels() []string {
	return []string{"mock-model"}
}

func (m *MockProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	m.CallCount++

	if m.Delay > 0 {
		select {
		case <-time.After(m.Delay):
			// Continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Return error for the first FailCount calls
	if m.CallCount <= m.FailCount {
		if m.Error != nil {
			return nil, m.Error
		}
		return nil, ErrRequestFailed
	}

	// Return response on successful calls
	return m.Response, nil
}
