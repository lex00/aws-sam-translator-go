// Package apigateway provides CloudFormation resource models for AWS API Gateway (REST API).
package apigateway

// Stage represents an AWS::ApiGateway::Stage CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-stage.html
type Stage struct {
	// AccessLogSetting specifies settings for logging access in this stage.
	AccessLogSetting *AccessLogSetting `json:"AccessLogSetting,omitempty" yaml:"AccessLogSetting,omitempty"`

	// CacheClusterEnabled indicates whether cache clustering is enabled for the stage.
	CacheClusterEnabled interface{} `json:"CacheClusterEnabled,omitempty" yaml:"CacheClusterEnabled,omitempty"`

	// CacheClusterSize specifies the stage's cache cluster size.
	CacheClusterSize interface{} `json:"CacheClusterSize,omitempty" yaml:"CacheClusterSize,omitempty"`

	// CanarySetting specifies settings for the canary deployment in this stage.
	CanarySetting *CanarySetting `json:"CanarySetting,omitempty" yaml:"CanarySetting,omitempty"`

	// ClientCertificateId is the ID of the client certificate for the stage.
	ClientCertificateId interface{} `json:"ClientCertificateId,omitempty" yaml:"ClientCertificateId,omitempty"`

	// DeploymentId is the ID of the deployment that the stage points to.
	DeploymentId interface{} `json:"DeploymentId,omitempty" yaml:"DeploymentId,omitempty"`

	// Description is the description of the stage.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// DocumentationVersion specifies the version of the associated API documentation.
	DocumentationVersion interface{} `json:"DocumentationVersion,omitempty" yaml:"DocumentationVersion,omitempty"`

	// MethodSettings specifies settings for all methods in the stage.
	MethodSettings []MethodSetting `json:"MethodSettings,omitempty" yaml:"MethodSettings,omitempty"`

	// RestApiId is the ID of the RestApi resource that you're deploying with this stage.
	RestApiId interface{} `json:"RestApiId" yaml:"RestApiId"`

	// StageName is the name of the stage.
	StageName interface{} `json:"StageName,omitempty" yaml:"StageName,omitempty"`

	// Tags is a map of key-value pairs to associate with the stage.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// TracingEnabled indicates whether active tracing with X-Ray is enabled for the stage.
	TracingEnabled interface{} `json:"TracingEnabled,omitempty" yaml:"TracingEnabled,omitempty"`

	// Variables is a map of stage variables.
	Variables map[string]interface{} `json:"Variables,omitempty" yaml:"Variables,omitempty"`
}

// AccessLogSetting represents access logging settings for a stage.
type AccessLogSetting struct {
	// DestinationArn is the ARN of the CloudWatch Logs log group or Kinesis Data Firehose delivery stream.
	DestinationArn interface{} `json:"DestinationArn,omitempty" yaml:"DestinationArn,omitempty"`

	// Format is the format of the access logs.
	Format interface{} `json:"Format,omitempty" yaml:"Format,omitempty"`
}

// CanarySetting represents canary deployment settings for a stage.
type CanarySetting struct {
	// DeploymentId is the ID of the canary deployment.
	DeploymentId interface{} `json:"DeploymentId,omitempty" yaml:"DeploymentId,omitempty"`

	// PercentTraffic is the percentage of traffic to divert to the canary deployment.
	PercentTraffic interface{} `json:"PercentTraffic,omitempty" yaml:"PercentTraffic,omitempty"`

	// StageVariableOverrides are stage variable overrides for the canary deployment.
	StageVariableOverrides map[string]interface{} `json:"StageVariableOverrides,omitempty" yaml:"StageVariableOverrides,omitempty"`

	// UseStageCache indicates whether to use the stage cache for the canary deployment.
	UseStageCache interface{} `json:"UseStageCache,omitempty" yaml:"UseStageCache,omitempty"`
}

// MethodSetting represents settings for a method in a stage.
type MethodSetting struct {
	// CacheDataEncrypted indicates whether the cached responses are encrypted.
	CacheDataEncrypted interface{} `json:"CacheDataEncrypted,omitempty" yaml:"CacheDataEncrypted,omitempty"`

	// CacheTtlInSeconds is the time-to-live (TTL) period, in seconds, for cached responses.
	CacheTtlInSeconds interface{} `json:"CacheTtlInSeconds,omitempty" yaml:"CacheTtlInSeconds,omitempty"`

	// CachingEnabled indicates whether responses are cached.
	CachingEnabled interface{} `json:"CachingEnabled,omitempty" yaml:"CachingEnabled,omitempty"`

	// DataTraceEnabled indicates whether data trace logging is enabled for methods in the stage.
	DataTraceEnabled interface{} `json:"DataTraceEnabled,omitempty" yaml:"DataTraceEnabled,omitempty"`

	// HttpMethod is the HTTP method. Use * for all methods.
	HttpMethod interface{} `json:"HttpMethod,omitempty" yaml:"HttpMethod,omitempty"`

	// LoggingLevel specifies the logging level for this method.
	// Valid values: OFF, INFO, ERROR
	LoggingLevel interface{} `json:"LoggingLevel,omitempty" yaml:"LoggingLevel,omitempty"`

	// MetricsEnabled indicates whether CloudWatch metrics are enabled for methods in the stage.
	MetricsEnabled interface{} `json:"MetricsEnabled,omitempty" yaml:"MetricsEnabled,omitempty"`

	// ResourcePath is the resource path for this method. Use /* for all paths.
	ResourcePath interface{} `json:"ResourcePath,omitempty" yaml:"ResourcePath,omitempty"`

	// ThrottlingBurstLimit is the throttling burst limit.
	ThrottlingBurstLimit interface{} `json:"ThrottlingBurstLimit,omitempty" yaml:"ThrottlingBurstLimit,omitempty"`

	// ThrottlingRateLimit is the throttling rate limit.
	ThrottlingRateLimit interface{} `json:"ThrottlingRateLimit,omitempty" yaml:"ThrottlingRateLimit,omitempty"`
}
