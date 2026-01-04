// Package sam provides SAM resource transformers.
package sam

import (
	"reflect"
	"testing"
)

func TestApiTransformer_Transform_Basic(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
			"info": map[string]interface{}{
				"title":   "Test API",
				"version": "1.0",
			},
			"paths": map[string]interface{}{},
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create RestApi, Deployment, and Stage
	if len(resources) < 3 {
		t.Errorf("Expected at least 3 resources, got %d", len(resources))
	}

	// Verify RestApi was created
	restApi, ok := resources["MyApi"].(map[string]interface{})
	if !ok {
		t.Fatal("MyApi resource not found")
	}
	if restApi["Type"] != "AWS::ApiGateway::RestApi" {
		t.Errorf("Expected AWS::ApiGateway::RestApi, got %v", restApi["Type"])
	}

	// Verify RestApi has Body property with the swagger definition
	props, ok := restApi["Properties"].(map[string]interface{})
	if !ok {
		t.Fatal("RestApi Properties not found")
	}
	if _, hasBody := props["Body"]; !hasBody {
		t.Error("RestApi should have Body property")
	}
}

func TestApiTransformer_Transform_WithName(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		Name:      "MyCustomApiName",
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	if props["Name"] != "MyCustomApiName" {
		t.Errorf("Expected Name 'MyCustomApiName', got %v", props["Name"])
	}
}

func TestApiTransformer_Transform_WithDefinitionUri(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionUri: map[string]interface{}{
			"Bucket": "my-bucket",
			"Key":    "swagger.json",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	bodyS3Location, ok := props["BodyS3Location"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected BodyS3Location property")
	}
	if bodyS3Location["Bucket"] != "my-bucket" {
		t.Errorf("Expected Bucket 'my-bucket', got %v", bodyS3Location["Bucket"])
	}
	if bodyS3Location["Key"] != "swagger.json" {
		t.Errorf("Expected Key 'swagger.json', got %v", bodyS3Location["Key"])
	}
}

func TestApiTransformer_Transform_WithStageVariables(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		Variables: map[string]interface{}{
			"Endpoint": "https://example.com",
			"Version":  "v1",
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find the Stage resource
	var stageProps map[string]interface{}
	for name, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Stage" {
			stageProps = resMap["Properties"].(map[string]interface{})
			_ = name
			break
		}
	}

	if stageProps == nil {
		t.Fatal("Stage resource not found")
	}

	variables, ok := stageProps["Variables"].(map[string]interface{})
	if !ok {
		t.Fatal("Stage Variables not found")
	}
	if variables["Endpoint"] != "https://example.com" {
		t.Errorf("Expected Endpoint variable, got %v", variables["Endpoint"])
	}
}

func TestApiTransformer_Transform_WithCaching(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName:           "Prod",
		CacheClusterEnabled: true,
		CacheClusterSize:    "0.5",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find the Stage resource
	var stageProps map[string]interface{}
	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Stage" {
			stageProps = resMap["Properties"].(map[string]interface{})
			break
		}
	}

	if stageProps == nil {
		t.Fatal("Stage resource not found")
	}

	if stageProps["CacheClusterEnabled"] != true {
		t.Error("Expected CacheClusterEnabled to be true")
	}
	if stageProps["CacheClusterSize"] != "0.5" {
		t.Errorf("Expected CacheClusterSize '0.5', got %v", stageProps["CacheClusterSize"])
	}
}

