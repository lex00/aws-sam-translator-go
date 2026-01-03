// Package apigatewayv2 provides CloudFormation resource models for AWS API Gateway V2 (HTTP API/WebSocket API).
package apigatewayv2

// Route represents an AWS::ApiGatewayV2::Route CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-route.html
type Route struct {
	// ApiId is the API identifier.
	ApiId interface{} `json:"ApiId" yaml:"ApiId"`

	// ApiKeyRequired specifies whether an API key is required for the route.
	ApiKeyRequired interface{} `json:"ApiKeyRequired,omitempty" yaml:"ApiKeyRequired,omitempty"`

	// AuthorizationScopes is a list of authorization scopes for the route.
	AuthorizationScopes []interface{} `json:"AuthorizationScopes,omitempty" yaml:"AuthorizationScopes,omitempty"`

	// AuthorizationType is the authorization type for the route.
	// Valid values: NONE, AWS_IAM, CUSTOM, JWT
	AuthorizationType interface{} `json:"AuthorizationType,omitempty" yaml:"AuthorizationType,omitempty"`

	// AuthorizerId is the identifier of the Authorizer resource.
	AuthorizerId interface{} `json:"AuthorizerId,omitempty" yaml:"AuthorizerId,omitempty"`

	// ModelSelectionExpression is the model selection expression for the route.
	ModelSelectionExpression interface{} `json:"ModelSelectionExpression,omitempty" yaml:"ModelSelectionExpression,omitempty"`

	// OperationName is the operation name for the route.
	OperationName interface{} `json:"OperationName,omitempty" yaml:"OperationName,omitempty"`

	// RequestModels specifies the request models for the route.
	RequestModels map[string]interface{} `json:"RequestModels,omitempty" yaml:"RequestModels,omitempty"`

	// RequestParameters specifies the request parameters for the route.
	RequestParameters map[string]interface{} `json:"RequestParameters,omitempty" yaml:"RequestParameters,omitempty"`

	// RouteKey is the route key for the route.
	RouteKey interface{} `json:"RouteKey" yaml:"RouteKey"`

	// RouteResponseSelectionExpression is the route response selection expression for the route.
	RouteResponseSelectionExpression interface{} `json:"RouteResponseSelectionExpression,omitempty" yaml:"RouteResponseSelectionExpression,omitempty"`

	// Target is the target for the route.
	Target interface{} `json:"Target,omitempty" yaml:"Target,omitempty"`
}
