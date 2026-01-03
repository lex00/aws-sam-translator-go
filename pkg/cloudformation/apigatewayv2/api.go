// Package apigatewayv2 provides CloudFormation resource models for AWS API Gateway V2 (HTTP API/WebSocket API).
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html
package apigatewayv2

// Api represents an AWS::ApiGatewayV2::Api CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-api.html
type Api struct {
	// ApiKeySelectionExpression is the API key selection expression for a WebSocket API.
	ApiKeySelectionExpression interface{} `json:"ApiKeySelectionExpression,omitempty" yaml:"ApiKeySelectionExpression,omitempty"`

	// BasePath specifies how to interpret the base path of the API during import.
	BasePath interface{} `json:"BasePath,omitempty" yaml:"BasePath,omitempty"`

	// Body is the OpenAPI definition. Used for import operations.
	Body interface{} `json:"Body,omitempty" yaml:"Body,omitempty"`

	// BodyS3Location specifies the Amazon S3 location of the OpenAPI definition.
	BodyS3Location *BodyS3Location `json:"BodyS3Location,omitempty" yaml:"BodyS3Location,omitempty"`

	// CorsConfiguration specifies a CORS configuration for an API.
	CorsConfiguration *Cors `json:"CorsConfiguration,omitempty" yaml:"CorsConfiguration,omitempty"`

	// CredentialsArn specifies the credentials required for the integration.
	CredentialsArn interface{} `json:"CredentialsArn,omitempty" yaml:"CredentialsArn,omitempty"`

	// Description is the description of the API.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// DisableExecuteApiEndpoint specifies whether clients can invoke your API using the default execute-api endpoint.
	DisableExecuteApiEndpoint interface{} `json:"DisableExecuteApiEndpoint,omitempty" yaml:"DisableExecuteApiEndpoint,omitempty"`

	// DisableSchemaValidation specifies whether schema validation is disabled.
	DisableSchemaValidation interface{} `json:"DisableSchemaValidation,omitempty" yaml:"DisableSchemaValidation,omitempty"`

	// FailOnWarnings specifies whether to rollback the API creation if a warning is encountered.
	FailOnWarnings interface{} `json:"FailOnWarnings,omitempty" yaml:"FailOnWarnings,omitempty"`

	// Name is the name of the API.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// ProtocolType is the API protocol. Valid values: HTTP, WEBSOCKET
	ProtocolType interface{} `json:"ProtocolType,omitempty" yaml:"ProtocolType,omitempty"`

	// RouteKey is the route key for the route (for quick create APIs).
	RouteKey interface{} `json:"RouteKey,omitempty" yaml:"RouteKey,omitempty"`

	// RouteSelectionExpression is the route selection expression for the API.
	RouteSelectionExpression interface{} `json:"RouteSelectionExpression,omitempty" yaml:"RouteSelectionExpression,omitempty"`

	// Tags is a collection of tags associated with the API.
	Tags map[string]interface{} `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// Target is the target URL for quick create APIs.
	Target interface{} `json:"Target,omitempty" yaml:"Target,omitempty"`

	// Version is the version identifier for the API.
	Version interface{} `json:"Version,omitempty" yaml:"Version,omitempty"`
}

// BodyS3Location represents an Amazon S3 location for the OpenAPI definition.
type BodyS3Location struct {
	// Bucket is the name of the S3 bucket where the OpenAPI file is stored.
	Bucket interface{} `json:"Bucket,omitempty" yaml:"Bucket,omitempty"`

	// Etag is the ETag of the S3 object.
	Etag interface{} `json:"Etag,omitempty" yaml:"Etag,omitempty"`

	// Key is the name of the S3 object.
	Key interface{} `json:"Key,omitempty" yaml:"Key,omitempty"`

	// Version is the version ID of the S3 object.
	Version interface{} `json:"Version,omitempty" yaml:"Version,omitempty"`
}

// Cors represents CORS configuration for an API.
type Cors struct {
	// AllowCredentials specifies whether credentials are included in the CORS request.
	AllowCredentials interface{} `json:"AllowCredentials,omitempty" yaml:"AllowCredentials,omitempty"`

	// AllowHeaders specifies the allowed headers.
	AllowHeaders []interface{} `json:"AllowHeaders,omitempty" yaml:"AllowHeaders,omitempty"`

	// AllowMethods specifies the allowed HTTP methods.
	AllowMethods []interface{} `json:"AllowMethods,omitempty" yaml:"AllowMethods,omitempty"`

	// AllowOrigins specifies the allowed origins.
	AllowOrigins []interface{} `json:"AllowOrigins,omitempty" yaml:"AllowOrigins,omitempty"`

	// ExposeHeaders specifies the exposed headers.
	ExposeHeaders []interface{} `json:"ExposeHeaders,omitempty" yaml:"ExposeHeaders,omitempty"`

	// MaxAge specifies the number of seconds that the browser should cache preflight results.
	MaxAge interface{} `json:"MaxAge,omitempty" yaml:"MaxAge,omitempty"`
}
