package ai

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/report"
)

// ComplianceMetadata contains privacy-safe metadata for AI analysis
// IMPORTANT: This struct must NEVER contain:
// - Source code
// - File contents
// - Absolute file paths
// - Sensitive configuration values
// - API keys or credentials
type ComplianceMetadata struct {
	// App identification (safe)
	AppName        string `json:"app_name"`
	Version        string `json:"version"`
	Description    string `json:"description,omitempty"`
	FlutterVersion string `json:"flutter_version,omitempty"`

	// Compliance summary (safe)
	ComplianceScore int `json:"compliance_score"`
	TotalChecks     int `json:"total_checks"`
	PassedChecks    int `json:"passed_checks"`
	HighCount       int `json:"high_count"`
	WarningCount    int `json:"warning_count"`
	InfoCount       int `json:"info_count"`

	// Platform versions (safe)
	AndroidTargetSDK int    `json:"android_target_sdk,omitempty"`
	AndroidMinSDK    int    `json:"android_min_sdk,omitempty"`
	IOSDeploymentTarget string `json:"ios_deployment_target,omitempty"`

	// Finding summaries (safe - minimal info)
	Findings []FindingMeta `json:"findings"`

	// Permissions (safe - just names)
	AndroidPermissions []string `json:"android_permissions,omitempty"`
	IOSPermissions     []string `json:"ios_permissions,omitempty"`

	// Dependencies (safe - names only, no versions for security)
	Dependencies    []string `json:"dependencies,omitempty"`
	DevDependencies []string `json:"dev_dependencies,omitempty"`

	// Feature flags (safe - booleans only)
	Features AppFeatures `json:"features"`

	// Configuration flags (safe)
	Config AppConfig `json:"config"`

	// Security flags (safe)
	Security SecurityFlags `json:"security"`
}

// FindingMeta contains minimal finding information (privacy-safe)
type FindingMeta struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Severity string `json:"severity"`
	Title    string `json:"title"`
	// NO file paths, NO line numbers, NO code snippets
}

// AppFeatures contains boolean feature flags
type AppFeatures struct {
	HasLogin          bool `json:"has_login"`
	HasCamera         bool `json:"has_camera"`
	HasLocation       bool `json:"has_location"`
	HasMicrophone     bool `json:"has_microphone"`
	HasPhotoLibrary   bool `json:"has_photo_library"`
	HasContacts       bool `json:"has_contacts"`
	HasCalendar       bool `json:"has_calendar"`
	HasBluetooth      bool `json:"has_bluetooth"`
	HasNotifications  bool `json:"has_notifications"`
	HasInAppPurchase  bool `json:"has_in_app_purchase"`
	HasSocialLogin    bool `json:"has_social_login"`
	HasEmailLogin     bool `json:"has_email_login"`
	HasDataDeletion   bool `json:"has_data_deletion"`
	HasPrivacyPolicy  bool `json:"has_privacy_policy"`
	HasTermsOfService bool `json:"has_terms_of_service"`
}

// AppConfig contains app configuration flags
type AppConfig struct {
	UsesMaterial3     bool `json:"uses_material_3"`
	HasLinter         bool `json:"has_linter"`
	HasIconConfig     bool `json:"has_icon_config"`
	HasSplashConfig   bool `json:"has_splash_config"`
	HasDeprecatedDeps bool `json:"has_deprecated_deps"`
	HasDebugDeps      bool `json:"has_debug_deps"`
}

// SecurityFlags contains security-related boolean flags
type SecurityFlags struct {
	IsDebuggable        bool `json:"is_debuggable"`
	AllowsBackup        bool `json:"allows_backup"`
	HasInsecureHTTP     bool `json:"has_insecure_http"`
	HasHardcodedKeys    bool `json:"has_hardcoded_keys"`
	MissingEncryption   bool `json:"missing_encryption_decl"`
}

// ExtractMetadata extracts privacy-safe metadata from a project
func ExtractMetadata(project *checker.Project, findings []report.Finding) *ComplianceMetadata {
	meta := &ComplianceMetadata{
		Findings:           make([]FindingMeta, 0, len(findings)),
		AndroidPermissions: make([]string, 0),
		IOSPermissions:     make([]string, 0),
		Dependencies:       make([]string, 0),
		DevDependencies:    make([]string, 0),
	}

	// Extract basic app info
	meta.extractAppInfo(project)

	// Extract findings (sanitized)
	meta.extractFindings(findings)

	// Calculate scores
	meta.calculateScores(findings)

	// Extract platform info
	meta.extractPlatformInfo(project)

	// Extract permissions
	meta.extractPermissions(project)

	// Extract dependencies
	meta.extractDependencies(project)

	// Extract feature flags
	meta.extractFeatures(project)

	// Extract config flags
	meta.extractConfig(project)

	// Extract security flags
	meta.extractSecurityFlags(project, findings)

	return meta
}

