// Package push provides push event source handlers for AWS SAM.
package push

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/sns"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// SNSEvent represents a SAM SNS event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-sns.html
type SNSEvent struct {
	// Topic is the ARN of the SNS topic to subscribe to (required).
	// Can be a string ARN, a Ref intrinsic, or other CloudFormation reference.
	Topic interface{} `json:"Topic" yaml:"Topic"`

	// Region is the region of the SNS topic (for cross-region subscriptions).
	Region interface{} `json:"Region,omitempty" yaml:"Region,omitempty"`

	// FilterPolicy is the SNS subscription filter policy for message filtering.
	// See: https://docs.aws.amazon.com/sns/latest/dg/sns-subscription-filter-policies.html
	FilterPolicy interface{} `json:"FilterPolicy,omitempty" yaml:"FilterPolicy,omitempty"`

	// FilterPolicyScope determines where the filter policy is applied.
	// Valid values: MessageAttributes | MessageBody
	FilterPolicyScope string `json:"FilterPolicyScope,omitempty" yaml:"FilterPolicyScope,omitempty"`

	// SqsSubscription indicates whether to create an SQS queue between SNS and Lambda.
	SqsSubscription bool `json:"SqsSubscription,omitempty" yaml:"SqsSubscription,omitempty"`

	// RedrivePolicy is the dead-letter queue redrive policy.
	RedrivePolicy interface{} `json:"RedrivePolicy,omitempty" yaml:"RedrivePolicy,omitempty"`
}

// SNSEventSourceHandler handles SNS event sources.
type SNSEventSourceHandler struct{}

// NewSNSEventSourceHandler creates a new SNS event source handler.
func NewSNSEventSourceHandler() *SNSEventSourceHandler {
	return &SNSEventSourceHandler{}
}

// GenerateResources generates CloudFormation resources for an SNS event source.
// It creates:
// 1. AWS::SNS::Subscription - subscribes the Lambda function to the SNS topic
// 2. AWS::Lambda::Permission - grants SNS permission to invoke the Lambda function
func (h *SNSEventSourceHandler) GenerateResources(
	functionLogicalID string,
	eventLogicalID string,
	event *SNSEvent,
) (map[string]interface{}, error) {
	if event.Topic == nil {
		return nil, fmt.Errorf("SNS event source requires a Topic property")
	}

	resources := make(map[string]interface{})

	// Generate logical IDs for the resources
	subscriptionLogicalID := fmt.Sprintf("%s%s", functionLogicalID, eventLogicalID)
	permissionLogicalID := fmt.Sprintf("%s%sPermission", functionLogicalID, eventLogicalID)

	// Create SNS Subscription
	subscription := &sns.Subscription{
		TopicArn: event.Topic,
		Protocol: "lambda",
		Endpoint: map[string]interface{}{
			"Fn::GetAtt": []interface{}{functionLogicalID, "Arn"},
		},
	}

	// Add optional properties
	if event.Region != nil {
		subscription.Region = event.Region
	}
	if event.FilterPolicy != nil {
		subscription.FilterPolicy = event.FilterPolicy
	}
	if event.FilterPolicyScope != "" {
		subscription.FilterPolicyScope = event.FilterPolicyScope
	}
	if event.RedrivePolicy != nil {
		subscription.RedrivePolicy = event.RedrivePolicy
	}

	// Convert subscription to CloudFormation format
	resources[subscriptionLogicalID] = h.subscriptionToCloudFormation(subscription)

	// Create Lambda Permission
	permission := lambda.NewSNSPermission(
		map[string]interface{}{"Ref": functionLogicalID},
		event.Topic,
	)

	// Convert permission to CloudFormation format
	resources[permissionLogicalID] = permission.ToCloudFormation()

	return resources, nil
}

// subscriptionToCloudFormation converts an SNS Subscription to CloudFormation format.
func (h *SNSEventSourceHandler) subscriptionToCloudFormation(sub *sns.Subscription) map[string]interface{} {
	properties := make(map[string]interface{})

	properties["TopicArn"] = sub.TopicArn
	properties["Protocol"] = sub.Protocol
	properties["Endpoint"] = sub.Endpoint

	if sub.Region != nil {
		properties["Region"] = sub.Region
	}
	if sub.FilterPolicy != nil {
		properties["FilterPolicy"] = sub.FilterPolicy
	}
	if sub.FilterPolicyScope != "" {
		properties["FilterPolicyScope"] = sub.FilterPolicyScope
	}
	if sub.RedrivePolicy != nil {
		properties["RedrivePolicy"] = sub.RedrivePolicy
	}

	return map[string]interface{}{
		"Type":       "AWS::SNS::Subscription",
		"Properties": properties,
	}
}
