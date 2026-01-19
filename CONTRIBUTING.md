# Contributing to FSCT

Thank you for your interest in contributing to FSCT! This document outlines the process for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Setting Up Development Environment

1. Fork the repository on GitHub

2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/fsct.git
cd fsct
```

3. Add the upstream remote:

```bash
git remote add upstream https://github.com/rickyirfandi/fsct.git
```

4. Create a feature branch:

```bash
git checkout -b feature/your-feature-name
```

5. Verify the build:

```bash
make build
```

6. Run tests:

```bash
make test
```

## Development Workflow

### Finding Issues to Work On

- Look for issues labeled `good first issue` or `help wanted`
- Check the [GitHub Issues](https://github.com/rickyirfandi/fsct/issues) page
- Discuss your approach before implementing

### Making Changes

1. Keep changes focused and small
2. Follow Go coding conventions
3. Add tests for new functionality
4. Update documentation as needed

### Code Style

- Use `gofmt` for formatting: `make format`
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use meaningful variable and function names
- Add comments for exported types and functions

### Commit Messages

Follow conventional commits:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance tasks

Example:

```
feat(android): add SDK version range check

Implement AND-013 to validate targetSdkVersion against
Google Play Store requirements.

Closes #42
```

### Pull Request Process

1. Ensure all tests pass: `make test`
2. Update the CHANGELOG.md
3. Push your changes and create a PR
4. Address review feedback
5. Squash commits before merging

### Adding New Checks

1. Choose the appropriate category in `internal/checker/<category>/`

2. Create a new file `<check_name>.go`:

```go
package mycategory

import (
    "github.com/ricky-irfandi/fsct/internal/checker"
    "github.com/ricky-irfandi/fsct/internal/report"
)

type MyNewCheck struct{}

func (c *MyNewCheck) ID() string {
    return "CAT-XXX"
}

func (c *MyNewCheck) Name() string {
    return "Human Readable Name"
}

func (c *MyNewCheck) Run(project *checker.Project) []report.Finding {
    var findings []report.Finding

    // Check implementation
    if /* condition */ {
        findings = append(findings, project.AddFinding(
            c.ID(),
            c.Name(),
            "Issue message",
            "file.gradle",
            "Suggested fix",
            report.SeverityWarning,
            0,
        ))
    }

    return findings
}
```

3. Register in `internal/registry/registry.go`:

```go
func (r *CheckerRegistry) registerCategoryChecks() {
    r.checks["CAT-XXX"] = &mycategory.MyNewCheck{}
}
```

4. Add tests in `<check_name>_test.go`:

```go
package mycategory

import (
    "testing"

    "github.com/ricky-irfandi/fsct/internal/checker"
)

func TestMyNewCheck(t *testing.T) {
    check := &MyNewCheck{}

    t.Run("condition should generate finding", func(t *testing.T) {
        project := &checker.Project{
            // Setup project data
        }

        findings := check.Run(project)

        if len(findings) != 1 {
            t.Errorf("expected 1 finding, got %d", len(findings))
        }
    })
}
```

### Adding New Parsers

1. Create `internal/parser/<format>.go`:

```go
package parser

type MyData struct {
    Field string `yaml:"field"`
}

func ParseMyData(path string) (*MyData, error) {
    // Parse implementation
}
```

2. Update `internal/parser/parser.go` if needed

3. Add tests

### Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported functions
- Update docs/ for new features

## Testing

### Running Tests

```bash
make test              # Run all tests
make test-cover        # Run with coverage
make test-race         # Run with race detector
```

### Writing Tests

- Use table-driven tests when possible
- Test both positive and negative cases
- Test edge cases
- Use meaningful test names

Example:

```go
func TestMyCheck_Run(t *testing.T) {
    tests := []struct {
        name     string
        project  *checker.Project
        expected int
    }{
        {
            name:     "valid case - no findings",
            project:  validProject(),
            expected: 0,
        },
        {
            name:     "invalid case - one finding",
            project:  invalidProject(),
            expected: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            check := &MyCheck{}
            findings := check.Run(tt.project)
            if len(findings) != tt.expected {
                t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
            }
        })
    }
}
```

## Code Review

Expect feedback on:
- Code clarity and maintainability
- Test coverage
- Documentation completeness
- Adherence to project conventions

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create git tag
4. Build and test on all platforms
5. Create GitHub release

## Questions?

- Open an issue for bugs
- Start a discussion for feature ideas
- Reach out via GitHub Discussions
