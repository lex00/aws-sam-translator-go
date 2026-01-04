// Package sam provides SAM resource transformers.
package sam

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Api represents an AWS::Serverless::Api resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-api.html
type Api struct {
	// Name is the name of the API Gateway RestApi resource.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// StageName is the name of the stage (required).
	StageName interface{} `json:"StageName,omitempty" yaml:"StageName,omitempty"`

	// Description is a description of the RestApi resource.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// DefinitionBody is an OpenAPI specification that defines the API.
	DefinitionBody map[string]interface{} `json:"DefinitionBody,omitempty" yaml:"DefinitionBody,omitempty"`

	// DefinitionUri is an S3 URI or local path to an OpenAPI specification file.
	DefinitionUri interface{} `json:"DefinitionUri,omitempty" yaml:"DefinitionUri,omitempty"`

	// BinaryMediaTypes is a list of MIME types to treat as binary.
	BinaryMediaTypes []interface{} `json:"BinaryMediaTypes,omitempty" yaml:"BinaryMediaTypes,omitempty"`

	// MinimumCompressionSize allows compression for responses larger than specified bytes.
	MinimumCompressionSize int `json:"MinimumCompressionSize,omitempty" yaml:"MinimumCompressionSize,omitempty"`

	// EndpointConfiguration specifies the endpoint type (EDGE, REGIONAL, PRIVATE).
	EndpointConfiguration *EndpointConfig `json:"EndpointConfiguration,omitempty" yaml:"EndpointConfiguration,omitempty"`

	// CacheClusterEnabled enables a cache cluster for the stage.
	CacheClusterEnabled bool `json:"CacheClusterEnabled,omitempty" yaml:"CacheClusterEnabled,omitempty"`

	// CacheClusterSize specifies the cache cluster size.
	CacheClusterSize string `json:"CacheClusterSize,omitempty" yaml:"CacheClusterSize,omitempty"`

	// Variables is a map of stage variables.
	Variables map[string]interface{} `json:"Variables,omitempty" yaml:"Variables,omitempty"`

	// TracingEnabled enables X-Ray tracing for the stage.
	TracingEnabled bool `json:"TracingEnabled,omitempty" yaml:"TracingEnabled,omitempty"`

	// Tags is a map of key-value pairs to apply to the RestApi.
	Tags map[string]interface{} `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// AccessLogSetting specifies access logging configuration.
	AccessLogSetting *AccessLogSetting `json:"AccessLogSetting,omitempty" yaml:"AccessLogSetting,omitempty"`

	// MethodSettings specifies settings for API methods.
	MethodSettings []MethodSettingConfig `json:"MethodSettings,omitempty" yaml:"MethodSettings,omitempty"`

	// FailOnWarnings indicates whether to fail on warnings during API import.
	FailOnWarnings bool `json:"FailOnWarnings,omitempty" yaml:"FailOnWarnings,omitempty"`

	// DisableExecuteApiEndpoint disables the default execute-api endpoint.
	DisableExecuteApiEndpoint bool `json:"DisableExecuteApiEndpoint,omitempty" yaml:"DisableExecuteApiEndpoint,omitempty"`

	// OpenApiVersion specifies the OpenAPI version (e.g., "3.0.1").
	OpenApiVersion string `json:"OpenApiVersion,omitempty" yaml:"OpenApiVersion,omitempty"`

	// Cors specifies CORS configuration.
	Cors *CorsConfig `json:"Cors,omitempty" yaml:"Cors,omitempty"`

	// Auth specifies authentication configuration.
	Auth *ApiAuth `json:"Auth,omitempty" yaml:"Auth,omitempty"`

	// CanarySetting specifies canary deployment configuration.
	CanarySetting *CanarySettingConfig `json:"CanarySetting,omitempty" yaml:"CanarySetting,omitempty"`

	// GatewayResponses configures gateway responses.
	GatewayResponses map[string]interface{} `json:"GatewayResponses,omitempty" yaml:"GatewayResponses,omitempty"`

	// Models defines request/response models.
	Models map[string]interface{} `json:"Models,omitempty" yaml:"Models,omitempty"`

	// Domain specifies custom domain configuration.
	Domain *DomainConfig `json:"Domain,omitempty" yaml:"Domain,omitempty"`

	// ApiKeySourceType specifies the source of API keys (HEADER or AUTHORIZER).
	ApiKeySourceType string `json:"ApiKeySourceType,omitempty" yaml:"ApiKeySourceType,omitempty"`
}

// EndpointConfig represents API endpoint configuration.
type EndpointConfig struct {
	// Type is the endpoint type: EDGE, REGIONAL, or PRIVATE.
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`

	// VPCEndpointIds is a list of VPC endpoint IDs (for PRIVATE endpoints).
	VPCEndpointIds []interface{} `json:"VPCEndpointIds,omitempty" yaml:"VPCEndpointIds,omitempty"`
}

