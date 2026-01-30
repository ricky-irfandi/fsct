package ai

import (
	"context"
	"fmt"
	"strings"
)

// Provider defines the interface for AI providers
type Provider interface {
	// Name returns the provider name
	Name() string

	// Complete sends a completion request and returns the response
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// SetAPIKey sets the API key for the provider
	SetAPIKey(key string)

	// SetModel sets the model to use
	SetModel(model string)

	// AvailableModels returns the list of available models
	AvailableModels() []string
}

// CompletionRequest represents a request to the AI API
type CompletionRequest struct {
	// SystemPrompt is the system message/prompt
	SystemPrompt string

	// UserPrompt is the user message/prompt
	UserPrompt string

	// MaxTokens is the maximum number of tokens to generate
	MaxTokens int

	// Temperature controls randomness (0.0 - 2.0)
	Temperature float64

	// Model overrides the default model
	Model string
}

// CompletionResponse represents a response from the AI API
type CompletionResponse struct {
	// Content is the generated text
	Content string

	// FinishReason indicates why the generation stopped
	FinishReason string

	// Usage contains token usage information
	Usage TokenUsage

	// Model is the model used for generation
	Model string
}

// TokenUsage represents token usage information
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// ProviderFactory creates providers
type ProviderFactory struct {
	apiKey  string
	baseURL string
	model   string
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory(apiKey, baseURL, model string) *ProviderFactory {
	return &ProviderFactory{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
	}
}

// Create creates a provider based on the provider name
func (f *ProviderFactory) Create(providerName string) (Provider, error) {
	switch strings.ToLower(providerName) {
	case "minimax":
		return f.createMiniMaxProvider()
	case "openai":
		return f.createOpenAIProvider()
	case "custom", "":
		return f.createCustomProvider()
	default:
		return nil, fmt.Errorf("%w: %s", ErrProviderNotFound, providerName)
	}
}

// createMiniMaxProvider creates a MiniMax provider
func (f *ProviderFactory) createMiniMaxProvider() (Provider, error) {
	provider := &MiniMaxProvider{
		baseURL: "https://api.minimax.chat/v1",
		model:   "abab6.5s-chat",
	}
	if f.apiKey != "" {
		provider.SetAPIKey(f.apiKey)
	}
	if f.model != "" {
		provider.SetModel(f.model)
	}
	return provider, nil
}

// createOpenAIProvider creates an OpenAI provider
func (f *ProviderFactory) createOpenAIProvider() (Provider, error) {
	provider := &OpenAIProvider{
		baseURL: "https://api.openai.com/v1",
		model:   "gpt-4",
	}
	if f.apiKey != "" {
		provider.SetAPIKey(f.apiKey)
	}
	if f.model != "" {
		provider.SetModel(f.model)
	}
	if f.baseURL != "" {
		provider.baseURL = f.baseURL
	}
	return provider, nil
}

// createCustomProvider creates a custom OpenAI-compatible provider
func (f *ProviderFactory) createCustomProvider() (Provider, error) {
	if f.baseURL == "" {
		return nil, NewValidationError("baseURL", "custom provider requires baseURL")
	}
	provider := &CustomProvider{
		baseURL: f.baseURL,
		model:   f.model,
	}
	if f.apiKey != "" {
		provider.SetAPIKey(f.apiKey)
	}
	return provider, nil
}

// DetectProvider attempts to detect the provider from the URL
func DetectProvider(url string) string {
	url = strings.ToLower(url)
	switch {
	case strings.Contains(url, "minimax"):
		return "minimax"
	case strings.Contains(url, "openai"):
		return "openai"
	default:
		return "custom"
	}
}
