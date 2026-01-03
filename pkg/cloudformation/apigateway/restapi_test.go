package apigateway

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRestApi_JSONSerialization(t *testing.T) {
	restApi := RestApi{
		Name:             "MyRestApi",
		Description:      "My REST API description",
		ApiKeySourceType: "HEADER",
		BinaryMediaTypes: []interface{}{"image/png", "application/octet-stream"},
		EndpointConfiguration: &EndpointConfiguration{
			Types: []interface{}{"REGIONAL"},
		},
		DisableExecuteApiEndpoint: false,
		FailOnWarnings:            true,
		MinimumCompressionSize:    1024,
		Policy:                    map[string]interface{}{"Version": "2012-10-17"},
		Tags: []Tag{
			{Key: "Environment", Value: "Production"},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(restApi)
	if err != nil {
		t.Fatalf("Failed to marshal RestApi to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled RestApi
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RestApi from JSON: %v", err)
	}

	if unmarshaled.Name != restApi.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, restApi.Name)
	}

	if unmarshaled.Description != restApi.Description {
		t.Errorf("Description mismatch: got %v, want %v", unmarshaled.Description, restApi.Description)
	}

	if len(unmarshaled.BinaryMediaTypes) != len(restApi.BinaryMediaTypes) {
		t.Errorf("BinaryMediaTypes length mismatch: got %d, want %d",
			len(unmarshaled.BinaryMediaTypes), len(restApi.BinaryMediaTypes))
	}

	if unmarshaled.EndpointConfiguration == nil {
		t.Error("EndpointConfiguration should not be nil")
	}
}

func TestRestApi_YAMLSerialization(t *testing.T) {
	restApi := RestApi{
		Name:        "MyRestApi",
		Description: "Test API",
		BodyS3Location: &S3Location{
			Bucket:  "my-bucket",
			Key:     "api-spec.yaml",
			Version: "v1",
		},
	}

	// Test YAML marshaling
	data, err := yaml.Marshal(restApi)
	if err != nil {
		t.Fatalf("Failed to marshal RestApi to YAML: %v", err)
	}

	// Test YAML unmarshaling
	var unmarshaled RestApi
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RestApi from YAML: %v", err)
	}

	if unmarshaled.Name != restApi.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, restApi.Name)
	}

	if unmarshaled.BodyS3Location == nil {
		t.Error("BodyS3Location should not be nil")
	} else if unmarshaled.BodyS3Location.Bucket != restApi.BodyS3Location.Bucket {
		t.Errorf("BodyS3Location.Bucket mismatch: got %v, want %v",
			unmarshaled.BodyS3Location.Bucket, restApi.BodyS3Location.Bucket)
	}
}

func TestRestApi_WithIntrinsicFunctions(t *testing.T) {
	// Test with intrinsic functions (e.g., Ref, Fn::Sub)
	restApi := RestApi{
		Name:        map[string]interface{}{"Fn::Sub": "${AWS::StackName}-api"},
		Description: map[string]interface{}{"Ref": "ApiDescription"},
		EndpointConfiguration: &EndpointConfiguration{
			Types: []interface{}{map[string]interface{}{"Ref": "EndpointType"}},
			VpcEndpointIds: []interface{}{
				map[string]interface{}{"Ref": "VpcEndpoint1"},
				map[string]interface{}{"Ref": "VpcEndpoint2"},
			},
		},
	}

	data, err := json.Marshal(restApi)
	if err != nil {
		t.Fatalf("Failed to marshal RestApi with intrinsic functions: %v", err)
	}

	var unmarshaled RestApi
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RestApi with intrinsic functions: %v", err)
	}

	// Verify the intrinsic function structure is preserved
	nameMap, ok := unmarshaled.Name.(map[string]interface{})
	if !ok {
		t.Error("Name should be a map for intrinsic function")
	} else if _, exists := nameMap["Fn::Sub"]; !exists {
		t.Error("Name should contain Fn::Sub intrinsic function")
	}
}

func TestRestApi_OmitEmpty(t *testing.T) {
	restApi := RestApi{
		Name: "MinimalApi",
	}

	data, err := json.Marshal(restApi)
	if err != nil {
		t.Fatalf("Failed to marshal RestApi: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	// These should be omitted because they're empty
	emptyFields := []string{"Description", "BinaryMediaTypes", "Body", "BodyS3Location",
		"EndpointConfiguration", "Policy", "Tags", "Parameters"}
	for _, field := range emptyFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestS3Location_JSONSerialization(t *testing.T) {
	s3Location := S3Location{
		Bucket:  "my-bucket",
		Key:     "openapi.yaml",
		ETag:    "abc123",
		Version: "v1",
	}

	data, err := json.Marshal(s3Location)
	if err != nil {
		t.Fatalf("Failed to marshal S3Location: %v", err)
	}

	var unmarshaled S3Location
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal S3Location: %v", err)
	}

	if unmarshaled.Bucket != s3Location.Bucket {
		t.Errorf("Bucket mismatch: got %v, want %v", unmarshaled.Bucket, s3Location.Bucket)
	}

	if unmarshaled.Key != s3Location.Key {
		t.Errorf("Key mismatch: got %v, want %v", unmarshaled.Key, s3Location.Key)
	}
}

func TestEndpointConfiguration_JSONSerialization(t *testing.T) {
	config := EndpointConfiguration{
		Types:          []interface{}{"PRIVATE"},
		VpcEndpointIds: []interface{}{"vpce-12345", "vpce-67890"},
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal EndpointConfiguration: %v", err)
	}

	var unmarshaled EndpointConfiguration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal EndpointConfiguration: %v", err)
	}

	if len(unmarshaled.Types) != len(config.Types) {
		t.Errorf("Types length mismatch: got %d, want %d",
			len(unmarshaled.Types), len(config.Types))
	}

	if len(unmarshaled.VpcEndpointIds) != len(config.VpcEndpointIds) {
		t.Errorf("VpcEndpointIds length mismatch: got %d, want %d",
			len(unmarshaled.VpcEndpointIds), len(config.VpcEndpointIds))
	}
}

func TestTag_JSONSerialization(t *testing.T) {
	tag := Tag{
		Key:   "Environment",
		Value: "Production",
	}

	data, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal Tag: %v", err)
	}

	var unmarshaled Tag
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Tag: %v", err)
	}

	if unmarshaled.Key != tag.Key {
		t.Errorf("Key mismatch: got %v, want %v", unmarshaled.Key, tag.Key)
	}

	if unmarshaled.Value != tag.Value {
		t.Errorf("Value mismatch: got %v, want %v", unmarshaled.Value, tag.Value)
	}
}
