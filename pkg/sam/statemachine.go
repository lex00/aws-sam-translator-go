// Package sam provides SAM resource transformers.
package sam

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/model/iam"
)

// StateMachine represents an AWS::Serverless::StateMachine resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-statemachine.html
type StateMachine struct {
	// Name is the name of the state machine.
	Name string `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Type is the state machine type. Valid values: STANDARD | EXPRESS
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`

	// Definition is the Amazon States Language definition object.
	Definition map[string]interface{} `json:"Definition,omitempty" yaml:"Definition,omitempty"`

	// DefinitionUri is the S3 URI or local path to the definition file.
	DefinitionUri interface{} `json:"DefinitionUri,omitempty" yaml:"DefinitionUri,omitempty"`

	// DefinitionSubstitutions is a map of key-value pairs for definition substitutions.
	DefinitionSubstitutions map[string]interface{} `json:"DefinitionSubstitutions,omitempty" yaml:"DefinitionSubstitutions,omitempty"`

	// Role is the ARN of the IAM role to use for execution. If not specified, a role is generated.
	Role string `json:"Role,omitempty" yaml:"Role,omitempty"`

	// RolePath is the path for the generated IAM role.
	RolePath string `json:"RolePath,omitempty" yaml:"RolePath,omitempty"`

	// Policies are policies to attach to the generated role.
	Policies interface{} `json:"Policies,omitempty" yaml:"Policies,omitempty"`

	// PermissionsBoundary is the ARN of the permissions boundary for the role.
	PermissionsBoundary string `json:"PermissionsBoundary,omitempty" yaml:"PermissionsBoundary,omitempty"`

	// Tracing specifies X-Ray tracing configuration.
	Tracing *TracingConfig `json:"Tracing,omitempty" yaml:"Tracing,omitempty"`

	// Logging specifies logging configuration.
	Logging *LoggingConfig `json:"Logging,omitempty" yaml:"Logging,omitempty"`

	// Tags is a map of key-value pairs to apply to the state machine.
	Tags map[string]interface{} `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// Events is a map of event sources that trigger the state machine.
	Events map[string]interface{} `json:"Events,omitempty" yaml:"Events,omitempty"`
}

// TracingConfig specifies X-Ray tracing configuration.
type TracingConfig struct {
	Enabled bool `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`
}

// LoggingConfig specifies logging configuration for state machines.
type LoggingConfig struct {
	Level                string        `json:"Level,omitempty" yaml:"Level,omitempty"`
	IncludeExecutionData bool          `json:"IncludeExecutionData,omitempty" yaml:"IncludeExecutionData,omitempty"`
	Destinations         []interface{} `json:"Destinations,omitempty" yaml:"Destinations,omitempty"`
}

// StateMachineTransformer transforms AWS::Serverless::StateMachine to CloudFormation.
type StateMachineTransformer struct{}

// NewStateMachineTransformer creates a new StateMachineTransformer.
func NewStateMachineTransformer() *StateMachineTransformer {
	return &StateMachineTransformer{}
}

// Transform converts a SAM StateMachine to CloudFormation resources.
func (t *StateMachineTransformer) Transform(logicalID string, sm *StateMachine) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Build the state machine properties
	props := make(map[string]interface{})

	// Set StateMachineName if specified
	if sm.Name != "" {
		props["StateMachineName"] = sm.Name
	}

	// Set StateMachineType if specified
	if sm.Type != "" {
		props["StateMachineType"] = sm.Type
	}

	// Handle definition
	if sm.Definition != nil {
		// Convert inline definition to DefinitionString with Fn::Join
		defString, err := t.convertDefinitionToString(sm.Definition)
		if err != nil {
			return nil, fmt.Errorf("failed to convert definition to string: %w", err)
		}
		props["DefinitionString"] = defString
	} else if sm.DefinitionUri != nil {
		// Handle DefinitionUri (S3 location)
		s3Location, err := t.buildDefinitionS3Location(sm.DefinitionUri)
		if err != nil {
			return nil, fmt.Errorf("failed to build DefinitionS3Location: %w", err)
		}
		props["DefinitionS3Location"] = s3Location
	}

	// Add DefinitionSubstitutions if specified
	if sm.DefinitionSubstitutions != nil {
		props["DefinitionSubstitutions"] = sm.DefinitionSubstitutions
	}

	// Resolve role (explicit or generated)
	roleLogicalID := ""
	if sm.Role != "" {
		// Use explicit role
		props["RoleArn"] = sm.Role
	} else {
		// Generate a role
		roleLogicalID = logicalID + "Role"
		role, err := t.generateRole(logicalID, sm)
		if err != nil {
			return nil, fmt.Errorf("failed to generate role: %w", err)
		}
		resources[roleLogicalID] = role.ToResource()

		// Reference the generated role
		props["RoleArn"] = map[string]interface{}{
			"Fn::GetAtt": []string{roleLogicalID, "Arn"},
		}
	}

	// Build tags (always include SAM tag first)
	tags := []map[string]interface{}{
		{"Key": "stateMachine:createdBy", "Value": "SAM"},
	}
	if sm.Tags != nil {
		for k, v := range sm.Tags {
			tags = append(tags, map[string]interface{}{"Key": k, "Value": v})
		}
	}
	props["Tags"] = tags

	// Set tracing configuration
	if sm.Tracing != nil {
		props["TracingConfiguration"] = map[string]interface{}{
			"Enabled": sm.Tracing.Enabled,
		}
		// Add X-Ray policy to the generated role
		if roleLogicalID != "" {
			t.addXRayPolicyToRole(resources, roleLogicalID)
		}
	}

	// Set logging configuration
	if sm.Logging != nil {
		loggingConfig := make(map[string]interface{})
		if sm.Logging.Level != "" {
			loggingConfig["Level"] = sm.Logging.Level
		}
		loggingConfig["IncludeExecutionData"] = sm.Logging.IncludeExecutionData
		if len(sm.Logging.Destinations) > 0 {
			loggingConfig["Destinations"] = sm.Logging.Destinations
		}
		props["LoggingConfiguration"] = loggingConfig
	}

	// Add the state machine resource
	resources[logicalID] = map[string]interface{}{
		"Type":       "AWS::StepFunctions::StateMachine",
		"Properties": props,
	}

	// Process events if any
	if len(sm.Events) > 0 {
		if err := t.processEvents(logicalID, sm.Events, resources); err != nil {
			return nil, fmt.Errorf("failed to process events: %w", err)
		}
	}

	return resources, nil
}

// convertDefinitionToString converts an inline definition to a DefinitionString with Fn::Join.
// This preserves intrinsic functions and provides proper JSON formatting.
func (t *StateMachineTransformer) convertDefinitionToString(definition map[string]interface{}) (interface{}, error) {
	// Convert to JSON with indentation
	jsonBytes, err := json.MarshalIndent(definition, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal definition: %w", err)
	}

	// Split by lines for Fn::Join
	lines := strings.Split(string(jsonBytes), "\n")
	lineInterfaces := make([]interface{}, len(lines))
	for i, line := range lines {
		lineInterfaces[i] = line
	}

	// Return Fn::Join to join lines with newlines
	return map[string]interface{}{
		"Fn::Join": []interface{}{
			"\n",
			lineInterfaces,
		},
	}, nil
}

// buildDefinitionS3Location builds an S3 location from DefinitionUri.
func (t *StateMachineTransformer) buildDefinitionS3Location(definitionUri interface{}) (map[string]interface{}, error) {
	switch uri := definitionUri.(type) {
	case string:
		return t.parseDefinitionUri(uri)
	case map[string]interface{}:
		// Already an S3 location object
		s3Location := make(map[string]interface{})
		if bucket, ok := uri["Bucket"]; ok {
			s3Location["Bucket"] = bucket
		}
		if key, ok := uri["Key"]; ok {
			s3Location["Key"] = key
		}
		if version, ok := uri["Version"]; ok {
			s3Location["Version"] = version
		}
		return s3Location, nil
	default:
		return nil, fmt.Errorf("unsupported DefinitionUri type: %T", definitionUri)
	}
}

// parseDefinitionUri parses an S3 URI string into an S3 location map.
func (t *StateMachineTransformer) parseDefinitionUri(uri string) (map[string]interface{}, error) {
	if !strings.HasPrefix(uri, "s3://") {
		return nil, fmt.Errorf("DefinitionUri must be an S3 URI (s3://...): %s", uri)
	}

	// Remove s3:// prefix
	path := strings.TrimPrefix(uri, "s3://")

	// Split into bucket and key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid S3 URI (missing key): %s", uri)
	}

	return map[string]interface{}{
		"Bucket": parts[0],
		"Key":    parts[1],
	}, nil
}

// generateRole creates an IAM role for the state machine.
func (t *StateMachineTransformer) generateRole(logicalID string, sm *StateMachine) (*iam.Role, error) {
	// Create assume role policy for Step Functions
	trustPolicy := iam.NewServiceTrustRelationship(iam.ServiceStepFunctions).ToPolicyDocument()
	role := iam.NewRole(trustPolicy)

	// Set role path if specified
	if sm.RolePath != "" {
		role.WithPath(sm.RolePath)
	}

	// Set permissions boundary if specified
	if sm.PermissionsBoundary != "" {
		role.WithPermissionsBoundary(sm.PermissionsBoundary)
	}

	// Process policies
	if sm.Policies != nil {
		if err := t.addPoliciesToRole(role, logicalID, sm.Policies); err != nil {
			return nil, err
		}
	}

	return role, nil
}

// addPoliciesToRole adds policies from the SAM StateMachine to the role.
func (t *StateMachineTransformer) addPoliciesToRole(role *iam.Role, logicalID string, policies interface{}) error {
	switch p := policies.(type) {
	case string:
		// Single managed policy ARN
		if strings.HasPrefix(p, "arn:") {
			role.AddManagedPolicyArn(p)
		}

	case []interface{}:
		// Array of policies
		for i, policy := range p {
			if err := t.addSinglePolicyToRole(role, logicalID, policy, i); err != nil {
				return err
			}
		}

	case map[string]interface{}:
		// Single inline policy document
		if err := t.addSinglePolicyToRole(role, logicalID, p, 0); err != nil {
			return err
		}
	}

	return nil
}

// addSinglePolicyToRole adds a single policy to the role.
func (t *StateMachineTransformer) addSinglePolicyToRole(role *iam.Role, logicalID string, policy interface{}, index int) error {
	switch p := policy.(type) {
	case string:
		// Managed policy ARN
		role.AddManagedPolicyArn(p)

	case map[string]interface{}:
		// Check if it's an inline policy document (has Statement)
		if _, hasStatement := p["Statement"]; hasStatement {
			policyName := fmt.Sprintf("%sRolePolicy%d", logicalID, index)
			policyDoc := t.buildPolicyDocumentFromMap(p)
			role.AddInlinePolicy(policyName, policyDoc)
		}
	}

	return nil
}

// buildPolicyDocumentFromMap creates a PolicyDocument from a map.
func (t *StateMachineTransformer) buildPolicyDocumentFromMap(m map[string]interface{}) *iam.PolicyDocument {
	doc := iam.NewPolicyDocument()

	if version, ok := m["Version"].(string); ok {
		doc.Version = version
	}

	if statements, ok := m["Statement"].([]interface{}); ok {
		for _, s := range statements {
			if stmtMap, ok := s.(map[string]interface{}); ok {
				stmt := t.buildStatementFromMap(stmtMap)
				doc.AddStatement(stmt)
			}
		}
	}

	return doc
}

// buildStatementFromMap creates a Statement from a map.
func (t *StateMachineTransformer) buildStatementFromMap(m map[string]interface{}) *iam.Statement {
	effect := iam.EffectAllow
	if e, ok := m["Effect"].(string); ok {
		effect = e
	}

	stmt := iam.NewStatement(effect)

	if sid, ok := m["Sid"].(string); ok {
		stmt.WithSid(sid)
	}
	if action, ok := m["Action"]; ok {
		stmt.WithAction(action)
	}
	if resource, ok := m["Resource"]; ok {
		stmt.WithResource(resource)
	}
	if principal, ok := m["Principal"]; ok {
		stmt.WithPrincipal(principal)
	}
	if condition, ok := m["Condition"].(map[string]interface{}); ok {
		stmt.WithConditions(condition)
	}

	return stmt
}

// addXRayPolicyToRole adds the X-Ray managed policy to the role.
func (t *StateMachineTransformer) addXRayPolicyToRole(resources map[string]interface{}, roleLogicalID string) {
	roleResource, ok := resources[roleLogicalID].(map[string]interface{})
	if !ok {
		return
	}

	props, ok := roleResource["Properties"].(map[string]interface{})
	if !ok {
		return
	}

	// Get or create ManagedPolicyArns
	var managedPolicies []interface{}
	if existing, ok := props["ManagedPolicyArns"].([]interface{}); ok {
		managedPolicies = existing
	}

	// Add X-Ray policy
	managedPolicies = append(managedPolicies, "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess")
	props["ManagedPolicyArns"] = managedPolicies
}

// processEvents processes event sources for the state machine.
func (t *StateMachineTransformer) processEvents(logicalID string, events map[string]interface{}, resources map[string]interface{}) error {
	for eventName, eventConfig := range events {
		eventMap, ok := eventConfig.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid event configuration for %s", eventName)
		}

		eventType, ok := eventMap["Type"].(string)
		if !ok {
			return fmt.Errorf("event %s missing Type", eventName)
		}

		eventResourceID := logicalID + eventName
		eventProps := make(map[string]interface{})

		// Get properties if present
		if props, ok := eventMap["Properties"].(map[string]interface{}); ok {
			eventProps = props
		}

		switch eventType {
		case "Schedule":
			if err := t.processScheduleEvent(logicalID, eventResourceID, eventProps, resources); err != nil {
				return fmt.Errorf("failed to process Schedule event %s: %w", eventName, err)
			}
		case "CloudWatchEvent":
			if err := t.processCloudWatchEvent(logicalID, eventResourceID, eventProps, resources); err != nil {
				return fmt.Errorf("failed to process CloudWatchEvent %s: %w", eventName, err)
			}
		case "ScheduleV2":
			if err := t.processScheduleV2Event(logicalID, eventResourceID, eventProps, resources); err != nil {
				return fmt.Errorf("failed to process ScheduleV2 event %s: %w", eventName, err)
			}
		default:
			return fmt.Errorf("unsupported event type: %s", eventType)
		}
	}

	return nil
}

// processScheduleEvent processes a Schedule event source.
func (t *StateMachineTransformer) processScheduleEvent(stateMachineID, eventID string, props map[string]interface{}, resources map[string]interface{}) error {
	// Create invocation role for the schedule
	roleID := eventID + "Role"
	invocationRole := iam.NewEventsInvocationRole()

	// Add policy to start state machine execution
	startPolicy := iam.NewPolicyDocument().AddStatement(
		iam.NewStatement(iam.EffectAllow).
			WithAction("states:StartExecution").
			WithResource(map[string]interface{}{
				"Ref": stateMachineID,
			}),
	)
	invocationRole.AddInlinePolicy(roleID+"Policy", startPolicy)
	resources[roleID] = invocationRole.ToResource()

	// Build the Events::Rule
	ruleProps := make(map[string]interface{})

	if schedule, ok := props["Schedule"].(string); ok {
		ruleProps["ScheduleExpression"] = schedule
	}
	if name, ok := props["Name"].(string); ok {
		ruleProps["Name"] = name
	}
	if description, ok := props["Description"].(string); ok {
		ruleProps["Description"] = description
	}

	// Handle Enabled property
	if enabled, ok := props["Enabled"].(bool); ok {
		if enabled {
			ruleProps["State"] = "ENABLED"
		} else {
			ruleProps["State"] = "DISABLED"
		}
	} else {
		ruleProps["State"] = "ENABLED"
	}

	// Build target
	target := map[string]interface{}{
		"Id": eventID + "StepFunctionsTarget",
		"Arn": map[string]interface{}{
			"Ref": stateMachineID,
		},
		"RoleArn": map[string]interface{}{
			"Fn::GetAtt": []string{roleID, "Arn"},
		},
	}

	// Add input if specified
	if input, ok := props["Input"]; ok {
		target["Input"] = input
	}

	ruleProps["Targets"] = []interface{}{target}

	resources[eventID] = map[string]interface{}{
		"Type":       "AWS::Events::Rule",
		"Properties": ruleProps,
	}

	return nil
}

// processCloudWatchEvent processes a CloudWatchEvent event source.
func (t *StateMachineTransformer) processCloudWatchEvent(stateMachineID, eventID string, props map[string]interface{}, resources map[string]interface{}) error {
	// Create invocation role
	roleID := eventID + "Role"
	invocationRole := iam.NewEventsInvocationRole()

	startPolicy := iam.NewPolicyDocument().AddStatement(
		iam.NewStatement(iam.EffectAllow).
			WithAction("states:StartExecution").
			WithResource(map[string]interface{}{
				"Ref": stateMachineID,
			}),
	)
	invocationRole.AddInlinePolicy(roleID+"Policy", startPolicy)
	resources[roleID] = invocationRole.ToResource()

	// Build the Events::Rule
	ruleProps := make(map[string]interface{})

	if pattern, ok := props["Pattern"].(map[string]interface{}); ok {
		ruleProps["EventPattern"] = pattern
	}
	if ruleName, ok := props["RuleName"].(string); ok {
		ruleProps["Name"] = ruleName
	}
	if state, ok := props["State"].(string); ok {
		ruleProps["State"] = state
	} else {
		ruleProps["State"] = "ENABLED"
	}

	// Build target
	target := map[string]interface{}{
		"Id": eventID + "StepFunctionsTarget",
		"Arn": map[string]interface{}{
			"Ref": stateMachineID,
		},
		"RoleArn": map[string]interface{}{
			"Fn::GetAtt": []string{roleID, "Arn"},
		},
	}

	ruleProps["Targets"] = []interface{}{target}

	resources[eventID] = map[string]interface{}{
		"Type":       "AWS::Events::Rule",
		"Properties": ruleProps,
	}

	return nil
}

// processScheduleV2Event processes a ScheduleV2 event source (EventBridge Scheduler).
func (t *StateMachineTransformer) processScheduleV2Event(stateMachineID, eventID string, props map[string]interface{}, resources map[string]interface{}) error {
	// Create invocation role for the scheduler
	roleID := eventID + "Role"

	// ScheduleV2 uses scheduler.amazonaws.com as the service principal
	trustPolicy := iam.NewServiceTrustRelationship("scheduler.amazonaws.com").ToPolicyDocument()
	invocationRole := iam.NewRole(trustPolicy)

	// Add policy to start state machine execution
	startPolicy := iam.NewPolicyDocument().AddStatement(
		iam.NewStatement(iam.EffectAllow).
			WithAction("states:StartExecution").
			WithResource(map[string]interface{}{
				"Ref": stateMachineID,
			}),
	)
	invocationRole.AddInlinePolicy(roleID+"Policy", startPolicy)
	resources[roleID] = invocationRole.ToResource()

	// Build the Scheduler::Schedule
	scheduleProps := make(map[string]interface{})

	if scheduleExpr, ok := props["ScheduleExpression"].(string); ok {
		scheduleProps["ScheduleExpression"] = scheduleExpr
	}
	if name, ok := props["Name"].(string); ok {
		scheduleProps["Name"] = name
	}
	if groupName, ok := props["GroupName"].(string); ok {
		scheduleProps["GroupName"] = groupName
	}
	if description, ok := props["Description"].(string); ok {
		scheduleProps["Description"] = description
	}

	// Handle FlexibleTimeWindow (defaults to OFF)
	if ftw, ok := props["FlexibleTimeWindow"].(map[string]interface{}); ok {
		scheduleProps["FlexibleTimeWindow"] = ftw
	} else {
		scheduleProps["FlexibleTimeWindow"] = map[string]interface{}{
			"Mode": "OFF",
		}
	}

	// Handle State
	if state, ok := props["State"].(string); ok {
		scheduleProps["State"] = state
	}

	// Build target
	target := map[string]interface{}{
		"Arn": map[string]interface{}{
			"Ref": stateMachineID,
		},
		"RoleArn": map[string]interface{}{
			"Fn::GetAtt": []string{roleID, "Arn"},
		},
	}

	// Add input if specified
	if input, ok := props["Input"]; ok {
		target["Input"] = input
	}

	scheduleProps["Target"] = target

	resources[eventID] = map[string]interface{}{
		"Type":       "AWS::Scheduler::Schedule",
		"Properties": scheduleProps,
	}

	return nil
}
