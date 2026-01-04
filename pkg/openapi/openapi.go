// Package openapi provides OpenAPI/Swagger generation for API Gateway.
package openapi

import (
	"fmt"
	"sort"
	"strings"
)

// Generator generates OpenAPI specifications for API Gateway.
type Generator struct {
	// Title is the API title used in the info section.
	Title string

	// Version is the API version used in the info section.
	Version string
}

// New creates a new OpenAPI Generator with default settings.
func New() *Generator {
	return &Generator{
		Title:   "API",
		Version: "1.0",
	}
}

// NewWithOptions creates a new OpenAPI Generator with custom settings.
func NewWithOptions(title, version string) *Generator {
	g := New()
	if title != "" {
		g.Title = title
	}
	if version != "" {
		g.Version = version
	}
	return g
}

// Route represents an API route with Lambda integration.
type Route struct {
	// Path is the API path (e.g., "/users/{id}").
	Path string

	// Method is the HTTP method (GET, POST, PUT, DELETE, etc.).
	Method string

	// FunctionLogicalID is the CloudFormation logical ID of the Lambda function.
	FunctionLogicalID string

	// FunctionArn is the Lambda function ARN (can be an intrinsic function).
	FunctionArn interface{}

	// OperationID is an optional unique operation identifier.
	OperationID string

	// Summary is an optional short summary of the operation.
	Summary string

	// Description is an optional longer description of the operation.
	Description string

	// Consumes is a list of MIME types the operation can consume.
	Consumes []string

	// Produces is a list of MIME types the operation can produce.
	Produces []string

	// Auth contains optional authorization configuration.
	Auth *RouteAuth

	// RequestParameters contains request parameter definitions.
	RequestParameters map[string]RequestParameter

	// PayloadFormatVersion is the payload format version (1.0 or 2.0).
	PayloadFormatVersion string
}

// RouteAuth contains authorization configuration for a route.
type RouteAuth struct {
	// Authorizer is the name of the authorizer to use.
	Authorizer string

	// Scopes is a list of OAuth2 scopes required.
	Scopes []string

	// ApiKeyRequired indicates if an API key is required.
	ApiKeyRequired bool
}

// RequestParameter represents a request parameter.
type RequestParameter struct {
	// Required indicates if the parameter is required.
	Required bool

	// Caching indicates if the parameter should be used for caching.
	Caching bool
}

// GenerateSwagger generates a Swagger 2.0 specification.
func (g *Generator) GenerateSwagger(routes []Route) (map[string]interface{}, error) {
	spec := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":   g.Title,
			"version": g.Version,
		},
		"paths": make(map[string]interface{}),
	}

	paths := spec["paths"].(map[string]interface{})

	for _, route := range routes {
		if err := g.addSwaggerRoute(paths, route); err != nil {
			return nil, fmt.Errorf("failed to add route %s %s: %w", route.Method, route.Path, err)
		}
	}

	return spec, nil
}

// GenerateOpenAPI3 generates an OpenAPI 3.0 specification.
func (g *Generator) GenerateOpenAPI3(routes []Route) (map[string]interface{}, error) {
	spec := map[string]interface{}{
		"openapi": "3.0.1",
		"info": map[string]interface{}{
			"title":   g.Title,
			"version": g.Version,
		},
		"paths": make(map[string]interface{}),
	}

	paths := spec["paths"].(map[string]interface{})

	for _, route := range routes {
		if err := g.addOpenAPI3Route(paths, route); err != nil {
			return nil, fmt.Errorf("failed to add route %s %s: %w", route.Method, route.Path, err)
		}
	}

	return spec, nil
}

// MergeRoutes merges routes into an existing OpenAPI specification.
// It detects whether the spec is Swagger 2.0 or OpenAPI 3.0 and uses the appropriate format.
func (g *Generator) MergeRoutes(spec map[string]interface{}, routes []Route) error {
	if spec == nil {
		return fmt.Errorf("spec cannot be nil")
	}

	// Detect spec version
	isOpenAPI3 := false
	if _, ok := spec["openapi"]; ok {
		isOpenAPI3 = true
	}

	// Ensure paths exists
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		paths = make(map[string]interface{})
		spec["paths"] = paths
	}

	for _, route := range routes {
		var err error
		if isOpenAPI3 {
			err = g.addOpenAPI3Route(paths, route)
		} else {
			err = g.addSwaggerRoute(paths, route)
		}
		if err != nil {
			return fmt.Errorf("failed to merge route %s %s: %w", route.Method, route.Path, err)
		}
	}

	return nil
}

