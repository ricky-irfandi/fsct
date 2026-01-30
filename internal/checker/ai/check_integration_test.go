package ai

import (
	"context"
	"testing"

	aipkg "github.com/ricky-irfandi/fsct/internal/ai"
	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

// mockAIProvider is a mock AI provider for testing
type mockAIProvider struct {
	response string
	err      error
}

func (m *mockAIProvider) Name() string { return "mock" }
func (m *mockAIProvider) SetAPIKey(key string) {}
func (m *mockAIProvider) SetModel(model string) {}
func (m *mockAIProvider) AvailableModels() []string { return []string{"mock"} }

func (m *mockAIProvider) Complete(ctx context.Context, req *aipkg.CompletionRequest) (*aipkg.CompletionResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &aipkg.CompletionResponse{
		Content:      m.response,
		FinishReason: "stop",
		Model:        "mock-model",
	}, nil
}

func TestAI001WithMockResponse(t *testing.T) {
	mockResponse := `{
		"risk_level": "medium",
		"confidence": "high",
		"compliance_score": 75,
		"store_readiness": {
			"app_store": false,
			"play_store": true,
			"reasoning": "Missing iOS permission descriptions"
		},
		"insights": [
			{
				"category": "permissions",
				"severity": "high",
				"title": "Missing Camera Permission Description",
				"description": "NSCameraUsageDescription is required for camera access",
				"confidence": "high"
			}
		],
		"suggestions": [
			{
				"priority": 1,
				"category": "permissions",
				"issue": "Missing camera description",
				"action": "Add NSCameraUsageDescription to Info.plist",
				"file_path": "ios/Runner/Info.plist"
			}
		],
		"reviewer_notes": ["Camera feature requires permission grant"]
	}`

	provider := &mockAIProvider{response: mockResponse}
	client := aipkg.NewClient(provider, nil)
	check := AI001PermissionJustificationCheck(client)

	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:    "TestApp",
			Version: "1.0.0",
		},
		HasCameraDeps: true,
	}

	findings := check.Run(project)

	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}

	// Check first finding
	if findings[0].ID != "AI-001-001" {
		t.Errorf("expected ID 'AI-001-001', got '%s'", findings[0].ID)
	}

	if findings[0].Severity != report.SeverityHigh {
		t.Errorf("expected HIGH severity, got %s", findings[0].Severity)
	}

	if !contains(findings[0].Message, "Camera") {
		t.Errorf("finding should mention Camera, got: %s", findings[0].Message)
	}

	if !contains(findings[0].Suggestion, "Info.plist") {
		t.Errorf("suggestion should mention Info.plist, got: %s", findings[0].Suggestion)
	}
}

func TestAI002WithMockResponse(t *testing.T) {
	mockResponse := `{
		"risk_level": "low",
		"confidence": "medium",
		"compliance_score": 90,
		"store_readiness": {
			"app_store": true,
			"play_store": true
		},
		"insights": [
			{
				"category": "policy",
				"severity": "info",
				"title": "Privacy Policy Present",
				"description": "Privacy policy URL is configured"
			}
		],
		"suggestions": [],
		"reviewer_notes": []
	}`

	provider := &mockAIProvider{response: mockResponse}
	client := aipkg.NewClient(provider, nil)
	check := AI002PolicyComplianceCheck(client)

	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:    "TestApp",
			Version: "1.0.0",
		},
	}

	findings := check.Run(project)

	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}

	if findings[0].Severity != report.SeverityInfo {
		t.Errorf("expected INFO severity, got %s", findings[0].Severity)
	}
}

