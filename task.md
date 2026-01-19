# FSCT — Atomic Task List for Multi-Agent Development

> **Project**: Flutter Store Compliance Tool (FSCT)
> **Language**: Go
> **Total Checks**: 52 (42 static + 5 AI + 5 reviewer verification)
> **Total Estimated Days**: 20 working days

---
## TUI
Make the terminal interface clean and beautiful, something that steve jobs and jony ive would approve

## 📋 HOW TO USE THIS TASK LIST

1. **Each task is ATOMIC** — can be completed independently
2. **Dependencies are EXPLICIT** — check `DEPENDS:` before starting
3. **Files are SPECIFIC** — exact file paths are provided
4. **Acceptance criteria** — each task has clear "DONE WHEN" conditions
5. **Assign by EPIC** — each agent should own one EPIC at a time

### Task Status Legend
- `[ ]` Not started
- `[/]` In progress (agent assigned)
- `[x]` Completed
- `[!]` Blocked (see notes)

---

# EPIC 00 — Project Setup & Foundation
> **ETD**: Day 1-2 (2 days)
> **Owner**: _unassigned_
> **Dependencies**: None

## E00-T01: Initialize Go Module
- **Status**: `[x]`
- **File**: `go.mod`, `go.sum`
- **Action**:
  ```bash
  mkdir -p /path/to/fsct
  cd /path/to/fsct
  go mod init github.com/ricky-irfandi/fsct
  ```
- **DONE WHEN**: `go.mod` exists with module name `github.com/ricky-irfandi/fsct`

## E00-T02: Create Directory Structure
- **Status**: `[x]`
- **DEPENDS**: E00-T01
- **Action**: Create all directories as specified
- **Files to create**:
  ```
  cmd/fsct/main.go
  internal/parser/.gitkeep
  internal/checker/android/.gitkeep
  internal/checker/ios/.gitkeep
  internal/checker/flutter/.gitkeep
  internal/checker/common/.gitkeep
  internal/checker/reviewer/.gitkeep
  internal/report/.gitkeep
  internal/config/.gitkeep
  internal/ai/.gitkeep
  testdata/.gitkeep
  ```
- **DONE WHEN**: All directories exist, `tree` command shows correct structure

## E00-T03: Add Core Dependencies
- **Status**: `[x]`
- **DEPENDS**: E00-T01
- **File**: `go.mod`
- **Action**: Add dependencies
  ```bash
  go get github.com/spf13/cobra@v1.8.0
  go get github.com/fatih/color@v1.16.0
  go get gopkg.in/yaml.v3@v3.0.1
  go get howett.net/plist@v1.0.1
  ```
- **DONE WHEN**: `go mod tidy` succeeds, all deps in `go.sum`

## E00-T04: Create CLI Entrypoint with Root Command
- **Status**: `[x]`
- **DEPENDS**: E00-T02, E00-T03
- **File**: `cmd/fsct/main.go`
- **Action**: Create main.go with cobra root command
- **DONE WHEN**: `go run cmd/fsct/main.go --help` prints usage

## E00-T05: Create `check` Subcommand Skeleton
- **Status**: `[x]`
- **DEPENDS**: E00-T04
- **File**: `cmd/fsct/check.go`
- **Flags to implement**:
  - `--project <path>` (required)
  - `--platform <android|ios|both>` (default: both)
  - `--format <text|json>` (default: text)
  - `--output <path>` (default: stdout)
  - `--ci` (boolean)
  - `--verbose` (boolean)
  - `--skip <check-ids>` (comma-separated)
- **DONE WHEN**: `go run cmd/fsct/main.go check --help` shows all flags

## E00-T06: Create `checks` Subcommand (List All Checks)
- **Status**: `[x]`
- **DEPENDS**: E00-T04
- **File**: `cmd/fsct/checks.go`
- **Action**: Create command that lists all available checks
- **DONE WHEN**: `go run cmd/fsct/main.go checks` prints table of check IDs

## E00-T07: Create Report Model
- **Status**: `[x]`
- **DEPENDS**: E00-T02
- **File**: `internal/report/model.go`
- **Structs to define**:
  ```go
  type Severity string // INFO, WARNING, HIGH
  type Finding struct { ID, Severity, Title, Message, File, Suggestion string; Line int }
  type Report struct { Version, Timestamp, Project string; Summary Summary; Findings []Finding }
  type Summary struct { High, Warning, Info, Passed int }
  ```
- **DONE WHEN**: Structs compile, can be JSON marshaled

## E00-T08: Create Config Model
- **Status**: `[x]`
- **DEPENDS**: E00-T02
- **File**: `internal/config/config.go`
- **Structs to define**: Config, AIConfig, ReviewerConfig
- **DONE WHEN**: Can parse `.fsct.yaml` into Config struct

