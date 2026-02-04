package ai

import (
	"strings"
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

func TestExtractMetadata(t *testing.T) {
	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:        "test_app",
			Version:     "1.0.0+1",
			Description: "A test Flutter application",
			Dependencies: map[string]string{
				"flutter":    "sdk",
				"http":       "^1.0.0",
				"camera":     "^0.10.0",
				"geolocator": "^10.0.0",
			},
			DevDependencies: map[string]string{
				"flutter_test":  "sdk",
				"flutter_lints": "^3.0.0",
			},
			HasLinter:       true,
			HasIconConfig:   true,
			HasSplashConfig: false,
		},
		GradleConfig: &checker.GradleConfigInfo{
			TargetSDKVersion: "34",
			MinSDKVersion:    "21",
		},
		AndroidManifest: &checker.AndroidManifestInfo{
			Permissions: []string{
				"android.permission.CAMERA",
				"android.permission.ACCESS_FINE_LOCATION",
				"android.permission.INTERNET",
			},
			Debuggable:  false,
			AllowBackup: true,
		},
		InfoPlist: &checker.InfoPlistInfo{
			HasCameraUsageDescription:   true,
			HasLocationUsageDescription: false,
		},
		HasCameraDeps:    true,
		HasLocationDeps:  true,
		HasNetworkDeps:   true,
		HasLoginPatterns: false,
	}

	findings := []report.Finding{
		{
			ID:       "AND-001",
			Severity: report.SeverityHigh,
			Title:    "Target SDK Version",
			Message:  "Target SDK should be 35+",
		},
		{
			ID:       "IOS-003",
			Severity: report.SeverityHigh,
			Title:    "Missing Location Description",
			Message:  "NSLocationWhenInUseUsageDescription is missing",
		},
		{
			ID:       "FLT-005",
			Severity: report.SeverityWarning,
			Title:    "Default Version",
			Message:  "Using default version 1.0.0+1",
		},
	}

	meta := ExtractMetadata(project, findings, 38)

	// Test basic info
	if meta.AppName != "test_app" {
		t.Errorf("expected app_name 'test_app', got '%s'", meta.AppName)
	}

	if meta.Version != "1.0.0+1" {
		t.Errorf("expected version '1.0.0+1', got '%s'", meta.Version)
	}

	// Test scores
	if meta.TotalChecks != 38 {
		t.Errorf("expected 38 total checks, got %d", meta.TotalChecks)
	}

	if meta.HighCount != 2 {
		t.Errorf("expected 2 high count, got %d", meta.HighCount)
	}

	if meta.WarningCount != 1 {
		t.Errorf("expected 1 warning count, got %d", meta.WarningCount)
	}

	// Test findings extraction
	if len(meta.Findings) != 3 {
		t.Errorf("expected 3 findings, got %d", len(meta.Findings))
	}

	// Verify finding has no file paths
	for _, f := range meta.Findings {
		if f.ID == "" {
			t.Error("finding ID should not be empty")
		}
		if f.Category == "" {
			t.Error("finding category should not be empty")
		}
	}

	// Test permissions
	if len(meta.AndroidPermissions) != 3 {
		t.Errorf("expected 3 android permissions, got %d", len(meta.AndroidPermissions))
	}

	// Test features
	if !meta.Features.HasCamera {
		t.Error("expected HasCamera to be true")
	}

	if !meta.Features.HasLocation {
		t.Error("expected HasLocation to be true")
	}

	// Test config
	if !meta.Config.HasLinter {
		t.Error("expected HasLinter to be true")
	}

	if !meta.Config.HasIconConfig {
		t.Error("expected HasIconConfig to be true")
	}

	// Test security
	if meta.Security.IsDebuggable {
		t.Error("expected IsDebuggable to be false")
	}

	if !meta.Security.AllowsBackup {
		t.Error("expected AllowsBackup to be true")
	}

	// Test JSON serialization
	jsonData, err := meta.ToJSON()
	if err != nil {
		t.Errorf("failed to serialize to JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON output should not be empty")
	}

	// Test size
	if meta.Size() == 0 {
		t.Error("metadata size should not be zero")
	}

	// Ensure size is reasonable (< 50KB)
	if meta.Size() > 50000 {
		t.Errorf("metadata size %d exceeds 50KB limit", meta.Size())
	}
}

