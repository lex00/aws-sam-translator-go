package apigatewayv2

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRoute_JSONSerialization(t *testing.T) {
	route := Route{
		ApiId:             "api123",
		RouteKey:          "GET /users",
		Target:            "integrations/int123",
		AuthorizationType: "JWT",
		AuthorizerId:      "auth123",
		AuthorizationScopes: []interface{}{
			"read:users",
			"admin",
		},
		OperationName: "GetUsers",
	}

	// Test JSON marshaling
	data, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal Route to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Route
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Route from JSON: %v", err)
	}

	if unmarshaled.ApiId != route.ApiId {
		t.Errorf("ApiId mismatch: got %v, want %v", unmarshaled.ApiId, route.ApiId)
	}

	if unmarshaled.RouteKey != route.RouteKey {
		t.Errorf("RouteKey mismatch: got %v, want %v", unmarshaled.RouteKey, route.RouteKey)
	}

	if unmarshaled.AuthorizationType != route.AuthorizationType {
		t.Errorf("AuthorizationType mismatch: got %v, want %v",
			unmarshaled.AuthorizationType, route.AuthorizationType)
	}

	if len(unmarshaled.AuthorizationScopes) != 2 {
		t.Errorf("AuthorizationScopes length mismatch: got %d, want 2",
			len(unmarshaled.AuthorizationScopes))
	}
}

func TestRoute_YAMLSerialization(t *testing.T) {
	route := Route{
		ApiId:             "api456",
		RouteKey:          "POST /items",
		Target:            "integrations/int456",
		AuthorizationType: "NONE",
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal Route to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Route
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Route from YAML: %v", err)
	}

	if unmarshaled.RouteKey != route.RouteKey {
		t.Errorf("RouteKey mismatch: got %v, want %v", unmarshaled.RouteKey, route.RouteKey)
	}
}

func TestRoute_WithIntrinsicFunctions(t *testing.T) {
	route := Route{
		ApiId:    map[string]interface{}{"Ref": "HttpApi"},
		RouteKey: "GET /users",
		Target: map[string]interface{}{
			"Fn::Sub": "integrations/${LambdaIntegration}",
		},
		AuthorizerId: map[string]interface{}{"Ref": "JwtAuthorizer"},
	}

	data, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal Route with intrinsic functions: %v", err)
	}

	var unmarshaled Route
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Route with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	targetMap, ok := unmarshaled.Target.(map[string]interface{})
	if !ok {
		t.Error("Target should be a map for intrinsic function")
	} else if _, exists := targetMap["Fn::Sub"]; !exists {
		t.Error("Target should contain Fn::Sub intrinsic function")
	}
}

func TestRoute_WebSocketRoutes(t *testing.T) {
	testCases := []struct {
		name     string
		routeKey string
	}{
		{"Connect", "$connect"},
		{"Disconnect", "$disconnect"},
		{"Default", "$default"},
		{"Custom", "sendMessage"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route := Route{
				ApiId:    "ws-api123",
				RouteKey: tc.routeKey,
				Target:   "integrations/int123",
			}

			data, err := json.Marshal(route)
			if err != nil {
				t.Fatalf("Failed to marshal WebSocket Route: %v", err)
			}

			var unmarshaled Route
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal WebSocket Route: %v", err)
			}

			if unmarshaled.RouteKey != tc.routeKey {
				t.Errorf("RouteKey mismatch: got %v, want %v",
					unmarshaled.RouteKey, tc.routeKey)
			}
		})
	}
}

func TestRoute_HTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD", "ANY"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			route := Route{
				ApiId:    "api123",
				RouteKey: method + " /test",
				Target:   "integrations/int123",
			}

			data, err := json.Marshal(route)
			if err != nil {
				t.Fatalf("Failed to marshal Route with %s method: %v", method, err)
			}

			var unmarshaled Route
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal Route with %s method: %v", method, err)
			}

			expectedRouteKey := method + " /test"
			if unmarshaled.RouteKey != expectedRouteKey {
				t.Errorf("RouteKey mismatch: got %v, want %v",
					unmarshaled.RouteKey, expectedRouteKey)
			}
		})
	}
}

func TestRoute_WithRequestModels(t *testing.T) {
	route := Route{
		ApiId:                    "api123",
		RouteKey:                 "POST /users",
		Target:                   "integrations/int123",
		ModelSelectionExpression: "$request.body.type",
		RequestModels: map[string]interface{}{
			"user": "UserModel",
		},
	}

	data, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal Route with RequestModels: %v", err)
	}

	var unmarshaled Route
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Route with RequestModels: %v", err)
	}

	if unmarshaled.ModelSelectionExpression != route.ModelSelectionExpression {
		t.Errorf("ModelSelectionExpression mismatch: got %v, want %v",
			unmarshaled.ModelSelectionExpression, route.ModelSelectionExpression)
	}
}

func TestRoute_WithRequestParameters(t *testing.T) {
	route := Route{
		ApiId:    "api123",
		RouteKey: "GET /users/{userId}",
		Target:   "integrations/int123",
		RequestParameters: map[string]interface{}{
			"userId": map[string]interface{}{
				"Required": true,
			},
		},
	}

	data, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal Route with RequestParameters: %v", err)
	}

	var unmarshaled Route
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Route with RequestParameters: %v", err)
	}

	if unmarshaled.RequestParameters == nil {
		t.Error("RequestParameters should not be nil")
	}
}

func TestRoute_AuthorizationTypes(t *testing.T) {
	authTypes := []string{"NONE", "AWS_IAM", "CUSTOM", "JWT"}

	for _, authType := range authTypes {
		t.Run(authType, func(t *testing.T) {
			route := Route{
				ApiId:             "api123",
				RouteKey:          "GET /test",
				AuthorizationType: authType,
			}

			data, err := json.Marshal(route)
			if err != nil {
				t.Fatalf("Failed to marshal Route with %s auth: %v", authType, err)
			}

			var unmarshaled Route
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal Route with %s auth: %v", authType, err)
			}

			if unmarshaled.AuthorizationType != authType {
				t.Errorf("AuthorizationType mismatch: got %v, want %v",
					unmarshaled.AuthorizationType, authType)
			}
		})
	}
}
