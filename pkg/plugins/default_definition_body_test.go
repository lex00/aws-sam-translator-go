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
	if plugin.Priority() != 500 {
		t.Errorf("Expected priority 500, got %d", plugin.Priority())
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
	// Title is now the logical ID
	if info["title"] != "MyApi" {
		t.Errorf("Expected title 'MyApi', got %v", info["title"])
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
	defBody, ok := myApi.Properties["DefinitionBody"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected DefinitionBody to be set")
	}

	// HttpApi uses OpenAPI 3.0
	if defBody["openapi"] != "3.0.1" {
		t.Errorf("Expected openapi 3.0.1, got %v", defBody["openapi"])
	}

	info, ok := defBody["info"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected info to be set")
	}
	if info["title"] != "MyHttpApi" {
		t.Errorf("Expected title 'MyHttpApi', got %v", info["title"])
	}
}

func TestDefaultDefinitionBodyPlugin_SkipsIfDefinitionBodyExists(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	existingDefBody := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title": "My Custom API",
		},
		"paths": map[string]interface{}{},
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
	if defBody["swagger"] != "2.0" {
		t.Errorf("Expected swagger 2.0, got %v", defBody["swagger"])
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

func TestDefaultDefinitionBodyPlugin_CollectsRoutesFromFunctionEvents(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"ServerlessRestApi": {
				Type: "AWS::Serverless::Api",
				Properties: map[string]interface{}{
					"StageName": "Prod",
				},
			},
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"GetUsers": map[string]interface{}{
							"Type": "Api",
							"Properties": map[string]interface{}{
								"Path":   "/users",
								"Method": "GET",
							},
						},
						"CreateUser": map[string]interface{}{
							"Type": "Api",
							"Properties": map[string]interface{}{
								"Path":   "/users",
								"Method": "POST",
							},
						},
					},
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	api := template.Resources["ServerlessRestApi"]
	defBody, ok := api.Properties["DefinitionBody"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected DefinitionBody to be set")
	}

	paths, ok := defBody["paths"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected paths to be set")
	}

	usersPath, ok := paths["/users"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected /users path to be set")
	}

	if _, exists := usersPath["get"]; !exists {
		t.Errorf("Expected GET method on /users")
	}

	if _, exists := usersPath["post"]; !exists {
		t.Errorf("Expected POST method on /users")
	}

	// Verify x-amazon-apigateway-integration is set
	getMethod, ok := usersPath["get"].(map[string]interface{})
	if ok {
		if _, hasIntegration := getMethod["x-amazon-apigateway-integration"]; !hasIntegration {
			t.Errorf("Expected x-amazon-apigateway-integration on GET method")
		}
	}
}

func TestDefaultDefinitionBodyPlugin_CollectsRoutesForExplicitApi(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyCustomApi": {
				Type: "AWS::Serverless::Api",
				Properties: map[string]interface{}{
					"StageName": "prod",
				},
			},
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"GetData": map[string]interface{}{
							"Type": "Api",
							"Properties": map[string]interface{}{
								"Path":      "/data",
								"Method":    "GET",
								"RestApiId": map[string]interface{}{"Ref": "MyCustomApi"},
							},
						},
					},
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	api := template.Resources["MyCustomApi"]
	defBody, ok := api.Properties["DefinitionBody"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected DefinitionBody to be set")
	}

	paths, ok := defBody["paths"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected paths to be set")
	}

	if _, ok := paths["/data"]; !ok {
		t.Errorf("Expected /data path to be set")
	}
}

func TestDefaultDefinitionBodyPlugin_HttpApiRoutesWithPayloadFormat(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"ServerlessHttpApi": {
				Type: "AWS::Serverless::HttpApi",
				Properties: map[string]interface{}{
					"StageName": "$default",
				},
			},
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"GetItems": map[string]interface{}{
							"Type": "HttpApi",
							"Properties": map[string]interface{}{
								"Path":                 "/items",
								"Method":               "GET",
								"PayloadFormatVersion": "1.0",
							},
						},
					},
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	api := template.Resources["ServerlessHttpApi"]
	defBody, ok := api.Properties["DefinitionBody"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected DefinitionBody to be set")
	}

	// Should be OpenAPI 3.0
	if defBody["openapi"] != "3.0.1" {
		t.Errorf("Expected openapi 3.0.1, got %v", defBody["openapi"])
	}

	paths, ok := defBody["paths"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected paths to be set")
	}

	if _, ok := paths["/items"]; !ok {
		t.Errorf("Expected /items path to be set")
	}
}

func TestDefaultDefinitionBodyPlugin_MergesRoutesIntoExistingSpec(t *testing.T) {
	plugin := NewDefaultDefinitionBodyPlugin()

	existingDefBody := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":   "My API",
			"version": "1.0",
		},
		"paths": map[string]interface{}{
			"/existing": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Existing endpoint",
				},
			},
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
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"NewEndpoint": map[string]interface{}{
							"Type": "Api",
							"Properties": map[string]interface{}{
								"Path":      "/new",
								"Method":    "POST",
								"RestApiId": map[string]interface{}{"Ref": "MyApi"},
							},
						},
					},
				},
			},
		},
	}

	err := plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	api := template.Resources["MyApi"]
	defBody := api.Properties["DefinitionBody"].(map[string]interface{})

	// Should preserve original info
	info := defBody["info"].(map[string]interface{})
	if info["title"] != "My API" {
		t.Errorf("Expected title 'My API', got %v", info["title"])
	}

	paths := defBody["paths"].(map[string]interface{})

	// Should preserve existing endpoint
	if _, ok := paths["/existing"]; !ok {
		t.Errorf("Expected /existing path to be preserved")
	}

	// Should add new endpoint
	if _, ok := paths["/new"]; !ok {
		t.Errorf("Expected /new path to be added")
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
