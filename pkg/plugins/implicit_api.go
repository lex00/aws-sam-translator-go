package plugins

import (
	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// ImplicitRestApiPlugin creates an implicit AWS::Serverless::Api when functions have
// Api events without an explicit RestApiId.
type ImplicitRestApiPlugin struct{}

// NewImplicitRestApiPlugin creates a new ImplicitRestApiPlugin.
func NewImplicitRestApiPlugin() *ImplicitRestApiPlugin {
	return &ImplicitRestApiPlugin{}
}

// Name returns the plugin name.
func (p *ImplicitRestApiPlugin) Name() string {
	return "ImplicitRestApiPlugin"
}

// Priority returns the execution priority (300).
func (p *ImplicitRestApiPlugin) Priority() int {
	return 300
}

// BeforeTransform creates the implicit ServerlessRestApi if needed.
func (p *ImplicitRestApiPlugin) BeforeTransform(template *types.Template) error {
	needsImplicitApi := false

	// Check if any function has an Api event without RestApiId
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
				if !ok || eventType != "Api" {
					continue
				}

				// Check if Properties.RestApiId is set
				props, ok := eventMap["Properties"].(map[string]interface{})
				if !ok {
					// No properties means implicit API
					needsImplicitApi = true
					break
				}

				if _, hasRestApiId := props["RestApiId"]; !hasRestApiId {
					needsImplicitApi = true
					break
				}
			}

			if needsImplicitApi {
				break
			}
		}
	}

	// Create ServerlessRestApi if needed and doesn't already exist
	if needsImplicitApi {
		if _, exists := template.Resources["ServerlessRestApi"]; !exists {
			if template.Resources == nil {
				template.Resources = make(map[string]types.Resource)
			}

			template.Resources["ServerlessRestApi"] = types.Resource{
				Type: "AWS::Serverless::Api",
				Properties: map[string]interface{}{
					"StageName": "Prod",
				},
			}
		}
	}

	return nil
}

// AfterTransform does nothing for this plugin.
func (p *ImplicitRestApiPlugin) AfterTransform(template *types.Template) error {
	return nil
}
