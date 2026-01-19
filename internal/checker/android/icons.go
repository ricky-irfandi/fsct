package android

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type MissingAppIconCheck struct{}

func (c *MissingAppIconCheck) ID() string {
	return "AND-007"
}

func (c *MissingAppIconCheck) Name() string {
	return "Missing App Icon Check"
}

func (c *MissingAppIconCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.AndroidPath == "" {
		return findings
	}

	resPath := filepath.Join(project.AndroidPath, "app", "src", "main", "res")

	mipmapDirs := []string{
		"mipmap-hdpi",
		"mipmap-mdpi",
		"mipmap-xhdpi",
		"mipmap-xxhdpi",
		"mipmap-xxxhdpi",
	}

	foundIcons := false
	for _, dir := range mipmapDirs {
		mipmapPath := filepath.Join(resPath, dir)
		files, err := os.ReadDir(mipmapPath)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !f.IsDir() && (strings.HasSuffix(f.Name(), ".png") || strings.HasSuffix(f.Name(), ".webp")) {
				if strings.HasPrefix(f.Name(), "ic_launcher") {
					foundIcons = true
					break
				}
			}
		}
		if foundIcons {
			break
		}
	}

	if !foundIcons {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"No ic_launcher icon found in res/mipmap-* directories",
			"android/app/src/main/res/",
			"Add ic_launcher.png or ic_launcher.webp to all mipmap directories",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type PlaceholderIconCheck struct{}

func (c *PlaceholderIconCheck) ID() string {
	return "AND-008"
}

func (c *PlaceholderIconCheck) Name() string {
	return "Placeholder Icon Check"
}

func (c *PlaceholderIconCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if project.AndroidPath == "" {
		return findings
	}

	resPath := filepath.Join(project.AndroidPath, "app", "src", "main", "res")

	mipmapPath := filepath.Join(resPath, "mipmap-hdpi")
	files, err := os.ReadDir(mipmapPath)
	if err != nil {
		return findings
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		content, err := os.ReadFile(filepath.Join(mipmapPath, f.Name()))
		if err != nil {
			continue
		}

		if isDefaultFlutterIcon(content) {
			findings = append(findings, project.AddFinding(
				c.ID(),
				c.Name(),
				"Default Flutter launcher icon detected",
				"android/app/src/main/res/mipmap-hdpi/"+f.Name(),
				"Replace with your app's custom launcher icon",
				report.SeverityWarning,
				0,
			))
			break
		}
	}

	return findings
}

func isDefaultFlutterIcon(data []byte) bool {
	if len(data) < 100 {
		return false
	}

	defaultIconSignatures := [][]byte{
		{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
	}

	for _, sig := range defaultIconSignatures {
		if len(data) >= len(sig) {
			match := true
			for i, b := range sig {
				if data[i] != b {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}

	return false
}
