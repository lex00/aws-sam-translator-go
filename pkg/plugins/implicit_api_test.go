package plugins

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestImplicitRestApiPlugin_Name(t *testing.T) {
	plugin := NewImplicitRestApiPlugin()
	if plugin.Name() != "ImplicitRestApiPlugin" {
		t.Errorf("Expected name 'ImplicitRestApiPlugin', got '%s'", plugin.Name())
	}
}

func TestImplicitRestApiPlugin_Priority(t *testing.T) {
	plugin := NewImplicitRestApiPlugin()
	if plugin.Priority() != 300 {
		t.Errorf("Expected priority 300, got %d", plugin.Priority())
	}
}

func TestImplicitRestApiPlugin_CreatesImplicitApi(t *testing.T) {
	plugin := NewImplicitRestApiPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"GetApi": map[string]interface{}{
							"Type": "Api",
							"Properties": map[string]interface{}{
								"Path":   "/hello",
								"Method": "get",
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

	// Should create ServerlessRestApi
	api, ok := template.Resources["ServerlessRestApi"]
	if !ok {
		t.Fatalf("Expected ServerlessRestApi to be created")
	}

	if api.Type != "AWS::Serverless::Api" {
		t.Errorf("Expected type AWS::Serverless::Api, got %s", api.Type)
	}

	if api.Properties["StageName"] != "Prod" {
		t.Errorf("Expected StageName 'Prod', got %v", api.Properties["StageName"])
	}
}

func TestImplicitRestApiPlugin_SkipsIfRestApiIdSpecified(t *testing.T) {
	plugin := NewImplicitRestApiPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyApi": {
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
						"GetApi": map[string]interface{}{
							"Type": "Api",
							"Properties": map[string]interface{}{
								"Path":      "/hello",
								"Method":    "get",
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

	// Should not create ServerlessRestApi
	_, ok := template.Resources["ServerlessRestApi"]
	if ok {
		t.Errorf("Should not create ServerlessRestApi when RestApiId is specified")
	}
}

func TestImplicitRestApiPlugin_SkipsIfNoApiEvents(t *testing.T) {
	plugin := NewImplicitRestApiPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"MyEvent": map[string]interface{}{
							"Type": "S3",
							"Properties": map[string]interface{}{
								"Bucket": "my-bucket",
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

	// Should not create ServerlessRestApi
	_, ok := template.Resources["ServerlessRestApi"]
	if ok {
		t.Errorf("Should not create ServerlessRestApi when no Api events exist")
	}
}

func TestImplicitRestApiPlugin_SkipsIfAlreadyExists(t *testing.T) {
	plugin := NewImplicitRestApiPlugin()

	existingApi := types.Resource{
		Type: "AWS::Serverless::Api",
		Properties: map[string]interface{}{
			"StageName": "custom",
		},
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"ServerlessRestApi": existingApi,
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"GetApi": map[string]interface{}{
							"Type": "Api",
							"Properties": map[string]interface{}{
								"Path":   "/hello",
								"Method": "get",
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

	// Should keep existing ServerlessRestApi
	api := template.Resources["ServerlessRestApi"]
	if api.Properties["StageName"] != "custom" {
		t.Errorf("Expected StageName 'custom', got %v", api.Properties["StageName"])
	}
}

func TestImplicitRestApiPlugin_AfterTransform(t *testing.T) {
	plugin := NewImplicitRestApiPlugin()
	template := &types.Template{}

	err := plugin.AfterTransform(template)
	if err != nil {
		t.Errorf("AfterTransform should not error, got: %v", err)
	}
}
