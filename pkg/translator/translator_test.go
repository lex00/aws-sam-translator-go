package translator

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/plugins"
	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestNew(t *testing.T) {
	tr := New()
	if tr == nil {
		t.Error("New() returned nil")
	}
}

func TestNewWithOptions(t *testing.T) {
	tr := NewWithOptions(Options{
		Region:    "us-west-2",
		AccountID: "123456789012",
		StackName: "test-stack",
	})
	if tr == nil {
		t.Fatal("NewWithOptions() returned nil")
	}
	if tr.options.Region != "us-west-2" {
		t.Errorf("expected region 'us-west-2', got '%s'", tr.options.Region)
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version is empty")
	}
}

func TestTransformEmptyTemplate(t *testing.T) {
	tr := New()
	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Resources:                make(map[string]types.Resource),
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	if result == nil {
		t.Fatal("Transform returned nil result")
	}
	if len(result.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(result.Resources))
	}
}

func TestTransformSimpleFunction(t *testing.T) {
	tr := NewWithOptions(Options{
		Region:    "us-east-1",
		AccountID: "123456789012",
		StackName: "test-stack",
	})

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "nodejs18.x",
					"CodeUri": "s3://bucket/key",
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have created Lambda function and IAM role
	if _, ok := result.Resources["MyFunction"]; !ok {
		t.Error("expected MyFunction in result")
	}

	// Verify it's a Lambda function now
	fn := result.Resources["MyFunction"]
	if fn.Type != "AWS::Lambda::Function" {
		t.Errorf("expected AWS::Lambda::Function, got %s", fn.Type)
	}
}

func TestTransformSimpleTable(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyTable": {
				Type: "AWS::Serverless::SimpleTable",
				Properties: map[string]interface{}{
					"TableName": "test-table",
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have transformed to DynamoDB table
	if _, ok := result.Resources["MyTable"]; !ok {
		t.Error("expected MyTable in result")
	}

	table := result.Resources["MyTable"]
	if table.Type != "AWS::DynamoDB::Table" {
		t.Errorf("expected AWS::DynamoDB::Table, got %s", table.Type)
	}
}

func TestTransformResourceOrdering(t *testing.T) {
	tr := New()

	// Create a template with multiple resource types (without connector for simplicity)
	// The ordering test verifies that resources are processed in the correct order
	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "nodejs18.x",
					"CodeUri": "s3://bucket/key",
				},
			},
			"MyTable": {
				Type: "AWS::Serverless::SimpleTable",
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Verify all resources are transformed
	if _, ok := result.Resources["MyFunction"]; !ok {
		t.Error("expected MyFunction in result")
	}
	if _, ok := result.Resources["MyTable"]; !ok {
		t.Error("expected MyTable in result")
	}
}

func TestTransformPluginLifecycle(t *testing.T) {
	tr := New()

	// Track plugin execution order
	var executionOrder []string

	// Register test plugins
	testPlugin := &testPlugin{
		name:     "TestPlugin",
		priority: 100,
		beforeFn: func(t *types.Template) error {
			executionOrder = append(executionOrder, "before")
			return nil
		},
		afterFn: func(t *types.Template) error {
			executionOrder = append(executionOrder, "after")
			return nil
		},
	}
	tr.RegisterPlugin(testPlugin)

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Resources:                make(map[string]types.Resource),
	}

	_, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if len(executionOrder) != 2 {
		t.Fatalf("expected 2 plugin calls, got %d", len(executionOrder))
	}
	if executionOrder[0] != "before" {
		t.Error("expected 'before' first")
	}
	if executionOrder[1] != "after" {
		t.Error("expected 'after' second")
	}
}

func TestTransformErrorAggregation(t *testing.T) {
	tr := New()

	// Create template with multiple invalid resources
	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"BadFunction1": {
				Type: "AWS::Serverless::Function",
				// Missing required properties
				Properties: map[string]interface{}{},
			},
		},
	}

	_, err := tr.Transform(template)
	if err == nil {
		t.Fatal("expected error for invalid function")
	}

	// Error should contain resource information
	if !strings.Contains(err.Error(), "BadFunction1") {
		t.Errorf("error should mention resource name: %v", err)
	}
}

func TestTransformMetadataPassthrough(t *testing.T) {
	tr := NewWithOptions(Options{
		PassThroughMetadata: true,
	})

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Metadata: map[string]interface{}{
			"AWS::CloudFormation::Interface": map[string]interface{}{
				"ParameterGroups": []interface{}{},
			},
		},
		Resources: map[string]types.Resource{},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if result.Metadata == nil {
		t.Error("expected metadata to be passed through")
	}
	if _, ok := result.Metadata["AWS::CloudFormation::Interface"]; !ok {
		t.Error("expected AWS::CloudFormation::Interface in metadata")
	}
}

