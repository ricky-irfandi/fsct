package ai

import (
	"errors"
	"testing"

	aipkg "github.com/ricky-irfandi/fsct/internal/ai"
	"github.com/ricky-irfandi/fsct/internal/report"
)

func TestNewCheck(t *testing.T) {
	client := &aipkg.Client{}
	check := NewCheck("AI-001", "Test Check", "AI Analysis", client)

	if check.ID() != "AI-001" {
		t.Errorf("expected ID 'AI-001', got '%s'", check.ID())
	}

	if check.Name() != "Test Check" {
		t.Errorf("expected name 'Test Check', got '%s'", check.Name())
	}

	if check.Category() != "AI Analysis" {
		t.Errorf("expected category 'AI Analysis', got '%s'", check.Category())
	}
}

func TestCheckConvertToFindings(t *testing.T) {
	client := &aipkg.Client{}
	check := NewCheck("AI-001", "Test Check", "AI Analysis", client)

	analysis := &aipkg.AIAnalysis{
		RiskLevel: "medium",
		Insights: []aipkg.AIInsight{
			{
				Category:    "permissions",
				Severity:    "high",
				Title:       "Missing Camera Description",
				Description: "Camera permission needs description",
			},
			{
				Category:    "permissions",
				Severity:    "warning",
				Title:       "Location Description Vague",
				Description: "Location description could be clearer",
			},
		},
		Suggestions: []aipkg.AISuggestion{
			{
				Priority: 1,
				Issue:    "Missing description",
				Action:   "Add NSCameraUsageDescription",
			},
			{
				Priority: 2,
				Issue:    "Vague description",
				Action:   "Clarify location usage",
			},
		},
	}

	findings := check.convertToFindings(analysis)

	if len(findings) != 2 {
		t.Errorf("expected 2 findings, got %d", len(findings))
	}

	// Check first finding
	if findings[0].ID != "AI-001-001" {
		t.Errorf("expected ID 'AI-001-001', got '%s'", findings[0].ID)
	}

	if findings[0].Severity != report.SeverityHigh {
		t.Errorf("expected HIGH severity, got %s", findings[0].Severity)
	}

	if findings[0].Suggestion != "Add NSCameraUsageDescription" {
		t.Errorf("expected suggestion 'Add NSCameraUsageDescription', got '%s'", findings[0].Suggestion)
	}

	// Check second finding
	if findings[1].Severity != report.SeverityWarning {
		t.Errorf("expected WARNING severity, got %s", findings[1].Severity)
	}
}

func TestCheckConvertToFindingsEmpty(t *testing.T) {
	client := &aipkg.Client{}
	check := NewCheck("AI-001", "Test Check", "AI Analysis", client)

	// Test with no insights but suggestions
	analysis := &aipkg.AIAnalysis{
		RiskLevel: "low",
		Suggestions: []aipkg.AISuggestion{
			{
				Priority: 1,
				Issue:    "Consider adding tests",
				Action:   "Add unit tests",
			},
		},
	}

	findings := check.convertToFindings(analysis)

	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}

	// Test with completely empty analysis
	emptyAnalysis := &aipkg.AIAnalysis{
		RiskLevel:   "low",
		Insights:    []aipkg.AIInsight{},
		Suggestions: []aipkg.AISuggestion{},
	}

	findings = check.convertToFindings(emptyAnalysis)

	if len(findings) != 1 {
		t.Errorf("expected 1 fallback finding, got %d", len(findings))
	}

	if findings[0].ID != "AI-001-001" {
		t.Errorf("expected ID 'AI-001-001', got '%s'", findings[0].ID)
	}
}

func TestCreateOfflineFinding(t *testing.T) {
	client := &aipkg.Client{}
	check := NewCheck("AI-001", "Test Check", "AI Analysis", client)

	err := errors.New("connection timeout")
	finding := check.createOfflineFinding(err)

	if finding.ID != "AI-001" {
		t.Errorf("expected ID 'AI-001', got '%s'", finding.ID)
	}

	if finding.Severity != report.SeverityInfo {
		t.Errorf("expected INFO severity, got %s", finding.Severity)
	}

	if finding.Title != "[Test Check] AI Analysis Unavailable" {
		t.Errorf("unexpected title: %s", finding.Title)
	}

	if finding.Suggestion == "" {
		t.Error("expected non-empty suggestion for offline finding")
	}
}

