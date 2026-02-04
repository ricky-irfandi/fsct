# FSCT - Flutter Store Compliance Tool

[![Go Version](https://img.shields.io/badge/Go-1.21-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-200%2B-success)](internal/checker)

A comprehensive static analysis tool for Flutter mobile applications to ensure compliance with Google Play Store and Apple App Store guidelines.

## Features

- **38 Core Compliance Checks** focused on store review compliance
- **Optional AI & Reviewer Checks** when configured
- **Multi-platform Support** - Android, iOS, and Flutter
- **Multiple Output Formats** - Console, JSON, YAML, HTML
- **CI/CD Ready** - Exit codes for automated pipelines
- **Severity Filtering** - Filter by high, warning, or info severity
- **Fast & Parallel** - Concurrent check execution

## Quick Start

### Installation

```bash
go install github.com/rickyirfandi/fsct@latest
```

### Basic Usage

```bash
# Run checks on a Flutter project
fsct check ./my_flutter_app

# JSON output for CI/CD
fsct check ./my_flutter_app --format json --output report.json

# Show only high severity issues
fsct check ./my_flutter_app --severity high

# Verbose mode
fsct check ./my_flutter_app -v

# CI mode (fails on high severity issues)
fsct check ./my_flutter_app --ci
```

## Command Options

```bash
fsct check [path] [flags]

Flags:
  --platform string   Platform to check: android, ios, or both (default "both")
  --format string     Output format: console, json, yaml, html (default "console")
  --output string     Output file path (default: stdout)
  --ci                CI mode: minimal output with exit codes
  -v, --verbose       Verbose output
  --skip string       Comma-separated list of check IDs to skip
  --severity string   Minimum severity to report: info, warning, high (default "info")
  --checks string     Comma-separated list of check IDs to run
```

## Check Categories

| Category | Description | Checks |
|----------|-------------|--------|
| Android | Google Play Store requirements | 12 |
| iOS | Apple App Store requirements | 12 |
| Flutter | Store-critical Flutter config | 4 |
| Security | Security vulnerabilities | 5 |
| Policy | Policy compliance | 5 |

### Android Checks (AND-001 to AND-012)

- **AND-001**: Target SDK Version (requires 34+)
- **AND-002**: Minimum SDK Version (recommends 21+)
- **AND-003**: Internet Permission
- **AND-004**: Dangerous Permissions
- **AND-005**: Debuggable flag
- **AND-006**: Exported Activities
- **AND-007**: Missing App Icons
- **AND-008**: Placeholder Icons
- **AND-009**: Application ID Format
- **AND-010**: Version Code
- **AND-011**: Package Visibility
- **AND-012**: Allow Backup

### iOS Checks (IOS-001 to IOS-012)

- **IOS-001**: Camera Usage Description
- **IOS-002**: Photo Library Usage Description
- **IOS-003**: Location Usage Description
- **IOS-004**: Microphone Usage Description
- **IOS-005**: Contacts Usage Description
- **IOS-006**: Calendars Usage Description
- **IOS-007**: Privacy Usage Text
- **IOS-008**: Launch Screen
- **IOS-009**: App Icon (1024px)
- **IOS-010**: Full Screen Conflict
- **IOS-011**: Encryption Declaration
- **IOS-012**: Deployment Target

### Flutter Checks (Store-Critical)

- **FLT-001**: Flutter SDK Version Constraint
- **FLT-003**: Min SDK Version
- **FLT-004**: Package Name Validation
- **FLT-005**: Version Management

### Security Checks (SEC-001 to SEC-005)

- **SEC-001**: Hardcoded Credentials
- **SEC-002**: Debug Mode
- **SEC-003**: Insecure HTTP URLs
- **SEC-004**: Exported Activities
- **SEC-005**: SQL Injection

## Output Examples

### Console Output

```
FSCT Report
==========

Summary:
  High:    2
  Warning: 5
  Info:    3
  Passed:  88

Findings:

[✗] AND-001 (HIGH)
    Title: Target SDK Version Check
    Message: Target SDK version is 31. Google Play Store requires targetSdkVersion 35+.
    File: android/app/build.gradle:24
    Suggestion: Update targetSdkVersion to 34 or higher

[⚠] AND-002 (WARNING)
    Title: Minimum SDK Version Check
    Message: Minimum SDK version is 16. Consider updating to API 21+ for better security.
    File: android/app/build.gradle:25
```

### JSON Output

```json
{
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z",
  "summary": {
    "high": 2,
    "warning": 5,
    "info": 3,
    "passed": 88
  },
  "findings": [
    {
      "id": "AND-001",
      "severity": "HIGH",
      "title": "Target SDK Version Check",
      "message": "Target SDK version is 31...",
      "file": "android/app/build.gradle",
      "line": 24
    }
  ]
}
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Compliance Check
on: [push, pull_request]

jobs:
  fsct:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run FSCT
        run: |
          go install github.com/rickyirfandi/fsct@latest
          fsct check . --ci --format json --output fsct-report.json
      - name: Upload Report
        if: failure()
        uses: actions/upload-artifact@v3
        with:
          name: fsct-report
          path: fsct-report.json
```

### GitLab CI

```yaml
fsct:
  image: golang:1.21
  script:
    - go install github.com/rickyirfandi/fsct@latest
    - fsct check . --ci || echo "Compliance issues found"
  artifacts:
    when: on_failure
    paths:
      - fsct-report.json
```

## Configuration

### Ignoring Checks

```bash
# Skip specific checks
fsct check . --skip AND-001,AND-002

# Run only specific checks
fsct check . --checks AND-001,FLT-001
```

### Severity Filtering

```bash
# Only show high severity
fsct check . --severity high

# Show high and warning
fsct check . --severity warning
```

## Development

### Building

```bash
git clone https://github.com/rickyirfandi/fsct.git
cd fsct
make build
```

### Testing

```bash
make test        # Run all tests
make test-cover  # Run tests with coverage
make lint        # Run linter
```

### Adding New Checks

1. Create a new file in `internal/checker/<category>/`
2. Implement the `Check` interface:

```go
type MyCheck struct{}

func (c *MyCheck) ID() string {
    return "CAT-001"
}

func (c *MyCheck) Name() string {
    return "My Check Title"
}

func (c *MyCheck) Run(project *checker.Project) []report.Finding {
    var findings []report.Finding
    // Check logic here
    return findings
}
```

3. Register in `internal/registry/registry.go`
4. Add tests in `<category>_test.go`

## Architecture

```
fsct/
├── cmd/
│   └── fsct/           # CLI entrypoint
├── internal/
│   ├── checker/        # Check implementations
│   │   ├── android/    # Android checks
│   │   ├── ios/        # iOS checks
│   │   ├── flutter/    # Flutter checks
│   │   ├── security/   # Security checks
│   │   ├── policy/     # Policy checks
│   │   ├── code/       # Code quality checks
│   │   ├── testing/    # Testing checks
│   │   ├── linting/    # Linting checks
│   │   ├── docs/       # Documentation checks
│   │   └── perf/       # Performance checks
│   ├── parser/         # File parsers
│   ├── registry/       # Check registry
│   ├── formatter/      # Output formatters
│   ├── filter/         # Result filters
│   ├── report/         # Report models
│   └── config/         # Configuration
├── docs/               # Documentation
├── examples/           # Example projects
├── testdata/           # Test fixtures
└── Makefile
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- [GitHub Issues](https://github.com/rickyirfandi/fsct/issues)
- [Documentation](docs/)
