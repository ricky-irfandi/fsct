package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ricky-irfandi/fsct/internal/checker"
	"github.com/ricky-irfandi/fsct/internal/parser"
	"github.com/ricky-irfandi/fsct/internal/report"
)

// Builder builds prompt data from a project
type Builder struct {
	project *checker.Project
	data    *PromptData
}

// NewBuilder creates a new prompt builder
func NewBuilder(project *checker.Project) *Builder {
	return &Builder{
		project: project,
		data:    NewPromptData(),
	}
}

// Build builds complete prompt data
func (b *Builder) Build(findings []report.Finding) *PromptData {
	b.buildMetadata()
	b.buildFindings(findings)
	b.buildPlatformConfig()
	b.buildPermissions()
	b.buildDependencies()
	b.buildSecurityAndPolicy()
	b.buildReviewerAccount()
	b.buildBlockers()

	// Calculate final values
	b.data.ComplianceScore = b.data.CalculateComplianceScore()
	b.data.Status = b.data.GetStatus()

	return b.data
}

// buildMetadata extracts app metadata
func (b *Builder) buildMetadata() {
	// Try to read pubspec.yaml for Flutter project info
	pubspecPath := filepath.Join(b.project.FlutterPath, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); err == nil {
		pubspec, err := parser.ParsePubspec(pubspecPath)
		if err == nil && pubspec != nil {
			b.data.AppName = pubspec.Name
			b.data.Version = pubspec.Version
			b.data.Description = pubspec.Description
			b.data.Repository = pubspec.Repository
			b.data.Homepage = pubspec.Homepage
			b.data.FlutterVersion = pubspec.GetSdkVersion()
		}
	}

	// Fallback to project path name
	if b.data.AppName == "" {
		b.data.AppName = filepath.Base(b.project.Path)
	}
	if b.data.Version == "" {
		b.data.Version = "1.0.0"
	}
}

// buildFindings processes findings into summary format
func (b *Builder) buildFindings(findings []report.Finding) {
	for _, f := range findings {
		category := extractCategory(f.ID)
		summary := FindingSummary{
			ID:         f.ID,
			Severity:   f.Severity,
			Category:   category,
			Title:      f.Title,
			Message:    f.Message,
			File:       f.File,
			Suggestion: f.Suggestion,
		}
		b.data.AddFinding(summary)
	}
}

// buildPlatformConfig extracts platform-specific configuration
func (b *Builder) buildPlatformConfig() {
	// Android config from Gradle
	if b.project.GradleConfig != nil {
		b.data.AndroidConfig.PackageName = b.project.GradleConfig.ApplicationID
		if v, err := strconv.Atoi(b.project.GradleConfig.TargetSDKVersion); err == nil {
			b.data.AndroidConfig.TargetSDKVersion = v
		}
		if v, err := strconv.Atoi(b.project.GradleConfig.MinSDKVersion); err == nil {
			b.data.AndroidConfig.MinSDKVersion = v
		}
		if v, err := strconv.Atoi(b.project.GradleConfig.VersionCode); err == nil {
			b.data.AndroidConfig.VersionCode = v
		}
		b.data.AndroidConfig.VersionName = b.project.GradleConfig.VersionName
	}

	// Android config from Manifest
	if b.project.AndroidManifest != nil {
		if b.data.AndroidConfig.PackageName == "" {
			b.data.AndroidConfig.PackageName = b.project.AndroidManifest.PackageName
		}
		b.data.AndroidConfig.IsDebuggable = b.project.AndroidManifest.Debuggable
		b.data.AndroidConfig.AllowBackup = b.project.AndroidManifest.AllowBackup
		for _, perm := range b.project.AndroidManifest.Permissions {
			if perm == "android.permission.INTERNET" {
				b.data.AndroidConfig.HasInternetPerm = true
			}
		}
	}

	// iOS config from Info.plist
	if b.project.InfoPlist != nil {
		b.data.IOSConfig.BundleIdentifier = b.project.InfoPlist.CFBundleIdentifier
		b.data.IOSConfig.BundleVersion = b.project.InfoPlist.CFBundleVersion
		b.data.IOSConfig.BundleShortVersion = b.project.InfoPlist.CFBundleShortVersionString
		b.data.IOSConfig.EncryptionDeclared = b.project.InfoPlist.EncryptionDeclarationSet
		b.data.IOSConfig.RequiresFullScreen = b.project.InfoPlist.RequiresFullScreen
		b.data.IOSConfig.HasCameraDescription = b.project.InfoPlist.HasCameraUsageDescription
		b.data.IOSConfig.HasLocationDescription = b.project.InfoPlist.HasLocationUsageDescription
		b.data.IOSConfig.HasPhotoLibDescription = b.project.InfoPlist.HasPhotoLibraryUsageDescription
		b.data.IOSConfig.HasMicrophoneDescription = b.project.InfoPlist.HasMicrophoneUsageDescription
	}

	// Flutter config from pubspec
	if b.project.Pubspec != nil {
		// UsesMaterial3 is not directly available, set to false
		b.data.FlutterConfig.UsesMaterial3 = false
		b.data.FlutterConfig.HasLinter = b.project.Pubspec.HasLinter
		b.data.FlutterConfig.HasIconConfig = b.project.Pubspec.HasIconConfig
		b.data.FlutterConfig.HasSplashConfig = b.project.Pubspec.HasSplashConfig
		b.data.FlutterConfig.HasDeprecatedPkg = b.project.Pubspec.HasDeprecatedPkg
		b.data.FlutterConfig.HasDebugDeps = b.project.Pubspec.HasDebugDeps
	}
}

