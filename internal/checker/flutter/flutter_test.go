package flutter

import (
	"testing"

	"github.com/ricky-irfandi/fsct/internal/checker"
)

func TestFlutterSDKVersionCheck_ID(t *testing.T) {
	c := &FlutterSDKVersionCheck{}
	if c.ID() != "FLT-001" {
		t.Errorf("expected FLT-001, got %s", c.ID())
	}
}

func TestFlutterSDKVersionCheck_Name(t *testing.T) {
	c := &FlutterSDKVersionCheck{}
	if c.Name() != "Flutter SDK Version Constraint" {
		t.Errorf("unexpected name: %s", c.Name())
	}
}

func TestFlutterSDKVersionCheck_Run(t *testing.T) {
	c := &FlutterSDKVersionCheck{}

	t.Run("no pubspec", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("no flutter sdk", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Dependencies: map[string]string{
					"http": "^0.13.0",
				},
			},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("has flutter sdk", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Dependencies: map[string]string{
					"flutter": "any",
					"http":    "^0.13.0",
				},
			},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}

func TestMaterial3Check_ID(t *testing.T) {
	c := &Material3Check{}
	if c.ID() != "FLT-002" {
		t.Errorf("expected FLT-002, got %s", c.ID())
	}
}

func TestMinSDKVersionCheck_ID(t *testing.T) {
	c := &MinSDKVersionCheck{}
	if c.ID() != "FLT-003" {
		t.Errorf("expected FLT-003, got %s", c.ID())
	}
}

func TestMinSDKVersionCheck_Run(t *testing.T) {
	c := &MinSDKVersionCheck{}

	t.Run("no gradle config", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("low sdk version", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				MinSDKVersion: "19",
			},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("valid sdk version", func(t *testing.T) {
		project := &checker.Project{
			GradleConfig: &checker.GradleConfigInfo{
				MinSDKVersion: "21",
			},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}

func TestPackageNameCheck_ID(t *testing.T) {
	c := &PackageNameCheck{}
	if c.ID() != "FLT-004" {
		t.Errorf("expected FLT-004, got %s", c.ID())
	}
}

func TestPackageNameCheck_Run(t *testing.T) {
	c := &PackageNameCheck{}

	t.Run("no pubspec", func(t *testing.T) {
		project := &checker.Project{}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("invalid package name", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Name: "invalidpackagename",
			},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("valid package name", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Name: "com.example.app",
			},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}

func TestVersionCheck_ID(t *testing.T) {
	c := &VersionCheck{}
	if c.ID() != "FLT-005" {
		t.Errorf("expected FLT-005, got %s", c.ID())
	}
}

func TestVersionCheck_Run(t *testing.T) {
	c := &VersionCheck{}

	t.Run("no version", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})

	t.Run("has version", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				Version: "1.0.0+1",
			},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})
}

func TestDependencyConstraintCheck_ID(t *testing.T) {
	c := &DependencyConstraintCheck{}
	if c.ID() != "FLT-006" {
		t.Errorf("expected FLT-006, got %s", c.ID())
	}
}

func TestDeprecatedPackageCheck_ID(t *testing.T) {
	c := &DeprecatedPackageCheck{}
	if c.ID() != "FLT-007" {
		t.Errorf("expected FLT-007, got %s", c.ID())
	}
}

func TestDeprecatedPackageCheck_Run(t *testing.T) {
	c := &DeprecatedPackageCheck{}

	t.Run("no deprecated", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				HasDeprecatedPkg: false,
			},
		}
		results := c.Run(project)
		if len(results) != 0 {
			t.Errorf("expected 0 findings, got %d", len(results))
		}
	})

	t.Run("has deprecated", func(t *testing.T) {
		project := &checker.Project{
			Pubspec: &checker.PubspecInfo{
				HasDeprecatedPkg: true,
			},
		}
		results := c.Run(project)
		if len(results) != 1 {
			t.Errorf("expected 1 finding, got %d", len(results))
		}
	})
}

func TestProjectStructureCheck_ID(t *testing.T) {
	c := &ProjectStructureCheck{}
	if c.ID() != "FLT-008" {
		t.Errorf("expected FLT-008, got %s", c.ID())
	}
}
