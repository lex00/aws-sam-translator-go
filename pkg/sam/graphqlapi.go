// Package sam provides SAM resource transformers.
package sam

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/model/iam"
)

// AWS AppSync resource type constants
const (
	TypeAppSyncGraphQLAPI      = "AWS::AppSync::GraphQLApi"
	TypeAppSyncApiKey          = "AWS::AppSync::ApiKey"
	TypeAppSyncDataSource      = "AWS::AppSync::DataSource"
	TypeAppSyncResolver        = "AWS::AppSync::Resolver"
	TypeAppSyncFunctionConfig  = "AWS::AppSync::FunctionConfiguration"
	TypeAppSyncGraphQLSchema   = "AWS::AppSync::GraphQLSchema"
	TypeAppSyncApiCache        = "AWS::AppSync::ApiCache"
	TypeAppSyncDomainName      = "AWS::AppSync::DomainName"
	TypeAppSyncDomainNameAssoc = "AWS::AppSync::DomainNameApiAssociation"
)

// GraphQLAuth represents an authentication configuration for GraphQL API.
type GraphQLAuth struct {
	// Type is the authentication type: API_KEY, AWS_IAM, OPENID_CONNECT, AMAZON_COGNITO_USER_POOLS, AWS_LAMBDA
	Type string `json:"Type" yaml:"Type"`

	// UserPoolConfig is the Cognito User Pool configuration (for AMAZON_COGNITO_USER_POOLS).
	UserPoolConfig map[string]interface{} `json:"UserPoolConfig,omitempty" yaml:"UserPoolConfig,omitempty"`

	// OpenIDConnectConfig is the OIDC configuration (for OPENID_CONNECT).
	OpenIDConnectConfig map[string]interface{} `json:"OpenIDConnectConfig,omitempty" yaml:"OpenIDConnectConfig,omitempty"`

	// LambdaAuthorizerConfig is the Lambda authorizer configuration (for AWS_LAMBDA).
	LambdaAuthorizerConfig map[string]interface{} `json:"LambdaAuthorizerConfig,omitempty" yaml:"LambdaAuthorizerConfig,omitempty"`
}

// GraphQLApiDataSource represents a data source configuration.
type GraphQLApiDataSource struct {
	// Type is the data source type: NONE, AWS_LAMBDA, AMAZON_DYNAMODB, AMAZON_ELASTICSEARCH,
	// AMAZON_OPENSEARCH_SERVICE, HTTP, RELATIONAL_DATABASE, AMAZON_EVENTBRIDGE
	Type string `json:"Type" yaml:"Type"`

	// Name is the data source name.
	Name string `json:"Name" yaml:"Name"`

	// Description is the data source description.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// ServiceRoleArn is the IAM role ARN for the data source.
	ServiceRoleArn interface{} `json:"ServiceRoleArn,omitempty" yaml:"ServiceRoleArn,omitempty"`

	// LambdaConfig is the Lambda configuration (for AWS_LAMBDA).
	LambdaConfig map[string]interface{} `json:"LambdaConfig,omitempty" yaml:"LambdaConfig,omitempty"`

	// DynamoDBConfig is the DynamoDB configuration (for AMAZON_DYNAMODB).
	DynamoDBConfig map[string]interface{} `json:"DynamoDBConfig,omitempty" yaml:"DynamoDBConfig,omitempty"`

	// HttpConfig is the HTTP configuration (for HTTP).
	HttpConfig map[string]interface{} `json:"HttpConfig,omitempty" yaml:"HttpConfig,omitempty"`

	// ElasticsearchConfig is the Elasticsearch configuration (for AMAZON_ELASTICSEARCH).
	ElasticsearchConfig map[string]interface{} `json:"ElasticsearchConfig,omitempty" yaml:"ElasticsearchConfig,omitempty"`

	// OpenSearchServiceConfig is the OpenSearch configuration (for AMAZON_OPENSEARCH_SERVICE).
	OpenSearchServiceConfig map[string]interface{} `json:"OpenSearchServiceConfig,omitempty" yaml:"OpenSearchServiceConfig,omitempty"`

	// RelationalDatabaseConfig is the RDS configuration (for RELATIONAL_DATABASE).
	RelationalDatabaseConfig map[string]interface{} `json:"RelationalDatabaseConfig,omitempty" yaml:"RelationalDatabaseConfig,omitempty"`

	// EventBridgeConfig is the EventBridge configuration (for AMAZON_EVENTBRIDGE).
	EventBridgeConfig map[string]interface{} `json:"EventBridgeConfig,omitempty" yaml:"EventBridgeConfig,omitempty"`
}

