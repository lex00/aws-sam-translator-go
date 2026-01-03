// Package apigateway provides CloudFormation resource models for AWS API Gateway (REST API).
package apigateway

// Method represents an AWS::ApiGateway::Method CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-method.html
type Method struct {
	// ApiKeyRequired indicates whether the method requires an API key.
	ApiKeyRequired interface{} `json:"ApiKeyRequired,omitempty" yaml:"ApiKeyRequired,omitempty"`

	// AuthorizationScopes is a list of authorization scopes configured on the method.
	AuthorizationScopes []interface{} `json:"AuthorizationScopes,omitempty" yaml:"AuthorizationScopes,omitempty"`

	// AuthorizationType is the method's authorization type.
	// Valid values: NONE, AWS_IAM, CUSTOM, COGNITO_USER_POOLS
	AuthorizationType interface{} `json:"AuthorizationType,omitempty" yaml:"AuthorizationType,omitempty"`

	// AuthorizerId is the identifier of an Authorizer resource to use on this method.
	AuthorizerId interface{} `json:"AuthorizerId,omitempty" yaml:"AuthorizerId,omitempty"`

	// HttpMethod is the HTTP method that clients use to call this method.
	HttpMethod interface{} `json:"HttpMethod" yaml:"HttpMethod"`

	// Integration specifies the backend system that the method calls.
	Integration *Integration `json:"Integration,omitempty" yaml:"Integration,omitempty"`

	// MethodResponses specifies the responses that can be sent to the client.
	MethodResponses []MethodResponse `json:"MethodResponses,omitempty" yaml:"MethodResponses,omitempty"`

	// OperationName is a human-friendly operation identifier for the method.
	OperationName interface{} `json:"OperationName,omitempty" yaml:"OperationName,omitempty"`

	// RequestModels specifies the resources that are used for the request's content type.
	RequestModels map[string]interface{} `json:"RequestModels,omitempty" yaml:"RequestModels,omitempty"`

	// RequestParameters specifies the request parameters that API Gateway accepts.
	RequestParameters map[string]interface{} `json:"RequestParameters,omitempty" yaml:"RequestParameters,omitempty"`

	// RequestValidatorId is the ID of the associated request validator.
	RequestValidatorId interface{} `json:"RequestValidatorId,omitempty" yaml:"RequestValidatorId,omitempty"`

	// ResourceId is the ID of an API Gateway resource.
	ResourceId interface{} `json:"ResourceId" yaml:"ResourceId"`

	// RestApiId is the ID of the RestApi resource that API Gateway creates the method in.
	RestApiId interface{} `json:"RestApiId" yaml:"RestApiId"`
}

// Integration represents the backend integration for an API Gateway method.
type Integration struct {
	// CacheKeyParameters is a list of request parameters whose values API Gateway caches.
	CacheKeyParameters []interface{} `json:"CacheKeyParameters,omitempty" yaml:"CacheKeyParameters,omitempty"`

	// CacheNamespace is an API-specific tag group of related cached parameters.
	CacheNamespace interface{} `json:"CacheNamespace,omitempty" yaml:"CacheNamespace,omitempty"`

	// ConnectionId is the ID of the VpcLink used for the integration.
	ConnectionId interface{} `json:"ConnectionId,omitempty" yaml:"ConnectionId,omitempty"`

	// ConnectionType is the type of the network connection to the integration endpoint.
	// Valid values: INTERNET, VPC_LINK
	ConnectionType interface{} `json:"ConnectionType,omitempty" yaml:"ConnectionType,omitempty"`

	// ContentHandling specifies how to handle request payload content type conversions.
	// Valid values: CONVERT_TO_BINARY, CONVERT_TO_TEXT
	ContentHandling interface{} `json:"ContentHandling,omitempty" yaml:"ContentHandling,omitempty"`

	// Credentials specifies the credentials required for the integration.
	Credentials interface{} `json:"Credentials,omitempty" yaml:"Credentials,omitempty"`

	// IntegrationHttpMethod is the HTTP method used in the integration request.
	IntegrationHttpMethod interface{} `json:"IntegrationHttpMethod,omitempty" yaml:"IntegrationHttpMethod,omitempty"`

	// IntegrationResponses specifies the integration's responses.
	IntegrationResponses []IntegrationResponse `json:"IntegrationResponses,omitempty" yaml:"IntegrationResponses,omitempty"`

	// PassthroughBehavior indicates when API Gateway passes requests to the backend.
	// Valid values: WHEN_NO_MATCH, WHEN_NO_TEMPLATES, NEVER
	PassthroughBehavior interface{} `json:"PassthroughBehavior,omitempty" yaml:"PassthroughBehavior,omitempty"`

	// RequestParameters specifies the request parameters that API Gateway sends to the backend.
	RequestParameters map[string]interface{} `json:"RequestParameters,omitempty" yaml:"RequestParameters,omitempty"`

	// RequestTemplates specifies request payload mapping templates.
	RequestTemplates map[string]interface{} `json:"RequestTemplates,omitempty" yaml:"RequestTemplates,omitempty"`

	// TimeoutInMillis is the custom timeout in milliseconds (50 to 29000).
	TimeoutInMillis interface{} `json:"TimeoutInMillis,omitempty" yaml:"TimeoutInMillis,omitempty"`

	// Type is the type of backend that your method calls.
	// Valid values: AWS, AWS_PROXY, HTTP, HTTP_PROXY, MOCK
	Type interface{} `json:"Type,omitempty" yaml:"Type,omitempty"`

	// Uri is the Uniform Resource Identifier (URI) for the integration.
	Uri interface{} `json:"Uri,omitempty" yaml:"Uri,omitempty"`
}

// IntegrationResponse represents a response that API Gateway sends after it receives a response from the backend.
type IntegrationResponse struct {
	// ContentHandling specifies how to handle response payload content type conversions.
	// Valid values: CONVERT_TO_BINARY, CONVERT_TO_TEXT
	ContentHandling interface{} `json:"ContentHandling,omitempty" yaml:"ContentHandling,omitempty"`

	// ResponseParameters specifies response parameters that API Gateway sends to the client.
	ResponseParameters map[string]interface{} `json:"ResponseParameters,omitempty" yaml:"ResponseParameters,omitempty"`

	// ResponseTemplates specifies response payload mapping templates.
	ResponseTemplates map[string]interface{} `json:"ResponseTemplates,omitempty" yaml:"ResponseTemplates,omitempty"`

	// SelectionPattern is a regular expression that specifies which responses this response applies to.
	SelectionPattern interface{} `json:"SelectionPattern,omitempty" yaml:"SelectionPattern,omitempty"`

	// StatusCode is the status code that API Gateway uses to map the response to a template.
	StatusCode interface{} `json:"StatusCode" yaml:"StatusCode"`
}

// MethodResponse represents a response that API Gateway sends to the client.
type MethodResponse struct {
	// ResponseModels specifies the resources used for the response's content type.
	ResponseModels map[string]interface{} `json:"ResponseModels,omitempty" yaml:"ResponseModels,omitempty"`

	// ResponseParameters specifies response parameters that API Gateway can send to the client.
	ResponseParameters map[string]interface{} `json:"ResponseParameters,omitempty" yaml:"ResponseParameters,omitempty"`

	// StatusCode is the HTTP status code that this response corresponds to.
	StatusCode interface{} `json:"StatusCode" yaml:"StatusCode"`
}
