// Package apigatewayv2 provides CloudFormation resource models for AWS API Gateway V2 (HTTP API/WebSocket API).
package apigatewayv2

// Authorizer represents an AWS::ApiGatewayV2::Authorizer CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-authorizer.html
type Authorizer struct {
	// ApiId is the API identifier.
	ApiId interface{} `json:"ApiId" yaml:"ApiId"`

	// AuthorizerCredentialsArn specifies the required credentials as an IAM role for API Gateway.
	AuthorizerCredentialsArn interface{} `json:"AuthorizerCredentialsArn,omitempty" yaml:"AuthorizerCredentialsArn,omitempty"`

	// AuthorizerPayloadFormatVersion specifies the format of the payload sent to an HTTP API Lambda authorizer.
	// Valid values: 1.0, 2.0
	AuthorizerPayloadFormatVersion interface{} `json:"AuthorizerPayloadFormatVersion,omitempty" yaml:"AuthorizerPayloadFormatVersion,omitempty"`

	// AuthorizerResultTtlInSeconds is the time to live (TTL), in seconds, of cached authorizer results.
	AuthorizerResultTtlInSeconds interface{} `json:"AuthorizerResultTtlInSeconds,omitempty" yaml:"AuthorizerResultTtlInSeconds,omitempty"`

	// AuthorizerType is the authorizer type.
	// Valid values: REQUEST, JWT
	AuthorizerType interface{} `json:"AuthorizerType" yaml:"AuthorizerType"`

	// AuthorizerUri is the authorizer's Uniform Resource Identifier (URI).
	AuthorizerUri interface{} `json:"AuthorizerUri,omitempty" yaml:"AuthorizerUri,omitempty"`

	// EnableSimpleResponses specifies whether a Lambda authorizer returns a response in a simple format.
	EnableSimpleResponses interface{} `json:"EnableSimpleResponses,omitempty" yaml:"EnableSimpleResponses,omitempty"`

	// IdentitySource is the identity source for which authorization is requested.
	IdentitySource []interface{} `json:"IdentitySource,omitempty" yaml:"IdentitySource,omitempty"`

	// IdentityValidationExpression is the validation expression for the identity.
	IdentityValidationExpression interface{} `json:"IdentityValidationExpression,omitempty" yaml:"IdentityValidationExpression,omitempty"`

	// JwtConfiguration specifies the configuration of a JWT authorizer.
	JwtConfiguration *JWTConfiguration `json:"JwtConfiguration,omitempty" yaml:"JwtConfiguration,omitempty"`

	// Name is the name of the authorizer.
	Name interface{} `json:"Name" yaml:"Name"`
}

// JWTConfiguration represents the JWT configuration for an authorizer.
type JWTConfiguration struct {
	// Audience is a list of the intended recipients of the JWT.
	Audience []interface{} `json:"Audience,omitempty" yaml:"Audience,omitempty"`

	// Issuer is the base domain of the identity provider that issues JWTs.
	Issuer interface{} `json:"Issuer,omitempty" yaml:"Issuer,omitempty"`
}