// GraphQLApiResolver represents a resolver configuration.
type GraphQLApiResolver struct {
	// FieldName is the field name in the GraphQL schema.
	FieldName string `json:"FieldName" yaml:"FieldName"`

	// TypeName is the type name in the GraphQL schema.
	TypeName string `json:"TypeName" yaml:"TypeName"`

	// DataSourceName is the name of the data source.
	DataSourceName string `json:"DataSourceName,omitempty" yaml:"DataSourceName,omitempty"`

	// Kind is the resolver kind: UNIT or PIPELINE.
	Kind string `json:"Kind,omitempty" yaml:"Kind,omitempty"`

	// RequestMappingTemplate is the VTL request mapping template.
	RequestMappingTemplate string `json:"RequestMappingTemplate,omitempty" yaml:"RequestMappingTemplate,omitempty"`

	// ResponseMappingTemplate is the VTL response mapping template.
	ResponseMappingTemplate string `json:"ResponseMappingTemplate,omitempty" yaml:"ResponseMappingTemplate,omitempty"`

	// RequestMappingTemplateS3Location is the S3 location of the request mapping template.
	RequestMappingTemplateS3Location string `json:"RequestMappingTemplateS3Location,omitempty" yaml:"RequestMappingTemplateS3Location,omitempty"`

	// ResponseMappingTemplateS3Location is the S3 location of the response mapping template.
	ResponseMappingTemplateS3Location string `json:"ResponseMappingTemplateS3Location,omitempty" yaml:"ResponseMappingTemplateS3Location,omitempty"`

	// PipelineConfig is the pipeline configuration (for PIPELINE kind).
	PipelineConfig map[string]interface{} `json:"PipelineConfig,omitempty" yaml:"PipelineConfig,omitempty"`

	// CachingConfig is the caching configuration.
	CachingConfig map[string]interface{} `json:"CachingConfig,omitempty" yaml:"CachingConfig,omitempty"`

	// SyncConfig is the sync configuration for Lambda resolvers.
	SyncConfig map[string]interface{} `json:"SyncConfig,omitempty" yaml:"SyncConfig,omitempty"`

	// MaxBatchSize is the maximum batch size for batch invocation.
	MaxBatchSize int `json:"MaxBatchSize,omitempty" yaml:"MaxBatchSize,omitempty"`

	// Runtime is the runtime configuration for JavaScript resolvers.
	Runtime map[string]interface{} `json:"Runtime,omitempty" yaml:"Runtime,omitempty"`

	// Code is the inline resolver code.
	Code string `json:"Code,omitempty" yaml:"Code,omitempty"`

	// CodeS3Location is the S3 location of the resolver code.
	CodeS3Location string `json:"CodeS3Location,omitempty" yaml:"CodeS3Location,omitempty"`
}

// GraphQLApiFunction represents an AppSync function configuration.
type GraphQLApiFunction struct {
	// Name is the function name.
	Name string `json:"Name" yaml:"Name"`

	// Description is the function description.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// DataSourceName is the name of the data source.
	DataSourceName string `json:"DataSourceName" yaml:"DataSourceName"`

	// RequestMappingTemplate is the VTL request mapping template.
	RequestMappingTemplate string `json:"RequestMappingTemplate,omitempty" yaml:"RequestMappingTemplate,omitempty"`

	// ResponseMappingTemplate is the VTL response mapping template.
	ResponseMappingTemplate string `json:"ResponseMappingTemplate,omitempty" yaml:"ResponseMappingTemplate,omitempty"`

	// RequestMappingTemplateS3Location is the S3 location of the request mapping template.
	RequestMappingTemplateS3Location string `json:"RequestMappingTemplateS3Location,omitempty" yaml:"RequestMappingTemplateS3Location,omitempty"`

	// ResponseMappingTemplateS3Location is the S3 location of the response mapping template.
	ResponseMappingTemplateS3Location string `json:"ResponseMappingTemplateS3Location,omitempty" yaml:"ResponseMappingTemplateS3Location,omitempty"`

	// Runtime is the runtime configuration for JavaScript functions.
	Runtime map[string]interface{} `json:"Runtime,omitempty" yaml:"Runtime,omitempty"`

	// Code is the inline function code.
	Code string `json:"Code,omitempty" yaml:"Code,omitempty"`

	// CodeS3Location is the S3 location of the function code.
	CodeS3Location string `json:"CodeS3Location,omitempty" yaml:"CodeS3Location,omitempty"`

	// SyncConfig is the sync configuration.
	SyncConfig map[string]interface{} `json:"SyncConfig,omitempty" yaml:"SyncConfig,omitempty"`

	// MaxBatchSize is the maximum batch size for batch invocation.
	MaxBatchSize int `json:"MaxBatchSize,omitempty" yaml:"MaxBatchSize,omitempty"`
}

