package reviewer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/config"
	"github.com/ricky-irfandi/fsct/internal/report"
)

// LoginVerificationCheck verifies reviewer credentials work (REV-004, REV-005)
type LoginVerificationCheck struct {
	config *config.ReviewerConfig
}

// NewLoginVerificationCheck creates a new login verification check
func NewLoginVerificationCheck(cfg *config.ReviewerConfig) *LoginVerificationCheck {
	return &LoginVerificationCheck{config: cfg}
}

// ID returns the check ID (uses both REV-004 and REV-005 based on result)
func (c *LoginVerificationCheck) ID() string {
	return "REV-004"
}

// Name returns the check name
func (c *LoginVerificationCheck) Name() string {
	return "Reviewer Login Verification"
}

// Run executes the check
func (c *LoginVerificationCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	// Skip if verification is not enabled or not configured
	if c.config == nil || c.config.Verification == nil || !c.config.Verification.Enabled {
		return findings
	}

	// Get credentials
	email := c.getEmail()
	password := c.getPassword()

	if email == "" || password == "" {
		findings = append(findings, report.Finding{
			ID:       "REV-004",
			Severity: report.SeverityWarning,
			Title:    "Cannot Verify Login - Missing Credentials",
			Message:  "Login verification is enabled but credentials are missing.",
			Suggestion: "Set REVIEWER_EMAIL and REVIEWER_PASSWORD environment variables.",
		})
		return findings
	}

	// Perform login verification
	result := c.verifyLogin(email, password)

	if result.Error != nil {
		// Network or system error (REV-005)
		findings = append(findings, report.Finding{
			ID:       "REV-005",
			Severity: report.SeverityHigh,
			Title:    "Login Verification Failed",
			Message:  fmt.Sprintf("Could not verify login: %v", result.Error),
			Suggestion: "Check the auth endpoint URL and network connectivity. " +
				"Ensure the endpoint is accessible from this environment.",
		})
		return findings
	}

	if !result.Success {
		// Login failed - invalid credentials (REV-005)
		findings = append(findings, report.Finding{
			ID:       "REV-005",
			Severity: report.SeverityHigh,
			Title:    "Reviewer Login Failed",
			Message:  "The provided reviewer credentials could not authenticate successfully.",
			Suggestion: "Verify the email and password are correct. " +
				"Try logging in manually to confirm the credentials work.",
		})
		return findings
	}

	if result.TokenExpired {
		// Token expired (REV-004)
		findings = append(findings, report.Finding{
			ID:       "REV-004",
			Severity: report.SeverityHigh,
			Title:    "Reviewer Token Expired",
			Message:  "The authentication token for the reviewer account has expired.",
			Suggestion: "Generate fresh credentials for app reviewers. " +
				"Update the REVIEWER_EMAIL and REVIEWER_PASSWORD environment variables.",
		})
		return findings
	}

	// Login successful
	findings = append(findings, report.Finding{
		ID:       "REV-004",
		Severity: report.SeverityInfo,
		Title:    "Reviewer Login Verified",
		Message:  "Reviewer credentials are valid and working.",
		Suggestion: "The test account is ready for app reviewers.",
	})

	return findings
}

// VerificationResult holds the result of a login verification
type VerificationResult struct {
	Success      bool
	TokenExpired bool
	Error        error
	ResponseBody string
}

