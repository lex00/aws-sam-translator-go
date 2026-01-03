package plugins

import (
	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// DefaultDefinitionBodyPlugin handles empty DefinitionBody for API Gateway resources.
type DefaultDefinitionBodyPlugin struct{}

// NewDefaultDefinitionBodyPlugin creates a new DefaultDefinitionBodyPlugin.
func NewDefaultDefinitionBodyPlugin() *DefaultDefinitionBodyPlugin {
	return &DefaultDefinitionBodyPlugin{}
}

// Name returns the plugin name.
func (p *DefaultDefinitionBodyPlugin) Name() string {
	return "DefaultDefinitionBodyPlugin"
}

// Priority returns the execution priority (200).
func (p *DefaultDefinitionBodyPlugin) Priority() int {
	return 200
}

// BeforeTransform handles empty DefinitionBody before transformation.
func (p *DefaultDefinitionBodyPlugin) BeforeTransform(template *types.Template) error {
	for _, resource := range template.Resources {
		if resource.Type == "AWS::Serverless::Api" || resource.Type == "AWS::Serverless::HttpApi" {
			if resource.Properties == nil {
				continue
			}

			// If DefinitionBody is not set and DefinitionUri is not set, add default DefinitionBody
			_, hasDefinitionBody := resource.Properties["DefinitionBody"]
			_, hasDefinitionUri := resource.Properties["DefinitionUri"]

			if !hasDefinitionBody && !hasDefinitionUri {
				// Set a minimal default DefinitionBody
				resource.Properties["DefinitionBody"] = map[string]interface{}{
					"swagger": "2.0",
					"info": map[string]interface{}{
						"title": "API",
					},
					"paths": map[string]interface{}{},
				}
			}
		}
	}

	return nil
}

// AfterTransform does nothing for this plugin.
func (p *DefaultDefinitionBodyPlugin) AfterTransform(template *types.Template) error {
	return nil
}