// GraphQLApi represents an AWS::Serverless::GraphQLApi resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-graphqlapi.html
type GraphQLApi struct {
	// Name is the name of the GraphQL API.
	Name string `json:"Name,omitempty" yaml:"Name,omitempty"`

	// SchemaInline is the inline GraphQL schema definition.
	SchemaInline string `json:"SchemaInline,omitempty" yaml:"SchemaInline,omitempty"`

	// SchemaUri is the S3 URI of the GraphQL schema.
	SchemaUri interface{} `json:"SchemaUri,omitempty" yaml:"SchemaUri,omitempty"`

	// Auth configures authentication for the API.
	Auth *GraphQLAuth `json:"Auth,omitempty" yaml:"Auth,omitempty"`

	// AdditionalAuthenticationProviders is a list of additional authentication providers.
	AdditionalAuthenticationProviders []GraphQLAuth `json:"AdditionalAuthenticationProviders,omitempty" yaml:"AdditionalAuthenticationProviders,omitempty"`

	// DataSources is a map of data source configurations.
	DataSources map[string]GraphQLApiDataSource `json:"DataSources,omitempty" yaml:"DataSources,omitempty"`

	// Resolvers is a map of resolver configurations.
	Resolvers map[string]GraphQLApiResolver `json:"Resolvers,omitempty" yaml:"Resolvers,omitempty"`

	// Functions is a map of function configurations.
	Functions map[string]GraphQLApiFunction `json:"Functions,omitempty" yaml:"Functions,omitempty"`

	// Logging configures CloudWatch logging.
	Logging map[string]interface{} `json:"Logging,omitempty" yaml:"Logging,omitempty"`

	// XrayEnabled enables AWS X-Ray tracing.
	XrayEnabled bool `json:"XrayEnabled,omitempty" yaml:"XrayEnabled,omitempty"`

	// Cache configures API caching.
	Cache map[string]interface{} `json:"Cache,omitempty" yaml:"Cache,omitempty"`

	// Tags is a map of key-value pairs to apply to the API.
	Tags map[string]string `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// DomainName configures a custom domain.
	DomainName map[string]interface{} `json:"DomainName,omitempty" yaml:"DomainName,omitempty"`

	// ApiKeys is a list of API key configurations.
	ApiKeys []map[string]interface{} `json:"ApiKeys,omitempty" yaml:"ApiKeys,omitempty"`

	// Condition is a CloudFormation condition name.
	Condition string `json:"Condition,omitempty" yaml:"Condition,omitempty"`

	// DependsOn specifies resource dependencies.
	DependsOn interface{} `json:"DependsOn,omitempty" yaml:"DependsOn,omitempty"`

	// Metadata is custom metadata for the resource.
	Metadata map[string]interface{} `json:"Metadata,omitempty" yaml:"Metadata,omitempty"`
}

// GraphQLApiTransformer transforms AWS::Serverless::GraphQLApi to CloudFormation.
type GraphQLApiTransformer struct{}

// NewGraphQLApiTransformer creates a new GraphQLApiTransformer.
func NewGraphQLApiTransformer() *GraphQLApiTransformer {
	return &GraphQLApiTransformer{}
}

// Transform converts a SAM GraphQLApi to CloudFormation resources.
func (t *GraphQLApiTransformer) Transform(logicalID string, api *GraphQLApi, ctx *TransformContext) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Build the main GraphQL API
	apiProps, err := t.buildApiProperties(logicalID, api)
	if err != nil {
		return nil, fmt.Errorf("failed to build API properties: %w", err)
	}

	apiResource := map[string]interface{}{
		"Type":       TypeAppSyncGraphQLAPI,
		"Properties": apiProps,
	}

	if api.Condition != "" {
		apiResource["Condition"] = api.Condition
	}
	if api.DependsOn != nil {
		apiResource["DependsOn"] = api.DependsOn
	}
	if api.Metadata != nil {
		apiResource["Metadata"] = api.Metadata
	}

	resources[logicalID] = apiResource

	// Build the schema
	schemaResource, err := t.buildSchema(logicalID, api)
	if err != nil {
		return nil, fmt.Errorf("failed to build schema: %w", err)
	}
	if schemaResource != nil {
		resources[logicalID+"Schema"] = schemaResource
	}

	// Build data sources
	dataSourceResources, roleResources, err := t.buildDataSources(logicalID, api)
	if err != nil {
		return nil, fmt.Errorf("failed to build data sources: %w", err)
	}
	for k, v := range dataSourceResources {
		resources[k] = v
	}
	for k, v := range roleResources {
		resources[k] = v
	}

	// Build functions
	functionResources, err := t.buildFunctions(logicalID, api)
	if err != nil {
		return nil, fmt.Errorf("failed to build functions: %w", err)
	}
	for k, v := range functionResources {
		resources[k] = v
	}

	// Build resolvers
	resolverResources, err := t.buildResolvers(logicalID, api)
	if err != nil {
		return nil, fmt.Errorf("failed to build resolvers: %w", err)
	}
	for k, v := range resolverResources {
		resources[k] = v
	}

	// Build API keys
	apiKeyResources, err := t.buildApiKeys(logicalID, api)
	if err != nil {
		return nil, fmt.Errorf("failed to build API keys: %w", err)
	}
	for k, v := range apiKeyResources {
		resources[k] = v
	}

	// Build cache
	if api.Cache != nil {
		cacheResource, err := t.buildCache(logicalID, api)
		if err != nil {
			return nil, fmt.Errorf("failed to build cache: %w", err)
		}
		resources[logicalID+"Cache"] = cacheResource
	}

	// Build domain name
	if api.DomainName != nil {
		domainResources, err := t.buildDomainName(logicalID, api)
		if err != nil {
			return nil, fmt.Errorf("failed to build domain name: %w", err)
		}
		for k, v := range domainResources {
			resources[k] = v
		}
	}

	// Build logging role
	if api.Logging != nil {
		loggingResources, err := t.buildLogging(logicalID, api)
		if err != nil {
			return nil, fmt.Errorf("failed to build logging: %w", err)
		}
		for k, v := range loggingResources {
			resources[k] = v
		}
	}

	return resources, nil
}

// buildApiProperties builds the AppSync GraphQL API properties.
func (t *GraphQLApiTransformer) buildApiProperties(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	props := make(map[string]interface{})

	// Name (default to logical ID)
	if api.Name != "" {
		props["Name"] = api.Name
	} else {
		props["Name"] = logicalID
	}

	// Authentication
	if api.Auth != nil {
		props["AuthenticationType"] = api.Auth.Type

		if api.Auth.UserPoolConfig != nil {
			props["UserPoolConfig"] = api.Auth.UserPoolConfig
		}
		if api.Auth.OpenIDConnectConfig != nil {
			props["OpenIDConnectConfig"] = api.Auth.OpenIDConnectConfig
		}
		if api.Auth.LambdaAuthorizerConfig != nil {
			props["LambdaAuthorizerConfig"] = api.Auth.LambdaAuthorizerConfig
		}
	} else {
		// Default to API_KEY if no auth specified
		props["AuthenticationType"] = "API_KEY"
	}

	// Additional authentication providers
	if len(api.AdditionalAuthenticationProviders) > 0 {
		additionalProviders := make([]interface{}, len(api.AdditionalAuthenticationProviders))
		for i, provider := range api.AdditionalAuthenticationProviders {
			providerMap := map[string]interface{}{
				"AuthenticationType": provider.Type,
			}
			if provider.UserPoolConfig != nil {
				providerMap["UserPoolConfig"] = provider.UserPoolConfig
			}
			if provider.OpenIDConnectConfig != nil {
				providerMap["OpenIDConnectConfig"] = provider.OpenIDConnectConfig
			}
			if provider.LambdaAuthorizerConfig != nil {
				providerMap["LambdaAuthorizerConfig"] = provider.LambdaAuthorizerConfig
			}
			additionalProviders[i] = providerMap
		}
		props["AdditionalAuthenticationProviders"] = additionalProviders
	}

	// X-Ray tracing
	if api.XrayEnabled {
		props["XrayEnabled"] = true
	}

	// Logging configuration (reference to CloudWatch role)
	if api.Logging != nil {
		logConfig := map[string]interface{}{
			"CloudWatchLogsRoleArn": map[string]interface{}{
				"Fn::GetAtt": []string{logicalID + "LoggingRole", "Arn"},
			},
		}
		if fieldLogLevel, ok := api.Logging["FieldLogLevel"]; ok {
			logConfig["FieldLogLevel"] = fieldLogLevel
		}
		if excludeVerbose, ok := api.Logging["ExcludeVerboseContent"]; ok {
			logConfig["ExcludeVerboseContent"] = excludeVerbose
		}
		props["LogConfig"] = logConfig
	}

	// Tags
	if len(api.Tags) > 0 {
		tags := make([]interface{}, 0, len(api.Tags))
		for k, v := range api.Tags {
			tags = append(tags, map[string]interface{}{
				"Key":   k,
				"Value": v,
			})
		}
		props["Tags"] = tags
	}

	return props, nil
}

// buildSchema builds the GraphQL schema resource.
func (t *GraphQLApiTransformer) buildSchema(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	// Need either SchemaInline or SchemaUri
	if api.SchemaInline == "" && api.SchemaUri == nil {
		return nil, nil // No schema specified, might be defined elsewhere
	}

	props := map[string]interface{}{
		"ApiId": map[string]interface{}{
			"Fn::GetAtt": []string{logicalID, "ApiId"},
		},
	}

	if api.SchemaInline != "" {
		props["Definition"] = api.SchemaInline
	} else if api.SchemaUri != nil {
		switch uri := api.SchemaUri.(type) {
		case string:
			props["DefinitionS3Location"] = uri
		case map[string]interface{}:
			// Could be an intrinsic function or S3 location object
			if bucket, hasBucket := uri["Bucket"]; hasBucket {
				s3Location := fmt.Sprintf("s3://%v/%v", bucket, uri["Key"])
				props["DefinitionS3Location"] = s3Location
			} else {
				props["DefinitionS3Location"] = uri
			}
		}
	}

	return map[string]interface{}{
		"Type":       TypeAppSyncGraphQLSchema,
		"Properties": props,
	}, nil
}

// buildDataSources builds data source resources.
func (t *GraphQLApiTransformer) buildDataSources(logicalID string, api *GraphQLApi) (map[string]interface{}, map[string]interface{}, error) {
	dataSourceResources := make(map[string]interface{})
	roleResources := make(map[string]interface{})

	for dsName, ds := range api.DataSources {
		dsLogicalID := logicalID + dsName + "DataSource"
		roleLogicalID := logicalID + dsName + "DataSourceRole"

		props := map[string]interface{}{
			"ApiId": map[string]interface{}{
				"Fn::GetAtt": []string{logicalID, "ApiId"},
			},
			"Name": ds.Name,
			"Type": ds.Type,
		}

		if ds.Description != "" {
			props["Description"] = ds.Description
		}

		// Handle service role
		if ds.ServiceRoleArn != nil {
			props["ServiceRoleArn"] = ds.ServiceRoleArn
		} else if ds.Type != "NONE" {
			// Create a role for the data source
			roleResource := t.buildDataSourceRole(ds, roleLogicalID)
			roleResources[roleLogicalID] = roleResource
			props["ServiceRoleArn"] = map[string]interface{}{
				"Fn::GetAtt": []string{roleLogicalID, "Arn"},
			}
		}

		// Type-specific configuration
		switch ds.Type {
		case "AWS_LAMBDA":
			if ds.LambdaConfig != nil {
				props["LambdaConfig"] = ds.LambdaConfig
			}
		case "AMAZON_DYNAMODB":
			if ds.DynamoDBConfig != nil {
				props["DynamoDBConfig"] = ds.DynamoDBConfig
			}
		case "HTTP":
			if ds.HttpConfig != nil {
				props["HttpConfig"] = ds.HttpConfig
			}
		case "AMAZON_ELASTICSEARCH":
			if ds.ElasticsearchConfig != nil {
				props["ElasticsearchConfig"] = ds.ElasticsearchConfig
			}
		case "AMAZON_OPENSEARCH_SERVICE":
			if ds.OpenSearchServiceConfig != nil {
				props["OpenSearchServiceConfig"] = ds.OpenSearchServiceConfig
			}
		case "RELATIONAL_DATABASE":
			if ds.RelationalDatabaseConfig != nil {
				props["RelationalDatabaseConfig"] = ds.RelationalDatabaseConfig
			}
		case "AMAZON_EVENTBRIDGE":
			if ds.EventBridgeConfig != nil {
				props["EventBridgeConfig"] = ds.EventBridgeConfig
			}
		}

		dataSourceResources[dsLogicalID] = map[string]interface{}{
			"Type":       TypeAppSyncDataSource,
			"Properties": props,
			"DependsOn":  logicalID + "Schema",
		}
	}

	return dataSourceResources, roleResources, nil
}

// buildDataSourceRole creates an IAM role for a data source.
func (t *GraphQLApiTransformer) buildDataSourceRole(ds GraphQLApiDataSource, roleLogicalID string) map[string]interface{} {
	trustPolicy := iam.NewAssumeRolePolicyForService("appsync.amazonaws.com")
	role := iam.NewRole(trustPolicy)

	// Add permissions based on data source type
	var actions []interface{}
	var resources []interface{}

	switch ds.Type {
	case "AWS_LAMBDA":
		actions = []interface{}{"lambda:InvokeFunction"}
		if ds.LambdaConfig != nil {
			if lambdaArn, ok := ds.LambdaConfig["LambdaFunctionArn"]; ok {
				resources = []interface{}{lambdaArn}
			}
		}
		if resources == nil {
			resources = []interface{}{"*"}
		}
	case "AMAZON_DYNAMODB":
		actions = []interface{}{
			"dynamodb:GetItem",
			"dynamodb:PutItem",
			"dynamodb:DeleteItem",
			"dynamodb:UpdateItem",
			"dynamodb:Query",
			"dynamodb:Scan",
			"dynamodb:BatchGetItem",
			"dynamodb:BatchWriteItem",
		}
		if ds.DynamoDBConfig != nil {
			if tableName, ok := ds.DynamoDBConfig["TableName"]; ok {
				resources = []interface{}{
					map[string]interface{}{
						"Fn::Sub": fmt.Sprintf("arn:aws:dynamodb:${AWS::Region}:${AWS::AccountId}:table/%v", tableName),
					},
					map[string]interface{}{
						"Fn::Sub": fmt.Sprintf("arn:aws:dynamodb:${AWS::Region}:${AWS::AccountId}:table/%v/*", tableName),
					},
				}
			}
		}
		if resources == nil {
			resources = []interface{}{"*"}
		}
	default:
		actions = []interface{}{"*"}
		resources = []interface{}{"*"}
	}

	stmt := iam.NewStatement(iam.EffectAllow)
	stmt.Action = actions
	stmt.Resource = resources

	doc := iam.NewPolicyDocument()
	doc.AddStatement(stmt)

	role.Policies = []iam.InlinePolicy{
		{
			PolicyName:     roleLogicalID + "Policy",
			PolicyDocument: doc,
		},
	}

	return map[string]interface{}{
		"Type":       "AWS::IAM::Role",
		"Properties": role.ToCloudFormation(),
	}
}

// buildFunctions builds AppSync function resources.
func (t *GraphQLApiTransformer) buildFunctions(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	for fnName, fn := range api.Functions {
		fnLogicalID := logicalID + fnName + "Function"

		props := map[string]interface{}{
			"ApiId": map[string]interface{}{
				"Fn::GetAtt": []string{logicalID, "ApiId"},
			},
			"Name":           fn.Name,
			"DataSourceName": fn.DataSourceName,
		}

		if fn.Description != "" {
			props["Description"] = fn.Description
		}

		// VTL templates
		if fn.RequestMappingTemplate != "" {
			props["RequestMappingTemplate"] = fn.RequestMappingTemplate
		} else if fn.RequestMappingTemplateS3Location != "" {
			props["RequestMappingTemplateS3Location"] = fn.RequestMappingTemplateS3Location
		}

		if fn.ResponseMappingTemplate != "" {
			props["ResponseMappingTemplate"] = fn.ResponseMappingTemplate
		} else if fn.ResponseMappingTemplateS3Location != "" {
			props["ResponseMappingTemplateS3Location"] = fn.ResponseMappingTemplateS3Location
		}

		// JavaScript runtime
		if fn.Runtime != nil {
			props["Runtime"] = fn.Runtime
		}

		if fn.Code != "" {
			props["Code"] = fn.Code
		} else if fn.CodeS3Location != "" {
			props["CodeS3Location"] = fn.CodeS3Location
		}

		if fn.SyncConfig != nil {
			props["SyncConfig"] = fn.SyncConfig
		}

		if fn.MaxBatchSize > 0 {
			props["MaxBatchSize"] = fn.MaxBatchSize
		}

		resources[fnLogicalID] = map[string]interface{}{
			"Type":       TypeAppSyncFunctionConfig,
			"Properties": props,
			"DependsOn":  logicalID + fn.DataSourceName + "DataSource",
		}
	}

	return resources, nil
}

// buildResolvers builds resolver resources.
func (t *GraphQLApiTransformer) buildResolvers(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	for resolverName, resolver := range api.Resolvers {
		resolverLogicalID := logicalID + resolverName + "Resolver"

		props := map[string]interface{}{
			"ApiId": map[string]interface{}{
				"Fn::GetAtt": []string{logicalID, "ApiId"},
			},
			"TypeName":  resolver.TypeName,
			"FieldName": resolver.FieldName,
		}

		if resolver.Kind != "" {
			props["Kind"] = resolver.Kind
		}

		if resolver.DataSourceName != "" {
			props["DataSourceName"] = resolver.DataSourceName
		}

		// VTL templates
		if resolver.RequestMappingTemplate != "" {
			props["RequestMappingTemplate"] = resolver.RequestMappingTemplate
		} else if resolver.RequestMappingTemplateS3Location != "" {
			props["RequestMappingTemplateS3Location"] = resolver.RequestMappingTemplateS3Location
		}

		if resolver.ResponseMappingTemplate != "" {
			props["ResponseMappingTemplate"] = resolver.ResponseMappingTemplate
		} else if resolver.ResponseMappingTemplateS3Location != "" {
			props["ResponseMappingTemplateS3Location"] = resolver.ResponseMappingTemplateS3Location
		}

		// JavaScript runtime
		if resolver.Runtime != nil {
			props["Runtime"] = resolver.Runtime
		}

		if resolver.Code != "" {
			props["Code"] = resolver.Code
		} else if resolver.CodeS3Location != "" {
			props["CodeS3Location"] = resolver.CodeS3Location
		}

		// Pipeline configuration
		if resolver.PipelineConfig != nil {
			props["PipelineConfig"] = resolver.PipelineConfig
		}

		// Caching
		if resolver.CachingConfig != nil {
			props["CachingConfig"] = resolver.CachingConfig
		}

		// Sync config
		if resolver.SyncConfig != nil {
			props["SyncConfig"] = resolver.SyncConfig
		}

		if resolver.MaxBatchSize > 0 {
			props["MaxBatchSize"] = resolver.MaxBatchSize
		}

		// Build dependencies
		depends := []string{logicalID + "Schema"}
		if resolver.DataSourceName != "" {
			depends = append(depends, logicalID+resolver.DataSourceName+"DataSource")
		}

		resources[resolverLogicalID] = map[string]interface{}{
			"Type":       TypeAppSyncResolver,
			"Properties": props,
			"DependsOn":  depends,
		}
	}

	return resources, nil
}

// buildApiKeys builds API key resources.
func (t *GraphQLApiTransformer) buildApiKeys(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Check if API_KEY auth is used (primary or additional)
	hasApiKeyAuth := api.Auth != nil && api.Auth.Type == "API_KEY"
	if !hasApiKeyAuth {
		for _, provider := range api.AdditionalAuthenticationProviders {
			if provider.Type == "API_KEY" {
				hasApiKeyAuth = true
				break
			}
		}
	}

	// If API_KEY auth is enabled and no explicit keys defined, create a default one
	if hasApiKeyAuth && len(api.ApiKeys) == 0 {
		resources[logicalID+"ApiKey"] = map[string]interface{}{
			"Type": TypeAppSyncApiKey,
			"Properties": map[string]interface{}{
				"ApiId": map[string]interface{}{
					"Fn::GetAtt": []string{logicalID, "ApiId"},
				},
			},
		}
		return resources, nil
	}

	for i, keyConfig := range api.ApiKeys {
		keyLogicalID := fmt.Sprintf("%sApiKey%d", logicalID, i)

		props := map[string]interface{}{
			"ApiId": map[string]interface{}{
				"Fn::GetAtt": []string{logicalID, "ApiId"},
			},
		}

		if desc, ok := keyConfig["Description"]; ok {
			props["Description"] = desc
		}
		if expires, ok := keyConfig["Expires"]; ok {
			props["Expires"] = expires
		}
		if apiKeyId, ok := keyConfig["ApiKeyId"]; ok {
			props["ApiKeyId"] = apiKeyId
		}

		resources[keyLogicalID] = map[string]interface{}{
			"Type":       TypeAppSyncApiKey,
			"Properties": props,
		}
	}

	return resources, nil
}

// buildCache builds the API cache resource.
func (t *GraphQLApiTransformer) buildCache(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	props := map[string]interface{}{
		"ApiId": map[string]interface{}{
			"Fn::GetAtt": []string{logicalID, "ApiId"},
		},
	}

	// Copy cache properties
	for k, v := range api.Cache {
		props[k] = v
	}

	// Ensure required properties
	if _, ok := props["ApiCachingBehavior"]; !ok {
		props["ApiCachingBehavior"] = "FULL_REQUEST_CACHING"
	}
	if _, ok := props["Type"]; !ok {
		props["Type"] = "SMALL"
	}
	if _, ok := props["Ttl"]; !ok {
		props["Ttl"] = 3600
	}

	return map[string]interface{}{
		"Type":       TypeAppSyncApiCache,
		"Properties": props,
	}, nil
}

// buildDomainName builds the custom domain name resources.
func (t *GraphQLApiTransformer) buildDomainName(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	domainProps := make(map[string]interface{})
	for k, v := range api.DomainName {
		domainProps[k] = v
	}

	resources[logicalID+"DomainName"] = map[string]interface{}{
		"Type":       TypeAppSyncDomainName,
		"Properties": domainProps,
	}

	// Create domain name association
	assocProps := map[string]interface{}{
		"ApiId": map[string]interface{}{
			"Fn::GetAtt": []string{logicalID, "ApiId"},
		},
		"DomainName": domainProps["DomainName"],
	}

	resources[logicalID+"DomainNameAssociation"] = map[string]interface{}{
		"Type":       TypeAppSyncDomainNameAssoc,
		"Properties": assocProps,
		"DependsOn":  logicalID + "DomainName",
	}

	return resources, nil
}

// buildLogging builds the CloudWatch logging role.
func (t *GraphQLApiTransformer) buildLogging(logicalID string, api *GraphQLApi) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	trustPolicy := iam.NewAssumeRolePolicyForService("appsync.amazonaws.com")
	role := iam.NewRole(trustPolicy)

	// Add CloudWatch Logs permissions
	stmt := iam.NewStatement(iam.EffectAllow)
	stmt.Action = []interface{}{
		"logs:CreateLogGroup",
		"logs:CreateLogStream",
		"logs:PutLogEvents",
	}
	stmt.Resource = "*"

	doc := iam.NewPolicyDocument()
	doc.AddStatement(stmt)

	role.Policies = []iam.InlinePolicy{
		{
			PolicyName:     logicalID + "LoggingPolicy",
			PolicyDocument: doc,
		},
	}

	resources[logicalID+"LoggingRole"] = map[string]interface{}{
		"Type":       "AWS::IAM::Role",
		"Properties": role.ToCloudFormation(),
	}

	return resources, nil
}
