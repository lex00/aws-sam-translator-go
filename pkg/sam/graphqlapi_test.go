package sam

import (
	"testing"
)

func TestGraphQLApiTransformer_BasicApi(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		Name:         "MyGraphQLApi",
		SchemaInline: "type Query { hello: String }",
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	if resources == nil {
		t.Fatal("resources should not be nil")
	}

	// Should create AppSync GraphQL API
	apiResource, ok := resources["MyApi"].(map[string]interface{})
	if !ok {
		t.Fatal("should have AppSync GraphQL API resource")
	}
	if apiResource["Type"] != TypeAppSyncGraphQLAPI {
		t.Errorf("expected Type %q, got %v", TypeAppSyncGraphQLAPI, apiResource["Type"])
	}

	props := apiResource["Properties"].(map[string]interface{})
	if props["Name"] != "MyGraphQLApi" {
		t.Errorf("expected Name 'MyGraphQLApi', got %v", props["Name"])
	}

	// Default auth should be API_KEY
	if props["AuthenticationType"] != "API_KEY" {
		t.Errorf("expected AuthenticationType 'API_KEY', got %v", props["AuthenticationType"])
	}

	// Should create schema
	schemaResource, ok := resources["MyApiSchema"].(map[string]interface{})
	if !ok {
		t.Fatal("should have schema resource")
	}
	if schemaResource["Type"] != TypeAppSyncGraphQLSchema {
		t.Errorf("expected schema Type %q, got %v", TypeAppSyncGraphQLSchema, schemaResource["Type"])
	}
}

