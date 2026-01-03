package plugins

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestGlobalsPlugin_Name(t *testing.T) {
	plugin := NewGlobalsPlugin()
	if plugin.Name() != "GlobalsPlugin" {
		t.Errorf("Expected name 'GlobalsPlugin', got '%s'", plugin.Name())
	}
}

func TestGlobalsPlugin_Priority(t *testing.T) {
	plugin := NewGlobalsPlugin()
	if plugin.Priority() != 100 {
		t.Errorf("Expected priority 100, got %d", plugin.Priority())
	}
}

func TestGlobalsPlugin_ApplyFunctionGlobals(t *testing.T) {
	plugin := NewGlobalsPlugin()

	template := &types.Template{
		Globals: map[string]interface{}{
			"Function": map[string]interface{}{
				"Runtime":    "python3.9",
				"Timeout":    30,
				"MemorySize": 512,
			},
		},
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
				},
			},
			"MyOtherFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "app.handler",
					"Runtime": "nodejs18.x", // Should not be overridden
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	// Check MyFunction got global properties
	myFunc := template.Resources["MyFunction"]
	if myFunc.Properties["Runtime"] != "python3.9" {
		t.Errorf("Expected Runtime 'python3.9', got '%v'", myFunc.Properties["Runtime"])
	}
	if myFunc.Properties["Timeout"] != 30 {
		t.Errorf("Expected Timeout 30, got %v", myFunc.Properties["Timeout"])
	}
	if myFunc.Properties["MemorySize"] != 512 {
		t.Errorf("Expected MemorySize 512, got %v", myFunc.Properties["MemorySize"])
	}
	if myFunc.Properties["Handler"] != "index.handler" {
		t.Errorf("Expected Handler 'index.handler', got '%v'", myFunc.Properties["Handler"])
	}

	// Check MyOtherFunction kept its own Runtime
	myOtherFunc := template.Resources["MyOtherFunction"]
	if myOtherFunc.Properties["Runtime"] != "nodejs18.x" {
		t.Errorf("Expected Runtime 'nodejs18.x', got '%v'", myOtherFunc.Properties["Runtime"])
	}
	// But got global Timeout and MemorySize
	if myOtherFunc.Properties["Timeout"] != 30 {
		t.Errorf("Expected Timeout 30, got %v", myOtherFunc.Properties["Timeout"])
	}
}

func TestGlobalsPlugin_ApplyApiGlobals(t *testing.T) {
	plugin := NewGlobalsPlugin()

	template := &types.Template{
		Globals: map[string]interface{}{
			"Api": map[string]interface{}{
				"Auth": map[string]interface{}{
					"ApiKeyRequired": true,
				},
			},
		},
		Resources: map[string]types.Resource{
			"MyApi": {
				Type: "AWS::Serverless::Api",
				Properties: map[string]interface{}{
					"StageName": "prod",
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myApi := template.Resources["MyApi"]
	auth, ok := myApi.Properties["Auth"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected Auth to be set")
	}
	if auth["ApiKeyRequired"] != true {
		t.Errorf("Expected ApiKeyRequired true, got %v", auth["ApiKeyRequired"])
	}
}

func TestGlobalsPlugin_NoGlobals(t *testing.T) {
	plugin := NewGlobalsPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	// Should not modify function without globals
	myFunc := template.Resources["MyFunction"]
	if len(myFunc.Properties) != 1 {
		t.Errorf("Expected 1 property, got %d", len(myFunc.Properties))
	}
}

func TestGlobalsPlugin_AfterTransform(t *testing.T) {
	plugin := NewGlobalsPlugin()
	template := &types.Template{}

	err := plugin.AfterTransform(template)
	if err != nil {
		t.Errorf("AfterTransform should not error, got: %v", err)
	}
}

func TestDeepCopy(t *testing.T) {
	original := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nested": "value2",
		},
		"key3": []interface{}{1, 2, 3},
	}

	copied := deepCopy(original)

	// Modify the copy
	copiedMap := copied.(map[string]interface{})
	copiedMap["key1"] = "modified"
	copiedNested := copiedMap["key2"].(map[string]interface{})
	copiedNested["nested"] = "modified"
	copiedArray := copiedMap["key3"].([]interface{})
	copiedArray[0] = 99

	// Original should be unchanged
	if original["key1"] != "value1" {
		t.Errorf("Original key1 was modified")
	}
	originalNested := original["key2"].(map[string]interface{})
	if originalNested["nested"] != "value2" {
		t.Errorf("Original nested value was modified")
	}
	originalArray := original["key3"].([]interface{})
	if originalArray[0] != 1 {
		t.Errorf("Original array was modified")
	}
}
