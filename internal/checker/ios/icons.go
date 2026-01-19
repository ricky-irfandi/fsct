package ios

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type MissingAppIconCheck struct{}

func (c *MissingAppIconCheck) ID() string {
	return "IOS-008"
}

func (c *MissingAppIconCheck) Name() string {
	return "Missing App Icon Check"
}

func (c *MissingAppIconCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.IOSPath == "" {
		return findings
	}

	appIconPath := filepath.Join(project.IOSPath, "Runner", "Assets.xcassets", "AppIcon.appiconset")

	if _, err := os.Stat(appIconPath); os.IsNotExist(err) {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"AppIcon.appiconset folder not found in Assets.xcassets",
			"ios/Runner/Assets.xcassets/AppIcon.appiconset",
			"Add AppIcon.appiconset with required icon images",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type Missing1024IconCheck struct{}

func (c *Missing1024IconCheck) ID() string {
	return "IOS-009"
}

func (c *Missing1024IconCheck) Name() string {
	return "Missing 1024x1024 Icon Check"
}

func (c *Missing1024IconCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.IOSPath == "" {
		return findings
	}

	appIconPath := filepath.Join(project.IOSPath, "Runner", "Assets.xcassets", "AppIcon.appiconset")
	contentsPath := filepath.Join(appIconPath, "Contents.json")

	content, err := os.ReadFile(contentsPath)
	if err != nil {
		return findings
	}

	var contents map[string]interface{}
	if err := json.Unmarshal(content, &contents); err != nil {
		return findings
	}

	images, ok := contents["images"].([]interface{})
	if !ok {
		return findings
	}

	has1024 := false
	for _, img := range images {
		image, ok := img.(map[string]interface{})
		if !ok {
			continue
		}

		size, ok := image["size"].(string)
		if !ok {
			continue
		}

		idiom, ok := image["idiom"].(string)
		if !ok {
			continue
		}

		if size == "1024x1024" && idiom == "ios-marketing" {
			has1024 = true
			break
		}
	}

	if !has1024 {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No 1024x1024 (ios-marketing) icon found in AppIcon.appiconset Contents.json",
			"ios/Runner/Assets.xcassets/AppIcon.appiconset/Contents.json",
			"Add a 1024x1024 icon for iOS marketing size",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}