## E00-T09: Add Makefile
- **Status**: `[x]`
- **DEPENDS**: E00-T04
- **File**: [Makefile](file:///Users/rickyirfandi/Documents/LOGIQUE/flutter-store-compliance-tool/Makefile)
- **Targets**: `build`, `test`, `lint`, `clean`, `install`
- **DONE WHEN**: `make build` produces `./bin/fsct` binary

## E00-T10: Add .gitignore
- **Status**: `[x]`
- **File**: [.gitignore](file:///Users/rickyirfandi/Documents/LOGIQUE/flutter-store-compliance-tool/.gitignore)
- **Content**: bin/, *.exe, .DS_Store, *.test, coverage.out
- **DONE WHEN**: File exists with appropriate entries

---

# EPIC 01 — File Parsers
> **ETD**: Day 3-4 (2 days)
> **Owner**: _completed_
> **Dependencies**: EPIC 00 complete

**Status**: All tasks completed ✅

## E01-T01: Android Manifest Parser
- **Status**: `[x]`
- **DEPENDS**: E00-T02
- **File**: `internal/parser/manifest.go`
- **Input**: `android/app/src/main/AndroidManifest.xml`
- **Must extract**:
  - Package name
  - All `<uses-permission>` elements
  - All `<activity>` elements with intent-filters
  - `android:exported` attributes
  - `android:debuggable` attribute
  - `android:allowBackup` attribute
- **DONE WHEN**: Unit test parses sample manifest correctly

## E01-T02: Android Manifest Parser Tests
- **Status**: `[x]`
- **DEPENDS**: E01-T01
- **File**: `internal/parser/manifest_test.go`
- **Fixtures**: `testdata/android/valid_manifest.xml`, `testdata/android/invalid_manifest.xml`
- **DONE WHEN**: `go test ./internal/parser/... -run Manifest` passes

## E01-T03: Gradle Parser
- **Status**: `[x]`
- **DEPENDS**: E00-T02
- **File**: `internal/parser/gradle.go`
- **Input**: `android/app/build.gradle` or `android/app/build.gradle.kts`
- **Must extract** (regex-based):
  - `applicationId`
  - `minSdkVersion`
  - `targetSdkVersion`
  - `versionCode`
  - `versionName`
- **DONE WHEN**: Unit test parses sample build.gradle correctly

## E01-T04: Gradle Parser Tests
- **Status**: `[x]`
- **DEPENDS**: E01-T03
- **File**: `internal/parser/gradle_test.go`
- **Fixtures**: `testdata/android/build.gradle`, `testdata/android/build.gradle.kts`
- **DONE WHEN**: `go test ./internal/parser/... -run Gradle` passes

## E01-T05: Info.plist Parser
- **Status**: `[x]`
- **DEPENDS**: E00-T02, E00-T03
- **File**: `internal/parser/plist.go`
- **Input**: `ios/Runner/Info.plist`
- **Must extract**:
  - `CFBundleIdentifier`
  - `CFBundleVersion`
  - `CFBundleShortVersionString`
  - All `NS*UsageDescription` keys
  - `ITSAppUsesNonExemptEncryption`
  - `UIRequiresFullScreen`
- **DONE WHEN**: Unit test parses sample Info.plist correctly

## E01-T06: Info.plist Parser Tests
- **Status**: `[x]`
- **DEPENDS**: E01-T05
- **File**: `internal/parser/plist_test.go`
- **Fixtures**: `testdata/ios/Info.plist`
- **DONE WHEN**: `go test ./internal/parser/... -run Plist` passes

## E01-T07: pubspec.yaml Parser
- **Status**: `[x]`
- **DEPENDS**: E00-T02
- **File**: `internal/parser/pubspec.go`
- **Input**: `pubspec.yaml`
- **Must extract**:
  - `name`
  - `version`
  - `description`
  - `homepage` / `repository`
  - `dependencies` (list)
  - `dev_dependencies` (list)
- **DONE WHEN**: Unit test parses sample pubspec.yaml correctly

## E01-T08: pubspec.yaml Parser Tests
- **Status**: `[x]`
- **DEPENDS**: E01-T07
- **File**: `internal/parser/pubspec_test.go`
- **Fixtures**: `testdata/flutter/pubspec.yaml`
- **DONE WHEN**: `go test ./internal/parser/... -run Pubspec` passes

## E01-T09: Dart File Scanner
- **Status**: `[x]`
- **DEPENDS**: E00-T02
- **File**: `internal/parser/dart.go`
- **Action**: Recursively scan `lib/**/*.dart` files
- **Must support**: Regex pattern matching with file:line context
- **DONE WHEN**: Can find pattern matches across multiple .dart files

## E01-T10: Dart File Scanner Tests
- **Status**: `[x]`
- **DEPENDS**: E01-T09
- **File**: `internal/parser/dart_test.go`
- **Fixtures**: `testdata/flutter/lib/sample.dart`
- **DONE WHEN**: `go test ./internal/parser/... -run Dart` passes

---

# EPIC 02 — Android Checks (12 checks)
> **ETD**: Day 5-6 (2 days)
> **Owner**: _completed_
> **Dependencies**: EPIC 01 complete

**Status**: All tasks completed ✅

## E02-T01: Check Interface Definition
- **Status**: `[x]`
- **DEPENDS**: E00-T07
- **File**: `internal/checker/checker.go`
- **File**: `internal/checker/checker.go`
- **Action**: Define common interface
  ```go
  type Check interface {
    ID() string
    Name() string
    Run(project *Project) []Finding
  }
  type Project struct { Path string; Manifest *Manifest; Plist *Plist; ... }
  ```
- **DONE WHEN**: Interface compiles, can be used by check implementations

## E02-T02: AND-001 Target SDK Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T03
- **File**: `internal/checker/android/sdk.go`
- **File**: `internal/checker/android/sdk.go`
- **Logic**: If `targetSdkVersion < 34` → HIGH
- **DONE WHEN**: Test confirms finding generated for SDK 31

## E02-T03: AND-002 Min SDK Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T03
- **File**: `internal/checker/android/sdk.go`
- **File**: `internal/checker/android/sdk.go`
- **Logic**: If `minSdkVersion < 21` → WARNING
- **DONE WHEN**: Test confirms finding generated for SDK 16

## E02-T04: AND-003 Internet Permission Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T01
- **File**: `internal/checker/android/permissions.go`
- **File**: `internal/checker/android/permissions.go`
- **Logic**: If app uses network deps but no INTERNET permission → WARNING
- **DONE WHEN**: Test with http package but no permission triggers warning

## E02-T05: AND-004 Dangerous Permissions Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T01
- **File**: `internal/checker/android/permissions.go`
- **File**: `internal/checker/android/permissions.go`
- **Logic**: CAMERA/LOCATION without uses-feature → WARNING
- **DONE WHEN**: Test confirms warning for CAMERA without uses-feature

## E02-T06: AND-005 Debuggable Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T01
- **File**: `internal/checker/android/metadata.go`
- **File**: `internal/checker/android/metadata.go`
- **Logic**: If `android:debuggable="true"` → HIGH
- **DONE WHEN**: Test confirms HIGH for debuggable manifest

## E02-T07: AND-006 Exported Attribute Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T01
- **File**: `internal/checker/android/metadata.go`
- **File**: `internal/checker/android/metadata.go`
- **Logic**: Activity with intent-filter but no `android:exported` → HIGH
- **DONE WHEN**: Test confirms HIGH for missing exported

## E02-T08: AND-007 Missing App Icon Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01
- **File**: `internal/checker/android/icons.go`
- **File**: `internal/checker/android/icons.go`
- **Logic**: No `ic_launcher` in any `res/mipmap-*` folder → HIGH
- **DONE WHEN**: Test confirms HIGH when no icons exist

## E02-T09: AND-008 Placeholder Icon Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01
- **File**: `internal/checker/android/icons.go`
- **File**: `internal/checker/android/icons.go`
- **Logic**: Default Flutter icon detected (hash comparison) → WARNING
- **DONE WHEN**: Test confirms WARNING for default Flutter icon

## E02-T10: AND-009 Application ID Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T03
- **File**: `internal/checker/android/metadata.go`
- **File**: `internal/checker/android/metadata.go`
- **Logic**: `applicationId` is `com.example.*` → HIGH
- **DONE WHEN**: Test confirms HIGH for com.example.myapp

## E02-T11: AND-010 Version Code Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T03
- **File**: `internal/checker/android/metadata.go`
- **File**: `internal/checker/android/metadata.go`
- **Logic**: `versionCode` is 1 → WARNING
- **DONE WHEN**: Test confirms WARNING for versionCode 1

## E02-T12: AND-011 Package Visibility Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T01
- **File**: `internal/checker/android/permissions.go`
- **File**: `internal/checker/android/permissions.go`
- **Logic**: Uses `url_launcher` but no `<queries>` → WARNING
- **DONE WHEN**: Test confirms WARNING when queries missing

## E02-T13: AND-012 Allow Backup Check
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T01
- **File**: `internal/checker/android/metadata.go`
- **File**: `internal/checker/android/metadata.go`
- **Logic**: `android:allowBackup="true"` without encryption → WARNING
- **DONE WHEN**: Test confirms WARNING for allowBackup=true

## E02-T14: Android Checks Unit Tests
- **Status**: `[x]`
- **DEPENDS**: E02-T02 through E02-T13
- **File**: `internal/checker/android/*_test.go`
- **File**: `internal/checker/android/*_test.go`
- **DONE WHEN**: `go test ./internal/checker/android/...` all pass

---

# EPIC 03 — iOS Checks (12 checks)
> **ETD**: Day 7-8 (2 days)
> **Owner**: _completed_
> **Dependencies**: EPIC 01 complete (can parallel with EPIC 02)

**Status**: All tasks completed ✅

## E03-T01: IOS-001 Camera Usage Description
- **Status**: `[x]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/privacy.go`
- **Logic**: Has camera dep but no `NSCameraUsageDescription` → HIGH
- **DONE WHEN**: Test confirms HIGH when description missing

## E03-T02: IOS-002 Photo Library Description
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/privacy.go`
- **Logic**: Has image_picker but no `NSPhotoLibraryUsageDescription` → HIGH
- **DONE WHEN**: Test confirms HIGH when description missing

## E03-T03: IOS-003 Location Description
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/privacy.go`
- **Logic**: Has location dep but no `NSLocationWhenInUseUsageDescription` → HIGH
- **DONE WHEN**: Test confirms HIGH when description missing

## E03-T04: IOS-004 Microphone Description
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/privacy.go`
- **Logic**: Has audio dep but no `NSMicrophoneUsageDescription` → HIGH
- **DONE WHEN**: Test confirms HIGH when description missing

## E03-T05: IOS-005 Contacts Description
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/privacy.go`
- **Logic**: Has contacts dep but no `NSContactsUsageDescription` → HIGH
- **DONE WHEN**: Test confirms HIGH when description missing

## E03-T06: IOS-006 Calendar Description
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/privacy.go`
- **Logic**: Has calendar dep but no `NSCalendarsUsageDescription` → HIGH
- **DONE WHEN**: Test confirms HIGH when description missing

## E03-T07: IOS-007 Empty Description Check
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/privacy.go`
- **Logic**: Any `NS*UsageDescription` is empty or placeholder → HIGH
- **Placeholder patterns**: "Add reason here", "TODO", "", "<string>"
- **DONE WHEN**: Test confirms HIGH for empty description

## E03-T08: IOS-008 Missing App Icon Check
- **Status**: `[ ]`
- **DEPENDS**: E02-T01
- **File**: `internal/checker/ios/icons.go`
- **Logic**: No `AppIcon.appiconset` folder → HIGH
- **DONE WHEN**: Test confirms HIGH when folder missing

## E03-T09: IOS-009 Missing 1024x1024 Icon
- **Status**: `[ ]`
- **DEPENDS**: E02-T01
- **File**: `internal/checker/ios/icons.go`
- **Logic**: Parse `Contents.json`, no 1024x1024 image → HIGH
- **DONE WHEN**: Test confirms HIGH when 1024 icon missing

## E03-T10: IOS-010 Full Screen Conflict
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/capabilities.go`
- **Logic**: `UIRequiresFullScreen` + iPad support → WARNING
- **DONE WHEN**: Test confirms WARNING for conflict

## E03-T11: IOS-011 Encryption Declaration
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T05
- **File**: `internal/checker/ios/capabilities.go`
- **Logic**: Missing `ITSAppUsesNonExemptEncryption` → WARNING
- **DONE WHEN**: Test confirms WARNING when key missing

## E03-T12: IOS-012 Deployment Target Check
- **Status**: `[ ]`
- **DEPENDS**: E02-T01
- **File**: `internal/checker/ios/capabilities.go`
- **Logic**: `IPHONEOS_DEPLOYMENT_TARGET < 12.0` → WARNING
- **DONE WHEN**: Test confirms WARNING for iOS 11

## E03-T13: iOS Checks Unit Tests
- **Status**: `[ ]`
- **DEPENDS**: E03-T01 through E03-T12
- **File**: `internal/checker/ios/*_test.go`
- **DONE WHEN**: `go test ./internal/checker/ios/...` all pass

---

# EPIC 04 — Flutter & Security Checks (18 checks)
> **ETD**: Day 9-10 (2 days)
> **Owner**: _unassigned_
> **Dependencies**: EPIC 01 complete

## E04-T01: FLT-001 Linter Check
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Logic**: No `flutter_lints` or `very_good_analysis` in dev_deps → INFO
- **DONE WHEN**: Test confirms INFO when no linter

## E04-T02: FLT-002 Default Version Check
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Logic**: `version: 1.0.0+1` → WARNING
- **DONE WHEN**: Test confirms WARNING for default version

## E04-T03: FLT-003 Missing Description
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Logic**: No `description` field → WARNING
- **DONE WHEN**: Test confirms WARNING when missing

## E04-T04: FLT-004 Missing Repository URL
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Logic**: No `homepage` and no `repository` → INFO
- **DONE WHEN**: Test confirms INFO when missing

## E04-T05: FLT-005 Deprecated Packages
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Deprecated list**: `package_info`, `android_alarm_manager`, `device_info`
- **Logic**: Uses deprecated package → WARNING
- **DONE WHEN**: Test confirms WARNING for package_info

## E04-T06: FLT-006 Debug Deps in Main
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Debug packages**: `flutter_driver`, `integration_test`, `mockito`
- **Logic**: Debug package in main dependencies → WARNING
- **DONE WHEN**: Test confirms WARNING for mockito in deps

## E04-T07: FLT-007 Missing Icon Config
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Logic**: No `flutter_launcher_icons` or manual icon config → WARNING
- **DONE WHEN**: Test confirms WARNING when no icon config

## E04-T08: FLT-008 No Splash Screen
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T07
- **File**: `internal/checker/flutter/pubspec.go`
- **Logic**: No `flutter_native_splash` → INFO
- **DONE WHEN**: Test confirms INFO when no splash

## E04-T09: SEC-001 Hardcoded API Keys
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/credentials.go`
- **Patterns**: `api_key`, `apiKey`, `API_KEY`, `sk-`, `pk_live_`
- **Logic**: Pattern match in .dart files → HIGH
- **DONE WHEN**: Test confirms HIGH for `apiKey = "sk-abc123"`

## E04-T10: SEC-002 Test Email Patterns
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/credentials.go`
- **Patterns**: `test@`, `example.com`, `@example.org`
- **Logic**: Pattern match in .dart files → WARNING
- **DONE WHEN**: Test confirms WARNING for test@example.com

## E04-T11: SEC-003 Hardcoded Passwords
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/credentials.go`
- **Patterns**: `password = "`, `password: "`, `"password123"`
- **Logic**: Pattern match → HIGH
- **DONE WHEN**: Test confirms HIGH for hardcoded password

## E04-T12: SEC-004 Debug Print Statements
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/debug.go`
- **Patterns**: `print(`, `debugPrint(`, `developer.log(`
- **Logic**: Pattern match → WARNING
- **DONE WHEN**: Test confirms WARNING for print()

## E04-T13: SEC-005 HTTP URLs
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/debug.go`
- **Patterns**: `http://` (exclude localhost, 10.*, 192.168.*)
- **Logic**: Non-local HTTP URL → WARNING
- **DONE WHEN**: Test confirms WARNING for http://api.example.com

## E04-T14: POL-001 Privacy Policy URL
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/policy.go`
- **Patterns**: `privacy`, `privacyPolicy`, `privacy_policy`
- **Logic**: No pattern found → WARNING
- **DONE WHEN**: Test confirms WARNING when no privacy URL

## E04-T15: POL-002 Terms of Service
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/policy.go`
- **Patterns**: `terms`, `termsOfService`, `tos`
- **Logic**: No pattern found → INFO
- **DONE WHEN**: Test confirms INFO when no TOS

## E04-T16: POL-003 Account Deletion
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/policy.go`
- **Patterns**: `deleteAccount`, `delete_account`, `removeAccount`
- **Logic**: Has login but no deletion pattern → WARNING
- **DONE WHEN**: Test confirms WARNING for login without deletion

## E04-T17: POL-004 Logout Pattern
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/policy.go`
- **Patterns**: `signOut`, `logout`, `logOut`
- **Logic**: Has login but no logout → WARNING
- **DONE WHEN**: Test confirms WARNING for login without logout

## E04-T18: POL-005 Password Recovery
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E01-T09
- **File**: `internal/checker/common/policy.go`
- **Patterns**: `forgotPassword`, `resetPassword`, `forgot_password`
- **Logic**: Has login but no recovery → WARNING
- **DONE WHEN**: Test confirms WARNING for login without recovery

## E04-T19: Flutter & Security Checks Unit Tests
- **Status**: `[ ]`
- **DEPENDS**: E04-T01 through E04-T18
- **File**: `internal/checker/flutter/*_test.go`, `internal/checker/common/*_test.go`
- **DONE WHEN**: `go test ./internal/checker/flutter/... ./internal/checker/common/...` all pass

---

# EPIC 05 — Reviewer Credentials Verification (5 checks)
> **ETD**: Day 11-12 (2 days)
> **Owner**: _unassigned_
> **Dependencies**: EPIC 01 complete

## E05-T01: Reviewer Config Model
- **Status**: `[ ]`
- **DEPENDS**: E00-T08
- **File**: `internal/config/reviewer.go`
- **Structs**:
  ```go
  type ReviewerConfig struct {
    EmailEnv     string `yaml:"email_env"`
    PasswordEnv  string `yaml:"password_env"`
    Verification *VerificationConfig `yaml:"verification"`
  }
  type VerificationConfig struct {
    Enabled         bool   `yaml:"enabled"`
    AuthEndpoint    string `yaml:"auth_endpoint"`
    Method          string `yaml:"method"`
    BodyTemplate    string `yaml:"body_template"`
    SuccessIndicator string `yaml:"success_indicator"`
  }
  ```
- **DONE WHEN**: Can parse reviewer config from .fsct.yaml

## E05-T02: REV-001 No Credentials Configured
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E05-T01
- **File**: `internal/checker/reviewer/credentials.go`
- **Logic**: No `reviewer_account` in config → WARNING
- **DONE WHEN**: Test confirms WARNING when config missing

## E05-T03: REV-002 Placeholder Email Check
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E05-T01
- **File**: `internal/checker/reviewer/credentials.go`
- **Patterns**: `test@`, `example.com`, `your-email@`
- **Logic**: Placeholder email detected → HIGH
- **DONE WHEN**: Test confirms HIGH for test@example.com

## E05-T04: REV-003 Weak Password Check
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E05-T01
- **File**: `internal/checker/reviewer/credentials.go`
- **Patterns**: `password`, `123456`, `test`, `qwerty`
- **Logic**: Weak password detected → HIGH
- **DONE WHEN**: Test confirms HIGH for password123

## E05-T05: REV-004/005 Login Verification (API Check)
- **Status**: `[ ]`
- **DEPENDS**: E02-T01, E05-T01
- **File**: `internal/checker/reviewer/verify.go`
- **Flags**: `--verify-login`
- **Logic**:
  1. Read credentials from env vars
  2. Make HTTP request to auth endpoint
  3. Check for success indicator in response
  4. REV-004: Token expired → HIGH
  5. REV-005: Login failed → HIGH
- **DONE WHEN**: Test with mock server confirms both scenarios

## E05-T06: Add --verify-login Flag
- **Status**: `[ ]`
- **DEPENDS**: E05-T05, E00-T05
- **File**: `cmd/fsct/check.go`
- **Flag**: `--verify-login` (boolean)
- **DONE WHEN**: Flag available and triggers verification

## E05-T07: Pre-Submit Checklist Generator
- **Status**: `[ ]`
- **DEPENDS**: E05-T01
- **File**: `internal/report/checklist.go`
- **Action**: Generate human-readable checklist for reviewer credentials
- **DONE WHEN**: Checklist prints credentials reminder with masked password

## E05-T08: Reviewer Checks Unit Tests
- **Status**: `[ ]`
- **DEPENDS**: E05-T02 through E05-T07
- **File**: `internal/checker/reviewer/*_test.go`
- **DONE WHEN**: `go test ./internal/checker/reviewer/...` all pass

---

# EPIC 06 — Output & Polish
> **ETD**: Day 13-14 (2 days)
> **Owner**: _unassigned_
> **Dependencies**: EPIC 02, 03, 04, 05 complete

## E06-T01: Terminal Output Formatter
- **Status**: `[ ]`
- **DEPENDS**: E00-T07
- **File**: `internal/report/terminal.go`
- **Features**:
  - Color-coded severity (red=HIGH, yellow=WARNING, blue=INFO)
  - Grouped by category (Android/iOS/Flutter/Security/Reviewer)
  - File:line references
  - Summary footer
- **DONE WHEN**: Output matches spec example

## E06-T02: JSON Output Formatter
- **Status**: `[ ]`
- **DEPENDS**: E00-T07
- **File**: `internal/report/json.go`
- **Features**:
  - Schema matches spec
  - Stable ordering (sorted by ID)
  - Pretty-printed with indent
- **DONE WHEN**: Valid JSON matching schema

## E06-T03: CI Mode Output
- **Status**: `[ ]`
- **DEPENDS**: E06-T01, E06-T02
- **File**: `internal/report/ci.go`
- **Features**:
  - Minimal, no colors
  - Exit codes: 0=OK, 1=WARNING, 2=HIGH
  - Single summary line
- **DONE WHEN**: Exit codes match severity

## E06-T04: Verbose Mode
- **Status**: `[ ]`
- **DEPENDS**: E06-T01
- **File**: `internal/report/terminal.go`
- **Features**:
  - Show all checks run (including passed)
  - Show timing per check
  - Show file scan stats
- **DONE WHEN**: `--verbose` shows extra details

## E06-T05: Check Skip Functionality
- **Status**: `[ ]`
- **DEPENDS**: E02-T01
- **File**: `internal/checker/checker.go`
- **Features**:
  - Parse `--skip AND-001,SEC-004`
  - Skip specified checks
  - Report skipped checks in output
- **DONE WHEN**: Skipped checks don't run

## E06-T06: Config File Loading
- **Status**: `[ ]`
- **DEPENDS**: E00-T08
- **File**: `internal/config/loader.go`
- **Features**:
  - Load `.fsct.yaml` from project root
  - Merge with CLI flags (CLI wins)
  - Support all config options
- **DONE WHEN**: Config file overrides work

## E06-T07: Output Formatters Unit Tests
- **Status**: `[ ]`
- **DEPENDS**: E06-T01 through E06-T06
- **File**: `internal/report/*_test.go`
- **DONE WHEN**: `go test ./internal/report/...` all pass

---

# EPIC 07 — AI Integration
> **ETD**: Day 15-17 (3 days)
> **Owner**: _unassigned_
> **Dependencies**: EPIC 06 complete

## E07-T01: OpenAI-Compatible HTTP Client
- **Status**: `[ ]`
- **DEPENDS**: E00-T02
- **File**: `internal/ai/client.go`
- **Features**:
  - Configurable base URL
  - API key from env var
  - Chat completions endpoint
  - Timeout handling
- **DONE WHEN**: Can call any OpenAI-compatible API

## E07-T02: MiniMax Provider Config
- **Status**: `[ ]`
- **DEPENDS**: E07-T01
- **File**: `internal/ai/providers/minimax.go`
- **Defaults**:
  - URL: `https://api.minimax.io/v1`
  - Model: `abab6.5s-chat`
- **DONE WHEN**: MiniMax config loads correctly

## E07-T03: Prompt Engineering
- **Status**: `[ ]`
- **DEPENDS**: E07-T01
- **File**: `internal/ai/prompt.go`
- **System prompt**: Store compliance expert persona
- **User prompt template**: Findings + context
- **DONE WHEN**: Prompt generates valid request

## E07-T04: Privacy-Safe Metadata Extraction
- **Status**: `[ ]`
- **DEPENDS**: E07-T03
- **File**: `internal/ai/metadata.go`
- **Extract**:
  - Findings (ID, severity, title only)
  - Permissions list
  - Dependencies list
  - Boolean flags (has_login, has_deletion)
- **NO SOURCE CODE**
- **DONE WHEN**: Metadata contains no code snippets

## E07-T05: AI Response Parser
- **Status**: `[ ]`
- **DEPENDS**: E07-T01
- **File**: `internal/ai/response.go`
- **Parse**: JSON response from AI
- **Extract**: insights, suggestions, risk level
- **DONE WHEN**: Can parse structured AI response

## E07-T06: AI Checks Implementation
- **Status**: `[ ]`
- **DEPENDS**: E07-T01 through E07-T05
- **File**: `internal/checker/ai/checks.go`
- **Checks**:
  - AI-001: Permission justification quality
  - AI-002: Policy compliance gaps
  - AI-003: Dependency risk assessment
  - AI-004: Store-specific guidance
  - AI-005: Suggested reviewer notes
- **DONE WHEN**: 5 AI checks generate findings

## E07-T07: AI CLI Flags
- **Status**: `[ ]`
- **DEPENDS**: E07-T06, E00-T05
- **File**: `cmd/fsct/check.go`
- **Flags**:
  - `--ai` (enable AI)
  - `--ai-url <url>`
  - `--ai-model <model>`
  - `--ai-key <key>` (shows warning, prefer env var)
  - `--offline` (skip AI)
- **DONE WHEN**: All flags work

## E07-T08: AI Output Section
- **Status**: `[ ]`
- **DEPENDS**: E07-T06, E06-T01
- **File**: `internal/report/terminal.go`
- **Features**:
  - Separate "AI Analysis" section
  - Show provider name
  - Format insights nicely
- **DONE WHEN**: AI output renders correctly

## E07-T09: AI Integration Unit Tests
- **Status**: `[ ]`
- **DEPENDS**: E07-T01 through E07-T08
- **File**: `internal/ai/*_test.go`
- **Mock**: Use mock HTTP server
- **DONE WHEN**: `go test ./internal/ai/...` all pass

---

# EPIC 08 — Distribution & Documentation
> **ETD**: Day 18-20 (3 days)
> **Owner**: _unassigned_
> **Dependencies**: All previous EPICs complete

## E08-T01: Makefile Enhancements
- **Status**: `[ ]`
- **File**: [Makefile](file:///Users/rickyirfandi/Documents/LOGIQUE/flutter-store-compliance-tool/Makefile)
- **Targets**:
  - `build`: Build for current OS
  - `build-all`: Build for linux/darwin/windows
  - `test`: Run all tests
  - `test-coverage`: Generate coverage report
  - `lint`: Run golangci-lint
  - `install`: Install to GOPATH/bin
- **DONE WHEN**: All targets work

## E08-T02: GoReleaser Configuration
- **Status**: `[ ]`
- **File**: `.goreleaser.yaml`
- **Features**:
  - Cross-compile for linux/darwin/windows (amd64/arm64)
  - Generate checksums
  - GitHub release with notes
  - Homebrew formula generation
- **DONE WHEN**: `goreleaser check` passes

## E08-T03: GitHub Actions CI
- **Status**: `[ ]`
- **File**: `.github/workflows/ci.yml`
- **Jobs**:
  - Lint
  - Test (go test ./...)
  - Build (all platforms)
- **DONE WHEN**: Workflow runs on push

## E08-T04: GitHub Actions Release
- **Status**: `[ ]`
- **File**: `.github/workflows/release.yml`
- **Trigger**: On tag push (v*)
- **Action**: Run GoReleaser
- **DONE WHEN**: Tagging creates release with binaries

## E08-T05: README.md
- **Status**: `[ ]`
- **File**: [README.md](file:///Users/rickyirfandi/Documents/LOGIQUE/flutter-store-compliance-tool/README.md)
- **Sections**:
  - Overview
  - Installation (brew, go install, binary)
  - Quick Start
  - All Checks Reference
  - Configuration
  - AI Integration
  - CI Usage
  - Contributing
- **DONE WHEN**: README is complete and accurate

## E08-T06: Example .fsct.yaml
- **Status**: `[ ]`
- **File**: `examples/.fsct.yaml`
- **Content**: All configuration options with comments
- **DONE WHEN**: Example is valid and documented

## E08-T07: Integration Tests
- **Status**: `[ ]`
- **File**: `tests/integration_test.go`
- **Tests**:
  - Full scan of sample Flutter project
  - All 47 static checks run
  - JSON output valid
  - Exit codes correct
- **DONE WHEN**: Integration tests pass

## E08-T08: Sample Flutter Project Fixture
- **Status**: `[ ]`
- **Directory**: `testdata/sample_flutter_app/`
- **Content**: Minimal Flutter project with intentional issues
- **DONE WHEN**: Can run `fsct check testdata/sample_flutter_app`

---

# 📊 EPIC Summary

| EPIC | Name | Tasks | ETD | Dependencies |
|------|------|-------|-----|--------------|
| 00 | Project Setup | 10 | Day 1-2 | None |
| 01 | File Parsers | 10 | Day 3-4 | EPIC 00 |
| 02 | Android Checks | 14 | Day 5-6 | EPIC 01 |
| 03 | iOS Checks | 13 | Day 7-8 | EPIC 01 |
| 04 | Flutter & Security | 19 | Day 9-10 | EPIC 01 |
| 05 | Reviewer Verification | 8 | Day 11-12 | EPIC 01 |
| 06 | Output & Polish | 7 | Day 13-14 | EPIC 02-05 |
| 07 | AI Integration | 9 | Day 15-17 | EPIC 06 |
| 08 | Distribution | 8 | Day 18-20 | All |

**Total**: 98 atomic tasks across 8 EPICs

---

# 🔀 Parallel Execution Guide

These EPICs can run in parallel:

```
Day 1-2:   [EPIC 00 - Setup]
              │
Day 3-4:   [EPIC 01 - Parsers]
              │
              ├───────────────┬───────────────┬───────────────┐
Day 5-6:   [EPIC 02]       [EPIC 03]       [EPIC 04]       [EPIC 05]
           (Android)        (iOS)          (Flutter)       (Reviewer)
              │               │               │               │
              └───────────────┴───────────────┴───────────────┘
                                    │
Day 13-14:                    [EPIC 06 - Output]
                                    │
Day 15-17:                    [EPIC 07 - AI]
                                    │
Day 18-20:                    [EPIC 08 - Distribution]
```

With 4 parallel agents on EPIC 02-05, timeline reduces to:
- **Sequential**: 20 days
- **4 Parallel Agents**: ~14 days
