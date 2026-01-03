package lambda

// Permission represents an AWS::Lambda::Permission CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-permission.html
type Permission struct {
	// Action is the action that the principal can use on the function (required).
	// Typically "lambda:InvokeFunction".
	Action string `json:"Action" yaml:"Action"`

	// EventSourceToken is a unique token for Alexa Smart Home functions.
	EventSourceToken string `json:"EventSourceToken,omitempty" yaml:"EventSourceToken,omitempty"`

	// FunctionName is the name or ARN of the Lambda function (required).
	FunctionName interface{} `json:"FunctionName" yaml:"FunctionName"`

	// FunctionUrlAuthType specifies the auth type for Function URL invocations.
	// Valid values: AWS_IAM, NONE
	FunctionUrlAuthType string `json:"FunctionUrlAuthType,omitempty" yaml:"FunctionUrlAuthType,omitempty"`

	// Principal is the AWS service or account invoking the function (required).
	Principal string `json:"Principal" yaml:"Principal"`

	// PrincipalOrgID restricts access to accounts in the specified organization.
	PrincipalOrgID string `json:"PrincipalOrgID,omitempty" yaml:"PrincipalOrgID,omitempty"`

	// SourceAccount is the account ID for S3 bucket invocations.
	SourceAccount interface{} `json:"SourceAccount,omitempty" yaml:"SourceAccount,omitempty"`

	// SourceArn is the ARN of the invoking resource.
	SourceArn interface{} `json:"SourceArn,omitempty" yaml:"SourceArn,omitempty"`
}

// NewPermission creates a new Permission with the required properties.
func NewPermission(action string, functionName interface{}, principal string) *Permission {
	return &Permission{
		Action:       action,
		FunctionName: functionName,
		Principal:    principal,
	}
}

// NewInvokePermission creates a permission for lambda:InvokeFunction.
func NewInvokePermission(functionName interface{}, principal string) *Permission {
	return &Permission{
		Action:       "lambda:InvokeFunction",
		FunctionName: functionName,
		Principal:    principal,
	}
}

// NewAPIGatewayPermission creates a permission for API Gateway invocations.
func NewAPIGatewayPermission(functionName interface{}, sourceArn interface{}) *Permission {
	return &Permission{
		Action:       "lambda:InvokeFunction",
		FunctionName: functionName,
		Principal:    "apigateway.amazonaws.com",
		SourceArn:    sourceArn,
	}
}

// NewS3Permission creates a permission for S3 bucket invocations.
func NewS3Permission(functionName interface{}, sourceArn interface{}, sourceAccount interface{}) *Permission {
	return &Permission{
		Action:        "lambda:InvokeFunction",
		FunctionName:  functionName,
		Principal:     "s3.amazonaws.com",
		SourceArn:     sourceArn,
		SourceAccount: sourceAccount,
	}
}

// NewSNSPermission creates a permission for SNS invocations.
func NewSNSPermission(functionName interface{}, sourceArn interface{}) *Permission {
	return &Permission{
		Action:       "lambda:InvokeFunction",
		FunctionName: functionName,
		Principal:    "sns.amazonaws.com",
		SourceArn:    sourceArn,
	}
}

// NewEventsPermission creates a permission for EventBridge invocations.
func NewEventsPermission(functionName interface{}, sourceArn interface{}) *Permission {
	return &Permission{
		Action:       "lambda:InvokeFunction",
		FunctionName: functionName,
		Principal:    "events.amazonaws.com",
		SourceArn:    sourceArn,
	}
}

// NewCloudWatchLogsPermission creates a permission for CloudWatch Logs invocations.
func NewCloudWatchLogsPermission(functionName interface{}, sourceArn interface{}) *Permission {
	return &Permission{
		Action:       "lambda:InvokeFunction",
		FunctionName: functionName,
		Principal:    "logs.amazonaws.com",
		SourceArn:    sourceArn,
	}
}

// NewCognitoPermission creates a permission for Cognito User Pool invocations.
func NewCognitoPermission(functionName interface{}, sourceArn interface{}) *Permission {
	return &Permission{
		Action:       "lambda:InvokeFunction",
		FunctionName: functionName,
		Principal:    "cognito-idp.amazonaws.com",
		SourceArn:    sourceArn,
	}
}

// NewIoTPermission creates a permission for IoT invocations.
func NewIoTPermission(functionName interface{}, sourceArn interface{}) *Permission {
	return &Permission{
		Action:       "lambda:InvokeFunction",
		FunctionName: functionName,
		Principal:    "iot.amazonaws.com",
		SourceArn:    sourceArn,
	}
}

// NewAlexaPermission creates a permission for Alexa Smart Home invocations.
func NewAlexaPermission(functionName interface{}, eventSourceToken string) *Permission {
	return &Permission{
		Action:           "lambda:InvokeFunction",
		FunctionName:     functionName,
		Principal:        "alexa-appkit.amazon.com",
		EventSourceToken: eventSourceToken,
	}
}

// WithSourceArn sets the source ARN for the permission.
func (p *Permission) WithSourceArn(sourceArn interface{}) *Permission {
	p.SourceArn = sourceArn
	return p
}

// WithSourceAccount sets the source account for the permission.
func (p *Permission) WithSourceAccount(sourceAccount interface{}) *Permission {
	p.SourceAccount = sourceAccount
	return p
}

// WithEventSourceToken sets the event source token for the permission.
func (p *Permission) WithEventSourceToken(token string) *Permission {
	p.EventSourceToken = token
	return p
}

// WithPrincipalOrgID restricts access to accounts in the organization.
func (p *Permission) WithPrincipalOrgID(orgID string) *Permission {
	p.PrincipalOrgID = orgID
	return p
}

// WithFunctionUrlAuthType sets the auth type for Function URL invocations.
func (p *Permission) WithFunctionUrlAuthType(authType string) *Permission {
	p.FunctionUrlAuthType = authType
	return p
}

// ToCloudFormation converts the Permission to a CloudFormation resource.
func (p *Permission) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	properties["Action"] = p.Action
	properties["FunctionName"] = p.FunctionName
	properties["Principal"] = p.Principal

	if p.EventSourceToken != "" {
		properties["EventSourceToken"] = p.EventSourceToken
	}
	if p.FunctionUrlAuthType != "" {
		properties["FunctionUrlAuthType"] = p.FunctionUrlAuthType
	}
	if p.PrincipalOrgID != "" {
		properties["PrincipalOrgID"] = p.PrincipalOrgID
	}
	if p.SourceAccount != nil {
		properties["SourceAccount"] = p.SourceAccount
	}
	if p.SourceArn != nil {
		properties["SourceArn"] = p.SourceArn
	}

	return map[string]interface{}{
		"Type":       ResourceTypePermission,
		"Properties": properties,
	}
}
