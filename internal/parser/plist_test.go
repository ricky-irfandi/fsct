package parser

import (
	"path/filepath"
	"testing"
)

func TestParseInfoPlist(t *testing.T) {
	testdataDir := getTestdataDir(t)

	t.Run("Info.plist", func(t *testing.T) {
		path := filepath.Join(testdataDir, "ios", "Info.plist")
		plist, err := ParseInfoPlist(path)
		if err != nil {
			t.Fatalf("Failed to parse Info.plist: %v", err)
		}

		if plist.CFBundleIdentifier != "com.example.myapp" {
			t.Errorf("Expected CFBundleIdentifier 'com.example.myapp', got '%s'", plist.CFBundleIdentifier)
		}

		if plist.CFBundleVersion != "1" {
			t.Errorf("Expected CFBundleVersion '1', got '%s'", plist.CFBundleVersion)
		}

		if plist.CFBundleShortVersionString != "1.0.0" {
			t.Errorf("Expected CFBundleShortVersionString '1.0.0', got '%s'", plist.CFBundleShortVersionString)
		}

		if !plist.HasCameraUsageDescription() {
			t.Error("Expected camera usage description")
		}

		if !plist.HasPhotoLibraryUsageDescription() {
			t.Error("Expected photo library usage description")
		}

		if !plist.HasLocationUsageDescription() {
			t.Error("Expected location usage description")
		}

		if !plist.HasMicrophoneUsageDescription() {
			t.Error("Expected microphone usage description")
		}

		if !plist.HasContactsUsageDescription() {
			t.Error("Expected contacts usage description")
		}

		if !plist.HasCalendarsUsageDescription() {
			t.Error("Expected calendars usage description")
		}

		if !plist.IsEncryptionDeclarationSet() {
			t.Error("Expected encryption declaration to be set")
		}

		if !plist.IsEncryptionExempt() {
			t.Error("Expected encryption to be exempt")
		}

		if plist.GetFullScreenRequirement() {
			t.Error("Expected full screen requirement to be false")
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := ParseInfoPlist("/nonexistent/path/Info.plist")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})
}