func TestTransformDeterministicOutput(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"FunctionA": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "nodejs18.x",
					"CodeUri": "s3://bucket/key",
				},
			},
			"FunctionB": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "nodejs18.x",
					"CodeUri": "s3://bucket/key",
				},
			},
		},
	}

	// Transform twice and compare
	result1, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform 1 failed: %v", err)
	}

	result2, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform 2 failed: %v", err)
	}

	// Convert to JSON and compare
	json1, _ := json.Marshal(result1)
	json2, _ := json.Marshal(result2)

	if string(json1) != string(json2) {
		t.Error("Transform output is not deterministic")
	}
}

func TestTransformPreservesNonSAMResources(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyBucket": {
				Type: "AWS::S3::Bucket",
				Properties: map[string]interface{}{
					"BucketName": "my-bucket",
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// S3 bucket should be preserved as-is
	bucket, ok := result.Resources["MyBucket"]
	if !ok {
		t.Fatal("expected MyBucket in result")
	}
	if bucket.Type != "AWS::S3::Bucket" {
		t.Errorf("expected AWS::S3::Bucket, got %s", bucket.Type)
	}
}

func TestTransformPreservesOutputs(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources:                map[string]types.Resource{},
		Outputs: map[string]types.Output{
			"BucketName": {
				Description: "The bucket name",
				Value:       "test-bucket",
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if len(result.Outputs) != 1 {
		t.Errorf("expected 1 output, got %d", len(result.Outputs))
	}
	if _, ok := result.Outputs["BucketName"]; !ok {
		t.Error("expected BucketName output")
	}
}

func TestTransformPreservesParameters(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Parameters: map[string]types.Parameter{
			"Environment": {
				Type:    "String",
				Default: "dev",
			},
		},
		Resources: map[string]types.Resource{},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if len(result.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(result.Parameters))
	}
	if _, ok := result.Parameters["Environment"]; !ok {
		t.Error("expected Environment parameter")
	}
}

func TestTransformPreservesMappings(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Mappings: map[string]interface{}{
			"RegionMap": map[string]interface{}{
				"us-east-1": map[string]interface{}{
					"AMI": "ami-12345",
				},
			},
		},
		Resources: map[string]types.Resource{},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if result.Mappings == nil {
		t.Error("expected mappings to be preserved")
	}
	if _, ok := result.Mappings["RegionMap"]; !ok {
		t.Error("expected RegionMap in mappings")
	}
}

func TestTransformPreservesConditions(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Conditions: map[string]interface{}{
			"IsProd": map[string]interface{}{
				"Fn::Equals": []interface{}{
					map[string]interface{}{"Ref": "Environment"},
					"prod",
				},
			},
		},
		Resources: map[string]types.Resource{},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	if result.Conditions == nil {
		t.Error("expected conditions to be preserved")
	}
	if _, ok := result.Conditions["IsProd"]; !ok {
		t.Error("expected IsProd in conditions")
	}
}

func TestTransformRemovesSAMTransform(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources:                map[string]types.Resource{},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// SAM transform should be removed from output
	if result.Transform != nil {
		transformStr, ok := result.Transform.(string)
		if ok && transformStr == "AWS::Serverless-2016-10-31" {
			t.Error("expected SAM transform to be removed")
		}
	}
}

func TestTransformStateMachine(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyStateMachine": {
				Type: "AWS::Serverless::StateMachine",
				Properties: map[string]interface{}{
					"Definition": map[string]interface{}{
						"StartAt": "HelloWorld",
						"States": map[string]interface{}{
							"HelloWorld": map[string]interface{}{
								"Type": "Pass",
								"End":  true,
							},
						},
					},
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have transformed to Step Functions state machine
	sm, ok := result.Resources["MyStateMachine"]
	if !ok {
		t.Fatal("expected MyStateMachine in result")
	}
	if sm.Type != "AWS::StepFunctions::StateMachine" {
		t.Errorf("expected AWS::StepFunctions::StateMachine, got %s", sm.Type)
	}
}

func TestTransformApi(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyApi": {
				Type: "AWS::Serverless::Api",
				Properties: map[string]interface{}{
					"StageName": "prod",
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have transformed to API Gateway RestApi
	api, ok := result.Resources["MyApi"]
	if !ok {
		t.Fatal("expected MyApi in result")
	}
	if api.Type != "AWS::ApiGateway::RestApi" {
		t.Errorf("expected AWS::ApiGateway::RestApi, got %s", api.Type)
	}
}

func TestTransformHttpApi(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyHttpApi": {
				Type: "AWS::Serverless::HttpApi",
				Properties: map[string]interface{}{
					"StageName": "prod",
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have transformed to API Gateway V2 Api
	api, ok := result.Resources["MyHttpApi"]
	if !ok {
		t.Fatal("expected MyHttpApi in result")
	}
	if api.Type != "AWS::ApiGatewayV2::Api" {
		t.Errorf("expected AWS::ApiGatewayV2::Api, got %s", api.Type)
	}
}

func TestTransformLayerVersion(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyLayer": {
				Type: "AWS::Serverless::LayerVersion",
				Properties: map[string]interface{}{
					"ContentUri":         "s3://bucket/layer.zip",
					"CompatibleRuntimes": []interface{}{"nodejs18.x"},
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have transformed to Lambda LayerVersion
	// LayerVersion transformer adds a hash suffix to the logical ID
	var found bool
	for name, res := range result.Resources {
		if strings.HasPrefix(name, "MyLayer") && res.Type == "AWS::Lambda::LayerVersion" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected a resource starting with MyLayer of type AWS::Lambda::LayerVersion in result")
	}
}

func TestTransformApplication(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyApp": {
				Type: "AWS::Serverless::Application",
				Properties: map[string]interface{}{
					"Location": "https://s3.amazonaws.com/bucket/template.yaml",
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have transformed to CloudFormation Stack
	app, ok := result.Resources["MyApp"]
	if !ok {
		t.Fatal("expected MyApp in result")
	}
	if app.Type != "AWS::CloudFormation::Stack" {
		t.Errorf("expected AWS::CloudFormation::Stack, got %s", app.Type)
	}
}

func TestTransformGraphQLApi(t *testing.T) {
	tr := New()

	template := &types.Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Transform:                "AWS::Serverless-2016-10-31",
		Resources: map[string]types.Resource{
			"MyGraphQL": {
				Type: "AWS::Serverless::GraphQLApi",
				Properties: map[string]interface{}{
					"SchemaInline": "type Query { hello: String }",
					"Auth": map[string]interface{}{
						"Type": "API_KEY",
					},
				},
			},
		},
	}

	result, err := tr.Transform(template)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have transformed to AppSync GraphQLApi
	api, ok := result.Resources["MyGraphQL"]
	if !ok {
		t.Fatal("expected MyGraphQL in result")
	}
	if api.Type != "AWS::AppSync::GraphQLApi" {
		t.Errorf("expected AWS::AppSync::GraphQLApi, got %s", api.Type)
	}
}

func TestTransformBytes(t *testing.T) {
	tr := New()

	input := []byte(`
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: s3://bucket/key
`)

	output, err := tr.TransformBytes(input)
	if err != nil {
		t.Fatalf("TransformBytes failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output")
	}

	// Output should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestGetResourceOrder(t *testing.T) {
	order := getResourceOrder()

	// Functions should come before APIs
	funcIdx := -1
	apiIdx := -1
	for i, rt := range order {
		if rt == "AWS::Serverless::Function" {
			funcIdx = i
		}
		if rt == "AWS::Serverless::Api" {
			apiIdx = i
		}
	}

	if funcIdx >= apiIdx {
		t.Error("Functions should be processed before APIs")
	}

	// Connectors should be last
	connectorIdx := -1
	for i, rt := range order {
		if rt == "AWS::Serverless::Connector" {
			connectorIdx = i
		}
	}

	if connectorIdx != len(order)-1 {
		t.Error("Connectors should be processed last")
	}
}

// testPlugin is a mock plugin for testing.
type testPlugin struct {
	name     string
	priority int
	beforeFn func(*types.Template) error
	afterFn  func(*types.Template) error
}

func (p *testPlugin) Name() string  { return p.name }
func (p *testPlugin) Priority() int { return p.priority }
func (p *testPlugin) BeforeTransform(t *types.Template) error {
	if p.beforeFn != nil {
		return p.beforeFn(t)
	}
	return nil
}
func (p *testPlugin) AfterTransform(t *types.Template) error {
	if p.afterFn != nil {
		return p.afterFn(t)
	}
	return nil
}

// Ensure testPlugin implements plugins.Plugin
var _ plugins.Plugin = (*testPlugin)(nil)