// extractAppInfo extracts basic app information
func (m *ComplianceMetadata) extractAppInfo(project *checker.Project) {
	if project.Pubspec != nil {
		m.AppName = project.Pubspec.Name
		m.Version = project.Pubspec.Version
		m.Description = sanitizeDescription(project.Pubspec.Description)
	}

	// Fallback to path name if no app name
	if m.AppName == "" && project.Path != "" {
		m.AppName = filepath.Base(project.Path)
	}

	// Default version if not set
	if m.Version == "" {
		m.Version = "1.0.0"
	}

	// Flutter SDK version not directly available in PubspecInfo
	// Could be extracted from Environment map if needed
}

// extractFindings converts findings to metadata format (privacy-safe)
func (m *ComplianceMetadata) extractFindings(findings []report.Finding) {
	for _, f := range findings {
		meta := FindingMeta{
			ID:       f.ID,
			Category: extractCategory(f.ID),
			Severity: string(f.Severity),
			Title:    f.Title,
		}
		m.Findings = append(m.Findings, meta)

		// Count by severity
		switch f.Severity {
		case report.SeverityHigh:
			m.HighCount++
		case report.SeverityWarning:
			m.WarningCount++
		case report.SeverityInfo:
			m.InfoCount++
		}
	}
}

// calculateScores calculates compliance scores
func (m *ComplianceMetadata) calculateScores(findings []report.Finding) {
	m.TotalChecks = 83 // Total number of checks in the system
	m.PassedChecks = m.TotalChecks - len(findings)

	// Simple scoring: base 100, deduct for issues
	score := 100
	score -= m.HighCount * 10
	score -= m.WarningCount * 3
	score -= m.InfoCount

	if score < 0 {
		score = 0
	}
	m.ComplianceScore = score
}

// extractPlatformInfo extracts platform-specific version info
func (m *ComplianceMetadata) extractPlatformInfo(project *checker.Project) {
	// Android
	if project.GradleConfig != nil {
		if v := parseInt(project.GradleConfig.TargetSDKVersion); v > 0 {
			m.AndroidTargetSDK = v
		}
		if v := parseInt(project.GradleConfig.MinSDKVersion); v > 0 {
			m.AndroidMinSDK = v
		}
	}

	// iOS
	if project.InfoPlist != nil {
		m.IOSDeploymentTarget = "12.0" // Default assumption
	}
}

// extractPermissions extracts permission information
func (m *ComplianceMetadata) extractPermissions(project *checker.Project) {
	// Android permissions
	if project.AndroidManifest != nil {
		for _, perm := range project.AndroidManifest.Permissions {
			// Extract just the permission name (not the full android.permission.*)
			parts := strings.Split(perm, ".")
			if len(parts) > 0 {
				m.AndroidPermissions = append(m.AndroidPermissions, parts[len(parts)-1])
			}
		}
	}

	// iOS permissions (based on presence in InfoPlist)
	if project.InfoPlist != nil {
		if project.InfoPlist.HasCameraUsageDescription {
			m.IOSPermissions = append(m.IOSPermissions, "Camera")
		}
		if project.InfoPlist.HasLocationUsageDescription {
			m.IOSPermissions = append(m.IOSPermissions, "Location")
		}
		if project.InfoPlist.HasPhotoLibraryUsageDescription {
			m.IOSPermissions = append(m.IOSPermissions, "PhotoLibrary")
		}
		if project.InfoPlist.HasMicrophoneUsageDescription {
			m.IOSPermissions = append(m.IOSPermissions, "Microphone")
		}
		if project.InfoPlist.HasContactsUsageDescription {
			m.IOSPermissions = append(m.IOSPermissions, "Contacts")
		}
		if project.InfoPlist.HasCalendarsUsageDescription {
			m.IOSPermissions = append(m.IOSPermissions, "Calendar")
		}
	}
}

// extractDependencies extracts dependency names only (no versions)
func (m *ComplianceMetadata) extractDependencies(project *checker.Project) {
	if project.Pubspec == nil {
		return
	}

	// Main dependencies
	for name := range project.Pubspec.Dependencies {
		m.Dependencies = append(m.Dependencies, name)
	}

	// Dev dependencies
	for name := range project.Pubspec.DevDependencies {
		m.DevDependencies = append(m.DevDependencies, name)
	}
}

