package openapi

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	g := New()
	if g.Title != "API" {
		t.Errorf("expected Title 'API', got %q", g.Title)
	}
	if g.Version != "1.0" {
		t.Errorf("expected Version '1.0', got %q", g.Version)
	}
}

func TestNewWithOptions(t *testing.T) {
	g := NewWithOptions("My API", "2.0")
	if g.Title != "My API" {
		t.Errorf("expected Title 'My API', got %q", g.Title)
	}
	if g.Version != "2.0" {
		t.Errorf("expected Version '2.0', got %q", g.Version)
	}
}

func TestGenerateSwagger(t *testing.T) {
	g := New()

	routes := []Route{
		{
			Path:              "/users",
			Method:            "GET",
			FunctionLogicalID: "ListUsersFunction",
		},
		{
			Path:              "/users/{id}",
			Method:            "GET",
			FunctionLogicalID: "GetUserFunction",
		},
		{
			Path:              "/users",
			Method:            "POST",
			FunctionLogicalID: "CreateUserFunction",
		},
	}

	spec, err := g.GenerateSwagger(routes)
	if err != nil {
		t.Fatalf("GenerateSwagger failed: %v", err)
	}

	// Check swagger version
	if spec["swagger"] != "2.0" {
		t.Errorf("expected swagger '2.0', got %v", spec["swagger"])
	}

	// Check info
	info, ok := spec["info"].(map[string]interface{})
	if !ok {
		t.Fatal("expected info to be a map")
	}
	if info["title"] != "API" {
		t.Errorf("expected title 'API', got %v", info["title"])
	}

	// Check paths
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		t.Fatal("expected paths to be a map")
	}

	// Check /users path
	usersPath, ok := paths["/users"].(map[string]interface{})
	if !ok {
		t.Fatal("expected /users path")
	}

	// Check GET method
	getMethod, ok := usersPath["get"].(map[string]interface{})
	if !ok {
		t.Fatal("expected GET method on /users")
	}

	// Check integration
	integration, ok := getMethod["x-amazon-apigateway-integration"].(map[string]interface{})
	if !ok {
		t.Fatal("expected x-amazon-apigateway-integration")
	}
	if integration["type"] != "aws_proxy" {
		t.Errorf("expected type 'aws_proxy', got %v", integration["type"])
	}
	if integration["httpMethod"] != "POST" {
		t.Errorf("expected httpMethod 'POST', got %v", integration["httpMethod"])
	}

	// Check POST method exists
	if _, ok := usersPath["post"]; !ok {
		t.Error("expected POST method on /users")
	}

	// Check /users/{id} path
	userByIdPath, ok := paths["/users/{id}"].(map[string]interface{})
	if !ok {
		t.Fatal("expected /users/{id} path")
	}

	getByIdMethod, ok := userByIdPath["get"].(map[string]interface{})
	if !ok {
		t.Fatal("expected GET method on /users/{id}")
	}

	// Check path parameters
	params, ok := getByIdMethod["parameters"].([]map[string]interface{})
	if !ok || len(params) != 1 {
		t.Fatal("expected 1 path parameter")
	}
	if params[0]["name"] != "id" {
		t.Errorf("expected parameter name 'id', got %v", params[0]["name"])
	}
	if params[0]["in"] != "path" {
		t.Errorf("expected parameter in 'path', got %v", params[0]["in"])
	}
}

func TestGenerateOpenAPI3(t *testing.T) {
	g := NewWithOptions("My HTTP API", "1.0")

	routes := []Route{
		{
			Path:              "/items",
			Method:            "GET",
			FunctionLogicalID: "ListItemsFunction",
		},
		{
			Path:                 "/items/{id}",
			Method:               "GET",
			FunctionLogicalID:    "GetItemFunction",
			PayloadFormatVersion: "1.0",
		},
	}

	spec, err := g.GenerateOpenAPI3(routes)
	if err != nil {
		t.Fatalf("GenerateOpenAPI3 failed: %v", err)
	}

	// Check openapi version
	if spec["openapi"] != "3.0.1" {
		t.Errorf("expected openapi '3.0.1', got %v", spec["openapi"])
	}

	// Check info
	info := spec["info"].(map[string]interface{})
	if info["title"] != "My HTTP API" {
		t.Errorf("expected title 'My HTTP API', got %v", info["title"])
	}

	// Check paths
	paths := spec["paths"].(map[string]interface{})

	// Check /items path
	itemsPath := paths["/items"].(map[string]interface{})
	getMethod := itemsPath["get"].(map[string]interface{})

	// Check integration with default payload version
	integration := getMethod["x-amazon-apigateway-integration"].(map[string]interface{})
	if integration["payloadFormatVersion"] != "2.0" {
		t.Errorf("expected payloadFormatVersion '2.0', got %v", integration["payloadFormatVersion"])
	}

	// Check /items/{id} with custom payload version
	itemByIdPath := paths["/items/{id}"].(map[string]interface{})
	getByIdMethod := itemByIdPath["get"].(map[string]interface{})
	integrationById := getByIdMethod["x-amazon-apigateway-integration"].(map[string]interface{})
	if integrationById["payloadFormatVersion"] != "1.0" {
		t.Errorf("expected payloadFormatVersion '1.0', got %v", integrationById["payloadFormatVersion"])
	}

	// Check OpenAPI 3 parameter format
	params := getByIdMethod["parameters"].([]map[string]interface{})
	if len(params) != 1 {
		t.Fatal("expected 1 path parameter")
	}
	schema, ok := params[0]["schema"].(map[string]interface{})
	if !ok {
		t.Fatal("expected schema in parameter (OpenAPI 3 format)")
	}
	if schema["type"] != "string" {
		t.Errorf("expected schema type 'string', got %v", schema["type"])
	}
}

