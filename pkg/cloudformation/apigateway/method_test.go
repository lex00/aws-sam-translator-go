package apigateway

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMethod_JSONSerialization(t *testing.T) {
	method := Method{
		HttpMethod:        "GET",
		ResourceId:        "resource123",
		RestApiId:         "api123",
		AuthorizationType: "NONE",
		ApiKeyRequired:    false,
		OperationName:     "GetUsers",
		RequestParameters: map[string]interface{}{
			"method.request.querystring.page": true,
		},
		Integration: &Integration{
			Type:                  "AWS_PROXY",
			IntegrationHttpMethod: "POST",
			Uri:                   "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789:function:myFunction/invocations",
		},
		MethodResponses: []MethodResponse{
			{
				StatusCode: "200",
				ResponseModels: map[string]interface{}{
					"application/json": "Empty",
				},
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(method)
	if err != nil {
		t.Fatalf("Failed to marshal Method to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Method
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Method from JSON: %v", err)
	}

	if unmarshaled.HttpMethod != method.HttpMethod {
		t.Errorf("HttpMethod mismatch: got %v, want %v", unmarshaled.HttpMethod, method.HttpMethod)
	}

	if unmarshaled.ResourceId != method.ResourceId {
		t.Errorf("ResourceId mismatch: got %v, want %v", unmarshaled.ResourceId, method.ResourceId)
	}

	if unmarshaled.Integration == nil {
		t.Error("Integration should not be nil")
	}
}

func TestMethod_YAMLSerialization(t *testing.T) {
	method := Method{
		HttpMethod:        "POST",
		ResourceId:        "resource456",
		RestApiId:         "api789",
		AuthorizationType: "AWS_IAM",
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(method)
	if err != nil {
		t.Fatalf("Failed to marshal Method to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Method
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Method from YAML: %v", err)
	}

	if unmarshaled.HttpMethod != method.HttpMethod {
		t.Errorf("HttpMethod mismatch: got %v, want %v", unmarshaled.HttpMethod, method.HttpMethod)
	}
}

func TestMethod_WithIntrinsicFunctions(t *testing.T) {
	method := Method{
		HttpMethod:        "GET",
		ResourceId:        map[string]interface{}{"Ref": "ApiResource"},
		RestApiId:         map[string]interface{}{"Ref": "RestApi"},
		AuthorizationType: "CUSTOM",
		AuthorizerId:      map[string]interface{}{"Ref": "Authorizer"},
		Integration: &Integration{
			Type: "AWS_PROXY",
			Uri: map[string]interface{}{
				"Fn::Sub": "arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${LambdaFunction.Arn}/invocations",
			},
		},
	}

	data, err := json.Marshal(method)
	if err != nil {
		t.Fatalf("Failed to marshal Method with intrinsic functions: %v", err)
	}

	var unmarshaled Method
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Method with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	resourceIdMap, ok := unmarshaled.ResourceId.(map[string]interface{})
	if !ok {
		t.Error("ResourceId should be a map for intrinsic function")
	} else if _, exists := resourceIdMap["Ref"]; !exists {
		t.Error("ResourceId should contain Ref intrinsic function")
	}
}

func TestIntegration_JSONSerialization(t *testing.T) {
	integration := Integration{
		Type:                  "HTTP_PROXY",
		IntegrationHttpMethod: "GET",
		Uri:                   "https://api.example.com/users",
		ConnectionType:        "INTERNET",
		PassthroughBehavior:   "WHEN_NO_MATCH",
		TimeoutInMillis:       29000,
		CacheKeyParameters:    []interface{}{"method.request.querystring.page"},
		CacheNamespace:        "users",
		RequestParameters: map[string]interface{}{
			"integration.request.querystring.page": "method.request.querystring.page",
		},
		RequestTemplates: map[string]interface{}{
			"application/json": "{}",
		},
		IntegrationResponses: []IntegrationResponse{
			{
				StatusCode:       "200",
				SelectionPattern: "",
				ResponseTemplates: map[string]interface{}{
					"application/json": "",
				},
			},
		},
	}

	data, err := json.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal Integration: %v", err)
	}

	var unmarshaled Integration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Integration: %v", err)
	}

	if unmarshaled.Type != integration.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaled.Type, integration.Type)
	}

	if unmarshaled.ConnectionType != integration.ConnectionType {
		t.Errorf("ConnectionType mismatch: got %v, want %v",
			unmarshaled.ConnectionType, integration.ConnectionType)
	}
}

func TestIntegration_VpcLink(t *testing.T) {
	integration := Integration{
		Type:                  "HTTP_PROXY",
		IntegrationHttpMethod: "GET",
		Uri:                   "http://internal-nlb.example.com",
		ConnectionType:        "VPC_LINK",
		ConnectionId:          "vpclink123",
	}

	data, err := json.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal VPC Link Integration: %v", err)
	}

	var unmarshaled Integration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal VPC Link Integration: %v", err)
	}

	if unmarshaled.ConnectionType != "VPC_LINK" {
		t.Errorf("ConnectionType mismatch: got %v, want VPC_LINK", unmarshaled.ConnectionType)
	}

	if unmarshaled.ConnectionId != "vpclink123" {
		t.Errorf("ConnectionId mismatch: got %v, want vpclink123", unmarshaled.ConnectionId)
	}
}

func TestIntegrationResponse_JSONSerialization(t *testing.T) {
	response := IntegrationResponse{
		StatusCode:       "200",
		SelectionPattern: "2\\d{2}",
		ContentHandling:  "CONVERT_TO_TEXT",
		ResponseParameters: map[string]interface{}{
			"method.response.header.Content-Type": "'application/json'",
		},
		ResponseTemplates: map[string]interface{}{
			"application/json": "#set($inputRoot = $input.path('$'))\n$inputRoot",
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal IntegrationResponse: %v", err)
	}

	var unmarshaled IntegrationResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal IntegrationResponse: %v", err)
	}

	if unmarshaled.StatusCode != response.StatusCode {
		t.Errorf("StatusCode mismatch: got %v, want %v",
			unmarshaled.StatusCode, response.StatusCode)
	}
}

func TestMethodResponse_JSONSerialization(t *testing.T) {
	response := MethodResponse{
		StatusCode: "200",
		ResponseModels: map[string]interface{}{
			"application/json": "UserModel",
		},
		ResponseParameters: map[string]interface{}{
			"method.response.header.Content-Type": true,
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal MethodResponse: %v", err)
	}

	var unmarshaled MethodResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal MethodResponse: %v", err)
	}

	if unmarshaled.StatusCode != response.StatusCode {
		t.Errorf("StatusCode mismatch: got %v, want %v",
			unmarshaled.StatusCode, response.StatusCode)
	}
}

func TestMethod_WithAuthorizationScopes(t *testing.T) {
	method := Method{
		HttpMethod:          "GET",
		ResourceId:          "resource123",
		RestApiId:           "api123",
		AuthorizationType:   "COGNITO_USER_POOLS",
		AuthorizerId:        "authorizer123",
		AuthorizationScopes: []interface{}{"read:users", "admin"},
	}

	data, err := json.Marshal(method)
	if err != nil {
		t.Fatalf("Failed to marshal Method with scopes: %v", err)
	}

	var unmarshaled Method
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Method with scopes: %v", err)
	}

	if len(unmarshaled.AuthorizationScopes) != 2 {
		t.Errorf("AuthorizationScopes length mismatch: got %d, want 2",
			len(unmarshaled.AuthorizationScopes))
	}
}
