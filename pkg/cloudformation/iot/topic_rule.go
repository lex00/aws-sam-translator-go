// Package iot provides CloudFormation resource models for AWS IoT.
package iot

// ResourceTypeTopicRule is the CloudFormation resource type for AWS::IoT::TopicRule.
const ResourceTypeTopicRule = "AWS::IoT::TopicRule"

// TopicRule represents an AWS::IoT::TopicRule CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot-topicrule.html
type TopicRule struct {
	// RuleName is the name of the rule (optional).
	// If not specified, AWS CloudFormation generates a unique ID.
	RuleName interface{} `json:"RuleName,omitempty" yaml:"RuleName,omitempty"`

	// TopicRulePayload is the rule payload (required).
	TopicRulePayload *TopicRulePayload `json:"TopicRulePayload" yaml:"TopicRulePayload"`

	// Tags are the tags to associate with the topic rule (optional).
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// TopicRulePayload describes the payload of an IoT topic rule.
type TopicRulePayload struct {
	// Actions is the list of actions to perform when the rule is triggered (required).
	Actions []Action `json:"Actions" yaml:"Actions"`

	// AwsIotSqlVersion is the version of the SQL rules engine to use (optional).
	// If not specified, the latest version is used.
	AwsIotSqlVersion interface{} `json:"AwsIotSqlVersion,omitempty" yaml:"AwsIotSqlVersion,omitempty"`

	// Description is a description of the rule (optional).
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// ErrorAction is the action to take when an error occurs (optional).
	ErrorAction *Action `json:"ErrorAction,omitempty" yaml:"ErrorAction,omitempty"`

	// RuleDisabled indicates whether the rule is disabled (optional).
	// Defaults to false.
	RuleDisabled interface{} `json:"RuleDisabled,omitempty" yaml:"RuleDisabled,omitempty"`

	// Sql is the SQL statement used to query the topic (required).
	Sql interface{} `json:"Sql" yaml:"Sql"`
}

// Action represents an action that can be performed by an IoT rule.
type Action struct {
	// Lambda specifies a Lambda action.
	Lambda *LambdaAction `json:"Lambda,omitempty" yaml:"Lambda,omitempty"`
}

// LambdaAction describes an action that invokes a Lambda function.
type LambdaAction struct {
	// FunctionArn is the ARN of the Lambda function (required).
	FunctionArn interface{} `json:"FunctionArn" yaml:"FunctionArn"`
}

// Tag represents a tag key-value pair.
type Tag struct {
	Key   string `json:"Key" yaml:"Key"`
	Value string `json:"Value" yaml:"Value"`
}

// NewTopicRule creates a new TopicRule with required parameters.
func NewTopicRule(sql interface{}, functionArn interface{}) *TopicRule {
	return &TopicRule{
		TopicRulePayload: &TopicRulePayload{
			Sql: sql,
			Actions: []Action{
				{
					Lambda: &LambdaAction{
						FunctionArn: functionArn,
					},
				},
			},
		},
	}
}

// WithRuleName sets the rule name.
func (t *TopicRule) WithRuleName(name interface{}) *TopicRule {
	t.RuleName = name
	return t
}

// WithAwsIotSqlVersion sets the SQL version.
func (t *TopicRule) WithAwsIotSqlVersion(version interface{}) *TopicRule {
	t.TopicRulePayload.AwsIotSqlVersion = version
	return t
}

// WithDescription sets the rule description.
func (t *TopicRule) WithDescription(description interface{}) *TopicRule {
	t.TopicRulePayload.Description = description
	return t
}

// ToCloudFormation converts the TopicRule to a CloudFormation resource.
func (t *TopicRule) ToCloudFormation() map[string]interface{} {
	actions := make([]interface{}, len(t.TopicRulePayload.Actions))
	for i, action := range t.TopicRulePayload.Actions {
		actionMap := make(map[string]interface{})
		if action.Lambda != nil {
			actionMap["Lambda"] = map[string]interface{}{
				"FunctionArn": action.Lambda.FunctionArn,
			}
		}
		actions[i] = actionMap
	}

	payload := map[string]interface{}{
		"Sql":     t.TopicRulePayload.Sql,
		"Actions": actions,
	}

	if t.TopicRulePayload.AwsIotSqlVersion != nil {
		payload["AwsIotSqlVersion"] = t.TopicRulePayload.AwsIotSqlVersion
	}
	if t.TopicRulePayload.Description != nil {
		payload["Description"] = t.TopicRulePayload.Description
	}
	if t.TopicRulePayload.RuleDisabled != nil {
		payload["RuleDisabled"] = t.TopicRulePayload.RuleDisabled
	}
	if t.TopicRulePayload.ErrorAction != nil {
		errorAction := make(map[string]interface{})
		if t.TopicRulePayload.ErrorAction.Lambda != nil {
			errorAction["Lambda"] = map[string]interface{}{
				"FunctionArn": t.TopicRulePayload.ErrorAction.Lambda.FunctionArn,
			}
		}
		payload["ErrorAction"] = errorAction
	}

	properties := map[string]interface{}{
		"TopicRulePayload": payload,
	}

	if t.RuleName != nil {
		properties["RuleName"] = t.RuleName
	}

	if len(t.Tags) > 0 {
		tags := make([]interface{}, len(t.Tags))
		for i, tag := range t.Tags {
			tags[i] = map[string]interface{}{
				"Key":   tag.Key,
				"Value": tag.Value,
			}
		}
		properties["Tags"] = tags
	}

	return map[string]interface{}{
		"Type":       ResourceTypeTopicRule,
		"Properties": properties,
	}
}
