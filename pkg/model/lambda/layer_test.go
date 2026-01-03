package lambda

import (
	"testing"
)

func TestNewLayerVersion(t *testing.T) {
	content := &LayerContent{S3Bucket: "my-bucket", S3Key: "my-key"}
	layer := NewLayerVersion(content)

	if layer.Content == nil {
		t.Error("expected Content to be set")
	}
	if layer.Content.S3Bucket != "my-bucket" {
		t.Errorf("expected S3Bucket 'my-bucket', got %v", layer.Content.S3Bucket)
	}
}

func TestNewLayerVersionFromS3(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key")

	if layer.Content == nil {
		t.Fatal("expected Content to be set")
	}
	if layer.Content.S3Bucket != "my-bucket" {
		t.Errorf("expected S3Bucket 'my-bucket', got %v", layer.Content.S3Bucket)
	}
	if layer.Content.S3Key != "my-key" {
		t.Errorf("expected S3Key 'my-key', got %v", layer.Content.S3Key)
	}
}

func TestNewLayerVersionFromS3WithVersion(t *testing.T) {
	layer := NewLayerVersionFromS3WithVersion("my-bucket", "my-key", "v1")

	if layer.Content.S3ObjectVersion != "v1" {
		t.Errorf("expected S3ObjectVersion 'v1', got %v", layer.Content.S3ObjectVersion)
	}
}

func TestLayerVersionWithLayerName(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key").WithLayerName("my-layer")

	if layer.LayerName != "my-layer" {
		t.Errorf("expected LayerName 'my-layer', got %s", layer.LayerName)
	}
}

func TestLayerVersionWithDescription(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key").WithDescription("My test layer")

	if layer.Description != "My test layer" {
		t.Errorf("expected Description 'My test layer', got %s", layer.Description)
	}
}

func TestLayerVersionWithLicenseInfo(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key").WithLicenseInfo("MIT")

	if layer.LicenseInfo != "MIT" {
		t.Errorf("expected LicenseInfo 'MIT', got %s", layer.LicenseInfo)
	}
}

func TestLayerVersionWithCompatibleRuntimes(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key").
		WithCompatibleRuntimes("nodejs18.x", "nodejs20.x")

	if len(layer.CompatibleRuntimes) != 2 {
		t.Errorf("expected 2 compatible runtimes, got %d", len(layer.CompatibleRuntimes))
	}
	if layer.CompatibleRuntimes[0] != "nodejs18.x" {
		t.Errorf("expected first runtime 'nodejs18.x', got %s", layer.CompatibleRuntimes[0])
	}
}

func TestLayerVersionWithCompatibleArchitectures(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key").
		WithCompatibleArchitectures("x86_64", "arm64")

	if len(layer.CompatibleArchitectures) != 2 {
		t.Errorf("expected 2 compatible architectures, got %d", len(layer.CompatibleArchitectures))
	}
}

func TestLayerVersionAddCompatibleRuntime(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key").
		AddCompatibleRuntime("python3.9").
		AddCompatibleRuntime("python3.10").
		AddCompatibleRuntime("python3.11")

	if len(layer.CompatibleRuntimes) != 3 {
		t.Errorf("expected 3 compatible runtimes, got %d", len(layer.CompatibleRuntimes))
	}
}

func TestLayerVersionAddCompatibleArchitecture(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key").
		AddCompatibleArchitecture("x86_64").
		AddCompatibleArchitecture("arm64")

	if len(layer.CompatibleArchitectures) != 2 {
		t.Errorf("expected 2 compatible architectures, got %d", len(layer.CompatibleArchitectures))
	}
}

func TestLayerVersionToCloudFormation_Minimal(t *testing.T) {
	layer := NewLayerVersionFromS3("my-bucket", "my-key")

	result := layer.ToCloudFormation()

	if result["Type"] != ResourceTypeLayerVersion {
		t.Errorf("expected Type %s, got %v", ResourceTypeLayerVersion, result["Type"])
	}

	props := result["Properties"].(map[string]interface{})
	content := props["Content"].(map[string]interface{})
	if content["S3Bucket"] != "my-bucket" {
		t.Errorf("expected S3Bucket in Content")
	}
}

