package android

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

func TestTargetSDKCheck(t *testing.T) {
	check := &TargetSDKCheck{}

	t.Run("SDK 31 should generate HIGH finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				TargetSDKVersion: "31",
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}

		if len(findings) > 0 {
			if findings[0].Severity != report.SeverityHigh {
				t.Errorf("Expected HIGH severity, got %s", findings[0].Severity)
			}
			if findings[0].ID != "AND-001" {
				t.Errorf("Expected ID AND-001, got %s", findings[0].ID)
			}
		}
	})

	t.Run("SDK 34 should generate HIGH finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				TargetSDKVersion: "34",
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})

	t.Run("SDK 35 should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				TargetSDKVersion: "35",
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})

	t.Run("empty gradle config should not panic", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings for empty config, got %d", len(findings))
		}
	})
}

func TestMinSDKCheck(t *testing.T) {
	check := &MinSDKCheck{}

	t.Run("SDK 16 should generate WARNING finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				MinSDKVersion: "16",
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}

		if len(findings) > 0 {
			if findings[0].Severity != report.SeverityWarning {
				t.Errorf("Expected WARNING severity, got %s", findings[0].Severity)
			}
		}
	})

	t.Run("SDK 21 should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				MinSDKVersion: "21",
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestInternetPermissionCheck(t *testing.T) {
	check := &InternetPermissionCheck{}

	t.Run("network deps without permission should generate WARNING", func(t *testing.T) {
		project := &checker.Project{
			HasNetworkDeps: true,
			AndroidManifest: &checker.AndroidManifestInfo{
				Permissions: []string{},
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})

	t.Run("network deps with permission should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasNetworkDeps: true,
			AndroidManifest: &checker.AndroidManifestInfo{
				Permissions: []string{"android.permission.INTERNET"},
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})

	t.Run("no network deps should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasNetworkDeps: false,
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestDebuggableCheck(t *testing.T) {
	check := &DebuggableCheck{}

	t.Run("debuggable true should generate HIGH finding", func(t *testing.T) {
		project := &checker.Project{
			AndroidManifest: &checker.AndroidManifestInfo{
				Debuggable: true,
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

	t.Run("debuggable false should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			AndroidManifest: &checker.AndroidManifestInfo{
				Debuggable: false,
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestApplicationIDCheck(t *testing.T) {
	check := &ApplicationIDCheck{}

	t.Run("com.example should generate HIGH finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				ApplicationID: "com.example.myapp",
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

	t.Run("unique app ID should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				ApplicationID: "com.mycompany.myapp",
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestVersionCodeCheck(t *testing.T) {
	check := &VersionCodeCheck{}

	t.Run("versionCode 1 should generate WARNING finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				VersionCode: "1",
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

	t.Run("versionCode 2 should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				VersionCode: "2",
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestAllowBackupCheck(t *testing.T) {
	check := &AllowBackupCheck{}

	t.Run("allowBackup true should generate WARNING finding", func(t *testing.T) {
		project := &checker.Project{
			AndroidManifest: &checker.AndroidManifestInfo{
				AllowBackup: true,
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
}

func TestExportedAttributeCheck(t *testing.T) {
	check := &ExportedAttributeCheck{}

	t.Run("activity with intent-filter but no exported should generate HIGH", func(t *testing.T) {
		project := &checker.Project{
			AndroidManifest: &checker.AndroidManifestInfo{
				Activities: []checker.ActivityInfo{
					{
						Name:            ".MainActivity",
						Exported:        false,
						HasIntentFilter: true,
					},
				},
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})

	t.Run("activity with exported set should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			AndroidManifest: &checker.AndroidManifestInfo{
				Activities: []checker.ActivityInfo{
					{
						Name:            ".MainActivity",
						Exported:        true,
						HasIntentFilter: true,
					},
				},
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}

func TestPackageVisibilityCheck(t *testing.T) {
	check := &PackageVisibilityCheck{}

	t.Run("url_launcher without queries should generate WARNING", func(t *testing.T) {
		project := &checker.Project{
			HasURLLauncher: true,
			AndroidManifest: &checker.AndroidManifestInfo{
				QueriesPackages: []string{},
			},
		}

		findings := check.Run(project)

		if len(findings) != 1 {
			t.Errorf("Expected 1 finding, got %d", len(findings))
		}
	})

	t.Run("url_launcher with queries should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasURLLauncher: true,
			AndroidManifest: &checker.AndroidManifestInfo{
				QueriesPackages: []string{"com.google.android.apps.maps"},
			},
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})

	t.Run("no url_launcher should not generate finding", func(t *testing.T) {
		project := &checker.Project{
			HasURLLauncher: false,
		}

		findings := check.Run(project)

		if len(findings) != 0 {
			t.Errorf("Expected 0 findings, got %d", len(findings))
		}
	})
}