// addSwaggerRoute adds a route to a Swagger 2.0 paths object.
func (g *Generator) addSwaggerRoute(paths map[string]interface{}, route Route) error {
	method := strings.ToLower(route.Method)
	if method == "" {
		return fmt.Errorf("method is required")
	}
	if route.Path == "" {
		return fmt.Errorf("path is required")
	}

	// Get or create path item
	pathItem, ok := paths[route.Path].(map[string]interface{})
	if !ok {
		pathItem = make(map[string]interface{})
		paths[route.Path] = pathItem
	}

	// Build operation
	operation := make(map[string]interface{})

	// Add operation metadata
	if route.OperationID != "" {
		operation["operationId"] = route.OperationID
	}
	if route.Summary != "" {
		operation["summary"] = route.Summary
	}
	if route.Description != "" {
		operation["description"] = route.Description
	}
	if len(route.Consumes) > 0 {
		operation["consumes"] = route.Consumes
	}
	if len(route.Produces) > 0 {
		operation["produces"] = route.Produces
	}

	// Add parameters from path
	params := g.extractPathParameters(route.Path)
	if len(params) > 0 {
		operation["parameters"] = params
	}

	// Add default responses
	operation["responses"] = map[string]interface{}{
		"200": map[string]interface{}{
			"description": "Success",
		},
	}

	// Add x-amazon-apigateway-integration
	operation["x-amazon-apigateway-integration"] = g.buildSwaggerIntegration(route)

	// Add security if auth is configured
	if route.Auth != nil {
		if route.Auth.Authorizer != "" {
			security := map[string]interface{}{
				route.Auth.Authorizer: route.Auth.Scopes,
			}
			if route.Auth.Scopes == nil {
				security[route.Auth.Authorizer] = []interface{}{}
			}
			operation["security"] = []interface{}{security}
		}
		if route.Auth.ApiKeyRequired {
			// Add API key security
			if existing, ok := operation["security"].([]interface{}); ok {
				existing = append(existing, map[string]interface{}{
					"api_key": []interface{}{},
				})
				operation["security"] = existing
			} else {
				operation["security"] = []interface{}{
					map[string]interface{}{
						"api_key": []interface{}{},
					},
				}
			}
		}
	}

	pathItem[method] = operation
	return nil
}

// addOpenAPI3Route adds a route to an OpenAPI 3.0 paths object.
func (g *Generator) addOpenAPI3Route(paths map[string]interface{}, route Route) error {
	method := strings.ToLower(route.Method)
	if method == "" {
		return fmt.Errorf("method is required")
	}
	if route.Path == "" {
		return fmt.Errorf("path is required")
	}

	// Get or create path item
	pathItem, ok := paths[route.Path].(map[string]interface{})
	if !ok {
		pathItem = make(map[string]interface{})
		paths[route.Path] = pathItem
	}

	// Build operation
	operation := make(map[string]interface{})

	// Add operation metadata
	if route.OperationID != "" {
		operation["operationId"] = route.OperationID
	}
	if route.Summary != "" {
		operation["summary"] = route.Summary
	}
	if route.Description != "" {
		operation["description"] = route.Description
	}

	// Add parameters from path
	params := g.extractPathParametersOpenAPI3(route.Path)
	if len(params) > 0 {
		operation["parameters"] = params
	}

	// Add default responses
	operation["responses"] = map[string]interface{}{
		"200": map[string]interface{}{
			"description": "Success",
		},
	}

	// Add x-amazon-apigateway-integration
	operation["x-amazon-apigateway-integration"] = g.buildOpenAPI3Integration(route)

	// Add security if auth is configured
	if route.Auth != nil {
		if route.Auth.Authorizer != "" {
			security := map[string]interface{}{
				route.Auth.Authorizer: route.Auth.Scopes,
			}
			if route.Auth.Scopes == nil {
				security[route.Auth.Authorizer] = []interface{}{}
			}
			operation["security"] = []interface{}{security}
		}
	}

	pathItem[method] = operation
	return nil
}

