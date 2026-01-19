package testing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type TestDirectoryCheck struct{}

func (c *TestDirectoryCheck) ID() string {
	return "TST-001"
}

func (c *TestDirectoryCheck) Name() string {
	return "Test Directory Existence"
}

func (c *TestDirectoryCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	hasTestDir := false
	for _, file := range project.DartFiles {
		if strings.Contains(file, "test/") || strings.Contains(file, "_test.dart") {
			hasTestDir = true
			break
		}
	}

	if !hasTestDir {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No test directory found in project",
			"project root",
			"Add a test/ directory with unit and widget tests",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type TestFileNamingCheck struct{}

func (c *TestFileNamingCheck) ID() string {
	return "TST-002"
}

func (c *TestFileNamingCheck) Name() string {
	return "Test File Naming Convention"
}

func (c *TestFileNamingCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	testPattern := regexp.MustCompile(`_test\.dart$`)

	for _, file := range project.DartFiles {
		if strings.Contains(file, "test/") && !testPattern.MatchString(file) {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Test file should end with _test.dart: "+file,
				file,
				"Rename file to follow Flutter test naming conventions",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type TestCoverageCheck struct{}

func (c *TestCoverageCheck) ID() string {
	return "TST-003"
}

func (c *TestCoverageCheck) Name() string {
	return "Test Coverage Check"
}

func (c *TestCoverageCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	testCount := 0
	for _, file := range project.DartFiles {
		if strings.Contains(file, "test/") {
			testCount++
		}
	}

	if testCount == 0 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No test files found",
			"test/",
			"Add unit and widget tests for better code coverage",
			report.SeverityWarning,
			0,
		))
	} else if testCount < 3 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"Low number of test files: "+strconv.Itoa(testCount),
			"test/",
			"Consider adding more tests for better coverage",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type WidgetTestCheck struct{}

func (c *WidgetTestCheck) ID() string {
	return "TST-004"
}

func (c *WidgetTestCheck) Name() string {
	return "Widget Test Presence"
}

func (c *WidgetTestCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	widgetTestPattern := regexp.MustCompile(`(?i)testWidgets`)
	hasWidgetTests := false

	for _, file := range project.DartFiles {
		if strings.Contains(file, "test/") && widgetTestPattern.MatchString(file) {
			hasWidgetTests = true
			break
		}
	}

	if !hasWidgetTests {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No widget tests found",
			"test/",
			"Add widget tests (testWidgets) for UI components",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type MockDependenciesCheck struct{}

func (c *MockDependenciesCheck) ID() string {
	return "TST-005"
}

func (c *MockDependenciesCheck) Name() string {
	return "Mock Dependencies Usage"
}

func (c *MockDependenciesCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	mockPattern := regexp.MustCompile(`(?i)(mockito|Mock|when|verify)`)
	hasMocks := false

	for _, file := range project.DartFiles {
		if strings.Contains(file, "test/") && mockPattern.MatchString(file) {
			hasMocks = true
			break
		}
	}

	if !hasMocks {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No mock dependencies found in tests",
			"test/",
			"Consider using mockito for mocking dependencies in tests",
			report.SeverityWarning,
			0,
		))
	}

	return findings
}

type GoldenTestCheck struct{}

func (c *GoldenTestCheck) ID() string {
	return "TST-006"
}

func (c *GoldenTestCheck) Name() string {
	return "Golden Test Presence"
}

func (c *GoldenTestCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	goldenPattern := regexp.MustCompile(`(?i)matchesGoldenFile`)
	hasGoldenTests := false

	for _, file := range project.DartFiles {
		if strings.Contains(file, "test/") && goldenPattern.MatchString(file) {
			hasGoldenTests = true
			break
		}
	}

	if !hasGoldenTests {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No golden tests found",
			"test/",
			"Consider adding golden tests for UI regression testing",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}
