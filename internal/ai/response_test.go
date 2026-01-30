package ai

import (
	"strings"
	"testing"
)

func TestParseResponse(t *testing.T) {
	jsonResponse := `{
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
				"issue": "Missing iOS permission descriptions",
				"action": "Add NSCameraUsageDescription to Info.plist",
				"file_path": "ios/Runner/Info.plist"
			}
		],
		"reviewer_notes": [
			"Test account credentials provided",
			"Camera feature requires permission grant"
		]
	}`

	analysis, err := ParseResponse(jsonResponse)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if analysis.RiskLevel != "medium" {
		t.Errorf("expected risk_level 'medium', got '%s'", analysis.RiskLevel)
	}

	if analysis.Confidence != "high" {
		t.Errorf("expected confidence 'high', got '%s'", analysis.Confidence)
	}

	if analysis.ComplianceScore != 75 {
		t.Errorf("expected compliance_score 75, got %d", analysis.ComplianceScore)
	}

	if analysis.StoreReadiness.AppStore {
		t.Error("expected app_store readiness to be false")
	}

	if !analysis.StoreReadiness.PlayStore {
		t.Error("expected play_store readiness to be true")
	}

	if len(analysis.Insights) != 1 {
		t.Errorf("expected 1 insight, got %d", len(analysis.Insights))
	}

	if len(analysis.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(analysis.Suggestions))
	}

	if analysis.Suggestions[0].Priority != 1 {
		t.Errorf("expected priority 1, got %d", analysis.Suggestions[0].Priority)
	}
}

