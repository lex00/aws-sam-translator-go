package intrinsics

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestNewResolver(t *testing.T) {
	ctx := NewResolveContext(nil)
	resolver := NewResolver(ctx)

	if resolver == nil {
		t.Fatal("expected resolver to be non-nil")
	}
	if resolver.context != ctx {
		t.Error("expected resolver context to match")
	}
	if resolver.dependencies == nil {
		t.Error("expected dependencies to be initialized")
	}
}

func TestResolverWithOptions(t *testing.T) {
	ctx := NewResolveContext(nil)

	idMap := map[string]string{"OldFunc": "NewFunc"}
	placeholders := map[string]interface{}{"PLACEHOLDER": true}
	tracker := NewDependencyTracker()

	resolver := NewResolver(ctx,
		WithLogicalIDMap(idMap),
		WithPlaceholders(placeholders),
		WithDependencyTracker(tracker),
	)

	if resolver.logicalIDMap["OldFunc"] != "NewFunc" {
		t.Error("expected logical ID map to be set")
	}
	if _, ok := resolver.placeholders["PLACEHOLDER"]; !ok {
		t.Error("expected placeholders to be set")
	}
	if resolver.dependencies != tracker {
		t.Error("expected dependency tracker to be set")
	}
}

func TestResolverResolveSimpleValue(t *testing.T) {
	ctx := NewResolveContext(nil)
	resolver := NewResolver(ctx)

	// Test simple string
	result, err := resolver.Resolve("simple string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "simple string" {
		t.Errorf("expected 'simple string', got %v", result)
	}

	// Test number
	result, err = resolver.Resolve(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 42 {
		t.Errorf("expected 42, got %v", result)
	}
}

func TestResolverResolveRef(t *testing.T) {
	template := &types.Template{
		Parameters: map[string]types.Parameter{
			"Environment": {
				Type:    "String",
				Default: "prod",
			},
		},
	}
	ctx := NewResolveContext(template)
	resolver := NewResolver(ctx)

	// Resolve parameter reference
	input := map[string]interface{}{"Ref": "Environment"}
	result, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "prod" {
		t.Errorf("expected 'prod', got %v", result)
	}
}

func TestResolverResolveNestedIntrinsics(t *testing.T) {
	ctx := NewResolveContext(nil)
	ctx.SetParameter("Stage", "dev")
	resolver := NewResolver(ctx)

	// Nested: Fn::Sub with parameter
	input := map[string]interface{}{
		"Fn::Sub": "Stage is ${Stage}",
	}

	result, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Stage is dev" {
		t.Errorf("expected 'Stage is dev', got %v", result)
	}
}

func TestResolverLogicalIDMutation(t *testing.T) {
	template := &types.Template{
		Resources: map[string]types.Resource{
			"NewFunction": {
				Type: "AWS::Lambda::Function",
			},
		},
	}
	ctx := NewResolveContext(template)
	resolver := NewResolver(ctx, WithLogicalIDMap(map[string]string{
		"OldFunction": "NewFunction",
	}))

	// Ref to old logical ID should be mutated
	input := map[string]interface{}{"Ref": "OldFunction"}
	result, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should preserve the Ref (resource reference) with new ID
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T", result)
	}
	if resultMap["Ref"] != "NewFunction" {
		t.Errorf("expected Ref to 'NewFunction', got %v", resultMap["Ref"])
	}
}

func TestResolverPlaceholderProtection(t *testing.T) {
	ctx := NewResolveContext(nil)
	resolver := NewResolver(ctx, WithPlaceholders(map[string]interface{}{
		"PROTECTED_VALUE": true,
	}))

	// Ref to placeholder should be preserved
	input := map[string]interface{}{"Ref": "PROTECTED_VALUE"}
	result, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be preserved as-is
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map result, got %T", result)
	}
	if resultMap["Ref"] != "PROTECTED_VALUE" {
		t.Errorf("expected Ref to 'PROTECTED_VALUE', got %v", resultMap["Ref"])
	}
}

func TestResolverDependencyTracking(t *testing.T) {
	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyBucket": {
				Type: "AWS::S3::Bucket",
			},
			"MyFunction": {
				Type: "AWS::Lambda::Function",
			},
		},
	}
	ctx := NewResolveContext(template)
	resolver := NewResolver(ctx)

	// Simulate resolving a resource property with a Ref
	input := map[string]interface{}{
		"Properties": map[string]interface{}{
			"BucketName": map[string]interface{}{"Ref": "MyBucket"},
		},
	}
	_, err := resolver.resolveValue(input, "Resources.MyFunction")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check dependencies were tracked
	deps := resolver.GetDependencies()
	if !deps.HasDependency("MyFunction", "MyBucket") {
		t.Error("expected MyFunction to depend on MyBucket")
	}
}

