package checker

import (
	"strings"

	"github.com/ricky-irfandi/fsct/internal/report"
)

type Check interface {
	ID() string
	Name() string
	Run(project *Project) []report.Finding
}

type Category string

const (
	CategoryAndroid  Category = "Android"
	CategoryIOS      Category = "iOS"
	CategoryFlutter  Category = "Flutter"
	CategorySecurity Category = "Security"
	CategoryReviewer Category = "Reviewer"
	CategoryAI       Category = "AI"
)

type Project struct {
	Path        string
	AndroidPath string
	IOSPath     string
	FlutterPath string

	AndroidManifest *AndroidManifestInfo
	GradleConfig    *GradleConfigInfo
	InfoPlist       *InfoPlistInfo
	Pubspec         *PubspecInfo
	DartFiles       []string

	HasNetworkDeps   bool
	HasCameraDeps    bool
	HasLocationDeps  bool
	HasImagePicker   bool
	HasURLLauncher   bool
	HasLoginPatterns bool
}

type AndroidManifestInfo struct {
	PackageName     string
	VersionCode     string
	VersionName     string
	Debuggable      bool
	AllowBackup     bool
	Permissions     []string
	Activities      []ActivityInfo
	QueriesPackages []string
}

type ActivityInfo struct {
	Name            string
	Exported        bool
	HasIntentFilter bool
}

type GradleConfigInfo struct {
	ApplicationID    string
	MinSDKVersion    string
	TargetSDKVersion string
	VersionCode      string
	VersionName      string
}

type InfoPlistInfo struct {
	CFBundleIdentifier         string
	CFBundleVersion            string
	CFBundleShortVersionString string

	HasCameraUsageDescription       bool
	HasPhotoLibraryUsageDescription bool
	HasLocationUsageDescription     bool
	HasMicrophoneUsageDescription   bool
	HasContactsUsageDescription     bool
	HasCalendarsUsageDescription    bool

	EncryptionDeclarationSet bool
	EncryptionExempt         bool
	RequiresFullScreen       bool
}

type PubspecInfo struct {
	Name        string
	Version     string
	Description string
	Homepage    string
	Repository  string

	Dependencies    map[string]string
	DevDependencies map[string]string

	HasLinter        bool
	HasIconConfig    bool
	HasSplashConfig  bool
	HasDeprecatedPkg bool
	HasDebugDeps     bool
}

func NewProject(path string) *Project {
	return &Project{
		Path:        path,
		AndroidPath: path + "/android",
		IOSPath:     path + "/ios",
		FlutterPath: path,

		AndroidManifest: &AndroidManifestInfo{},
		GradleConfig:    &GradleConfigInfo{},
		InfoPlist:       &InfoPlistInfo{},
		Pubspec:         &PubspecInfo{},

		DartFiles: make([]string, 0),
	}
}

func (p *Project) AddFinding(id, title, message, file, suggestion string, severity report.Severity, line int) report.Finding {
	return report.Finding{
		ID:         id,
		Severity:   severity,
		Title:      title,
		Message:    message,
		File:       file,
		Line:       line,
		Suggestion: suggestion,
	}
}

type CheckRegistry struct {
	checks map[string]Check
}

func NewCheckRegistry() *CheckRegistry {
	return &CheckRegistry{
		checks: make(map[string]Check),
	}
}

func (r *CheckRegistry) Register(check Check) {
	r.checks[check.ID()] = check
}

func (r *CheckRegistry) Get(id string) (Check, bool) {
	check, ok := r.checks[id]
	return check, ok
}

func (r *CheckRegistry) GetAll() []Check {
	checks := make([]Check, 0, len(r.checks))
	for _, check := range r.checks {
		checks = append(checks, check)
	}
	return checks
}

func (r *CheckRegistry) GetByCategory(prefix string) []Check {
	checks := make([]Check, 0)
	for _, check := range r.checks {
		if strings.HasPrefix(check.ID(), prefix) {
			checks = append(checks, check)
		}
	}
	return checks
}
