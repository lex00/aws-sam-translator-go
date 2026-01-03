// Package stepfunctions provides CloudFormation resource models for AWS Step Functions.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-stepfunctions-statemachine.html
package stepfunctions

// StateMachine represents an AWS::StepFunctions::StateMachine resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-stepfunctions-statemachine.html
type StateMachine struct {
	// StateMachineName is the name of the state machine.
	StateMachineName interface{} `json:"StateMachineName,omitempty" yaml:"StateMachineName,omitempty"`

	// Definition is the Amazon States Language definition of the state machine.
	Definition interface{} `json:"Definition,omitempty" yaml:"Definition,omitempty"`

	// DefinitionString is the Amazon States Language definition as a string.
	DefinitionString interface{} `json:"DefinitionString,omitempty" yaml:"DefinitionString,omitempty"`

	// DefinitionS3Location specifies the S3 location of the state machine definition.
	DefinitionS3Location *S3Location `json:"DefinitionS3Location,omitempty" yaml:"DefinitionS3Location,omitempty"`

	// DefinitionSubstitutions is a map of key-value pairs for definition substitutions.
	DefinitionSubstitutions map[string]interface{} `json:"DefinitionSubstitutions,omitempty" yaml:"DefinitionSubstitutions,omitempty"`

	// RoleArn is the ARN of the IAM role used for executions.
	RoleArn interface{} `json:"RoleArn,omitempty" yaml:"RoleArn,omitempty"`

	// StateMachineType is the type of state machine. Valid values: STANDARD | EXPRESS
	StateMachineType string `json:"StateMachineType,omitempty" yaml:"StateMachineType,omitempty"`

	// LoggingConfiguration specifies the logging configuration.
	LoggingConfiguration *LoggingConfiguration `json:"LoggingConfiguration,omitempty" yaml:"LoggingConfiguration,omitempty"`

	// TracingConfiguration specifies the tracing configuration for X-Ray.
	TracingConfiguration *TracingConfiguration `json:"TracingConfiguration,omitempty" yaml:"TracingConfiguration,omitempty"`

	// Tags is a list of tags to apply to the state machine.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// EncryptionConfiguration specifies encryption settings.
	EncryptionConfiguration *EncryptionConfiguration `json:"EncryptionConfiguration,omitempty" yaml:"EncryptionConfiguration,omitempty"`
}

// S3Location specifies an S3 location for the definition.
type S3Location struct {
	// Bucket is the name of the S3 bucket.
	Bucket interface{} `json:"Bucket" yaml:"Bucket"`

	// Key is the name of the file in the S3 bucket.
	Key interface{} `json:"Key" yaml:"Key"`

	// Version is the version of the file.
	Version interface{} `json:"Version,omitempty" yaml:"Version,omitempty"`
}

// LoggingConfiguration specifies the logging configuration for a state machine.
type LoggingConfiguration struct {
	// Level is the logging level. Valid values: ALL | ERROR | FATAL | OFF
	Level string `json:"Level,omitempty" yaml:"Level,omitempty"`

	// IncludeExecutionData indicates whether to include execution data.
	IncludeExecutionData interface{} `json:"IncludeExecutionData,omitempty" yaml:"IncludeExecutionData,omitempty"`

	// Destinations is a list of logging destinations.
	Destinations []LogDestination `json:"Destinations,omitempty" yaml:"Destinations,omitempty"`
}

// LogDestination specifies a logging destination.
type LogDestination struct {
	// CloudWatchLogsLogGroup specifies the CloudWatch Logs log group.
	CloudWatchLogsLogGroup *CloudWatchLogsLogGroup `json:"CloudWatchLogsLogGroup,omitempty" yaml:"CloudWatchLogsLogGroup,omitempty"`
}

// CloudWatchLogsLogGroup specifies a CloudWatch Logs log group.
type CloudWatchLogsLogGroup struct {
	// LogGroupArn is the ARN of the CloudWatch Logs log group.
	LogGroupArn interface{} `json:"LogGroupArn" yaml:"LogGroupArn"`
}

// TracingConfiguration specifies X-Ray tracing configuration.
type TracingConfiguration struct {
	// Enabled indicates whether X-Ray tracing is enabled.
	Enabled interface{} `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`
}

// Tag represents a key-value pair tag.
type Tag struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}

// EncryptionConfiguration specifies encryption settings for the state machine.
type EncryptionConfiguration struct {
	// KmsKeyId is the ARN of the KMS key.
	KmsKeyId interface{} `json:"KmsKeyId,omitempty" yaml:"KmsKeyId,omitempty"`

	// KmsDataKeyReusePeriodSeconds is the maximum duration for reusing data keys.
	KmsDataKeyReusePeriodSeconds interface{} `json:"KmsDataKeyReusePeriodSeconds,omitempty" yaml:"KmsDataKeyReusePeriodSeconds,omitempty"`

	// Type is the encryption type. Valid values: AWS_OWNED_KEY | CUSTOMER_MANAGED_KMS_KEY
	Type string `json:"Type" yaml:"Type"`
}
