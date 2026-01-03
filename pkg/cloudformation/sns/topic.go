// Package sns provides CloudFormation resource models for Amazon SNS.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-sns-topic.html
package sns

// Topic represents an AWS::SNS::Topic resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-sns-topic.html
type Topic struct {
	// TopicName is the name of the topic.
	TopicName interface{} `json:"TopicName,omitempty" yaml:"TopicName,omitempty"`

	// DisplayName is the display name for the topic.
	DisplayName interface{} `json:"DisplayName,omitempty" yaml:"DisplayName,omitempty"`

	// FifoTopic indicates whether this is a FIFO topic.
	FifoTopic interface{} `json:"FifoTopic,omitempty" yaml:"FifoTopic,omitempty"`

	// ContentBasedDeduplication enables content-based deduplication for FIFO topics.
	ContentBasedDeduplication interface{} `json:"ContentBasedDeduplication,omitempty" yaml:"ContentBasedDeduplication,omitempty"`

	// KmsMasterKeyId is the ID of an AWS KMS key for server-side encryption.
	KmsMasterKeyId interface{} `json:"KmsMasterKeyId,omitempty" yaml:"KmsMasterKeyId,omitempty"`

	// DataProtectionPolicy is the body of the data protection policy for the topic.
	DataProtectionPolicy interface{} `json:"DataProtectionPolicy,omitempty" yaml:"DataProtectionPolicy,omitempty"`

	// Subscription is a list of subscriptions to create when the topic is created.
	Subscription []TopicSubscription `json:"Subscription,omitempty" yaml:"Subscription,omitempty"`

	// Tags is a list of key-value pairs to apply to the topic.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// TracingConfig specifies tracing mode for the topic.
	TracingConfig interface{} `json:"TracingConfig,omitempty" yaml:"TracingConfig,omitempty"`

	// SignatureVersion is the signature version for the topic. Valid values: 1 | 2
	SignatureVersion interface{} `json:"SignatureVersion,omitempty" yaml:"SignatureVersion,omitempty"`

	// ArchivePolicy is the archive policy for the topic (for FIFO topics).
	ArchivePolicy interface{} `json:"ArchivePolicy,omitempty" yaml:"ArchivePolicy,omitempty"`

	// DeliveryStatusLogging specifies delivery status logging configuration.
	DeliveryStatusLogging []LoggingConfig `json:"DeliveryStatusLogging,omitempty" yaml:"DeliveryStatusLogging,omitempty"`
}

// TopicSubscription represents a subscription embedded in a topic.
type TopicSubscription struct {
	// Endpoint is the endpoint that receives notifications.
	Endpoint interface{} `json:"Endpoint" yaml:"Endpoint"`

	// Protocol is the subscription's protocol.
	Protocol string `json:"Protocol" yaml:"Protocol"`
}

// Tag represents a key-value pair tag.
type Tag struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}

// LoggingConfig specifies delivery status logging configuration.
type LoggingConfig struct {
	// Protocol is the protocol for the endpoint. Valid values: http/s | sqs | lambda | firehose | application
	Protocol string `json:"Protocol" yaml:"Protocol"`

	// SuccessFeedbackRoleArn is the IAM role ARN for success feedback logging.
	SuccessFeedbackRoleArn interface{} `json:"SuccessFeedbackRoleArn,omitempty" yaml:"SuccessFeedbackRoleArn,omitempty"`

	// SuccessFeedbackSampleRate is the percentage of successful deliveries to log.
	SuccessFeedbackSampleRate interface{} `json:"SuccessFeedbackSampleRate,omitempty" yaml:"SuccessFeedbackSampleRate,omitempty"`

	// FailureFeedbackRoleArn is the IAM role ARN for failure feedback logging.
	FailureFeedbackRoleArn interface{} `json:"FailureFeedbackRoleArn,omitempty" yaml:"FailureFeedbackRoleArn,omitempty"`
}