func TestAI003WithMockResponse(t *testing.T) {
	mockResponse := `{
		"risk_level": "medium",
		"confidence": "high",
		"insights": [
			{
				"category": "dependencies",
				"severity": "warning",
				"title": "Outdated HTTP Package",
				"description": "Consider updating to http 1.0.0 or later"
			}
		],
		"suggestions": [
			{
				"priority": 2,
				"category": "dependencies",
				"issue": "Outdated http package",
				"action": "Update pubspec.yaml to use http: ^1.0.0"
			}
		]
	}`

	provider := &mockAIProvider{response: mockResponse}
	client := aipkg.NewClient(provider, nil)
	check := AI003DependencyRiskCheck(client)

	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:    "TestApp",
			Version: "1.0.0",
			Dependencies: map[string]string{
				"http": "^0.13.0",
			},
		},
	}

	findings := check.Run(project)

	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}

	// Find the warning finding
	var warningFound bool
	for _, f := range findings {
		if f.Severity == report.SeverityWarning {
			warningFound = true
			break
		}
	}

	if !warningFound {
		t.Error("expected at least one WARNING severity finding")
	}
}

func TestAI004WithMockResponse(t *testing.T) {
	mockResponse := `{
		"risk_level": "low",
		"confidence": "high",
		"store_readiness": {
			"app_store": true,
			"play_store": true
		},
		"insights": [
			{
				"category": "store",
				"severity": "info",
				"title": "Ready for Submission",
				"description": "App meets both store requirements"
			}
		],
		"suggestions": [
			{
				"priority": 3,
				"category": "store",
				"issue": "Consider adding screenshots",
				"action": "Add screenshots for both stores",
				"platform": "both"
			}
		]
	}`

	provider := &mockAIProvider{response: mockResponse}
	client := aipkg.NewClient(provider, nil)
	check := AI004StoreGuidanceCheck(client)

	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:    "TestApp",
			Version: "1.0.0",
		},
	}

	findings := check.Run(project)

	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}
}

func TestAI005WithMockResponse(t *testing.T) {
	mockResponse := `{
		"risk_level": "low",
		"confidence": "high",
		"insights": [],
		"reviewer_notes": [
			"Test account: reviewer@example.com / ReviewPass123!",
			"All features accessible after login",
			"Demo data pre-populated"
		],
		"test_account": {
			"needed": true,
			"reason": "App requires login to access features",
			"setup_instructions": "Use provided test account credentials"
		},
		"demo_data": [
			"Sample user profile with photo",
			"3 pre-populated content items"
		],
		"special_instructions": [
			"Enable location services for map feature"
		]
	}`

	provider := &mockAIProvider{response: mockResponse}
	client := aipkg.NewClient(provider, nil)
	check := AI005ReviewerNotesCheck(client)

	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:    "TestApp",
			Version: "1.0.0",
		},
		HasLoginPatterns: true,
	}

	findings := check.Run(project)

	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}
}

func TestAIWithError(t *testing.T) {
	provider := &mockAIProvider{err: aipkg.ErrNoAPIKey}
	client := aipkg.NewClient(provider, nil)
	check := AI001PermissionJustificationCheck(client)

	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:    "TestApp",
			Version: "1.0.0",
		},
	}

	findings := check.Run(project)

	if len(findings) != 1 {
		t.Fatalf("expected 1 offline finding, got %d", len(findings))
	}

	if findings[0].Severity != report.SeverityInfo {
		t.Errorf("expected INFO severity for offline finding, got %s", findings[0].Severity)
	}

	if !contains(findings[0].Title, "Unavailable") {
		t.Errorf("offline finding title should mention 'Unavailable', got: %s", findings[0].Title)
	}
}

func TestAIWithPlainTextResponse(t *testing.T) {
	// Test fallback parsing for non-JSON response
	mockResponse := `Analysis Results:

Risk Level: medium

## Issues Found

- Critical: Missing camera permission description
- Warning: Location description is vague

## Recommendations

1. Add NSCameraUsageDescription to Info.plist
2. Clarify why location is needed in the description

Store Readiness:
- App Store: No (missing descriptions)
- Play Store: Yes`

	provider := &mockAIProvider{response: mockResponse}
	client := aipkg.NewClient(provider, nil)
	check := AI001PermissionJustificationCheck(client)

	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:    "TestApp",
			Version: "1.0.0",
		},
	}

	findings := check.Run(project)

	// Should still produce findings even without JSON
	if len(findings) == 0 {
		t.Fatal("expected at least one finding from plain text response")
	}
}
