package intrinsics

import (
	"reflect"
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestFindInMapAction_Name(t *testing.T) {
	action := &FindInMapAction{}
	if action.Name() != "Fn::FindInMap" {
		t.Errorf("expected 'Fn::FindInMap', got %s", action.Name())
	}
}

func TestFindInMapAction_ResolveBasic(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionMap": map[string]interface{}{
				"us-east-1": map[string]interface{}{
					"AMI": "ami-12345678",
				},
				"us-west-2": map[string]interface{}{
					"AMI": "ami-87654321",
				},
			},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"RegionMap", "us-east-1", "AMI"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "ami-12345678" {
		t.Errorf("expected 'ami-12345678', got %v", result)
	}
}

func TestFindInMapAction_ResolveNestedValue(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"EnvironmentConfig": map[string]interface{}{
				"prod": map[string]interface{}{
					"InstanceType": "m5.xlarge",
					"MinInstances": 3,
				},
				"dev": map[string]interface{}{
					"InstanceType": "t3.micro",
					"MinInstances": 1,
				},
			},
		},
	}
	ctx := NewResolveContext(template)

	// Test string value
	result, err := action.Resolve(ctx, []interface{}{"EnvironmentConfig", "prod", "InstanceType"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "m5.xlarge" {
		t.Errorf("expected 'm5.xlarge', got %v", result)
	}

	// Test numeric value
	result, err = action.Resolve(ctx, []interface{}{"EnvironmentConfig", "dev", "MinInstances"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 1 {
		t.Errorf("expected 1, got %v", result)
	}
}

func TestFindInMapAction_ResolveMappingNotFound(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionMap": map[string]interface{}{},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"NonExistentMap", "key1", "key2"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-existent mapping")
	}
}

func TestFindInMapAction_ResolveTopLevelKeyNotFound(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionMap": map[string]interface{}{
				"us-east-1": map[string]interface{}{
					"AMI": "ami-12345678",
				},
			},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"RegionMap", "eu-west-1", "AMI"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-existent top level key")
	}
}

func TestFindInMapAction_ResolveSecondLevelKeyNotFound(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionMap": map[string]interface{}{
				"us-east-1": map[string]interface{}{
					"AMI": "ami-12345678",
				},
			},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"RegionMap", "us-east-1", "InstanceType"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-existent second level key")
	}
}

func TestFindInMapAction_ResolveInvalidArrayLength(t *testing.T) {
	action := &FindInMapAction{}
	ctx := NewResolveContext(nil)

	tests := []struct {
		name  string
		input []interface{}
	}{
		{"too few", []interface{}{"MapName", "Key1"}},
		{"too many", []interface{}{"MapName", "Key1", "Key2", "Key3"}},
		{"empty", []interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := action.Resolve(ctx, tt.input)
			if err == nil {
				t.Error("expected error for invalid array length")
			}
		})
	}
}

func TestFindInMapAction_ResolveInvalidType(t *testing.T) {
	action := &FindInMapAction{}
	ctx := NewResolveContext(nil)

	_, err := action.Resolve(ctx, "not an array")
	if err == nil {
		t.Error("expected error for non-array input")
	}
}

func TestFindInMapAction_ResolveWithIntrinsicKey(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionMap": map[string]interface{}{
				"us-east-1": map[string]interface{}{
					"AMI": "ami-12345678",
				},
			},
		},
	}
	ctx := NewResolveContext(template)

	// Key is an unresolved intrinsic
	input := []interface{}{
		"RegionMap",
		map[string]interface{}{"Ref": "AWS::Region"}, // Unresolved intrinsic
		"AMI",
	}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve for CloudFormation when key is intrinsic
	expected := map[string]interface{}{"Fn::FindInMap": input}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestFindInMapAction_ResolveNilMappings(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: nil,
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"RegionMap", "us-east-1", "AMI"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for nil mappings")
	}
}

func TestFindInMapAction_ResolveNilTemplate(t *testing.T) {
	action := &FindInMapAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"RegionMap", "us-east-1", "AMI"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for nil template")
	}
}

func TestFindInMapAction_ResolveInvalidMappingType(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionMap": "not a map",
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"RegionMap", "us-east-1", "AMI"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-map mapping")
	}
}

func TestFindInMapAction_ResolveInvalidTopLevelValueType(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionMap": map[string]interface{}{
				"us-east-1": "not a map",
			},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"RegionMap", "us-east-1", "AMI"}

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for non-map top level value")
	}
}

func TestFindInMapAction_ResolveComplexValue(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"Config": map[string]interface{}{
				"prod": map[string]interface{}{
					"Settings": map[string]interface{}{
						"Enabled": true,
						"Count":   5,
					},
				},
			},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"Config", "prod", "Settings"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{
		"Enabled": true,
		"Count":   5,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestFindInMapAction_ResolveListValue(t *testing.T) {
	action := &FindInMapAction{}
	template := &types.Template{
		Mappings: map[string]interface{}{
			"RegionAZs": map[string]interface{}{
				"us-east-1": map[string]interface{}{
					"AZs": []interface{}{"us-east-1a", "us-east-1b", "us-east-1c"},
				},
			},
		},
	}
	ctx := NewResolveContext(template)

	input := []interface{}{"RegionAZs", "us-east-1", "AZs"}

	result, err := action.Resolve(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []interface{}{"us-east-1a", "us-east-1b", "us-east-1c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestFindInMapAction_ResolveInvalidKeyType(t *testing.T) {
	action := &FindInMapAction{}
	ctx := NewResolveContext(nil)

	input := []interface{}{"RegionMap", 123, "AMI"} // Numeric key

	_, err := action.Resolve(ctx, input)
	if err == nil {
		t.Error("expected error for invalid key type")
	}
}
