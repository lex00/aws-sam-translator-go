// Package apigatewayv2 provides CloudFormation resource models for AWS API Gateway V2 (HTTP API/WebSocket API).
package apigatewayv2

// Stage represents an AWS::ApiGatewayV2::Stage CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-stage.html
type Stage struct {
	// AccessLogSettings specifies settings for logging access in this stage.
	AccessLogSettings *AccessLogSettings `json:"AccessLogSettings,omitempty" yaml:"AccessLogSettings,omitempty"`

	// AccessPolicyId is the identifier of the AccessPolicy that is used for this stage.
	AccessPolicyId interface{} `json:"AccessPolicyId,omitempty" yaml:"AccessPolicyId,omitempty"`

	// ApiId is the API identifier.
	ApiId interface{} `json:"ApiId" yaml:"ApiId"`

	// AutoDeploy specifies whether updates to an API automatically trigger a new deployment.
	AutoDeploy interface{} `json:"AutoDeploy,omitempty" yaml:"AutoDeploy,omitempty"`

	// ClientCertificateId is the ID of the client certificate for the stage.
	ClientCertificateId interface{} `json:"ClientCertificateId,omitempty" yaml:"ClientCertificateId,omitempty"`

	// DefaultRouteSettings specifies the default route settings for the stage.
	DefaultRouteSettings *RouteSettings `json:"DefaultRouteSettings,omitempty" yaml:"DefaultRouteSettings,omitempty"`

	// DeploymentId is the deployment identifier for the API stage.
	DeploymentId interface{} `json:"DeploymentId,omitempty" yaml:"DeploymentId,omitempty"`

	// Description is the description of the stage.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// RouteSettings is a map of route settings for the stage.
	RouteSettings map[string]interface{} `json:"RouteSettings,omitempty" yaml:"RouteSettings,omitempty"`

	// StageName is the name of the stage.
	StageName interface{} `json:"StageName" yaml:"StageName"`

	// StageVariables is a map of stage variables.
	StageVariables map[string]interface{} `json:"StageVariables,omitempty" yaml:"StageVariables,omitempty"`

	// Tags is a collection of tags associated with the stage.
	Tags map[string]interface{} `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// AccessLogSettings represents access log settings for a stage.
type AccessLogSettings struct {
	// DestinationArn is the ARN of the CloudWatch Logs log group to receive access logs.
	DestinationArn interface{} `json:"DestinationArn,omitempty" yaml:"DestinationArn,omitempty"`

	// Format is the format of the access logs.
	Format interface{} `json:"Format,omitempty" yaml:"Format,omitempty"`
}

// RouteSettings represents route settings for a stage.
type RouteSettings struct {
	// DataTraceEnabled specifies whether data trace logging is enabled.
	DataTraceEnabled interface{} `json:"DataTraceEnabled,omitempty" yaml:"DataTraceEnabled,omitempty"`

	// DetailedMetricsEnabled specifies whether detailed metrics are enabled.
	DetailedMetricsEnabled interface{} `json:"DetailedMetricsEnabled,omitempty" yaml:"DetailedMetricsEnabled,omitempty"`

	// LoggingLevel specifies the logging level.
	// Valid values: ERROR, INFO, OFF
	LoggingLevel interface{} `json:"LoggingLevel,omitempty" yaml:"LoggingLevel,omitempty"`

	// ThrottlingBurstLimit is the throttling burst limit.
	ThrottlingBurstLimit interface{} `json:"ThrottlingBurstLimit,omitempty" yaml:"ThrottlingBurstLimit,omitempty"`

	// ThrottlingRateLimit is the throttling rate limit.
	ThrottlingRateLimit interface{} `json:"ThrottlingRateLimit,omitempty" yaml:"ThrottlingRateLimit,omitempty"`
}
