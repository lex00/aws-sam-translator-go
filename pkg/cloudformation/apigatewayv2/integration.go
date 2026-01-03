// Package apigatewayv2 provides CloudFormation resource models for AWS API Gateway V2 (HTTP API/WebSocket API).
package apigatewayv2

// Integration represents an AWS::ApiGatewayV2::Integration CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integration.html
type Integration struct {
	// ApiId is the API identifier.
	ApiId interface{} `json:"ApiId" yaml:"ApiId"`

	// ConnectionId is the ID of the VPC link for a private integration.
	ConnectionId interface{} `json:"ConnectionId,omitempty" yaml:"ConnectionId,omitempty"`

	// ConnectionType is the type of the network connection to the integration endpoint.
	// Valid values: INTERNET, VPC_LINK
	ConnectionType interface{} `json:"ConnectionType,omitempty" yaml:"ConnectionType,omitempty"`

	// ContentHandlingStrategy specifies how to handle response payload content type conversions.
	// Valid values: CONVERT_TO_BINARY, CONVERT_TO_TEXT
	ContentHandlingStrategy interface{} `json:"ContentHandlingStrategy,omitempty" yaml:"ContentHandlingStrategy,omitempty"`

	// CredentialsArn specifies the credentials required for the integration.
	CredentialsArn interface{} `json:"CredentialsArn,omitempty" yaml:"CredentialsArn,omitempty"`

	// Description is the description of the integration.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// IntegrationMethod specifies the integration's HTTP method type.
	IntegrationMethod interface{} `json:"IntegrationMethod,omitempty" yaml:"IntegrationMethod,omitempty"`

	// IntegrationSubtype specifies the AWS service action to invoke.
	IntegrationSubtype interface{} `json:"IntegrationSubtype,omitempty" yaml:"IntegrationSubtype,omitempty"`

	// IntegrationType is the integration type.
	// Valid values: AWS, AWS_PROXY, HTTP, HTTP_PROXY, MOCK
	IntegrationType interface{} `json:"IntegrationType" yaml:"IntegrationType"`

	// IntegrationUri is the URI of the Lambda function for a Lambda proxy integration.
	IntegrationUri interface{} `json:"IntegrationUri,omitempty" yaml:"IntegrationUri,omitempty"`

	// PassthroughBehavior specifies the pass-through behavior for incoming requests.
	// Valid values: WHEN_NO_MATCH, NEVER, WHEN_NO_TEMPLATES
	PassthroughBehavior interface{} `json:"PassthroughBehavior,omitempty" yaml:"PassthroughBehavior,omitempty"`

	// PayloadFormatVersion specifies the format of the payload sent to an integration.
	// Valid values: 1.0, 2.0
	PayloadFormatVersion interface{} `json:"PayloadFormatVersion,omitempty" yaml:"PayloadFormatVersion,omitempty"`

	// RequestParameters is a key-value map specifying request parameters that are passed from the method request.
	RequestParameters map[string]interface{} `json:"RequestParameters,omitempty" yaml:"RequestParameters,omitempty"`

	// RequestTemplates is a map of Velocity templates that are applied on the request payload.
	RequestTemplates map[string]interface{} `json:"RequestTemplates,omitempty" yaml:"RequestTemplates,omitempty"`

	// ResponseParameters specifies response parameters.
	ResponseParameters map[string]interface{} `json:"ResponseParameters,omitempty" yaml:"ResponseParameters,omitempty"`

	// TemplateSelectionExpression is the template selection expression for the integration.
	TemplateSelectionExpression interface{} `json:"TemplateSelectionExpression,omitempty" yaml:"TemplateSelectionExpression,omitempty"`

	// TimeoutInMillis is the custom timeout in milliseconds.
	TimeoutInMillis interface{} `json:"TimeoutInMillis,omitempty" yaml:"TimeoutInMillis,omitempty"`

	// TlsConfig is the TLS configuration for a private integration.
	TlsConfig *TlsConfig `json:"TlsConfig,omitempty" yaml:"TlsConfig,omitempty"`
}

// TlsConfig represents TLS configuration for a private integration.
type TlsConfig struct {
	// ServerNameToVerify specifies the server name to verify.
	ServerNameToVerify interface{} `json:"ServerNameToVerify,omitempty" yaml:"ServerNameToVerify,omitempty"`
}
