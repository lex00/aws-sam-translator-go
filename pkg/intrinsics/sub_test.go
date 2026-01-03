package intrinsics

import (
	"reflect"
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestSubAction_Name(t *testing.T) {
	action := &SubAction{}
	if action.Name() != "Fn::Sub" {
		t.Errorf("expected 'Fn::Sub', got %s", action.Name())
	}
}

func TestSubAction_ResolveStringFormSimple(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	result, err := action.Resolve(ctx, "Hello World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("expected 'Hello World', got %v", result)
	}
}

func TestSubAction_ResolveStringFormWithPseudoParam(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	result, err := action.Resolve(ctx, "Region: ${AWS::Region}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Region: us-east-1" {
		t.Errorf("expected 'Region: us-east-1', got %v", result)
	}
}

func TestSubAction_ResolveStringFormWithParameter(t *testing.T) {
	action := &SubAction{}
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"Environment": {
				Type:    "String",
				Default: "prod",
			},
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "Env: ${Environment}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Env: prod" {
		t.Errorf("expected 'Env: prod', got %v", result)
	}
}

func TestSubAction_ResolveStringFormMultipleVars(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)
	ctx.SetParameter("Stage", "dev")

	result, err := action.Resolve(ctx, "arn:aws:s3:::${Stage}-bucket-${AWS::Region}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "arn:aws:s3:::dev-bucket-us-east-1"
	if result != expected {
		t.Errorf("expected %q, got %v", expected, result)
	}
}

func TestSubAction_ResolveArrayFormSimple(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{
		"Hello ${Name}",
		map[string]interface{}{"Name": "World"},
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello World" {
		t.Errorf("expected 'Hello World', got %v", result)
	}
}

func TestSubAction_ResolveArrayFormWithPseudoAndLocal(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{
		"arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:${FunctionName}",
		map[string]interface{}{"FunctionName": "my-function"},
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "arn:aws:lambda:us-east-1:123456789012:function:my-function"
	if result != expected {
		t.Errorf("expected %q, got %v", expected, result)
	}
}

func TestSubAction_ResolveArrayFormLocalOverridesPseudo(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	// Local variable should take precedence
	input := []interface{}{
		"Region: ${AWS::Region}",
		map[string]interface{}{"AWS::Region": "custom-region"},
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Region: custom-region" {
		t.Errorf("expected 'Region: custom-region', got %v", result)
	}
}

func TestSubAction_ResolveWithResourceRef(t *testing.T) {
	action := &SubAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Lambda::Function",
			},
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "Function: ${MyFunction}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Resource reference should be preserved
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if _, ok := resultMap["Fn::Sub"]; !ok {
		t.Errorf("expected Fn::Sub in result, got %v", result)
	}
}

func TestSubAction_ResolveWithGetAttSyntax(t *testing.T) {
	action := &SubAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Lambda::Function",
			},
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "ARN: ${MyFunction.Arn}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// GetAtt reference should be preserved
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if _, ok := resultMap["Fn::Sub"]; !ok {
		t.Errorf("expected Fn::Sub in result, got %v", result)
	}
}

func TestSubAction_ResolveWithCachedResourceAttr(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)
	ctx.Resources["MyFunction"] = map[string]interface{}{
		"Arn": "arn:aws:lambda:us-east-1:123456789012:function:my-func",
	}

	result, err := action.Resolve(ctx, "ARN: ${MyFunction.Arn}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "ARN: arn:aws:lambda:us-east-1:123456789012:function:my-func"
	if result != expected {
		t.Errorf("expected %q, got %v", expected, result)
	}
}

func TestSubAction_ResolveArrayFormInvalidLength(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"just one element"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for array with 1 element")
	}
}

func TestSubAction_ResolveArrayFormInvalidFirstElement(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{
		123, // not a string
		map[string]interface{}{},
	}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-string first element")
	}
}

func TestSubAction_ResolveArrayFormInvalidSecondElement(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{
		"Hello ${Name}",
		"not a map",
	}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-map second element")
	}
}

func TestSubAction_ResolveInvalidType(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, 123)
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestSubAction_ResolveNumericParameter(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)
	ctx.SetParameter("Port", 8080)

	result, err := action.Resolve(ctx, "Port: ${Port}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Port: 8080" {
		t.Errorf("expected 'Port: 8080', got %v", result)
	}
}