// extractFeatures extracts feature flags from project
func (m *ComplianceMetadata) extractFeatures(project *checker.Project) {
	// From project flags
	m.Features.HasLogin = project.HasLoginPatterns
	m.Features.HasCamera = project.HasCameraDeps
	m.Features.HasLocation = project.HasLocationDeps
	m.Features.HasPhotoLibrary = project.HasImagePicker

	// Detect from dependencies
	if project.Pubspec != nil {
		for dep := range project.Pubspec.Dependencies {
			depLower := strings.ToLower(dep)
			switch {
			case depLower == "microphone" || depLower == "record":
				m.Features.HasMicrophone = true
			case depLower == "contacts_service" || depLower == "contacts":
				m.Features.HasContacts = true
			case depLower == "device_calendar":
				m.Features.HasCalendar = true
			case depLower == "flutter_blue" || depLower == "bluetooth":
				m.Features.HasBluetooth = true
			case depLower == "firebase_messaging" || depLower == "flutter_local_notifications":
				m.Features.HasNotifications = true
			case depLower == "in_app_purchase" || depLower == "purchases_flutter":
				m.Features.HasInAppPurchase = true
			case depLower == "google_sign_in" || depLower == "sign_in_with_apple":
				m.Features.HasSocialLogin = true
			case depLower == "firebase_auth":
				m.Features.HasEmailLogin = true
			}
		}
	}

	// Detect from findings
	for _, finding := range m.Findings {
		switch finding.ID {
		case "POL-003":
			m.Features.HasDataDeletion = true
		case "POL-001":
			m.Features.HasPrivacyPolicy = true
		case "POL-002":
			m.Features.HasTermsOfService = true
		}
	}
}

// extractConfig extracts configuration flags
func (m *ComplianceMetadata) extractConfig(project *checker.Project) {
	if project.Pubspec != nil {
		m.Config.HasLinter = project.Pubspec.HasLinter
		m.Config.HasIconConfig = project.Pubspec.HasIconConfig
		m.Config.HasSplashConfig = project.Pubspec.HasSplashConfig
		m.Config.HasDeprecatedDeps = project.Pubspec.HasDeprecatedPkg
		m.Config.HasDebugDeps = project.Pubspec.HasDebugDeps
	}

	// Material 3 detection - assume true for newer Flutter projects
	m.Config.UsesMaterial3 = true
}

// extractSecurityFlags extracts security-related flags
func (m *ComplianceMetadata) extractSecurityFlags(project *checker.Project, findings []report.Finding) {
	// From manifest
	if project.AndroidManifest != nil {
		m.Security.IsDebuggable = project.AndroidManifest.Debuggable
		m.Security.AllowsBackup = project.AndroidManifest.AllowBackup
	}

	// From findings
	for _, f := range findings {
		switch f.ID {
		case "SEC-003": // Insecure HTTP
			m.Security.HasInsecureHTTP = true
		case "SEC-001": // Hardcoded credentials
			m.Security.HasHardcodedKeys = true
		case "IOS-011": // Missing encryption declaration
			m.Security.MissingEncryption = true
		}
	}
}

// ToJSON serializes metadata to JSON
func (m *ComplianceMetadata) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// ToJSONCompact serializes metadata to compact JSON
func (m *ComplianceMetadata) ToJSONCompact() ([]byte, error) {
	return json.Marshal(m)
}

// Size returns the approximate size of the metadata in bytes
func (m *ComplianceMetadata) Size() int {
	data, _ := m.ToJSONCompact()
	return len(data)
}

// Helper functions

func extractCategory(id string) string {
	if len(id) < 3 {
		return "Other"
	}

	switch id[:3] {
	case "AND":
		return "Android"
	case "IOS":
		return "iOS"
	case "FLT":
		return "Flutter"
	case "SEC":
		return "Security"
	case "POL":
		return "Policy"
	case "COD":
		return "Code Quality"
	case "TST":
		return "Testing"
	case "LIN":
		return "Linting"
	case "DOC":
		return "Documentation"
	case "PER":
		return "Performance"
	case "REV":
		return "Reviewer"
	case "AI-":
		return "AI Analysis"
	default:
		return "Other"
	}
}

func parseInt(s string) int {
	var result int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int(ch-'0')
		} else {
			break
		}
	}
	return result
}

func sanitizeDescription(desc string) string {
	// Limit description length
	if len(desc) > 200 {
		desc = desc[:200] + "..."
	}
	// Remove potentially sensitive info (case-insensitive)
	lower := strings.ToLower(desc)
	if strings.Contains(lower, "password") || strings.Contains(lower, "secret") {
		return "[Description contains sensitive keywords]"
	}
	return desc
}
