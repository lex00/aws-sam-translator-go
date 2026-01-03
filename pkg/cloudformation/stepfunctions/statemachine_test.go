package stepfunctions

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStateMachine_JSONSerialization(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyStateMachine",
		DefinitionString: `{"StartAt": "HelloWorld", "States": {"HelloWorld": {"Type": "Pass", "End": true}}}`,
		RoleArn:          "arn:aws:iam::123456789012:role/StepFunctionsRole",
		StateMachineType: "STANDARD",
		Tags: []Tag{
			{Key: "Environment", Value: "Production"},
		},
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine to JSON: %v", err)
	}

	var unmarshaled StateMachine
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal state machine from JSON: %v", err)
	}

	if unmarshaled.StateMachineName != sm.StateMachineName {
		t.Errorf("StateMachineName mismatch: got %v, want %v", unmarshaled.StateMachineName, sm.StateMachineName)
	}

	if unmarshaled.StateMachineType != sm.StateMachineType {
		t.Errorf("StateMachineType mismatch: got %v, want %v", unmarshaled.StateMachineType, sm.StateMachineType)
	}
}

func TestStateMachine_YAMLSerialization(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyStateMachine",
		Definition: map[string]interface{}{
			"StartAt": "HelloWorld",
			"States": map[string]interface{}{
				"HelloWorld": map[string]interface{}{
					"Type": "Pass",
					"End":  true,
				},
			},
		},
		RoleArn: "arn:aws:iam::123456789012:role/StepFunctionsRole",
	}

	data, err := yaml.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine to YAML: %v", err)
	}

	var unmarshaled StateMachine
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal state machine from YAML: %v", err)
	}

	if unmarshaled.StateMachineName != sm.StateMachineName {
		t.Errorf("StateMachineName mismatch: got %v, want %v", unmarshaled.StateMachineName, sm.StateMachineName)
	}
}

func TestStateMachine_WithIntrinsicFunctions(t *testing.T) {
	sm := StateMachine{
		StateMachineName: map[string]interface{}{
			"Fn::Sub": "${AWS::StackName}-state-machine",
		},
		RoleArn: map[string]interface{}{
			"Fn::GetAtt": []string{"StepFunctionsRole", "Arn"},
		},
		DefinitionSubstitutions: map[string]interface{}{
			"LambdaArn": map[string]interface{}{
				"Fn::GetAtt": []string{"MyFunction", "Arn"},
			},
		},
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine with intrinsics: %v", err)
	}

	var unmarshaled StateMachine
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal state machine with intrinsics: %v", err)
	}
}

func TestStateMachine_WithS3Location(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyStateMachine",
		DefinitionS3Location: &S3Location{
			Bucket:  "my-definitions-bucket",
			Key:     "state-machines/my-definition.json",
			Version: "v1",
		},
		RoleArn: "arn:aws:iam::123456789012:role/StepFunctionsRole",
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine with S3 location: %v", err)
	}

	var unmarshaled StateMachine
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal state machine with S3 location: %v", err)
	}

	if unmarshaled.DefinitionS3Location == nil {
		t.Error("DefinitionS3Location should not be nil")
	}
}

func TestStateMachine_WithLoggingConfiguration(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyStateMachine",
		LoggingConfiguration: &LoggingConfiguration{
			Level:                "ALL",
			IncludeExecutionData: true,
			Destinations: []LogDestination{
				{
					CloudWatchLogsLogGroup: &CloudWatchLogsLogGroup{
						LogGroupArn: "arn:aws:logs:us-east-1:123456789012:log-group:/aws/stepfunctions/MyStateMachine:*",
					},
				},
			},
		},
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine with logging: %v", err)
	}

	var unmarshaled StateMachine
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal state machine with logging: %v", err)
	}

	if unmarshaled.LoggingConfiguration == nil {
		t.Error("LoggingConfiguration should not be nil")
	}

	if len(unmarshaled.LoggingConfiguration.Destinations) != 1 {
		t.Errorf("Destinations length mismatch: got %d, want 1", len(unmarshaled.LoggingConfiguration.Destinations))
	}
}

func TestStateMachine_WithTracingConfiguration(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyStateMachine",
		TracingConfiguration: &TracingConfiguration{
			Enabled: true,
		},
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine with tracing: %v", err)
	}

	var unmarshaled StateMachine
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal state machine with tracing: %v", err)
	}

	if unmarshaled.TracingConfiguration == nil {
		t.Error("TracingConfiguration should not be nil")
	}
}

func TestStateMachine_WithEncryptionConfiguration(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyStateMachine",
		EncryptionConfiguration: &EncryptionConfiguration{
			Type:                         "CUSTOMER_MANAGED_KMS_KEY",
			KmsKeyId:                     "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
			KmsDataKeyReusePeriodSeconds: 300,
		},
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine with encryption: %v", err)
	}

	var unmarshaled StateMachine
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal state machine with encryption: %v", err)
	}

	if unmarshaled.EncryptionConfiguration == nil {
		t.Error("EncryptionConfiguration should not be nil")
	}

	if unmarshaled.EncryptionConfiguration.Type != sm.EncryptionConfiguration.Type {
		t.Errorf("Encryption Type mismatch: got %v, want %v", unmarshaled.EncryptionConfiguration.Type, sm.EncryptionConfiguration.Type)
	}
}

func TestStateMachine_ExpressType(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyExpressStateMachine",
		StateMachineType: "EXPRESS",
		DefinitionString: `{"StartAt": "HelloWorld", "States": {"HelloWorld": {"Type": "Pass", "End": true}}}`,
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal express state machine: %v", err)
	}

	var unmarshaled StateMachine
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal express state machine: %v", err)
	}

	if unmarshaled.StateMachineType != "EXPRESS" {
		t.Errorf("StateMachineType mismatch: got %v, want EXPRESS", unmarshaled.StateMachineType)
	}
}

func TestStateMachine_OmitEmpty(t *testing.T) {
	sm := StateMachine{
		StateMachineName: "MyStateMachine",
	}

	data, err := json.Marshal(sm)
	if err != nil {
		t.Fatalf("Failed to marshal state machine: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	omittedFields := []string{
		"Definition",
		"DefinitionString",
		"DefinitionS3Location",
		"DefinitionSubstitutions",
		"RoleArn",
		"StateMachineType",
		"LoggingConfiguration",
		"TracingConfiguration",
		"Tags",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestS3Location_JSONSerialization(t *testing.T) {
	loc := S3Location{
		Bucket:  "my-bucket",
		Key:     "my-key.json",
		Version: "abc123",
	}

	data, err := json.Marshal(loc)
	if err != nil {
		t.Fatalf("Failed to marshal S3Location: %v", err)
	}

	var unmarshaled S3Location
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal S3Location: %v", err)
	}

	if unmarshaled.Bucket != loc.Bucket {
		t.Errorf("Bucket mismatch: got %v, want %v", unmarshaled.Bucket, loc.Bucket)
	}
}
