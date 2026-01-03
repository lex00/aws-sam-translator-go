package sam

import (
	"testing"
)

func TestLayerVersionTransformer_Transform_Minimal(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	lv := &LayerVersion{
		ContentUri: "s3://sam-demo-bucket/layer.zip",
	}

	resources, newLogicalID, err := transformer.Transform("MinimalLayer", lv)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Check that logical ID has hash suffix
	if len(newLogicalID) <= len("MinimalLayer") {
		t.Errorf("expected logical ID with hash suffix, got %q", newLogicalID)
	}

	// Check resource exists with new logical ID
	resource, ok := resources[newLogicalID].(map[string]interface{})
	if !ok {
		t.Fatalf("resource not found with logical ID %q", newLogicalID)
	}

	// Check type
	if resource["Type"] != "AWS::Lambda::LayerVersion" {
		t.Errorf("expected Type 'AWS::Lambda::LayerVersion', got %v", resource["Type"])
	}

	// Check DeletionPolicy - default is Retain
	if resource["DeletionPolicy"] != "Retain" {
		t.Errorf("expected DeletionPolicy 'Retain', got %v", resource["DeletionPolicy"])
	}

	// Check properties
	props := resource["Properties"].(map[string]interface{})

	// Check Content
	content := props["Content"].(map[string]interface{})
	if content["S3Bucket"] != "sam-demo-bucket" {
		t.Errorf("expected S3Bucket 'sam-demo-bucket', got %v", content["S3Bucket"])
	}
	if content["S3Key"] != "layer.zip" {
		t.Errorf("expected S3Key 'layer.zip', got %v", content["S3Key"])
	}

	// LayerName should default to original logical ID
	if props["LayerName"] != "MinimalLayer" {
		t.Errorf("expected LayerName 'MinimalLayer', got %v", props["LayerName"])
	}
}

func TestLayerVersionTransformer_Transform_WithAllProperties(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	lv := &LayerVersion{
		ContentUri:         "s3://sam-demo-bucket/layer.zip",
		LayerName:          "MyAwesomeLayer",
		Description:        "Starter Lambda Layer",
		CompatibleRuntimes: []string{"python3.9"},
		LicenseInfo:        "License information",
		RetentionPolicy:    "Retain",
	}

	resources, newLogicalID, err := transformer.Transform("CompleteLayer", lv)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources[newLogicalID].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	// Check LayerName
	if props["LayerName"] != "MyAwesomeLayer" {
		t.Errorf("expected LayerName 'MyAwesomeLayer', got %v", props["LayerName"])
	}

	// Check Description
	if props["Description"] != "Starter Lambda Layer" {
		t.Errorf("expected Description 'Starter Lambda Layer', got %v", props["Description"])
	}

	// Check CompatibleRuntimes
	runtimes := props["CompatibleRuntimes"].([]string)
	if len(runtimes) != 1 || runtimes[0] != "python3.9" {
		t.Errorf("expected CompatibleRuntimes ['python3.9'], got %v", runtimes)
	}

	// Check LicenseInfo
	if props["LicenseInfo"] != "License information" {
		t.Errorf("expected LicenseInfo 'License information', got %v", props["LicenseInfo"])
	}

	// Check DeletionPolicy
	if resource["DeletionPolicy"] != "Retain" {
		t.Errorf("expected DeletionPolicy 'Retain', got %v", resource["DeletionPolicy"])
	}
}

func TestLayerVersionTransformer_Transform_WithContentUriObject(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	lv := &LayerVersion{
		ContentUri: map[string]interface{}{
			"Bucket":  "somebucket",
			"Key":     "somekey",
			"Version": "v1",
		},
		RetentionPolicy: "Delete",
	}

	resources, newLogicalID, err := transformer.Transform("LayerWithContentUriObject", lv)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources[newLogicalID].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	content := props["Content"].(map[string]interface{})
	if content["S3Bucket"] != "somebucket" {
		t.Errorf("expected S3Bucket 'somebucket', got %v", content["S3Bucket"])
	}
	if content["S3Key"] != "somekey" {
		t.Errorf("expected S3Key 'somekey', got %v", content["S3Key"])
	}
	if content["S3ObjectVersion"] != "v1" {
		t.Errorf("expected S3ObjectVersion 'v1', got %v", content["S3ObjectVersion"])
	}

	// Check DeletionPolicy
	if resource["DeletionPolicy"] != "Delete" {
		t.Errorf("expected DeletionPolicy 'Delete', got %v", resource["DeletionPolicy"])
	}
}

func TestLayerVersionTransformer_Transform_CaseInsensitiveRetentionPolicy(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	lv := &LayerVersion{
		ContentUri:      "s3://sam-demo-bucket/layer.zip",
		RetentionPolicy: "DeleTe",
	}

	resources, newLogicalID, err := transformer.Transform("LayerWithCaseInsensitiveRetentionPolicy", lv)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources[newLogicalID].(map[string]interface{})

	// Case-insensitive "DeleTe" should map to "Delete"
	if resource["DeletionPolicy"] != "Delete" {
		t.Errorf("expected DeletionPolicy 'Delete', got %v", resource["DeletionPolicy"])
	}
}