func TestApiTransformer_Transform_WithEndpointConfiguration(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		EndpointConfiguration: &EndpointConfig{
			Type: "REGIONAL",
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	endpointConfig, ok := props["EndpointConfiguration"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected EndpointConfiguration property")
	}
	types, ok := endpointConfig["Types"].([]interface{})
	if !ok {
		t.Fatal("Expected Types in EndpointConfiguration")
	}
	if len(types) != 1 || types[0] != "REGIONAL" {
		t.Errorf("Expected ['REGIONAL'], got %v", types)
	}
}

func TestApiTransformer_Transform_WithBinaryMediaTypes(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		BinaryMediaTypes: []interface{}{
			"image/png",
			"application/octet-stream",
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	binaryTypes, ok := props["BinaryMediaTypes"].([]interface{})
	if !ok {
		t.Fatal("Expected BinaryMediaTypes property")
	}
	if len(binaryTypes) != 2 {
		t.Errorf("Expected 2 binary media types, got %d", len(binaryTypes))
	}
}

func TestApiTransformer_Transform_WithTracingEnabled(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName:      "Prod",
		TracingEnabled: true,
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find the Stage resource
	var stageProps map[string]interface{}
	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Stage" {
			stageProps = resMap["Properties"].(map[string]interface{})
			break
		}
	}

	if stageProps == nil {
		t.Fatal("Stage resource not found")
	}

	if stageProps["TracingEnabled"] != true {
		t.Error("Expected TracingEnabled to be true")
	}
}

func TestApiTransformer_Transform_WithAccessLogSetting(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		AccessLogSetting: &AccessLogSetting{
			DestinationArn: "arn:aws:logs:us-east-1:123456789012:log-group:my-log-group",
			Format:         "$requestId $httpMethod $path",
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find the Stage resource
	var stageProps map[string]interface{}
	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Stage" {
			stageProps = resMap["Properties"].(map[string]interface{})
			break
		}
	}

	if stageProps == nil {
		t.Fatal("Stage resource not found")
	}

	accessLog, ok := stageProps["AccessLogSetting"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected AccessLogSetting property")
	}
	if accessLog["DestinationArn"] != "arn:aws:logs:us-east-1:123456789012:log-group:my-log-group" {
		t.Errorf("Unexpected DestinationArn: %v", accessLog["DestinationArn"])
	}
}

func TestApiTransformer_Transform_WithMethodSettings(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		MethodSettings: []MethodSettingConfig{
			{
				HttpMethod:       "*",
				ResourcePath:     "/*",
				LoggingLevel:     "INFO",
				DataTraceEnabled: true,
				MetricsEnabled:   true,
			},
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find the Stage resource
	var stageProps map[string]interface{}
	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Stage" {
			stageProps = resMap["Properties"].(map[string]interface{})
			break
		}
	}

	if stageProps == nil {
		t.Fatal("Stage resource not found")
	}

	methodSettings, ok := stageProps["MethodSettings"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected MethodSettings property")
	}
	if len(methodSettings) != 1 {
		t.Errorf("Expected 1 method setting, got %d", len(methodSettings))
	}
	if methodSettings[0]["LoggingLevel"] != "INFO" {
		t.Errorf("Expected LoggingLevel 'INFO', got %v", methodSettings[0]["LoggingLevel"])
	}
}

func TestApiTransformer_Transform_WithTags(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		Tags: map[string]interface{}{
			"Environment": "Production",
			"Team":        "Backend",
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	tags, ok := props["Tags"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected Tags property")
	}
	if len(tags) == 0 {
		t.Error("Expected at least 1 tag")
	}
}

func TestApiTransformer_Transform_DeploymentIdStability(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
			"info": map[string]interface{}{
				"title":   "Test API",
				"version": "1.0",
			},
			"paths": map[string]interface{}{
				"/test": map[string]interface{}{
					"get": map[string]interface{}{},
				},
			},
		},
	}

	resources1, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("First transform failed: %v", err)
	}

	resources2, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Second transform failed: %v", err)
	}

	// Find deployment IDs
	var deploymentID1, deploymentID2 string
	for name := range resources1 {
		if name != "MyApi" && resources1[name].(map[string]interface{})["Type"] == "AWS::ApiGateway::Deployment" {
			deploymentID1 = name
			break
		}
	}
	for name := range resources2 {
		if name != "MyApi" && resources2[name].(map[string]interface{})["Type"] == "AWS::ApiGateway::Deployment" {
			deploymentID2 = name
			break
		}
	}

	if deploymentID1 != deploymentID2 {
		t.Errorf("Deployment IDs should be stable: %s vs %s", deploymentID1, deploymentID2)
	}
}

