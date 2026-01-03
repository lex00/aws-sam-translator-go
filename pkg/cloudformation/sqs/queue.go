// Package sqs provides CloudFormation resource models for Amazon SQS.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-sqs-queue.html
package sqs

// Queue represents an AWS::SQS::Queue resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-sqs-queue.html
type Queue struct {
	// QueueName is the name of the queue.
	QueueName interface{} `json:"QueueName,omitempty" yaml:"QueueName,omitempty"`

	// ContentBasedDeduplication enables content-based deduplication for FIFO queues.
	ContentBasedDeduplication interface{} `json:"ContentBasedDeduplication,omitempty" yaml:"ContentBasedDeduplication,omitempty"`

	// DeduplicationScope specifies the scope of deduplication for FIFO queues.
	// Valid values: messageGroup | queue
	DeduplicationScope string `json:"DeduplicationScope,omitempty" yaml:"DeduplicationScope,omitempty"`

	// DelaySeconds is the time in seconds that delivery of messages is delayed.
	DelaySeconds interface{} `json:"DelaySeconds,omitempty" yaml:"DelaySeconds,omitempty"`

	// FifoQueue indicates whether this is a FIFO queue.
	FifoQueue interface{} `json:"FifoQueue,omitempty" yaml:"FifoQueue,omitempty"`

	// FifoThroughputLimit specifies the throughput limit for FIFO queues.
	// Valid values: perQueue | perMessageGroupId
	FifoThroughputLimit string `json:"FifoThroughputLimit,omitempty" yaml:"FifoThroughputLimit,omitempty"`

	// KmsMasterKeyId is the ID of an AWS KMS key for server-side encryption.
	KmsMasterKeyId interface{} `json:"KmsMasterKeyId,omitempty" yaml:"KmsMasterKeyId,omitempty"`

	// KmsDataKeyReusePeriodSeconds is the length of time for which the KMS key is reused.
	KmsDataKeyReusePeriodSeconds interface{} `json:"KmsDataKeyReusePeriodSeconds,omitempty" yaml:"KmsDataKeyReusePeriodSeconds,omitempty"`

	// MaximumMessageSize is the maximum size of a message in bytes (1024-262144).
	MaximumMessageSize interface{} `json:"MaximumMessageSize,omitempty" yaml:"MaximumMessageSize,omitempty"`

	// MessageRetentionPeriod is the number of seconds to retain messages (60-1209600).
	MessageRetentionPeriod interface{} `json:"MessageRetentionPeriod,omitempty" yaml:"MessageRetentionPeriod,omitempty"`

	// ReceiveMessageWaitTimeSeconds is the duration for long polling (0-20 seconds).
	ReceiveMessageWaitTimeSeconds interface{} `json:"ReceiveMessageWaitTimeSeconds,omitempty" yaml:"ReceiveMessageWaitTimeSeconds,omitempty"`

	// RedriveAllowPolicy specifies which source queues can use this queue as a DLQ.
	RedriveAllowPolicy interface{} `json:"RedriveAllowPolicy,omitempty" yaml:"RedriveAllowPolicy,omitempty"`

	// RedrivePolicy specifies the dead-letter queue configuration.
	RedrivePolicy interface{} `json:"RedrivePolicy,omitempty" yaml:"RedrivePolicy,omitempty"`

	// SqsManagedSseEnabled enables SQS managed server-side encryption.
	SqsManagedSseEnabled interface{} `json:"SqsManagedSseEnabled,omitempty" yaml:"SqsManagedSseEnabled,omitempty"`

	// Tags is a list of key-value pairs to apply to the queue.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// VisibilityTimeout is the visibility timeout in seconds (0-43200).
	VisibilityTimeout interface{} `json:"VisibilityTimeout,omitempty" yaml:"VisibilityTimeout,omitempty"`
}

// Tag represents a key-value pair tag.
type Tag struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}

// RedrivePolicy represents the dead-letter queue configuration.
// This is used when creating the RedrivePolicy as a structured object.
type RedrivePolicy struct {
	// DeadLetterTargetArn is the ARN of the dead-letter queue.
	DeadLetterTargetArn interface{} `json:"deadLetterTargetArn" yaml:"deadLetterTargetArn"`

	// MaxReceiveCount is the number of times a message can be received before being sent to the DLQ.
	MaxReceiveCount interface{} `json:"maxReceiveCount" yaml:"maxReceiveCount"`
}

// RedriveAllowPolicy represents which source queues can use this queue as a DLQ.
// This is used when creating the RedriveAllowPolicy as a structured object.
type RedriveAllowPolicy struct {
	// RedrivePermission specifies who can use this queue as a DLQ.
	// Valid values: allowAll | denyAll | byQueue
	RedrivePermission string `json:"redrivePermission" yaml:"redrivePermission"`

	// SourceQueueArns is the list of ARNs of source queues (when redrivePermission is byQueue).
	SourceQueueArns []interface{} `json:"sourceQueueArns,omitempty" yaml:"sourceQueueArns,omitempty"`
}