// AccessLogSetting represents access logging configuration.
type AccessLogSetting struct {
	// DestinationArn is the ARN of the CloudWatch Logs log group or Firehose stream.
	DestinationArn interface{} `json:"DestinationArn,omitempty" yaml:"DestinationArn,omitempty"`

	// Format is the log format string.
	Format interface{} `json:"Format,omitempty" yaml:"Format,omitempty"`
}

// MethodSettingConfig represents method-level settings.
type MethodSettingConfig struct {
	// HttpMethod is the HTTP method (* for all methods).
	HttpMethod string `json:"HttpMethod,omitempty" yaml:"HttpMethod,omitempty"`

	// ResourcePath is the resource path (/* for all paths).
	ResourcePath string `json:"ResourcePath,omitempty" yaml:"ResourcePath,omitempty"`

	// CachingEnabled enables caching for this method.
	CachingEnabled bool `json:"CachingEnabled,omitempty" yaml:"CachingEnabled,omitempty"`

	// CacheTtlInSeconds is the TTL for cached responses.
	CacheTtlInSeconds int `json:"CacheTtlInSeconds,omitempty" yaml:"CacheTtlInSeconds,omitempty"`

	// DataTraceEnabled enables data trace logging.
	DataTraceEnabled bool `json:"DataTraceEnabled,omitempty" yaml:"DataTraceEnabled,omitempty"`

	// LoggingLevel is the logging level (OFF, INFO, ERROR).
	LoggingLevel string `json:"LoggingLevel,omitempty" yaml:"LoggingLevel,omitempty"`

	// MetricsEnabled enables CloudWatch metrics.
	MetricsEnabled bool `json:"MetricsEnabled,omitempty" yaml:"MetricsEnabled,omitempty"`

	// ThrottlingBurstLimit is the throttling burst limit.
	ThrottlingBurstLimit int `json:"ThrottlingBurstLimit,omitempty" yaml:"ThrottlingBurstLimit,omitempty"`

	// ThrottlingRateLimit is the throttling rate limit.
	ThrottlingRateLimit float64 `json:"ThrottlingRateLimit,omitempty" yaml:"ThrottlingRateLimit,omitempty"`
}

// CorsConfig represents CORS configuration.
type CorsConfig struct {
	// AllowOrigin is the Access-Control-Allow-Origin header value.
	AllowOrigin interface{} `json:"AllowOrigin,omitempty" yaml:"AllowOrigin,omitempty"`

	// AllowMethods is the Access-Control-Allow-Methods header value.
	AllowMethods interface{} `json:"AllowMethods,omitempty" yaml:"AllowMethods,omitempty"`

	// AllowHeaders is the Access-Control-Allow-Headers header value.
	AllowHeaders interface{} `json:"AllowHeaders,omitempty" yaml:"AllowHeaders,omitempty"`

	// MaxAge is the Access-Control-Max-Age header value.
	MaxAge interface{} `json:"MaxAge,omitempty" yaml:"MaxAge,omitempty"`

	// AllowCredentials indicates if credentials are allowed.
	AllowCredentials bool `json:"AllowCredentials,omitempty" yaml:"AllowCredentials,omitempty"`
}