// verifyLogin attempts to verify login credentials
func (c *LoginVerificationCheck) verifyLogin(email, password string) *VerificationResult {
	cfg := c.config.Verification

	// Build request body
	bodyTemplate := cfg.BodyTemplate
	if bodyTemplate == "" {
		// Default JSON body
		bodyTemplate = `{"email":"{{email}}","password":"{{password}}"}`
	}

	body := strings.NewReplacer(
		"{{email}}", email,
		"{{password}}", password,
	).Replace(bodyTemplate)

	// Create HTTP request
	method := cfg.Method
	if method == "" {
		method = "POST"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, cfg.AuthEndpoint, bytes.NewBufferString(body))
	if err != nil {
		return &VerificationResult{Error: fmt.Errorf("failed to create request: %w", err)}
	}

	req.Header.Set("Content-Type", "application/json")

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &VerificationResult{Error: fmt.Errorf("request failed: %w", err)}
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &VerificationResult{Error: fmt.Errorf("failed to read response: %w", err)}
	}
	responseBody := string(bodyBytes)

	// Check for success indicator
	successIndicator := cfg.SuccessIndicator
	if successIndicator == "" {
		successIndicator = "token"
	}

	// Determine result based on response
	result := &VerificationResult{
		ResponseBody: responseBody,
	}

	// Check for success
	if strings.Contains(responseBody, successIndicator) {
		result.Success = true
	}

	// Check for token expiration indicators
	expiredIndicators := []string{
		"expired", "invalid_token", "token_expired",
		"session_expired", "unauthorized", "401",
	}

	responseLower := strings.ToLower(responseBody)
	for _, indicator := range expiredIndicators {
		if strings.Contains(responseLower, indicator) {
			result.TokenExpired = true
			result.Success = false
			break
		}
	}

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		result.Success = false
		if resp.StatusCode == 401 {
			result.TokenExpired = true
		}
	}

	return result
}

func (c *LoginVerificationCheck) getEmail() string {
	if c.config.Email != "" {
		return c.config.Email
	}
	if c.config.EmailEnv != "" {
		return getEnvOrDefault(c.config.EmailEnv, "")
	}
	return ""
}

func (c *LoginVerificationCheck) getPassword() string {
	if c.config.Password != "" {
		return c.config.Password
	}
	if c.config.PasswordEnv != "" {
		return getEnvOrDefault(c.config.PasswordEnv, "")
	}
	return ""
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := getEnv(key); value != "" {
		return value
	}
	return defaultValue
}

// GenerateReviewerInstructions generates reviewer instructions document
func GenerateReviewerInstructions(config *config.ReviewerConfig) string {
	if config == nil {
		return ""
	}

	email := config.Email
	if email == "" && config.EmailEnv != "" {
		email = getEnv(config.EmailEnv)
	}

	var sb strings.Builder
	sb.WriteString("═══════════════════════════════════════════════════════════════════\n")
	sb.WriteString("REVIEWER TEST ACCOUNT INFORMATION\n")
	sb.WriteString("═══════════════════════════════════════════════════════════════════\n\n")

	sb.WriteString("App Store Connect / Play Console Reviewer Information:\n\n")

	sb.WriteString("Test Account Credentials:\n")
	if email != "" {
		sb.WriteString(fmt.Sprintf("• Email: %s\n", maskEmail(email)))
	} else {
		sb.WriteString("• Email: [Not configured]\n")
	}
	sb.WriteString("• Password: [REDACTED - provided separately]\n\n")

	sb.WriteString("Test Account Setup:\n")
	sb.WriteString("• Account has pre-populated data for testing\n")
	sb.WriteString("• All app features are accessible\n")
	sb.WriteString("• No 2FA required on this account\n\n")

	sb.WriteString("Special Instructions for Reviewers:\n")
	sb.WriteString("1. Login with provided credentials\n")
	sb.WriteString("2. Main features are available immediately after login\n")
	sb.WriteString("3. Test data is clearly marked as 'Demo' or 'Test'\n\n")

	sb.WriteString("If you encounter any issues:\n")
	sb.WriteString("• Contact: developer@example.com\n")
	sb.WriteString("• Reference: Reviewer Account for App Submission\n\n")

	sb.WriteString("═══════════════════════════════════════════════════════════════════\n")

	return sb.String()
}

// ParseVerificationResponse parses a JSON verification response
func ParseVerificationResponse(body string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}
	return result, nil
}