func TestMergeRoutes(t *testing.T) {
	g := New()

	// Start with an existing Swagger spec
	existingSpec := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":   "Existing API",
			"version": "1.0",
		},
		"paths": map[string]interface{}{
			"/existing": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Existing endpoint",
				},
			},
		},
	}

	routes := []Route{
		{
			Path:              "/new",
			Method:            "POST",
			FunctionLogicalID: "NewFunction",
		},
	}

	err := g.MergeRoutes(existingSpec, routes)
	if err != nil {
		t.Fatalf("MergeRoutes failed: %v", err)
	}

	paths := existingSpec["paths"].(map[string]interface{})

	// Check existing path is preserved
	if _, ok := paths["/existing"]; !ok {
		t.Error("expected /existing path to be preserved")
	}

	// Check new path was added
	newPath, ok := paths["/new"].(map[string]interface{})
	if !ok {
		t.Fatal("expected /new path to be added")
	}
	if _, ok := newPath["post"]; !ok {
		t.Error("expected POST method on /new")
	}
}

func TestMergeRoutesOpenAPI3(t *testing.T) {
	g := New()

	// Start with an existing OpenAPI 3 spec
	existingSpec := map[string]interface{}{
		"openapi": "3.0.1",
		"info": map[string]interface{}{
			"title":   "Existing API",
			"version": "1.0",
		},
		"paths": map[string]interface{}{},
	}

	routes := []Route{
		{
			Path:              "/items",
			Method:            "GET",
			FunctionLogicalID: "ListItemsFunction",
		},
	}

	err := g.MergeRoutes(existingSpec, routes)
	if err != nil {
		t.Fatalf("MergeRoutes failed: %v", err)
	}

	paths := existingSpec["paths"].(map[string]interface{})
	itemsPath := paths["/items"].(map[string]interface{})
	getMethod := itemsPath["get"].(map[string]interface{})

	// Should use OpenAPI 3 format (payloadFormatVersion)
	integration := getMethod["x-amazon-apigateway-integration"].(map[string]interface{})
	if _, ok := integration["payloadFormatVersion"]; !ok {
		t.Error("expected payloadFormatVersion for OpenAPI 3")
	}
}

func TestRouteWithAuth(t *testing.T) {
	g := New()

	routes := []Route{
		{
			Path:              "/secure",
			Method:            "GET",
			FunctionLogicalID: "SecureFunction",
			Auth: &RouteAuth{
				Authorizer: "MyCognitoAuthorizer",
				Scopes:     []string{"read", "write"},
			},
		},
	}

	spec, err := g.GenerateSwagger(routes)
	if err != nil {
		t.Fatalf("GenerateSwagger failed: %v", err)
	}

	paths := spec["paths"].(map[string]interface{})
	securePath := paths["/secure"].(map[string]interface{})
	getMethod := securePath["get"].(map[string]interface{})

	security, ok := getMethod["security"].([]interface{})
	if !ok || len(security) == 0 {
		t.Fatal("expected security configuration")
	}

	securityItem := security[0].(map[string]interface{})
	scopes, ok := securityItem["MyCognitoAuthorizer"].([]string)
	if !ok {
		t.Fatal("expected MyCognitoAuthorizer security")
	}
	if len(scopes) != 2 || scopes[0] != "read" || scopes[1] != "write" {
		t.Errorf("expected scopes [read, write], got %v", scopes)
	}
}

