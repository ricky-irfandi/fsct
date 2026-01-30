package reviewer

import (
	"strings"
	"testing"

	"github.com/ricky-irfandi/fsct/internal/config"
)

func TestGenerateChecklist(t *testing.T) {
	cfg := &config.ReviewerConfig{
		Email:    "reviewer@example.com",
		Password: "SecurePass123!",
	}

	checklist := GenerateChecklist(cfg)

	if checklist == nil {
		t.Fatal("expected checklist, got nil")
	}

	if len(checklist.Items) == 0 {
		t.Error("expected items in checklist")
	}

	if !checklist.EmailConfigured {
		t.Error("expected EmailConfigured to be true")
	}

	if !checklist.PasswordConfigured {
		t.Error("expected PasswordConfigured to be true")
	}
}

func TestGenerateChecklistWithNilConfig(t *testing.T) {
	checklist := GenerateChecklist(nil)

	if checklist == nil {
		t.Fatal("expected checklist, got nil")
	}

	// Should have at least the config item
	found := false
	for _, item := range checklist.Items {
		if item.ID == "config" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'config' item in checklist")
	}
}

func TestGenerateChecklistWithPlaceholderEmail(t *testing.T) {
	cfg := &config.ReviewerConfig{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	checklist := GenerateChecklist(cfg)

	// Find the valid item
	var validItem *ChecklistItem
	for i := range checklist.Items {
		if checklist.Items[i].ID == "valid" {
			validItem = &checklist.Items[i]
			break
		}
	}

	if validItem == nil {
		t.Fatal("expected 'valid' item in checklist")
	}

	if validItem.Checked {
		t.Error("expected valid item to be unchecked for placeholder email")
	}
}

func TestGenerateChecklistWithWeakPassword(t *testing.T) {
	cfg := &config.ReviewerConfig{
		Email:    "reviewer@example.com",
		Password: "password123",
	}

	checklist := GenerateChecklist(cfg)

	// Find the strong item
	var strongItem *ChecklistItem
	for i := range checklist.Items {
		if checklist.Items[i].ID == "strong" {
			strongItem = &checklist.Items[i]
			break
		}
	}

	if strongItem == nil {
		t.Fatal("expected 'strong' item in checklist")
	}

	if strongItem.Checked {
		t.Error("expected strong item to be unchecked for weak password")
	}
}

func TestToMarkdown(t *testing.T) {
	cfg := &config.ReviewerConfig{
		Email:    "reviewer@example.com",
		Password: "SecurePass123!",
	}

	checklist := GenerateChecklist(cfg)
	markdown := checklist.ToMarkdown()

	if markdown == "" {
		t.Error("expected non-empty markdown")
	}

	// Check for expected content
	expectedContent := []string{
		"# Reviewer Account Checklist",
		"Required Items",
		"Recommended Items",
		"Status Summary",
	}

	for _, content := range expectedContent {
		if !strings.Contains(markdown, content) {
			t.Errorf("markdown should contain '%s'", content)
		}
	}
}

func TestToConsole(t *testing.T) {
	cfg := &config.ReviewerConfig{
		Email:    "reviewer@example.com",
		Password: "SecurePass123!",
	}

	checklist := GenerateChecklist(cfg)
	console := checklist.ToConsole()

	if console == "" {
		t.Error("expected non-empty console output")
	}

	// Check for expected content
	if !strings.Contains(console, "REVIEWER ACCOUNT CHECKLIST") {
		t.Error("console output should contain 'REVIEWER ACCOUNT CHECKLIST'")
	}
}

func TestIsReady(t *testing.T) {
	// Test with all required items checked
	checklist := &ReviewerChecklist{
		Items: []ChecklistItem{
			{ID: "1", Required: true, Checked: true},
			{ID: "2", Required: true, Checked: true},
			{ID: "3", Required: false, Checked: false},
		},
	}

	if !checklist.IsReady() {
		t.Error("expected IsReady to be true when all required items are checked")
	}

	// Test with missing required item
	checklist.Items[1].Checked = false
	if checklist.IsReady() {
		t.Error("expected IsReady to be false when required item is unchecked")
	}
}

func TestGetCompletionPercentage(t *testing.T) {
	// Test with 50% completion
	checklist := &ReviewerChecklist{
		Items: []ChecklistItem{
			{ID: "1", Checked: true},
			{ID: "2", Checked: false},
		},
	}

	percentage := checklist.GetCompletionPercentage()
	if percentage != 50.0 {
		t.Errorf("expected 50.0%%, got %.1f%%", percentage)
	}

	// Test with empty checklist
	emptyChecklist := &ReviewerChecklist{
		Items: []ChecklistItem{},
	}

	percentage = emptyChecklist.GetCompletionPercentage()
	if percentage != 100.0 {
		t.Errorf("expected 100.0%% for empty checklist, got %.1f%%", percentage)
	}

	// Test with all checked
	allChecked := &ReviewerChecklist{
		Items: []ChecklistItem{
			{ID: "1", Checked: true},
			{ID: "2", Checked: true},
			{ID: "3", Checked: true},
		},
	}

	percentage = allChecked.GetCompletionPercentage()
	if percentage != 100.0 {
		t.Errorf("expected 100.0%%, got %.1f%%", percentage)
	}
}

func TestChecklistItemTypes(t *testing.T) {
	checklist := &ReviewerChecklist{
		Items: []ChecklistItem{
			{ID: "email", Description: "Email configured", Required: true, Checked: true},
			{ID: "password", Description: "Password configured", Required: true, Checked: true},
			{ID: "test_data", Description: "Test data ready", Required: false, Checked: false},
		},
	}

	// Count required vs optional
	requiredCount := 0
	optionalCount := 0
	for _, item := range checklist.Items {
		if item.Required {
			requiredCount++
		} else {
			optionalCount++
		}
	}

	if requiredCount != 2 {
		t.Errorf("expected 2 required items, got %d", requiredCount)
	}

	if optionalCount != 1 {
		t.Errorf("expected 1 optional item, got %d", optionalCount)
	}
}

func TestGenerateReviewerInstructions(t *testing.T) {
	cfg := &config.ReviewerConfig{
		Email: "reviewer@example.com",
	}

	instructions := GenerateReviewerInstructions(cfg)

	if instructions == "" {
		t.Error("expected non-empty instructions")
	}

	// Check for expected content
	expectedContent := []string{
		"REVIEWER TEST ACCOUNT INFORMATION",
		"Test Account Credentials",
		"App Store Connect",
		"Play Console",
	}

	for _, content := range expectedContent {
		if !strings.Contains(instructions, content) {
			t.Errorf("instructions should contain '%s'", content)
		}
	}

	// Email should be masked
	if strings.Contains(instructions, "reviewer@example.com") {
		t.Error("email should be masked in instructions")
	}
}

func TestGenerateReviewerInstructionsWithNilConfig(t *testing.T) {
	instructions := GenerateReviewerInstructions(nil)
	if instructions != "" {
		t.Error("expected empty instructions for nil config")
	}
}