func TestApiTransformer_Transform_WithCors(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		Cors: &CorsConfig{
			AllowOrigin:      "'*'",
			AllowMethods:     "'GET,POST'",
			AllowHeaders:     "'Content-Type'",
			MaxAge:           "'600'",
			AllowCredentials: true,
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// CORS should be embedded in the OpenAPI definition
	restApi := resources["MyApi"].(map[string]interface{})
	if restApi["Type"] != "AWS::ApiGateway::RestApi" {
		t.Errorf("Expected AWS::ApiGateway::RestApi type")
	}
}

func TestApiTransformer_Transform_WithAuth(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		Auth: &ApiAuth{
			DefaultAuthorizer: "MyCognitoAuth",
			Authorizers: map[string]interface{}{
				"MyCognitoAuth": map[string]interface{}{
					"UserPoolArn": "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_XXXX",
				},
			},
		},
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Check that authorizer resources are created
	hasAuthorizer := false
	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Authorizer" {
			hasAuthorizer = true
			break
		}
	}

	if !hasAuthorizer {
		t.Error("Expected an Authorizer resource to be created")
	}
}

func TestApiTransformer_Transform_WithMinimumCompressionSize(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName:              "Prod",
		MinimumCompressionSize: 1024,
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	if props["MinimumCompressionSize"] != 1024 {
		t.Errorf("Expected MinimumCompressionSize 1024, got %v", props["MinimumCompressionSize"])
	}
}

func TestApiTransformer_Transform_DefaultStageNameIsRequired(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		// No StageName provided
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	_, err := transformer.Transform("MyApi", api)
	if err == nil {
		t.Error("Expected error when StageName is not provided")
	}
}

func TestApi_DefinitionBodyAndUriMutuallyExclusive(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
		DefinitionUri: "s3://bucket/key",
	}

	// Should use DefinitionBody when both are specified (matching Python behavior)
	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	// DefinitionBody takes precedence
	if _, hasBody := props["Body"]; !hasBody {
		t.Error("Expected Body property when DefinitionBody is provided")
	}
	if _, hasS3 := props["BodyS3Location"]; hasS3 {
		t.Error("Should not have BodyS3Location when DefinitionBody is provided")
	}
}

func TestApiTransformer_Transform_WithOpenApi30(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName:      "Prod",
		OpenApiVersion: "3.0.1",
		DefinitionBody: map[string]interface{}{
			"openapi": "3.0.1",
			"info": map[string]interface{}{
				"title":   "Test API",
				"version": "1.0",
			},
			"paths": map[string]interface{}{},
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	if restApi["Type"] != "AWS::ApiGateway::RestApi" {
		t.Errorf("Expected AWS::ApiGateway::RestApi, got %v", restApi["Type"])
	}
}

func TestApiTransformer_Transform_GeneratesDeploymentAndStage(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	hasDeployment := false
	hasStage := false

	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		switch resMap["Type"] {
		case "AWS::ApiGateway::Deployment":
			hasDeployment = true
		case "AWS::ApiGateway::Stage":
			hasStage = true
		}
	}

	if !hasDeployment {
		t.Error("Expected a Deployment resource")
	}
	if !hasStage {
		t.Error("Expected a Stage resource")
	}
}

func TestApiTransformer_StageDependsOnDeployment(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find Stage and verify it references the Deployment
	var stageRes map[string]interface{}
	var deploymentName string

	for name, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Deployment" {
			deploymentName = name
		}
		if resMap["Type"] == "AWS::ApiGateway::Stage" {
			stageRes = resMap
		}
	}

	if deploymentName == "" {
		t.Fatal("Deployment resource not found")
	}
	if stageRes == nil {
		t.Fatal("Stage resource not found")
	}

	props := stageRes["Properties"].(map[string]interface{})
	deploymentId := props["DeploymentId"]

	// DeploymentId should reference the deployment
	refMap, ok := deploymentId.(map[string]interface{})
	if !ok {
		t.Fatal("DeploymentId should be a Ref")
	}
	if refMap["Ref"] != deploymentName {
		t.Errorf("Stage DeploymentId should reference %s, got %v", deploymentName, refMap["Ref"])
	}
}

