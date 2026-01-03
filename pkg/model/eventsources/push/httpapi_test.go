package push

import (
	"reflect"
	"testing"
)

func TestNewHttpApiEventHandler(t *testing.T) {
	event := &HttpApiEvent{
		Path:   "/hello",
		Method: "GET",
	}

	handler := NewHttpApiEventHandler("MyFunction", "MyEvent", event)

	if handler == nil {
		t.Fatal("expected handler to be non-nil")
	}
	if handler.functionLogicalID != "MyFunction" {
		t.Errorf("expected functionLogicalID 'MyFunction', got %s", handler.functionLogicalID)
	}
	if handler.eventLogicalID != "MyEvent" {
		t.Errorf("expected eventLogicalID 'MyEvent', got %s", handler.eventLogicalID)
	}
	if handler.event != event {
		t.Error("expected event to match")
	}
}

func TestHttpApiEventHandler_GenerateResources_MinimalConfig(t *testing.T) {
	event := &HttpApiEvent{}
	handler := NewHttpApiEventHandler("MyFunction", "HttpApiEvent", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resources == nil {
		t.Fatal("expected resources to be non-nil")
	}

	// Should generate Integration, Route, and Permission
	if len(resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(resources))
	}

	// Check Integration
	integration, ok := resources["MyFunctionHttpApiEventIntegration"]
	if !ok {
		t.Fatal("Integration resource should exist")
	}
	integrationMap, ok := integration.(map[string]interface{})
	if !ok {
		t.Fatal("Integration should be a map")
	}
	if integrationMap["Type"] != "AWS::ApiGatewayV2::Integration" {
		t.Errorf("expected Integration type 'AWS::ApiGatewayV2::Integration', got %v", integrationMap["Type"])
	}

	integrationProps, ok := integrationMap["Properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Integration Properties should be a map")
	}
	if integrationProps["IntegrationType"] != "AWS_PROXY" {
		t.Errorf("expected IntegrationType 'AWS_PROXY', got %v", integrationProps["IntegrationType"])
	}
	if !reflect.DeepEqual(integrationProps["IntegrationUri"], functionArn) {
		t.Errorf("expected IntegrationUri to match functionArn")
	}
	if integrationProps["PayloadFormatVersion"] != "2.0" {
		t.Errorf("expected PayloadFormatVersion '2.0', got %v", integrationProps["PayloadFormatVersion"])
	}

	// Check that ApiId defaults to ServerlessHttpApi
	apiID, ok := integrationProps["ApiId"].(map[string]interface{})
	if !ok {
		t.Fatal("ApiId should be a map")
	}
	if apiID["Ref"] != "ServerlessHttpApi" {
		t.Errorf("expected ApiId Ref 'ServerlessHttpApi', got %v", apiID["Ref"])
	}

	// Check Route
	route, ok := resources["MyFunctionHttpApiEventRoute"]
	if !ok {
		t.Fatal("Route resource should exist")
	}
	routeMap, ok := route.(map[string]interface{})
	if !ok {
		t.Fatal("Route should be a map")
	}
	if routeMap["Type"] != "AWS::ApiGatewayV2::Route" {
		t.Errorf("expected Route type 'AWS::ApiGatewayV2::Route', got %v", routeMap["Type"])
	}

	routeProps, ok := routeMap["Properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Route Properties should be a map")
	}
	// Default route key should be "ANY /"
	if routeProps["RouteKey"] != "ANY /" {
		t.Errorf("expected RouteKey 'ANY /', got %v", routeProps["RouteKey"])
	}

	// Check Permission
	permission, ok := resources["MyFunctionHttpApiEventPermission"]
	if !ok {
		t.Fatal("Permission resource should exist")
	}
	permissionMap, ok := permission.(map[string]interface{})
	if !ok {
		t.Fatal("Permission should be a map")
	}
	if permissionMap["Type"] != "AWS::Lambda::Permission" {
		t.Errorf("expected Permission type 'AWS::Lambda::Permission', got %v", permissionMap["Type"])
	}

	permissionProps, ok := permissionMap["Properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Permission Properties should be a map")
	}
	if permissionProps["Action"] != "lambda:InvokeFunction" {
		t.Errorf("expected Action 'lambda:InvokeFunction', got %v", permissionProps["Action"])
	}
	if permissionProps["Principal"] != "apigateway.amazonaws.com" {
		t.Errorf("expected Principal 'apigateway.amazonaws.com', got %v", permissionProps["Principal"])
	}
}

