package apigateway

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDeployment_JSONSerialization(t *testing.T) {
	deployment := Deployment{
		RestApiId:   "api123",
		StageName:   "prod",
		Description: "Production deployment",
		DeploymentCanarySettings: &DeploymentCanarySettings{
			PercentTraffic: 10.0,
			StageVariableOverrides: map[string]interface{}{
				"lambdaAlias": "canary",
			},
			UseStageCache: false,
		},
		StageDescription: &StageDescription{
			Description:         "Production stage",
			CacheClusterEnabled: true,
			CacheClusterSize:    "0.5",
			LoggingLevel:        "INFO",
			MetricsEnabled:      true,
			Variables: map[string]interface{}{
				"environment": "production",
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(deployment)
	if err != nil {
		t.Fatalf("Failed to marshal Deployment to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Deployment
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Deployment from JSON: %v", err)
	}

	if unmarshaled.RestApiId != deployment.RestApiId {
		t.Errorf("RestApiId mismatch: got %v, want %v", unmarshaled.RestApiId, deployment.RestApiId)
	}

	if unmarshaled.StageName != deployment.StageName {
		t.Errorf("StageName mismatch: got %v, want %v", unmarshaled.StageName, deployment.StageName)
	}

	if unmarshaled.DeploymentCanarySettings == nil {
		t.Error("DeploymentCanarySettings should not be nil")
	}

	if unmarshaled.StageDescription == nil {
		t.Error("StageDescription should not be nil")
	}
}

func TestDeployment_YAMLSerialization(t *testing.T) {
	deployment := Deployment{
		RestApiId:   "api123",
		Description: "Test deployment",
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(deployment)
	if err != nil {
		t.Fatalf("Failed to marshal Deployment to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Deployment
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Deployment from YAML: %v", err)
	}

	if unmarshaled.RestApiId != deployment.RestApiId {
		t.Errorf("RestApiId mismatch: got %v, want %v", unmarshaled.RestApiId, deployment.RestApiId)
	}
}

func TestDeployment_WithIntrinsicFunctions(t *testing.T) {
	deployment := Deployment{
		RestApiId:   map[string]interface{}{"Ref": "RestApi"},
		StageName:   map[string]interface{}{"Ref": "StageName"},
		Description: map[string]interface{}{"Fn::Sub": "Deployment for ${AWS::StackName}"},
	}

	data, err := json.Marshal(deployment)
	if err != nil {
		t.Fatalf("Failed to marshal Deployment with intrinsic functions: %v", err)
	}

	var unmarshaled Deployment
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Deployment with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	restApiIdMap, ok := unmarshaled.RestApiId.(map[string]interface{})
	if !ok {
		t.Error("RestApiId should be a map for intrinsic function")
	} else if _, exists := restApiIdMap["Ref"]; !exists {
		t.Error("RestApiId should contain Ref intrinsic function")
	}
}

func TestDeploymentCanarySettings_JSONSerialization(t *testing.T) {
	settings := DeploymentCanarySettings{
		PercentTraffic: 5.0,
		StageVariableOverrides: map[string]interface{}{
			"alias": "canary",
		},
		UseStageCache: true,
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Failed to marshal DeploymentCanarySettings: %v", err)
	}

	var unmarshaled DeploymentCanarySettings
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal DeploymentCanarySettings: %v", err)
	}

	if unmarshaled.PercentTraffic != settings.PercentTraffic {
		t.Errorf("PercentTraffic mismatch: got %v, want %v",
			unmarshaled.PercentTraffic, settings.PercentTraffic)
	}
}

func TestStageDescription_JSONSerialization(t *testing.T) {
	desc := StageDescription{
		Description:          "Stage description",
		CacheClusterEnabled:  true,
		CacheClusterSize:     "1.6",
		CacheDataEncrypted:   true,
		CacheTtlInSeconds:    600,
		CachingEnabled:       true,
		DataTraceEnabled:     false,
		LoggingLevel:         "ERROR",
		MetricsEnabled:       true,
		ThrottlingBurstLimit: 1000,
		ThrottlingRateLimit:  500.0,
		TracingEnabled:       true,
		Variables: map[string]interface{}{
			"key": "value",
		},
		AccessLogSetting: &AccessLogSetting{
			DestinationArn: "arn:aws:logs:us-east-1:123456789:log-group:test",
			Format:         "$requestId",
		},
		MethodSettings: []MethodSetting{
			{
				HttpMethod:     "*",
				ResourcePath:   "/*",
				LoggingLevel:   "INFO",
				MetricsEnabled: true,
			},
		},
	}

	data, err := json.Marshal(desc)
	if err != nil {
		t.Fatalf("Failed to marshal StageDescription: %v", err)
	}

	var unmarshaled StageDescription
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal StageDescription: %v", err)
	}

	if unmarshaled.Description != desc.Description {
		t.Errorf("Description mismatch: got %v, want %v",
			unmarshaled.Description, desc.Description)
	}

	if unmarshaled.AccessLogSetting == nil {
		t.Error("AccessLogSetting should not be nil")
	}

	if len(unmarshaled.MethodSettings) != len(desc.MethodSettings) {
		t.Errorf("MethodSettings length mismatch: got %d, want %d",
			len(unmarshaled.MethodSettings), len(desc.MethodSettings))
	}
}

func TestDeployment_MinimalConfiguration(t *testing.T) {
	deployment := Deployment{
		RestApiId: "api123",
	}

	data, err := json.Marshal(deployment)
	if err != nil {
		t.Fatalf("Failed to marshal minimal Deployment: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	// These should be omitted because they're empty
	emptyFields := []string{"StageName", "Description", "DeploymentCanarySettings", "StageDescription"}
	for _, field := range emptyFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}