func TestSubAction_ResolveArrayFormPreservesUnresolvedVars(t *testing.T) {
	action := &SubAction{}
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyBucket": {Type: "AWS::S3::Bucket"},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{
		"s3://${MyBucket}/${Prefix}",
		map[string]interface{}{"Prefix": map[string]interface{}{"Ref": "PrefixParam"}},
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve as Fn::Sub with filtered variables
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}

	subValue, ok := resultMap["Fn::Sub"]
	if !ok {
		t.Fatalf("expected Fn::Sub key in result")
	}

	// Should be array form with filtered vars
	subArr, ok := subValue.([]interface{})
	if ok && len(subArr) == 2 {
		// Check that filtered vars contain only unresolved ones
		if varMap, ok := subArr[1].(map[string]interface{}); ok {
			if _, hasPrefix := varMap["Prefix"]; !hasPrefix {
				t.Error("expected Prefix in filtered variables")
			}
		}
	}
}

func TestSubAction_ValueToString(t *testing.T) {
	action := &SubAction{}

	tests := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "hello"},
		{42, "42"},
		{int64(100), "100"},
		{3.14, "3.14"},
		{10.0, "10"}, // Whole number should be formatted as int
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		result := action.valueToString(tt.input)
		if result != tt.expected {
			t.Errorf("valueToString(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestSubAction_ResolveNoVariables(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	result, err := action.Resolve(ctx, "No variables here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "No variables here" {
		t.Errorf("expected 'No variables here', got %v", result)
	}
}

func TestSubAction_ResolveEmptyString(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	result, err := action.Resolve(ctx, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %v", result)
	}
}

func TestSubAction_ResolvePartialSubstitution(t *testing.T) {
	action := &SubAction{}
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"KnownParam":   {Type: "String", Default: "known-value"},
			"UnknownParam": {Type: "String"}, // No default
		},
	}
	ctx := NewResolveContext(template)

	result, err := action.Resolve(ctx, "${KnownParam}-${UnknownParam}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve Fn::Sub with partial substitution
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T: %v", result, result)
	}

	subValue, ok := resultMap["Fn::Sub"]
	if !ok {
		t.Fatalf("expected Fn::Sub in result")
	}

	// The known param should be substituted, unknown preserved
	subStr, ok := subValue.(string)
	if !ok {
		t.Fatalf("expected string sub value, got %T", subValue)
	}

	expected := "known-value-${UnknownParam}"
	if subStr != expected {
		t.Errorf("expected %q, got %q", expected, subStr)
	}
}

func TestSubAction_ResolveComplexARN(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)
	ctx.SetPseudoParameter("AWS::Region", "us-west-2")
	ctx.SetPseudoParameter("AWS::AccountId", "111122223333")

	input := []interface{}{
		"arn:${AWS::Partition}:lambda:${AWS::Region}:${AWS::AccountId}:function:${FunctionName}",
		map[string]interface{}{"FunctionName": "MyLambda"},
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "arn:aws:lambda:us-west-2:111122223333:function:MyLambda"
	if result != expected {
		t.Errorf("expected %q, got %v", expected, result)
	}
}

func TestSubAction_ResolveFullyResolved(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)
	ctx.SetParameter("Env", "prod")
	ctx.SetParameter("Name", "myapp")

	result, err := action.Resolve(ctx, "${Env}-${Name}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When fully resolved, should return plain string
	resultStr, ok := result.(string)
	if !ok {
		t.Fatalf("expected string result when fully resolved, got %T: %v", result, result)
	}

	expected := "prod-myapp"
	if resultStr != expected {
		t.Errorf("expected %q, got %q", expected, resultStr)
	}
}

func TestSubAction_ResolveArrayFullyResolved(t *testing.T) {
	action := &SubAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{
		"Region: ${Region}, Account: ${Account}",
		map[string]interface{}{
			"Region":  "eu-west-1",
			"Account": "999888777666",
		},
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When fully resolved from local vars, should return plain string
	expected := "Region: eu-west-1, Account: 999888777666"
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %q, got %v", expected, result)
	}
}