// buildPermissions extracts permission information
func (b *Builder) buildPermissions() {
	// Android permissions
	if b.project.AndroidManifest != nil {
		dangerousPerms := map[string]string{
			"android.permission.CAMERA":                    "high",
			"android.permission.ACCESS_FINE_LOCATION":      "high",
			"android.permission.ACCESS_COARSE_LOCATION":    "medium",
			"android.permission.RECORD_AUDIO":              "high",
			"android.permission.READ_CONTACTS":             "high",
			"android.permission.WRITE_CONTACTS":            "high",
			"android.permission.READ_CALENDAR":             "medium",
			"android.permission.WRITE_CALENDAR":            "medium",
			"android.permission.READ_EXTERNAL_STORAGE":     "medium",
			"android.permission.WRITE_EXTERNAL_STORAGE":    "medium",
			"android.permission.INTERNET":                  "low",
			"android.permission.ACCESS_NETWORK_STATE":      "low",
			"android.permission.RECEIVE_BOOT_COMPLETED":    "low",
			"android.permission.VIBRATE":                   "low",
		}

		for _, perm := range b.project.AndroidManifest.Permissions {
			riskLevel := "low"
			if level, ok := dangerousPerms[perm]; ok {
				riskLevel = level
			}
			b.data.AndroidPermissions = append(b.data.AndroidPermissions, PermissionInfo{
				Name:      perm,
				Platform:  "android",
				Present:   true,
				RiskLevel: riskLevel,
			})
		}
	}

	// iOS permissions
	iosPerms := []struct {
		name    string
		present bool
		desc    string
	}{
		{"NSCameraUsageDescription", b.project.InfoPlist != nil && b.project.InfoPlist.HasCameraUsageDescription, ""},
		{"NSPhotoLibraryUsageDescription", b.project.InfoPlist != nil && b.project.InfoPlist.HasPhotoLibraryUsageDescription, ""},
		{"NSLocationWhenInUseUsageDescription", b.project.InfoPlist != nil && b.project.InfoPlist.HasLocationUsageDescription, ""},
		{"NSLocationAlwaysUsageDescription", false, ""},
		{"NSMicrophoneUsageDescription", b.project.InfoPlist != nil && b.project.InfoPlist.HasMicrophoneUsageDescription, ""},
		{"NSContactsUsageDescription", b.project.InfoPlist != nil && b.project.InfoPlist.HasContactsUsageDescription, ""},
		{"NSCalendarsUsageDescription", b.project.InfoPlist != nil && b.project.InfoPlist.HasCalendarsUsageDescription, ""},
	}

	for _, perm := range iosPerms {
		b.data.IOSPermissions = append(b.data.IOSPermissions, PermissionInfo{
			Name:     perm.name,
			Platform: "ios",
			Present:  perm.present,
		})
	}
}

