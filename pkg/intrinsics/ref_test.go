package intrinsics

import (
	"reflect"
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestRefAction_Name(t *testing.T) {
	action := &RefAction{}
	if action.Name() != "Ref" {
		t.Errorf("expected 'Ref', got %s", action.Name())
	}
}

func TestRefAction_ResolvePseudoParameter(t *testing.T) {
	action := &RefAction{}
	ctx := NewResolveContext(nil)

	tests := []struct {
		name     string
		refValue string
		expected string
	}{
		{"AWS::Region", "AWS::Region", "us-east-1"},
		{"AWS::AccountId", "AWS::AccountId", "123456789012"},
		{"AWS::StackName", "AWS::StackName", "sam-app"},
		{"AWS::Partition", "AWS::Partition", "aws"},
		{"AWS::URLSuffix", "AWS::URLSuffix", "amazonaws.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := action.Resolve(ctx, tt.refValue)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %v", tt.expected, result)
			}
		})
	}
}

func TestRefAction_ResolveAWSNoValue(t *testing.T) {
	action := &RefAction{}
	ctx := NewResolveContext(nil)

	result, err := action.Resolve(ctx, "AWS::NoValue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsNoValue(result) {
		t.Errorf("expected NoValue sentinel for AWS::NoValue, got %v", result)
	}
}

func TestRefAction_ResolveParameter(t *testing.T) {
	action := &RefAction{}
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"Environment": {
				Type:    "String",
				Default: "production",
			},
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "Environment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "production" {
		t.Errorf("expected 'production', got %v", result)
	}
}

func TestRefAction_ResolveParameterWithOverride(t *testing.T) {
	action := &RefAction{}
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"Environment": {
				Type:    "String",
				Default: "production",
			},
		},
	}
	ctx := NewResolveContext(template)
	ctx.SetParameter("Environment", "staging")

	result, err := action.Resolve(ctx, "Environment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "staging" {
		t.Errorf("expected 'staging', got %v", result)
	}
}

func TestRefAction_ResolveParameterNoValue(t *testing.T) {
	action := &RefAction{}
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"BucketName": {
				Type: "String",
				// No default
			},
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "BucketName")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve the Ref for CloudFormation
	expected := map[string]interface{}{"Ref": "BucketName"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestRefAction_ResolveResource(t *testing.T) {
	action := &RefAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
			},
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "MyFunction")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve the Ref for CloudFormation
	expected := map[string]interface{}{"Ref": "MyFunction"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestRefAction_ResolveUnknown(t *testing.T) {
	action := &RefAction{}
	ctx := NewResolveContext(nil)

	result, err := action.Resolve(ctx, "UnknownRef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve the Ref for CloudFormation
	expected := map[string]interface{}{"Ref": "UnknownRef"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestRefAction_ResolveInvalidType(t *testing.T) {
	action := &RefAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, 123)
	if err == nil {
		t.Error("expected error for non-string value")
	}
}

func TestRefAction_ResolveCustomPseudoParameter(t *testing.T) {
	action := &RefAction{}
	ctx := NewResolveContext(nil)
	ctx.SetPseudoParameter("AWS::Region", "eu-central-1")

	result, err := action.Resolve(ctx, "AWS::Region")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "eu-central-1" {
		t.Errorf("expected 'eu-central-1', got %v", result)
	}
}
