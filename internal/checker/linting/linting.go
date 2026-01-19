package linting

import (
	"regexp"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type AnalysisOptionsCheck struct{}

func (c *AnalysisOptionsCheck) ID() string {
	return "LINT-001"
}

func (c *AnalysisOptionsCheck) Name() string {
	return "Analysis Options File Presence"
}

func (c *AnalysisOptionsCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasAnalysisOptions := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "analysis_options.yaml") || strings.Contains(file, "analysis_options.yml") {
			hasAnalysisOptions = true
			break
		}
	}

	if !hasAnalysisOptions {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No analysis_options.yaml found",
			"project root",
			"Create analysis_options.yaml for lint rules configuration",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type LinterRulesCheck struct{}

func (c *LinterRulesCheck) ID() string {
	return "LINT-002"
}

func (c *LinterRulesCheck) Name() string {
	return "Linter Rules Configuration"
}

func (c *LinterRulesCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasLinterRules := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "analysis_options") {
			if strings.Contains(file, "linter:") || strings.Contains(file, "rules:") {
				hasLinterRules = true
				break
			}
		}
	}

	if !hasLinterRules {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No linter rules configured",
			"analysis_options.yaml",
			"Add linter rules to enforce code quality standards",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type StrongModeCheck struct{}

func (c *StrongModeCheck) ID() string {
	return "LINT-003"
}

func (c *StrongModeCheck) Name() string {
	return "Strong Mode Analysis"
}

func (c *StrongModeCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasStrongMode := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "analysis_options") {
			if strings.Contains(file, "strong-mode") || strings.Contains(file, "implicit-casts") {
				hasStrongMode = true
				break
			}
		}
	}

	if !hasStrongMode {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Strong mode analysis not explicitly configured",
			"analysis_options.yaml",
			"Enable strong mode for better type safety",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}

type FileNamingCheck struct{}

func (c *FileNamingCheck) ID() string {
	return "LINT-004"
}

func (c *FileNamingCheck) Name() string {
	return "File Naming Rules"
}

func (c *FileNamingCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasFileNamingRules := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "analysis_options") {
			if strings.Contains(file, "file_names") || strings.Contains(file, "camel_case_types") {
				hasFileNamingRules = true
				break
			}
		}
	}

	if !hasFileNamingRules {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"File naming rules not configured",
			"analysis_options.yaml",
			"Add file naming rules (camel_case_types, library_names)",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}

type StyleGuideCheck struct{}

func (c *StyleGuideCheck) ID() string {
	return "LINT-005"
}

func (c *StyleGuideCheck) Name() string {
	return "Style Guide Rules"
}

func (c *StyleGuideCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	styleRules := []string{"lines_longer_than_80_chars", "avoid_as", "prefer_const_constructors"}
	hasStyleRules := false

	for _, file := range project.DartFiles {
		if strings.Contains(file, "analysis_options") {
			for _, rule := range styleRules {
				if strings.Contains(file, rule) {
					hasStyleRules = true
					break
				}
			}
		}
	}

	if !hasStyleRules {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Style guide rules not fully configured",
			"analysis_options.yaml",
			"Add style rules like prefer_const_constructors, avoid_as",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}

type PublicAPIDocCheck struct{}

func (c *PublicAPIDocCheck) ID() string {
	return "LINT-006"
}

func (c *PublicAPIDocCheck) Name() string {
	return "Public API Documentation Rules"
}

func (c *PublicAPIDocCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasDocRules := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "analysis_options") {
			if strings.Contains(file, "public_member_api_docs") || strings.Contains(file, "comment_references") {
				hasDocRules = true
				break
			}
		}
	}

	if !hasDocRules {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Public API documentation rules not configured",
			"analysis_options.yaml",
			"Add public_member_api_docs rule for documentation enforcement",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}

type IgnoreCommentsCheck struct{}

func (c *IgnoreCommentsCheck) ID() string {
	return "LINT-007"
}

func (c *IgnoreCommentsCheck) Name() string {
	return "Ignore Comments Configuration"
}

func (c *IgnoreCommentsCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	ignoreCount := 0
	ignorePattern := regexp.MustCompile(`//\s*ignore:`)

	for _, file := range project.DartFiles {
		ignoreCount += len(ignorePattern.FindAllString(file, -1))
	}

	if ignoreCount > 10 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Many ignore comments found in project",
			"source files",
			"Review and address lint violations instead of ignoring them",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}