func TestRouteWithApiKey(t *testing.T) {
	g := New()

	routes := []Route{
		{
			Path:              "/api-key-required",
			Method:            "GET",
			FunctionLogicalID: "ApiKeyFunction",
			Auth: &RouteAuth{
				ApiKeyRequired: true,
			},
		},
	}

	spec, err := g.GenerateSwagger(routes)
	if err != nil {
		t.Fatalf("GenerateSwagger failed: %v", err)
	}

	paths := spec["paths"].(map[string]interface{})
	apiKeyPath := paths["/api-key-required"].(map[string]interface{})
	getMethod := apiKeyPath["get"].(map[string]interface{})

	security, ok := getMethod["security"].([]interface{})
	if !ok || len(security) == 0 {
		t.Fatal("expected security configuration")
	}

	// Should have api_key security
	found := false
	for _, s := range security {
		if secMap, ok := s.(map[string]interface{}); ok {
			if _, hasApiKey := secMap["api_key"]; hasApiKey {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("expected api_key security")
	}
}

func TestAddCorsToSpecOpenAPI3(t *testing.T) {
	g := New()

	spec := map[string]interface{}{
		"openapi": "3.0.1",
		"paths":   map[string]interface{}{},
	}

	cors := &CorsConfig{
		AllowOrigins:     []string{"https://example.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600,
	}

	err := g.AddCorsToSpec(spec, cors)
	if err != nil {
		t.Fatalf("AddCorsToSpec failed: %v", err)
	}

	corsConfig, ok := spec["x-amazon-apigateway-cors"].(map[string]interface{})
	if !ok {
		t.Fatal("expected x-amazon-apigateway-cors")
	}

	allowOrigins := corsConfig["allowOrigins"].([]string)
	if len(allowOrigins) != 1 || allowOrigins[0] != "https://example.com" {
		t.Errorf("unexpected allowOrigins: %v", allowOrigins)
	}

	if corsConfig["allowCredentials"] != true {
		t.Error("expected allowCredentials to be true")
	}

	if corsConfig["maxAge"] != 3600 {
		t.Errorf("expected maxAge 3600, got %v", corsConfig["maxAge"])
	}
}

func TestAddCorsToSpecSwagger(t *testing.T) {
	g := New()

	spec := map[string]interface{}{
		"swagger": "2.0",
		"paths": map[string]interface{}{
			"/users": map[string]interface{}{
				"get": map[string]interface{}{},
			},
		},
	}

	cors := &CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
	}

	err := g.AddCorsToSpec(spec, cors)
	if err != nil {
		t.Fatalf("AddCorsToSpec failed: %v", err)
	}

	paths := spec["paths"].(map[string]interface{})
	usersPath := paths["/users"].(map[string]interface{})

	// Should have added OPTIONS method
	optionsMethod, ok := usersPath["options"].(map[string]interface{})
	if !ok {
		t.Fatal("expected OPTIONS method for CORS preflight")
	}

	integration := optionsMethod["x-amazon-apigateway-integration"].(map[string]interface{})
	if integration["type"] != "mock" {
		t.Errorf("expected mock integration for CORS, got %v", integration["type"])
	}
}

func TestAddSecurityDefinitions(t *testing.T) {
	g := New()

	spec := map[string]interface{}{
		"swagger": "2.0",
		"paths":   map[string]interface{}{},
	}

	authorizers := map[string]interface{}{
		"MyCognitoAuth": map[string]interface{}{
			"UserPoolArn": "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_abc123",
		},
		"MyLambdaAuth": map[string]interface{}{
			"FunctionArn": map[string]interface{}{
				"Fn::GetAtt": []string{"AuthFunction", "Arn"},
			},
		},
	}

	err := g.AddSecurityDefinitions(spec, authorizers)
	if err != nil {
		t.Fatalf("AddSecurityDefinitions failed: %v", err)
	}

	securityDefs, ok := spec["securityDefinitions"].(map[string]interface{})
	if !ok {
		t.Fatal("expected securityDefinitions")
	}

	// Check Cognito authorizer
	cognitoAuth, ok := securityDefs["MyCognitoAuth"].(map[string]interface{})
	if !ok {
		t.Fatal("expected MyCognitoAuth")
	}
	if cognitoAuth["x-amazon-apigateway-authtype"] != "cognito_user_pools" {
		t.Errorf("expected cognito_user_pools authtype, got %v", cognitoAuth["x-amazon-apigateway-authtype"])
	}

	// Check Lambda authorizer
	lambdaAuth, ok := securityDefs["MyLambdaAuth"].(map[string]interface{})
	if !ok {
		t.Fatal("expected MyLambdaAuth")
	}
	if lambdaAuth["x-amazon-apigateway-authtype"] != "custom" {
		t.Errorf("expected custom authtype, got %v", lambdaAuth["x-amazon-apigateway-authtype"])
	}
}

func TestAddSecurityDefinitionsOpenAPI3(t *testing.T) {
	g := New()

	spec := map[string]interface{}{
		"openapi": "3.0.1",
		"paths":   map[string]interface{}{},
	}

	authorizers := map[string]interface{}{
		"JwtAuth": map[string]interface{}{
			"JwtConfiguration": map[string]interface{}{
				"issuer":   "https://auth.example.com",
				"audience": []string{"api"},
			},
		},
	}

	err := g.AddSecurityDefinitions(spec, authorizers)
	if err != nil {
		t.Fatalf("AddSecurityDefinitions failed: %v", err)
	}

	components, ok := spec["components"].(map[string]interface{})
	if !ok {
		t.Fatal("expected components")
	}

	securitySchemes, ok := components["securitySchemes"].(map[string]interface{})
	if !ok {
		t.Fatal("expected securitySchemes")
	}

	jwtAuth, ok := securitySchemes["JwtAuth"].(map[string]interface{})
	if !ok {
		t.Fatal("expected JwtAuth")
	}

	if jwtAuth["type"] != "oauth2" {
		t.Errorf("expected oauth2 type, got %v", jwtAuth["type"])
	}

	authorizerConfig := jwtAuth["x-amazon-apigateway-authorizer"].(map[string]interface{})
	if authorizerConfig["type"] != "jwt" {
		t.Errorf("expected jwt type, got %v", authorizerConfig["type"])
	}
}

func TestIsOpenAPI3(t *testing.T) {
	tests := []struct {
		name     string
		spec     map[string]interface{}
		expected bool
	}{
		{
			name:     "OpenAPI 3 spec",
			spec:     map[string]interface{}{"openapi": "3.0.1"},
			expected: true,
		},
		{
			name:     "Swagger 2 spec",
			spec:     map[string]interface{}{"swagger": "2.0"},
			expected: false,
		},
		{
			name:     "Empty spec",
			spec:     map[string]interface{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsOpenAPI3(tt.spec)
			if result != tt.expected {
				t.Errorf("IsOpenAPI3(%v) = %v, want %v", tt.spec, result, tt.expected)
			}
		})
	}
}

func TestIsSwagger(t *testing.T) {
	tests := []struct {
		name     string
		spec     map[string]interface{}
		expected bool
	}{
		{
			name:     "Swagger 2 spec",
			spec:     map[string]interface{}{"swagger": "2.0"},
			expected: true,
		},
		{
			name:     "OpenAPI 3 spec",
			spec:     map[string]interface{}{"openapi": "3.0.1"},
			expected: false,
		},
		{
			name:     "Empty spec",
			spec:     map[string]interface{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSwagger(tt.spec)
			if result != tt.expected {
				t.Errorf("IsSwagger(%v) = %v, want %v", tt.spec, result, tt.expected)
			}
		})
	}
}

func TestExtractPathParameters(t *testing.T) {
	g := New()

	tests := []struct {
		path           string
		expectedParams []string
	}{
		{"/users", nil},
		{"/users/{id}", []string{"id"}},
		{"/users/{userId}/posts/{postId}", []string{"userId", "postId"}},
		{"/proxy/{proxy+}", []string{"proxy"}},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			params := g.extractPathParameters(tt.path)
			if len(params) != len(tt.expectedParams) {
				t.Errorf("expected %d params, got %d", len(tt.expectedParams), len(params))
				return
			}
			for i, expected := range tt.expectedParams {
				if params[i]["name"] != expected {
					t.Errorf("expected param %q, got %q", expected, params[i]["name"])
				}
			}
		})
	}
}

