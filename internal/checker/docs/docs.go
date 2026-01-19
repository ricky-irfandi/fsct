package docs

import (
	"regexp"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type ReadmePresenceCheck struct{}

func (c *ReadmePresenceCheck) ID() string {
	return "DOC-001"
}

func (c *ReadmePresenceCheck) Name() string {
	return "README.md Presence"
}

func (c *ReadmePresenceCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasReadme := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "README.md") || strings.Contains(file, "README") {
			hasReadme = true
			break
		}
	}

	if !hasReadme {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"README.md not found",
			"project root",
			"Add README.md with project description, setup instructions, and usage examples",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type ReadmeContentCheck struct{}

func (c *ReadmeContentCheck) ID() string {
	return "DOC-002"
}

func (c *ReadmeContentCheck) Name() string {
	return "README.md Content Quality"
}

func (c *ReadmeContentCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	readmeContentPattern := regexp.MustCompile(`(?i)(install|setup|usage|example|feature)`)
	hasContent := false

	for _, file := range project.DartFiles {
		if strings.Contains(file, "README") && readmeContentPattern.MatchString(file) {
			hasContent = true
			break
		}
	}

	if !hasContent {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"README.md may be missing content",
			"README.md",
			"Add sections for installation, usage, and examples",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}

type ChangelogPresenceCheck struct{}

func (c *ChangelogPresenceCheck) ID() string {
	return "DOC-003"
}

func (c *ChangelogPresenceCheck) Name() string {
	return "CHANGELOG.md Presence"
}

func (c *ChangelogPresenceCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasChangelog := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "CHANGELOG") || strings.Contains(file, "CHANGES") {
			hasChangelog = true
			break
		}
	}

	if !hasChangelog {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"CHANGELOG.md not found",
			"project root",
			"Add CHANGELOG.md to track version changes",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}

type LicensePresenceCheck struct{}

func (c *LicensePresenceCheck) ID() string {
	return "DOC-004"
}

func (c *LicensePresenceCheck) Name() string {
	return "LICENSE File Presence"
}

func (c *LicensePresenceCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasLicense := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "LICENSE") || strings.Contains(file, "COPYING") {
			hasLicense = true
			break
		}
	}

	if !hasLicense {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"LICENSE file not found",
			"project root",
			"Add LICENSE file (MIT, Apache 2.0, etc.)",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type ApiDocumentationCheck struct{}

func (c *ApiDocumentationCheck) ID() string {
	return "DOC-005"
}

func (c *ApiDocumentationCheck) Name() string {
	return "API Documentation"
}

func (c *ApiDocumentationCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	docPattern := regexp.MustCompile(`///`)
	undocumentedCount := 0

	for _, file := range project.DartFiles {
		if !strings.Contains(file, "test/") {
			if !docPattern.MatchString(file) {
				undocumentedCount++
			}
		}
	}

	if undocumentedCount > 10 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Many files lack documentation comments",
			"lib/",
			"Add /// documentation comments for public APIs",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type CodeCommentsCheck struct{}

func (c *CodeCommentsCheck) ID() string {
	return "DOC-006"
}

func (c *CodeCommentsCheck) Name() string {
	return "Code Comment Quality"
}

func (c *CodeCommentsCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasIssues := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "TODO") || strings.Contains(file, "FIXME") {
			hasIssues = true
			break
		}
	}

	if hasIssues {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"TODO/FIXME comments found in code",
			"source files",
			"Review and address TODO items or create issues",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}
