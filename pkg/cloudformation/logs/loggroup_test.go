package logs

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLogGroup_JSONSerialization(t *testing.T) {
	logGroup := LogGroup{
		LogGroupName:    "/aws/lambda/MyFunction",
		RetentionInDays: 30,
		Tags: []Tag{
			{Key: "Environment", Value: "Production"},
		},
	}

	data, err := json.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal log group to JSON: %v", err)
	}

	var unmarshaled LogGroup
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal log group from JSON: %v", err)
	}

	if unmarshaled.LogGroupName != logGroup.LogGroupName {
		t.Errorf("LogGroupName mismatch: got %v, want %v", unmarshaled.LogGroupName, logGroup.LogGroupName)
	}

	// Compare as strings since JSON unmarshaling converts numbers to float64
	if fmt.Sprintf("%v", unmarshaled.RetentionInDays) != fmt.Sprintf("%v", logGroup.RetentionInDays) {
		t.Errorf("RetentionInDays mismatch: got %v, want %v", unmarshaled.RetentionInDays, logGroup.RetentionInDays)
	}
}

func TestLogGroup_YAMLSerialization(t *testing.T) {
	logGroup := LogGroup{
		LogGroupName:    "/aws/lambda/MyFunction",
		RetentionInDays: 14,
	}

	data, err := yaml.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal log group to YAML: %v", err)
	}

	var unmarshaled LogGroup
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal log group from YAML: %v", err)
	}

	if unmarshaled.LogGroupName != logGroup.LogGroupName {
		t.Errorf("LogGroupName mismatch: got %v, want %v", unmarshaled.LogGroupName, logGroup.LogGroupName)
	}
}

func TestLogGroup_WithIntrinsicFunctions(t *testing.T) {
	logGroup := LogGroup{
		LogGroupName: map[string]interface{}{
			"Fn::Sub": "/aws/lambda/${AWS::StackName}-function",
		},
		KmsKeyId: map[string]interface{}{
			"Fn::GetAtt": []string{"MyKey", "Arn"},
		},
	}

	data, err := json.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal log group with intrinsics: %v", err)
	}

	var unmarshaled LogGroup
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal log group with intrinsics: %v", err)
	}
}

func TestLogGroup_WithEncryption(t *testing.T) {
	logGroup := LogGroup{
		LogGroupName: "/aws/lambda/MyFunction",
		KmsKeyId:     "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
	}

	data, err := json.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal encrypted log group: %v", err)
	}

	var unmarshaled LogGroup
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal encrypted log group: %v", err)
	}

	if unmarshaled.KmsKeyId != logGroup.KmsKeyId {
		t.Errorf("KmsKeyId mismatch: got %v, want %v", unmarshaled.KmsKeyId, logGroup.KmsKeyId)
	}
}

func TestLogGroup_WithDataProtectionPolicy(t *testing.T) {
	logGroup := LogGroup{
		LogGroupName: "/aws/lambda/MyFunction",
		DataProtectionPolicy: map[string]interface{}{
			"Name":        "data-protection-policy",
			"Description": "Protect sensitive data",
			"Version":     "2021-06-01",
			"Statement": []map[string]interface{}{
				{
					"Sid":            "audit-policy",
					"DataIdentifier": []string{"arn:aws:dataprotection::aws:data-identifier/CreditCardNumber"},
					"Operation": map[string]interface{}{
						"Audit": map[string]interface{}{
							"FindingsDestination": map[string]interface{}{
								"CloudWatchLogs": map[string]interface{}{
									"LogGroup": "/aws/logs/audit",
								},
							},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal log group with data protection policy: %v", err)
	}

	var unmarshaled LogGroup
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal log group with data protection policy: %v", err)
	}

	if unmarshaled.DataProtectionPolicy == nil {
		t.Error("DataProtectionPolicy should not be nil")
	}
}

func TestLogGroup_WithLogGroupClass(t *testing.T) {
	logGroup := LogGroup{
		LogGroupName:  "/aws/lambda/MyFunction",
		LogGroupClass: "INFREQUENT_ACCESS",
	}

	data, err := json.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal log group with class: %v", err)
	}

	var unmarshaled LogGroup
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal log group with class: %v", err)
	}

	if unmarshaled.LogGroupClass != logGroup.LogGroupClass {
		t.Errorf("LogGroupClass mismatch: got %v, want %v", unmarshaled.LogGroupClass, logGroup.LogGroupClass)
	}
}

func TestLogGroup_AllRetentionValues(t *testing.T) {
	// Test that we can use various valid retention values
	retentionValues := []int{1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653}

	for _, retention := range retentionValues {
		logGroup := LogGroup{
			LogGroupName:    "/aws/test/retention",
			RetentionInDays: retention,
		}

		data, err := json.Marshal(logGroup)
		if err != nil {
			t.Fatalf("Failed to marshal log group with retention %d: %v", retention, err)
		}

		var unmarshaled LogGroup
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal log group with retention %d: %v", retention, err)
		}
	}
}

func TestLogGroup_OmitEmpty(t *testing.T) {
	logGroup := LogGroup{
		LogGroupName: "/aws/lambda/MyFunction",
	}

	data, err := json.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal log group: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	omittedFields := []string{
		"RetentionInDays",
		"KmsKeyId",
		"DataProtectionPolicy",
		"LogGroupClass",
		"Tags",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestLogGroup_MinimalConfiguration(t *testing.T) {
	// Test that a log group with just a name serializes correctly
	logGroup := LogGroup{
		LogGroupName: "/my/log/group",
	}

	data, err := json.Marshal(logGroup)
	if err != nil {
		t.Fatalf("Failed to marshal minimal log group: %v", err)
	}

	var unmarshaled LogGroup
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal minimal log group: %v", err)
	}

	if unmarshaled.LogGroupName != logGroup.LogGroupName {
		t.Errorf("LogGroupName mismatch: got %v, want %v", unmarshaled.LogGroupName, logGroup.LogGroupName)
	}
}