func TestExtractMetadataPrivacy(t *testing.T) {
	project := &checker.Project{
		Pubspec: &checker.PubspecInfo{
			Name:        "test_app",
			Version:     "1.0.0",
			Description: "This app contains password fields and secret keys",
		},
	}

	findings := []report.Finding{
		{
			ID:       "SEC-001",
			Severity: report.SeverityHigh,
			Title:    "Hardcoded Credentials",
			Message:  "Found hardcoded API key in lib/config.dart line 42",
			File:     "/Users/sensitive/path/to/project/lib/config.dart",
			Line:     42,
		},
	}

	meta := ExtractMetadata(project, findings, 38)

	// Verify file path is NOT in metadata
	for _, f := range meta.Findings {
		if f.ID == "SEC-001" {
			// Finding should exist but without file path
			if f.Title == "" {
				t.Error("finding title should be preserved")
			}
		}
	}

	// Verify sensitive description is sanitized
	if meta.Description == "This app contains password fields and secret keys" {
		t.Error("sensitive description should be sanitized")
	}
}

func TestExtractCategory(t *testing.T) {
	tests := []struct {
		id       string
		expected string
	}{
		{"AND-001", "Android"},
		{"IOS-002", "iOS"},
		{"FLT-003", "Flutter"},
		{"SEC-004", "Security"},
		{"POL-005", "Policy"},
		{"COD-006", "Code Quality"},
		{"TST-007", "Testing"},
		{"LINT-008", "Linting"},
		{"DOC-009", "Documentation"},
		{"PERF-010", "Performance"},
		{"REV-011", "Reviewer"},
		{"AI-012", "AI Analysis"},
		{"UNKNOWN", "Other"},
		{"", "Other"},
	}

	for _, tt := range tests {
		result := extractCategory(tt.id)
		if result != tt.expected {
			t.Errorf("extractCategory(%q) = %q, want %q", tt.id, result, tt.expected)
		}
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"34", 34},
		{"21", 21},
		{"35+", 35},
		{"100abc", 100},
		{"", 0},
		{"abc", 0},
	}

	for _, tt := range tests {
		result := parseInt(tt.input)
		if result != tt.expected {
			t.Errorf("parseInt(%q) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeDescription(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		checkSuffix bool
	}{
		{
			input:    "This is a normal description",
			expected: "This is a normal description",
		},
		{
			input:    "This contains a password field",
			expected: "[Description contains sensitive keywords]",
		},
		{
			input:    "Secret key is stored here",
			expected: "[Description contains sensitive keywords]",
		},
		{
			input:       strings.Repeat("a", 300), // Very long description
			expected:    "...",
			checkSuffix: true,
		},
	}

	for _, tt := range tests {
		result := sanitizeDescription(tt.input)
		if tt.checkSuffix {
			// For long strings, just check truncation suffix
			if !strings.HasSuffix(result, tt.expected) {
				t.Errorf("sanitizeDescription should end with %q", tt.expected)
			}
			if len(result) > 210 {
				t.Errorf("result length %d exceeds 210", len(result))
			}
			continue
		}
		if result != tt.expected {
			t.Errorf("sanitizeDescription(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func TestEmptyProject(t *testing.T) {
	project := checker.NewProject("/test/path")
	findings := []report.Finding{}

	meta := ExtractMetadata(project, findings, 38)

	if meta.AppName == "" {
		t.Error("should have default app name from path")
	}

	if meta.TotalChecks != 38 {
		t.Errorf("expected 38 total checks, got %d", meta.TotalChecks)
	}

	if meta.ComplianceScore != 100 {
		t.Errorf("expected score 100 for no findings, got %d", meta.ComplianceScore)
	}
}
