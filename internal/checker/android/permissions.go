package android

import (
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type InternetPermissionCheck struct{}

func (c *InternetPermissionCheck) ID() string {
	return "AND-003"
}

func (c *InternetPermissionCheck) Name() string {
	return "Internet Permission Check"
}

func (c *InternetPermissionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if !project.HasNetworkDeps {
		return findings
	}

	hasInternet := false
	if project.AndroidManifest != nil {
		for _, perm := range project.AndroidManifest.Permissions {
			if strings.Contains(perm, "INTERNET") {
				hasInternet = true
				break
			}
		}
	}

	if !hasInternet {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses network dependencies but does not have INTERNET permission declared in AndroidManifest.xml",
			"android/app/src/main/AndroidManifest.xml",
			"Add <uses-permission android:name=\"android.permission.INTERNET\" /> to AndroidManifest.xml",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type DangerousPermissionsCheck struct{}

func (c *DangerousPermissionsCheck) ID() string {
	return "AND-004"
}

func (c *DangerousPermissionsCheck) Name() string {
	return "Dangerous Permissions Check"
}

func (c *DangerousPermissionsCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.AndroidManifest == nil {
		return findings
	}

	dangerousPermissions := map[string]string{
		"CAMERA":                 "android.hardware.camera",
		"ACCESS_FINE_LOCATION":   "android.hardware.location.gps",
		"ACCESS_COARSE_LOCATION": "android.hardware.location",
		"RECORD_AUDIO":           "android.hardware.microphone",
	}

	for perm, feature := range dangerousPermissions {
		hasPermission := false
		hasFeature := false

		for _, p := range project.AndroidManifest.Permissions {
			if strings.Contains(p, perm) {
				hasPermission = true
				break
			}
		}

		for _, activity := range project.AndroidManifest.Activities {
			if strings.Contains(activity.Name, feature) {
				hasFeature = true
				break
			}
		}

		if hasPermission && !hasFeature {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"App uses "+perm+" permission without corresponding uses-feature declaration",
				"android/app/src/main/AndroidManifest.xml",
				"Add <uses-feature android:name=\""+feature+"\" android:required=\"false\" /> to AndroidManifest.xml",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type PackageVisibilityCheck struct{}

func (c *PackageVisibilityCheck) ID() string {
	return "AND-011"
}

func (c *PackageVisibilityCheck) Name() string {
	return "Package Visibility Check"
}

func (c *PackageVisibilityCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if !project.HasURLLauncher {
		return findings
	}

	hasQueries := false
	if project.AndroidManifest != nil && len(project.AndroidManifest.QueriesPackages) > 0 {
		hasQueries = true
	}

	if !hasQueries {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses url_launcher but does not have <queries> declaration in AndroidManifest.xml for package visibility",
			"android/app/src/main/AndroidManifest.xml",
			"Add <queries> element with the packages your app needs to query",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}