// ApiAuth represents API authentication configuration.
type ApiAuth struct {
	// DefaultAuthorizer is the default authorizer for all methods.
	DefaultAuthorizer string `json:"DefaultAuthorizer,omitempty" yaml:"DefaultAuthorizer,omitempty"`

	// Authorizers is a map of authorizer configurations.
	Authorizers map[string]interface{} `json:"Authorizers,omitempty" yaml:"Authorizers,omitempty"`

	// ApiKeyRequired indicates if API key is required by default.
	ApiKeyRequired bool `json:"ApiKeyRequired,omitempty" yaml:"ApiKeyRequired,omitempty"`

	// ResourcePolicy is the resource policy for the API.
	ResourcePolicy interface{} `json:"ResourcePolicy,omitempty" yaml:"ResourcePolicy,omitempty"`

	// UsagePlan configures usage plans.
	UsagePlan interface{} `json:"UsagePlan,omitempty" yaml:"UsagePlan,omitempty"`
}

// CanarySettingConfig represents canary deployment configuration.
type CanarySettingConfig struct {
	// PercentTraffic is the percentage of traffic to route to canary.
	PercentTraffic float64 `json:"PercentTraffic,omitempty" yaml:"PercentTraffic,omitempty"`

	// DeploymentId is the deployment ID for the canary.
	DeploymentId interface{} `json:"DeploymentId,omitempty" yaml:"DeploymentId,omitempty"`

	// StageVariableOverrides are stage variable overrides for canary.
	StageVariableOverrides map[string]interface{} `json:"StageVariableOverrides,omitempty" yaml:"StageVariableOverrides,omitempty"`

	// UseStageCache indicates if canary should use stage cache.
	UseStageCache bool `json:"UseStageCache,omitempty" yaml:"UseStageCache,omitempty"`
}

// DomainConfig represents custom domain configuration.
type DomainConfig struct {
	// DomainName is the custom domain name.
	DomainName interface{} `json:"DomainName,omitempty" yaml:"DomainName,omitempty"`

	// CertificateArn is the ACM certificate ARN.
	CertificateArn interface{} `json:"CertificateArn,omitempty" yaml:"CertificateArn,omitempty"`

	// EndpointConfiguration specifies the endpoint type.
	EndpointConfiguration interface{} `json:"EndpointConfiguration,omitempty" yaml:"EndpointConfiguration,omitempty"`

	// Route53 specifies Route53 configuration.
	Route53 interface{} `json:"Route53,omitempty" yaml:"Route53,omitempty"`

	// BasePath is the base path mapping.
	BasePath interface{} `json:"BasePath,omitempty" yaml:"BasePath,omitempty"`
}

// ApiTransformer transforms AWS::Serverless::Api to CloudFormation resources.
type ApiTransformer struct{}

// NewApiTransformer creates a new ApiTransformer.
func NewApiTransformer() *ApiTransformer {
	return &ApiTransformer{}
}

