package ios

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

func TestCameraUsageDescriptionCheck(t *testing.T) {
	check := &CameraUsageDescriptionCheck{}

	t.Run("camera deps without description should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			HasCameraDeps: true,
			InfoPlist: &checker.InfoPlistInfo{
				HasCameraUsageDescription: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}

		if len(findings) > 0 && findings[0].Severity != report.SeverityHigh {
			t.Errorf("Expected HIGH severity, got %s", findings[0].Severity)
		}
	})

	t.Run("camera deps with description should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasCameraDeps: true,
			InfoPlist: &checker.InfoPlistInfo{
				HasCameraUsageDescription: true,
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})

	t.Run("no camera deps should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasCameraDeps: false,
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestPhotoLibraryUsageDescriptionCheck(t *testing.T) {
	check := &PhotoLibraryUsageDescriptionCheck{}

	t.Run("image_picker without description should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			HasImagePicker: true,
			InfoPlist: &checker.InfoPlistInfo{
				HasPhotoLibraryUsageDescription: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})

	t.Run("image_picker with description should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasImagePicker: true,
			InfoPlist: &checker.InfoPlistInfo{
				HasPhotoLibraryUsageDescription: true,
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestLocationUsageDescriptionCheck(t *testing.T) {
	check := &LocationUsageDescriptionCheck{}

	t.Run("location deps without description should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			HasLocationDeps: true,
			InfoPlist: &checker.InfoPlistInfo{
				HasLocationUsageDescription: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})

	t.Run("location deps with description should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasLocationDeps: true,
			InfoPlist: &checker.InfoPlistInfo{
				HasLocationUsageDescription: true,
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestEncryptionDeclarationCheck(t *testing.T) {
	check := &EncryptionDeclarationCheck{}

	t.Run("missing encryption declaration should generate WARNING", func(t *testing.T) {
		project := &checker.Project{
			InfoPlist: &checker.InfoPlistInfo{
				EncryptionDeclarationSet: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}

		if len(findings) > 0 && findings[0].Severity != report.SeverityWarning {
			t.Errorf("Expected WARNING severity, got %s", findings[0].Severity)
		}
	})

	t.Run("encryption declaration set should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			InfoPlist: &checker.InfoPlistInfo{
				EncryptionDeclarationSet: true,
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestFullScreenConflictCheck(t *testing.T) {
	check := &FullScreenConflictCheck{}

	t.Run("full screen true should generate WARNING", func(t *testing.T) {
		project := &checker.Project{
			InfoPlist: &checker.InfoPlistInfo{
				RequiresFullScreen: true,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}

		if len(findings) > 0 && findings[0].Severity != report.SeverityWarning {
			t.Errorf("Expected WARNING severity, got %s", findings[0].Severity)
		}
	})

	t.Run("full screen false should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			InfoPlist: &checker.InfoPlistInfo{
				RequiresFullScreen: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestEmptyUsageDescriptionCheck(t *testing.T) {
	check := &EmptyUsageDescriptionCheck{}

	t.Run("check is disabled - no findings", func(t *testing.T) {
		project := &checker.Project{
			InfoPlist: &checker.InfoPlistInfo{
				CFBundleIdentifier: "Add reason here for camera access",
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings (check disabled), got %d", len(findings))
		}
	})

	t.Run("no findings for any input", func(t *testing.T) {
		project := &checker.Project{
			InfoPlist: &checker.InfoPlistInfo{
				CFBundleIdentifier: "We need camera access to take photos of your documents",
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings (check disabled), got %d", len(findings))
		}
	})
}

func TestMissingAppIconCheck(t *testing.T) {
	check := &MissingAppIconCheck{}

	t.Run("missing AppIcon folder should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			IOSPath: "/nonexistent",
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}

		if len(findings) > 0 && findings[0].Severity != report.SeverityHigh {
			t.Errorf("Expected HIGH severity, got %s", findings[0].Severity)
		}
	})
}

func TestDeploymentTargetCheck(t *testing.T) {
	check := &DeploymentTargetCheck{}

	t.Run("deployment target less than 12 should generate WARNING", func(t *testing.T) {
		project := &checker.Project{
			IOSPath: "/nonexistent",
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings for nonexistent path, got %d", len(findings))
		}
	})
}

func TestMicrophoneUsageDescriptionCheck(t *testing.T) {
	check := &MicrophoneUsageDescriptionCheck{}

	t.Run("mic deps without description should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Dependencies: map[string]string{
					"flutter_sound": "^9.0.0",
				},
			},
			InfoPlist: &checker.InfoPlistInfo{
				HasMicrophoneUsageDescription: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})

	t.Run("mic deps with description should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Dependencies: map[string]string{
					"flutter_sound": "^9.0.0",
				},
			},
			InfoPlist: &checker.InfoPlistInfo{
				HasMicrophoneUsageDescription: true,
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestContactsUsageDescriptionCheck(t *testing.T) {
	check := &ContactsUsageDescriptionCheck{}

	t.Run("contacts deps without description should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Dependencies: map[string]string{
					"contacts_service": "^0.6.0",
				},
			},
			InfoPlist: &checker.InfoPlistInfo{
				HasContactsUsageDescription: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})
}

func TestCalendarsUsageDescriptionCheck(t *testing.T) {
	check := &CalendarsUsageDescriptionCheck{}

	t.Run("calendar deps without description should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Dependencies: map[string]string{
					"event_calendar": "^1.0.0",
				},
			},
			InfoPlist: &checker.InfoPlistInfo{
				HasCalendarsUsageDescription: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})
}