func TestGenerateSwaggerWithFunctionArn(t *testing.T) {
	g := New()

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []string{"MyFunction", "Arn"},
	}

	routes := []Route{
		{
			Path:        "/test",
			Method:      "GET",
			FunctionArn: functionArn,
		},
	}

	spec, err := g.GenerateSwagger(routes)
	if err != nil {
		t.Fatalf("GenerateSwagger failed: %v", err)
	}

	paths := spec["paths"].(map[string]interface{})
	testPath := paths["/test"].(map[string]interface{})
	getMethod := testPath["get"].(map[string]interface{})
	integration := getMethod["x-amazon-apigateway-integration"].(map[string]interface{})

	uri := integration["uri"].(map[string]interface{})
	fnSub := uri["Fn::Sub"].([]interface{})
	vars := fnSub[1].(map[string]interface{})

	if vars["FunctionArn"] == nil {
		t.Error("expected FunctionArn in Fn::Sub variables")
	}
}

func TestValidationErrors(t *testing.T) {
	g := New()

	// Test missing method
	_, err := g.GenerateSwagger([]Route{{Path: "/test"}})
	if err == nil {
		t.Error("expected error for missing method")
	}

	// Test missing path
	_, err = g.GenerateSwagger([]Route{{Method: "GET"}})
	if err == nil {
		t.Error("expected error for missing path")
	}
}

// Helper function to pretty print JSON for debugging
func prettyJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