// Transform converts a SAM Api to CloudFormation resources.
func (t *ApiTransformer) Transform(logicalID string, api *Api) (map[string]interface{}, error) {
	// Validate required fields
	if api.StageName == nil || api.StageName == "" {
		return nil, fmt.Errorf("StageName is required for AWS::Serverless::Api")
	}

	resources := make(map[string]interface{})

	// Build RestApi resource
	restApiProps := make(map[string]interface{})

	// Set Name if provided
	if api.Name != nil {
		restApiProps["Name"] = api.Name
	}

	// Set Description if provided
	if api.Description != nil {
		restApiProps["Description"] = api.Description
	}

	// Handle DefinitionBody or DefinitionUri
	if api.DefinitionBody != nil {
		// DefinitionBody takes precedence
		restApiProps["Body"] = api.DefinitionBody
	} else if api.DefinitionUri != nil {
		s3Location, err := t.buildS3Location(api.DefinitionUri)
		if err != nil {
			return nil, fmt.Errorf("failed to build DefinitionUri: %w", err)
		}
		restApiProps["BodyS3Location"] = s3Location
	}

	// Set BinaryMediaTypes
	if len(api.BinaryMediaTypes) > 0 {
		restApiProps["BinaryMediaTypes"] = api.BinaryMediaTypes
	}

	// Set MinimumCompressionSize
	if api.MinimumCompressionSize > 0 {
		restApiProps["MinimumCompressionSize"] = api.MinimumCompressionSize
	}

	// Set EndpointConfiguration
	if api.EndpointConfiguration != nil {
		endpointConfig := make(map[string]interface{})
		if api.EndpointConfiguration.Type != "" {
			endpointConfig["Types"] = []interface{}{api.EndpointConfiguration.Type}
		}
		if len(api.EndpointConfiguration.VPCEndpointIds) > 0 {
			endpointConfig["VpcEndpointIds"] = api.EndpointConfiguration.VPCEndpointIds
		}
		restApiProps["EndpointConfiguration"] = endpointConfig
	}

	// Set FailOnWarnings
	if api.FailOnWarnings {
		restApiProps["FailOnWarnings"] = true
	}

	// Set DisableExecuteApiEndpoint
	if api.DisableExecuteApiEndpoint {
		restApiProps["DisableExecuteApiEndpoint"] = true
	}

	// Set ApiKeySourceType
	if api.ApiKeySourceType != "" {
		restApiProps["ApiKeySourceType"] = api.ApiKeySourceType
	}

	// Build Tags
	if len(api.Tags) > 0 {
		tags := make([]map[string]interface{}, 0, len(api.Tags))
		// Sort keys for deterministic output
		keys := make([]string, 0, len(api.Tags))
		for k := range api.Tags {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			tags = append(tags, map[string]interface{}{
				"Key":   k,
				"Value": api.Tags[k],
			})
		}
		restApiProps["Tags"] = tags
	}

	// Add RestApi resource
	resources[logicalID] = map[string]interface{}{
		"Type":       "AWS::ApiGateway::RestApi",
		"Properties": restApiProps,
	}

	// Generate deployment ID based on definition content for stability
	deploymentLogicalID := t.generateDeploymentLogicalID(logicalID, api)

	// Build Deployment resource
	deploymentProps := map[string]interface{}{
		"RestApiId": map[string]interface{}{
			"Ref": logicalID,
		},
	}

	resources[deploymentLogicalID] = map[string]interface{}{
		"Type":       "AWS::ApiGateway::Deployment",
		"Properties": deploymentProps,
	}

	// Build Stage resource
	stageLogicalID := logicalID + "Stage"
	stageProps := map[string]interface{}{
		"RestApiId": map[string]interface{}{
			"Ref": logicalID,
		},
		"StageName": api.StageName,
		"DeploymentId": map[string]interface{}{
			"Ref": deploymentLogicalID,
		},
	}

	// Set stage variables
	if len(api.Variables) > 0 {
		stageProps["Variables"] = api.Variables
	}

	// Set cache settings
	if api.CacheClusterEnabled {
		stageProps["CacheClusterEnabled"] = true
	}
	if api.CacheClusterSize != "" {
		stageProps["CacheClusterSize"] = api.CacheClusterSize
	}

	// Set tracing
	if api.TracingEnabled {
		stageProps["TracingEnabled"] = true
	}

	// Set access log settings
	if api.AccessLogSetting != nil {
		accessLogSetting := make(map[string]interface{})
		if api.AccessLogSetting.DestinationArn != nil {
			accessLogSetting["DestinationArn"] = api.AccessLogSetting.DestinationArn
		}
		if api.AccessLogSetting.Format != nil {
			accessLogSetting["Format"] = api.AccessLogSetting.Format
		}
		stageProps["AccessLogSetting"] = accessLogSetting
	}

	// Set method settings
	if len(api.MethodSettings) > 0 {
		methodSettings := make([]map[string]interface{}, 0, len(api.MethodSettings))
		for _, ms := range api.MethodSettings {
			setting := make(map[string]interface{})
			if ms.HttpMethod != "" {
				setting["HttpMethod"] = ms.HttpMethod
			}
			if ms.ResourcePath != "" {
				setting["ResourcePath"] = ms.ResourcePath
			}
			if ms.CachingEnabled {
				setting["CachingEnabled"] = true
			}
			if ms.CacheTtlInSeconds > 0 {
				setting["CacheTtlInSeconds"] = ms.CacheTtlInSeconds
			}
			if ms.DataTraceEnabled {
				setting["DataTraceEnabled"] = true
			}
			if ms.LoggingLevel != "" {
				setting["LoggingLevel"] = ms.LoggingLevel
			}
			if ms.MetricsEnabled {
				setting["MetricsEnabled"] = true
			}
			if ms.ThrottlingBurstLimit > 0 {
				setting["ThrottlingBurstLimit"] = ms.ThrottlingBurstLimit
			}
			if ms.ThrottlingRateLimit > 0 {
				setting["ThrottlingRateLimit"] = ms.ThrottlingRateLimit
			}
			methodSettings = append(methodSettings, setting)
		}
		stageProps["MethodSettings"] = methodSettings
	}

	// Set canary settings
	if api.CanarySetting != nil {
		canarySetting := make(map[string]interface{})
		if api.CanarySetting.PercentTraffic > 0 {
			canarySetting["PercentTraffic"] = api.CanarySetting.PercentTraffic
		}
		if api.CanarySetting.DeploymentId != nil {
			canarySetting["DeploymentId"] = api.CanarySetting.DeploymentId
		}
		if len(api.CanarySetting.StageVariableOverrides) > 0 {
			canarySetting["StageVariableOverrides"] = api.CanarySetting.StageVariableOverrides
		}
		if api.CanarySetting.UseStageCache {
			canarySetting["UseStageCache"] = true
		}
		stageProps["CanarySetting"] = canarySetting
	}

	resources[stageLogicalID] = map[string]interface{}{
		"Type":       "AWS::ApiGateway::Stage",
		"Properties": stageProps,
	}

	// Process authorizers if Auth is configured
	if api.Auth != nil && len(api.Auth.Authorizers) > 0 {
		if err := t.processAuthorizers(logicalID, api.Auth, resources); err != nil {
			return nil, fmt.Errorf("failed to process authorizers: %w", err)
		}
	}

	return resources, nil
}

