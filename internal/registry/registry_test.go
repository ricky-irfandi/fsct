package registry

import (
	"testing"

	aipkg "github.com/ricky-irfandi/fsct/internal/ai"
)

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Fatal("expected registry to be created")
	}

	if reg.checks == nil {
		t.Error("expected checks map to be initialized")
	}
}

func TestRegisterAll(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	// Should have 75 checks (no AI checks yet)
	if reg.Count() != 75 {
		t.Errorf("expected 75 checks, got %d", reg.Count())
	}
}

func TestRegisterAIChecks(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	// Create a mock AI client
	client := createMockAIClient()

	// Register AI checks
	reg.RegisterAIChecks(client)

	// Should now have 80 checks (75 + 5 AI)
	if reg.Count() != 80 {
		t.Errorf("expected 80 checks after AI registration, got %d", reg.Count())
	}

	// Check specific AI checks
	aiChecks := []string{"AI-001", "AI-002", "AI-003", "AI-004", "AI-005"}
	for _, id := range aiChecks {
		check, ok := reg.Get(id)
		if !ok {
			t.Errorf("expected AI check %s to be registered", id)
		}
		if check.ID() != id {
			t.Errorf("expected check ID %s, got %s", id, check.ID())
		}
	}
}

func TestRegisterAIChecksWithNilClient(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	// Register with nil client
	reg.RegisterAIChecks(nil)

	// Count should remain 75
	if reg.Count() != 75 {
		t.Errorf("expected 75 checks with nil AI client, got %d", reg.Count())
	}
}

func TestRegisterAIChecksWithUnavailableClient(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	// Register AI checks with an unconfigured client
	unconfiguredClient := &aipkg.Client{}
	reg.RegisterAIChecks(unconfiguredClient)

	// Count should remain 75 since client is not available
	if reg.Count() != 75 {
		t.Errorf("expected 75 checks with unavailable AI client, got %d", reg.Count())
	}
}

func TestGet(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	// Test getting an existing check
	check, ok := reg.Get("AND-001")
	if !ok {
		t.Error("expected to find AND-001")
	}
	if check.ID() != "AND-001" {
		t.Errorf("expected ID AND-001, got %s", check.ID())
	}

	// Test getting a non-existent check
	_, ok = reg.Get("NON-EXISTENT")
	if ok {
		t.Error("expected not to find non-existent check")
	}
}

func TestGetAll(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	checks := reg.GetAll()
	if len(checks) != 75 {
		t.Errorf("expected 75 checks, got %d", len(checks))
	}
}

func TestGetAllByID(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	checks := reg.GetAllByID()
	if len(checks) != 75 {
		t.Errorf("expected 75 checks, got %d", len(checks))
	}

	// Verify specific check exists
	if _, ok := checks["AND-001"]; !ok {
		t.Error("expected AND-001 in map")
	}
}

func TestCount(t *testing.T) {
	reg := NewRegistry()
	
	if reg.Count() != 0 {
		t.Errorf("expected 0 checks initially, got %d", reg.Count())
	}

	reg.RegisterAll()
	
	if reg.Count() != 75 {
		t.Errorf("expected 75 checks after registration, got %d", reg.Count())
	}
}

func TestGetCategories(t *testing.T) {
	reg := NewRegistry()
	categories := reg.GetCategories()

	expectedCategories := []string{
		"Android", "iOS", "Flutter", "Security", "Policy",
		"Code Quality", "Testing", "Linting", "Documentation", "Performance",
		"AI Analysis", "Reviewer",
	}

	if len(categories) != len(expectedCategories) {
		t.Errorf("expected %d categories, got %d", len(expectedCategories), len(categories))
	}

	// Check for AI Analysis category
	hasAI := false
	for _, cat := range categories {
		if cat == "AI Analysis" {
			hasAI = true
			break
		}
	}
	if !hasAI {
		t.Error("expected 'AI Analysis' in categories")
	}
}

func TestHasAIChecks(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterAll()

	// Initially should not have AI checks
	if reg.HasAIChecks() {
		t.Error("expected HasAIChecks to be false initially")
	}

	// Register AI checks
	client := createMockAIClient()
	reg.RegisterAIChecks(client)

	// Now should have AI checks
	if !reg.HasAIChecks() {
		t.Error("expected HasAIChecks to be true after registration")
	}
}

// createMockAIClient creates a mock AI client for testing
func createMockAIClient() *aipkg.Client {
	// Create a mock provider using the factory with test values
	factory := aipkg.NewProviderFactory("test-key", "", "")
	provider, _ := factory.Create("minimax")
	return aipkg.NewClient(provider, nil)
}