func TestLayerVersionToCloudFormation_Full(t *testing.T) {
	layer := NewLayerVersionFromS3WithVersion("my-bucket", "my-key", "v1").
		WithLayerName("my-layer").
		WithDescription("Test layer").
		WithLicenseInfo("MIT").
		WithCompatibleRuntimes("nodejs18.x", "nodejs20.x").
		WithCompatibleArchitectures("x86_64", "arm64")

	result := layer.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["LayerName"] != "my-layer" {
		t.Errorf("expected LayerName 'my-layer', got %v", props["LayerName"])
	}
	if props["Description"] != "Test layer" {
		t.Errorf("expected Description 'Test layer', got %v", props["Description"])
	}
	if props["LicenseInfo"] != "MIT" {
		t.Errorf("expected LicenseInfo 'MIT', got %v", props["LicenseInfo"])
	}

	runtimes := props["CompatibleRuntimes"].([]string)
	if len(runtimes) != 2 {
		t.Errorf("expected 2 runtimes, got %d", len(runtimes))
	}

	archs := props["CompatibleArchitectures"].([]string)
	if len(archs) != 2 {
		t.Errorf("expected 2 architectures, got %d", len(archs))
	}

	content := props["Content"].(map[string]interface{})
	if content["S3ObjectVersion"] != "v1" {
		t.Errorf("expected S3ObjectVersion 'v1', got %v", content["S3ObjectVersion"])
	}
}

func TestLayerContentToMap(t *testing.T) {
	tests := []struct {
		name     string
		content  *LayerContent
		expected map[string]interface{}
	}{
		{
			name:    "S3 content",
			content: &LayerContent{S3Bucket: "bucket", S3Key: "key"},
			expected: map[string]interface{}{
				"S3Bucket": "bucket",
				"S3Key":    "key",
			},
		},
		{
			name:    "S3 content with version",
			content: &LayerContent{S3Bucket: "bucket", S3Key: "key", S3ObjectVersion: "v1"},
			expected: map[string]interface{}{
				"S3Bucket":        "bucket",
				"S3Key":           "key",
				"S3ObjectVersion": "v1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.content.toMap()
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("expected %s=%v, got %v", k, v, result[k])
				}
			}
		})
	}
}

// LayerVersionPermission tests

func TestNewLayerVersionPermission(t *testing.T) {
	perm := NewLayerVersionPermission("arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1", "123456789012")

	if perm.Action != "lambda:GetLayerVersion" {
		t.Errorf("expected Action 'lambda:GetLayerVersion', got %s", perm.Action)
	}
	if perm.LayerVersionArn != "arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1" {
		t.Errorf("unexpected LayerVersionArn: %v", perm.LayerVersionArn)
	}
	if perm.Principal != "123456789012" {
		t.Errorf("expected Principal '123456789012', got %s", perm.Principal)
	}
}

func TestNewLayerVersionPermissionPublic(t *testing.T) {
	perm := NewLayerVersionPermissionPublic("arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1")

	if perm.Principal != "*" {
		t.Errorf("expected Principal '*', got %s", perm.Principal)
	}
}

func TestNewLayerVersionPermissionOrg(t *testing.T) {
	perm := NewLayerVersionPermissionOrg("arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1", "o-1234567890")

	if perm.Principal != "*" {
		t.Errorf("expected Principal '*', got %s", perm.Principal)
	}
	if perm.OrganizationId != "o-1234567890" {
		t.Errorf("expected OrganizationId 'o-1234567890', got %s", perm.OrganizationId)
	}
}

func TestLayerVersionPermissionWithOrganizationId(t *testing.T) {
	perm := NewLayerVersionPermissionPublic("arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1").
		WithOrganizationId("o-9876543210")

	if perm.OrganizationId != "o-9876543210" {
		t.Errorf("expected OrganizationId 'o-9876543210', got %s", perm.OrganizationId)
	}
}

func TestLayerVersionPermissionToCloudFormation_Minimal(t *testing.T) {
	perm := NewLayerVersionPermission("arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1", "123456789012")

	result := perm.ToCloudFormation()

	if result["Type"] != ResourceTypeLayerVersionPermission {
		t.Errorf("expected Type %s, got %v", ResourceTypeLayerVersionPermission, result["Type"])
	}

	props := result["Properties"].(map[string]interface{})
	if props["Action"] != "lambda:GetLayerVersion" {
		t.Errorf("expected Action in properties")
	}
	if props["LayerVersionArn"] != "arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1" {
		t.Errorf("expected LayerVersionArn in properties")
	}
	if props["Principal"] != "123456789012" {
		t.Errorf("expected Principal in properties")
	}
}

func TestLayerVersionPermissionToCloudFormation_WithOrg(t *testing.T) {
	perm := NewLayerVersionPermissionOrg("arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1", "o-1234567890")

	result := perm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["OrganizationId"] != "o-1234567890" {
		t.Errorf("expected OrganizationId in properties")
	}
}

func TestLayerVersionPermissionWithIntrinsicArn(t *testing.T) {
	arnRef := map[string]interface{}{"Ref": "MyLayerVersion"}
	perm := NewLayerVersionPermission(arnRef, "123456789012")

	result := perm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	arn := props["LayerVersionArn"].(map[string]interface{})
	if arn["Ref"] != "MyLayerVersion" {
		t.Errorf("expected Ref to MyLayerVersion, got %v", arn)
	}
}
