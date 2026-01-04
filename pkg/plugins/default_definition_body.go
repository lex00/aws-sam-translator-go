package plugins

import (
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/openapi"
	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// DefaultDefinitionBodyPlugin handles empty DefinitionBody for API Gateway resources.
// It collects routes from function events and generates OpenAPI specifications.
type DefaultDefinitionBodyPlugin struct{}

// NewDefaultDefinitionBodyPlugin creates a new DefaultDefinitionBodyPlugin.
func NewDefaultDefinitionBodyPlugin() *DefaultDefinitionBodyPlugin {
	return &DefaultDefinitionBodyPlugin{}
}

// Name returns the plugin name.
func (p *DefaultDefinitionBodyPlugin) Name() string {
	return "DefaultDefinitionBodyPlugin"
}

// Priority returns the execution priority (500 - runs after implicit API plugins).
func (p *DefaultDefinitionBodyPlugin) Priority() int {
	return 500
}

// apiRoutes tracks routes for each API resource.
type apiRoutes struct {
	isHttpApi bool
	routes    []openapi.Route
}

// BeforeTransform handles empty DefinitionBody before transformation.
func (p *DefaultDefinitionBodyPlugin) BeforeTransform(template *types.Template) error {
	// Collect routes from function events for each API
	routesByApi := p.collectRoutes(template)

	for logicalID, resource := range template.Resources {
		if resource.Type != "AWS::Serverless::Api" && resource.Type != "AWS::Serverless::HttpApi" {
			continue
		}

		if resource.Properties == nil {
			resource.Properties = make(map[string]interface{})
			template.Resources[logicalID] = resource
		}

		// If DefinitionBody is not set and DefinitionUri is not set, add default DefinitionBody
		_, hasDefinitionBody := resource.Properties["DefinitionBody"]
		_, hasDefinitionUri := resource.Properties["DefinitionUri"]

		if !hasDefinitionBody && !hasDefinitionUri {
			isHttpApi := resource.Type == "AWS::Serverless::HttpApi"
			apiName := logicalID

			// Get routes for this API
			routes := []openapi.Route{}
			if collected, ok := routesByApi[logicalID]; ok {
				routes = collected.routes
			}

			// Also check for implicit API references
			implicitApiID := "ServerlessRestApi"
			if isHttpApi {
				implicitApiID = "ServerlessHttpApi"
			}
			if logicalID == implicitApiID {
				if collected, ok := routesByApi[""]; ok {
					routes = append(routes, collected.routes...)
				}
			}

			// Generate the OpenAPI spec
			generator := openapi.New()
			generator.Title = apiName

			var spec map[string]interface{}
			var err error

			if isHttpApi {
				// HttpApi uses OpenAPI 3.0
				spec, err = generator.GenerateOpenAPI3(routes)
			} else {
				// Api uses Swagger 2.0
				spec, err = generator.GenerateSwagger(routes)
			}

			if err != nil {
				// Fall back to minimal spec on error
				if isHttpApi {
					spec = map[string]interface{}{
						"openapi": "3.0.1",
						"info": map[string]interface{}{
							"title":   apiName,
							"version": "1.0",
						},
						"paths": map[string]interface{}{},
					}
				} else {
					spec = map[string]interface{}{
						"swagger": "2.0",
						"info": map[string]interface{}{
							"title":   apiName,
							"version": "1.0",
						},
						"paths": map[string]interface{}{},
					}
				}
			}

			resource.Properties["DefinitionBody"] = spec
			template.Resources[logicalID] = resource
		} else if hasDefinitionBody {
			// Merge routes into existing DefinitionBody
			defBody, ok := resource.Properties["DefinitionBody"].(map[string]interface{})
			if !ok {
				continue
			}

			routes := []openapi.Route{}
			if collected, ok := routesByApi[logicalID]; ok {
				routes = collected.routes
			}

			// Also check for implicit API references
			isHttpApi := resource.Type == "AWS::Serverless::HttpApi"
			implicitApiID := "ServerlessRestApi"
			if isHttpApi {
				implicitApiID = "ServerlessHttpApi"
			}
			if logicalID == implicitApiID {
				if collected, ok := routesByApi[""]; ok {
					routes = append(routes, collected.routes...)
				}
			}

			if len(routes) > 0 {
				generator := openapi.New()
				if err := generator.MergeRoutes(defBody, routes); err == nil {
					resource.Properties["DefinitionBody"] = defBody
					template.Resources[logicalID] = resource
				}
			}
		}
	}

	return nil
}

// collectRoutes extracts routes from function events.
func (p *DefaultDefinitionBodyPlugin) collectRoutes(template *types.Template) map[string]*apiRoutes {
	routesByApi := make(map[string]*apiRoutes)

	for funcLogicalID, resource := range template.Resources {
		if resource.Type != "AWS::Serverless::Function" {
			continue
		}

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
			if !ok {
				continue
			}

			if eventType != "Api" && eventType != "HttpApi" {
				continue
			}

			isHttpApi := eventType == "HttpApi"

			props, ok := eventMap["Properties"].(map[string]interface{})
			if !ok {
				props = make(map[string]interface{})
			}

			// Get the API reference
			apiRef := ""
			if isHttpApi {
				if apiID, ok := props["ApiId"]; ok {
					apiRef = p.extractRef(apiID)
				}
			} else {
				if restApiID, ok := props["RestApiId"]; ok {
					apiRef = p.extractRef(restApiID)
				}
			}

			// Get route details
			path := "/"
			if pathVal, ok := props["Path"].(string); ok {
				path = pathVal
			}

			method := "GET"
			if methodVal, ok := props["Method"].(string); ok {
				method = strings.ToUpper(methodVal)
			}

			// Build the route
			route := openapi.Route{
				Path:              path,
				Method:            method,
				FunctionLogicalID: funcLogicalID,
				// Build Lambda invocation ARN reference
				FunctionArn: map[string]interface{}{
					"Fn::Sub": "arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${" + funcLogicalID + ".Arn}/invocations",
				},
			}

			// Set payload format for HttpApi
			if isHttpApi {
				if payloadFormat, ok := props["PayloadFormatVersion"].(string); ok {
					route.PayloadFormatVersion = payloadFormat
				} else {
					route.PayloadFormatVersion = "2.0"
				}
			}

			// Check for auth settings
			if auth, ok := props["Auth"].(map[string]interface{}); ok {
				routeAuth := &openapi.RouteAuth{}
				if authorizer, ok := auth["Authorizer"].(string); ok {
					routeAuth.Authorizer = authorizer
				}
				if apiKeyRequired, ok := auth["ApiKeyRequired"].(bool); ok {
					routeAuth.ApiKeyRequired = apiKeyRequired
				}
				if scopes, ok := auth["AuthorizationScopes"].([]interface{}); ok {
					for _, s := range scopes {
						if str, ok := s.(string); ok {
							routeAuth.Scopes = append(routeAuth.Scopes, str)
						}
					}
				}
				route.Auth = routeAuth
			}

			// Add route to the appropriate API
			if _, exists := routesByApi[apiRef]; !exists {
				routesByApi[apiRef] = &apiRoutes{
					isHttpApi: isHttpApi,
					routes:    []openapi.Route{},
				}
			}
			routesByApi[apiRef].routes = append(routesByApi[apiRef].routes, route)
		}
	}

	return routesByApi
}

// extractRef extracts a logical ID from a Ref intrinsic or returns the string value.
func (p *DefaultDefinitionBodyPlugin) extractRef(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	if m, ok := val.(map[string]interface{}); ok {
		if ref, ok := m["Ref"].(string); ok {
			return ref
		}
	}
	return ""
}

// AfterTransform does nothing for this plugin.
func (p *DefaultDefinitionBodyPlugin) AfterTransform(template *types.Template) error {
	return nil
}
