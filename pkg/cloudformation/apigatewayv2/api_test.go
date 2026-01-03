package apigatewayv2

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestApi_JSONSerialization(t *testing.T) {
	api := Api{
		Name:                      "MyHttpApi",
		Description:               "My HTTP API",
		ProtocolType:              "HTTP",
		RouteSelectionExpression:  "$request.method $request.path",
		DisableExecuteApiEndpoint: false,
		CorsConfiguration: &Cors{
			AllowOrigins:     []interface{}{"https://example.com"},
			AllowMethods:     []interface{}{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []interface{}{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           3600,
		},
		Tags: map[string]interface{}{
			"Environment": "Production",
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("Failed to marshal Api to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Api
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Api from JSON: %v", err)
	}

	if unmarshaled.Name != api.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, api.Name)
	}

	if unmarshaled.ProtocolType != api.ProtocolType {
		t.Errorf("ProtocolType mismatch: got %v, want %v", unmarshaled.ProtocolType, api.ProtocolType)
	}

	if unmarshaled.CorsConfiguration == nil {
		t.Error("CorsConfiguration should not be nil")
	}
}

func TestApi_YAMLSerialization(t *testing.T) {
	api := Api{
		Name:                      "WebSocketApi",
		Description:               "My WebSocket API",
		ProtocolType:              "WEBSOCKET",
		RouteSelectionExpression:  "$request.body.action",
		ApiKeySelectionExpression: "$request.header.x-api-key",
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(api)
	if err != nil {
		t.Fatalf("Failed to marshal Api to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Api
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Api from YAML: %v", err)
	}

	if unmarshaled.Name != api.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, api.Name)
	}

	if unmarshaled.ProtocolType != api.ProtocolType {
		t.Errorf("ProtocolType mismatch: got %v, want %v", unmarshaled.ProtocolType, api.ProtocolType)
	}
}

func TestApi_WithIntrinsicFunctions(t *testing.T) {
	api := Api{
		Name:         map[string]interface{}{"Fn::Sub": "${AWS::StackName}-api"},
		Description:  map[string]interface{}{"Ref": "ApiDescription"},
		ProtocolType: "HTTP",
		CorsConfiguration: &Cors{
			AllowOrigins: []interface{}{
				map[string]interface{}{"Ref": "AllowedOrigin"},
			},
		},
	}

	data, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("Failed to marshal Api with intrinsic functions: %v", err)
	}

	var unmarshaled Api
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Api with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	nameMap, ok := unmarshaled.Name.(map[string]interface{})
	if !ok {
		t.Error("Name should be a map for intrinsic function")
	} else if _, exists := nameMap["Fn::Sub"]; !exists {
		t.Error("Name should contain Fn::Sub intrinsic function")
	}
}

func TestApi_QuickCreate(t *testing.T) {
	// Quick create API with Target
	api := Api{
		Name:         "QuickCreateApi",
		ProtocolType: "HTTP",
		Target:       "arn:aws:lambda:us-east-1:123456789:function:myFunction",
		RouteKey:     "GET /items",
	}

	data, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("Failed to marshal Quick Create Api: %v", err)
	}

	var unmarshaled Api
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Quick Create Api: %v", err)
	}

	if unmarshaled.Target != api.Target {
		t.Errorf("Target mismatch: got %v, want %v", unmarshaled.Target, api.Target)
	}

	if unmarshaled.RouteKey != api.RouteKey {
		t.Errorf("RouteKey mismatch: got %v, want %v", unmarshaled.RouteKey, api.RouteKey)
	}
}

func TestApi_WithBodyS3Location(t *testing.T) {
	api := Api{
		Name:         "ImportedApi",
		ProtocolType: "HTTP",
		BodyS3Location: &BodyS3Location{
			Bucket:  "my-bucket",
			Key:     "openapi.yaml",
			Version: "v1",
		},
	}

	data, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("Failed to marshal Api with BodyS3Location: %v", err)
	}

	var unmarshaled Api
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Api with BodyS3Location: %v", err)
	}

	if unmarshaled.BodyS3Location == nil {
		t.Error("BodyS3Location should not be nil")
	} else if unmarshaled.BodyS3Location.Bucket != api.BodyS3Location.Bucket {
		t.Errorf("BodyS3Location.Bucket mismatch: got %v, want %v",
			unmarshaled.BodyS3Location.Bucket, api.BodyS3Location.Bucket)
	}
}

func TestApi_OmitEmpty(t *testing.T) {
	api := Api{
		Name:         "MinimalApi",
		ProtocolType: "HTTP",
	}

	data, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("Failed to marshal Api: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	// These should be omitted because they're empty
	emptyFields := []string{"Description", "Body", "BodyS3Location", "CorsConfiguration", "Tags"}
	for _, field := range emptyFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestBodyS3Location_JSONSerialization(t *testing.T) {
	location := BodyS3Location{
		Bucket:  "my-bucket",
		Key:     "openapi.json",
		Etag:    "abc123",
		Version: "v1",
	}

	data, err := json.Marshal(location)
	if err != nil {
		t.Fatalf("Failed to marshal BodyS3Location: %v", err)
	}

	var unmarshaled BodyS3Location
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal BodyS3Location: %v", err)
	}

	if unmarshaled.Bucket != location.Bucket {
		t.Errorf("Bucket mismatch: got %v, want %v", unmarshaled.Bucket, location.Bucket)
	}
}

func TestCors_JSONSerialization(t *testing.T) {
	cors := Cors{
		AllowCredentials: true,
		AllowHeaders:     []interface{}{"Content-Type", "X-Amz-Date"},
		AllowMethods:     []interface{}{"GET", "POST"},
		AllowOrigins:     []interface{}{"https://example.com", "https://api.example.com"},
		ExposeHeaders:    []interface{}{"X-Custom-Header"},
		MaxAge:           7200,
	}

	data, err := json.Marshal(cors)
	if err != nil {
		t.Fatalf("Failed to marshal Cors: %v", err)
	}

	var unmarshaled Cors
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Cors: %v", err)
	}

	if len(unmarshaled.AllowOrigins) != len(cors.AllowOrigins) {
		t.Errorf("AllowOrigins length mismatch: got %d, want %d",
			len(unmarshaled.AllowOrigins), len(cors.AllowOrigins))
	}

	// MaxAge is interface{}, so after JSON unmarshal it becomes float64
	if maxAge, ok := unmarshaled.MaxAge.(float64); !ok || maxAge != 7200 {
		t.Errorf("MaxAge mismatch: got %v, want 7200", unmarshaled.MaxAge)
	}
}
