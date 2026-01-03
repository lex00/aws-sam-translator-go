package plugins

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestImplicitHttpApiPlugin_Name(t *testing.T) {
	plugin := NewImplicitHttpApiPlugin()
	if plugin.Name() != "ImplicitHttpApiPlugin" {
		t.Errorf("Expected name 'ImplicitHttpApiPlugin', got '%s'", plugin.Name())
	}
}

func TestImplicitHttpApiPlugin_Priority(t *testing.T) {
	plugin := NewImplicitHttpApiPlugin()
	if plugin.Priority() != 310 {
		t.Errorf("Expected priority 310, got %d", plugin.Priority())
	}
}

func TestImplicitHttpApiPlugin_CreatesImplicitHttpApi(t *testing.T) {
	plugin := NewImplicitHttpApiPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"GetApi": map[string]interface{}{
							"Type": "HttpApi",
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

	// Should create ServerlessHttpApi
	api, ok := template.Resources["ServerlessHttpApi"]
	if !ok {
		t.Fatalf("Expected ServerlessHttpApi to be created")
	}

	if api.Type != "AWS::Serverless::HttpApi" {
		t.Errorf("Expected type AWS::Serverless::HttpApi, got %s", api.Type)
	}

	if api.Properties["StageName"] != "$default" {
		t.Errorf("Expected StageName '$default', got %v", api.Properties["StageName"])
	}
}

func TestImplicitHttpApiPlugin_SkipsIfApiIdSpecified(t *testing.T) {
	plugin := NewImplicitHttpApiPlugin()

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyHttpApi": {
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
						"GetApi": map[string]interface{}{
							"Type": "HttpApi",
							"Properties": map[string]interface{}{
								"Path":   "/hello",
								"Method": "get",
								"ApiId":  map[string]interface{}{"Ref": "MyHttpApi"},
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

	// Should not create ServerlessHttpApi
	_, ok := template.Resources["ServerlessHttpApi"]
	if ok {
		t.Errorf("Should not create ServerlessHttpApi when ApiId is specified")
	}
}

func TestImplicitHttpApiPlugin_SkipsIfNoHttpApiEvents(t *testing.T) {
	plugin := NewImplicitHttpApiPlugin()

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

	// Should not create ServerlessHttpApi
	_, ok := template.Resources["ServerlessHttpApi"]
	if ok {
		t.Errorf("Should not create ServerlessHttpApi when no HttpApi events exist")
	}
}

func TestImplicitHttpApiPlugin_SkipsIfAlreadyExists(t *testing.T) {
	plugin := NewImplicitHttpApiPlugin()

	existingApi := types.Resource{
		Type: "AWS::Serverless::HttpApi",
		Properties: map[string]interface{}{
			"StageName": "custom",
		},
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"ServerlessHttpApi": existingApi,
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Events": map[string]interface{}{
						"GetApi": map[string]interface{}{
							"Type": "HttpApi",
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

	// Should keep existing ServerlessHttpApi
	api := template.Resources["ServerlessHttpApi"]
	if api.Properties["StageName"] != "custom" {
		t.Errorf("Expected StageName 'custom', got %v", api.Properties["StageName"])
	}
}

func TestImplicitHttpApiPlugin_AfterTransform(t *testing.T) {
	plugin := NewImplicitHttpApiPlugin()
	template := &types.Template{}

	err := plugin.AfterTransform(template)
	if err != nil {
		t.Errorf("AfterTransform should not error, got: %v", err)
	}
}
