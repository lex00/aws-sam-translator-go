package intrinsics

import (
	"reflect"
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestGetAttAction_Name(t *testing.T) {
	action := &GetAttAction{}
	if action.Name() != "Fn::GetAtt" {
		t.Errorf("expected 'Fn::GetAtt', got %s", action.Name())
	}
}

func TestGetAttAction_ResolveArrayForm(t *testing.T) {
	action := &GetAttAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {Type: "AWS::Lambda::Function"},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"MyFunction", "Arn"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve for CloudFormation
	expected := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetAttAction_ResolveStringForm(t *testing.T) {
	action := &GetAttAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyBucket": {Type: "AWS::S3::Bucket"},
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "MyBucket.Arn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyBucket", "Arn"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetAttAction_ResolveWithCachedAttribute(t *testing.T) {
	action := &GetAttAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {Type: "AWS::Lambda::Function"},
		},
	}
	ctx := NewResolveContext(template)
	ctx.Resources["MyFunction"] = map[string]interface{}{
		"Arn": "arn:aws:lambda:us-east-1:123456789012:function:my-func",
	}

	input := []interface{}{"MyFunction", "Arn"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "arn:aws:lambda:us-east-1:123456789012:function:my-func"
	if result != expected {
		t.Errorf("expected %q, got %v", expected, result)
	}
}

func TestGetAttAction_ResolveNestedAttribute(t *testing.T) {
	action := &GetAttAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyApi": {Type: "AWS::ApiGateway::RestApi"},
		},
	}
	ctx := NewResolveContext(template)

	// Nested attribute like ["MyApi", "RootResourceId", "Value"]
	input := []interface{}{"MyApi", "RootResourceId", "Value"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyApi", "RootResourceId.Value"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetAttAction_ResolveResourceNotFound(t *testing.T) {
	action := &GetAttAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {Type: "AWS::Lambda::Function"},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"NonExistentResource", "Arn"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-existent resource")
	}
}

func TestGetAttAction_ResolveInvalidArrayLength(t *testing.T) {
	action := &GetAttAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"OnlyOne"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for array with less than 2 elements")
	}
}

func TestGetAttAction_ResolveInvalidResourceNameType(t *testing.T) {
	action := &GetAttAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{123, "Arn"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-string resource name")
	}
}

func TestGetAttAction_ResolveInvalidAttributeNameType(t *testing.T) {
	action := &GetAttAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"MyFunction", 123}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-string attribute name")
	}
}

func TestGetAttAction_ResolveInvalidStringFormat(t *testing.T) {
	action := &GetAttAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, "NoDotsHere")
	if err == nil {
		t.Error("expected error for invalid string format")
	}
}

func TestGetAttAction_ResolveInvalidType(t *testing.T) {
	action := &GetAttAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, 123)
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestGetAttAction_ResolveWithNilTemplate(t *testing.T) {
	action := &GetAttAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"MyFunction", "Arn"}

	// Should still work (just preserve the GetAtt)
	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetAttAction_ResolveStringFormWithNestedAttribute(t *testing.T) {
	action := &GetAttAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyResource": {Type: "AWS::CloudFormation::CustomResource"},
		},
	}
	ctx := NewResolveContext(template)

	// String form: "Resource.Attr1.Attr2" is split into ["Resource", "Attr1.Attr2"]
	result, err := action.Resolve(ctx, "MyResource.Nested.Attribute")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyResource", "Nested.Attribute"},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetAttAction_ResolveArrayFormPreservesFormat(t *testing.T) {
	action := &GetAttAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyTable": {Type: "AWS::DynamoDB::Table"},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"MyTable", "StreamArn"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the result preserves the proper structure
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T", result)
	}

	getAttValue := resultMap["Fn::GetAtt"]
	getAttArr, ok := getAttValue.([]interface{})
	if !ok {
		t.Fatalf("expected array value, got %T", getAttValue)
	}

	if len(getAttArr) != 2 {
		t.Errorf("expected 2 elements, got %d", len(getAttArr))
	}
	if getAttArr[0] != "MyTable" {
		t.Errorf("expected 'MyTable', got %v", getAttArr[0])
	}
	if getAttArr[1] != "StreamArn" {
		t.Errorf("expected 'StreamArn', got %v", getAttArr[1])
	}
}
