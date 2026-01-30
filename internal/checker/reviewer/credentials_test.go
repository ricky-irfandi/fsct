package reviewer

import (
	"os"
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/config"
	"github.com/ricky-irfandi/fsct/internal/report"
)

func TestNoCredentialsCheck(t *testing.T) {
	// Test with nil config
	check := NewNoCredentialsCheck(nil)
	findings := check.Run(&checker.Project{})

	if len(findings) != 1 {
		t.Errorf("expected 1 finding for nil config, got %d", len(findings))
	}

	if findings[0].ID != "REV-001" {
		t.Errorf("expected ID REV-001, got %s", findings[0].ID)
	}

	// Test with empty config
	check = NewNoCredentialsCheck(&config.ReviewerConfig{})
	findings = check.Run(&checker.Project{})

	if len(findings) != 1 {
		t.Errorf("expected 1 finding for empty config, got %d", len(findings))
	}
}

func TestNoCredentialsCheckWithValidCredentials(t *testing.T) {
	// Set env vars
	os.Setenv("TEST_EMAIL", "reviewer@example.com")
	os.Setenv("TEST_PASSWORD", "SecurePass123!")
	defer os.Unsetenv("TEST_EMAIL")
	defer os.Unsetenv("TEST_PASSWORD")

	cfg := &config.ReviewerConfig{
		EmailEnv:    "TEST_EMAIL",
		PasswordEnv: "TEST_PASSWORD",
	}

	check := NewNoCredentialsCheck(cfg)
	findings := check.Run(&checker.Project{})

	if len(findings) != 0 {
		t.Errorf("expected 0 findings with valid credentials, got %d", len(findings))
	}
}

func TestNoCredentialsCheckWithDirectValues(t *testing.T) {
	cfg := &config.ReviewerConfig{
		Email:    "reviewer@example.com",
		Password: "SecurePass123!",
	}

	check := NewNoCredentialsCheck(cfg)
	findings := check.Run(&checker.Project{})

	if len(findings) != 0 {
		t.Errorf("expected 0 findings with direct credentials, got %d", len(findings))
	}
}

func TestPlaceholderEmailCheck(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected int // number of findings
	}{
		{"valid email", "reviewer@company.com", 0},
		{"test@example.com", "test@example.com", 1},
		{"your-email@example.com", "your-email@example.com", 1},
		{"demo@example.org", "demo@example.org", 1},
		{"admin@example", "admin@example", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.ReviewerConfig{
				Email: tt.email,
			}

			check := NewPlaceholderEmailCheck(cfg)
			findings := check.Run(&checker.Project{})

			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}

			if tt.expected > 0 && findings[0].ID != "REV-002" {
				t.Errorf("expected ID REV-002, got %s", findings[0].ID)
			}

			if tt.expected > 0 && findings[0].Severity != report.SeverityHigh {
				t.Errorf("expected HIGH severity, got %s", findings[0].Severity)
			}
		})
	}
}

func TestPlaceholderEmailCheckWithEmptyConfig(t *testing.T) {
	check := NewPlaceholderEmailCheck(nil)
	findings := check.Run(&checker.Project{})

	if len(findings) != 0 {
		t.Errorf("expected 0 findings with nil config, got %d", len(findings))
	}
}

func TestWeakPasswordCheck(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected int
		severity report.Severity
	}{
		{"strong password", "SecureP@ssw0rd!", 0, ""},
		{"too short", "short", 1, report.SeverityHigh},
		{"password", "password123", 1, report.SeverityHigh},
		{"qwerty", "qwerty123", 1, report.SeverityHigh},
		{"admin", "admin2024", 1, report.SeverityHigh},
		{"test", "testflutter", 1, report.SeverityHigh},
		{"repeated chars", "passsword123", 1, report.SeverityWarning},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.ReviewerConfig{
				Password: tt.password,
			}

			check := NewWeakPasswordCheck(cfg)
			findings := check.Run(&checker.Project{})

			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}

			if tt.expected > 0 && findings[0].ID != "REV-003" {
				t.Errorf("expected ID REV-003, got %s", findings[0].ID)
			}

			if tt.expected > 0 && tt.severity != "" && findings[0].Severity != tt.severity {
				t.Errorf("expected %s severity, got %s", tt.severity, findings[0].Severity)
			}
		})
	}
}

func TestWeakPasswordCheckWithEnvVar(t *testing.T) {
	os.Setenv("TEST_WEAK_PASS", "password123")
	defer os.Unsetenv("TEST_WEAK_PASS")

	cfg := &config.ReviewerConfig{
		PasswordEnv: "TEST_WEAK_PASS",
	}

	check := NewWeakPasswordCheck(cfg)
	findings := check.Run(&checker.Project{})

	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"reviewer@example.com", "re***@example.com"},
		{"ab@example.com", "***@example.com"},
		{"", ""},
		{"invalid", "***"},
		{"@example.com", "***@example.com"}, // Empty local part gets masked
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := maskEmail(tt.input)
			if result != tt.expected {
				t.Errorf("maskEmail(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHasRepeatedChars(t *testing.T) {
	tests := []struct {
		input     string
		threshold int
		expected  bool
	}{
		{"aaabbb", 3, true},
		{"password", 3, false},
		{"passsword", 3, true},
		{"abc", 3, false},
		{"aaa", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hasRepeatedChars(tt.input, tt.threshold)
			if result != tt.expected {
				t.Errorf("hasRepeatedChars(%q, %d) = %v, want %v",
					tt.input, tt.threshold, result, tt.expected)
			}
		})
	}
}

func TestGetConfigFromEnv(t *testing.T) {
	// Set test env vars
	os.Setenv("REVIEWER_EMAIL", "test@example.com")
	os.Setenv("REVIEWER_PASSWORD", "testpass123")
	defer os.Unsetenv("REVIEWER_EMAIL")
	defer os.Unsetenv("REVIEWER_PASSWORD")

	cfg := GetConfigFromEnv()

	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	if cfg.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", cfg.Email)
	}

	if cfg.Password != "testpass123" {
		t.Errorf("expected password 'testpass123', got '%s'", cfg.Password)
	}
}

func TestGetConfigFromEnvDefaults(t *testing.T) {
	// Clear env vars
	os.Unsetenv("REVIEWER_EMAIL")
	os.Unsetenv("REVIEWER_PASSWORD")

	cfg := GetConfigFromEnv()

	if cfg == nil {
		t.Fatal("expected config, got nil")
	}

	// Should have default env var names
	if cfg.EmailEnv != "REVIEWER_EMAIL" {
		t.Errorf("expected EmailEnv 'REVIEWER_EMAIL', got '%s'", cfg.EmailEnv)
	}

	if cfg.PasswordEnv != "REVIEWER_PASSWORD" {
		t.Errorf("expected PasswordEnv 'REVIEWER_PASSWORD', got '%s'", cfg.PasswordEnv)
	}
}