// buildS3Location builds an S3 location from DefinitionUri.
func (t *ApiTransformer) buildS3Location(definitionUri interface{}) (map[string]interface{}, error) {
	switch uri := definitionUri.(type) {
	case string:
		return t.parseS3Uri(uri)
	case map[string]interface{}:
		s3Location := make(map[string]interface{})
		if bucket, ok := uri["Bucket"]; ok {
			s3Location["Bucket"] = bucket
		}
		if key, ok := uri["Key"]; ok {
			s3Location["Key"] = key
		}
		if version, ok := uri["Version"]; ok {
			s3Location["Version"] = version
		}
		return s3Location, nil
	default:
		return nil, fmt.Errorf("unsupported DefinitionUri type: %T", definitionUri)
	}
}

// parseS3Uri parses an S3 URI string into an S3 location map.
func (t *ApiTransformer) parseS3Uri(uri string) (map[string]interface{}, error) {
	if !strings.HasPrefix(uri, "s3://") {
		return nil, fmt.Errorf("DefinitionUri must be an S3 URI (s3://...): %s", uri)
	}

	path := strings.TrimPrefix(uri, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid S3 URI (missing key): %s", uri)
	}

	return map[string]interface{}{
		"Bucket": parts[0],
		"Key":    parts[1],
	}, nil
}

