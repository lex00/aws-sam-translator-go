package apigatewayv2

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestAuthorizer_JSONSerialization(t *testing.T) {
	authorizer := Authorizer{
		ApiId:                        "api123",
		AuthorizerType:               "JWT",
		Name:                         "JwtAuthorizer",
		IdentitySource:               []interface{}{"$request.header.Authorization"},
		AuthorizerResultTtlInSeconds: 300,
		JwtConfiguration: &JWTConfiguration{
			Audience: []interface{}{"https://api.example.com"},
			Issuer:   "https://auth.example.com/",
		},
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

	if unmarshaled.AuthorizerType != authorizer.AuthorizerType {
		t.Errorf("AuthorizerType mismatch: got %v, want %v",
			unmarshaled.AuthorizerType, authorizer.AuthorizerType)
	}

	if unmarshaled.JwtConfiguration == nil {
		t.Error("JwtConfiguration should not be nil")
	}
}

func TestAuthorizer_YAMLSerialization(t *testing.T) {
	authorizer := Authorizer{
		ApiId:          "api456",
		AuthorizerType: "REQUEST",
		Name:           "LambdaAuthorizer",
		AuthorizerUri:  "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789:function:authorizer/invocations",
		IdentitySource: []interface{}{
			"$request.header.Authorization",
			"$context.httpMethod",
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

	if unmarshaled.AuthorizerType != authorizer.AuthorizerType {
		t.Errorf("AuthorizerType mismatch: got %v, want %v",
			unmarshaled.AuthorizerType, authorizer.AuthorizerType)
	}
}

func TestAuthorizer_WithIntrinsicFunctions(t *testing.T) {
	authorizer := Authorizer{
		ApiId:          map[string]interface{}{"Ref": "HttpApi"},
		AuthorizerType: "REQUEST",
		Name:           map[string]interface{}{"Fn::Sub": "${AWS::StackName}-authorizer"},
		AuthorizerUri: map[string]interface{}{
			"Fn::Sub": "arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${AuthorizerFunction.Arn}/invocations",
		},
		AuthorizerCredentialsArn: map[string]interface{}{"Fn::GetAtt": []string{"AuthorizerRole", "Arn"}},
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

func TestAuthorizer_JWTType(t *testing.T) {
	authorizer := Authorizer{
		ApiId:          "api123",
		AuthorizerType: "JWT",
		Name:           "CognitoAuthorizer",
		IdentitySource: []interface{}{"$request.header.Authorization"},
		JwtConfiguration: &JWTConfiguration{
			Audience: []interface{}{
				"client-id-1",
				"client-id-2",
			},
			Issuer: "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_XXXXX",
		},
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal JWT Authorizer: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JWT Authorizer: %v", err)
	}

	if unmarshaled.AuthorizerType != "JWT" {
		t.Errorf("AuthorizerType mismatch: got %v, want JWT", unmarshaled.AuthorizerType)
	}

	if unmarshaled.JwtConfiguration == nil {
		t.Error("JwtConfiguration should not be nil")
	} else if len(unmarshaled.JwtConfiguration.Audience) != 2 {
		t.Errorf("JwtConfiguration.Audience length mismatch: got %d, want 2",
			len(unmarshaled.JwtConfiguration.Audience))
	}
}

func TestAuthorizer_REQUESTType(t *testing.T) {
	authorizer := Authorizer{
		ApiId:                          "api123",
		AuthorizerType:                 "REQUEST",
		Name:                           "LambdaAuthorizer",
		AuthorizerUri:                  "arn:aws:apigateway:us-east-1:lambda:path/invocations",
		AuthorizerPayloadFormatVersion: "2.0",
		EnableSimpleResponses:          true,
		IdentitySource: []interface{}{
			"$request.header.Authorization",
		},
		AuthorizerResultTtlInSeconds: 0,
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal REQUEST Authorizer: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal REQUEST Authorizer: %v", err)
	}

	if unmarshaled.AuthorizerType != "REQUEST" {
		t.Errorf("AuthorizerType mismatch: got %v, want REQUEST", unmarshaled.AuthorizerType)
	}

	if unmarshaled.AuthorizerPayloadFormatVersion != "2.0" {
		t.Errorf("AuthorizerPayloadFormatVersion mismatch: got %v, want 2.0",
			unmarshaled.AuthorizerPayloadFormatVersion)
	}
}

func TestAuthorizer_WithIdentityValidation(t *testing.T) {
	authorizer := Authorizer{
		ApiId:                        "api123",
		AuthorizerType:               "REQUEST",
		Name:                         "ValidatingAuthorizer",
		AuthorizerUri:                "arn:aws:apigateway:us-east-1:lambda:path/invocations",
		IdentityValidationExpression: "^Bearer [-0-9a-zA-Z._]+$",
		IdentitySource:               []interface{}{"$request.header.Authorization"},
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal Authorizer with identity validation: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Authorizer with identity validation: %v", err)
	}

	if unmarshaled.IdentityValidationExpression != authorizer.IdentityValidationExpression {
		t.Errorf("IdentityValidationExpression mismatch: got %v, want %v",
			unmarshaled.IdentityValidationExpression, authorizer.IdentityValidationExpression)
	}
}

func TestJWTConfiguration_JSONSerialization(t *testing.T) {
	config := JWTConfiguration{
		Audience: []interface{}{"audience1", "audience2"},
		Issuer:   "https://issuer.example.com/",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal JWTConfiguration: %v", err)
	}

	var unmarshaled JWTConfiguration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JWTConfiguration: %v", err)
	}

	if len(unmarshaled.Audience) != len(config.Audience) {
		t.Errorf("Audience length mismatch: got %d, want %d",
			len(unmarshaled.Audience), len(config.Audience))
	}

	if unmarshaled.Issuer != config.Issuer {
		t.Errorf("Issuer mismatch: got %v, want %v",
			unmarshaled.Issuer, config.Issuer)
	}
}

func TestAuthorizer_MultipleIdentitySources(t *testing.T) {
	authorizer := Authorizer{
		ApiId:          "api123",
		AuthorizerType: "REQUEST",
		Name:           "MultiSourceAuthorizer",
		AuthorizerUri:  "arn:aws:apigateway:us-east-1:lambda:path/invocations",
		IdentitySource: []interface{}{
			"$request.header.Authorization",
			"$request.querystring.token",
			"$context.httpMethod",
		},
	}

	data, err := json.Marshal(authorizer)
	if err != nil {
		t.Fatalf("Failed to marshal Authorizer with multiple identity sources: %v", err)
	}

	var unmarshaled Authorizer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Authorizer with multiple identity sources: %v", err)
	}

	if len(unmarshaled.IdentitySource) != 3 {
		t.Errorf("IdentitySource length mismatch: got %d, want 3",
			len(unmarshaled.IdentitySource))
	}
}