func TestHttpApiEventHandler_GenerateResources_WithPathAndMethod(t *testing.T) {
	event := &HttpApiEvent{
		Path:   "/users/{id}",
		Method: "GET",
	}
	handler := NewHttpApiEventHandler("MyFunction", "GetUser", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	route := resources["MyFunctionGetUserRoute"].(map[string]interface{})
	routeProps := route["Properties"].(map[string]interface{})

	if routeProps["RouteKey"] != "GET /users/{id}" {
		t.Errorf("expected RouteKey 'GET /users/{id}', got %v", routeProps["RouteKey"])
	}
}

func TestHttpApiEventHandler_GenerateResources_WithExplicitApiId(t *testing.T) {
	event := &HttpApiEvent{
		ApiId:  map[string]interface{}{"Ref": "MyHttpApi"},
		Path:   "/hello",
		Method: "POST",
	}
	handler := NewHttpApiEventHandler("MyFunction", "PostHello", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	integration := resources["MyFunctionPostHelloIntegration"].(map[string]interface{})
	integrationProps := integration["Properties"].(map[string]interface{})

	apiID, ok := integrationProps["ApiId"].(map[string]interface{})
	if !ok {
		t.Fatal("ApiId should be a map")
	}
	if apiID["Ref"] != "MyHttpApi" {
		t.Errorf("expected ApiId Ref 'MyHttpApi', got %v", apiID["Ref"])
	}
}

func TestHttpApiEventHandler_GenerateResources_WithPayloadFormatVersion(t *testing.T) {
	event := &HttpApiEvent{
		PayloadFormatVersion: "1.0",
	}
	handler := NewHttpApiEventHandler("MyFunction", "Event", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	integration := resources["MyFunctionEventIntegration"].(map[string]interface{})
	integrationProps := integration["Properties"].(map[string]interface{})

	if integrationProps["PayloadFormatVersion"] != "1.0" {
		t.Errorf("expected PayloadFormatVersion '1.0', got %v", integrationProps["PayloadFormatVersion"])
	}
}

func TestHttpApiEventHandler_GenerateResources_WithTimeout(t *testing.T) {
	event := &HttpApiEvent{
		TimeoutInMillis: 5000,
	}
	handler := NewHttpApiEventHandler("MyFunction", "Event", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	integration := resources["MyFunctionEventIntegration"].(map[string]interface{})
	integrationProps := integration["Properties"].(map[string]interface{})

	if integrationProps["TimeoutInMillis"] != 5000 {
		t.Errorf("expected TimeoutInMillis 5000, got %v", integrationProps["TimeoutInMillis"])
	}
}

func TestHttpApiEventHandler_GenerateResources_WithAuthNone(t *testing.T) {
	event := &HttpApiEvent{
		Auth: &HttpApiAuth{
			Authorizer: "NONE",
		},
	}
	handler := NewHttpApiEventHandler("MyFunction", "Event", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	route := resources["MyFunctionEventRoute"].(map[string]interface{})
	routeProps := route["Properties"].(map[string]interface{})

	if routeProps["AuthorizationType"] != "NONE" {
		t.Errorf("expected AuthorizationType 'NONE', got %v", routeProps["AuthorizationType"])
	}
}

func TestHttpApiEventHandler_GenerateResources_WithAuthIAM(t *testing.T) {
	event := &HttpApiEvent{
		Auth: &HttpApiAuth{
			Authorizer: "AWS_IAM",
		},
	}
	handler := NewHttpApiEventHandler("MyFunction", "Event", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	route := resources["MyFunctionEventRoute"].(map[string]interface{})
	routeProps := route["Properties"].(map[string]interface{})

	if routeProps["AuthorizationType"] != "AWS_IAM" {
		t.Errorf("expected AuthorizationType 'AWS_IAM', got %v", routeProps["AuthorizationType"])
	}
}

func TestHttpApiEventHandler_GenerateResources_WithAuthJWT(t *testing.T) {
	event := &HttpApiEvent{
		Auth: &HttpApiAuth{
			Authorizer:          "MyAuthorizer",
			AuthorizationScopes: []interface{}{"scope1", "scope2"},
		},
	}
	handler := NewHttpApiEventHandler("MyFunction", "Event", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	route := resources["MyFunctionEventRoute"].(map[string]interface{})
	routeProps := route["Properties"].(map[string]interface{})

	if routeProps["AuthorizationType"] != "JWT" {
		t.Errorf("expected AuthorizationType 'JWT', got %v", routeProps["AuthorizationType"])
	}
	if routeProps["AuthorizerId"] != "MyAuthorizer" {
		t.Errorf("expected AuthorizerId 'MyAuthorizer', got %v", routeProps["AuthorizerId"])
	}
	if !reflect.DeepEqual(routeProps["AuthorizationScopes"], []interface{}{"scope1", "scope2"}) {
		t.Errorf("expected AuthorizationScopes to match")
	}
}

func TestHttpApiEventHandler_GenerateResources_WithInvokeRole(t *testing.T) {
	event := &HttpApiEvent{
		Auth: &HttpApiAuth{
			InvokeRole: map[string]interface{}{
				"Fn::GetAtt": []interface{}{"MyRole", "Arn"},
			},
		},
	}
	handler := NewHttpApiEventHandler("MyFunction", "Event", event)

	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	resources, err := handler.GenerateResources(functionArn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	integration := resources["MyFunctionEventIntegration"].(map[string]interface{})
	integrationProps := integration["Properties"].(map[string]interface{})

	expectedRole := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyRole", "Arn"},
	}
	if !reflect.DeepEqual(integrationProps["CredentialsArn"], expectedRole) {
		t.Errorf("expected CredentialsArn to match")
	}
}

func TestHttpApiEventHandler_BuildRouteKey_SimpleStrings(t *testing.T) {
	testCases := []struct {
		name     string
		method   interface{}
		path     interface{}
		expected string
	}{
		{
			name:     "GET /users",
			method:   "GET",
			path:     "/users",
			expected: "GET /users",
		},
		{
			name:     "POST /items",
			method:   "POST",
			path:     "/items",
			expected: "POST /items",
		},
		{
			name:     "ANY /catch-all",
			method:   "ANY",
			path:     "/catch-all",
			expected: "ANY /catch-all",
		},
		{
			name:     "* wildcard",
			method:   "*",
			path:     "/wildcard",
			expected: "ANY /wildcard",
		},
		{
			name:     "$default route",
			method:   "ANY",
			path:     "/$default",
			expected: "$default",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := &HttpApiEventHandler{}
			result := handler.buildRouteKey(tc.method, tc.path)
			if result != tc.expected {
				t.Errorf("expected '%s', got '%v'", tc.expected, result)
			}
		})
	}
}

func TestHttpApiEventHandler_BuildRouteKey_DynamicValues(t *testing.T) {
	handler := &HttpApiEventHandler{}

	// Test with dynamic method
	result := handler.buildRouteKey(
		map[string]interface{}{"Ref": "HttpMethod"},
		"/users",
	)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected result to be a map")
	}

	fnSub, ok := resultMap["Fn::Sub"].([]interface{})
	if !ok {
		t.Fatal("expected Fn::Sub to be a slice")
	}
	if fnSub[0] != "${Method} ${Path}" {
		t.Errorf("expected template '${Method} ${Path}', got %v", fnSub[0])
	}

	params, ok := fnSub[1].(map[string]interface{})
	if !ok {
		t.Fatal("expected params to be a map")
	}
	if params["Method"] == nil {
		t.Error("expected Method to be set")
	}
	if params["Path"] != "/users" {
		t.Errorf("expected Path '/users', got %v", params["Path"])
	}
}

func TestHttpApiEventHandler_BuildPermissionSourceArn(t *testing.T) {
	handler := &HttpApiEventHandler{}

	apiID := map[string]interface{}{"Ref": "MyApi"}
	result := handler.buildPermissionSourceArn(apiID)

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("expected result to be a map")
	}

	fnSub, ok := resultMap["Fn::Sub"].([]interface{})
	if !ok {
		t.Fatal("expected Fn::Sub to be a slice")
	}
	if fnSub[0] != "arn:${AWS::Partition}:execute-api:${AWS::Region}:${AWS::AccountId}:${ApiId}/*" {
		t.Errorf("unexpected ARN template: %v", fnSub[0])
	}

	params, ok := fnSub[1].(map[string]interface{})
	if !ok {
		t.Fatal("expected params to be a map")
	}

	apiIDParam, ok := params["ApiId"].(map[string]interface{})
	if !ok {
		t.Fatal("expected ApiId param to be a map")
	}
	if apiIDParam["Ref"] != "MyApi" {
		t.Errorf("expected ApiId Ref 'MyApi', got %v", apiIDParam["Ref"])
	}
}

func TestHttpApiEvent_AllFields(t *testing.T) {
	// Test that all fields can be set without errors
	event := &HttpApiEvent{
		ApiId:  map[string]interface{}{"Ref": "MyApi"},
		Method: "GET",
		Path:   "/users",
		Auth: &HttpApiAuth{
			Authorizer:          "MyAuthorizer",
			AuthorizationScopes: []interface{}{"read:users"},
			InvokeRole:          "arn:aws:iam::123456789012:role/MyRole",
		},
		PayloadFormatVersion: "2.0",
		RouteSettings: &HttpApiRouteSettings{
			DataTraceEnabled:       true,
			DetailedMetricsEnabled: true,
			LoggingLevel:           "INFO",
			ThrottlingBurstLimit:   100,
			ThrottlingRateLimit:    50,
		},
		TimeoutInMillis: 10000,
	}

	if event.Method != "GET" {
		t.Errorf("expected Method 'GET', got %v", event.Method)
	}
	if event.Path != "/users" {
		t.Errorf("expected Path '/users', got %v", event.Path)
	}
	if event.Auth == nil {
		t.Fatal("Auth should not be nil")
	}
	if event.Auth.Authorizer != "MyAuthorizer" {
		t.Errorf("expected Authorizer 'MyAuthorizer', got %v", event.Auth.Authorizer)
	}
	if event.PayloadFormatVersion != "2.0" {
		t.Errorf("expected PayloadFormatVersion '2.0', got %v", event.PayloadFormatVersion)
	}
	if event.RouteSettings == nil {
		t.Fatal("RouteSettings should not be nil")
	}
	if event.TimeoutInMillis != 10000 {
		t.Errorf("expected TimeoutInMillis 10000, got %v", event.TimeoutInMillis)
	}
}

func TestHttpApiEventHandler_GenerateResources_WithAnyMethod(t *testing.T) {
	testCases := []struct {
		name     string
		method   interface{}
		expected string
	}{
		{
			name:     "ANY method",
			method:   "ANY",
			expected: "ANY /test",
		},
		{
			name:     "* wildcard",
			method:   "*",
			expected: "ANY /test",
		},
		{
			name:     "lowercase any",
			method:   "any",
			expected: "any /test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &HttpApiEvent{
				Method: tc.method,
				Path:   "/test",
			}
			handler := NewHttpApiEventHandler("MyFunction", "Event", event)

			functionArn := map[string]interface{}{
				"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
			}

			resources, err := handler.GenerateResources(functionArn)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			route := resources["MyFunctionEventRoute"].(map[string]interface{})
			routeProps := route["Properties"].(map[string]interface{})

			if routeProps["RouteKey"] != tc.expected {
				t.Errorf("expected RouteKey '%s', got %v", tc.expected, routeProps["RouteKey"])
			}
		})
	}
}