func TestGetSystemPrompt(t *testing.T) {
	client := &aipkg.Client{}
	prompts := &SystemPrompts{}

	tests := []struct {
		id       string
		expected string
	}{
		{"AI-001", prompts.AI001PermissionJustification()},
		{"AI-002", prompts.AI002PolicyCompliance()},
		{"AI-003", prompts.AI003DependencyRisk()},
		{"AI-004", prompts.AI004StoreGuidance()},
		{"AI-005", prompts.AI005ReviewerNotes()},
		{"AI-999", prompts.AI001PermissionJustification()}, // Default fallback
	}

	for _, tt := range tests {
		check := NewCheck(tt.id, "Test", "AI Analysis", client)
		prompt := check.getSystemPrompt()

		if prompt == "" {
			t.Errorf("check %s returned empty prompt", tt.id)
		}

		// Verify prompt contains expected content
		if !contains(prompt, "JSON") {
			t.Errorf("check %s prompt should mention JSON", tt.id)
		}
	}
}

func TestAllChecks(t *testing.T) {
	// Test with nil client
	checks := AllChecks(nil)
	if checks != nil {
		t.Error("expected nil checks for nil client")
	}

	// Test with unavailable client
	unavailableClient := &aipkg.Client{}
	checks = AllChecks(unavailableClient)
	if checks != nil {
		t.Error("expected nil checks for unavailable client")
	}
}

func TestIndividualCheckConstructors(t *testing.T) {
	client := &aipkg.Client{}

	check1 := AI001PermissionJustificationCheck(client)
	if check1.ID() != "AI-001" {
		t.Errorf("AI001: expected ID 'AI-001', got '%s'", check1.ID())
	}

	check2 := AI002PolicyComplianceCheck(client)
	if check2.ID() != "AI-002" {
		t.Errorf("AI002: expected ID 'AI-002', got '%s'", check2.ID())
	}

	check3 := AI003DependencyRiskCheck(client)
	if check3.ID() != "AI-003" {
		t.Errorf("AI003: expected ID 'AI-003', got '%s'", check3.ID())
	}

	check4 := AI004StoreGuidanceCheck(client)
	if check4.ID() != "AI-004" {
		t.Errorf("AI004: expected ID 'AI-004', got '%s'", check4.ID())
	}

	check5 := AI005ReviewerNotesCheck(client)
	if check5.ID() != "AI-005" {
		t.Errorf("AI005: expected ID 'AI-005', got '%s'", check5.ID())
	}
}

func TestBuildUserPrompt(t *testing.T) {
	metadata := &aipkg.ComplianceMetadata{
		AppName:         "TestApp",
		Version:         "1.0.0",
		Description:     "A test app",
		ComplianceScore: 85,
		Findings: []aipkg.FindingMeta{
			{ID: "AND-001", Severity: "HIGH", Title: "Target SDK", Category: "Android"},
		},
		AndroidTargetSDK: 34,
		AndroidMinSDK:    21,
		AndroidPermissions: []string{"CAMERA", "LOCATION"},
		Dependencies:   []string{"http", "camera"},
		Features: aipkg.AppFeatures{
			HasCamera:   true,
			HasLocation: true,
		},
		Security: aipkg.SecurityFlags{
			IsDebuggable: false,
		},
	}

	prompt := BuildUserPrompt("AI Permission Analysis", metadata)

	// Check that prompt contains key information
	checks := []string{
		"TestApp",
		"1.0.0",
		"AND-001",
		"CAMERA",
		"LOCATION",
		"34",
		"JSON",
	}

	for _, check := range checks {
		if !contains(prompt, check) {
			t.Errorf("prompt should contain '%s'", check)
		}
	}
}

func TestSystemPrompts(t *testing.T) {
	prompts := &SystemPrompts{}

	// Test each prompt returns non-empty string
	if prompts.AI001PermissionJustification() == "" {
		t.Error("AI001 prompt should not be empty")
	}

	if prompts.AI002PolicyCompliance() == "" {
		t.Error("AI002 prompt should not be empty")
	}

	if prompts.AI003DependencyRisk() == "" {
		t.Error("AI003 prompt should not be empty")
	}

	if prompts.AI004StoreGuidance() == "" {
		t.Error("AI004 prompt should not be empty")
	}

	if prompts.AI005ReviewerNotes() == "" {
		t.Error("AI005 prompt should not be empty")
	}

	// Verify prompts contain expected keywords
	p1 := prompts.AI001PermissionJustification()
	if !contains(p1, "permission") {
		t.Error("AI001 prompt should mention permissions")
	}

	p2 := prompts.AI002PolicyCompliance()
	if !contains(p2, "policy") {
		t.Error("AI002 prompt should mention policy")
	}

	p3 := prompts.AI003DependencyRisk()
	if !contains(p3, "dependenc") {
		t.Error("AI003 prompt should mention dependencies")
	}

	p4 := prompts.AI004StoreGuidance()
	if !contains(p4, "store") {
		t.Error("AI004 prompt should mention store")
	}

	p5 := prompts.AI005ReviewerNotes()
	if !contains(p5, "reviewer") {
		t.Error("AI005 prompt should mention reviewer")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
