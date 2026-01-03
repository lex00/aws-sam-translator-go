// Package sns provides CloudFormation resource models for Amazon SNS.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-sns-subscription.html
package sns

// Subscription represents an AWS::SNS::Subscription resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-sns-subscription.html
type Subscription struct {
	// TopicArn is the ARN of the topic to subscribe to.
	TopicArn interface{} `json:"TopicArn" yaml:"TopicArn"`

	// Protocol is the subscription's protocol.
	// Valid values: http | https | email | email-json | sms | sqs | application | lambda | firehose
	Protocol string `json:"Protocol" yaml:"Protocol"`

	// Endpoint is the endpoint that receives notifications.
	Endpoint interface{} `json:"Endpoint,omitempty" yaml:"Endpoint,omitempty"`

	// DeliveryPolicy is the delivery policy JSON assigned to the subscription.
	DeliveryPolicy interface{} `json:"DeliveryPolicy,omitempty" yaml:"DeliveryPolicy,omitempty"`

	// FilterPolicy is the filter policy JSON assigned to the subscription.
	FilterPolicy interface{} `json:"FilterPolicy,omitempty" yaml:"FilterPolicy,omitempty"`

	// FilterPolicyScope determines where the filter policy is applied.
	// Valid values: MessageAttributes | MessageBody
	FilterPolicyScope string `json:"FilterPolicyScope,omitempty" yaml:"FilterPolicyScope,omitempty"`

	// RawMessageDelivery indicates whether raw message delivery is enabled.
	RawMessageDelivery interface{} `json:"RawMessageDelivery,omitempty" yaml:"RawMessageDelivery,omitempty"`

	// RedrivePolicy is the dead-letter queue redrive policy.
	RedrivePolicy interface{} `json:"RedrivePolicy,omitempty" yaml:"RedrivePolicy,omitempty"`

	// Region is the region for cross-region subscriptions.
	Region interface{} `json:"Region,omitempty" yaml:"Region,omitempty"`

	// SubscriptionRoleArn is the ARN of the IAM role for Firehose subscriptions.
	SubscriptionRoleArn interface{} `json:"SubscriptionRoleArn,omitempty" yaml:"SubscriptionRoleArn,omitempty"`

	// ReplayPolicy is the replay policy for the subscription (for FIFO topics).
	ReplayPolicy interface{} `json:"ReplayPolicy,omitempty" yaml:"ReplayPolicy,omitempty"`
}