// buildSwaggerIntegration builds the x-amazon-apigateway-integration for Swagger 2.0.
func (g *Generator) buildSwaggerIntegration(route Route) map[string]interface{} {
	integration := map[string]interface{}{
		"type":       "aws_proxy",
		"httpMethod": "POST",
		"uri":        g.buildLambdaIntegrationUri(route),
	}

	// Add passthroughBehavior for proxy integrations
	integration["passthroughBehavior"] = "when_no_match"

	return integration
}

// buildOpenAPI3Integration builds the x-amazon-apigateway-integration for OpenAPI 3.0.
func (g *Generator) buildOpenAPI3Integration(route Route) map[string]interface{} {
	integration := map[string]interface{}{
		"type":       "aws_proxy",
		"httpMethod": "POST",
		"uri":        g.buildLambdaIntegrationUri(route),
	}

	// Set payload format version for HTTP APIs
	payloadVersion := route.PayloadFormatVersion
	if payloadVersion == "" {
		payloadVersion = "2.0"
	}
	integration["payloadFormatVersion"] = payloadVersion

	return integration
}

// buildLambdaIntegrationUri builds the Lambda integration URI.
func (g *Generator) buildLambdaIntegrationUri(route Route) interface{} {
	// If FunctionArn is provided directly, use it
	if route.FunctionArn != nil {
		return map[string]interface{}{
			"Fn::Sub": []interface{}{
				"arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${FunctionArn}/invocations",
				map[string]interface{}{
					"FunctionArn": route.FunctionArn,
				},
			},
		}
	}

	// Otherwise use the logical ID to reference the function
	if route.FunctionLogicalID != "" {
		return map[string]interface{}{
			"Fn::Sub": []interface{}{
				"arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${FunctionArn}/invocations",
				map[string]interface{}{
					"FunctionArn": map[string]interface{}{
						"Fn::GetAtt": []interface{}{route.FunctionLogicalID, "Arn"},
					},
				},
			},
		}
	}

	// Fallback - should not happen in normal usage
	return ""
}

// extractPathParameters extracts path parameters for Swagger 2.0 format.
func (g *Generator) extractPathParameters(path string) []map[string]interface{} {
	var params []map[string]interface{}

	// Find all {param} patterns
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName := strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
			// Handle proxy parameters like {proxy+}
			paramName = strings.TrimSuffix(paramName, "+")

			params = append(params, map[string]interface{}{
				"name":     paramName,
				"in":       "path",
				"required": true,
				"type":     "string",
			})
		}
	}

	return params
}

// extractPathParametersOpenAPI3 extracts path parameters for OpenAPI 3.0 format.
func (g *Generator) extractPathParametersOpenAPI3(path string) []map[string]interface{} {
	var params []map[string]interface{}

	// Find all {param} patterns
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName := strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
			// Handle proxy parameters like {proxy+}
			paramName = strings.TrimSuffix(paramName, "+")

			params = append(params, map[string]interface{}{
				"name":     paramName,
				"in":       "path",
				"required": true,
				"schema": map[string]interface{}{
					"type": "string",
				},
			})
		}
	}

	return params
}

// AddCorsToSpec adds CORS configuration to an OpenAPI spec.
func (g *Generator) AddCorsToSpec(spec map[string]interface{}, cors *CorsConfig) error {
	if spec == nil || cors == nil {
		return nil
	}

	// For OpenAPI 3.0, use x-amazon-apigateway-cors
	if _, ok := spec["openapi"]; ok {
		corsConfig := make(map[string]interface{})

		if cors.AllowOrigins != nil {
			corsConfig["allowOrigins"] = cors.AllowOrigins
		}
		if cors.AllowMethods != nil {
			corsConfig["allowMethods"] = cors.AllowMethods
		}
		if cors.AllowHeaders != nil {
			corsConfig["allowHeaders"] = cors.AllowHeaders
		}
		if cors.ExposeHeaders != nil {
			corsConfig["exposeHeaders"] = cors.ExposeHeaders
		}
		if cors.MaxAge > 0 {
			corsConfig["maxAge"] = cors.MaxAge
		}
		if cors.AllowCredentials {
			corsConfig["allowCredentials"] = true
		}

		spec["x-amazon-apigateway-cors"] = corsConfig
		return nil
	}

	// For Swagger 2.0, add OPTIONS methods to each path
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		return nil
	}

	// Get sorted paths for deterministic output
	pathKeys := make([]string, 0, len(paths))
	for path := range paths {
		pathKeys = append(pathKeys, path)
	}
	sort.Strings(pathKeys)

	for _, pathKey := range pathKeys {
		pathItem, ok := paths[pathKey].(map[string]interface{})
		if !ok {
			continue
		}

		// Add OPTIONS method for CORS preflight
		if _, hasOptions := pathItem["options"]; !hasOptions {
			pathItem["options"] = g.buildCorsOptionsMethod(cors)
		}
	}

	return nil
}

