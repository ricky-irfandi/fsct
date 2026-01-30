package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from .fsct.yaml file if it exists
func LoadConfig(projectPath string) *Config {
	cfg := &Config{}

	configPath := filepath.Join(projectPath, ".fsct.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg
	}

	_ = yaml.Unmarshal(data, cfg)
	return cfg
}

type Config struct {
	AI        *AIConfig        `yaml:"ai,omitempty"`
	Reviewer  *ReviewerConfig  `yaml:"reviewer,omitempty"`
	Checks    *ChecksConfig    `yaml:"checks,omitempty"`
	Platforms *PlatformsConfig `yaml:"platforms,omitempty"`
}

type AIConfig struct {
	Enabled   bool   `yaml:"enabled"`
	URL       string `yaml:"url"`
	Model     string `yaml:"model"`
	APIKey    string `yaml:"api_key,omitempty"`
	APIKeyEnv string `yaml:"api_key_env"`
	Offline   bool   `yaml:"offline"`
}

type ReviewerConfig struct {
	Email        string              `yaml:"email"`
	EmailEnv     string              `yaml:"email_env"`
	Password     string              `yaml:"password,omitempty"`
	PasswordEnv  string              `yaml:"password_env"`
	Verification *VerificationConfig `yaml:"verification,omitempty"`
}

type VerificationConfig struct {
	Enabled          bool   `yaml:"enabled"`
	AuthEndpoint     string `yaml:"auth_endpoint"`
	Method           string `yaml:"method"`
	BodyTemplate     string `yaml:"body_template,omitempty"`
	SuccessIndicator string `yaml:"success_indicator"`
}

type ChecksConfig struct {
	Skip    []string `yaml:"skip"`
	Include []string `yaml:"include"`
}

type PlatformsConfig struct {
	Android bool `yaml:"android"`
	IOS     bool `yaml:"ios"`
}
