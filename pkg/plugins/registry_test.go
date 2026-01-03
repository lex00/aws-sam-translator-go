package plugins

import (
	"errors"
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// mockPlugin is a test plugin for testing priority ordering.
type mockPlugin struct {
	name           string
	priority       int
	beforeCalled   *[]string
	afterCalled    *[]string
	beforeErr      error
	afterErr       error
}

func (p *mockPlugin) Name() string {
	return p.name
}

func (p *mockPlugin) Priority() int {
	return p.priority
}

func (p *mockPlugin) BeforeTransform(template *types.Template) error {
	if p.beforeCalled != nil {
		*p.beforeCalled = append(*p.beforeCalled, p.name)
	}
	return p.beforeErr
}

func (p *mockPlugin) AfterTransform(template *types.Template) error {
	if p.afterCalled != nil {
		*p.afterCalled = append(*p.afterCalled, p.name)
	}
	return p.afterErr
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	plugin1 := &mockPlugin{name: "plugin1", priority: 100}
	plugin2 := &mockPlugin{name: "plugin2", priority: 200}

	registry.Register(plugin1)
	registry.Register(plugin2)

	plugins := registry.Plugins()
	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(plugins))
	}
}

func TestRegistry_RunBeforeTransform_PriorityOrder(t *testing.T) {
	registry := NewRegistry()

	var callOrder []string

	// Register plugins in reverse priority order
	plugin3 := &mockPlugin{name: "plugin3", priority: 300, beforeCalled: &callOrder}
	plugin1 := &mockPlugin{name: "plugin1", priority: 100, beforeCalled: &callOrder}
	plugin2 := &mockPlugin{name: "plugin2", priority: 200, beforeCalled: &callOrder}

	registry.Register(plugin3)
	registry.Register(plugin1)
	registry.Register(plugin2)

	template := &types.Template{}
	err := registry.RunBeforeTransform(template)
	if err != nil {
		t.Fatalf("RunBeforeTransform failed: %v", err)
	}

	// Verify plugins were called in priority order (100, 200, 300)
	expectedOrder := []string{"plugin1", "plugin2", "plugin3"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d calls, got %d", len(expectedOrder), len(callOrder))
	}

	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("Expected plugin %s at position %d, got %s", expected, i, callOrder[i])
		}
	}
}

func TestRegistry_RunAfterTransform_PriorityOrder(t *testing.T) {
	registry := NewRegistry()

	var callOrder []string

	// Register plugins in reverse priority order
	plugin3 := &mockPlugin{name: "plugin3", priority: 300, afterCalled: &callOrder}
	plugin1 := &mockPlugin{name: "plugin1", priority: 100, afterCalled: &callOrder}
	plugin2 := &mockPlugin{name: "plugin2", priority: 200, afterCalled: &callOrder}

	registry.Register(plugin3)
	registry.Register(plugin1)
	registry.Register(plugin2)

	template := &types.Template{}
	err := registry.RunAfterTransform(template)
	if err != nil {
		t.Fatalf("RunAfterTransform failed: %v", err)
	}

	// Verify plugins were called in priority order (100, 200, 300)
	expectedOrder := []string{"plugin1", "plugin2", "plugin3"}
	if len(callOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d calls, got %d", len(expectedOrder), len(callOrder))
	}

	for i, expected := range expectedOrder {
		if callOrder[i] != expected {
			t.Errorf("Expected plugin %s at position %d, got %s", expected, i, callOrder[i])
		}
	}
}

func TestRegistry_RunBeforeTransform_Error(t *testing.T) {
	registry := NewRegistry()

	expectedErr := errors.New("before transform error")
	plugin := &mockPlugin{name: "plugin1", priority: 100, beforeErr: expectedErr}

	registry.Register(plugin)

	template := &types.Template{}
	err := registry.RunBeforeTransform(template)
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestRegistry_RunAfterTransform_Error(t *testing.T) {
	registry := NewRegistry()

	expectedErr := errors.New("after transform error")
	plugin := &mockPlugin{name: "plugin1", priority: 100, afterErr: expectedErr}

	registry.Register(plugin)

	template := &types.Template{}
	err := registry.RunAfterTransform(template)
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestRegistry_EmptyRegistry(t *testing.T) {
	registry := NewRegistry()
	template := &types.Template{}

	// Should not error with no plugins
	err := registry.RunBeforeTransform(template)
	if err != nil {
		t.Errorf("RunBeforeTransform with empty registry failed: %v", err)
	}

	err = registry.RunAfterTransform(template)
	if err != nil {
		t.Errorf("RunAfterTransform with empty registry failed: %v", err)
	}
}