// CorsConfig represents CORS configuration for OpenAPI.
type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           int
	AllowCredentials bool
}

// buildCorsOptionsMethod builds an OPTIONS method for CORS preflight.
func (g *Generator) buildCorsOptionsMethod(cors *CorsConfig) map[string]interface{} {
	allowOrigin := "*"
	if len(cors.AllowOrigins) > 0 {
		allowOrigin = strings.Join(cors.AllowOrigins, ",")
	}

	allowMethods := "GET,POST,PUT,DELETE,OPTIONS"
	if len(cors.AllowMethods) > 0 {
		allowMethods = strings.Join(cors.AllowMethods, ",")
	}

	allowHeaders := "Content-Type,Authorization,X-Amz-Date,X-Api-Key,X-Amz-Security-Token"
	if len(cors.AllowHeaders) > 0 {
		allowHeaders = strings.Join(cors.AllowHeaders, ",")
	}

	return map[string]interface{}{
		"summary": "CORS support",
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"description": "Default response for CORS method",
				"headers": map[string]interface{}{
					"Access-Control-Allow-Origin": map[string]interface{}{
						"type": "string",
					},
					"Access-Control-Allow-Methods": map[string]interface{}{
						"type": "string",
					},
					"Access-Control-Allow-Headers": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
		"x-amazon-apigateway-integration": map[string]interface{}{
			"type": "mock",
			"requestTemplates": map[string]interface{}{
				"application/json": `{"statusCode": 200}`,
			},
			"responses": map[string]interface{}{
				"default": map[string]interface{}{
					"statusCode": "200",
					"responseParameters": map[string]interface{}{
						"method.response.header.Access-Control-Allow-Headers": fmt.Sprintf("'%s'", allowHeaders),
						"method.response.header.Access-Control-Allow-Methods": fmt.Sprintf("'%s'", allowMethods),
						"method.response.header.Access-Control-Allow-Origin":  fmt.Sprintf("'%s'", allowOrigin),
					},
				},
			},
		},
	}
}

// AddSecurityDefinitions adds security definitions to the spec.
func (g *Generator) AddSecurityDefinitions(spec map[string]interface{}, authorizers map[string]interface{}) error {
	if spec == nil || len(authorizers) == 0 {
		return nil
	}

	// Detect spec version
	isOpenAPI3 := false
	if _, ok := spec["openapi"]; ok {
		isOpenAPI3 = true
	}

	if isOpenAPI3 {
		// OpenAPI 3.0 uses components.securitySchemes
		components, ok := spec["components"].(map[string]interface{})
		if !ok {
			components = make(map[string]interface{})
			spec["components"] = components
		}

		securitySchemes, ok := components["securitySchemes"].(map[string]interface{})
		if !ok {
			securitySchemes = make(map[string]interface{})
			components["securitySchemes"] = securitySchemes
		}

		for name, config := range authorizers {
			securitySchemes[name] = g.buildSecuritySchemeOpenAPI3(config)
		}
	} else {
		// Swagger 2.0 uses securityDefinitions
		securityDefs, ok := spec["securityDefinitions"].(map[string]interface{})
		if !ok {
			securityDefs = make(map[string]interface{})
			spec["securityDefinitions"] = securityDefs
		}

		for name, config := range authorizers {
			securityDefs[name] = g.buildSecuritySchemeSwagger(config)
		}
	}

	return nil
}

// buildSecuritySchemeSwagger builds a security scheme for Swagger 2.0.
func (g *Generator) buildSecuritySchemeSwagger(config interface{}) map[string]interface{} {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"type": "apiKey",
			"name": "Authorization",
			"in":   "header",
		}
	}

	scheme := make(map[string]interface{})

	// Check for Cognito authorizer
	if userPoolArn, ok := configMap["UserPoolArn"]; ok {
		scheme["type"] = "apiKey"
		scheme["name"] = "Authorization"
		scheme["in"] = "header"
		scheme["x-amazon-apigateway-authtype"] = "cognito_user_pools"
		scheme["x-amazon-apigateway-authorizer"] = map[string]interface{}{
			"type":         "cognito_user_pools",
			"providerARNs": []interface{}{userPoolArn},
		}
		return scheme
	}

	// Check for Lambda authorizer
	if functionArn, ok := configMap["FunctionArn"]; ok {
		scheme["type"] = "apiKey"
		scheme["name"] = "Authorization"
		scheme["in"] = "header"
		scheme["x-amazon-apigateway-authtype"] = "custom"
		scheme["x-amazon-apigateway-authorizer"] = map[string]interface{}{
			"type":                         "token",
			"authorizerUri":                g.buildAuthorizerUri(functionArn),
			"authorizerResultTtlInSeconds": 300,
		}
		return scheme
	}

	// Default API key scheme
	scheme["type"] = "apiKey"
	scheme["name"] = "x-api-key"
	scheme["in"] = "header"

	return scheme
}

