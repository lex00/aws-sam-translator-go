// Package cognito provides CloudFormation resource models for Amazon Cognito.
package cognito

// ResourceTypeUserPool is the CloudFormation resource type for AWS::Cognito::UserPool.
const ResourceTypeUserPool = "AWS::Cognito::UserPool"

// LambdaConfig represents the Lambda trigger configuration for a Cognito User Pool.
// This maps to the LambdaConfig property of AWS::Cognito::UserPool.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-lambdaconfig.html
type LambdaConfig struct {
	// CreateAuthChallenge is the ARN of the Lambda function to invoke for custom challenge creation.
	CreateAuthChallenge interface{} `json:"CreateAuthChallenge,omitempty" yaml:"CreateAuthChallenge,omitempty"`

	// CustomEmailSender configures a custom email sender Lambda.
	CustomEmailSender *CustomEmailSender `json:"CustomEmailSender,omitempty" yaml:"CustomEmailSender,omitempty"`

	// CustomMessage is the ARN of the Lambda function for custom messages.
	CustomMessage interface{} `json:"CustomMessage,omitempty" yaml:"CustomMessage,omitempty"`

	// CustomSMSSender configures a custom SMS sender Lambda.
	CustomSMSSender *CustomSMSSender `json:"CustomSMSSender,omitempty" yaml:"CustomSMSSender,omitempty"`

	// DefineAuthChallenge is the ARN of the Lambda function to define auth challenges.
	DefineAuthChallenge interface{} `json:"DefineAuthChallenge,omitempty" yaml:"DefineAuthChallenge,omitempty"`

	// KMSKeyID is the KMS key ID for custom sender Lambdas.
	KMSKeyID interface{} `json:"KMSKeyID,omitempty" yaml:"KMSKeyID,omitempty"`

	// PostAuthentication is the ARN of the Lambda function for post-authentication.
	PostAuthentication interface{} `json:"PostAuthentication,omitempty" yaml:"PostAuthentication,omitempty"`

	// PostConfirmation is the ARN of the Lambda function for post-confirmation.
	PostConfirmation interface{} `json:"PostConfirmation,omitempty" yaml:"PostConfirmation,omitempty"`

	// PreAuthentication is the ARN of the Lambda function for pre-authentication.
	PreAuthentication interface{} `json:"PreAuthentication,omitempty" yaml:"PreAuthentication,omitempty"`

	// PreSignUp is the ARN of the Lambda function for pre-sign-up.
	PreSignUp interface{} `json:"PreSignUp,omitempty" yaml:"PreSignUp,omitempty"`

	// PreTokenGeneration is the ARN of the Lambda function for pre-token generation.
	PreTokenGeneration interface{} `json:"PreTokenGeneration,omitempty" yaml:"PreTokenGeneration,omitempty"`

	// PreTokenGenerationConfig configures advanced pre-token generation.
	PreTokenGenerationConfig *PreTokenGenerationConfig `json:"PreTokenGenerationConfig,omitempty" yaml:"PreTokenGenerationConfig,omitempty"`

	// UserMigration is the ARN of the Lambda function for user migration.
	UserMigration interface{} `json:"UserMigration,omitempty" yaml:"UserMigration,omitempty"`

	// VerifyAuthChallengeResponse is the ARN of the Lambda function to verify auth challenge responses.
	VerifyAuthChallengeResponse interface{} `json:"VerifyAuthChallengeResponse,omitempty" yaml:"VerifyAuthChallengeResponse,omitempty"`
}

// CustomEmailSender configures a custom email sender Lambda function.
type CustomEmailSender struct {
	LambdaArn     interface{} `json:"LambdaArn" yaml:"LambdaArn"`
	LambdaVersion interface{} `json:"LambdaVersion" yaml:"LambdaVersion"`
}

// CustomSMSSender configures a custom SMS sender Lambda function.
type CustomSMSSender struct {
	LambdaArn     interface{} `json:"LambdaArn" yaml:"LambdaArn"`
	LambdaVersion interface{} `json:"LambdaVersion" yaml:"LambdaVersion"`
}

// PreTokenGenerationConfig configures advanced pre-token generation.
type PreTokenGenerationConfig struct {
	LambdaArn     interface{} `json:"LambdaArn" yaml:"LambdaArn"`
	LambdaVersion interface{} `json:"LambdaVersion" yaml:"LambdaVersion"`
}

// CognitoTrigger represents the valid Cognito trigger types.
type CognitoTrigger string

const (
	TriggerPreSignUp                    CognitoTrigger = "PreSignUp"
	TriggerPostConfirmation             CognitoTrigger = "PostConfirmation"
	TriggerPreAuthentication            CognitoTrigger = "PreAuthentication"
	TriggerPostAuthentication           CognitoTrigger = "PostAuthentication"
	TriggerPreTokenGeneration           CognitoTrigger = "PreTokenGeneration"
	TriggerCustomMessage                CognitoTrigger = "CustomMessage"
	TriggerUserMigration                CognitoTrigger = "UserMigration"
	TriggerDefineAuthChallenge          CognitoTrigger = "DefineAuthChallenge"
	TriggerCreateAuthChallenge          CognitoTrigger = "CreateAuthChallenge"
	TriggerVerifyAuthChallengeResponse  CognitoTrigger = "VerifyAuthChallengeResponse"
)

// GetLambdaConfigProperty returns the LambdaConfig property name for a trigger.
func GetLambdaConfigProperty(trigger CognitoTrigger) string {
	return string(trigger)
}
