package intrinsics

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestNewResolveContext(t *testing.T) {
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"Environment": {
				Type:    "String",
				Default: "dev",
			},
			"BucketName": {
				Type: "String",
				// No default
			},
		},
	}

	ctx := NewResolveContext(template)

	// Check pseudo-parameters are initialized
	if ctx.PseudoParameters["AWS::Region"] != "us-east-1" {
		t.Errorf("expected AWS::Region = us-east-1, got %s", ctx.PseudoParameters["AWS::Region"])
	}
	if ctx.PseudoParameters["AWS::AccountId"] != "123456789012" {
		t.Errorf("expected AWS::AccountId = 123456789012, got %s", ctx.PseudoParameters["AWS::AccountId"])
	}

	// Check parameter defaults are extracted
	if ctx.Parameters["Environment"] != "dev" {
		t.Errorf("expected Environment = dev, got %v", ctx.Parameters["Environment"])
	}
	if _, ok := ctx.Parameters["BucketName"]; ok {
		t.Error("BucketName should not have a value (no default)")
	}
}

func TestNewResolveContextNilTemplate(t *testing.T) {
	ctx := NewResolveContext(nil)

	if ctx.PseudoParameters == nil {
		t.Error("PseudoParameters should not be nil")
	}
	if ctx.Parameters == nil {
		t.Error("Parameters should not be nil")
	}
}

func TestResolveContextSetPseudoParameter(t *testing.T) {
	ctx := NewResolveContext(nil)
	ctx.SetPseudoParameter("AWS::Region", "eu-west-1")

	if ctx.PseudoParameters["AWS::Region"] != "eu-west-1" {
		t.Errorf("expected AWS::Region = eu-west-1, got %s", ctx.PseudoParameters["AWS::Region"])
	}
}

func TestResolveContextSetParameter(t *testing.T) {
	ctx := NewResolveContext(nil)
	ctx.SetParameter("MyParam", "value123")

	if ctx.Parameters["MyParam"] != "value123" {
		t.Errorf("expected MyParam = value123, got %v", ctx.Parameters["MyParam"])
	}
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	// Check that default actions are registered
	actions := []string{
		"Ref", "Fn::Sub", "Fn::GetAtt", "Fn::FindInMap",
		"Fn::Join", "Fn::If", "Fn::Select", "Fn::Base64",
		"Fn::GetAZs", "Fn::Split", "Fn::ImportValue", "Condition",
	}

	for _, name := range actions {
		if _, ok := registry.Get(name); !ok {
			t.Errorf("expected action '%s' to be registered", name)
		}
	}
}

func TestRegistryResolveNestedIntrinsics(t *testing.T) {
	registry := NewRegistry()
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"Stage": {
				Type:    "String",
				Default: "prod",
			},
		},
	}
	ctx := NewResolveContext(template)

	// Test nested intrinsic: Fn::Sub with Ref
	input := map[string]interface{}{
		"Fn::Sub": "Environment: ${Stage}",
	}

	result, err := registry.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Environment: prod"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRegistryResolveMapWithMultipleKeys(t *testing.T) {
	registry := NewRegistry()
	ctx := NewResolveContext(nil)
	ctx.SetParameter("Env", "test")

	// Non-intrinsic map with nested intrinsics
	input := map[string]interface{}{
		"Key1": map[string]interface{}{"Ref": "AWS::Region"},
		"Key2": "static",
	}

	result, err := registry.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}

	if resultMap["Key1"] != "us-east-1" {
		t.Errorf("expected Key1 = us-east-1, got %v", resultMap["Key1"])
	}
	if resultMap["Key2"] != "static" {
		t.Errorf("expected Key2 = static, got %v", resultMap["Key2"])
	}
}

func TestRegistryResolveSlice(t *testing.T) {
	registry := NewRegistry()
	ctx := NewResolveContext(nil)

	input := []interface{}{
		map[string]interface{}{"Ref": "AWS::Region"},
		"static",
		map[string]interface{}{"Ref": "AWS::AccountId"},
	}

	result, err := registry.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultSlice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(resultSlice) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(resultSlice))
	}
	if resultSlice[0] != "us-east-1" {
		t.Errorf("expected element 0 = us-east-1, got %v", resultSlice[0])
	}
	if resultSlice[1] != "static" {
		t.Errorf("expected element 1 = static, got %v", resultSlice[1])
	}
	if resultSlice[2] != "123456789012" {
		t.Errorf("expected element 2 = 123456789012, got %v", resultSlice[2])
	}
}

func TestIntrinsicError(t *testing.T) {
	err := NewIntrinsicError("Fn::Test", "test message")

	expected := "Fn::Test: test message"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
