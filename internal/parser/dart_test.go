package parser

import (
	"path/filepath"
	"testing"
)

func TestDartScanner(t *testing.T) {
	testdataDir := getTestdataDir(t)
	libDir := filepath.Join(testdataDir, "flutter", "lib")

	t.Run("scan for print statements", func(t *testing.T) {
		matches, err := FindPrintStatements(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) < 2 {
			t.Errorf("Expected at least 2 print statements, got %d", len(matches))
		}

		found := false
		for _, m := range matches {
			if m.Line == 20 || m.Line == 28 || m.Line == 55 || m.Line == 59 || m.Line == 63 || m.Line == 67 {
				found = true
				break
			}
		}
		if !found {
			t.Logf("All matches: %+v", matches)
			t.Error("Expected to find print statements")
		}
	})

	t.Run("scan for api key", func(t *testing.T) {
		matches, err := FindAPIKeys(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find api key")
		}

		found := false
		for _, m := range matches {
			if m.File == "sample.dart" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected to find api key in sample.dart")
		}
	})

	t.Run("scan for passwords", func(t *testing.T) {
		matches, err := FindHardcodedPasswords(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find hardcoded password")
		}
	})

	t.Run("scan for http urls", func(t *testing.T) {
		matches, err := FindHTTPURLs(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find http url")
		}
	})

	t.Run("scan for privacy patterns", func(t *testing.T) {
		matches, err := FindPrivacyPatterns(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find privacy patterns")
		}
	})

	t.Run("scan for logout patterns", func(t *testing.T) {
		matches, err := FindLogoutPatterns(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find logout patterns")
		}
	})

	t.Run("scan for delete account patterns", func(t *testing.T) {
		matches, err := FindDeleteAccountPatterns(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find delete account patterns")
		}
	})

	t.Run("scan for password recovery patterns", func(t *testing.T) {
		matches, err := FindPasswordRecoveryPatterns(libDir)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find password recovery patterns")
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		matches, err := FindPrintStatements("/nonexistent")
		if err != nil {
			t.Logf("Got error (expected): %v", err)
		}
		if len(matches) != 0 {
			t.Errorf("Expected 0 matches for nonexistent directory, got %d", len(matches))
		}
	})
}

func TestDartScannerFindPatternInFile(t *testing.T) {
	testdataDir := getTestdataDir(t)
	sampleFile := filepath.Join(testdataDir, "flutter", "lib", "sample.dart")

	t.Run("find specific pattern", func(t *testing.T) {
		scanner := NewDartScanner(filepath.Dir(sampleFile))
		matches, err := scanner.FindPatternInFile(sampleFile, `print\(.*\)`)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) < 2 {
			t.Errorf("Expected at least 2 print matches, got %d", len(matches))
		}
	})
}

func TestScanDartFiles(t *testing.T) {
	testdataDir := getTestdataDir(t)
	libDir := filepath.Join(testdataDir, "flutter", "lib")

	t.Run("scan with single pattern", func(t *testing.T) {
		matches, err := ScanDartFiles(libDir, `api[keyK]?`)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) == 0 {
			t.Error("Expected to find api key pattern")
		}
	})

	t.Run("scan with multiple patterns", func(t *testing.T) {
		patterns := []string{
			`print\(`,
			`debugPrint\(`,
		}
		matches, err := ScanDartFilesMulti(libDir, patterns)
		if err != nil {
			t.Fatalf("Failed to scan: %v", err)
		}

		if len(matches) < 2 {
			t.Errorf("Expected at least 2 matches, got %d", len(matches))
		}
	})
}
