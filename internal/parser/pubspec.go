package parser

import (
	"os"
	"regexp"
	"strings"
)

type Pubspec struct {
	Name            string
	Version         string
	Description     string
	Homepage        string
	Repository      string
	Dependencies    map[string]string
	DevDependencies map[string]string
	Flutter         *FlutterConfig
	Environment     map[string]string
}

type FlutterConfig struct {
	UsesMaterialDesign bool
	Plugin             *PluginConfig
	Assets             []string
	Fonts              []FontConfig
}

type PluginConfig struct {
	Platforms map[string]PluginPlatformConfig
}

type PluginPlatformConfig struct {
	PluginClass string
}

type FontConfig struct {
	Family string
	Fonts  []FontEntry
}

type FontEntry struct {
	Weight int
	Style  string
}

func isTopLevel(line string) bool {
	return len(line) > 0 && (line[0] == ' ' || line[0] == '\t') == false
}

func ParsePubspec(path string) (*Pubspec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	pubspec := &Pubspec{
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
		Environment:     make(map[string]string),
	}

	namePattern := regexp.MustCompile(`^name:\s*(.+)$`)
	versionPattern := regexp.MustCompile(`^version:\s*(.+)$`)
	descriptionPattern := regexp.MustCompile(`^description:\s*(.+)$`)
	homepagePattern := regexp.MustCompile(`^homepage:\s*(.+)$`)
	repositoryPattern := regexp.MustCompile(`^repository:\s*(.+)$`)
	usesMaterialDesignPattern := regexp.MustCompile(`uses-material-design:\s*(true|false)`)

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)

		if matches := namePattern.FindStringSubmatch(line); len(matches) > 1 {
			pubspec.Name = strings.TrimSpace(matches[1])
		}
		if matches := versionPattern.FindStringSubmatch(line); len(matches) > 1 {
			pubspec.Version = strings.TrimSpace(matches[1])
		}
		if matches := descriptionPattern.FindStringSubmatch(line); len(matches) > 1 {
			pubspec.Description = strings.TrimSpace(matches[1])
		}
		if matches := homepagePattern.FindStringSubmatch(line); len(matches) > 1 {
			pubspec.Homepage = strings.TrimSpace(matches[1])
		}
		if matches := repositoryPattern.FindStringSubmatch(line); len(matches) > 1 {
			pubspec.Repository = strings.TrimSpace(matches[1])
		}
		if matches := usesMaterialDesignPattern.FindStringSubmatch(line); len(matches) > 1 {
			pubspec.Flutter = &FlutterConfig{}
			pubspec.Flutter.UsesMaterialDesign = matches[1] == "true"
		}
	}

	depPattern := regexp.MustCompile(`^([a-zA-Z0-9_-]+):\s*(.+)$`)
	simpleDepPattern := regexp.MustCompile(`^-\s*([a-zA-Z0-9_-]+):\s*(.+)$`)

	currentSection := ""
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		sectionPatterns := map[string]string{
			"dependencies:":     "dependencies",
			"dev_dependencies:": "dev_dependencies",
			"flutter:":          "flutter",
			"environment:":      "environment",
		}

		if isTopLevel(line) {
			for pattern, section := range sectionPatterns {
				if trimmed == pattern {
					currentSection = section
					break
				}
			}
		}

		if currentSection == "dependencies" || currentSection == "dev_dependencies" {
			if matches := depPattern.FindStringSubmatch(trimmed); len(matches) > 1 {
				depName := matches[1]
				depVersion := strings.TrimSpace(matches[2])
				if currentSection == "dependencies" {
					pubspec.Dependencies[depName] = depVersion
				} else {
					pubspec.DevDependencies[depName] = depVersion
				}
			} else if matches := simpleDepPattern.FindStringSubmatch(trimmed); len(matches) > 1 {
				depName := matches[1]
				depVersion := strings.TrimSpace(matches[2])
				if currentSection == "dependencies" {
					pubspec.Dependencies[depName] = depVersion
				} else {
					pubspec.DevDependencies[depName] = depVersion
				}
			}
		}

		if currentSection == "environment" {
			envPattern := regexp.MustCompile(`^([a-zA-Z0-9_]+):\s*(.+)$`)
			if matches := envPattern.FindStringSubmatch(trimmed); len(matches) > 1 {
				pubspec.Environment[matches[1]] = strings.TrimSpace(matches[2])
			}
		}
	}

	if pubspec.Flutter == nil {
		usesMaterialDesignPattern := regexp.MustCompile(`uses-material-design:\s*(true|false)`)
		for _, line := range strings.Split(content, "\n") {
			if strings.Contains(line, "flutter:") {
				pubspec.Flutter = &FlutterConfig{UsesMaterialDesign: false}
				break
			}
		}
		for _, line := range strings.Split(content, "\n") {
			if matches := usesMaterialDesignPattern.FindStringSubmatch(line); len(matches) > 1 {
				if pubspec.Flutter == nil {
					pubspec.Flutter = &FlutterConfig{}
				}
				pubspec.Flutter.UsesMaterialDesign = matches[1] == "true"
				break
			}
		}
	}

	return pubspec, nil
}

func (p *Pubspec) HasDependency(dep string) bool {
	_, ok := p.Dependencies[dep]
	return ok
}

func (p *Pubspec) HasDevDependency(dep string) bool {
	_, ok := p.DevDependencies[dep]
	return ok
}

func (p *Pubspec) HasLinter() bool {
	linters := []string{"flutter_lints", "very_good_analysis", "pedantic", "dart_style"}
	for _, linter := range linters {
		if p.HasDevDependency(linter) {
			return true
		}
	}
	return false
}

func (p *Pubspec) HasIconConfig() bool {
	if p.HasDependency("flutter_launcher_icons") {
		return true
	}
	if _, ok := p.DevDependencies["flutter_launcher_icons"]; ok {
		return true
	}
	return false
}

func (p *Pubspec) HasSplashConfig() bool {
	return p.HasDependency("flutter_native_splash")
}

func (p *Pubspec) HasDeprecatedPackage() bool {
	deprecatedPackages := []string{"package_info", "android_alarm_manager", "device_info"}
	for _, pkg := range deprecatedPackages {
		if p.HasDependency(pkg) || p.HasDevDependency(pkg) {
			return true
		}
	}
	return false
}

func (p *Pubspec) HasDebugDepInMain() bool {
	debugPackages := []string{"flutter_driver", "integration_test", "mockito"}
	for _, pkg := range debugPackages {
		if p.HasDependency(pkg) {
			return true
		}
	}
	return false
}

func (p *Pubspec) IsDefaultVersion() bool {
	return p.Version == "1.0.0+1" || p.Version == "1.0.0"
}

func (p *Pubspec) HasDescription() bool {
	return p.Description != ""
}

func (p *Pubspec) HasRepository() bool {
	return p.Repository != "" || p.Homepage != ""
}

func (p *Pubspec) GetSdkVersion() string {
	if val, ok := p.Environment["sdk"]; ok {
		return val
	}
	return ""
}
