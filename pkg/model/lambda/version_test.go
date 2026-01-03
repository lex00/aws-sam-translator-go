package lambda

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
	version := NewVersion("my-function")

	if version.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", version.FunctionName)
	}
}

func TestNewVersionWithDescription(t *testing.T) {
	version := NewVersionWithDescription("my-function", "Initial version")

	if version.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", version.FunctionName)
	}
	if version.Description != "Initial version" {
		t.Errorf("expected Description 'Initial version', got %s", version.Description)
	}
}

func TestVersionWithCodeSha256(t *testing.T) {
	version := NewVersion("my-function").WithCodeSha256("abc123def456")

	if version.CodeSha256 != "abc123def456" {
		t.Errorf("expected CodeSha256 'abc123def456', got %s", version.CodeSha256)
	}
}

func TestVersionWithProvisionedConcurrency(t *testing.T) {
	version := NewVersion("my-function").WithProvisionedConcurrency(100)

	if version.ProvisionedConcurrencyConfig == nil {
		t.Fatal("expected ProvisionedConcurrencyConfig to be set")
	}
	if version.ProvisionedConcurrencyConfig.ProvisionedConcurrentExecutions != 100 {
		t.Errorf("expected ProvisionedConcurrentExecutions 100, got %d",
			version.ProvisionedConcurrencyConfig.ProvisionedConcurrentExecutions)
	}
}

func TestVersionWithRuntimePolicy(t *testing.T) {
	runtimeArn := "arn:aws:lambda:us-east-1::runtime:python3.9"
	version := NewVersion("my-function").WithRuntimePolicy("Manual", runtimeArn)

	if version.RuntimePolicy == nil {
		t.Fatal("expected RuntimePolicy to be set")
	}
	if version.RuntimePolicy.UpdateRuntimeOn != "Manual" {
		t.Errorf("expected UpdateRuntimeOn 'Manual', got %s", version.RuntimePolicy.UpdateRuntimeOn)
	}
	if version.RuntimePolicy.RuntimeVersionArn != runtimeArn {
		t.Errorf("expected RuntimeVersionArn %s, got %v", runtimeArn, version.RuntimePolicy.RuntimeVersionArn)
	}
}

func TestVersionToCloudFormation_Minimal(t *testing.T) {
	version := NewVersion("my-function")

	result := version.ToCloudFormation()

	if result["Type"] != ResourceTypeVersion {
		t.Errorf("expected Type %s, got %v", ResourceTypeVersion, result["Type"])
	}

	props := result["Properties"].(map[string]interface{})
	if props["FunctionName"] != "my-function" {
		t.Errorf("expected FunctionName in properties")
	}
}

func TestVersionToCloudFormation_Full(t *testing.T) {
	version := NewVersion("my-function").
		WithCodeSha256("abc123").
		WithProvisionedConcurrency(50).
		WithRuntimePolicy("Auto", nil)
	version.Description = "Production version"
	version.Policy = map[string]interface{}{"Version": "2012-10-17"}

	result := version.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["CodeSha256"] != "abc123" {
		t.Errorf("expected CodeSha256 'abc123', got %v", props["CodeSha256"])
	}

	if props["Description"] != "Production version" {
		t.Errorf("expected Description 'Production version', got %v", props["Description"])
	}

	if props["Policy"] == nil {
		t.Error("expected Policy in properties")
	}

	provConfig := props["ProvisionedConcurrencyConfig"].(map[string]interface{})
	if provConfig["ProvisionedConcurrentExecutions"] != 50 {
		t.Errorf("expected ProvisionedConcurrentExecutions 50, got %v", provConfig["ProvisionedConcurrentExecutions"])
	}

	runtimePolicy := props["RuntimePolicy"].(map[string]interface{})
	if runtimePolicy["UpdateRuntimeOn"] != "Auto" {
		t.Errorf("expected UpdateRuntimeOn 'Auto', got %v", runtimePolicy["UpdateRuntimeOn"])
	}
}

func TestVersionWithIntrinsicFunctionName(t *testing.T) {
	// Test with intrinsic function reference
	fnRef := map[string]interface{}{"Ref": "MyLambdaFunction"}
	version := NewVersion(fnRef)

	result := version.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	fnName := props["FunctionName"].(map[string]interface{})
	if fnName["Ref"] != "MyLambdaFunction" {
		t.Errorf("expected Ref to MyLambdaFunction, got %v", fnName)
	}
}

func TestProvisionedConcurrencyConfigToMap(t *testing.T) {
	config := &ProvisionedConcurrencyConfig{
		ProvisionedConcurrentExecutions: 25,
	}

	result := config.toMap()
	if result["ProvisionedConcurrentExecutions"] != 25 {
		t.Errorf("expected ProvisionedConcurrentExecutions 25, got %v", result["ProvisionedConcurrentExecutions"])
	}
}

func TestRuntimePolicyToMap(t *testing.T) {
	tests := []struct {
		name   string
		policy *RuntimePolicy
		check  func(map[string]interface{}) bool
	}{
		{
			name:   "Auto update",
			policy: &RuntimePolicy{UpdateRuntimeOn: "Auto"},
			check: func(m map[string]interface{}) bool {
				return m["UpdateRuntimeOn"] == "Auto" && m["RuntimeVersionArn"] == nil
			},
		},
		{
			name:   "Manual with ARN",
			policy: &RuntimePolicy{UpdateRuntimeOn: "Manual", RuntimeVersionArn: "arn:aws:lambda:us-east-1::runtime:python3.9"},
			check: func(m map[string]interface{}) bool {
				return m["UpdateRuntimeOn"] == "Manual" && m["RuntimeVersionArn"] == "arn:aws:lambda:us-east-1::runtime:python3.9"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.policy.toMap()
			if !tt.check(result) {
				t.Errorf("unexpected result: %v", result)
			}
		})
	}
}
