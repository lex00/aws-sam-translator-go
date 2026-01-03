// Package logs provides CloudFormation resource models for Amazon CloudWatch Logs.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-logs-loggroup.html
package logs

// LogGroup represents an AWS::Logs::LogGroup resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-logs-loggroup.html
type LogGroup struct {
	// LogGroupName is the name of the log group.
	LogGroupName interface{} `json:"LogGroupName,omitempty" yaml:"LogGroupName,omitempty"`

	// RetentionInDays is the number of days to retain log events.
	// Valid values: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653
	RetentionInDays interface{} `json:"RetentionInDays,omitempty" yaml:"RetentionInDays,omitempty"`

	// KmsKeyId is the ARN of the AWS KMS key for encryption.
	KmsKeyId interface{} `json:"KmsKeyId,omitempty" yaml:"KmsKeyId,omitempty"`

	// DataProtectionPolicy is the data protection policy for the log group.
	DataProtectionPolicy interface{} `json:"DataProtectionPolicy,omitempty" yaml:"DataProtectionPolicy,omitempty"`

	// LogGroupClass specifies the log class. Valid values: STANDARD | INFREQUENT_ACCESS
	LogGroupClass string `json:"LogGroupClass,omitempty" yaml:"LogGroupClass,omitempty"`

	// Tags is a list of key-value pairs to apply to the log group.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// Tag represents a key-value pair tag.
type Tag struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}
