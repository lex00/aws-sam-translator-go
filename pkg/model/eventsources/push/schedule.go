// Package push provides CloudFormation resource generators for push event sources
// that trigger Lambda functions (EventBridge, SNS, S3, etc.).
package push

import (
	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/events"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// ScheduleEvent represents a SAM Schedule event source that triggers a Lambda function
// on a recurring schedule using Amazon EventBridge (formerly CloudWatch Events).
//
// SAM Template Syntax:
//
//	Events:
//	  MySchedule:
//	    Type: Schedule
//	    Properties:
//	      Schedule: rate(1 minute)  # or cron(0 12 * * ? *)
//	      Name: my-schedule-rule    # Optional
//	      Description: "My schedule" # Optional
//	      Enabled: true             # Optional, defaults to true
//	      State: ENABLED            # Optional, overrides Enabled
//	      Input: '{"key": "value"}' # Optional
//	      DeadLetterConfig:         # Optional
//	        Type: SQS               # or ARN
//	        TargetArn: ...          # if Type is ARN
//	        QueueLogicalId: ...     # if Type is SQS
//	      RetryPolicy:              # Optional
//	        MaximumEventAgeInSeconds: 3600
//	        MaximumRetryAttempts: 2
//
// CloudFormation Resources Generated:
//   - AWS::Events::Rule - The EventBridge rule with the schedule expression
//   - AWS::Lambda::Permission - Grants EventBridge permission to invoke the function
type ScheduleEvent struct {
	// Schedule is the schedule expression (required).
	// Can be a rate expression like "rate(1 minute)" or a cron expression like "cron(0 12 * * ? *)".
	Schedule interface{} `json:"Schedule" yaml:"Schedule"`

	// Name is the name of the EventBridge rule (optional).
	// If not specified, a name is auto-generated.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Description is the description of the rule (optional).
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// Enabled indicates whether the rule is enabled (optional).
	// Defaults to true. If State is also specified, State takes precedence.
	// Deprecated: Use State instead.
	Enabled interface{} `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`

	// State is the state of the rule (optional).
	// Valid values: ENABLED, DISABLED, ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS
	// If specified, this takes precedence over Enabled.
	State interface{} `json:"State,omitempty" yaml:"State,omitempty"`

	// Input is custom input to pass to the Lambda function (optional).
	// Must be valid JSON as a string.
	Input interface{} `json:"Input,omitempty" yaml:"Input,omitempty"`

	// DeadLetterConfig specifies the dead-letter queue for failed invocations (optional).
	DeadLetterConfig *DeadLetterConfig `json:"DeadLetterConfig,omitempty" yaml:"DeadLetterConfig,omitempty"`

	// RetryPolicy specifies retry settings for failed invocations (optional).
	RetryPolicy *RetryPolicy `json:"RetryPolicy,omitempty" yaml:"RetryPolicy,omitempty"`
}

// ToCloudFormationResources converts the Schedule event source to CloudFormation resources.
// It generates:
//  1. An AWS::Events::Rule resource with the schedule expression and Lambda target
//  2. An AWS::Lambda::Permission resource to allow EventBridge to invoke the function
//
// Parameters:
//   - functionLogicalId: The logical ID of the Lambda function resource
//   - functionArn: The ARN or Ref to the Lambda function
//   - eventLogicalId: The logical ID for the event source (e.g., "Schedule1")
//
// Returns:
//   - A map of CloudFormation resources keyed by logical ID
func (s *ScheduleEvent) ToCloudFormationResources(functionLogicalId string, functionArn interface{}, eventLogicalId string) map[string]interface{} {
	resources := make(map[string]interface{})

	// Generate logical IDs
	ruleLogicalId := functionLogicalId + eventLogicalId
	permissionLogicalId := ruleLogicalId + "Permission"
	targetId := ruleLogicalId + "LambdaTarget"

	// Determine the State value
	state := s.determineState()

	// Build the EventBridge Rule target
	target := events.Target{
		Id:  targetId,
		Arn: functionArn,
	}

	// Add optional target properties
	if s.Input != nil {
		target.Input = s.Input
	}

	if s.DeadLetterConfig != nil {
		dlqArn := s.resolveDLQArn()
		if dlqArn != nil {
			target.DeadLetterConfig = &events.DeadLetterConfig{
				Arn: dlqArn,
			}
		}
	}

	if s.RetryPolicy != nil {
		target.RetryPolicy = &events.RetryPolicy{
			MaximumEventAgeInSeconds: s.RetryPolicy.MaximumEventAgeInSeconds,
			MaximumRetryAttempts:     s.RetryPolicy.MaximumRetryAttempts,
		}
	}

	// Build the EventBridge Rule
	rule := events.Rule{
		ScheduleExpression: s.Schedule,
		Targets:            []events.Target{target},
	}

	// Add optional rule properties
	if s.Name != nil {
		rule.Name = s.Name
	}

	if s.Description != nil {
		rule.Description = s.Description
	}

	// Convert rule to CloudFormation format
	resources[ruleLogicalId] = s.ruleToCloudFormation(rule, state)

	// Generate Lambda permission
	ruleArnRef := map[string]interface{}{
		"Fn::GetAtt": []string{ruleLogicalId, "Arn"},
	}

	permission := lambda.NewEventsPermission(
		map[string]interface{}{"Ref": functionLogicalId},
		ruleArnRef,
	)

	resources[permissionLogicalId] = permission.ToCloudFormation()

	return resources
}

// determineState resolves the State value from State or Enabled properties.
// State takes precedence over Enabled. If neither is set, returns nil.
func (s *ScheduleEvent) determineState() interface{} {
	if s.State != nil {
		return s.State
	}

	// Handle Enabled property (legacy)
	if s.Enabled != nil {
		// If Enabled is a boolean value, convert to ENABLED/DISABLED
		switch v := s.Enabled.(type) {
		case bool:
			if v {
				return "ENABLED"
			}
			return "DISABLED"
		case string:
			// If it's already a string like "true" or "false", convert
			switch v {
			case "true":
				return "ENABLED"
			case "false":
				return "DISABLED"
			default:
				// Otherwise, assume it's an intrinsic function or similar, pass through
				return s.Enabled
			}
		default:
			// For intrinsic functions and other complex types, pass through
			return s.Enabled
		}
	}

	return nil
}

// resolveDLQArn resolves the dead-letter queue ARN from the DeadLetterConfig.
// Returns the ARN to use in the EventBridge target configuration.
func (s *ScheduleEvent) resolveDLQArn() interface{} {
	if s.DeadLetterConfig == nil {
		return nil
	}

	// If TargetArn is specified, use it directly
	if s.DeadLetterConfig.TargetArn != nil {
		return s.DeadLetterConfig.TargetArn
	}

	// If Type is SQS and QueueLogicalId is specified, generate a GetAtt
	if s.DeadLetterConfig.Type == "SQS" && s.DeadLetterConfig.QueueLogicalId != "" {
		return map[string]interface{}{
			"Fn::GetAtt": []string{s.DeadLetterConfig.QueueLogicalId, "Arn"},
		}
	}

	return nil
}

// ruleToCloudFormation converts an events.Rule to CloudFormation resource format.
// The state parameter is passed separately to support intrinsic functions.
func (s *ScheduleEvent) ruleToCloudFormation(rule events.Rule, state interface{}) map[string]interface{} {
	properties := make(map[string]interface{})

	if rule.Name != nil {
		properties["Name"] = rule.Name
	}

	if rule.Description != nil {
		properties["Description"] = rule.Description
	}

	if rule.EventBusName != nil {
		properties["EventBusName"] = rule.EventBusName
	}

	if rule.EventPattern != nil {
		properties["EventPattern"] = rule.EventPattern
	}

	if rule.ScheduleExpression != nil {
		properties["ScheduleExpression"] = rule.ScheduleExpression
	}

	// Add state if provided (supports both string and intrinsic functions)
	if state != nil {
		properties["State"] = state
	}

	if len(rule.Targets) > 0 {
		targets := make([]map[string]interface{}, len(rule.Targets))
		for i, target := range rule.Targets {
			targets[i] = s.targetToMap(target)
		}
		properties["Targets"] = targets
	}

	if rule.RoleArn != nil {
		properties["RoleArn"] = rule.RoleArn
	}

	return map[string]interface{}{
		"Type":       "AWS::Events::Rule",
		"Properties": properties,
	}
}

// targetToMap converts an events.Target to a map for CloudFormation.
func (s *ScheduleEvent) targetToMap(target events.Target) map[string]interface{} {
	result := make(map[string]interface{})

	result["Id"] = target.Id
	result["Arn"] = target.Arn

	if target.RoleArn != nil {
		result["RoleArn"] = target.RoleArn
	}

	if target.Input != nil {
		result["Input"] = target.Input
	}

	if target.InputPath != nil {
		result["InputPath"] = target.InputPath
	}

	if target.InputTransformer != nil {
		inputTransformer := map[string]interface{}{
			"InputTemplate": target.InputTransformer.InputTemplate,
		}
		if len(target.InputTransformer.InputPathsMap) > 0 {
			inputTransformer["InputPathsMap"] = target.InputTransformer.InputPathsMap
		}
		result["InputTransformer"] = inputTransformer
	}

	if target.DeadLetterConfig != nil {
		dlc := make(map[string]interface{})
		if target.DeadLetterConfig.Arn != nil {
			dlc["Arn"] = target.DeadLetterConfig.Arn
		}
		if len(dlc) > 0 {
			result["DeadLetterConfig"] = dlc
		}
	}

	if target.RetryPolicy != nil {
		rp := make(map[string]interface{})
		if target.RetryPolicy.MaximumEventAgeInSeconds != nil {
			rp["MaximumEventAgeInSeconds"] = target.RetryPolicy.MaximumEventAgeInSeconds
		}
		if target.RetryPolicy.MaximumRetryAttempts != nil {
			rp["MaximumRetryAttempts"] = target.RetryPolicy.MaximumRetryAttempts
		}
		if len(rp) > 0 {
			result["RetryPolicy"] = rp
		}
	}

	return result
}
