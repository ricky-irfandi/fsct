package ai

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Client is the main AI client with retry and timeout support
type Client struct {
	provider Provider
	config   *ClientConfig
}

// ClientConfig contains configuration for the AI client
type ClientConfig struct {
	// Timeout is the maximum time to wait for a request
	Timeout time.Duration

	// MaxRetries is the maximum number of retries
	MaxRetries int

	// RetryDelay is the initial delay between retries
	RetryDelay time.Duration

	// MaxRetryDelay is the maximum delay between retries
	MaxRetryDelay time.Duration

	// RetryBackoff is the exponential backoff multiplier
	RetryBackoff float64

	// Verbose enables detailed logging
	Verbose bool
}

// DefaultClientConfig returns the default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryDelay:    1 * time.Second,
		MaxRetryDelay: 30 * time.Second,
		RetryBackoff:  2.0,
		Verbose:       false,
	}
}

// NewClient creates a new AI client
func NewClient(provider Provider, config *ClientConfig) *Client {
	if config == nil {
		config = DefaultClientConfig()
	}
	return &Client{
		provider: provider,
		config:   config,
	}
}

// NewClientFromConfig creates a client from config values
func NewClientFromConfig(apiKey, providerName, baseURL, model string, config *ClientConfig) (*Client, error) {
	factory := NewProviderFactory(apiKey, baseURL, model)
	provider, err := factory.Create(providerName)
	if err != nil {
		return nil, err
	}
	return NewClient(provider, config), nil
}

// Complete sends a completion request with retry logic
func (c *Client) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	// Create a timeout context if not provided
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		// Check if context is canceled
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		if attempt > 0 && c.config.Verbose {
			fmt.Printf("AI client: retry attempt %d/%d\n", attempt, c.config.MaxRetries)
		}

		// Make the request
		resp, err := c.provider.Complete(ctx, req)
		if err == nil {
			if c.config.Verbose {
				fmt.Printf("AI client: request successful (tokens: %d)\n", resp.Usage.TotalTokens)
			}
			return resp, nil
		}

		lastErr = err

		// Check if we should retry
		if attempt < c.config.MaxRetries && IsRetryableError(err) {
			delay := c.calculateDelay(attempt)
			if c.config.Verbose {
				fmt.Printf("AI client: request failed (retryable), waiting %v: %v\n", delay, err)
			}

			select {
			case <-time.After(delay):
				// Continue to next retry
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		} else {
			// Not retryable or no more retries
			if c.config.Verbose {
				fmt.Printf("AI client: request failed (not retryable): %v\n", err)
			}
			return nil, err
		}
	}

	// All retries exhausted
	return nil, fmt.Errorf("all %d attempts failed: %w", c.config.MaxRetries+1, lastErr)
}

// calculateDelay calculates the delay for a retry attempt using exponential backoff
func (c *Client) calculateDelay(attempt int) time.Duration {
	delay := float64(c.config.RetryDelay) * math.Pow(c.config.RetryBackoff, float64(attempt))
	if delay > float64(c.config.MaxRetryDelay) {
		delay = float64(c.config.MaxRetryDelay)
	}
	return time.Duration(delay)
}

// Provider returns the underlying provider
func (c *Client) Provider() Provider {
	return c.provider
}

// SetProvider changes the provider
func (c *Client) SetProvider(provider Provider) {
	c.provider = provider
}

// IsAvailable checks if the AI client is properly configured and available
func (c *Client) IsAvailable() bool {
	return c.provider != nil
}

// CompleteSimple is a simplified version for basic use cases
func (c *Client) CompleteSimple(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	req := &CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    2000,
		Temperature:  0.7,
	}

	resp, err := c.Complete(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}
