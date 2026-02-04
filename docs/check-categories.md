# Check Categories

FSCT organizes its 38 core checks into 5 categories based on store review compliance requirements.
Optional AI and reviewer checks can be enabled when configured.

## Overview

| Category | ID Prefix | Checks | Severity |
|----------|-----------|--------|----------|
| Android | AND- | 12 | High, Warning |
| iOS | IOS- | 12 | High, Warning |
| Flutter | FLT- | 4 | High, Warning |
| Security | SEC- | 5 | Critical, High |
| Policy | POL- | 5 | High, Warning |

---

## Android Checks (AND-001 to AND-012)

These checks validate compliance with Google Play Store requirements.

### AND-001: Target SDK Version Check
- **Severity**: HIGH
- **Requirement**: targetSdkVersion must be 35 or higher
- **Google Play**: Requires API 35+ for new apps

### AND-002: Minimum SDK Version Check
- **Severity**: WARNING
- **Recommendation**: minSdkVersion 21 or higher
- **Rationale**: Better security and performance features

### AND-003: Internet Permission Check
- **Severity**: HIGH
- **Requirement**: Declare INTERNET permission if using network
- **Reference**: Android Manifest

### AND-004: Dangerous Permissions Check
- **Severity**: WARNING
- **Requirement**: Document need for dangerous permissions
- **Reference**: Android 6.0+ runtime permissions

### AND-005: Debuggable Check
- **Severity**: HIGH
- **Requirement**: debuggable must be false in release
- **Security**: Prevents debug access in production

### AND-006: Exported Activities Check
- **Severity**: HIGH
- **Requirement**: Set android:exported explicitly
- **Security**: Prevents unauthorized activity access

### AND-007: Missing App Icon Check
- **Severity**: WARNING
- **Requirement**: Provide mipmap launcher icons
- **Google Play**: Required for app listing

### AND-008: Placeholder Icon Check
- **Severity**: WARNING
- **Requirement**: Replace default launcher icons
- **Branding**: Default icons appear unprofessional

### AND-009: Application ID Check
- **Severity**: HIGH
- **Requirement**: Use valid reverse-domain notation
- **Format**: com.example.appname

### AND-010: Version Code Check
- **Severity**: WARNING
- **Requirement**: versionCode must increment for updates
- **Requirement**: versionName follows semantic versioning

### AND-011: Package Visibility Check
- **Severity**: HIGH
- **Requirement**: Declare package visibility needs
- **Android 11+**: New QUERY_ALL_PACKAGES permission

### AND-012: Allow Backup Check
- **Severity**: WARNING
- **Requirement**: Set allowBackup appropriately
- **Security**: Consider sensitive data exposure

---

## iOS Checks (IOS-001 to IOS-012)

These checks validate compliance with Apple App Store requirements.

### IOS-001: Camera Usage Description
- **Severity**: HIGH
- **Requirement**: NSCameraUsageDescription in Info.plist
- **Apple**: Required for camera access

### IOS-002: Photo Library Usage Description
- **Severity**: HIGH
- **Requirement**: NSPhotoLibraryUsageDescription
- **Apple**: Required for photo library access

### IOS-003: Location Usage Description
- **Severity**: HIGH
- **Requirement**: NSLocationUsageDescription or more specific
- **Apple**: Required for location access

### IOS-004: Microphone Usage Description
- **Severity**: HIGH
- **Requirement**: NSMicrophoneUsageDescription
- **Apple**: Required for microphone access

### IOS-005: Contacts Usage Description
- **Severity**: HIGH
- **Requirement**: NSContactsUsageDescription
- **Apple**: Required for contacts access

### IOS-006: Calendars Usage Description
- **Severity**: HIGH
- **Requirement**: NSCalendarsUsageDescription
- **Apple**: Required for calendar access

### IOS-007: Privacy Usage Text Check
- **Severity**: WARNING
- **Requirement**: Provide clear usage descriptions
- **Guidelines**: Apple reviews for clarity

