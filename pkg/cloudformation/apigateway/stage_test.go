package apigateway

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStage_JSONSerialization(t *testing.T) {
	stage := Stage{
		StageName:           "prod",
		RestApiId:           "abc123",
		DeploymentId:        "deploy123",
		Description:         "Production stage",
		CacheClusterEnabled: true,
		CacheClusterSize:    "0.5",
		TracingEnabled:      true,
		Variables: map[string]interface{}{
			"lambdaAlias": "live",
		},
		AccessLogSetting: &AccessLogSetting{
			DestinationArn: "arn:aws:logs:us-east-1:123456789:log-group:api-access-logs",
			Format:         "$requestId $httpMethod $resourcePath",
		},
		MethodSettings: []MethodSetting{
			{
				HttpMethod:     "*",
				ResourcePath:   "/*",
				LoggingLevel:   "INFO",
				MetricsEnabled: true,
			},
		},
		Tags: []Tag{
			{Key: "Stage", Value: "Production"},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(stage)
	if err != nil {
		t.Fatalf("Failed to marshal Stage to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Stage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Stage from JSON: %v", err)
	}

	if unmarshaled.StageName != stage.StageName {
		t.Errorf("StageName mismatch: got %v, want %v", unmarshaled.StageName, stage.StageName)
	}

	if unmarshaled.RestApiId != stage.RestApiId {
		t.Errorf("RestApiId mismatch: got %v, want %v", unmarshaled.RestApiId, stage.RestApiId)
	}

	if unmarshaled.AccessLogSetting == nil {
		t.Error("AccessLogSetting should not be nil")
	}

	if len(unmarshaled.MethodSettings) != len(stage.MethodSettings) {
		t.Errorf("MethodSettings length mismatch: got %d, want %d",
			len(unmarshaled.MethodSettings), len(stage.MethodSettings))
	}
}

func TestStage_YAMLSerialization(t *testing.T) {
	stage := Stage{
		StageName:    "dev",
		RestApiId:    "api123",
		DeploymentId: "deploy456",
		Variables: map[string]interface{}{
			"environment": "development",
		},
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(stage)
	if err != nil {
		t.Fatalf("Failed to marshal Stage to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Stage
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Stage from YAML: %v", err)
	}

	if unmarshaled.StageName != stage.StageName {
		t.Errorf("StageName mismatch: got %v, want %v", unmarshaled.StageName, stage.StageName)
	}
}

func TestStage_WithIntrinsicFunctions(t *testing.T) {
	stage := Stage{
		StageName:    map[string]interface{}{"Ref": "StageName"},
		RestApiId:    map[string]interface{}{"Ref": "RestApi"},
		DeploymentId: map[string]interface{}{"Ref": "Deployment"},
		AccessLogSetting: &AccessLogSetting{
			DestinationArn: map[string]interface{}{"Fn::GetAtt": []string{"LogGroup", "Arn"}},
			Format:         "$requestId",
		},
	}

	data, err := json.Marshal(stage)
	if err != nil {
		t.Fatalf("Failed to marshal Stage with intrinsic functions: %v", err)
	}

	var unmarshaled Stage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Stage with intrinsic functions: %v", err)
	}

	// Verify intrinsic function structure is preserved
	restApiIdMap, ok := unmarshaled.RestApiId.(map[string]interface{})
	if !ok {
		t.Error("RestApiId should be a map for intrinsic function")
	} else if _, exists := restApiIdMap["Ref"]; !exists {
		t.Error("RestApiId should contain Ref intrinsic function")
	}
}

func TestAccessLogSetting_JSONSerialization(t *testing.T) {
	setting := AccessLogSetting{
		DestinationArn: "arn:aws:logs:us-east-1:123456789:log-group:test",
		Format:         "$requestId $httpMethod",
	}

	data, err := json.Marshal(setting)
	if err != nil {
		t.Fatalf("Failed to marshal AccessLogSetting: %v", err)
	}

	var unmarshaled AccessLogSetting
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal AccessLogSetting: %v", err)
	}

	if unmarshaled.DestinationArn != setting.DestinationArn {
		t.Errorf("DestinationArn mismatch: got %v, want %v",
			unmarshaled.DestinationArn, setting.DestinationArn)
	}
}

func TestCanarySetting_JSONSerialization(t *testing.T) {
	setting := CanarySetting{
		DeploymentId:   "deploy789",
		PercentTraffic: 10.0,
		StageVariableOverrides: map[string]interface{}{
			"lambdaAlias": "canary",
		},
		UseStageCache: false,
	}

	data, err := json.Marshal(setting)
	if err != nil {
		t.Fatalf("Failed to marshal CanarySetting: %v", err)
	}

	var unmarshaled CanarySetting
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal CanarySetting: %v", err)
	}

	if unmarshaled.DeploymentId != setting.DeploymentId {
		t.Errorf("DeploymentId mismatch: got %v, want %v",
			unmarshaled.DeploymentId, setting.DeploymentId)
	}
}

func TestMethodSetting_JSONSerialization(t *testing.T) {
	setting := MethodSetting{
		HttpMethod:           "GET",
		ResourcePath:         "/users",
		LoggingLevel:         "ERROR",
		MetricsEnabled:       true,
		DataTraceEnabled:     false,
		CachingEnabled:       true,
		CacheTtlInSeconds:    300,
		ThrottlingBurstLimit: 500,
		ThrottlingRateLimit:  1000.0,
	}

	data, err := json.Marshal(setting)
	if err != nil {
		t.Fatalf("Failed to marshal MethodSetting: %v", err)
	}

	var unmarshaled MethodSetting
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal MethodSetting: %v", err)
	}

	if unmarshaled.HttpMethod != setting.HttpMethod {
		t.Errorf("HttpMethod mismatch: got %v, want %v",
			unmarshaled.HttpMethod, setting.HttpMethod)
	}

	if unmarshaled.ResourcePath != setting.ResourcePath {
		t.Errorf("ResourcePath mismatch: got %v, want %v",
			unmarshaled.ResourcePath, setting.ResourcePath)
	}

	if unmarshaled.LoggingLevel != setting.LoggingLevel {
		t.Errorf("LoggingLevel mismatch: got %v, want %v",
			unmarshaled.LoggingLevel, setting.LoggingLevel)
	}
}

func TestStage_WithCanarySetting(t *testing.T) {
	stage := Stage{
		StageName:    "prod",
		RestApiId:    "api123",
		DeploymentId: "deploy123",
		CanarySetting: &CanarySetting{
			PercentTraffic: 5.0,
			UseStageCache:  true,
		},
	}

	data, err := json.Marshal(stage)
	if err != nil {
		t.Fatalf("Failed to marshal Stage with CanarySetting: %v", err)
	}

	var unmarshaled Stage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Stage with CanarySetting: %v", err)
	}

	if unmarshaled.CanarySetting == nil {
		t.Error("CanarySetting should not be nil")
	}
}