// buildSecuritySchemeOpenAPI3 builds a security scheme for OpenAPI 3.0.
func (g *Generator) buildSecuritySchemeOpenAPI3(config interface{}) map[string]interface{} {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"type": "apiKey",
			"name": "Authorization",
			"in":   "header",
		}
	}

	scheme := make(map[string]interface{})

	// Check for JWT authorizer (HttpApi)
	if jwtConfig, ok := configMap["JwtConfiguration"]; ok {
		scheme["type"] = "oauth2"
		scheme["x-amazon-apigateway-authorizer"] = map[string]interface{}{
			"type":             "jwt",
			"jwtConfiguration": jwtConfig,
		}
		if identitySource, ok := configMap["IdentitySource"]; ok {
			scheme["x-amazon-apigateway-authorizer"].(map[string]interface{})["identitySource"] = identitySource
		}
		return scheme
	}

	// Check for Cognito authorizer
	if userPoolArn, ok := configMap["UserPoolArn"]; ok {
		scheme["type"] = "apiKey"
		scheme["name"] = "Authorization"
		scheme["in"] = "header"
		scheme["x-amazon-apigateway-authorizer"] = map[string]interface{}{
			"type":         "cognito_user_pools",
			"providerARNs": []interface{}{userPoolArn},
		}
		return scheme
	}

	// Check for Lambda authorizer
	if functionArn, ok := configMap["FunctionArn"]; ok {
		authType := "request"
		if _, hasToken := configMap["Identity"]; !hasToken {
			authType = "token"
		}
		scheme["type"] = "apiKey"
		scheme["name"] = "Authorization"
		scheme["in"] = "header"
		scheme["x-amazon-apigateway-authorizer"] = map[string]interface{}{
			"type":          authType,
			"authorizerUri": g.buildAuthorizerUri(functionArn),
		}
		return scheme
	}

	// Default API key scheme
	scheme["type"] = "apiKey"
	scheme["name"] = "x-api-key"
	scheme["in"] = "header"

	return scheme
}

// buildAuthorizerUri builds the Lambda authorizer URI.
func (g *Generator) buildAuthorizerUri(functionArn interface{}) interface{} {
	return map[string]interface{}{
		"Fn::Sub": []interface{}{
			"arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${FunctionArn}/invocations",
			map[string]interface{}{
				"FunctionArn": functionArn,
			},
		},
	}
}

// IsOpenAPI3 checks if the spec is OpenAPI 3.0.
func IsOpenAPI3(spec map[string]interface{}) bool {
	_, ok := spec["openapi"]
	return ok
}

// IsSwagger checks if the spec is Swagger 2.0.
func IsSwagger(spec map[string]interface{}) bool {
	_, ok := spec["swagger"]
	return ok
}
