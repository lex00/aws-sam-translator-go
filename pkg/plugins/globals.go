package plugins

import (
	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// GlobalsPlugin applies the Globals section to all SAM resources before transformation.
type GlobalsPlugin struct{}

// NewGlobalsPlugin creates a new GlobalsPlugin.
func NewGlobalsPlugin() *GlobalsPlugin {
	return &GlobalsPlugin{}
}

// Name returns the plugin name.
func (p *GlobalsPlugin) Name() string {
	return "GlobalsPlugin"
}

// Priority returns the execution priority (100 - runs first).
func (p *GlobalsPlugin) Priority() int {
	return 100
}

// BeforeTransform applies global properties to all SAM resources.
func (p *GlobalsPlugin) BeforeTransform(template *types.Template) error {
	if template.Globals == nil {
		return nil
	}

	// Apply Function globals
	if functionGlobals, ok := template.Globals["Function"].(map[string]interface{}); ok {
		p.applyFunctionGlobals(template, functionGlobals)
	}

	// Apply Api globals
	if apiGlobals, ok := template.Globals["Api"].(map[string]interface{}); ok {
		p.applyApiGlobals(template, apiGlobals)
	}

	// Apply HttpApi globals
	if httpApiGlobals, ok := template.Globals["HttpApi"].(map[string]interface{}); ok {
		p.applyHttpApiGlobals(template, httpApiGlobals)
	}

	// Apply SimpleTable globals
	if simpleTableGlobals, ok := template.Globals["SimpleTable"].(map[string]interface{}); ok {
		p.applySimpleTableGlobals(template, simpleTableGlobals)
	}

	return nil
}

// AfterTransform does nothing for this plugin.
func (p *GlobalsPlugin) AfterTransform(template *types.Template) error {
	return nil
}

// applyFunctionGlobals applies global Function properties to all AWS::Serverless::Function resources.
func (p *GlobalsPlugin) applyFunctionGlobals(template *types.Template, globals map[string]interface{}) {
	for _, resource := range template.Resources {
		if resource.Type == "AWS::Serverless::Function" {
			if resource.Properties == nil {
				resource.Properties = make(map[string]interface{})
			}
			mergeProperties(resource.Properties, globals)
		}
	}
}

// applyApiGlobals applies global Api properties to all AWS::Serverless::Api resources.
func (p *GlobalsPlugin) applyApiGlobals(template *types.Template, globals map[string]interface{}) {
	for _, resource := range template.Resources {
		if resource.Type == "AWS::Serverless::Api" {
			if resource.Properties == nil {
				resource.Properties = make(map[string]interface{})
			}
			mergeProperties(resource.Properties, globals)
		}
	}
}

// applyHttpApiGlobals applies global HttpApi properties to all AWS::Serverless::HttpApi resources.
func (p *GlobalsPlugin) applyHttpApiGlobals(template *types.Template, globals map[string]interface{}) {
	for _, resource := range template.Resources {
		if resource.Type == "AWS::Serverless::HttpApi" {
			if resource.Properties == nil {
				resource.Properties = make(map[string]interface{})
			}
			mergeProperties(resource.Properties, globals)
		}
	}
}

// applySimpleTableGlobals applies global SimpleTable properties to all AWS::Serverless::SimpleTable resources.
func (p *GlobalsPlugin) applySimpleTableGlobals(template *types.Template, globals map[string]interface{}) {
	for _, resource := range template.Resources {
		if resource.Type == "AWS::Serverless::SimpleTable" {
			if resource.Properties == nil {
				resource.Properties = make(map[string]interface{})
			}
			mergeProperties(resource.Properties, globals)
		}
	}
}

// mergeProperties merges global properties into resource properties.
// Resource-specific properties take precedence over global properties.
func mergeProperties(resourceProps, globalProps map[string]interface{}) {
	for key, value := range globalProps {
		// Only set if not already defined in resource
		if _, exists := resourceProps[key]; !exists {
			resourceProps[key] = deepCopy(value)
		}
	}
}

// deepCopy creates a deep copy of a value.
func deepCopy(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			result[key] = deepCopy(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = deepCopy(item)
		}
		return result
	default:
		return v
	}
}
