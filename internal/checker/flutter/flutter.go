package flutter

import (
	"strconv"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type FlutterSDKVersionCheck struct{}

func (c *FlutterSDKVersionCheck) ID() string {
	return "FLT-001"
}

func (c *FlutterSDKVersionCheck) Name() string {
	return "Flutter SDK Version Constraint"
}

func (c *FlutterSDKVersionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.Pubspec == nil || project.Pubspec.Dependencies == nil {
		return findings
	}

	hasFlutterSDK := false
	for dep := range project.Pubspec.Dependencies {
		if strings.Contains(strings.ToLower(dep), "flutter") {
			hasFlutterSDK = true
			break
		}
	}

	if !hasFlutterSDK {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Flutter SDK version constraint not found in pubspec.yaml",
			"pubspec.yaml",
			"Add sdk: flutter constraint to dependencies section",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type Material3Check struct{}

func (c *Material3Check) ID() string {
	return "FLT-002"
}

func (c *Material3Check) Name() string {
	return "Flutter Use Material Design 3"
}

func (c *Material3Check) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding
	return findings
}

type MinSDKVersionCheck struct{}

func (c *MinSDKVersionCheck) ID() string {
	return "FLT-003"
}

func (c *MinSDKVersionCheck) Name() string {
	return "Flutter Min SDK Version"
}

func (c *MinSDKVersionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.GradleConfig != nil && project.GradleConfig.MinSDKVersion != "" {
		minSDK := project.GradleConfig.MinSDKVersion
		minSDKInt, err := strconv.Atoi(minSDK)
		if err != nil {
			return findings
		}
		if minSDKInt < 21 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Minimum SDK version "+minSDK+" is too low (recommended: 21+)",
				"android/app/build.gradle",
				"Consider raising minSdkVersion to 21 or higher for better compatibility",
				report.SeverityHigh,
				0,
			))
		}
	}

	return findings
}

type PackageNameCheck struct{}

func (c *PackageNameCheck) ID() string {
	return "FLT-004"
}

func (c *PackageNameCheck) Name() string {
	return "Flutter Package Name Validation"
}

func (c *PackageNameCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.Pubspec != nil && project.Pubspec.Name != "" {
		if !strings.Contains(project.Pubspec.Name, ".") {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Invalid package name format: "+project.Pubspec.Name,
				"pubspec.yaml",
				"Package name should follow reverse domain notation (e.g., com.example.app)",
				report.SeverityHigh,
				0,
			))
		}
	}

	return findings
}

type VersionCheck struct{}

func (c *VersionCheck) ID() string {
	return "FLT-005"
}

func (c *VersionCheck) Name() string {
	return "Flutter Version Management"
}

func (c *VersionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.Pubspec != nil && project.Pubspec.Version == "" {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Version not specified in pubspec.yaml",
			"pubspec.yaml",
			"Add version field in format: 1.0.0+1",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type DependencyConstraintCheck struct{}

func (c *DependencyConstraintCheck) ID() string {
	return "FLT-006"
}

func (c *DependencyConstraintCheck) Name() string {
	return "Flutter Dependency Version Constraints"
}

func (c *DependencyConstraintCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.Pubspec != nil && project.Pubspec.Dependencies != nil {
		unpinnedCount := 0
		for _, version := range project.Pubspec.Dependencies {
			if version == "any" || version == "" {
				unpinnedCount++
			}
		}
		if unpinnedCount > 3 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Multiple dependencies without version constraints",
				"pubspec.yaml",
				"Use version constraints (^) to ensure reproducible builds",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type DeprecatedPackageCheck struct{}

func (c *DeprecatedPackageCheck) ID() string {
	return "FLT-007"
}

func (c *DeprecatedPackageCheck) Name() string {
	return "Flutter Deprecated Package Usage"
}

func (c *DeprecatedPackageCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.Pubspec != nil && project.Pubspec.HasDeprecatedPkg {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Potentially deprecated packages found",
			"pubspec.yaml",
			"Check for latest FlutterFire packages and migration guides",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type ProjectStructureCheck struct{}

func (c *ProjectStructureCheck) ID() string {
	return "FLT-008"
}

func (c *ProjectStructureCheck) Name() string {
	return "Flutter Project Structure"
}

func (c *ProjectStructureCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding
	return findings
}
