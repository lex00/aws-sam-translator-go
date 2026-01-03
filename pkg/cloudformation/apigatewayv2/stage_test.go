package apigatewayv2

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStage_JSONSerialization(t *testing.T) {
	stage := Stage{
		ApiId:       "api123",
		StageName:   "$default",
		AutoDeploy:  true,
		Description: "Default stage",
		AccessLogSettings: &AccessLogSettings{
			DestinationArn: "arn:aws:logs:us-east-1:123456789:log-group:api-logs",
			Format:         "$requestId $httpMethod $resourcePath $status",
		},
		DefaultRouteSettings: &RouteSettings{
			ThrottlingBurstLimit:   500,
			ThrottlingRateLimit:    1000.0,
			DetailedMetricsEnabled: true,
		},
		StageVariables: map[string]interface{}{
			"environment": "production",
		},
		Tags: map[string]interface{}{
			"Stage": "Default",
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

	if unmarshaled.ApiId != stage.ApiId {
		t.Errorf("ApiId mismatch: got %v, want %v", unmarshaled.ApiId, stage.ApiId)
	}

	if unmarshaled.StageName != stage.StageName {
		t.Errorf("StageName mismatch: got %v, want %v", unmarshaled.StageName, stage.StageName)
	}

	if unmarshaled.AccessLogSettings == nil {
		t.Error("AccessLogSettings should not be nil")
	}

	if unmarshaled.DefaultRouteSettings == nil {
		t.Error("DefaultRouteSettings should not be nil")
	}
}

func TestStage_YAMLSerialization(t *testing.T) {
	stage := Stage{
		ApiId:        "api456",
		StageName:    "prod",
		DeploymentId: "deploy123",
		Description:  "Production stage",
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

	if unmarshaled.ApiId != stage.ApiId {
		t.Errorf("ApiId mismatch: got %v, want %v", unmarshaled.ApiId, stage.ApiId)
	}
}

func TestStage_WithIntrinsicFunctions(t *testing.T) {
	stage := Stage{
		ApiId:        map[string]interface{}{"Ref": "HttpApi"},
		StageName:    "$default",
		DeploymentId: map[string]interface{}{"Ref": "ApiDeployment"},
		AccessLogSettings: &AccessLogSettings{
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
	apiIdMap, ok := unmarshaled.ApiId.(map[string]interface{})
	if !ok {
		t.Error("ApiId should be a map for intrinsic function")
	} else if _, exists := apiIdMap["Ref"]; !exists {
		t.Error("ApiId should contain Ref intrinsic function")
	}
}

func TestAccessLogSettings_JSONSerialization(t *testing.T) {
	settings := AccessLogSettings{
		DestinationArn: "arn:aws:logs:us-east-1:123456789:log-group:test",
		Format:         `{"requestId":"$context.requestId","ip":"$context.identity.sourceIp"}`,
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Failed to marshal AccessLogSettings: %v", err)
	}

	var unmarshaled AccessLogSettings
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal AccessLogSettings: %v", err)
	}

	if unmarshaled.DestinationArn != settings.DestinationArn {
		t.Errorf("DestinationArn mismatch: got %v, want %v",
			unmarshaled.DestinationArn, settings.DestinationArn)
	}
}

func TestRouteSettings_JSONSerialization(t *testing.T) {
	settings := RouteSettings{
		DataTraceEnabled:       true,
		DetailedMetricsEnabled: true,
		LoggingLevel:           "INFO",
		ThrottlingBurstLimit:   1000,
		ThrottlingRateLimit:    2000.0,
	}

	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Failed to marshal RouteSettings: %v", err)
	}

	var unmarshaled RouteSettings
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RouteSettings: %v", err)
	}

	if unmarshaled.LoggingLevel != settings.LoggingLevel {
		t.Errorf("LoggingLevel mismatch: got %v, want %v",
			unmarshaled.LoggingLevel, settings.LoggingLevel)
	}
}

func TestStage_WithRouteSettings(t *testing.T) {
	stage := Stage{
		ApiId:     "api123",
		StageName: "prod",
		RouteSettings: map[string]interface{}{
			"GET /users": map[string]interface{}{
				"ThrottlingBurstLimit": 100,
				"ThrottlingRateLimit":  50.0,
			},
			"POST /users": map[string]interface{}{
				"ThrottlingBurstLimit": 50,
				"ThrottlingRateLimit":  25.0,
			},
		},
	}

	data, err := json.Marshal(stage)
	if err != nil {
		t.Fatalf("Failed to marshal Stage with RouteSettings: %v", err)
	}

	var unmarshaled Stage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Stage with RouteSettings: %v", err)
	}

	if len(unmarshaled.RouteSettings) != 2 {
		t.Errorf("RouteSettings length mismatch: got %d, want 2",
			len(unmarshaled.RouteSettings))
	}
}

func TestStage_DefaultStage(t *testing.T) {
	stage := Stage{
		ApiId:      "api123",
		StageName:  "$default",
		AutoDeploy: true,
	}

	data, err := json.Marshal(stage)
	if err != nil {
		t.Fatalf("Failed to marshal $default Stage: %v", err)
	}

	var unmarshaled Stage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal $default Stage: %v", err)
	}

	if unmarshaled.StageName != "$default" {
		t.Errorf("StageName mismatch: got %v, want $default", unmarshaled.StageName)
	}
}
