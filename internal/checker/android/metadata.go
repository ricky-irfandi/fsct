package android

import (
	"strconv"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type DebuggableCheck struct{}

func (c *DebuggableCheck) ID() string {
	return "AND-005"
}

func (c *DebuggableCheck) Name() string {
	return "Debuggable Check"
}

func (c *DebuggableCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.AndroidManifest == nil {
		return findings
	}

	if project.AndroidManifest.Debuggable {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"android:debuggable is set to true. This should be false for release builds.",
			"android/app/src/main/AndroidManifest.xml",
			"Set android:debuggable=\"false\" or remove the attribute",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type ExportedAttributeCheck struct{}

func (c *ExportedAttributeCheck) ID() string {
	return "AND-006"
}

func (c *ExportedAttributeCheck) Name() string {
	return "Exported Attribute Check"
}

func (c *ExportedAttributeCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.AndroidManifest == nil {
		return findings
	}

	for _, activity := range project.AndroidManifest.Activities {
		if activity.HasIntentFilter && !activity.Exported {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Activity "+activity.Name+" has intent-filter but android:exported is not explicitly set",
				"android/app/src/main/AndroidManifest.xml",
				"Set android:exported=\"true\" or android:exported=\"false\" for the activity",
				report.SeverityHigh,
				0,
			))
		}
	}

	return findings
}

type ApplicationIDCheck struct{}

func (c *ApplicationIDCheck) ID() string {
	return "AND-009"
}

func (c *ApplicationIDCheck) Name() string {
	return "Application ID Check"
}

func (c *ApplicationIDCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.GradleConfig == nil || project.GradleConfig.ApplicationID == "" {
		return findings
	}

	if strings.HasPrefix(project.GradleConfig.ApplicationID, "com.example.") {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Application ID starts with com.example. Google Play Store requires a unique, valid package name.",
			"android/app/build.gradle",
			"Change applicationId to a unique package name (e.g., com.yourcompany.yourapp)",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type VersionCodeCheck struct{}

func (c *VersionCodeCheck) ID() string {
	return "AND-010"
}

func (c *VersionCodeCheck) Name() string {
	return "Version Code Check"
}

func (c *VersionCodeCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.GradleConfig == nil || project.GradleConfig.VersionCode == "" {
		return findings
	}

	versionCode, err := strconv.Atoi(project.GradleConfig.VersionCode)
	if err != nil {
		return findings
	}

	if versionCode == 1 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"versionCode is 1. Consider incrementing the version code for updates.",
			"android/app/build.gradle",
			"Increment versionCode for subsequent releases",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type AllowBackupCheck struct{}

func (c *AllowBackupCheck) ID() string {
	return "AND-012"
}

func (c *AllowBackupCheck) Name() string {
	return "Allow Backup Check"
}

func (c *AllowBackupCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.AndroidManifest == nil {
		return findings
	}

	if project.AndroidManifest.AllowBackup {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"android:allowBackup is set to true. Consider disabling if app handles sensitive data.",
			"android/app/src/main/AndroidManifest.xml",
			"Set android:allowBackup=\"false\" or implement encryption for sensitive data",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}