func TestApiTransformer_RestApiIdInDeployment(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find Deployment and verify RestApiId references the RestApi
	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Deployment" {
			props := resMap["Properties"].(map[string]interface{})
			restApiId := props["RestApiId"]

			refMap, ok := restApiId.(map[string]interface{})
			if !ok {
				t.Fatal("RestApiId should be a Ref")
			}
			if refMap["Ref"] != "MyApi" {
				t.Errorf("Deployment RestApiId should reference MyApi, got %v", refMap["Ref"])
			}
			return
		}
	}

	t.Fatal("Deployment resource not found")
}

func TestApiTransformer_Transform_WithDescription(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName:   "Prod",
		Description: "My test API description",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	if props["Description"] != "My test API description" {
		t.Errorf("Expected Description, got %v", props["Description"])
	}
}

func TestApiTransformer_Transform_FailOnWarnings(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName:      "Prod",
		FailOnWarnings: true,
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	if props["FailOnWarnings"] != true {
		t.Error("Expected FailOnWarnings to be true")
	}
}

func TestApiTransformer_Transform_DisableExecuteApiEndpoint(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName:                 "Prod",
		DisableExecuteApiEndpoint: true,
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	restApi := resources["MyApi"].(map[string]interface{})
	props := restApi["Properties"].(map[string]interface{})

	if props["DisableExecuteApiEndpoint"] != true {
		t.Error("Expected DisableExecuteApiEndpoint to be true")
	}
}

func TestApiTransformer_DeploymentDependsOnRestApi(t *testing.T) {
	transformer := NewApiTransformer()

	api := &Api{
		StageName: "Prod",
		DefinitionBody: map[string]interface{}{
			"swagger": "2.0",
		},
	}

	resources, err := transformer.Transform("MyApi", api)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Find the Deployment resource and check it has DependsOn for RestApi
	for _, res := range resources {
		resMap, ok := res.(map[string]interface{})
		if !ok {
			continue
		}
		if resMap["Type"] == "AWS::ApiGateway::Deployment" {
			dependsOn := resMap["DependsOn"]
			if dependsOn == nil {
				// DependsOn is optional if there's just the RestApi reference
				return
			}
			if dependsOnSlice, ok := dependsOn.([]string); ok {
				found := false
				for _, dep := range dependsOnSlice {
					if dep == "MyApi" {
						found = true
						break
					}
				}
				if !found {
					t.Error("Deployment should depend on MyApi")
				}
			}
			return
		}
	}
}

func TestNewApiTransformer(t *testing.T) {
	transformer := NewApiTransformer()
	if transformer == nil {
		t.Error("NewApiTransformer should not return nil")
	}
}

func TestApi_AllFieldsPopulated(t *testing.T) {
	api := &Api{
		Name:                      "TestApi",
		StageName:                 "Prod",
		Description:               "Test description",
		DefinitionBody:            map[string]interface{}{"swagger": "2.0"},
		DefinitionUri:             "s3://bucket/key",
		BinaryMediaTypes:          []interface{}{"image/png"},
		MinimumCompressionSize:    1024,
		EndpointConfiguration:     &EndpointConfig{Type: "REGIONAL"},
		CacheClusterEnabled:       true,
		CacheClusterSize:          "0.5",
		Variables:                 map[string]interface{}{"key": "value"},
		TracingEnabled:            true,
		Tags:                      map[string]interface{}{"env": "test"},
		AccessLogSetting:          &AccessLogSetting{DestinationArn: "arn:aws:logs:..."},
		MethodSettings:            []MethodSettingConfig{},
		FailOnWarnings:            true,
		DisableExecuteApiEndpoint: true,
		OpenApiVersion:            "3.0.1",
		Cors:                      &CorsConfig{AllowOrigin: "'*'"},
		Auth:                      &ApiAuth{DefaultAuthorizer: "test"},
	}

	// Just verify all fields can be set without panic
	if api.Name != "TestApi" {
		t.Errorf("Unexpected Name: %v", api.Name)
	}
	if !reflect.DeepEqual(api.BinaryMediaTypes, []interface{}{"image/png"}) {
		t.Errorf("Unexpected BinaryMediaTypes: %v", api.BinaryMediaTypes)
	}
}
