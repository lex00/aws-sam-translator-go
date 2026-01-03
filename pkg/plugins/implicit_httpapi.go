package plugins

import (
	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// ImplicitHttpApiPlugin creates an implicit AWS::Serverless::HttpApi when functions have
// HttpApi events without an explicit ApiId.
type ImplicitHttpApiPlugin struct{}

// NewImplicitHttpApiPlugin creates a new ImplicitHttpApiPlugin.
func NewImplicitHttpApiPlugin() *ImplicitHttpApiPlugin {
	return &ImplicitHttpApiPlugin{}
}

// Name returns the plugin name.
func (p *ImplicitHttpApiPlugin) Name() string {
	return "ImplicitHttpApiPlugin"
}

// Priority returns the execution priority (310).
func (p *ImplicitHttpApiPlugin) Priority() int {
	return 310
}

// BeforeTransform creates the implicit ServerlessHttpApi if needed.
func (p *ImplicitHttpApiPlugin) BeforeTransform(template *types.Template) error {
	needsImplicitHttpApi := false

	// Check if any function has an HttpApi event without ApiId
	for _, resource := range template.Resources {
		if resource.Type == "AWS::Serverless::Function" {
			if resource.Properties == nil {
				continue
			}

			events, ok := resource.Properties["Events"].(map[string]interface{})
			if !ok {
				continue
			}

			for _, eventDef := range events {
				eventMap, ok := eventDef.(map[string]interface{})
				if !ok {
					continue
				}

				eventType, ok := eventMap["Type"].(string)
				if !ok || eventType != "HttpApi" {
					continue
				}

				// Check if Properties.ApiId is set
				props, ok := eventMap["Properties"].(map[string]interface{})
				if !ok {
					// No properties means implicit HttpApi
					needsImplicitHttpApi = true
					break
				}

				if _, hasApiId := props["ApiId"]; !hasApiId {
					needsImplicitHttpApi = true
					break
				}
			}

			if needsImplicitHttpApi {
				break
			}
		}
	}

	// Create ServerlessHttpApi if needed and doesn't already exist
	if needsImplicitHttpApi {
		if _, exists := template.Resources["ServerlessHttpApi"]; !exists {
			if template.Resources == nil {
				template.Resources = make(map[string]types.Resource)
			}

			template.Resources["ServerlessHttpApi"] = types.Resource{
				Type: "AWS::Serverless::HttpApi",
				Properties: map[string]interface{}{
					"StageName": "$default",
				},
			}
		}
	}

	return nil
}

// AfterTransform does nothing for this plugin.
func (p *ImplicitHttpApiPlugin) AfterTransform(template *types.Template) error {
	return nil
}
