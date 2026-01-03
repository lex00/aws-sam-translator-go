package types

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTemplate_JSONSerialization(t *testing.T) {
	template := Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Description:              "Test template",
		Parameters: map[string]Parameter{
			"Stage": {
				Type:    "String",
				Default: "dev",
			},
		},
		Resources: map[string]Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Runtime": "go1.x",
					"Handler": "main",
				},
			},
		},
		Outputs: map[string]Output{
			"FunctionArn": {
				Value:       map[string]interface{}{"Fn::GetAtt": []string{"MyFunction", "Arn"}},
				Description: "Function ARN",
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal template to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Template
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal template from JSON: %v", err)
	}

	// Verify key fields
	if unmarshaled.AWSTemplateFormatVersion != template.AWSTemplateFormatVersion {
		t.Errorf("AWSTemplateFormatVersion mismatch: got %q, want %q",
			unmarshaled.AWSTemplateFormatVersion, template.AWSTemplateFormatVersion)
	}

	if unmarshaled.Description != template.Description {
		t.Errorf("Description mismatch: got %q, want %q",
			unmarshaled.Description, template.Description)
	}

	if len(unmarshaled.Resources) != len(template.Resources) {
		t.Errorf("Resources count mismatch: got %d, want %d",
			len(unmarshaled.Resources), len(template.Resources))
	}
}

func TestTemplate_YAMLSerialization(t *testing.T) {
	template := Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Description:              "Test template",
		Resources: map[string]Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Runtime": "go1.x",
				},
			},
		},
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal template to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled Template
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal template from YAML: %v", err)
	}

	if unmarshaled.AWSTemplateFormatVersion != template.AWSTemplateFormatVersion {
		t.Errorf("AWSTemplateFormatVersion mismatch: got %q, want %q",
			unmarshaled.AWSTemplateFormatVersion, template.AWSTemplateFormatVersion)
	}
}

func TestParameter_JSONSerialization(t *testing.T) {
	param := Parameter{
		Type:                  "String",
		Default:               "production",
		Description:           "Deployment stage",
		AllowedValues:         []string{"dev", "staging", "production"},
		AllowedPattern:        "^[a-z]+$",
		ConstraintDescription: "Must be lowercase letters",
		MaxLength:             50,
		MinLength:             1,
		NoEcho:                true,
	}

	data, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("Failed to marshal parameter: %v", err)
	}

	var unmarshaled Parameter
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal parameter: %v", err)
	}

	if unmarshaled.Type != param.Type {
		t.Errorf("Type mismatch: got %q, want %q", unmarshaled.Type, param.Type)
	}

	if len(unmarshaled.AllowedValues) != len(param.AllowedValues) {
		t.Errorf("AllowedValues length mismatch: got %d, want %d",
			len(unmarshaled.AllowedValues), len(param.AllowedValues))
	}

	if unmarshaled.NoEcho != param.NoEcho {
		t.Errorf("NoEcho mismatch: got %v, want %v", unmarshaled.NoEcho, param.NoEcho)
	}
}

func TestParameter_NumericConstraints(t *testing.T) {
	param := Parameter{
		Type:     "Number",
		Default:  10,
		MinValue: 1,
		MaxValue: 100,
	}

	data, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("Failed to marshal parameter: %v", err)
	}

	var unmarshaled Parameter
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal parameter: %v", err)
	}

	if unmarshaled.MinValue != param.MinValue {
		t.Errorf("MinValue mismatch: got %v, want %v", unmarshaled.MinValue, param.MinValue)
	}

	if unmarshaled.MaxValue != param.MaxValue {
		t.Errorf("MaxValue mismatch: got %v, want %v", unmarshaled.MaxValue, param.MaxValue)
	}
}

func TestResource_JSONSerialization(t *testing.T) {
	resource := Resource{
		Type: "AWS::Lambda::Function",
		Properties: map[string]interface{}{
			"FunctionName": "my-function",
			"Runtime":      "go1.x",
			"Handler":      "main",
			"Code": map[string]interface{}{
				"S3Bucket": "my-bucket",
				"S3Key":    "code.zip",
			},
		},
		Metadata: map[string]interface{}{
			"aws:sam:build": true,
		},
		DependsOn:      []interface{}{"MyRole"},
		Condition:      "IsProduction",
		DeletionPolicy: "Retain",
	}

	data, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal resource: %v", err)
	}

	var unmarshaled Resource
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal resource: %v", err)
	}

	if unmarshaled.Type != resource.Type {
		t.Errorf("Type mismatch: got %q, want %q", unmarshaled.Type, resource.Type)
	}

	if unmarshaled.Condition != resource.Condition {
		t.Errorf("Condition mismatch: got %q, want %q", unmarshaled.Condition, resource.Condition)
	}

	if unmarshaled.DeletionPolicy != resource.DeletionPolicy {
		t.Errorf("DeletionPolicy mismatch: got %q, want %q",
			unmarshaled.DeletionPolicy, resource.DeletionPolicy)
	}
}

func TestOutput_JSONSerialization(t *testing.T) {
	output := Output{
		Description: "The function ARN",
		Value:       map[string]interface{}{"Fn::GetAtt": []string{"MyFunction", "Arn"}},
		Export: &Export{
			Name: map[string]interface{}{"Fn::Sub": "${AWS::StackName}-FunctionArn"},
		},
		Condition: "IsProduction",
	}

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal output: %v", err)
	}

	var unmarshaled Output
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if unmarshaled.Description != output.Description {
		t.Errorf("Description mismatch: got %q, want %q",
			unmarshaled.Description, output.Description)
	}

	if unmarshaled.Condition != output.Condition {
		t.Errorf("Condition mismatch: got %q, want %q",
			unmarshaled.Condition, output.Condition)
	}

	if unmarshaled.Export == nil {
		t.Error("Export should not be nil")
	}
}

func TestExport_JSONSerialization(t *testing.T) {
	export := Export{
		Name: "my-export-name",
	}

	data, err := json.Marshal(export)
	if err != nil {
		t.Fatalf("Failed to marshal export: %v", err)
	}

	var unmarshaled Export
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal export: %v", err)
	}

	if unmarshaled.Name != export.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, export.Name)
	}
}

func TestTemplate_OmitEmpty(t *testing.T) {
	// Test that empty fields are omitted in JSON
	template := Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Resources: map[string]Resource{
			"MyResource": {
				Type: "AWS::S3::Bucket",
			},
		},
	}

	data, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Failed to marshal template: %v", err)
	}

	// Check that empty fields are not present
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	// These should be omitted because they're empty
	emptyFields := []string{"Description", "Parameters", "Mappings", "Conditions", "Outputs", "Globals", "Metadata"}
	for _, field := range emptyFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestResource_YAMLSerialization(t *testing.T) {
	resource := Resource{
		Type: "AWS::Lambda::Function",
		Properties: map[string]interface{}{
			"FunctionName": "test-function",
		},
		UpdatePolicy: map[string]interface{}{
			"AutoScalingRollingUpdate": map[string]interface{}{
				"MinInstancesInService": 1,
			},
		},
	}

	data, err := yaml.Marshal(resource)
	if err != nil {
		t.Fatalf("Failed to marshal resource to YAML: %v", err)
	}

	var unmarshaled Resource
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal resource from YAML: %v", err)
	}

	if unmarshaled.Type != resource.Type {
		t.Errorf("Type mismatch: got %q, want %q", unmarshaled.Type, resource.Type)
	}

	if unmarshaled.UpdatePolicy == nil {
		t.Error("UpdatePolicy should not be nil")
	}
}
