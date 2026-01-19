package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAndroidManifest(t *testing.T) {
	testdataDir := getTestdataDir(t)

	t.Run("valid_manifest.xml", func(t *testing.T) {
		path := filepath.Join(testdataDir, "android", "valid_manifest.xml")
		manifest, err := ParseAndroidManifest(path)
		if err != nil {
			t.Fatalf("Failed to parse manifest: %v", err)
		}

		if manifest.Package != "com.example.myapp" {
			t.Errorf("Expected package 'com.example.myapp', got '%s'", manifest.Package)
		}

		if manifest.VersionCode != "1" {
			t.Errorf("Expected versionCode '1', got '%s'", manifest.VersionCode)
		}

		if manifest.VersionName != "1.0.0" {
			t.Errorf("Expected versionName '1.0.0', got '%s'", manifest.VersionName)
		}

		if manifest.GetDebuggable() {
			t.Error("Expected debuggable to be false")
		}

		if !manifest.GetAllowBackup() {
			t.Error("Expected allowBackup to be true")
		}

		if len(manifest.UsesPermissions) != 4 {
			t.Errorf("Expected 4 permissions, got %d", len(manifest.UsesPermissions))
		}

		if !manifest.HasInternetPermission() {
			t.Error("Expected internet permission")
		}

		if !manifest.HasCameraPermission() {
			t.Error("Expected camera permission")
		}

		if !manifest.HasLocationPermission() {
			t.Error("Expected location permission")
		}

		activities := manifest.GetActivitiesWithIntentFilters()
		if len(activities) != 3 {
			t.Errorf("Expected 3 activities with intent filters, got %d", len(activities))
		}

		needsExported := manifest.ActivityNeedsExported()
		if len(needsExported) != 1 {
			t.Errorf("Expected 1 activity needing exported, got %d", len(needsExported))
		}

		if !manifest.HasQueryPackage("com.google.android.apps.maps") {
			t.Error("Expected query package")
		}
	})

	t.Run("invalid_manifest.xml", func(t *testing.T) {
		path := filepath.Join(testdataDir, "android", "invalid_manifest.xml")
		manifest, err := ParseAndroidManifest(path)
		if err != nil {
			t.Fatalf("Failed to parse manifest: %v", err)
		}

		if manifest.Package != "com.debug.app" {
			t.Errorf("Expected package 'com.debug.app', got '%s'", manifest.Package)
		}

		if !manifest.GetDebuggable() {
			t.Error("Expected debuggable to be true")
		}

		needsExported := manifest.ActivityNeedsExported()
		if len(needsExported) != 1 {
			t.Errorf("Expected 1 activity needing exported, got %d", len(needsExported))
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := ParseAndroidManifest("/nonexistent/path/manifest.xml")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})
}

func TestParseGradleFile(t *testing.T) {
	testdataDir := getTestdataDir(t)

	t.Run("build.gradle", func(t *testing.T) {
		path := filepath.Join(testdataDir, "android", "build.gradle")
		config, err := ParseGradleFile(path)
		if err != nil {
			t.Fatalf("Failed to parse gradle file: %v", err)
		}

		if config.ApplicationID != "com.example.myapp" {
			t.Errorf("Expected applicationId 'com.example.myapp', got '%s'", config.ApplicationID)
		}

		if config.MinSDKVersion != "21" {
			t.Errorf("Expected minSdkVersion '21', got '%s'", config.MinSDKVersion)
		}

		if config.TargetSDKVersion != "34" {
			t.Errorf("Expected targetSdkVersion '34', got '%s'", config.TargetSDKVersion)
		}

		if config.VersionCode != "1" {
			t.Errorf("Expected versionCode '1', got '%s'", config.VersionCode)
		}

		if config.VersionName != "1.0.0" {
			t.Errorf("Expected versionName '1.0.0', got '%s'", config.VersionName)
		}

		if config.NdkVersion != "25.2.9519653" {
			t.Errorf("Expected ndkVersion '25.2.9519653', got '%s'", config.NdkVersion)
		}
	})

	t.Run("build.gradle.kts", func(t *testing.T) {
		path := filepath.Join(testdataDir, "android", "build.gradle.kts")
		config, err := ParseGradleKtsFile(path)
		if err != nil {
			t.Fatalf("Failed to parse gradle.kts file: %v", err)
		}

		if config.ApplicationID != "com.example.myapp" {
			t.Errorf("Expected applicationId 'com.example.myapp', got '%s'", config.ApplicationID)
		}

		if config.MinSDKVersion != "21" {
			t.Errorf("Expected minSdkVersion '21', got '%s'", config.MinSDKVersion)
		}

		if config.TargetSDKVersion != "34" {
			t.Errorf("Expected targetSdkVersion '34', got '%s'", config.TargetSDKVersion)
		}

		if config.VersionCode != "1" {
			t.Errorf("Expected versionCode '1', got '%s'", config.VersionCode)
		}

		if config.VersionName != "1.0.0" {
			t.Errorf("Expected versionName '1.0.0', got '%s'", config.VersionName)
		}
	})
}

func getTestdataDir(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	dir := cwd
	for i := 0; i < 5; i++ {
		testdataPath := filepath.Join(dir, "testdata")
		if _, err := os.Stat(testdataPath); err == nil {
			return testdataPath
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatalf("Could not find testdata directory from %s", cwd)
	return ""
}