func TestResolverResolveTemplate(t *testing.T) {
	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Description:              "Test template",
		Parameters: map[string]types.Parameter{
			"Environment": {
				Type:    "String",
				Default: "prod",
			},
		},
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Lambda::Function",
				Properties: map[string]interface{}{
					"FunctionName": map[string]interface{}{
						"Fn::Sub": "${Environment}-my-function",
					},
				},
			},
		},
		Outputs: map[string]types.Output{
			"FunctionArn": {
				Value: map[string]interface{}{
					"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
				},
			},
		},
	}

	ctx := NewResolveContext(template)
	resolver := NewResolver(ctx)

	result, err := resolver.ResolveTemplate(template)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check resource property was resolved
	funcProps := result.Resources["MyFunction"].Properties
	if funcProps["FunctionName"] != "prod-my-function" {
		t.Errorf("expected 'prod-my-function', got %v", funcProps["FunctionName"])
	}
}

func TestResolverResolveTemplateWithLogicalIDMapping(t *testing.T) {
	template := &types.Template{
		Resources: map[string]types.Resource{
			"OldName": {
				Type:       "AWS::Lambda::Function",
				Properties: map[string]interface{}{},
			},
		},
	}

	ctx := NewResolveContext(template)
	resolver := NewResolver(ctx, WithLogicalIDMap(map[string]string{
		"OldName": "NewName",
	}))

	result, err := resolver.ResolveTemplate(template)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check resource was renamed
	if _, ok := result.Resources["NewName"]; !ok {
		t.Error("expected resource to be renamed to 'NewName'")
	}
	if _, ok := result.Resources["OldName"]; ok {
		t.Error("old resource name should not exist")
	}
}

func TestResolverSetMethods(t *testing.T) {
	ctx := NewResolveContext(nil)
	resolver := NewResolver(ctx)

	// Test SetParameter
	resolver.SetParameter("MyParam", "value")
	if ctx.Parameters["MyParam"] != "value" {
		t.Error("SetParameter did not work")
	}

	// Test SetResourceAttribute
	resolver.SetResourceAttribute("MyResource", "Arn", "arn:aws:...")
	if ctx.Resources["MyResource"]["Arn"] != "arn:aws:..." {
		t.Error("SetResourceAttribute did not work")
	}

	// Test AddLogicalIDMapping
	resolver.AddLogicalIDMapping("Old", "New")
	if resolver.logicalIDMap["Old"] != "New" {
		t.Error("AddLogicalIDMapping did not work")
	}

	// Test AddPlaceholder
	resolver.AddPlaceholder("PH", true)
	if _, ok := resolver.placeholders["PH"]; !ok {
		t.Error("AddPlaceholder did not work")
	}
}

func TestResolverPreserveCFRuntimeIntrinsics(t *testing.T) {
	ctx := NewResolveContext(nil)
	resolver := NewResolver(ctx)

	// Fn::ImportValue should always be preserved
	input := map[string]interface{}{
		"Fn::ImportValue": "SharedValue",
	}
	result, err := resolver.Resolve(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if resultMap["Fn::ImportValue"] != "SharedValue" {
		t.Error("Fn::ImportValue should be preserved")
	}
}

func TestResolverHandleNoValue(t *testing.T) {
	ctx := NewResolveContext(nil)
	resolver := NewResolver(ctx)

	// Map with NoValue property should have it removed
	input := map[string]interface{}{
		"Keep":   "value",
		"Remove": NoValue{},
	}

	result, err := resolver.resolveMap(input, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result["Remove"]; ok {
		t.Error("NoValue property should be removed")
	}
	if result["Keep"] != "value" {
		t.Error("Keep property should be preserved")
	}
}

func TestResolverResolveSlice(t *testing.T) {
	ctx := NewResolveContext(nil)
	ctx.SetParameter("Item1", "resolved1")
	resolver := NewResolver(ctx)

	input := []interface{}{
		map[string]interface{}{"Ref": "Item1"},
		"static",
		map[string]interface{}{"Ref": "AWS::Region"},
	}

	result, err := resolver.Resolve(input)
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
	if resultSlice[0] != "resolved1" {
		t.Errorf("expected 'resolved1', got %v", resultSlice[0])
	}
	if resultSlice[1] != "static" {
		t.Errorf("expected 'static', got %v", resultSlice[1])
	}
	if resultSlice[2] != "us-east-1" {
		t.Errorf("expected 'us-east-1', got %v", resultSlice[2])
	}
}

func TestResolverMapDependsOn(t *testing.T) {
	ctx := NewResolveContext(nil)
	resolver := NewResolver(ctx, WithLogicalIDMap(map[string]string{
		"OldDep": "NewDep",
	}))

	// Test string DependsOn
	result := resolver.mapDependsOn("OldDep")
	if result != "NewDep" {
		t.Errorf("expected 'NewDep', got %v", result)
	}

	// Test slice DependsOn
	sliceResult := resolver.mapDependsOn([]interface{}{"OldDep", "Other"})
	slice, ok := sliceResult.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", sliceResult)
	}
	if slice[0] != "NewDep" {
		t.Errorf("expected 'NewDep', got %v", slice[0])
	}
	if slice[1] != "Other" {
		t.Errorf("expected 'Other', got %v", slice[1])
	}
}
