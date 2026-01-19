package ios

import (
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

type CameraUsageDescriptionCheck struct{}

func (c *CameraUsageDescriptionCheck) ID() string {
	return "IOS-001"
}

func (c *CameraUsageDescriptionCheck) Name() string {
	return "Camera Usage Description Check"
}

func (c *CameraUsageDescriptionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if !project.HasCameraDeps {
		return findings
	}

	if project.InfoPlist == nil || !project.InfoPlist.HasCameraUsageDescription {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses camera dependencies but does not have NSCameraUsageDescription in Info.plist",
			"ios/Runner/Info.plist",
			"Add NSCameraUsageDescription with a clear explanation of why camera access is needed",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type PhotoLibraryUsageDescriptionCheck struct{}

func (c *PhotoLibraryUsageDescriptionCheck) ID() string {
	return "IOS-002"
}

func (c *PhotoLibraryUsageDescriptionCheck) Name() string {
	return "Photo Library Usage Description Check"
}

func (c *PhotoLibraryUsageDescriptionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if !project.HasImagePicker {
		return findings
	}

	if project.InfoPlist == nil || !project.InfoPlist.HasPhotoLibraryUsageDescription {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses image_picker but does not have NSPhotoLibraryUsageDescription in Info.plist",
			"ios/Runner/Info.plist",
			"Add NSPhotoLibraryUsageDescription with a clear explanation",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type LocationUsageDescriptionCheck struct{}

func (c *LocationUsageDescriptionCheck) ID() string {
	return "IOS-003"
}

func (c *LocationUsageDescriptionCheck) Name() string {
	return "Location Usage Description Check"
}

func (c *LocationUsageDescriptionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	if !project.HasLocationDeps {
		return findings
	}

	if project.InfoPlist == nil || !project.InfoPlist.HasLocationUsageDescription {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses location dependencies but does not have NSLocationWhenInUseUsageDescription in Info.plist",
			"ios/Runner/Info.plist",
			"Add NSLocationWhenInUseUsageDescription with a clear explanation",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type MicrophoneUsageDescriptionCheck struct{}

func (c *MicrophoneUsageDescriptionCheck) ID() string {
	return "IOS-004"
}

func (c *MicrophoneUsageDescriptionCheck) Name() string {
	return "Microphone Usage Description Check"
}

func (c *MicrophoneUsageDescriptionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	hasMicDep := false
	if project.Pubspec != nil {
		for dep := range project.Pubspec.Dependencies {
			lowerDep := strings.ToLower(dep)
			if strings.Contains(lowerDep, "mic") || strings.Contains(lowerDep, "audio") || strings.Contains(lowerDep, "record") || strings.Contains(lowerDep, "sound") || strings.Contains(lowerDep, "voice") {
				hasMicDep = true
				break
			}
		}
	}

	if !hasMicDep {
		return findings
	}

	if project.InfoPlist == nil || !project.InfoPlist.HasMicrophoneUsageDescription {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses microphone dependencies but does not have NSMicrophoneUsageDescription in Info.plist",
			"ios/Runner/Info.plist",
			"Add NSMicrophoneUsageDescription with a clear explanation",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type ContactsUsageDescriptionCheck struct{}

func (c *ContactsUsageDescriptionCheck) ID() string {
	return "IOS-005"
}

func (c *ContactsUsageDescriptionCheck) Name() string {
	return "Contacts Usage Description Check"
}

func (c *ContactsUsageDescriptionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	hasContactsDep := false
	if project.Pubspec != nil {
		for dep := range project.Pubspec.Dependencies {
			if strings.Contains(strings.ToLower(dep), "contact") {
				hasContactsDep = true
				break
			}
		}
	}

	if !hasContactsDep {
		return findings
	}

	if project.InfoPlist == nil || !project.InfoPlist.HasContactsUsageDescription {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses contacts dependencies but does not have NSContactsUsageDescription in Info.plist",
			"ios/Runner/Info.plist",
			"Add NSContactsUsageDescription with a clear explanation",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type CalendarsUsageDescriptionCheck struct{}

func (c *CalendarsUsageDescriptionCheck) ID() string {
	return "IOS-006"
}

func (c *CalendarsUsageDescriptionCheck) Name() string {
	return "Calendars Usage Description Check"
}

func (c *CalendarsUsageDescriptionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	hasCalendarDep := false
	if project.Pubspec != nil {
		for dep := range project.Pubspec.Dependencies {
			if strings.Contains(strings.ToLower(dep), "calendar") || strings.Contains(strings.ToLower(dep), "event") {
				hasCalendarDep = true
				break
			}
		}
	}

	if !hasCalendarDep {
		return findings
	}

	if project.InfoPlist == nil || !project.InfoPlist.HasCalendarsUsageDescription {
		findings = append(findings, project.AddFinding(
			c.ID(),
			c.Name(),
			"App uses calendar dependencies but does not have NSCalendarsUsageDescription in Info.plist",
			"ios/Runner/Info.plist",
			"Add NSCalendarsUsageDescription with a clear explanation",
			report.SeverityHigh,
			0,
		))
	}

	return findings
}

type EmptyUsageDescriptionCheck struct{}

func (c *EmptyUsageDescriptionCheck) ID() string {
	return "IOS-007"
}

func (c *EmptyUsageDescriptionCheck) Name() string {
	return "Empty Usage Description Check"
}

func (c *EmptyUsageDescriptionCheck) Run(project *checker.Project) []report.Finding {
	var findings []report.Finding

	return findings
}
