# Output Formats

FSCT supports multiple output formats for different use cases.

## Available Formats

| Format | Description | Extension | Use Case |
|--------|-------------|-----------|----------|
| console | Colored text output | - | Interactive use |
| json | Structured JSON | .json | CI/CD, automation |
| yaml | Human-readable YAML | .yaml | Configuration, logs |
| html | Styled HTML report | .html | Documentation, sharing |

---

## Console Format (Default)

Colored text output with icons and formatting for interactive use.

### Example Output

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
    Message: Minimum SDK version is 16. Consider updating to API 21+.
    File: android/app/build.gradle:25

[ℹ] FLT-001 (INFO)
    Title: Flutter SDK Version Constraint
    Message: Flutter SDK version constraint not found in pubspec.yaml
    File: pubspec.yaml
    Suggestion: Add sdk: flutter constraint to dependencies section
```

### Usage

```bash
fsct check . --format console
fsct check .  # Default format
```

---

## JSON Format

Machine-readable structured data perfect for CI/CD pipelines.

### Structure

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
      "line": 24,
      "suggestion": "Update targetSdkVersion to 34 or higher"
    }
  ]
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| version | string | FSCT version |
| timestamp | string | ISO 8601 timestamp |
| summary | object | Count of findings by severity |
| findings | array | Array of finding objects |

### Finding Object

| Field | Type | Description |
|-------|------|-------------|
| id | string | Check ID (e.g., AND-001) |
| severity | string | HIGH, WARNING, or INFO |
| title | string | Check title |
| message | string | Finding description |
| file | string | File path |
| line | int | Line number (0 if N/A) |
| suggestion | string | Recommended fix |

### Usage

```bash
# Output to file
fsct check . --format json --output report.json

# Pipeline usage
fsct check . --format json | jq '.summary.high'
```

---

## YAML Format

Human-readable serialization for configuration and logging.

### Example

```yaml
version: "1.0.0"
timestamp: "2024-01-15T10:30:00Z"
summary:
  high: 2
  warning: 5
  info: 3
  passed: 88
findings:
  - id: "AND-001"
    severity: "HIGH"
    title: "Target SDK Version Check"
    message: "Target SDK version is 31..."
    file: "android/app/build.gradle"
    line: 24
    suggestion: "Update targetSdkVersion to 34 or higher"
```

### Usage

```bash
fsct check . --format yaml --output report.yaml
```

---

## HTML Format

Styled web page with visual charts and filtering.

### Features

- Responsive design
- Color-coded severity badges
- Summary statistics cards
- Collapsible findings
- Print-friendly

### Example

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>FSCT Report</title>
    <style>
        /* Responsive CSS styling */
        .stat { padding: 20px; }
        .high { color: #f44336; }
        .warning { color: #ff9800; }
        .info { color: #2196F3; }
    </style>
</head>
<body>
    <header>
        <h1>FSCT Report</h1>
        <p>Generated at January 15, 2024 10:30</p>
    </header>
    <div class="summary">
        <div class="stat high">
            <div class="stat-value">2</div>
            <div class="stat-label">High</div>
        </div>
        <!-- More stats -->
    </div>
    <div class="findings">
        <div class="finding high">
            <span class="finding-id">AND-001</span>
            <span class="finding-severity severity-high">HIGH</span>
            <div class="finding-title">Target SDK Version Check</div>
            <!-- More details -->
        </div>
    </div>
</body>
</html>
```

### Usage

```bash
fsct check . --format html --output report.html
```

---

## Output File Naming

When using `--output`, FSCT automatically appends the appropriate extension:

```bash
# Console (no extension)
fsct check . --output report --format console

# JSON (adds .json)
fsct check . --output report.json
# or
fsct check . --output report --format json

# YAML (adds .yaml)
fsct check . --output report.yaml

# HTML (adds .html)
fsct check . --output report.html
```

---

## CI/CD Integration Examples

### GitHub Actions with JSON

```yaml
- name: Run FSCT
  run: |
    fsct check . --format json --output fsct-report.json
- name: Parse Results
  if: always()
  run: |
    HIGH=$(cat fsct-report.json | jq '.summary.high')
    if [ "$HIGH" -gt 0 ]; then
      echo "::warning::Found $HIGH high severity issues"
    fi
```

### GitLab CI with YAML

```yaml
fsct:
  script:
    - fsct check . --format yaml --output fsct-report.yaml
  artifacts:
    paths:
      - fsct-report.yaml
    when: always()
```

### Jenkins Pipeline

```groovy
stage('Compliance Check') {
    steps {
        sh 'fsct check . --format json --output fsct-report.json'
    }
    post {
        always {
            script {
                def report = readJSON file: 'fsct-report.json'
                echo "High: ${report.summary.high}"
                echo "Warning: ${report.summary.warning}"
            }
        }
    }
}
```

---

## Exit Codes

FSCT uses exit codes for CI/CD integration:

| Code | Meaning |
|------|---------|
| 0 | Success (no high severity issues) |
| 1 | Error or high severity issues found |

### CI Mode

```bash
# Use --ci flag for strict exit code behavior
fsct check . --ci

# Exits with code 1 if any HIGH severity issues
```

---

## Filtering Output

### By Severity

```bash
# Only show HIGH severity
fsct check . --severity high --format json

# Show HIGH and WARNING
fsct check . --severity warning --format json
```

### By Check ID

```bash
# Only specific checks
fsct check . --checks AND-001,IOS-001 --format json
```

### Skip Checks

```bash
# Exclude certain checks
fsct check . --skip AND-001,AND-002 --format json
```