func TestParseResponseWithCodeBlock(t *testing.T) {
	response := `Here's the analysis in JSON format:

` + "```json" + `
{
	"risk_level": "low",
	"store_readiness": {
		"app_store": true,
		"play_store": true
	},
	"insights": []
}
` + "```" + `

Let me know if you need more details!`

	analysis, err := ParseResponse(response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if analysis.RiskLevel != "low" {
		t.Errorf("expected risk_level 'low', got '%s'", analysis.RiskLevel)
	}

	if !analysis.StoreReadiness.AppStore {
		t.Error("expected app_store readiness to be true")
	}
}

func TestParsePlainTextResponse(t *testing.T) {
	response := `## Analysis Results

Risk Level: High

This app has critical issues that need to be addressed.

## Store Readiness

App Store: No
Play Store: Yes

## Issues

- Critical: Missing privacy policy URL
- Warning: Using outdated dependencies
- Info: Consider adding more tests

## Suggestions

1. Add a privacy policy URL to your app
2. Update dependencies to latest versions
3. Add unit tests for critical functionality`

	analysis := parsePlainTextResponse(response)

	if analysis.RiskLevel != "high" {
		t.Errorf("expected risk_level 'high', got '%s'", analysis.RiskLevel)
	}

	if analysis.StoreReadiness.AppStore {
		t.Error("expected app_store readiness to be false")
	}

	if !analysis.StoreReadiness.PlayStore {
		t.Error("expected play_store readiness to be true")
	}

	if len(analysis.Insights) < 2 {
		t.Errorf("expected at least 2 insights, got %d", len(analysis.Insights))
	}

	if len(analysis.Suggestions) < 2 {
		t.Errorf("expected at least 2 suggestions, got %d", len(analysis.Suggestions))
	}
}

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "json with backticks",
			input: "```json\n{\"key\": \"value\"}\n```",
			expected: `{"key": "value"}`,
		},
		{
			name: "json without language",
			input: "```\n{\"key\": \"value\"}\n```",
			expected: `{"key": "value"}`,
		},
		{
			name: "plain json",
			input: `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name: "text with json",
			input: "Here is the result: {\"key\": \"value\"} Thanks!",
			expected: `{"key": "value"}`,
		},
		{
			name:     "no json",
			input:    "Just plain text without JSON",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSON(tt.input)
			if result != tt.expected {
				t.Errorf("extractJSON() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestAIAnalysisConvertToFindings(t *testing.T) {
	analysis := &AIAnalysis{
		Insights: []AIInsight{
			{
				Category:    "permissions",
				Severity:    "high",
				Title:       "Missing Camera Permission",
				Description: "Camera permission description is missing",
			},
			{
				Category:    "policy",
				Severity:    "warning",
				Title:       "Privacy Policy",
				Description: "Privacy policy URL not found",
			},
			{
				Category:    "general",
				Severity:    "info",
				Title:       "Test Coverage",
				Description: "Consider adding more tests",
			},
		},
		Suggestions: []AISuggestion{
			{
				Priority: 1,
				Issue:    "Missing Camera Permission",
				Action:   "Add NSCameraUsageDescription",
			},
		},
	}

	findings := analysis.ConvertToFindings("AI-001")

	if len(findings) != 3 {
		t.Errorf("expected 3 findings, got %d", len(findings))
	}

	// Check first finding (high severity)
	if findings[0].ID != "AI-001-001" {
		t.Errorf("expected ID 'AI-001-001', got '%s'", findings[0].ID)
	}

	if findings[0].Title != "Missing Camera Permission" {
		t.Errorf("expected title 'Missing Camera Permission', got '%s'", findings[0].Title)
	}

	if findings[0].Suggestion != "Add NSCameraUsageDescription" {
		t.Errorf("expected suggestion 'Add NSCameraUsageDescription', got '%s'", findings[0].Suggestion)
	}

	// Check second finding (warning severity)
	if findings[1].ID != "AI-001-002" {
		t.Errorf("expected ID 'AI-001-002', got '%s'", findings[1].ID)
	}
}

func TestAIAnalysisIsReady(t *testing.T) {
	tests := []struct {
		name     string
		analysis *AIAnalysis
		expected bool
	}{
		{
			name: "ready",
			analysis: &AIAnalysis{
				RiskLevel: "low",
				StoreReadiness: StoreReadiness{
					AppStore:  true,
					PlayStore: true,
				},
			},
			expected: true,
		},
		{
			name: "not ready - high risk",
			analysis: &AIAnalysis{
				RiskLevel: "high",
				StoreReadiness: StoreReadiness{
					AppStore:  true,
					PlayStore: true,
				},
			},
			expected: false,
		},
		{
			name: "not ready - app store",
			analysis: &AIAnalysis{
				RiskLevel: "low",
				StoreReadiness: StoreReadiness{
					AppStore:  false,
					PlayStore: true,
				},
			},
			expected: false,
		},
		{
			name: "not ready - play store",
			analysis: &AIAnalysis{
				RiskLevel: "low",
				StoreReadiness: StoreReadiness{
					AppStore:  true,
					PlayStore: false,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.analysis.IsReady()
			if result != tt.expected {
				t.Errorf("IsReady() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPriorityFindings(t *testing.T) {
	analysis := &AIAnalysis{
		Insights: []AIInsight{
			{Severity: "info", Title: "Low priority 1"},
			{Severity: "high", Title: "High priority 1"},
			{Severity: "warning", Title: "Medium priority"},
			{Severity: "high", Title: "High priority 2"},
			{Severity: "info", Title: "Low priority 2"},
		},
	}

	// Get all
	all := analysis.GetPriorityFindings(0)
	if len(all) != 5 {
		t.Errorf("expected 5 findings, got %d", len(all))
	}

	// Check order (high -> warning -> info)
	if all[0].Severity != "high" {
		t.Errorf("first finding should be high severity, got %s", all[0].Severity)
	}

	if all[3].Severity != "info" {
		t.Errorf("fourth finding should be info severity, got %s", all[3].Severity)
	}

	// Get limited
	limited := analysis.GetPriorityFindings(2)
	if len(limited) != 2 {
		t.Errorf("expected 2 findings, got %d", len(limited))
	}
}

func TestExpectedJSONSchema(t *testing.T) {
	schema := ExpectedJSONSchema()

	if schema == "" {
		t.Error("schema should not be empty")
	}

	// Check for key fields
	requiredFields := []string{
		"risk_level",
		"store_readiness",
		"insights",
		"suggestions",
		"reviewer_notes",
	}

	for _, field := range requiredFields {
		if !strings.Contains(schema, field) {
			t.Errorf("schema should contain field %q", field)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input   string
		maxLen  int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"short", 5, "short"},
		{"longer text", 5, "lo..."},
	}

	for _, tt := range tests {
		result := truncateString(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}
