package intrinsics

import (
	"reflect"
	"testing"
)

func TestJoinAction_Name(t *testing.T) {
	action := &JoinAction{}
	if action.Name() != "Fn::Join" {
		t.Errorf("expected 'Fn::Join', got %s", action.Name())
	}
}

func TestJoinAction_ResolveStaticStrings(t *testing.T) {
	action := &JoinAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"-", []interface{}{"a", "b", "c"}}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Static string values are evaluated at transform time
	expected := "a-b-c"
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestJoinAction_ResolvePassthrough(t *testing.T) {
	action := &JoinAction{}
	ctx := NewResolveContext(nil)

	// Input with non-string value should pass through
	input := []interface{}{"-", []interface{}{"a", map[string]interface{}{"Ref": "SomeParam"}, "c"}}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::Join": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestJoinAction_ResolveInvalidType(t *testing.T) {
	action := &JoinAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, "not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestJoinAction_ResolveInvalidLength(t *testing.T) {
	action := &JoinAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, []interface{}{"only one"})
	if err == nil {
		t.Error("expected error for array with wrong length")
	}
}

func TestIfAction_Name(t *testing.T) {
	action := &IfAction{}
	if action.Name() != "Fn::If" {
		t.Errorf("expected 'Fn::If', got %s", action.Name())
	}
}

func TestIfAction_ResolveUnknownCondition(t *testing.T) {
	action := &IfAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"MyCondition", "TrueValue", "FalseValue"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::If": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestIfAction_ResolveKnownConditionTrue(t *testing.T) {
	action := &IfAction{}
	ctx := NewResolveContext(nil)
	ctx.Conditions["MyCondition"] = true

	input := []interface{}{"MyCondition", "TrueValue", "FalseValue"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "TrueValue" {
		t.Errorf("expected 'TrueValue', got %v", result)
	}
}

func TestIfAction_ResolveKnownConditionFalse(t *testing.T) {
	action := &IfAction{}
	ctx := NewResolveContext(nil)
	ctx.Conditions["MyCondition"] = false

	input := []interface{}{"MyCondition", "TrueValue", "FalseValue"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "FalseValue" {
		t.Errorf("expected 'FalseValue', got %v", result)
	}
}

func TestIfAction_ResolveInvalidType(t *testing.T) {
	action := &IfAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, "not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestIfAction_ResolveInvalidLength(t *testing.T) {
	action := &IfAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, []interface{}{"Cond", "TrueVal"})
	if err == nil {
		t.Error("expected error for array with wrong length")
	}
}

func TestIfAction_ResolveInvalidConditionName(t *testing.T) {
	action := &IfAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, []interface{}{123, "TrueVal", "FalseVal"})
	if err == nil {
		t.Error("expected error for non-string condition name")
	}
}

func TestSelectAction_Name(t *testing.T) {
	action := &SelectAction{}
	if action.Name() != "Fn::Select" {
		t.Errorf("expected 'Fn::Select', got %s", action.Name())
	}
}

func TestSelectAction_Resolve(t *testing.T) {
	action := &SelectAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{1, []interface{}{"a", "b", "c"}}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::Select": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSelectAction_ResolveInvalidType(t *testing.T) {
	action := &SelectAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, "not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestSelectAction_ResolveInvalidLength(t *testing.T) {
	action := &SelectAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, []interface{}{0})
	if err == nil {
		t.Error("expected error for array with wrong length")
	}
}

func TestBase64Action_Name(t *testing.T) {
	action := &Base64Action{}
	if action.Name() != "Fn::Base64" {
		t.Errorf("expected 'Fn::Base64', got %s", action.Name())
	}
}

func TestBase64Action_Resolve(t *testing.T) {
	action := &Base64Action{}
	ctx := NewResolveContext(nil)

	input := "Hello World"

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::Base64": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGetAZsAction_Name(t *testing.T) {
	action := &GetAZsAction{}
	if action.Name() != "Fn::GetAZs" {
		t.Errorf("expected 'Fn::GetAZs', got %s", action.Name())
	}
}

func TestGetAZsAction_Resolve(t *testing.T) {
	action := &GetAZsAction{}
	ctx := NewResolveContext(nil)

	input := "us-east-1"

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::GetAZs": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSplitAction_Name(t *testing.T) {
	action := &SplitAction{}
	if action.Name() != "Fn::Split" {
		t.Errorf("expected 'Fn::Split', got %s", action.Name())
	}
}

func TestSplitAction_Resolve(t *testing.T) {
	action := &SplitAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{",", "a,b,c"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::Split": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSplitAction_ResolveInvalidType(t *testing.T) {
	action := &SplitAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, "not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestSplitAction_ResolveInvalidLength(t *testing.T) {
	action := &SplitAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, []interface{}{","})
	if err == nil {
		t.Error("expected error for array with wrong length")
	}
}

func TestImportValueAction_Name(t *testing.T) {
	action := &ImportValueAction{}
	if action.Name() != "Fn::ImportValue" {
		t.Errorf("expected 'Fn::ImportValue', got %s", action.Name())
	}
}

func TestImportValueAction_Resolve(t *testing.T) {
	action := &ImportValueAction{}
	ctx := NewResolveContext(nil)

	input := "MyExportedValue"

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::ImportValue": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestConditionAction_Name(t *testing.T) {
	action := &ConditionAction{}
	if action.Name() != "Condition" {
		t.Errorf("expected 'Condition', got %s", action.Name())
	}
}

func TestConditionAction_ResolveUnknown(t *testing.T) {
	action := &ConditionAction{}
	ctx := NewResolveContext(nil)

	result, err := action.Resolve(ctx, "MyCondition")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Condition": "MyCondition"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestConditionAction_ResolveKnownTrue(t *testing.T) {
	action := &ConditionAction{}
	ctx := NewResolveContext(nil)
	ctx.Conditions["MyCondition"] = true

	result, err := action.Resolve(ctx, "MyCondition")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != true {
		t.Errorf("expected true, got %v", result)
	}
}

func TestConditionAction_ResolveKnownFalse(t *testing.T) {
	action := &ConditionAction{}
	ctx := NewResolveContext(nil)
	ctx.Conditions["MyCondition"] = false

	result, err := action.Resolve(ctx, "MyCondition")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != false {
		t.Errorf("expected false, got %v", result)
	}
}

func TestConditionAction_ResolveInvalidType(t *testing.T) {
	action := &ConditionAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, 123)
	if err == nil {
		t.Error("expected error for non-string input")
	}
}

func TestIfAction_ResolveWithComplexValues(t *testing.T) {
	action := &IfAction{}
	ctx := NewResolveContext(nil)
	ctx.Conditions["IsProd"] = true

	trueValue := map[string]interface{}{
		"Type":  "m5.xlarge",
		"Count": 3,
	}
	falseValue := map[string]interface{}{
		"Type":  "t3.micro",
		"Count": 1,
	}

	input := []interface{}{"IsProd", trueValue, falseValue}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, trueValue) {
		t.Errorf("expected %v, got %v", trueValue, result)
	}
}

func TestBase64Action_ResolveWithIntrinsic(t *testing.T) {
	action := &Base64Action{}
	ctx := NewResolveContext(nil)

	// Base64 with nested intrinsic
	input := map[string]interface{}{
		"Fn::Sub": "Hello ${AWS::Region}",
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"Fn::Base64": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
