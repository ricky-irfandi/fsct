package ai

import (
	"fmt"
	"os"
	"strings"
)

// Config holds AI configuration
type Config struct {
	// Enabled enables/disables AI integration
	Enabled bool

	// Provider is the AI provider name (minimax, openai, custom)
	Provider string

	// APIKey is the API key for the provider
	APIKey string

	// APIKeyEnv is the environment variable name for the API key
	APIKeyEnv string

	// BaseURL is the base URL for the API (optional, for custom providers)
	BaseURL string

	// Model is the model to use
	Model string

	// Timeout is the request timeout
	Timeout string

	// MaxRetries is the maximum number of retries
	MaxRetries int

	// Offline mode disables AI even if configured
	Offline bool
}

// NewConfig returns a default configuration
func NewConfig() *Config {
	return &Config{
		Enabled:    true,
		Provider:   "minimax",
		APIKeyEnv:  "AI_API_KEY",
		Model:      "",
		Timeout:    "30s",
		MaxRetries: 3,
		Offline:    false,
	}
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	// Check for API key in environment
	if c.APIKeyEnv != "" {
		if key := os.Getenv(c.APIKeyEnv); key != "" {
			c.APIKey = key
		}
	}

	// Check common env vars if still no API key
	if c.APIKey == "" {
		for _, envVar := range []string{"AI_API_KEY", "OPENAI_API_KEY", "MINIMAX_API_KEY"} {
			if key := os.Getenv(envVar); key != "" {
				c.APIKey = key
				break
			}
		}
	}

	// Provider from env
	if provider := os.Getenv("AI_PROVIDER"); provider != "" {
		c.Provider = provider
	}

	// Base URL from env
	if baseURL := os.Getenv("AI_BASE_URL"); baseURL != "" {
		c.BaseURL = baseURL
	}

	// Model from env
	if model := os.Getenv("AI_MODEL"); model != "" {
		c.Model = model
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Offline {
		return nil
	}

	if !c.Enabled {
		return nil
	}

	if c.APIKey == "" {
		return fmt.Errorf("%w: no API key configured (set %s env var)", ErrNoAPIKey, c.APIKeyEnv)
	}

	if c.Provider == "" {
		c.Provider = "minimax"
	}

	c.Provider = strings.ToLower(c.Provider)

	validProviders := []string{"minimax", "openai", "custom"}
	found := false
	for _, p := range validProviders {
		if c.Provider == p {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid provider: %s (must be one of: %v)", c.Provider, validProviders)
	}

	if c.Provider == "custom" && c.BaseURL == "" {
		return NewValidationError("baseURL", "custom provider requires baseURL")
	}

	return nil
}

// IsConfigured returns true if AI is properly configured
func (c *Config) IsConfigured() bool {
	if c.Offline || !c.Enabled {
		return false
	}
	return c.APIKey != ""
}

// GetProviderName returns the provider name with auto-detection
func (c *Config) GetProviderName() string {
	if c.Provider != "" {
		return c.Provider
	}

	// Auto-detect from base URL
	if c.BaseURL != "" {
		return DetectProvider(c.BaseURL)
	}

	return "minimax"
}

// MaskedAPIKey returns the API key masked for display
func (c *Config) MaskedAPIKey() string {
	if c.APIKey == "" {
		return ""
	}

	if len(c.APIKey) <= 8 {
		return "****"
	}

	return c.APIKey[:4] + "****" + c.APIKey[len(c.APIKey)-4:]
}

// NewClient creates a new AI client from this config
func (c *Config) NewClient() (*Client, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	factory := NewProviderFactory(c.APIKey, c.BaseURL, c.Model)
	provider, err := factory.Create(c.GetProviderName())
	if err != nil {
		return nil, err
	}

	clientConfig := DefaultClientConfig()
	if c.MaxRetries > 0 {
		clientConfig.MaxRetries = c.MaxRetries
	}

	return NewClient(provider, clientConfig), nil
}

// LoadConfig loads configuration from multiple sources
// Priority: 1) CLI flags, 2) Environment variables, 3) Config file, 4) Defaults
func LoadConfig(cliAPIKey, cliProvider, cliURL, cliModel string, cliOffline bool) (*Config, error) {
	config := NewConfig()

	// Load from environment
	config.LoadFromEnv()

	// Override with CLI flags (highest priority)
	if cliAPIKey != "" {
		config.APIKey = cliAPIKey
	}
	if cliProvider != "" {
		config.Provider = cliProvider
	}
	if cliURL != "" {
		config.BaseURL = cliURL
	}
	if cliModel != "" {
		config.Model = cliModel
	}
	if cliOffline {
		config.Offline = cliOffline
	}

	return config, nil
}
