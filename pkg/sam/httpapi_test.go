package sam

import (
	"testing"
)

func TestHttpApiTransformer_Transform_Minimal(t *testing.T) {
	transformer := NewHttpApiTransformer()

	// Minimal HttpApi with just defaults
	api := &HttpApi{}

	ctx := &TransformContext{
		Region:    "us-east-1",
		AccountID: "123456789012",
		StackName: "TestStack",
		Partition: "aws",
	}

	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have 2 resources: Api and Stage
	if len(resources) != 2 {
		t.Errorf("expected 2 resources (Api + Stage), got %d", len(resources))
	}

	// Check Api resource exists
	apiResource, ok := resources["MyHttpApi"].(map[string]interface{})
	if !ok {
		t.Fatal("MyHttpApi resource not found")
	}

	// Check type
	if apiResource["Type"] != "AWS::ApiGatewayV2::Api" {
		t.Errorf("expected Type 'AWS::ApiGatewayV2::Api', got %v", apiResource["Type"])
	}

	// Check properties
	props := apiResource["Properties"].(map[string]interface{})

	// ProtocolType should be HTTP
	if props["ProtocolType"] != "HTTP" {
		t.Errorf("expected ProtocolType 'HTTP', got %v", props["ProtocolType"])
	}

	// Name should default to logical ID
	if props["Name"] != "MyHttpApi" {
		t.Errorf("expected Name 'MyHttpApi', got %v", props["Name"])
	}

	// Check Stage resource exists
	stageResource, ok := resources["MyHttpApiStage"].(map[string]interface{})
	if !ok {
		t.Fatal("MyHttpApiStage resource not found")
	}

	if stageResource["Type"] != "AWS::ApiGatewayV2::Stage" {
		t.Errorf("expected Type 'AWS::ApiGatewayV2::Stage', got %v", stageResource["Type"])
	}

	stageProps := stageResource["Properties"].(map[string]interface{})

	// StageName should default to $default
	if stageProps["StageName"] != "$default" {
		t.Errorf("expected StageName '$default', got %v", stageProps["StageName"])
	}

	// AutoDeploy should be true
	if stageProps["AutoDeploy"] != true {
		t.Errorf("expected AutoDeploy true, got %v", stageProps["AutoDeploy"])
	}

	// ApiId should reference the API
	apiRef := stageProps["ApiId"].(map[string]interface{})
	if apiRef["Ref"] != "MyHttpApi" {
		t.Errorf("expected ApiId Ref 'MyHttpApi', got %v", apiRef["Ref"])
	}
}