func TestGraphQLApiTransformer_WithCognitoAuth(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Auth: &GraphQLAuth{
			Type: "AMAZON_COGNITO_USER_POOLS",
			UserPoolConfig: map[string]interface{}{
				"UserPoolId":    "us-east-1_abc123",
				"DefaultAction": "ALLOW",
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	if props["AuthenticationType"] != "AMAZON_COGNITO_USER_POOLS" {
		t.Errorf("expected AuthenticationType 'AMAZON_COGNITO_USER_POOLS', got %v", props["AuthenticationType"])
	}

	userPoolConfig := props["UserPoolConfig"].(map[string]interface{})
	if userPoolConfig["UserPoolId"] != "us-east-1_abc123" {
		t.Errorf("expected UserPoolId 'us-east-1_abc123', got %v", userPoolConfig["UserPoolId"])
	}
}

func TestGraphQLApiTransformer_WithIAMAuth(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Auth: &GraphQLAuth{
			Type: "AWS_IAM",
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	if props["AuthenticationType"] != "AWS_IAM" {
		t.Errorf("expected AuthenticationType 'AWS_IAM', got %v", props["AuthenticationType"])
	}
}

func TestGraphQLApiTransformer_WithOIDCAuth(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Auth: &GraphQLAuth{
			Type: "OPENID_CONNECT",
			OpenIDConnectConfig: map[string]interface{}{
				"Issuer":   "https://auth.example.com",
				"ClientId": "client-123",
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	if props["AuthenticationType"] != "OPENID_CONNECT" {
		t.Errorf("expected AuthenticationType 'OPENID_CONNECT', got %v", props["AuthenticationType"])
	}

	oidcConfig := props["OpenIDConnectConfig"].(map[string]interface{})
	if oidcConfig["Issuer"] != "https://auth.example.com" {
		t.Errorf("expected Issuer 'https://auth.example.com', got %v", oidcConfig["Issuer"])
	}
}

func TestGraphQLApiTransformer_WithLambdaAuth(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Auth: &GraphQLAuth{
			Type: "AWS_LAMBDA",
			LambdaAuthorizerConfig: map[string]interface{}{
				"AuthorizerUri": "arn:aws:lambda:us-east-1:123456789012:function:MyAuthorizerFunction",
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	if props["AuthenticationType"] != "AWS_LAMBDA" {
		t.Errorf("expected AuthenticationType 'AWS_LAMBDA', got %v", props["AuthenticationType"])
	}

	lambdaConfig := props["LambdaAuthorizerConfig"].(map[string]interface{})
	if lambdaConfig["AuthorizerUri"] == nil {
		t.Error("LambdaAuthorizerConfig should have AuthorizerUri")
	}
}

func TestGraphQLApiTransformer_WithAdditionalAuthProviders(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Auth: &GraphQLAuth{
			Type: "API_KEY",
		},
		AdditionalAuthenticationProviders: []GraphQLAuth{
			{
				Type: "AWS_IAM",
			},
			{
				Type: "AMAZON_COGNITO_USER_POOLS",
				UserPoolConfig: map[string]interface{}{
					"UserPoolId": "us-east-1_xyz789",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	additionalProviders := props["AdditionalAuthenticationProviders"].([]interface{})
	if len(additionalProviders) != 2 {
		t.Errorf("expected 2 additional auth providers, got %d", len(additionalProviders))
	}

	provider1 := additionalProviders[0].(map[string]interface{})
	if provider1["AuthenticationType"] != "AWS_IAM" {
		t.Errorf("expected first provider to be AWS_IAM, got %v", provider1["AuthenticationType"])
	}

	provider2 := additionalProviders[1].(map[string]interface{})
	if provider2["AuthenticationType"] != "AMAZON_COGNITO_USER_POOLS" {
		t.Errorf("expected second provider to be AMAZON_COGNITO_USER_POOLS, got %v", provider2["AuthenticationType"])
	}
}

func TestGraphQLApiTransformer_WithSchemaUri(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaUri: "s3://my-bucket/schemas/schema.graphql",
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	schemaResource := resources["MyApiSchema"].(map[string]interface{})
	props := schemaResource["Properties"].(map[string]interface{})

	if props["DefinitionS3Location"] != "s3://my-bucket/schemas/schema.graphql" {
		t.Errorf("expected DefinitionS3Location, got %v", props["DefinitionS3Location"])
	}
}

func TestGraphQLApiTransformer_WithLambdaDataSource(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		DataSources: map[string]GraphQLApiDataSource{
			"Lambda": {
				Type: "AWS_LAMBDA",
				Name: "LambdaDataSource",
				LambdaConfig: map[string]interface{}{
					"LambdaFunctionArn": "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create data source
	dsResource, ok := resources["MyApiLambdaDataSource"].(map[string]interface{})
	if !ok {
		t.Fatal("should have Lambda data source resource")
	}
	if dsResource["Type"] != TypeAppSyncDataSource {
		t.Errorf("expected Type %q, got %v", TypeAppSyncDataSource, dsResource["Type"])
	}

	props := dsResource["Properties"].(map[string]interface{})
	if props["Type"] != "AWS_LAMBDA" {
		t.Errorf("expected data source Type 'AWS_LAMBDA', got %v", props["Type"])
	}
	if props["Name"] != "LambdaDataSource" {
		t.Errorf("expected Name 'LambdaDataSource', got %v", props["Name"])
	}

	// Should create IAM role for data source
	roleResource, ok := resources["MyApiLambdaDataSourceRole"].(map[string]interface{})
	if !ok {
		t.Fatal("should have IAM role for data source")
	}
	if roleResource["Type"] != "AWS::IAM::Role" {
		t.Errorf("expected Type 'AWS::IAM::Role', got %v", roleResource["Type"])
	}
}

func TestGraphQLApiTransformer_WithDynamoDBDataSource(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { getItem(id: ID!): Item }",
		DataSources: map[string]GraphQLApiDataSource{
			"DynamoDB": {
				Type: "AMAZON_DYNAMODB",
				Name: "DynamoDBDataSource",
				DynamoDBConfig: map[string]interface{}{
					"TableName": "MyTable",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	dsResource := resources["MyApiDynamoDBDataSource"].(map[string]interface{})
	props := dsResource["Properties"].(map[string]interface{})

	if props["Type"] != "AMAZON_DYNAMODB" {
		t.Errorf("expected Type 'AMAZON_DYNAMODB', got %v", props["Type"])
	}

	dynamoDBConfig := props["DynamoDBConfig"].(map[string]interface{})
	if dynamoDBConfig["TableName"] != "MyTable" {
		t.Errorf("expected TableName 'MyTable', got %v", dynamoDBConfig["TableName"])
	}
}

func TestGraphQLApiTransformer_WithNoneDataSource(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		DataSources: map[string]GraphQLApiDataSource{
			"None": {
				Type: "NONE",
				Name: "NoneDataSource",
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	dsResource := resources["MyApiNoneDataSource"].(map[string]interface{})
	props := dsResource["Properties"].(map[string]interface{})

	if props["Type"] != "NONE" {
		t.Errorf("expected Type 'NONE', got %v", props["Type"])
	}

	// NONE data source should not have a role
	if _, ok := resources["MyApiNoneDataSourceRole"]; ok {
		t.Error("NONE data source should not have a role")
	}
}

func TestGraphQLApiTransformer_WithHTTPDataSource(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { fetch: String }",
		DataSources: map[string]GraphQLApiDataSource{
			"HTTP": {
				Type: "HTTP",
				Name: "HTTPDataSource",
				HttpConfig: map[string]interface{}{
					"Endpoint": "https://api.example.com",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	dsResource := resources["MyApiHTTPDataSource"].(map[string]interface{})
	props := dsResource["Properties"].(map[string]interface{})

	if props["Type"] != "HTTP" {
		t.Errorf("expected Type 'HTTP', got %v", props["Type"])
	}

	httpConfig := props["HttpConfig"].(map[string]interface{})
	if httpConfig["Endpoint"] != "https://api.example.com" {
		t.Errorf("expected Endpoint 'https://api.example.com', got %v", httpConfig["Endpoint"])
	}
}

func TestGraphQLApiTransformer_WithResolver(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		DataSources: map[string]GraphQLApiDataSource{
			"None": {
				Type: "NONE",
				Name: "NoneDataSource",
			},
		},
		Resolvers: map[string]GraphQLApiResolver{
			"QueryHello": {
				TypeName:                "Query",
				FieldName:               "hello",
				DataSourceName:          "None",
				RequestMappingTemplate:  "{ \"version\": \"2017-02-28\", \"payload\": {} }",
				ResponseMappingTemplate: "$util.toJson($ctx.result)",
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resolverResource, ok := resources["MyApiQueryHelloResolver"].(map[string]interface{})
	if !ok {
		t.Fatal("should have resolver resource")
	}
	if resolverResource["Type"] != TypeAppSyncResolver {
		t.Errorf("expected Type %q, got %v", TypeAppSyncResolver, resolverResource["Type"])
	}

	props := resolverResource["Properties"].(map[string]interface{})
	if props["TypeName"] != "Query" {
		t.Errorf("expected TypeName 'Query', got %v", props["TypeName"])
	}
	if props["FieldName"] != "hello" {
		t.Errorf("expected FieldName 'hello', got %v", props["FieldName"])
	}
	if props["DataSourceName"] != "None" {
		t.Errorf("expected DataSourceName 'None', got %v", props["DataSourceName"])
	}
}

func TestGraphQLApiTransformer_WithPipelineResolver(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { getData: String }",
		DataSources: map[string]GraphQLApiDataSource{
			"Lambda": {
				Type: "AWS_LAMBDA",
				Name: "LambdaDataSource",
				LambdaConfig: map[string]interface{}{
					"LambdaFunctionArn": "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
				},
			},
		},
		Functions: map[string]GraphQLApiFunction{
			"GetDataFunction": {
				Name:                   "GetDataFunction",
				DataSourceName:         "Lambda",
				RequestMappingTemplate: "{ \"version\": \"2017-02-28\" }",
			},
		},
		Resolvers: map[string]GraphQLApiResolver{
			"QueryGetData": {
				TypeName:  "Query",
				FieldName: "getData",
				Kind:      "PIPELINE",
				PipelineConfig: map[string]interface{}{
					"Functions": []interface{}{"GetDataFunction"},
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have function
	fnResource, ok := resources["MyApiGetDataFunctionFunction"].(map[string]interface{})
	if !ok {
		t.Fatal("should have function resource")
	}
	if fnResource["Type"] != TypeAppSyncFunctionConfig {
		t.Errorf("expected Type %q, got %v", TypeAppSyncFunctionConfig, fnResource["Type"])
	}

	// Should have resolver with pipeline config
	resolverResource := resources["MyApiQueryGetDataResolver"].(map[string]interface{})
	props := resolverResource["Properties"].(map[string]interface{})
	if props["Kind"] != "PIPELINE" {
		t.Errorf("expected Kind 'PIPELINE', got %v", props["Kind"])
	}
	if props["PipelineConfig"] == nil {
		t.Error("PipelineConfig should not be nil")
	}
}

func TestGraphQLApiTransformer_WithXRayTracing(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		XrayEnabled:  true,
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	if props["XrayEnabled"] != true {
		t.Errorf("expected XrayEnabled true, got %v", props["XrayEnabled"])
	}
}

func TestGraphQLApiTransformer_WithLogging(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Logging: map[string]interface{}{
			"FieldLogLevel":         "ALL",
			"ExcludeVerboseContent": false,
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have logging role
	loggingRole, ok := resources["MyApiLoggingRole"].(map[string]interface{})
	if !ok {
		t.Fatal("should have logging role")
	}
	if loggingRole["Type"] != "AWS::IAM::Role" {
		t.Errorf("expected Type 'AWS::IAM::Role', got %v", loggingRole["Type"])
	}

	// API should have LogConfig
	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})
	logConfig := props["LogConfig"].(map[string]interface{})
	if logConfig["FieldLogLevel"] != "ALL" {
		t.Errorf("expected FieldLogLevel 'ALL', got %v", logConfig["FieldLogLevel"])
	}
}

func TestGraphQLApiTransformer_WithCache(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Cache: map[string]interface{}{
			"Type":               "SMALL",
			"ApiCachingBehavior": "FULL_REQUEST_CACHING",
			"Ttl":                600,
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	cacheResource, ok := resources["MyApiCache"].(map[string]interface{})
	if !ok {
		t.Fatal("should have cache resource")
	}
	if cacheResource["Type"] != TypeAppSyncApiCache {
		t.Errorf("expected Type %q, got %v", TypeAppSyncApiCache, cacheResource["Type"])
	}

	props := cacheResource["Properties"].(map[string]interface{})
	if props["Type"] != "SMALL" {
		t.Errorf("expected cache Type 'SMALL', got %v", props["Type"])
	}
	if props["Ttl"] != 600 {
		t.Errorf("expected Ttl 600, got %v", props["Ttl"])
	}
}

func TestGraphQLApiTransformer_WithTags(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Tags: map[string]string{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	tags := props["Tags"].([]interface{})
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}

func TestGraphQLApiTransformer_WithApiKeys(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Auth: &GraphQLAuth{
			Type: "API_KEY",
		},
		ApiKeys: []map[string]interface{}{
			{
				"Description": "Key for frontend app",
				"Expires":     1735689600,
			},
			{
				"Description": "Key for mobile app",
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have 2 API keys
	key0, ok := resources["MyApiApiKey0"].(map[string]interface{})
	if !ok {
		t.Fatal("should have first API key")
	}
	if key0["Type"] != TypeAppSyncApiKey {
		t.Errorf("expected Type %q, got %v", TypeAppSyncApiKey, key0["Type"])
	}

	key1, ok := resources["MyApiApiKey1"].(map[string]interface{})
	if !ok {
		t.Fatal("should have second API key")
	}
	if key1["Type"] != TypeAppSyncApiKey {
		t.Errorf("expected Type %q, got %v", TypeAppSyncApiKey, key1["Type"])
	}
}

func TestGraphQLApiTransformer_DefaultApiKey(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Auth: &GraphQLAuth{
			Type: "API_KEY",
		},
		// No explicit ApiKeys, should create default
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have default API key
	apiKey, ok := resources["MyApiApiKey"].(map[string]interface{})
	if !ok {
		t.Fatal("should have default API key")
	}
	if apiKey["Type"] != TypeAppSyncApiKey {
		t.Errorf("expected Type %q, got %v", TypeAppSyncApiKey, apiKey["Type"])
	}
}

func TestGraphQLApiTransformer_WithDomainName(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		DomainName: map[string]interface{}{
			"DomainName":     "api.example.com",
			"CertificateArn": "arn:aws:acm:us-east-1:123456789012:certificate/abc123",
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have domain name
	domainName, ok := resources["MyApiDomainName"].(map[string]interface{})
	if !ok {
		t.Fatal("should have domain name resource")
	}
	if domainName["Type"] != TypeAppSyncDomainName {
		t.Errorf("expected Type %q, got %v", TypeAppSyncDomainName, domainName["Type"])
	}

	// Should have domain name association
	domainAssoc, ok := resources["MyApiDomainNameAssociation"].(map[string]interface{})
	if !ok {
		t.Fatal("should have domain name association resource")
	}
	if domainAssoc["Type"] != TypeAppSyncDomainNameAssoc {
		t.Errorf("expected Type %q, got %v", TypeAppSyncDomainNameAssoc, domainAssoc["Type"])
	}
}

func TestGraphQLApiTransformer_WithCondition(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Condition:    "IsProduction",
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	if apiResource["Condition"] != "IsProduction" {
		t.Errorf("expected Condition 'IsProduction', got %v", apiResource["Condition"])
	}
}

func TestGraphQLApiTransformer_WithDependsOn(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		DependsOn:    []string{"MyBucket", "MyTable"},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	dependsOn := apiResource["DependsOn"].([]string)
	if len(dependsOn) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(dependsOn))
	}
}

func TestGraphQLApiTransformer_WithMetadata(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		Metadata: map[string]interface{}{
			"cfn-lint": map[string]interface{}{
				"config": map[string]interface{}{
					"ignore_checks": []string{"W3002"},
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyApi"].(map[string]interface{})
	metadata := apiResource["Metadata"].(map[string]interface{})
	if metadata["cfn-lint"] == nil {
		t.Error("Metadata should contain cfn-lint")
	}
}

func TestGraphQLApiTransformer_DefaultName(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		// No Name specified
	}

	resources, err := transformer.Transform("MyGraphQLApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	apiResource := resources["MyGraphQLApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})

	// Should default to logical ID
	if props["Name"] != "MyGraphQLApi" {
		t.Errorf("expected Name 'MyGraphQLApi', got %v", props["Name"])
	}
}

func TestGraphQLApiTransformer_WithJavaScriptResolver(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		DataSources: map[string]GraphQLApiDataSource{
			"None": {
				Type: "NONE",
				Name: "NoneDataSource",
			},
		},
		Resolvers: map[string]GraphQLApiResolver{
			"QueryHello": {
				TypeName:       "Query",
				FieldName:      "hello",
				DataSourceName: "None",
				Runtime: map[string]interface{}{
					"Name":           "APPSYNC_JS",
					"RuntimeVersion": "1.0.0",
				},
				Code: `export function request(ctx) { return {}; }
export function response(ctx) { return "Hello, World!"; }`,
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resolverResource := resources["MyApiQueryHelloResolver"].(map[string]interface{})
	props := resolverResource["Properties"].(map[string]interface{})

	if props["Runtime"] == nil {
		t.Error("Runtime should not be nil")
	}
	if props["Code"] == nil {
		t.Error("Code should not be nil")
	}
}

func TestGraphQLApiTransformer_WithExplicitServiceRole(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		SchemaInline: "type Query { hello: String }",
		DataSources: map[string]GraphQLApiDataSource{
			"Lambda": {
				Type:           "AWS_LAMBDA",
				Name:           "LambdaDataSource",
				ServiceRoleArn: "arn:aws:iam::123456789012:role/MyExistingRole",
				LambdaConfig: map[string]interface{}{
					"LambdaFunctionArn": "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	dsResource := resources["MyApiLambdaDataSource"].(map[string]interface{})
	props := dsResource["Properties"].(map[string]interface{})

	if props["ServiceRoleArn"] != "arn:aws:iam::123456789012:role/MyExistingRole" {
		t.Errorf("expected explicit ServiceRoleArn, got %v", props["ServiceRoleArn"])
	}

	// Should NOT create a role when explicitly provided
	if _, ok := resources["MyApiLambdaDataSourceRole"]; ok {
		t.Error("should not create role when explicitly provided")
	}
}

func TestGraphQLApiTransformer_CompleteExample(t *testing.T) {
	transformer := NewGraphQLApiTransformer()

	api := &GraphQLApi{
		Name: "MyCompleteAPI",
		SchemaInline: `
			type Query {
				getItem(id: ID!): Item
				listItems: [Item]
			}
			type Mutation {
				createItem(input: CreateItemInput!): Item
			}
			type Item {
				id: ID!
				name: String!
			}
			input CreateItemInput {
				name: String!
			}
		`,
		Auth: &GraphQLAuth{
			Type: "AMAZON_COGNITO_USER_POOLS",
			UserPoolConfig: map[string]interface{}{
				"UserPoolId":    "us-east-1_abc123",
				"DefaultAction": "ALLOW",
			},
		},
		AdditionalAuthenticationProviders: []GraphQLAuth{
			{Type: "API_KEY"},
		},
		DataSources: map[string]GraphQLApiDataSource{
			"ItemsTable": {
				Type: "AMAZON_DYNAMODB",
				Name: "ItemsTableDataSource",
				DynamoDBConfig: map[string]interface{}{
					"TableName": "Items",
				},
			},
		},
		Resolvers: map[string]GraphQLApiResolver{
			"QueryGetItem": {
				TypeName:                "Query",
				FieldName:               "getItem",
				DataSourceName:          "ItemsTable",
				RequestMappingTemplate:  `{"version": "2017-02-28", "operation": "GetItem", "key": {"id": $util.dynamodb.toDynamoDBJson($ctx.args.id)}}`,
				ResponseMappingTemplate: "$util.toJson($ctx.result)",
			},
		},
		Logging: map[string]interface{}{
			"FieldLogLevel": "ERROR",
		},
		XrayEnabled: true,
		Tags: map[string]string{
			"Environment": "test",
		},
	}

	resources, err := transformer.Transform("MyCompleteApi", api, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have all expected resources
	expectedResources := []string{
		"MyCompleteApi",                         // GraphQL API
		"MyCompleteApiSchema",                   // Schema
		"MyCompleteApiItemsTableDataSource",     // Data source
		"MyCompleteApiItemsTableDataSourceRole", // Data source role
		"MyCompleteApiQueryGetItemResolver",     // Resolver
		"MyCompleteApiLoggingRole",              // Logging role
		"MyCompleteApiApiKey",                   // Default API key (from additional provider)
	}

	for _, resourceID := range expectedResources {
		if _, ok := resources[resourceID]; !ok {
			t.Errorf("expected resource %q not found", resourceID)
		}
	}

	// Verify API properties
	apiResource := resources["MyCompleteApi"].(map[string]interface{})
	props := apiResource["Properties"].(map[string]interface{})
	if props["Name"] != "MyCompleteAPI" {
		t.Errorf("expected Name 'MyCompleteAPI', got %v", props["Name"])
	}
	if props["XrayEnabled"] != true {
		t.Errorf("expected XrayEnabled true, got %v", props["XrayEnabled"])
	}
}
