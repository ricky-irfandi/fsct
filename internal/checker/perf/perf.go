package perf

import (
	"regexp"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type ConstConstructorCheck struct{}

func (c *ConstConstructorCheck) ID() string {
	return "PERF-001"
}

func (c *ConstConstructorCheck) Name() string {
	return "Const Constructor Usage"
}

func (c *ConstConstructorCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	constructorPattern := regexp.MustCompile(`(?i)class\s+\w+\s*\{[^}]*const\s+\w+\(`)

	for _, file := range project.DartFiles {
		if constructorPattern.MatchString(file) {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Consider using const constructors for immutable widgets",
				file,
				"Add 'const' keyword to constructors for better performance",
				report.SeverityInfo,
				0,
			))
		}
	}

	return findings
}

type BuildOptimizationCheck struct{}

func (c *BuildOptimizationCheck) ID() string {
	return "PERF-002"
}

func (c *BuildOptimizationCheck) Name() string {
	return "Build Method Optimization"
}

func (c *BuildOptimizationCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	heavyPattern := regexp.MustCompile(`(?i)(JSON\.decode|HttpClient|File\.read|Database\.query)`)

	for _, file := range project.DartFiles {
		if heavyPattern.MatchString(file) && !strings.Contains(file, "async") {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Heavy operations detected in file",
				file,
				"Move heavy operations out of build method, use async/await",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type ListBuilderCheck struct{}

func (c *ListBuilderCheck) ID() string {
	return "PERF-003"
}

func (c *ListBuilderCheck) Name() string {
	return "List Builder Usage"
}

func (c *ListBuilderCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	mapPattern := regexp.MustCompile(`\.map\s*\(\s*\w+\s*=>\s*[A-Z]`)
	forPattern := regexp.MustCompile(`(?i)for\s*\(`)

	for _, file := range project.DartFiles {
		childrenPattern := regexp.MustCompile(`children:\s*\[`)
		if childrenPattern.MatchString(file) {
			if mapPattern.MatchString(file) || forPattern.MatchString(file) {
				findings = append(findings, project.AddFinding(
					c.ID(),
					c.Name(),
					"Consider using ListView.builder for dynamic lists",
					file,
					"Use ListView.builder instead of .map/.for for better performance",
					report.SeverityInfo,
					0,
				))
			}
		}
	}

	return findings
}

type ImageOptimizationCheck struct{}

func (c *ImageOptimizationCheck) ID() string {
	return "PERF-004"
}

func (c *ImageOptimizationCheck) Name() string {
	return "Image Caching"
}

func (c *ImageOptimizationCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	imagePattern := regexp.MustCompile(`(?i)Image\.(asset|network|file)`)
	cachePattern := regexp.MustCompile(`(?i)precacheImage|CacheManager`)

	for _, file := range project.DartFiles {
		if imagePattern.MatchString(file) && !cachePattern.MatchString(file) {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Image loading without explicit caching",
				file,
				"Consider using cached_network_image or precaching",
				report.SeverityInfo,
				0,
			))
		}
	}

	return findings
}

type StateManagementCheck struct{}

func (c *StateManagementCheck) ID() string {
	return "PERF-005"
}

func (c *StateManagementCheck) Name() string {
	return "State Management Optimization"
}

func (c *StateManagementCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	setStatePattern := regexp.MustCompile(`setState\(\s*\(\)\s*=>\s*\{`)
	providerPattern := regexp.MustCompile(`(?i)Provider|Consumer|Selector`)
	riverpodPattern := regexp.MustCompile(`(?i)Riverpod|useProvider|StateProvider`)
	blocPattern := regexp.MustCompile(`(?i)Bloc|BlocProvider|BlocBuilder`)

	hasStateManagement := providerPattern.MatchString(strings.Join(project.DartFiles, " ")) ||
		riverpodPattern.MatchString(strings.Join(project.DartFiles, " ")) ||
		blocPattern.MatchString(strings.Join(project.DartFiles, " "))

	for _, file := range project.DartFiles {
		setStateCount := len(setStatePattern.FindAllString(file, -1))
		if setStateCount > 15 && !hasStateManagement {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Excessive setState usage detected",
				file,
				"Consider using Provider, Riverpod, or BLoc for complex state",
				report.SeverityWarning,
				0,
			))
		}
	}

	return findings
}

type DependencyOptimizationCheck struct{}

func (c *DependencyOptimizationCheck) ID() string {
	return "PERF-006"
}

func (c *DependencyOptimizationCheck) Name() string {
	return "Dependency Optimization"
}

func (c *DependencyOptimizationCheck) Run(project *checker.Project) []report.Finding {
	findings := []report.Finding{}

	heavyPatterns := []string{
		`(?i)firebase_?auth`,
		`(?i)cloud_?firestore`,
		`(?i)dio|http`,
		`(?i)shared_?preferences`,
	}

	foundHeavy := []string{}
	for _, pattern := range heavyPatterns {
		re := regexp.MustCompile(pattern)
		for _, file := range project.DartFiles {
			if re.MatchString(file) {
				foundHeavy = append(foundHeavy, pattern)
				break
			}
		}
	}

	if len(foundHeavy) == 0 && len(project.DartFiles) > 5 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No heavy dependencies found",
			"pubspec.yaml",
			"Consider if this is intentional for a lightweight app",
			report.SeverityInfo,
			0,
		))
	}

	return findings
}