// buildDependencies extracts dependency information
func (b *Builder) buildDependencies() {
	if b.project.Pubspec == nil {
		return
	}

	// Main dependencies
	highRiskPackages := map[string]string{
		"http":              "Uses insecure HTTP by default",
		"dio":               "Check for certificate validation",
		"permission_handler": "Ensure iOS descriptions are set",
	}

	for name, version := range b.project.Pubspec.Dependencies {
		dep := DependencyInfo{
			Name:    name,
			Version: version,
			IsDev:   false,
		}

		if _, ok := highRiskPackages[name]; ok {
			dep.RiskLevel = "high"
			b.data.HighRiskDeps = append(b.data.HighRiskDeps, dep)
		}

		b.data.Dependencies = append(b.data.Dependencies, dep)
	}

	// Dev dependencies
	for name, version := range b.project.Pubspec.DevDependencies {
		dep := DependencyInfo{
			Name:    name,
			Version: version,
			IsDev:   true,
		}
		b.data.DevDependencies = append(b.data.DevDependencies, dep)
	}
}

// buildSecurityAndPolicy extracts security and policy flags
func (b *Builder) buildSecurityAndPolicy() {
	// Security flags based on project state
	b.data.SecurityFlags = []SecurityFlag{
		{
			Type:        "Debug Mode",
			Present:     b.project.AndroidManifest != nil && b.project.AndroidManifest.Debuggable,
			Description: "App is debuggable",
			Severity:    report.SeverityHigh,
		},
		{
			Type:        "Insecure HTTP",
			Present:     b.project.HasNetworkDeps,
			Description: "App uses network dependencies",
			Severity:    report.SeverityWarning,
		},
		{
			Type:        "Exported Activities",
			Present:     false, // Will be set by security checks
			Description: "Activities are exported",
			Severity:    report.SeverityHigh,
		},
	}

	// Policy flags
	b.data.PolicyFlags = []PolicyFlag{
		{
			Type:        "Privacy Policy",
			Present:     false,
			Description: "Privacy policy URL configured",
		},
		{
			Type:        "Terms of Service",
			Present:     false,
			Description: "Terms of service URL configured",
		},
		{
			Type:        "Account Deletion",
			Present:     false,
			Description: "Account deletion functionality detected",
		},
		{
			Type:        "Logout Functionality",
			Present:     false,
			Description: "Logout functionality detected",
		},
	}

	// Detect login patterns
	if b.project.HasLoginPatterns {
		b.data.PolicyFlags[2].Present = true
		b.data.PolicyFlags[3].Present = true
	}
}

// buildReviewerAccount extracts reviewer account info
func (b *Builder) buildReviewerAccount() {
	// Check for reviewer account config in environment
	email := os.Getenv("REVIEWER_EMAIL")
	passwordEnv := "REVIEWER_PASSWORD"
	password := os.Getenv(passwordEnv)

	b.data.ReviewerAccount = &ReviewerAccountInfo{
		Configured: email != "" && password != "",
		Email:      email,
		PasswordEnv: passwordEnv,
	}

	if password != "" {
		// Check for weak password
		weakPatterns := []string{"password", "123456", "test", "qwerty", "admin", "flutter"}
		lowerPass := strings.ToLower(password)
		for _, pattern := range weakPatterns {
			if strings.Contains(lowerPass, pattern) {
				b.data.ReviewerAccount.HasWeakPassword = true
				break
			}
		}
	}
}

