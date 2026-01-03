// Package plugins provides the plugin system for template transformation hooks.
package plugins

import "github.com/lex00/aws-sam-translator-go/pkg/types"

// Plugin defines the interface for transformation plugins.
type Plugin interface {
	// Name returns the plugin name.
	Name() string
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