func TestLayerVersionTransformer_Transform_WithArchitectures(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	lv := &LayerVersion{
		ContentUri:              "s3://sam-demo-bucket/layer.zip",
		CompatibleArchitectures: []string{"x86_64", "arm64"},
	}

	resources, newLogicalID, err := transformer.Transform("LayerWithArchitectures", lv)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources[newLogicalID].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	archs := props["CompatibleArchitectures"].([]string)
	if len(archs) != 2 {
		t.Errorf("expected 2 architectures, got %d", len(archs))
	}
	if archs[0] != "x86_64" || archs[1] != "arm64" {
		t.Errorf("expected architectures [x86_64, arm64], got %v", archs)
	}
}

func TestLayerVersionTransformer_parseS3Uri(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	tests := []struct {
		name       string
		uri        string
		wantBucket string
		wantKey    string
		wantErr    bool
	}{
		{
			name:       "valid S3 URI",
			uri:        "s3://my-bucket/my-key.zip",
			wantBucket: "my-bucket",
			wantKey:    "my-key.zip",
		},
		{
			name:       "S3 URI with path",
			uri:        "s3://bucket/path/to/file.zip",
			wantBucket: "bucket",
			wantKey:    "path/to/file.zip",
		},
		{
			name:    "invalid - not S3 URI",
			uri:     "http://example.com/file.zip",
			wantErr: true,
		},
		{
			name:    "invalid - no key",
			uri:     "s3://bucket-only",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transformer.parseS3Uri(tt.uri)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result["S3Bucket"] != tt.wantBucket {
				t.Errorf("expected S3Bucket %q, got %v", tt.wantBucket, result["S3Bucket"])
			}
			if result["S3Key"] != tt.wantKey {
				t.Errorf("expected S3Key %q, got %v", tt.wantKey, result["S3Key"])
			}
		})
	}
}

func TestLayerVersionTransformer_mapRetentionPolicy(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil defaults to Retain", nil, "Retain"},
		{"Retain", "Retain", "Retain"},
		{"retain lowercase", "retain", "Retain"},
		{"Delete", "Delete", "Delete"},
		{"delete lowercase", "delete", "Delete"},
		{"DeleTe mixed case", "DeleTe", "Delete"},
		{"unknown defaults to Retain", "Unknown", "Retain"},
		{"intrinsic function defaults to Retain", map[string]interface{}{"Ref": "DeleteParam"}, "Retain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.mapRetentionPolicy(tt.input)
			if result != tt.expected {
				t.Errorf("mapRetentionPolicy(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLayerVersionTransformer_generateLogicalID(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	content := map[string]interface{}{
		"S3Bucket": "sam-demo-bucket",
		"S3Key":    "layer.zip",
	}

	id1 := transformer.generateLogicalID("MyLayer", content)
	id2 := transformer.generateLogicalID("MyLayer", content)

	// Same inputs should produce same output (deterministic)
	if id1 != id2 {
		t.Errorf("expected deterministic IDs, got %q and %q", id1, id2)
	}

	// ID should start with base and have hash suffix
	if len(id1) <= len("MyLayer") {
		t.Errorf("expected ID longer than base, got %q", id1)
	}

	// Different content should produce different ID
	content2 := map[string]interface{}{
		"S3Bucket": "different-bucket",
		"S3Key":    "layer.zip",
	}
	id3 := transformer.generateLogicalID("MyLayer", content2)
	if id1 == id3 {
		t.Errorf("different content should produce different IDs: %q == %q", id1, id3)
	}
}

func TestLayerVersionTransformer_Transform_Error_MissingContentUri(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	lv := &LayerVersion{
		// ContentUri is nil/missing
	}

	_, _, err := transformer.Transform("InvalidLayer", lv)
	if err == nil {
		t.Error("expected error for missing ContentUri, got nil")
	}
}

func TestLayerVersionTransformer_Transform_Error_InvalidContentUriObject(t *testing.T) {
	transformer := NewLayerVersionTransformer()

	// Missing Bucket
	lv := &LayerVersion{
		ContentUri: map[string]interface{}{
			"Key": "somekey",
		},
	}

	_, _, err := transformer.Transform("InvalidLayer", lv)
	if err == nil {
		t.Error("expected error for missing Bucket, got nil")
	}

	// Missing Key
	lv2 := &LayerVersion{
		ContentUri: map[string]interface{}{
			"Bucket": "somebucket",
		},
	}

	_, _, err = transformer.Transform("InvalidLayer2", lv2)
	if err == nil {
		t.Error("expected error for missing Key, got nil")
	}
}
