package reviewer

import (
	"os"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/config"
	"github.com/ricky-irfandi/fsct/internal/report"
)

// NoCredentialsCheck checks if reviewer credentials are configured (REV-001)
type NoCredentialsCheck struct {
	config *config.ReviewerConfig
}

// NewNoCredentialsCheck creates a new REV-001 check
func NewNoCredentialsCheck(cfg *config.ReviewerConfig) *NoCredentialsCheck {
	return &NoCredentialsCheck{config: cfg}
}

// ID returns the check ID
func (c *NoCredentialsCheck) ID() string {
	return "REV-001"
}

// Name returns the check name
func (c *NoCredentialsCheck) Name() string {
	return "No Reviewer Credentials Configured"
}

// Run executes the check
func (c *NoCredentialsCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if c.config == nil || !c.hasValidCredentials() {
		findings = append(findings, report.Finding{
			ID:       c.ID(),
			Severity: report.SeverityWarning,
			Title:    c.Name(),
			Message: "No reviewer test account is configured. App reviewers may reject the submission " +
				"if they cannot test login functionality or access premium features.",
			Suggestion: "Add reviewer account configuration to .fsct.yaml:\n" +
				"  reviewer:\n" +
				"    email_env: REVIEWER_EMAIL\n" +
				"    password_env: REVIEWER_PASSWORD\n" +
				"Or set environment variables: REVIEWER_EMAIL and REVIEWER_PASSWORD",
		})
	}

	return findings
}

func (c *NoCredentialsCheck) hasValidCredentials() bool {
	if c.config == nil {
		return false
	}

	// Check email
	email := c.getEmail()
	if email == "" {
		return false
	}

	// Check password
	password := c.getPassword()
	return password != ""
}

func (c *NoCredentialsCheck) getEmail() string {
	if c.config.Email != "" {
		return c.config.Email
	}
	if c.config.EmailEnv != "" {
		return os.Getenv(c.config.EmailEnv)
	}
	return ""
}

func (c *NoCredentialsCheck) getPassword() string {
	if c.config.Password != "" {
		return c.config.Password
	}
	if c.config.PasswordEnv != "" {
		return os.Getenv(c.config.PasswordEnv)
	}
	return ""
}

// PlaceholderEmailCheck checks for placeholder email addresses (REV-002)
type PlaceholderEmailCheck struct {
	config *config.ReviewerConfig
}

// NewPlaceholderEmailCheck creates a new REV-002 check
func NewPlaceholderEmailCheck(cfg *config.ReviewerConfig) *PlaceholderEmailCheck {
	return &PlaceholderEmailCheck{config: cfg}
}

// ID returns the check ID
func (c *PlaceholderEmailCheck) ID() string {
	return "REV-002"
}

// Name returns the check name
func (c *PlaceholderEmailCheck) Name() string {
	return "Placeholder Reviewer Email Detected"
}

// Run executes the check
func (c *PlaceholderEmailCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if c.config == nil {
		return findings
	}

	email := c.getEmail()
	if email == "" {
		return findings
	}

	// Check for common placeholder patterns
	placeholders := []string{
		"test@",
		"example.com",
		"your-email@",
		"user@example",
		"email@example",
		"test@example.com",
		"admin@example",
		"demo@",
		"sample@",
	}

	emailLower := strings.ToLower(email)
	for _, pattern := range placeholders {
		if strings.Contains(emailLower, pattern) {
			findings = append(findings, report.Finding{
				ID:       c.ID(),
				Severity: report.SeverityHigh,
				Title:    c.Name(),
				Message:  "The reviewer email '" + maskEmail(email) + "' appears to be a placeholder.",
				Suggestion: "Use a real email address for reviewer testing. " +
					"Create a dedicated test account (e.g., reviewer@yourcompany.com).",
			})
			break
		}
	}

	return findings
}

func (c *PlaceholderEmailCheck) getEmail() string {
	if c.config.Email != "" {
		return c.config.Email
	}
	if c.config.EmailEnv != "" {
		return os.Getenv(c.config.EmailEnv)
	}
	return ""
}

// WeakPasswordCheck checks for weak/common passwords (REV-003)
type WeakPasswordCheck struct {
	config *config.ReviewerConfig
}

// NewWeakPasswordCheck creates a new REV-003 check
func NewWeakPasswordCheck(cfg *config.ReviewerConfig) *WeakPasswordCheck {
	return &WeakPasswordCheck{config: cfg}
}

