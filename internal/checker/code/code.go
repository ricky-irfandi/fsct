package code

import (
	"regexp"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type FileLengthCheck struct{}

func (c *FileLengthCheck) ID() string {
	return "COD-001"
}

func (c *FileLengthCheck) Name() string {
	return "File Length Check"
}

func (c *FileLengthCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	for _, file := range project.DartFiles {
		lines := strings.Count(file, "\n")
		if lines > 400 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"File exceeds 400 lines: "+file,
				file,
				"Consider splitting this file into smaller, focused modules",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type ClassLengthCheck struct{}

func (c *ClassLengthCheck) ID() string {
	return "COD-002"
}

func (c *ClassLengthCheck) Name() string {
	return "Class Length Check"
}

func (c *ClassLengthCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	classPattern := regexp.MustCompile(`(?i)class\s+\w+`)
	functionPattern := regexp.MustCompile(`(?i)(void|String|int|bool|List|Map)\s+\w+\s*\(`)

	for _, file := range project.DartFiles {
		classCount := len(classPattern.FindAllString(file, -1))
		if classCount > 10 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"File may have too many classes: "+file,
				file,
				"Consider separating classes into different files",
				report.SeverityWarning,
				0,
			))
		}

		funcCount := len(functionPattern.FindAllString(file, -1))
		if funcCount > 20 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"File may have too many functions: "+file,
				file,
				"Consider grouping related functions or extracting to separate classes",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type MethodComplexityCheck struct{}

func (c *MethodComplexityCheck) ID() string {
	return "COD-003"
}

func (c *MethodComplexityCheck) Name() string {
	return "Method Complexity Check"
}

func (c *MethodComplexityCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	for _, file := range project.DartFiles {
		nestingLevel := 0
		for _, char := range file {
			if char == '{' {
				nestingLevel++
				if nestingLevel > 5 {
					findings = append(findings, project.AddFinding(
						c.ID(),
						c.Name(),
						"Deep nesting detected in file: "+file,
						file,
						"Consider extracting nested code into separate methods",
						report.SeverityWarning,
						0,
					))
					break
				}
			} else if char == '}' {
				nestingLevel--
			}
		}
	}

	return findings
}

type NamingConventionCheck struct{}

func (c *NamingConventionCheck) ID() string {
	return "COD-004"
}

func (c *NamingConventionCheck) Name() string {
	return "Naming Convention Check"
}

func (c *NamingConventionCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	classPattern := regexp.MustCompile(`(?i)class\s+[a-z]`)
	variablePattern := regexp.MustCompile(`(?i)\b(?:final|const|var|int|String|bool)\s+[A-Z]`)

	for _, file := range project.DartFiles {
		if classPattern.MatchString(file) {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Class name should start with capital letter in: "+file,
				file,
				"Follow Dart naming conventions (PascalCase for classes)",
				report.SeverityWarning,
				0,
			))
		}

		if variablePattern.MatchString(file) {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Variable name should start with lowercase in: "+file,
				file,
				"Follow Dart naming conventions (camelCase for variables)",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type ImportOrganizationCheck struct{}

func (c *ImportOrganizationCheck) ID() string {
	return "COD-005"
}

func (c *ImportOrganizationCheck) Name() string {
	return "Import Organization Check"
}

func (c *ImportOrganizationCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	importPattern := regexp.MustCompile(`^import\s+['"][^'"]+['"];`)

	for _, file := range project.DartFiles {
		imports := importPattern.FindAllString(file, -1)
		if len(imports) > 15 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Too many imports in file: "+file,
				file,
				"Consider using package: imports and consolidating related imports",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type CommentQualityCheck struct{}

func (c *CommentQualityCheck) ID() string {
	return "COD-006"
}

func (c *CommentQualityCheck) Name() string {
	return "Comment Quality Check"
}

func (c *CommentQualityCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	todoPattern := regexp.MustCompile(`(?i)//\s*TODO`)
	fixmePattern := regexp.MustCompile(`(?i)//\s*FIXME`)
	commentPattern := regexp.MustCompile(`//`)

	for _, file := range project.DartFiles {
		todos := len(todoPattern.FindAllString(file, -1))
		fixmes := len(fixmePattern.FindAllString(file, -1))
		comments := len(commentPattern.FindAllString(file, -1))

		if todos > 5 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Many TODO comments found in: "+file,
				file,
				"Address TODO items or create issues for tracking",
				report.SeverityWarning,
				0,
			))
		}

		if fixmes > 3 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Many FIXME comments found in: "+file,
				file,
				"Fixme items indicate technical debt that should be addressed",
				report.SeverityWarning,
				0,
			))
		}

		lines := strings.Count(file, "\n")
		if lines > 50 && comments == 0 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"No comments found in large file: "+file,
				file,
				"Consider adding documentation comments for public APIs",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type CyclomaticComplexityCheck struct{}

func (c *CyclomaticComplexityCheck) ID() string {
	return "COD-007"
}

func (c *CyclomaticComplexityCheck) Name() string {
	return "Cyclomatic Complexity Check"
}

func (c *CyclomaticComplexityCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	ifPattern := regexp.MustCompile(`(?i)\bif\s*\(`)
	forPattern := regexp.MustCompile(`(?i)\bfor\s*\(`)
	whilePattern := regexp.MustCompile(`(?i)\bwhile\s*\(`)
	casePattern := regexp.MustCompile(`(?i)\bcase\s+`)
	ternaryPattern := regexp.MustCompile(`\?[^:]+:`)

	for _, file := range project.DartFiles {
		ifCount := len(ifPattern.FindAllString(file, -1))
		forCount := len(forPattern.FindAllString(file, -1))
		whileCount := len(whilePattern.FindAllString(file, -1))
		caseCount := len(casePattern.FindAllString(file, -1))
		ternaryCount := len(ternaryPattern.FindAllString(file, -1))

		complexity := ifCount + forCount + whileCount + caseCount + ternaryCount
		if complexity > 15 {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"High cyclomatic complexity detected in: "+file,
				file,
				"Consider extracting complex logic into separate functions",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type DuplicateCodeCheck struct{}

func (c *DuplicateCodeCheck) ID() string {
	return "COD-008"
}

func (c *DuplicateCodeCheck) Name() string {
	return "Duplicate Code Detection"
}

func (c *DuplicateCodeCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	return findings
}
