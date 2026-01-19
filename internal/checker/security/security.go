package security

import (
	"regexp"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

var (
	credentialPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)api_?key['"]?\s*[:=]\s*['"]?[a-zA-Z0-9_-]{20,}['"]?`),
		regexp.MustCompile(`(?i)secret['"]?\s*[:=]\s*['"]?[a-zA-Z0-9_\-+=/]{16,}['"]?`),
		regexp.MustCompile(`(?i)password['"]?\s*[:=]\s*['"]?[^'"]{8,}['"]?`),
		regexp.MustCompile(`(?i)auth_?token['"]?\s*[:=]\s*['"]?[a-zA-Z0-9_\-\.]{20,}['"]?`),
	}
)

type HardcodedCredentialsCheck struct{}

func (c *HardcodedCredentialsCheck) ID() string {
	return "SEC-001"
}

func (c *HardcodedCredentialsCheck) Name() string {
	return "Hardcoded Credentials Detection"
}

func (c *HardcodedCredentialsCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	for _, file := range project.DartFiles {
		for _, re := range credentialPatterns {
			if re.MatchString(file) {
				findings = append(findings, project.AddFinding(
					c.ID(),
					c.Name(),
					"Potential hardcoded credentials found in file",
					file,
					"Use environment variables or secure configuration storage",
					report.SeverityHigh,
					0,
				))
			}
		}
	}

	return findings
}

type DebugModeCheck struct{}

func (c *DebugModeCheck) ID() string {
	return "SEC-002"
}

func (c *DebugModeCheck) Name() string {
	return "Debug Mode Check"
}

func (c *DebugModeCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	for _, file := range project.DartFiles {
		if strings.Contains(file, "print(") || strings.Contains(file, "debugPrint(") {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Debug statements found in file",
				file,
				"Remove debug print statements or use logger that respects build mode",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type InsecureHTTPCheck struct{}

func (c *InsecureHTTPCheck) ID() string {
	return "SEC-003"
}

func (c *InsecureHTTPCheck) Name() string {
	return "Insecure HTTP URL Usage"
}

func (c *InsecureHTTPCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	httpPattern := regexp.MustCompile(`(?i)http://[^\s"'<>]+`)

	for _, file := range project.DartFiles {
		matches := httpPattern.FindAllString(file, -1)
		for _, match := range matches {
			if !strings.Contains(match, "localhost") && !strings.Contains(match, "127.0.0.1") {
				findings = append(findings, project.AddFinding(
					c.ID(),
					c.Name(),
					"Found insecure HTTP URL: "+match,
					file,
					"Replace HTTP URLs with HTTPS for secure communication",
					report.SeverityHigh,
					0,
				))
			}
		}
	}

	return findings
}

type ExportedActivityCheck struct{}

func (c *ExportedActivityCheck) ID() string {
	return "SEC-004"
}

func (c *ExportedActivityCheck) Name() string {
	return "Android Exportable Activity Security"
}

func (c *ExportedActivityCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	if project.AndroidManifest == nil || project.AndroidManifest.Activities == nil {
		return findings
	}

	for _, activity := range project.AndroidManifest.Activities {
		if activity.Exported && !activity.HasIntentFilter {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Exported activity without intent filter: "+activity.Name,
				"AndroidManifest.xml",
				"Set android:exported=\"false\" for activities that don't need external access",
				report.SeverityHigh,
				0,
			))
		}
	}

	return findings
}

type SQLInjectionCheck struct{}

func (c *SQLInjectionCheck) ID() string {
	return "SEC-005"
}

func (c *SQLInjectionCheck) Name() string {
	return "SQL Injection Prevention"
}

func (c *SQLInjectionCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	sqlPattern := regexp.MustCompile(`(?i)rawQuery\s*\(\s*["']\s*(?:SELECT|INSERT|UPDATE|DELETE).*["']\s*\)`)

	for _, file := range project.DartFiles {
		if sqlPattern.MatchString(file) {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Potential SQL injection vulnerability in file",
				file,
				"Use parameterized queries instead of string concatenation",
				report.SeverityHigh,
				0,
			))
		}
	}

	return findings
}