// generateDeploymentLogicalID generates a stable deployment logical ID based on the API definition.
func (t *ApiTransformer) generateDeploymentLogicalID(logicalID string, api *Api) string {
	// Create a hash of the definition for stable deployment IDs
	h := sha256.New()

	// Include definition body if present
	if api.DefinitionBody != nil {
		jsonBytes, err := json.Marshal(api.DefinitionBody)
		if err == nil {
			h.Write(jsonBytes)
		}
	}

	// Include definition URI if present
	if api.DefinitionUri != nil {
		switch uri := api.DefinitionUri.(type) {
		case string:
			h.Write([]byte(uri))
		case map[string]interface{}:
			jsonBytes, err := json.Marshal(uri)
			if err == nil {
				h.Write(jsonBytes)
			}
		}
	}

	hash := hex.EncodeToString(h.Sum(nil))[:10]
	return logicalID + "Deployment" + hash
}

// processAuthorizers creates API Gateway Authorizer resources.
func (t *ApiTransformer) processAuthorizers(apiLogicalID string, auth *ApiAuth, resources map[string]interface{}) error {
	for name, config := range auth.Authorizers {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			continue
		}

		authorizerLogicalID := apiLogicalID + name + "Authorizer"
		authorizerProps := map[string]interface{}{
			"Name": name,
			"RestApiId": map[string]interface{}{
				"Ref": apiLogicalID,
			},
		}

		// Determine authorizer type
		if _, hasCognito := configMap["UserPoolArn"]; hasCognito {
			authorizerProps["Type"] = "COGNITO_USER_POOLS"
			authorizerProps["ProviderARNs"] = []interface{}{configMap["UserPoolArn"]}
			authorizerProps["IdentitySource"] = "method.request.header.Authorization"
		} else if _, hasLambda := configMap["FunctionArn"]; hasLambda {
			// Lambda authorizer
			if identity, ok := configMap["Identity"].(map[string]interface{}); ok {
				if _, hasHeaders := identity["Headers"]; hasHeaders {
					authorizerProps["Type"] = "REQUEST"
				} else {
					authorizerProps["Type"] = "TOKEN"
				}
			} else {
				authorizerProps["Type"] = "TOKEN"
			}
			authorizerProps["AuthorizerUri"] = t.buildLambdaAuthorizerUri(configMap["FunctionArn"])
		}

		// Set identity source if specified
		if identitySource, ok := configMap["IdentitySource"].(string); ok {
			authorizerProps["IdentitySource"] = identitySource
		}

		// Set identity validation expression if specified
		if validationExpr, ok := configMap["IdentityValidationExpression"].(string); ok {
			authorizerProps["IdentityValidationExpression"] = validationExpr
		}

		// Set TTL if specified
		if ttl, ok := configMap["AuthorizerResultTtlInSeconds"]; ok {
			authorizerProps["AuthorizerResultTtlInSeconds"] = ttl
		}

		resources[authorizerLogicalID] = map[string]interface{}{
			"Type":       "AWS::ApiGateway::Authorizer",
			"Properties": authorizerProps,
		}
	}

	return nil
}

// buildLambdaAuthorizerUri builds the Lambda authorizer URI.
func (t *ApiTransformer) buildLambdaAuthorizerUri(functionArn interface{}) interface{} {
	// Return Fn::Sub for the authorizer URI
	return map[string]interface{}{
		"Fn::Sub": []interface{}{
			"arn:aws:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${FunctionArn}/invocations",
			map[string]interface{}{
				"FunctionArn": functionArn,
			},
		},
	}
}
