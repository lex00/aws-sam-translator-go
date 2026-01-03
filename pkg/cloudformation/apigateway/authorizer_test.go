package apigateway

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestAuthorizer_JSONSerialization(t *testing.T) {
	authorizer := Authorizer{
		Name:                         "MyAuthorizer",
		RestApiId:                    "api123",
		Type:                         "TOKEN",
		AuthorizerUri:                "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789:function:authorizer/invocations",
		AuthorizerCredentials:        "arn:aws:iam::123456789:role/authorizer-role",
		IdentitySource:               "method.request.header.Authorization",
		IdentityValidationExpression: "^Bearer [-0-9a-zA-Z._]+$",
		AuthorizerResultTtlInSeconds: 300,
	}

	// Test JSON marshaling
	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal Authorizer to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Authorizer from JSON: %v", err)
	}

	if unmarshaled.Name != authorizer.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, authorizer.Name)
	}

	if unmarshaled.RestApiId != authorizer.RestApiId {
		t.Errorf("RestApiId mismatch: got %v, want %v", unmarshaled.RestApiId, authorizer.RestApiId)
	}

	if unmarshaled.Type != authorizer.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaled.Type, authorizer.Type)
	}
}

func TestAuthorizer_YAMLSerialization(t *testing.T) {
	authorizer := Authorizer{
		Name:           "CognitoAuthorizer",
		RestApiId:      "api123",
		Type:           "COGNITO_USER_POOLS",
		IdentitySource: "method.request.header.Authorization",
		ProviderARNs: []interface{}{
			"arn:aws:cognito-idp:us-east-1:123456789:userpool/us-east-1_XXXXX",
		},
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal Authorizer to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Authorizer
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Authorizer from YAML: %v", err)
	}

	if unmarshaled.Name != authorizer.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, authorizer.Name)
	}

	if unmarshaled.Type != authorizer.Type {
		t.Errorf("Type mismatch: got %v, want %v", unmarshaled.Type, authorizer.Type)
	}

	if len(unmarshaled.ProviderARNs) != len(authorizer.ProviderARNs) {
		t.Errorf("ProviderARNs length mismatch: got %d, want %d",
			len(unmarshaled.ProviderARNs), len(authorizer.ProviderARNs))
	}
}

func TestAuthorizer_WithIntrinsicFunctions(t *testing.T) {
	authorizer := Authorizer{
		Name:      map[string]interface{}{"Fn::Sub": "${AWS::StackName}-authorizer"},
		RestApiId: map[string]interface{}{"Ref": "RestApi"},
		Type:      "REQUEST",
		AuthorizerUri: map[string]interface{}{
			"Fn::Sub": "arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${AuthorizerFunction.Arn}/invocations",
		},
		AuthorizerCredentials: map[string]interface{}{"Fn::GetAtt": []string{"AuthorizerRole", "Arn"}},
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal Authorizer with intrinsic functions: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Authorizer with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	nameMap, ok := unmarshaled.Name.(map[string]interface{})
	if !ok {
		t.Error("Name should be a map for intrinsic function")
	} else if _, exists := nameMap["Fn::Sub"]; !exists {
		t.Error("Name should contain Fn::Sub intrinsic function")
	}
}

func TestAuthorizer_TokenType(t *testing.T) {
	authorizer := Authorizer{
		Name:                         "TokenAuthorizer",
		RestApiId:                    "api123",
		Type:                         "TOKEN",
		AuthorizerUri:                "arn:aws:apigateway:us-east-1:lambda:path/invocations",
		IdentitySource:               "method.request.header.Authorization",
		IdentityValidationExpression: "^Bearer .+$",
		AuthorizerResultTtlInSeconds: 0,
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal TOKEN authorizer: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal TOKEN authorizer: %v", err)
	}

	if unmarshaled.Type != "TOKEN" {
		t.Errorf("Type mismatch: got %v, want TOKEN", unmarshaled.Type)
	}
}

func TestAuthorizer_RequestType(t *testing.T) {
	authorizer := Authorizer{
		Name:           "RequestAuthorizer",
		RestApiId:      "api123",
		Type:           "REQUEST",
		AuthorizerUri:  "arn:aws:apigateway:us-east-1:lambda:path/invocations",
		IdentitySource: "method.request.header.Authorization,context.httpMethod",
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal REQUEST authorizer: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal REQUEST authorizer: %v", err)
	}

	if unmarshaled.Type != "REQUEST" {
		t.Errorf("Type mismatch: got %v, want REQUEST", unmarshaled.Type)
	}
}

func TestAuthorizer_CognitoUserPoolsType(t *testing.T) {
	authorizer := Authorizer{
		Name:           "CognitoAuthorizer",
		RestApiId:      "api123",
		Type:           "COGNITO_USER_POOLS",
		IdentitySource: "method.request.header.Authorization",
		ProviderARNs: []interface{}{
			"arn:aws:cognito-idp:us-east-1:123456789:userpool/us-east-1_ABC123",
			"arn:aws:cognito-idp:us-east-1:123456789:userpool/us-east-1_DEF456",
		},
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal COGNITO_USER_POOLS authorizer: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal COGNITO_USER_POOLS authorizer: %v", err)
	}

	if unmarshaled.Type != "COGNITO_USER_POOLS" {
		t.Errorf("Type mismatch: got %v, want COGNITO_USER_POOLS", unmarshaled.Type)
	}

	if len(unmarshaled.ProviderARNs) != 2 {
		t.Errorf("ProviderARNs length mismatch: got %d, want 2", len(unmarshaled.ProviderARNs))
	}
}
