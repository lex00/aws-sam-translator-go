package plugins

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestDefaultDefinitionBodyPlugin_Name(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()
	if plugin.Name() != "DefaultDefinitionBodyPlugin" {
		t.Errorf("Expected name 'DefaultDefinitionBodyPlugin', got '%s'", plugin.Name())
	}
}

func TestDefaultDefinitionBodyPlugin_Priority(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()
	if plugin.Priority() != 200 {
		t.Errorf("Expected priority 200, got %d", plugin.Priority())
	}
}

func TestDefaultDefinitionBodyPlugin_AddsDefaultForApi(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	template := &types.Template{
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
	defBody, ok := myApi.Properties["DefinitionBody"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected DefinitionBody to be set")
	}

	if defBody["swagger"] != "2.0" {
		t.Errorf("Expected swagger 2.0, got %v", defBody["swagger"])
	}

	info, ok := defBody["info"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected info to be set")
	}
	if info["title"] != "API" {
		t.Errorf("Expected title 'API', got %v", info["title"])
	}

	_, ok = defBody["paths"]
	if !ok {
		t.Errorf("Expected paths to be set")
	}
}

func TestDefaultDefinitionBodyPlugin_AddsDefaultForHttpApi(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyHttpApi": {
				Type: "AWS::Serverless::HttpApi",
				Properties: map[string]interface{}{
					"StageName": "$default",
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myApi := template.Resources["MyHttpApi"]
	_, ok := myApi.Properties["DefinitionBody"]
	if !ok {
		t.Fatalf("Expected DefinitionBody to be set")
	}
}

func TestDefaultDefinitionBodyPlugin_SkipsIfDefinitionBodyExists(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	existingDefBody := map[string]interface{}{
		"swagger": "3.0",
		"info": map[string]interface{}{
			"title": "My Custom API",
		},
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyApi": {
				Type: "AWS::Serverless::Api",
				Properties: map[string]interface{}{
					"StageName":      "prod",
					"DefinitionBody": existingDefBody,
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myApi := template.Resources["MyApi"]
	defBody := myApi.Properties["DefinitionBody"].(map[string]interface{})

	// Should keep existing definition
	if defBody["swagger"] != "3.0" {
		t.Errorf("Expected swagger 3.0, got %v", defBody["swagger"])
	}
	info := defBody["info"].(map[string]interface{})
	if info["title"] != "My Custom API" {
		t.Errorf("Expected title 'My Custom API', got %v", info["title"])
	}
}

func TestDefaultDefinitionBodyPlugin_SkipsIfDefinitionUriExists(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyApi": {
				Type: "AWS::Serverless::Api",
				Properties: map[string]interface{}{
					"StageName":     "prod",
					"DefinitionUri": "swagger.yaml",
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myApi := template.Resources["MyApi"]
	_, hasDefinitionBody := myApi.Properties["DefinitionBody"]
	if hasDefinitionBody {
		t.Errorf("Should not add DefinitionBody when DefinitionUri exists")
	}
}

func TestDefaultDefinitionBodyPlugin_AfterTransform(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()
	template := &types.Template{}

	err := plugin.AfterTransform(template)
	if err != nil {
		t.Errorf("AfterTransform should not error, got: %v", err)
	}
}
