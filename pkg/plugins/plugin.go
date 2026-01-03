// Package plugins provides the plugin system for template transformation hooks.
package plugins

import (
	"sort"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// Plugin defines the interface for transformation plugins.
type Plugin interface {
	// Name returns the plugin name.
	Name() string
	// Priority returns the execution priority (lower values run first).
	Priority() int
	// BeforeTransform is called before resource transformation.
	BeforeTransform(template *types.Template) error
	// AfterTransform is called after resource transformation.
	AfterTransform(template *types.Template) error
}

// Registry manages registered plugins.
type Registry struct {
	plugins []Plugin
}

// NewRegistry creates a new plugin registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins: make([]Plugin, 0),
	}
}

// Register adds a plugin to the registry.
func (r *Registry) Register(p Plugin) {
	r.plugins = append(r.plugins, p)
}

// Plugins returns all registered plugins.
func (r *Registry) Plugins() []Plugin {
	return r.plugins
}

// RunBeforeTransform executes all plugins' BeforeTransform hooks in priority order.
func (r *Registry) RunBeforeTransform(template *types.Template) error {
	// Sort plugins by priority
	sorted := make([]Plugin, len(r.plugins))
	copy(sorted, r.plugins)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority() < sorted[j].Priority()
	})

	// Execute BeforeTransform in priority order
	for _, plugin := range sorted {
		if err := plugin.BeforeTransform(template); err != nil {
			return err
		}
	}

	return nil
}

// RunAfterTransform executes all plugins' AfterTransform hooks in priority order.
func (r *Registry) RunAfterTransform(template *types.Template) error {
	// Sort plugins by priority
	sorted := make([]Plugin, len(r.plugins))
	copy(sorted, r.plugins)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority() < sorted[j].Priority()
	})

	// Execute AfterTransform in priority order
	for _, plugin := range sorted {
		if err := plugin.AfterTransform(template); err != nil {
			return err
		}
	}

	return nil
}
