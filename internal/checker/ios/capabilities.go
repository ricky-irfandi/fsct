package ios

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type FullScreenConflictCheck struct{}

func (c *FullScreenConflictCheck) ID() string {
	return "IOS-010"
}

func (c *FullScreenConflictCheck) Name() string {
	return "Full Screen Conflict Check"
}

func (c *FullScreenConflictCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.InfoPlist == nil {
		return findings
	}

	if project.InfoPlist.RequiresFullScreen {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"UIRequiresFullScreen is set to true. If your app supports iPad, this may cause App Store rejection.",
			"ios/Runner/Info.plist",
			"Either set UIRequiresFullScreen to false or ensure iPad is not in the target device list",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type EncryptionDeclarationCheck struct{}

func (c *EncryptionDeclarationCheck) ID() string {
	return "IOS-011"
}

func (c *EncryptionDeclarationCheck) Name() string {
	return "Encryption Declaration Check"
}

func (c *EncryptionDeclarationCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.InfoPlist == nil {
		return findings
	}

	if !project.InfoPlist.EncryptionDeclarationSet {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"ITSAppUsesNonExemptEncryption is not set in Info.plist. Apple requires this declaration.",
			"ios/Runner/Info.plist",
			"Add ITSAppUsesNonExemptEncryption and set it to false if not using encryption, or true and provide export compliance",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type DeploymentTargetCheck struct{}

func (c *DeploymentTargetCheck) ID() string {
	return "IOS-012"
}

func (c *DeploymentTargetCheck) Name() string {
	return "Deployment Target Check"
}

func (c *DeploymentTargetCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	iosPath := project.IOSPath
	if iosPath == "" {
		return findings
	}

	xcodeprojPath := filepath.Join(iosPath, "Runner.xcodeproj", "project.pbxproj")

	content, err := os.ReadFile(xcodeprojPath)
	if err != nil {
		return findings
	}

	pattern := regexp.MustCompile(`IPHONEOS_DEPLOYMENT_TARGET\s*=\s*([0-9.]+)`)
	matches := pattern.FindStringSubmatch(string(content))

	if len(matches) < 2 {
		return findings
	}

	deploymentTarget := matches[1]

	pattern = regexp.MustCompile(`^(\d+)\.(\d+)`)
	matches = pattern.FindStringSubmatch(deploymentTarget)

	if len(matches) < 3 {
		return findings
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])

	version := float64(major) + float64(minor)*0.1

	if version < 12.0 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"IPHONEOS_DEPLOYMENT_TARGET is less than 12.0. Consider updating to support modern iOS versions.",
			"ios/Runner.xcodeproj/project.pbxproj",
			"Update IPHONEOS_DEPLOYMENT_TARGET to 12.0 or higher",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}
