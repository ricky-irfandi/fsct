package android

import (
	"strconv"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type TargetSDKCheck struct{}

func (c *TargetSDKCheck) ID() string {
	return "AND-001"
}

func (c *TargetSDKCheck) Name() string {
	return "Target SDK Version Check"
}

func (c *TargetSDKCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.GradleConfig == nil || project.GradleConfig.TargetSDKVersion == "" {
		return findings
	}

	targetSDK, err := strconv.Atoi(project.GradleConfig.TargetSDKVersion)
	if err != nil {
		return findings
	}

	if targetSDK < 35 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Target SDK version is "+project.GradleConfig.TargetSDKVersion+". Google Play Store requires targetSdkVersion 35+.",
			"android/app/build.gradle",
			"Update targetSdkVersion to 35 or higher",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type MinSDKCheck struct{}

func (c *MinSDKCheck) ID() string {
	return "AND-002"
}

func (c *MinSDKCheck) Name() string {
	return "Minimum SDK Version Check"
}

func (c *MinSDKCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.GradleConfig == nil || project.GradleConfig.MinSDKVersion == "" {
		return findings
	}

	minSDK, err := strconv.Atoi(project.GradleConfig.MinSDKVersion)
	if err != nil {
		return findings
	}

	if minSDK < 21 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Minimum SDK version is "+project.GradleConfig.MinSDKVersion+". Consider updating to API 21+ for better security and performance.",
			"android/app/build.gradle",
			"Update minSdkVersion to 21 or higher",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}