// ID returns the check ID
func (c *WeakPasswordCheck) ID() string {
	return "REV-003"
}

// Name returns the check name
func (c *WeakPasswordCheck) Name() string {
	return "Weak Reviewer Password Detected"
}

// Run executes the check
func (c *WeakPasswordCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if c.config == nil {
		return findings
	}

	password := c.getPassword()
	if password == "" {
		return findings
	}

	// Check password length
	if len(password) < 8 {
		findings = append(findings, report.Finding{
			ID:       c.ID(),
			Severity: report.SeverityHigh,
			Title:    "Reviewer Password Too Short",
			Message:  "The reviewer password is less than 8 characters long.",
			Suggestion: "Use a stronger password with at least 8 characters, " +
				"including uppercase, lowercase, numbers, and special characters.",
		})
		return findings
	}

	// Check for common weak patterns
	weakPatterns := []string{
		"password", "123456", "qwerty", "admin", "test",
		"flutter", "app", "reviewer", "login", "user",
		"abc123", "welcome", "monkey", "dragon", "master",
	}

	passwordLower := strings.ToLower(password)
	for _, pattern := range weakPatterns {
		if strings.Contains(passwordLower, pattern) {
			findings = append(findings, report.Finding{
				ID:       c.ID(),
				Severity: report.SeverityHigh,
				Title:    c.Name(),
				Message:  "The reviewer password contains a common weak pattern.",
				Suggestion: "Use a unique, strong password that doesn't contain common words or patterns. " +
					"Consider using a password generator.",
			})
			return findings
		}
	}

	// Check for repeated characters
	if hasRepeatedChars(password, 3) {
		findings = append(findings, report.Finding{
			ID:       c.ID(),
			Severity: report.SeverityWarning,
			Title:    "Reviewer Password Has Repeated Characters",
			Message:  "The reviewer password contains repeated characters.",
			Suggestion: "Avoid using repeated characters in passwords for better security.",
		})
	}

	return findings
}

func (c *WeakPasswordCheck) getPassword() string {
	if c.config.Password != "" {
		return c.config.Password
	}
	if c.config.PasswordEnv != "" {
		return os.Getenv(c.config.PasswordEnv)
	}
	return ""
}

// Helper functions

func hasRepeatedChars(s string, threshold int) bool {
	if len(s) < threshold {
		return false
	}

	count := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			count++
			if count >= threshold {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}

func maskEmail(email string) string {
	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}

	local := parts[0]
	domain := parts[1]

	// Mask most of the local part
	if len(local) <= 2 {
		return "***@" + domain
	}

	return local[:2] + "***@" + domain
}

// getEnv is a helper to get environment variable (extracted for testability)
var getEnv = os.Getenv

// LoadConfigFromFile loads reviewer config from .fsct.yaml if it exists
func LoadConfigFromFile(path string) *config.ReviewerConfig {
	// This is a simplified version - in production, you'd properly parse the YAML
	// For now, we'll return nil and rely on environment variables
	return nil
}

// GetConfigFromEnv loads reviewer config from environment variables
func GetConfigFromEnv() *config.ReviewerConfig {
	cfg := &config.ReviewerConfig{}

	// Try to get from common env vars
	if email := os.Getenv("REVIEWER_EMAIL"); email != "" {
		cfg.Email = email
	}
	if password := os.Getenv("REVIEWER_PASSWORD"); password != "" {
		cfg.Password = password
	}

	// If not set directly, check for env var names
	if cfg.Email == "" && os.Getenv("REVIEWER_EMAIL_ENV") != "" {
		cfg.EmailEnv = os.Getenv("REVIEWER_EMAIL_ENV")
		cfg.Email = os.Getenv(cfg.EmailEnv)
	}
	if cfg.Password == "" && os.Getenv("REVIEWER_PASSWORD_ENV") != "" {
		cfg.PasswordEnv = os.Getenv("REVIEWER_PASSWORD_ENV")
		cfg.Password = os.Getenv(cfg.PasswordEnv)
	}

	// Default to standard env var names if nothing else
	if cfg.Email == "" && cfg.EmailEnv == "" {
		cfg.EmailEnv = "REVIEWER_EMAIL"
	}
	if cfg.Password == "" && cfg.PasswordEnv == "" {
		cfg.PasswordEnv = "REVIEWER_PASSWORD"
	}

	return cfg
}
