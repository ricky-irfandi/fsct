package parser

import (
	"path/filepath"
	"testing"
)

func TestParsePubspec(t *testing.T) {
	testdataDir := getTestdataDir(t)

	t.Run("pubspec.yaml", func(t *testing.T) {
		path := filepath.Join(testdataDir, "flutter", "pubspec.yaml")
		pubspec, err := ParsePubspec(path)
		if err != nil {
			t.Fatalf("Failed to parse pubspec.yaml: %v", err)
		}

		if pubspec.Name != "my_app" {
			t.Errorf("Expected name 'my_app', got '%s'", pubspec.Name)
		}

		if pubspec.Version != "1.0.0+1" {
			t.Errorf("Expected version '1.0.0+1', got '%s'", pubspec.Version)
		}

		if pubspec.Description != "A new Flutter project." {
			t.Errorf("Expected description 'A new Flutter project.', got '%s'", pubspec.Description)
		}

		if pubspec.Homepage != "https://example.com" {
			t.Errorf("Expected homepage 'https://example.com', got '%s'", pubspec.Homepage)
		}

		if !pubspec.HasDependency("http") {
			t.Error("Expected http dependency")
		}

		if !pubspec.HasDevDependency("flutter_lints") {
			t.Error("Expected flutter_lints dev dependency")
		}

		if !pubspec.HasDependency("url_launcher") {
			t.Error("Expected url_launcher dependency")
		}

		if !pubspec.HasDependency("image_picker") {
			t.Error("Expected image_picker dependency")
		}

		if !pubspec.HasLinter() {
			t.Error("Expected to have linter")
		}

		if !pubspec.IsDefaultVersion() {
			t.Error("Expected version to be default (1.0.0+1)")
		}

		if !pubspec.HasDescription() {
			t.Error("Expected to have description")
		}

		if !pubspec.HasRepository() {
			t.Error("Expected to have repository/homepage")
		}

		if pubspec.GetSdkVersion() != "'>=3.0.0 <4.0.0'" && pubspec.GetSdkVersion() != ">=3.0.0 <4.0.0" {
			t.Errorf("Expected sdk version, got '%s'", pubspec.GetSdkVersion())
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := ParsePubspec("/nonexistent/path/pubspec.yaml")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})
}
