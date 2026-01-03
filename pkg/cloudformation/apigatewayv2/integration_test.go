package apigatewayv2

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestIntegration_JSONSerialization(t *testing.T) {
	integration := Integration{
		ApiId:                "api123",
		IntegrationType:      "AWS_PROXY",
		IntegrationUri:       "arn:aws:lambda:us-east-1:123456789:function:myFunction",
		IntegrationMethod:    "POST",
		PayloadFormatVersion: "2.0",
		TimeoutInMillis:      30000,
		Description:          "Lambda integration",
	}

	// Test JSON marshaling
	data, err := json.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal Integration to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Integration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Integration from JSON: %v", err)
	}

	if unmarshaled.ApiId != integration.ApiId {
		t.Errorf("ApiId mismatch: got %v, want %v", unmarshaled.ApiId, integration.ApiId)
	}

	if unmarshaled.IntegrationType != integration.IntegrationType {
		t.Errorf("IntegrationType mismatch: got %v, want %v",
			unmarshaled.IntegrationType, integration.IntegrationType)
	}

	if unmarshaled.PayloadFormatVersion != integration.PayloadFormatVersion {
		t.Errorf("PayloadFormatVersion mismatch: got %v, want %v",
			unmarshaled.PayloadFormatVersion, integration.PayloadFormatVersion)
	}
}

func TestIntegration_YAMLSerialization(t *testing.T) {
	integration := Integration{
		ApiId:             "api456",
		IntegrationType:   "HTTP_PROXY",
		IntegrationUri:    "https://api.example.com",
		IntegrationMethod: "ANY",
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal Integration to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Integration
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Integration from YAML: %v", err)
	}

	if unmarshaled.IntegrationType != integration.IntegrationType {
		t.Errorf("IntegrationType mismatch: got %v, want %v",
			unmarshaled.IntegrationType, integration.IntegrationType)
	}
}

func TestIntegration_WithIntrinsicFunctions(t *testing.T) {
	integration := Integration{
		ApiId:           map[string]interface{}{"Ref": "HttpApi"},
		IntegrationType: "AWS_PROXY",
		IntegrationUri: map[string]interface{}{
			"Fn::Sub": "arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${LambdaFunction.Arn}/invocations",
		},
		CredentialsArn: map[string]interface{}{"Fn::GetAtt": []string{"IntegrationRole", "Arn"}},
	}

	data, err := json.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal Integration with intrinsic functions: %v", err)
	}

	var unmarshaled Integration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Integration with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	integrationUriMap, ok := unmarshaled.IntegrationUri.(map[string]interface{})
	if !ok {
		t.Error("IntegrationUri should be a map for intrinsic function")
	} else if _, exists := integrationUriMap["Fn::Sub"]; !exists {
		t.Error("IntegrationUri should contain Fn::Sub intrinsic function")
	}
}

func TestIntegration_VpcLink(t *testing.T) {
	integration := Integration{
		ApiId:             "api123",
		IntegrationType:   "HTTP_PROXY",
		IntegrationUri:    "http://internal-nlb.example.com",
		IntegrationMethod: "ANY",
		ConnectionType:    "VPC_LINK",
		ConnectionId:      "vpclink123",
		TlsConfig: &TlsConfig{
			ServerNameToVerify: "internal-nlb.example.com",
		},
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

	if unmarshaled.TlsConfig == nil {
		t.Error("TlsConfig should not be nil")
	}
}

func TestIntegration_AWSServiceIntegration(t *testing.T) {
	integration := Integration{
		ApiId:              "api123",
		IntegrationType:    "AWS_PROXY",
		IntegrationSubtype: "SQS-SendMessage",
		CredentialsArn:     "arn:aws:iam::123456789:role/api-sqs-role",
		RequestParameters: map[string]interface{}{
			"QueueUrl":    "https://sqs.us-east-1.amazonaws.com/123456789/my-queue",
			"MessageBody": "$request.body",
		},
		PayloadFormatVersion: "1.0",
	}

	data, err := json.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal AWS Service Integration: %v", err)
	}

	var unmarshaled Integration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal AWS Service Integration: %v", err)
	}

	if unmarshaled.IntegrationSubtype != "SQS-SendMessage" {
		t.Errorf("IntegrationSubtype mismatch: got %v, want SQS-SendMessage",
			unmarshaled.IntegrationSubtype)
	}
}

func TestIntegration_WithRequestParameters(t *testing.T) {
	integration := Integration{
		ApiId:             "api123",
		IntegrationType:   "HTTP_PROXY",
		IntegrationUri:    "https://api.example.com/users",
		IntegrationMethod: "GET",
		RequestParameters: map[string]interface{}{
			"overwrite:querystring.page": "$request.querystring.page",
			"overwrite:header.x-api-key": "$context.apiId",
		},
	}

	data, err := json.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal Integration with RequestParameters: %v", err)
	}

	var unmarshaled Integration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Integration with RequestParameters: %v", err)
	}

	if len(unmarshaled.RequestParameters) != 2 {
		t.Errorf("RequestParameters length mismatch: got %d, want 2",
			len(unmarshaled.RequestParameters))
	}
}

func TestIntegration_WithResponseParameters(t *testing.T) {
	integration := Integration{
		ApiId:             "api123",
		IntegrationType:   "HTTP_PROXY",
		IntegrationUri:    "https://api.example.com",
		IntegrationMethod: "GET",
		ResponseParameters: map[string]interface{}{
			"200": map[string]interface{}{
				"overwrite:header.x-custom": "$response.header.x-original",
			},
		},
	}

	data, err := json.Marshal(integration)
	if err != nil {
		t.Fatalf("Failed to marshal Integration with ResponseParameters: %v", err)
	}

	var unmarshaled Integration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Integration with ResponseParameters: %v", err)
	}

	if unmarshaled.ResponseParameters == nil {
		t.Error("ResponseParameters should not be nil")
	}
}

func TestTlsConfig_JSONSerialization(t *testing.T) {
	config := TlsConfig{
		ServerNameToVerify: "api.example.com",
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal TlsConfig: %v", err)
	}

	var unmarshaled TlsConfig
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal TlsConfig: %v", err)
	}

	if unmarshaled.ServerNameToVerify != config.ServerNameToVerify {
		t.Errorf("ServerNameToVerify mismatch: got %v, want %v",
			unmarshaled.ServerNameToVerify, config.ServerNameToVerify)
	}
}
