package prompt

import (
	"github.com/ricky-irfandi/fsct/internal/report"
)

// PromptData contains all data needed to generate a comprehensive AI prompt
type PromptData struct {
	// App Metadata
	AppName        string
	Version        string
	Description    string
	FlutterVersion string
	Repository     string
	Homepage       string

	// Compliance Score
	ComplianceScore int
	Status          string

	// Static Analysis Results
	FindingsBySeverity map[report.Severity][]FindingSummary
	FindingsByCategory map[string][]FindingSummary
	TotalChecks        int
	PassedChecks       int

	// Platform Configuration
	AndroidConfig AndroidSummary
	IOSConfig     IOSSummary
	FlutterConfig FlutterSummary

	// Permissions
	AndroidPermissions []PermissionInfo
	IOSPermissions     []PermissionInfo

	// Dependencies
	Dependencies    []DependencyInfo
	DevDependencies []DependencyInfo
	HighRiskDeps    []DependencyInfo
	OutdatedDeps    []DependencyInfo

	// Security & Policy
	SecurityFlags   []SecurityFlag
	PolicyFlags     []PolicyFlag
	MissingPolicies []string

	// Reviewer Account
	ReviewerAccount *ReviewerAccountInfo

	// Blockers
	Blockers []BlockerInfo
}

// FindingSummary condensed finding for prompt
type FindingSummary struct {
	ID         string
	Severity   report.Severity
	Category   string
	Title      string
	Message    string
	File       string
	Suggestion string
}

// AndroidSummary Android-specific configuration
type AndroidSummary struct {
	PackageName      string
	TargetSDKVersion int
	MinSDKVersion    int
	VersionCode      int
	VersionName      string
	IsDebuggable     bool
	AllowBackup      bool
	HasInternetPerm  bool
}

// IOSSummary iOS-specific configuration
type IOSSummary struct {
	BundleIdentifier         string
	BundleVersion            string
	BundleShortVersion       string
	DeploymentTarget         string
	RequiresFullScreen       bool
	EncryptionDeclared       bool
	HasCameraDescription     bool
	HasLocationDescription   bool
	HasPhotoLibDescription   bool
	HasMicrophoneDescription bool
}

// FlutterSummary Flutter-specific configuration
type FlutterSummary struct {
	UsesMaterial3    bool
	HasLinter        bool
	HasIconConfig    bool
	HasSplashConfig  bool
	HasDeprecatedPkg bool
	HasDebugDeps     bool
}

// PermissionInfo permission details
type PermissionInfo struct {
	Name          string
	Platform      string // android, ios
	Present       bool
	HasDescription bool
	Description   string
	RiskLevel     string // high, medium, low
}

// DependencyInfo dependency details
type DependencyInfo struct {
	Name       string
	Version    string
	IsDev      bool
	RiskLevel  string // high, medium, low
	IsOutdated bool
	LatestVer  string
}

// SecurityFlag security-related flag
type SecurityFlag struct {
	Type        string
	Present     bool
	Description string
	Severity    report.Severity
}

// PolicyFlag policy-related flag
type PolicyFlag struct {
	Type        string
	Present     bool
	Description string
}

// ReviewerAccountInfo reviewer test account configuration
type ReviewerAccountInfo struct {
	Configured     bool
	Email          string
	PasswordEnv    string
	HasWeakPassword bool
	CanVerify      bool
	LoginEndpoint  string
}

// BlockerInfo critical blocker for submission
type BlockerInfo struct {
	ID          string
	Category    string
	Title       string
	Description string
	Platform    string // android, ios, both
	FixFile     string
	FixSnippet  string
}

// NewPromptData creates a new PromptData instance
func NewPromptData() *PromptData {
	return &PromptData{
		FindingsBySeverity: make(map[report.Severity][]FindingSummary),
		FindingsByCategory: make(map[string][]FindingSummary),
		AndroidPermissions: make([]PermissionInfo, 0),
		IOSPermissions:     make([]PermissionInfo, 0),
		Dependencies:       make([]DependencyInfo, 0),
		DevDependencies:    make([]DependencyInfo, 0),
		HighRiskDeps:       make([]DependencyInfo, 0),
		OutdatedDeps:       make([]DependencyInfo, 0),
		SecurityFlags:      make([]SecurityFlag, 0),
		PolicyFlags:        make([]PolicyFlag, 0),
		MissingPolicies:    make([]string, 0),
		Blockers:           make([]BlockerInfo, 0),
	}
}

// AddFinding adds a finding to the appropriate buckets
func (p *PromptData) AddFinding(f FindingSummary) {
	// By severity
	p.FindingsBySeverity[f.Severity] = append(p.FindingsBySeverity[f.Severity], f)
	// By category
	p.FindingsByCategory[f.Category] = append(p.FindingsByCategory[f.Category], f)
}

// GetHighCount returns number of HIGH severity findings
func (p *PromptData) GetHighCount() int {
	return len(p.FindingsBySeverity[report.SeverityHigh])
}

// GetWarningCount returns number of WARNING severity findings
func (p *PromptData) GetWarningCount() int {
	return len(p.FindingsBySeverity[report.SeverityWarning])
}

// GetInfoCount returns number of INFO severity findings
func (p *PromptData) GetInfoCount() int {
	return len(p.FindingsBySeverity[report.SeverityInfo])
}

// IsReadyForAppStore checks if app is ready for App Store submission
func (p *PromptData) IsReadyForAppStore() bool {
	// Check for iOS-specific blockers
	for _, blocker := range p.Blockers {
		if blocker.Platform == "ios" || blocker.Platform == "both" {
			return false
		}
	}
	// Check for any HIGH severity iOS findings
	for _, f := range p.FindingsBySeverity[report.SeverityHigh] {
		if f.Category == "iOS" || f.Category == "Security" {
			return false
		}
	}
	return true
}

// IsReadyForPlayStore checks if app is ready for Play Store submission
func (p *PromptData) IsReadyForPlayStore() bool {
	// Check for Android-specific blockers
	for _, blocker := range p.Blockers {
		if blocker.Platform == "android" || blocker.Platform == "both" {
			return false
		}
	}
	// Check for any HIGH severity Android findings
	for _, f := range p.FindingsBySeverity[report.SeverityHigh] {
		if f.Category == "Android" || f.Category == "Security" {
			return false
		}
	}
	return true
}

// CalculateComplianceScore computes overall compliance score
func (p *PromptData) CalculateComplianceScore() int {
	if p.TotalChecks == 0 {
		return 0
	}
	// Score based on passed checks
	passed := p.PassedChecks
	// Deduct for high severity issues
	highPenalty := p.GetHighCount() * 10
	warningPenalty := p.GetWarningCount() * 3
	infoPenalty := p.GetInfoCount()

	score := (passed * 100 / p.TotalChecks) - highPenalty - warningPenalty - infoPenalty
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

// GetStatus returns overall compliance status
func (p *PromptData) GetStatus() string {
	score := p.CalculateComplianceScore()
	highCount := p.GetHighCount()

	if highCount > 0 {
		return "❌ NEEDS CRITICAL ATTENTION"
	}
	if score >= 90 {
		return "✅ EXCELLENT"
	}
	if score >= 75 {
		return "⚠️ GOOD WITH MINOR ISSUES"
	}
	if score >= 50 {
		return "⚠️ NEEDS ATTENTION"
	}
	return "❌ SIGNIFICANT WORK REQUIRED"
}
