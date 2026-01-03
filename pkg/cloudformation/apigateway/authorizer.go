// Package apigateway provides CloudFormation resource models for AWS API Gateway (REST API).
package apigateway

// Authorizer represents an AWS::ApiGateway::Authorizer CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html
type Authorizer struct {
	// AuthType is an optional customer-defined field used in OpenAPI imports and exports.
	AuthType interface{} `json:"AuthType,omitempty" yaml:"AuthType,omitempty"`

	// AuthorizerCredentials specifies the credentials required for the authorizer.
	AuthorizerCredentials interface{} `json:"AuthorizerCredentials,omitempty" yaml:"AuthorizerCredentials,omitempty"`

	// AuthorizerResultTtlInSeconds is the TTL of cached authorizer results in seconds.
	AuthorizerResultTtlInSeconds interface{} `json:"AuthorizerResultTtlInSeconds,omitempty" yaml:"AuthorizerResultTtlInSeconds,omitempty"`

	// AuthorizerUri is the authorizer's Uniform Resource Identifier (URI).
	AuthorizerUri interface{} `json:"AuthorizerUri,omitempty" yaml:"AuthorizerUri,omitempty"`

	// IdentitySource is the source of the identity in an incoming request.
	IdentitySource interface{} `json:"IdentitySource,omitempty" yaml:"IdentitySource,omitempty"`

	// IdentityValidationExpression is a validation expression for the incoming identity.
	IdentityValidationExpression interface{} `json:"IdentityValidationExpression,omitempty" yaml:"IdentityValidationExpression,omitempty"`

	// Name is the name of the authorizer.
	Name interface{} `json:"Name" yaml:"Name"`

	// ProviderARNs is a list of Amazon Cognito user pool ARNs for the COGNITO_USER_POOLS authorizer.
	ProviderARNs []interface{} `json:"ProviderARNs,omitempty" yaml:"ProviderARNs,omitempty"`

	// RestApiId is the ID of the RestApi resource that API Gateway creates the authorizer in.
	RestApiId interface{} `json:"RestApiId" yaml:"RestApiId"`

	// Type is the type of authorizer.
	// Valid values: TOKEN, REQUEST, COGNITO_USER_POOLS
	Type interface{} `json:"Type" yaml:"Type"`
}
