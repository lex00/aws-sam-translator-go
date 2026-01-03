// Package push provides event source handlers for push-based Lambda triggers.
package push

import (
	"fmt"

	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/events"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// CloudWatchEvent represents a CloudWatch Events (EventBridge) event source.
// This event source type creates an AWS::Events::Rule that triggers a Lambda function
// based on an event pattern or schedule expression.
//
// SAM Event Type: CloudWatchEvent or EventBridgeRule
// CloudFormation Resources Generated:
//   - AWS::Events::Rule - The EventBridge rule with pattern/schedule and targets
//   - AWS::Lambda::Permission - Permission for EventBridge to invoke the function
//
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-cloudwatchevent.html
type CloudWatchEvent struct {
	// EventBusName is the name or ARN of the event bus to associate with this rule.
	// If not specified, the default event bus is used.
	EventBusName interface{} `json:"EventBusName,omitempty" yaml:"EventBusName,omitempty"`

	// Pattern is the event pattern that describes which events to match.
	// This is an object that follows the EventBridge event pattern format.
	// Either Pattern or Schedule must be specified, but not both.
	Pattern interface{} `json:"Pattern,omitempty" yaml:"Pattern,omitempty"`

	// Schedule is the scheduling expression (cron or rate).
	// Either Pattern or Schedule must be specified, but not both.
	// Examples: "rate(5 minutes)", "cron(0 12 * * ? *)"
	Schedule interface{} `json:"Schedule,omitempty" yaml:"Schedule,omitempty"`

	// State indicates whether the rule is enabled. Valid values: DISABLED | ENABLED
	// Default: ENABLED
	State string `json:"State,omitempty" yaml:"State,omitempty"`

	// Input is the JSON text to pass to the target function.
	Input interface{} `json:"Input,omitempty" yaml:"Input,omitempty"`

	// InputPath is the JSONPath to extract from the event to pass to the target.
	InputPath interface{} `json:"InputPath,omitempty" yaml:"InputPath,omitempty"`

	// InputTransformer specifies settings for transforming input before passing to the target.
	InputTransformer *InputTransformer `json:"InputTransformer,omitempty" yaml:"InputTransformer,omitempty"`

	// Target allows customization of the EventBridge target configuration.
	Target *Target `json:"Target,omitempty" yaml:"Target,omitempty"`

	// DeadLetterConfig specifies the dead-letter queue configuration.
	DeadLetterConfig *DeadLetterConfig `json:"DeadLetterConfig,omitempty" yaml:"DeadLetterConfig,omitempty"`

	// RetryPolicy specifies the retry policy settings.
	RetryPolicy *RetryPolicy `json:"RetryPolicy,omitempty" yaml:"RetryPolicy,omitempty"`
}

// InputTransformer specifies settings for input transformation.
type InputTransformer struct {
	// InputPathsMap is a map of JSON paths to extract from the event.
	InputPathsMap map[string]interface{} `json:"InputPathsMap,omitempty" yaml:"InputPathsMap,omitempty"`

	// InputTemplate is the template to use for the transformed input.
	InputTemplate interface{} `json:"InputTemplate" yaml:"InputTemplate"`
}

// Target allows customization of the EventBridge target configuration.
type Target struct {
	// Id is a unique identifier for the target.
	// If not specified, a default ID will be generated.
	Id interface{} `json:"Id,omitempty" yaml:"Id,omitempty"`
}

// DeadLetterConfig specifies dead-letter queue configuration.
type DeadLetterConfig struct {
	// Arn is the ARN of the SQS queue to use as the dead-letter queue.
	Arn interface{} `json:"Arn,omitempty" yaml:"Arn,omitempty"`

	// Type is the type of dead-letter config. Currently only "SQS" is supported.
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

// RetryPolicy specifies retry policy settings.
type RetryPolicy struct {
	// MaximumEventAgeInSeconds is the maximum age of an event (60-86400).
	MaximumEventAgeInSeconds interface{} `json:"MaximumEventAgeInSeconds,omitempty" yaml:"MaximumEventAgeInSeconds,omitempty"`

	// MaximumRetryAttempts is the maximum number of retry attempts (0-185).
	MaximumRetryAttempts interface{} `json:"MaximumRetryAttempts,omitempty" yaml:"MaximumRetryAttempts,omitempty"`
}

// ToCloudFormation generates CloudFormation resources for the CloudWatch Event source.
// It returns a map containing:
//   - An AWS::Events::Rule resource
//   - An AWS::Lambda::Permission resource allowing EventBridge to invoke the function
//
// The logical IDs are generated as:
//   - Rule: {FunctionLogicalId}{EventLogicalId}
//   - Permission: {FunctionLogicalId}{EventLogicalId}Permission
func (e *CloudWatchEvent) ToCloudFormation(functionLogicalID string, eventLogicalID string) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Generate logical IDs for the resources
	ruleLogicalID := functionLogicalID + eventLogicalID
	permissionLogicalID := ruleLogicalID + "Permission"

	// Create the Events Rule
	rule := e.createEventsRule(functionLogicalID, eventLogicalID)
	resources[ruleLogicalID] = rule

	// Create Lambda Permission for EventBridge to invoke the function
	permission := e.createLambdaPermission(functionLogicalID, ruleLogicalID)
	resources[permissionLogicalID] = permission

	return resources, nil
}

// createEventsRule creates the AWS::Events::Rule CloudFormation resource.
func (e *CloudWatchEvent) createEventsRule(functionLogicalID string, eventLogicalID string) map[string]interface{} {
	ruleLogicalID := functionLogicalID + eventLogicalID

	rule := &events.Rule{
		EventBusName: e.EventBusName,
		State:        e.State,
	}

	// Set either EventPattern or ScheduleExpression
	if e.Pattern != nil {
		rule.EventPattern = e.Pattern
	} else if e.Schedule != nil {
		rule.ScheduleExpression = e.Schedule
	}

	// Create target configuration
	target := events.Target{
		Arn: map[string]interface{}{
			"Fn::GetAtt": []string{functionLogicalID, "Arn"},
		},
	}

	// Set target ID (use custom ID if provided, otherwise generate default)
	if e.Target != nil && e.Target.Id != nil {
		target.Id = e.Target.Id
	} else {
		target.Id = ruleLogicalID + "LambdaTarget"
	}

	// Add Input configuration if provided
	if e.Input != nil {
		target.Input = e.Input
	}

	// Add InputPath if provided
	if e.InputPath != nil {
		target.InputPath = e.InputPath
	}

	// Add InputTransformer if provided
	if e.InputTransformer != nil {
		target.InputTransformer = &events.InputTransformer{
			InputPathsMap: e.InputTransformer.InputPathsMap,
			InputTemplate: e.InputTransformer.InputTemplate,
		}
	}

	// Add DeadLetterConfig if provided
	if e.DeadLetterConfig != nil {
		target.DeadLetterConfig = &events.DeadLetterConfig{
			Arn: e.DeadLetterConfig.Arn,
		}
	}

	// Add RetryPolicy if provided
	if e.RetryPolicy != nil {
		target.RetryPolicy = &events.RetryPolicy{
			MaximumEventAgeInSeconds: e.RetryPolicy.MaximumEventAgeInSeconds,
			MaximumRetryAttempts:     e.RetryPolicy.MaximumRetryAttempts,
		}
	}

	rule.Targets = []events.Target{target}

	return ruleToCloudFormation(rule)
}

// createLambdaPermission creates the AWS::Lambda::Permission resource.
func (e *CloudWatchEvent) createLambdaPermission(functionLogicalID string, ruleLogicalID string) map[string]interface{} {
	permission := lambda.NewEventsPermission(
		map[string]interface{}{"Ref": functionLogicalID},
		map[string]interface{}{
			"Fn::GetAtt": []string{ruleLogicalID, "Arn"},
		},
	)

	return permission.ToCloudFormation()
}

// ruleToCloudFormation converts an events.Rule to a CloudFormation resource.
func ruleToCloudFormation(r *events.Rule) map[string]interface{} {
	properties := make(map[string]interface{})

	if r.Name != nil {
		properties["Name"] = r.Name
	}
	if r.Description != nil {
		properties["Description"] = r.Description
	}
	if r.EventBusName != nil {
		properties["EventBusName"] = r.EventBusName
	}
	if r.EventPattern != nil {
		properties["EventPattern"] = r.EventPattern
	}
	if r.ScheduleExpression != nil {
		properties["ScheduleExpression"] = r.ScheduleExpression
	}
	if r.State != "" {
		properties["State"] = r.State
	}
	if len(r.Targets) > 0 {
		targets := make([]map[string]interface{}, len(r.Targets))
		for i, target := range r.Targets {
			targets[i] = targetToMap(target)
		}
		properties["Targets"] = targets
	}
	if r.RoleArn != nil {
		properties["RoleArn"] = r.RoleArn
	}

	return map[string]interface{}{
		"Type":       "AWS::Events::Rule",
		"Properties": properties,
	}
}

// targetToMap converts an events.Target to a map for CloudFormation.
func targetToMap(t events.Target) map[string]interface{} {
	target := make(map[string]interface{})

	target["Id"] = t.Id
	target["Arn"] = t.Arn

	if t.RoleArn != nil {
		target["RoleArn"] = t.RoleArn
	}
	if t.Input != nil {
		target["Input"] = t.Input
	}
	if t.InputPath != nil {
		target["InputPath"] = t.InputPath
	}
	if t.InputTransformer != nil {
		inputTransformer := make(map[string]interface{})
		if t.InputTransformer.InputPathsMap != nil {
			inputTransformer["InputPathsMap"] = t.InputTransformer.InputPathsMap
		}
		inputTransformer["InputTemplate"] = t.InputTransformer.InputTemplate
		target["InputTransformer"] = inputTransformer
	}
	if t.DeadLetterConfig != nil {
		dlc := make(map[string]interface{})
		if t.DeadLetterConfig.Arn != nil {
			dlc["Arn"] = t.DeadLetterConfig.Arn
		}
		target["DeadLetterConfig"] = dlc
	}
	if t.RetryPolicy != nil {
		rp := make(map[string]interface{})
		if t.RetryPolicy.MaximumEventAgeInSeconds != nil {
			rp["MaximumEventAgeInSeconds"] = t.RetryPolicy.MaximumEventAgeInSeconds
		}
		if t.RetryPolicy.MaximumRetryAttempts != nil {
			rp["MaximumRetryAttempts"] = t.RetryPolicy.MaximumRetryAttempts
		}
		target["RetryPolicy"] = rp
	}

	return target
}

// Validate checks if the CloudWatch Event configuration is valid.
func (e *CloudWatchEvent) Validate() error {
	// Must have either Pattern or Schedule, but not both
	hasPattern := e.Pattern != nil
	hasSchedule := e.Schedule != nil

	if !hasPattern && !hasSchedule {
		return fmt.Errorf("CloudWatchEvent must specify either Pattern or Schedule")
	}

	if hasPattern && hasSchedule {
		return fmt.Errorf("CloudWatchEvent cannot specify both Pattern and Schedule")
	}

	return nil
}