// buildBlockers identifies critical submission blockers
func (b *Builder) buildBlockers() {
	blockers := []BlockerInfo{}

	// Android blockers
	if b.data.AndroidConfig.TargetSDKVersion > 0 && b.data.AndroidConfig.TargetSDKVersion < 35 {
		blockers = append(blockers, BlockerInfo{
			ID:          "AND-001",
			Category:    "Android SDK",
			Title:       "Target SDK Version Too Low",
			Description: fmt.Sprintf("Target SDK is %d, Play Store requires 35+", b.data.AndroidConfig.TargetSDKVersion),
			Platform:    "android",
			FixFile:     "android/app/build.gradle",
			FixSnippet:  "targetSdkVersion 35",
		})
	}

	if b.project.AndroidManifest != nil && b.project.AndroidManifest.Debuggable {
		blockers = append(blockers, BlockerInfo{
			ID:          "AND-005",
			Category:    "Android Security",
			Title:       "Debuggable Flag Enabled",
			Description: "App has android:debuggable=\"true\" in manifest",
			Platform:    "android",
			FixFile:     "android/app/src/main/AndroidManifest.xml",
			FixSnippet:  "Remove android:debuggable attribute",
		})
	}

	// iOS blockers
	if b.data.IOSConfig.BundleIdentifier == "" || strings.HasPrefix(b.data.IOSConfig.BundleIdentifier, "com.example") {
		blockers = append(blockers, BlockerInfo{
			ID:          "IOS-BUNDLE",
			Category:    "iOS Bundle ID",
			Title:       "Invalid Bundle Identifier",
			Description: "Bundle ID is missing or uses com.example prefix",
			Platform:    "ios",
			FixFile:     "ios/Runner.xcodeproj/project.pbxproj",
			FixSnippet:  "Set unique PRODUCT_BUNDLE_IDENTIFIER",
		})
	}

	// Check for missing critical iOS descriptions when permissions are present
	hasCameraDep := b.project.HasCameraDeps
	hasLocationDep := b.project.HasLocationDeps

	if hasCameraDep && !b.data.IOSConfig.HasCameraDescription {
		blockers = append(blockers, BlockerInfo{
			ID:          "IOS-001",
			Category:    "iOS Privacy",
			Title:       "Missing Camera Usage Description",
			Description: "App uses camera but NSCameraUsageDescription is missing",
			Platform:    "ios",
			FixFile:     "ios/Runner/Info.plist",
			FixSnippet:  "<key>NSCameraUsageDescription</key>\n<string>This app needs camera access to...</string>",
		})
	}

	if hasLocationDep && !b.data.IOSConfig.HasLocationDescription {
		blockers = append(blockers, BlockerInfo{
			ID:          "IOS-003",
			Category:    "iOS Privacy",
			Title:       "Missing Location Usage Description",
			Description: "App uses location but NSLocationWhenInUseUsageDescription is missing",
			Platform:    "ios",
			FixFile:     "ios/Runner/Info.plist",
			FixSnippet:  "<key>NSLocationWhenInUseUsageDescription</key>\n<string>This app needs location access to...</string>",
		})
	}

	b.data.Blockers = blockers
}

// Helper functions

func extractCategory(id string) string {
	if strings.HasPrefix(id, "AND-") {
		return "Android"
	}
	if strings.HasPrefix(id, "IOS-") {
		return "iOS"
	}
	if strings.HasPrefix(id, "FLT-") {
		return "Flutter"
	}
	if strings.HasPrefix(id, "SEC-") {
		return "Security"
	}
	if strings.HasPrefix(id, "POL-") {
		return "Policy"
	}
	if strings.HasPrefix(id, "COD-") {
		return "Code Quality"
	}
	if strings.HasPrefix(id, "TST-") {
		return "Testing"
	}
	if strings.HasPrefix(id, "LINT-") {
		return "Linting"
	}
	if strings.HasPrefix(id, "DOC-") {
		return "Documentation"
	}
	if strings.HasPrefix(id, "PERF-") {
		return "Performance"
	}
	if strings.HasPrefix(id, "REV-") {
		return "Reviewer"
	}
	if strings.HasPrefix(id, "AI-") {
		return "AI Analysis"
	}
	return "Other"
}
