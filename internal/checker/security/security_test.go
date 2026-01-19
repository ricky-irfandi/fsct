package security

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestHardcodedCredentialsCheck_ID(t *testing.T) {
	c := &HardcodedCredentialsCheck{}
	if c.ID() != "SEC-001" {
		t.Errorf("expected SEC-001, got %s", c.ID())
	}
}

func TestHardcodedCredentialsCheck_Name(t *testing.T) {
	c := &HardcodedCredentialsCheck{}
	if c.Name() != "Hardcoded Credentials Detection" {
		t.Errorf("unexpected name: %s", c.Name())
	}
}

func TestHardcodedCredentialsCheck_Run(t *testing.T) {
	c := &HardcodedCredentialsCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("empty dart files list", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("detects api_key", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"const String api_key = '12345678901234567890';"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("detects secret", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"final String secret = 'abc123XYZ456defGHI';"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("detects password", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"String password = 'securepass123';"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("detects auth_token", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"String auth_token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9';"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("no false positives", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"void main() { print('Hello'); }"},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}

func TestDebugModeCheck_ID(t *testing.T) {
	c := &DebugModeCheck{}
	if c.ID() != "SEC-002" {
		t.Errorf("expected SEC-002, got %s", c.ID())
	}
}

func TestDebugModeCheck_Run(t *testing.T) {
	c := &DebugModeCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("detects print statement", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"void main() { print('debug'); }"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("detects debugPrint", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"debugPrint('log message');"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("no false positives", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"void main() { return; }"},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}

func TestInsecureHTTPCheck_ID(t *testing.T) {
	c := &InsecureHTTPCheck{}
	if c.ID() != "SEC-003" {
		t.Errorf("expected SEC-003, got %s", c.ID())
	}
}

func TestInsecureHTTPCheck_Run(t *testing.T) {
	c := &InsecureHTTPCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("detects http url", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"String url = 'http://example.com/api';"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("ignores localhost", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"String url = 'http://localhost:8080/api';"},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("ignores 127.0.0.1", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"String url = 'http://127.0.0.1:3000/api';"},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("https is allowed", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"String url = 'https://api.example.com';"},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}

func TestExportedActivityCheck_ID(t *testing.T) {
	c := &ExportedActivityCheck{}
	if c.ID() != "SEC-004" {
		t.Errorf("expected SEC-004, got %s", c.ID())
	}
}

func TestExportedActivityCheck_Run(t *testing.T) {
	c := &ExportedActivityCheck{}

	t.Run("no manifest", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("safe activities", func(t *testing.T) {
		project := &checker.Project{
			AndroidManifest: &checker.AndroidManifestInfo{
				Activities: []checker.ActivityInfo{
					{
						Name:            ".MainActivity",
						Exported:        true,
						HasIntentFilter: true,
					},
				},
			},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("unsafe exported activity", func(t *testing.T) {
		project := &checker.Project{
			AndroidManifest: &checker.AndroidManifestInfo{
				Activities: []checker.ActivityInfo{
					{
						Name:            ".MainActivity",
						Exported:        true,
						HasIntentFilter: false,
					},
				},
			},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})
}

func TestSQLInjectionCheck_ID(t *testing.T) {
	c := &SQLInjectionCheck{}
	if c.ID() != "SEC-005" {
		t.Errorf("expected SEC-005, got %s", c.ID())
	}
}

func TestSQLInjectionCheck_Run(t *testing.T) {
	c := &SQLInjectionCheck{}

	t.Run("no dart files", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("detects rawQuery SELECT", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"await rawQuery('SELECT * FROM users');"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("detects rawQuery INSERT", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"rawQuery('INSERT INTO users (name) VALUES (?)');"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("detects rawQuery UPDATE", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"rawQuery('UPDATE users SET name = ?');"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("detects rawQuery DELETE", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"rawQuery('DELETE FROM users WHERE id = ?');"},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("no false positives", func(t *testing.T) {
		project := &checker.Project{
			DartFiles: []string{"final query = 'SELECT';"},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}