func TestHttpApiTransformer_Transform_WithStageName(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		StageName: "prod",
		Name:      "MyProductionApi",
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stageResource := resources["MyHttpApiStage"].(map[string]interface{})
	stageProps := stageResource["Properties"].(map[string]interface{})

	if stageProps["StageName"] != "prod" {
		t.Errorf("expected StageName 'prod', got %v", stageProps["StageName"])
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	if apiProps["Name"] != "MyProductionApi" {
		t.Errorf("expected Name 'MyProductionApi', got %v", apiProps["Name"])
	}
}

func TestHttpApiTransformer_Transform_WithDescription(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Description: "My test HTTP API",
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	if apiProps["Description"] != "My test HTTP API" {
		t.Errorf("expected Description 'My test HTTP API', got %v", apiProps["Description"])
	}
}

func TestHttpApiTransformer_Transform_WithCorsBoolean(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		CorsConfiguration: true,
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	corsConfig := apiProps["CorsConfiguration"].(map[string]interface{})
	if corsConfig == nil {
		t.Fatal("CorsConfiguration should be present")
	}

	// Check AllowOrigins
	origins := corsConfig["AllowOrigins"].([]interface{})
	if len(origins) != 1 || origins[0] != "*" {
		t.Errorf("expected AllowOrigins ['*'], got %v", origins)
	}

	// Check AllowMethods
	methods := corsConfig["AllowMethods"].([]interface{})
	if len(methods) < 5 {
		t.Errorf("expected multiple AllowMethods, got %d", len(methods))
	}
}

func TestHttpApiTransformer_Transform_WithCorsString(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		CorsConfiguration: "https://example.com",
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	corsConfig := apiProps["CorsConfiguration"].(map[string]interface{})
	origins := corsConfig["AllowOrigins"].([]interface{})
	if len(origins) != 1 || origins[0] != "https://example.com" {
		t.Errorf("expected AllowOrigins ['https://example.com'], got %v", origins)
	}
}

func TestHttpApiTransformer_Transform_WithCorsMap(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		CorsConfiguration: map[string]interface{}{
			"AllowOrigins":     []interface{}{"https://example.com", "https://other.com"},
			"AllowMethods":     []interface{}{"GET", "POST"},
			"AllowHeaders":     []interface{}{"Authorization", "Content-Type"},
			"ExposeHeaders":    []interface{}{"X-Custom-Header"},
			"AllowCredentials": true,
			"MaxAge":           600,
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	corsConfig := apiProps["CorsConfiguration"].(map[string]interface{})

	origins := corsConfig["AllowOrigins"].([]interface{})
	if len(origins) != 2 {
		t.Errorf("expected 2 AllowOrigins, got %d", len(origins))
	}

	methods := corsConfig["AllowMethods"].([]interface{})
	if len(methods) != 2 {
		t.Errorf("expected 2 AllowMethods, got %d", len(methods))
	}

	if corsConfig["AllowCredentials"] != true {
		t.Errorf("expected AllowCredentials true, got %v", corsConfig["AllowCredentials"])
	}

	if corsConfig["MaxAge"] != 600 {
		t.Errorf("expected MaxAge 600, got %v", corsConfig["MaxAge"])
	}
}

func TestHttpApiTransformer_Transform_WithAccessLogSettings(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		AccessLogSettings: &HttpApiAccessLogSettings{
			DestinationArn: map[string]interface{}{
				"Fn::GetAtt": []string{"LogGroup", "Arn"},
			},
			Format: `{"requestId":"$context.requestId"}`,
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stageResource := resources["MyHttpApiStage"].(map[string]interface{})
	stageProps := stageResource["Properties"].(map[string]interface{})

	accessLog := stageProps["AccessLogSettings"].(map[string]interface{})
	if accessLog == nil {
		t.Fatal("AccessLogSettings should be present")
	}

	destArn := accessLog["DestinationArn"].(map[string]interface{})
	getAtt := destArn["Fn::GetAtt"].([]string)
	if getAtt[0] != "LogGroup" || getAtt[1] != "Arn" {
		t.Errorf("expected DestinationArn Fn::GetAtt [LogGroup, Arn], got %v", getAtt)
	}

	if accessLog["Format"] != `{"requestId":"$context.requestId"}` {
		t.Errorf("unexpected Format: %v", accessLog["Format"])
	}
}

func TestHttpApiTransformer_Transform_WithAccessLogSettingsDefaultFormat(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		AccessLogSettings: &HttpApiAccessLogSettings{
			DestinationArn: "arn:aws:logs:us-east-1:123456789:log-group:my-log-group",
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stageResource := resources["MyHttpApiStage"].(map[string]interface{})
	stageProps := stageResource["Properties"].(map[string]interface{})

	accessLog := stageProps["AccessLogSettings"].(map[string]interface{})
	if accessLog["Format"] == nil || accessLog["Format"] == "" {
		t.Error("expected default Format to be set when DestinationArn is provided")
	}
}

func TestHttpApiTransformer_Transform_WithDefaultRouteSettings(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		DefaultRouteSettings: &HttpApiRouteSettings{
			ThrottlingBurstLimit:   100,
			ThrottlingRateLimit:    50.0,
			DetailedMetricsEnabled: true,
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stageResource := resources["MyHttpApiStage"].(map[string]interface{})
	stageProps := stageResource["Properties"].(map[string]interface{})

	defaultRouteSettings := stageProps["DefaultRouteSettings"].(map[string]interface{})
	if defaultRouteSettings["ThrottlingBurstLimit"] != 100 {
		t.Errorf("expected ThrottlingBurstLimit 100, got %v", defaultRouteSettings["ThrottlingBurstLimit"])
	}
	if defaultRouteSettings["ThrottlingRateLimit"] != 50.0 {
		t.Errorf("expected ThrottlingRateLimit 50.0, got %v", defaultRouteSettings["ThrottlingRateLimit"])
	}
	if defaultRouteSettings["DetailedMetricsEnabled"] != true {
		t.Errorf("expected DetailedMetricsEnabled true, got %v", defaultRouteSettings["DetailedMetricsEnabled"])
	}
}

func TestHttpApiTransformer_Transform_WithRouteSettings(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		RouteSettings: map[string]interface{}{
			"GET /users": map[string]interface{}{
				"ThrottlingBurstLimit": 200,
			},
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stageResource := resources["MyHttpApiStage"].(map[string]interface{})
	stageProps := stageResource["Properties"].(map[string]interface{})

	routeSettings := stageProps["RouteSettings"].(map[string]interface{})
	getUsersSettings := routeSettings["GET /users"].(map[string]interface{})
	if getUsersSettings["ThrottlingBurstLimit"] != 200 {
		t.Errorf("expected ThrottlingBurstLimit 200, got %v", getUsersSettings["ThrottlingBurstLimit"])
	}
}

func TestHttpApiTransformer_Transform_WithStageVariables(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		StageVariables: map[string]interface{}{
			"Environment": "production",
			"Version":     "v1",
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stageResource := resources["MyHttpApiStage"].(map[string]interface{})
	stageProps := stageResource["Properties"].(map[string]interface{})

	stageVars := stageProps["StageVariables"].(map[string]interface{})
	if stageVars["Environment"] != "production" {
		t.Errorf("expected Environment 'production', got %v", stageVars["Environment"])
	}
	if stageVars["Version"] != "v1" {
		t.Errorf("expected Version 'v1', got %v", stageVars["Version"])
	}
}

func TestHttpApiTransformer_Transform_WithTags(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Tags: map[string]interface{}{
			"Environment": "test",
			"Project":     "demo",
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	tags := apiProps["Tags"].(map[string]interface{})
	if tags["Environment"] != "test" {
		t.Errorf("expected Environment tag 'test', got %v", tags["Environment"])
	}
	if tags["Project"] != "demo" {
		t.Errorf("expected Project tag 'demo', got %v", tags["Project"])
	}
}

func TestHttpApiTransformer_Transform_WithDefinitionBodyInline(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		DefinitionBody: map[string]interface{}{
			"openapi": "3.0.1",
			"info": map[string]interface{}{
				"title":   "My API",
				"version": "1.0",
			},
			"paths": map[string]interface{}{
				"/users": map[string]interface{}{
					"get": map[string]interface{}{
						"responses": map[string]interface{}{
							"200": map[string]interface{}{
								"description": "OK",
							},
						},
					},
				},
			},
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	if apiProps["Body"] == nil {
		t.Fatal("Body should be set when DefinitionBody is provided")
	}
}

func TestHttpApiTransformer_Transform_WithDefinitionBodyIntrinsics(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		DefinitionBody: map[string]interface{}{
			"openapi": "3.0.1",
			"info": map[string]interface{}{
				"title": map[string]interface{}{
					"Fn::Sub": "${AWS::StackName}-api",
				},
			},
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	// Body should be preserved as a map when it contains intrinsics
	body := apiProps["Body"].(map[string]interface{})
	if body == nil {
		t.Fatal("Body should be preserved as map when containing intrinsics")
	}
}

func TestHttpApiTransformer_Transform_WithDefinitionUriString(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		DefinitionUri: "s3://my-bucket/openapi.yaml",
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	s3Location := apiProps["BodyS3Location"].(map[string]interface{})
	if s3Location["Bucket"] != "my-bucket" {
		t.Errorf("expected Bucket 'my-bucket', got %v", s3Location["Bucket"])
	}
	if s3Location["Key"] != "openapi.yaml" {
		t.Errorf("expected Key 'openapi.yaml', got %v", s3Location["Key"])
	}
}

func TestHttpApiTransformer_Transform_WithDefinitionUriObject(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		DefinitionUri: map[string]interface{}{
			"Bucket":  "my-bucket",
			"Key":     "path/to/openapi.yaml",
			"Version": "abc123",
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	s3Location := apiProps["BodyS3Location"].(map[string]interface{})
	if s3Location["Bucket"] != "my-bucket" {
		t.Errorf("expected Bucket 'my-bucket', got %v", s3Location["Bucket"])
	}
	if s3Location["Key"] != "path/to/openapi.yaml" {
		t.Errorf("expected Key 'path/to/openapi.yaml', got %v", s3Location["Key"])
	}
	if s3Location["Version"] != "abc123" {
		t.Errorf("expected Version 'abc123', got %v", s3Location["Version"])
	}
}

func TestHttpApiTransformer_Transform_WithJwtAuthorizer(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Auth: &HttpApiAuth{
			DefaultAuthorizer: "MyJwtAuthorizer",
			Authorizers: map[string]interface{}{
				"MyJwtAuthorizer": map[string]interface{}{
					"JwtConfiguration": map[string]interface{}{
						"Issuer":   "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_xxx",
						"Audience": []interface{}{"myclient"},
					},
				},
			},
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have Api, Stage, and Authorizer
	if len(resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(resources))
	}

	authResource, ok := resources["MyHttpApiMyJwtAuthorizerAuthorizer"].(map[string]interface{})
	if !ok {
		t.Fatal("Authorizer resource not found")
	}

	if authResource["Type"] != "AWS::ApiGatewayV2::Authorizer" {
		t.Errorf("expected Type 'AWS::ApiGatewayV2::Authorizer', got %v", authResource["Type"])
	}

	authProps := authResource["Properties"].(map[string]interface{})

	if authProps["AuthorizerType"] != "JWT" {
		t.Errorf("expected AuthorizerType 'JWT', got %v", authProps["AuthorizerType"])
	}

	if authProps["Name"] != "MyJwtAuthorizer" {
		t.Errorf("expected Name 'MyJwtAuthorizer', got %v", authProps["Name"])
	}

	// Check identity source
	identitySource := authProps["IdentitySource"].([]interface{})
	if len(identitySource) != 1 || identitySource[0] != "$request.header.Authorization" {
		t.Errorf("expected IdentitySource ['$request.header.Authorization'], got %v", identitySource)
	}
}

func TestHttpApiTransformer_Transform_WithLambdaAuthorizer(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Auth: &HttpApiAuth{
			Authorizers: map[string]interface{}{
				"MyLambdaAuth": map[string]interface{}{
					"FunctionArn": map[string]interface{}{
						"Fn::GetAtt": []interface{}{"AuthFunction", "Arn"},
					},
					"AuthorizerPayloadFormatVersion": "2.0",
					"EnableSimpleResponses":          true,
					"IdentitySource":                 []interface{}{"$request.header.Authorization"},
				},
			},
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	authResource := resources["MyHttpApiMyLambdaAuthAuthorizer"].(map[string]interface{})
	authProps := authResource["Properties"].(map[string]interface{})

	if authProps["AuthorizerType"] != "REQUEST" {
		t.Errorf("expected AuthorizerType 'REQUEST', got %v", authProps["AuthorizerType"])
	}

	if authProps["AuthorizerPayloadFormatVersion"] != "2.0" {
		t.Errorf("expected AuthorizerPayloadFormatVersion '2.0', got %v", authProps["AuthorizerPayloadFormatVersion"])
	}

	if authProps["EnableSimpleResponses"] != true {
		t.Errorf("expected EnableSimpleResponses true, got %v", authProps["EnableSimpleResponses"])
	}

	// AuthorizerUri should be constructed
	authUri := authProps["AuthorizerUri"].(map[string]interface{})
	if authUri == nil {
		t.Fatal("AuthorizerUri should be set")
	}
}

func TestHttpApiTransformer_Transform_WithDomain(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Domain: &HttpApiDomain{
			DomainName:     "api.example.com",
			CertificateArn: "arn:aws:acm:us-east-1:123456789:certificate/abc123",
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have Api, Stage, DomainName, and ApiMapping
	if len(resources) < 4 {
		t.Errorf("expected at least 4 resources, got %d", len(resources))
	}

	// Check Domain Name resource
	domainResource, ok := resources["MyHttpApiDomainName"].(map[string]interface{})
	if !ok {
		t.Fatal("DomainName resource not found")
	}

	if domainResource["Type"] != "AWS::ApiGatewayV2::DomainName" {
		t.Errorf("expected Type 'AWS::ApiGatewayV2::DomainName', got %v", domainResource["Type"])
	}

	domainProps := domainResource["Properties"].(map[string]interface{})
	if domainProps["DomainName"] != "api.example.com" {
		t.Errorf("expected DomainName 'api.example.com', got %v", domainProps["DomainName"])
	}

	// Check domain name configurations
	configs := domainProps["DomainNameConfigurations"].([]interface{})
	if len(configs) != 1 {
		t.Fatalf("expected 1 DomainNameConfiguration, got %d", len(configs))
	}

	config := configs[0].(map[string]interface{})
	if config["CertificateArn"] != "arn:aws:acm:us-east-1:123456789:certificate/abc123" {
		t.Errorf("expected CertificateArn, got %v", config["CertificateArn"])
	}
	if config["EndpointType"] != "REGIONAL" {
		t.Errorf("expected EndpointType 'REGIONAL', got %v", config["EndpointType"])
	}

	// Check API Mapping resource
	mappingResource, ok := resources["MyHttpApiApiMapping"].(map[string]interface{})
	if !ok {
		t.Fatal("ApiMapping resource not found")
	}

	if mappingResource["Type"] != "AWS::ApiGatewayV2::ApiMapping" {
		t.Errorf("expected Type 'AWS::ApiGatewayV2::ApiMapping', got %v", mappingResource["Type"])
	}

	// Check DependsOn
	dependsOn := mappingResource["DependsOn"].([]interface{})
	if len(dependsOn) != 2 {
		t.Errorf("expected 2 DependsOn entries, got %d", len(dependsOn))
	}
}

func TestHttpApiTransformer_Transform_WithDomainAndBasePath(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Domain: &HttpApiDomain{
			DomainName:     "api.example.com",
			CertificateArn: "arn:aws:acm:us-east-1:123456789:certificate/abc123",
			BasePath:       "v1",
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	mappingResource := resources["MyHttpApiApiMapping"].(map[string]interface{})
	mappingProps := mappingResource["Properties"].(map[string]interface{})

	if mappingProps["ApiMappingKey"] != "v1" {
		t.Errorf("expected ApiMappingKey 'v1', got %v", mappingProps["ApiMappingKey"])
	}
}

func TestHttpApiTransformer_Transform_WithDomainAndRoute53(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Domain: &HttpApiDomain{
			DomainName:     "api.example.com",
			CertificateArn: "arn:aws:acm:us-east-1:123456789:certificate/abc123",
			Route53: map[string]interface{}{
				"HostedZoneId": "Z123456789",
			},
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have Route53 RecordSet
	recordSetResource, ok := resources["MyHttpApiRecordSet"].(map[string]interface{})
	if !ok {
		t.Fatal("RecordSet resource not found")
	}

	if recordSetResource["Type"] != "AWS::Route53::RecordSet" {
		t.Errorf("expected Type 'AWS::Route53::RecordSet', got %v", recordSetResource["Type"])
	}

	recordProps := recordSetResource["Properties"].(map[string]interface{})
	if recordProps["Type"] != "A" {
		t.Errorf("expected record Type 'A', got %v", recordProps["Type"])
	}
	if recordProps["HostedZoneId"] != "Z123456789" {
		t.Errorf("expected HostedZoneId 'Z123456789', got %v", recordProps["HostedZoneId"])
	}
}

func TestHttpApiTransformer_Transform_WithDomainAndRoute53IPv6(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		Domain: &HttpApiDomain{
			DomainName:     "api.example.com",
			CertificateArn: "arn:aws:acm:us-east-1:123456789:certificate/abc123",
			Route53: map[string]interface{}{
				"HostedZoneId": "Z123456789",
				"IpV6":         true,
			},
		},
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have both A and AAAA records
	_, okA := resources["MyHttpApiRecordSet"]
	_, okAAAA := resources["MyHttpApiRecordSetV6"]

	if !okA {
		t.Error("A record (MyHttpApiRecordSet) should exist")
	}
	if !okAAAA {
		t.Error("AAAA record (MyHttpApiRecordSetV6) should exist")
	}

	aaaaRecord := resources["MyHttpApiRecordSetV6"].(map[string]interface{})
	aaaaProps := aaaaRecord["Properties"].(map[string]interface{})
	if aaaaProps["Type"] != "AAAA" {
		t.Errorf("expected record Type 'AAAA', got %v", aaaaProps["Type"])
	}
}

func TestHttpApiTransformer_Transform_WithFailOnWarnings(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		FailOnWarnings: true,
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	if apiProps["FailOnWarnings"] != true {
		t.Errorf("expected FailOnWarnings true, got %v", apiProps["FailOnWarnings"])
	}
}

func TestHttpApiTransformer_Transform_WithDisableExecuteApiEndpoint(t *testing.T) {
	transformer := NewHttpApiTransformer()

	api := &HttpApi{
		DisableExecuteApiEndpoint: true,
	}

	ctx := &TransformContext{}
	resources, err := transformer.Transform("MyHttpApi", api, ctx)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyHttpApi"].(map[string]interface{})
	apiProps := apiResource["Properties"].(map[string]interface{})

	if apiProps["DisableExecuteApiEndpoint"] != true {
		t.Errorf("expected DisableExecuteApiEndpoint true, got %v", apiProps["DisableExecuteApiEndpoint"])
	}
}

func TestHttpApiTransformer_parseS3Uri(t *testing.T) {
	transformer := NewHttpApiTransformer()

	// Valid S3 URI
	result, err := transformer.parseS3Uri("s3://my-bucket/path/to/file.yaml")
	if err != nil {
		t.Fatalf("parseS3Uri failed: %v", err)
	}
	if result["Bucket"] != "my-bucket" {
		t.Errorf("expected Bucket 'my-bucket', got %v", result["Bucket"])
	}
	if result["Key"] != "path/to/file.yaml" {
		t.Errorf("expected Key 'path/to/file.yaml', got %v", result["Key"])
	}

	// Invalid URI - not S3
	_, err = transformer.parseS3Uri("https://example.com/file.yaml")
	if err == nil {
		t.Error("expected error for non-S3 URI")
	}

	// Invalid URI - no key
	_, err = transformer.parseS3Uri("s3://bucket-only")
	if err == nil {
		t.Error("expected error for S3 URI without key")
	}
}

func TestHttpApiTransformer_containsIntrinsics(t *testing.T) {
	transformer := NewHttpApiTransformer()

	// Map with no intrinsics
	noIntrinsics := map[string]interface{}{
		"key": "value",
		"nested": map[string]interface{}{
			"inner": "value",
		},
	}
	if transformer.containsIntrinsics(noIntrinsics) {
		t.Error("expected no intrinsics")
	}

	// Map with Ref
	withRef := map[string]interface{}{
		"key": map[string]interface{}{
			"Ref": "SomeResource",
		},
	}
	if !transformer.containsIntrinsics(withRef) {
		t.Error("expected intrinsics with Ref")
	}

	// Map with Fn::Sub
	withFnSub := map[string]interface{}{
		"key": map[string]interface{}{
			"Fn::Sub": "${AWS::StackName}",
		},
	}
	if !transformer.containsIntrinsics(withFnSub) {
		t.Error("expected intrinsics with Fn::Sub")
	}

	// Nested intrinsics in array
	withNestedArray := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{
				"Fn::GetAtt": []string{"Resource", "Arn"},
			},
		},
	}
	if !transformer.containsIntrinsics(withNestedArray) {
		t.Error("expected intrinsics in nested array")
	}
}

func TestHttpApiTransformer_buildLambdaAuthorizerUri(t *testing.T) {
	transformer := NewHttpApiTransformer()

	// Test with string ARN
	result := transformer.buildLambdaAuthorizerUri("arn:aws:lambda:us-east-1:123456789:function:MyAuth")
	fnSub := result.(map[string]interface{})["Fn::Sub"]
	if fnSub == nil {
		t.Error("expected Fn::Sub for string ARN")
	}

	// Test with Ref
	refResult := transformer.buildLambdaAuthorizerUri(map[string]interface{}{
		"Ref": "MyAuthFunction",
	})
	refFnSub := refResult.(map[string]interface{})["Fn::Sub"]
	if refFnSub == nil {
		t.Error("expected Fn::Sub for Ref")
	}

	// Test with Fn::GetAtt
	getAttResult := transformer.buildLambdaAuthorizerUri(map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyAuthFunction", "Arn"},
	})
	getAttFnSub := getAttResult.(map[string]interface{})["Fn::Sub"]
	if getAttFnSub == nil {
		t.Error("expected Fn::Sub for Fn::GetAtt")
	}
}

func TestHttpApiTransformer_buildCorsConfiguration_False(t *testing.T) {
	transformer := NewHttpApiTransformer()

	result := transformer.buildCorsConfiguration(false)
	if result != nil {
		t.Error("expected nil for false CORS configuration")
	}
}

func TestHttpApiTransformer_defaultAccessLogFormat(t *testing.T) {
	transformer := NewHttpApiTransformer()

	format := transformer.defaultAccessLogFormat()
	if format == "" {
		t.Error("expected non-empty default access log format")
	}

	// Should contain common log fields
	if !contains(format, "requestId") {
		t.Error("expected requestId in format")
	}
	if !contains(format, "httpMethod") {
		t.Error("expected httpMethod in format")
	}
	if !contains(format, "status") {
		t.Error("expected status in format")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