### IOS-008: Launch Screen Check
- **Severity**: WARNING
- **Requirement**: Provide LaunchScreen.storyboard
- **Apple**: Required for app launch

### IOS-009: App Icon Check
- **Severity**: WARNING
- **Requirement**: 1024x1024 App Store icon
- **Apple**: Required for App Store listing

### IOS-010: Full Screen Conflict Check
- **Severity**: WARNING
- **Requirement**: UIRequiresFullScreen for iPad
- **App Store**: Multi-window requirements

### IOS-011: Encryption Declaration
- **Severity**: HIGH
- **Requirement**: Export Compliance encryption declaration
- **ITAR**: Required for apps using encryption

### IOS-012: Deployment Target Check
- **Severity**: WARNING
- **Recommendation**: iOS 13.0 or higher
- **Support**: Latest iOS features and security

---

## Flutter Checks (Store-Critical)

Flutter-specific configuration and best practices.

### FLT-001: Flutter SDK Version
- **Severity**: HIGH
- **Requirement**: Constrain Flutter SDK version
- **Build**: Reproducible builds

### FLT-003: Minimum SDK Version
- **Severity**: HIGH
- **Requirement**: Android minSdkVersion 21+
- **Compatibility**: Modern Android features

### FLT-004: Package Name Format
- **Severity**: HIGH
- **Requirement**: Valid reverse domain notation
- **Platform**: Required for app bundles

### FLT-005: Version Management
- **Severity**: WARNING
- **Requirement**: Specify version in pubspec.yaml
- **Tracking**: Semantic versioning

---

## Security Checks (SEC-001 to SEC-005)

Security vulnerability detection.

### SEC-001: Hardcoded Credentials
- **Severity**: CRITICAL
- **Detection**: API keys, passwords, tokens in code
- **Action**: Use environment variables

### SEC-002: Debug Mode Code
- **Severity**: HIGH
- **Detection**: print(), debugPrint() in production
- **Action**: Remove debug statements

### SEC-003: Insecure HTTP URLs
- **Severity**: HIGH
- **Detection**: HTTP:// instead of HTTPS://
- **Security**: Use encrypted connections

### SEC-004: Exported Activities
- **Severity**: HIGH
- **Detection**: Unprotected exported activities
- **Security**: Restrict component access

### SEC-005: SQL Injection
- **Severity**: HIGH
- **Detection**: Raw SQL with string concatenation
- **Security**: Use parameterized queries

---

## Policy Checks (POL-001 to POL-005)

App Store policy compliance.

### POL-001: Privacy Policy URL
- **Severity**: HIGH
- **Requirement**: Link to privacy policy
- **Compliance**: GDPR, CCPA requirements

### POL-002: Terms of Service URL
- **Severity**: WARNING
- **Recommendation**: Link to terms of service
- **Compliance**: App Store guidelines

### POL-003: Data Deletion Contact
- **Severity**: HIGH
- **Requirement**: Data deletion method
- **Compliance**: CCPA, GDPR Article 17

### POL-004: Logout Functionality
- **Severity**: WARNING
- **Recommendation**: User account logout
- **Guidelines**: App Store requirement

### POL-005: Account Recovery
- **Severity**: WARNING
- **Recommendation**: Password reset capability
- **Guidelines**: App Store requirement

---

## Severity Levels

### HIGH
Critical issues that may cause:
- App Store rejection
- Security vulnerabilities
- Policy violations

### WARNING
Important issues that may cause:
- App Store warnings
- Performance problems
- Maintainability issues

### INFO
Best practice recommendations:
- Code improvements
- Documentation suggestions
- Optimization tips

---

## Filtering Checks

### By Severity

```bash
# Only HIGH severity
fsct check . --severity high

# HIGH and WARNING
fsct check . --severity warning
```

### By Category

```bash
# Only Android checks
fsct check . --checks AND-001,AND-002

# Multiple categories
fsct check . --checks AND-001,IOS-001,FLT-001
```

### Skip Checks

```bash
# Skip specific checks
fsct check . --skip AND-001,IOS-001
```
