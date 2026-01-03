package push

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCloudWatchEvent_ToCloudFormation_EventPattern(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source":      []string{"aws.ec2"},
			"detail-type": []string{"EC2 Instance State-change Notification"},
		},
	}

	resources, err := event.ToCloudFormation("MyFunction", "MyEvent")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	// Verify we got 2 resources (Rule + Permission)
	if len(resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resources))
	}

	// Check Events Rule
	rule, ok := resources["MyFunctionMyEvent"]
	if !ok {
		t.Fatal("Events Rule resource not found")
	}

	ruleMap, ok := rule.(map[string]interface{})
	if !ok {
		t.Fatal("Rule is not a map")
	}

	if ruleMap["Type"] != "AWS::Events::Rule" {
		t.Errorf("Expected Type AWS::Events::Rule, got %v", ruleMap["Type"])
	}

	properties, ok := ruleMap["Properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Properties is not a map")
	}

	// Verify EventPattern is set
	if properties["EventPattern"] == nil {
		t.Error("EventPattern not set in rule properties")
	}

	// Verify ScheduleExpression is not set
	if properties["ScheduleExpression"] != nil {
		t.Error("ScheduleExpression should not be set when Pattern is used")
	}

	// Verify Targets
	targets, ok := properties["Targets"].([]map[string]interface{})
	if !ok || len(targets) != 1 {
		t.Fatal("Targets should be a slice with 1 target")
	}

	target := targets[0]
	if target["Id"] != "MyFunctionMyEventLambdaTarget" {
		t.Errorf("Expected target Id 'MyFunctionMyEventLambdaTarget', got %v", target["Id"])
	}

	// Check Lambda Permission
	permission, ok := resources["MyFunctionMyEventPermission"]
	if !ok {
		t.Fatal("Lambda Permission resource not found")
	}

	permMap, ok := permission.(map[string]interface{})
	if !ok {
		t.Fatal("Permission is not a map")
	}

	if permMap["Type"] != "AWS::Lambda::Permission" {
		t.Errorf("Expected Type AWS::Lambda::Permission, got %v", permMap["Type"])
	}

	permProps, ok := permMap["Properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Permission Properties is not a map")
	}

	if permProps["Action"] != "lambda:InvokeFunction" {
		t.Errorf("Expected Action lambda:InvokeFunction, got %v", permProps["Action"])
	}

	if permProps["Principal"] != "events.amazonaws.com" {
		t.Errorf("Expected Principal events.amazonaws.com, got %v", permProps["Principal"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_Schedule(t *testing.T) {
	event := &CloudWatchEvent{
		Schedule: "rate(5 minutes)",
		State:    "ENABLED",
	}

	resources, err := event.ToCloudFormation("ScheduledFunc", "Timer")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule, ok := resources["ScheduledFuncTimer"]
	if !ok {
		t.Fatal("Events Rule resource not found")
	}

	ruleMap := rule.(map[string]interface{})
	properties := ruleMap["Properties"].(map[string]interface{})

	// Verify ScheduleExpression is set
	if properties["ScheduleExpression"] != "rate(5 minutes)" {
		t.Errorf("Expected ScheduleExpression 'rate(5 minutes)', got %v", properties["ScheduleExpression"])
	}

	// Verify EventPattern is not set
	if properties["EventPattern"] != nil {
		t.Error("EventPattern should not be set when Schedule is used")
	}

	// Verify State is set
	if properties["State"] != "ENABLED" {
		t.Errorf("Expected State ENABLED, got %v", properties["State"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithEventBusName(t *testing.T) {
	event := &CloudWatchEvent{
		EventBusName: "CustomEventBus",
		Pattern: map[string]interface{}{
			"source": []string{"custom.app"},
		},
	}

	resources, err := event.ToCloudFormation("MyFunc", "CustomEvent")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncCustomEvent"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})

	if properties["EventBusName"] != "CustomEventBus" {
		t.Errorf("Expected EventBusName 'CustomEventBus', got %v", properties["EventBusName"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithCustomTargetId(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
		Target: &Target{
			Id: "CustomTargetId123",
		},
	}

	resources, err := event.ToCloudFormation("MyFunc", "Event1")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncEvent1"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})
	targets := properties["Targets"].([]map[string]interface{})

	if targets[0]["Id"] != "CustomTargetId123" {
		t.Errorf("Expected target Id 'CustomTargetId123', got %v", targets[0]["Id"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithInput(t *testing.T) {
	inputJSON := `{"key": "value"}`
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
		Input: inputJSON,
	}

	resources, err := event.ToCloudFormation("MyFunc", "Event1")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncEvent1"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})
	targets := properties["Targets"].([]map[string]interface{})

	if targets[0]["Input"] != inputJSON {
		t.Errorf("Expected Input to be set, got %v", targets[0]["Input"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithInputPath(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
		InputPath: "$.detail",
	}

	resources, err := event.ToCloudFormation("MyFunc", "Event1")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncEvent1"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})
	targets := properties["Targets"].([]map[string]interface{})

	if targets[0]["InputPath"] != "$.detail" {
		t.Errorf("Expected InputPath '$.detail', got %v", targets[0]["InputPath"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithInputTransformer(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
		InputTransformer: &InputTransformer{
			InputPathsMap: map[string]interface{}{
				"instance": "$.detail.instance-id",
				"state":    "$.detail.state",
			},
			InputTemplate: `{"instance": <instance>, "state": <state>}`,
		},
	}

	resources, err := event.ToCloudFormation("MyFunc", "Event1")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncEvent1"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})
	targets := properties["Targets"].([]map[string]interface{})

	inputTransformer, ok := targets[0]["InputTransformer"].(map[string]interface{})
	if !ok {
		t.Fatal("InputTransformer not found or not a map")
	}

	inputPathsMap, ok := inputTransformer["InputPathsMap"].(map[string]interface{})
	if !ok {
		t.Fatal("InputPathsMap not found or not a map")
	}

	if inputPathsMap["instance"] != "$.detail.instance-id" {
		t.Errorf("Expected instance path '$.detail.instance-id', got %v", inputPathsMap["instance"])
	}

	if inputTransformer["InputTemplate"] != `{"instance": <instance>, "state": <state>}` {
		t.Errorf("InputTemplate mismatch")
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithDeadLetterConfig(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
		DeadLetterConfig: &DeadLetterConfig{
			Arn: "arn:aws:sqs:us-east-1:123456789012:my-dlq",
		},
	}

	resources, err := event.ToCloudFormation("MyFunc", "Event1")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncEvent1"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})
	targets := properties["Targets"].([]map[string]interface{})

	dlc, ok := targets[0]["DeadLetterConfig"].(map[string]interface{})
	if !ok {
		t.Fatal("DeadLetterConfig not found or not a map")
	}

	if dlc["Arn"] != "arn:aws:sqs:us-east-1:123456789012:my-dlq" {
		t.Errorf("DeadLetterConfig Arn mismatch, got %v", dlc["Arn"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithRetryPolicy(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
		RetryPolicy: &RetryPolicy{
			MaximumEventAgeInSeconds: 3600,
			MaximumRetryAttempts:     2,
		},
	}

	resources, err := event.ToCloudFormation("MyFunc", "Event1")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncEvent1"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})
	targets := properties["Targets"].([]map[string]interface{})

	retryPolicy, ok := targets[0]["RetryPolicy"].(map[string]interface{})
	if !ok {
		t.Fatal("RetryPolicy not found or not a map")
	}

	if retryPolicy["MaximumEventAgeInSeconds"] != 3600 {
		t.Errorf("Expected MaximumEventAgeInSeconds 3600, got %v", retryPolicy["MaximumEventAgeInSeconds"])
	}

	if retryPolicy["MaximumRetryAttempts"] != 2 {
		t.Errorf("Expected MaximumRetryAttempts 2, got %v", retryPolicy["MaximumRetryAttempts"])
	}
}

func TestCloudWatchEvent_ToCloudFormation_WithIntrinsicFunctions(t *testing.T) {
	event := &CloudWatchEvent{
		EventBusName: map[string]interface{}{
			"Ref": "MyEventBus",
		},
		Pattern: map[string]interface{}{
			"source": []string{"custom.app"},
		},
	}

	resources, err := event.ToCloudFormation("MyFunc", "Event1")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["MyFuncEvent1"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})

	eventBusName, ok := properties["EventBusName"].(map[string]interface{})
	if !ok {
		t.Fatal("EventBusName should be a map (intrinsic function)")
	}

	if eventBusName["Ref"] != "MyEventBus" {
		t.Errorf("Expected Ref to MyEventBus, got %v", eventBusName["Ref"])
	}
}

func TestCloudWatchEvent_Validate_Valid_Pattern(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
	}

	err := event.Validate()
	if err != nil {
		t.Errorf("Validation should pass for valid Pattern, got error: %v", err)
	}
}

func TestCloudWatchEvent_Validate_Valid_Schedule(t *testing.T) {
	event := &CloudWatchEvent{
		Schedule: "rate(5 minutes)",
	}

	err := event.Validate()
	if err != nil {
		t.Errorf("Validation should pass for valid Schedule, got error: %v", err)
	}
}

func TestCloudWatchEvent_Validate_Missing_PatternAndSchedule(t *testing.T) {
	event := &CloudWatchEvent{}

	err := event.Validate()
	if err == nil {
		t.Error("Validation should fail when both Pattern and Schedule are missing")
	}
}

func TestCloudWatchEvent_Validate_Both_PatternAndSchedule(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
		Schedule: "rate(5 minutes)",
	}

	err := event.Validate()
	if err == nil {
		t.Error("Validation should fail when both Pattern and Schedule are specified")
	}
}

func TestCloudWatchEvent_JSONSerialization(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source":      []string{"aws.ec2"},
			"detail-type": []string{"EC2 Instance State-change Notification"},
		},
		State: "ENABLED",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal CloudWatchEvent to JSON: %v", err)
	}

	var unmarshaled CloudWatchEvent
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal CloudWatchEvent from JSON: %v", err)
	}

	if unmarshaled.State != event.State {
		t.Errorf("State mismatch after JSON round-trip: got %v, want %v", unmarshaled.State, event.State)
	}
}

func TestCloudWatchEvent_YAMLSerialization(t *testing.T) {
	event := &CloudWatchEvent{
		Schedule: "rate(1 hour)",
		State:    "DISABLED",
		Input:    `{"key": "value"}`,
	}

	data, err := yaml.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal CloudWatchEvent to YAML: %v", err)
	}

	var unmarshaled CloudWatchEvent
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal CloudWatchEvent from YAML: %v", err)
	}

	if unmarshaled.State != event.State {
		t.Errorf("State mismatch after YAML round-trip: got %v, want %v", unmarshaled.State, event.State)
	}

	if unmarshaled.Input != event.Input {
		t.Errorf("Input mismatch after YAML round-trip: got %v, want %v", unmarshaled.Input, event.Input)
	}
}

func TestCloudWatchEvent_ComplexEventPattern(t *testing.T) {
	event := &CloudWatchEvent{
		EventBusName: "ExternalEventBridge",
		Pattern: map[string]interface{}{
			"detail": map[string]interface{}{
				"state": []string{"terminated"},
			},
		},
	}

	resources, err := event.ToCloudFormation("TriggeredFunction", "OnTerminate")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	rule := resources["TriggeredFunctionOnTerminate"].(map[string]interface{})
	properties := rule["Properties"].(map[string]interface{})

	// Verify the event pattern is correctly preserved
	eventPattern, ok := properties["EventPattern"].(map[string]interface{})
	if !ok {
		t.Fatal("EventPattern should be a map")
	}

	detail, ok := eventPattern["detail"].(map[string]interface{})
	if !ok {
		t.Fatal("detail should be a map in EventPattern")
	}

	state, ok := detail["state"].([]string)
	if !ok {
		t.Fatal("state should be a string slice in detail")
	}

	if len(state) != 1 || state[0] != "terminated" {
		t.Errorf("Expected state to be ['terminated'], got %v", state)
	}
}

func TestInputTransformer_JSONSerialization(t *testing.T) {
	it := &InputTransformer{
		InputPathsMap: map[string]interface{}{
			"instance": "$.detail.instance-id",
		},
		InputTemplate: `{"instance": <instance>}`,
	}

	data, err := json.Marshal(it)
	if err != nil {
		t.Fatalf("Failed to marshal InputTransformer to JSON: %v", err)
	}

	var unmarshaled InputTransformer
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal InputTransformer from JSON: %v", err)
	}

	if unmarshaled.InputTemplate != it.InputTemplate {
		t.Errorf("InputTemplate mismatch after JSON round-trip")
	}
}

func TestCloudWatchEvent_PermissionSourceArn(t *testing.T) {
	event := &CloudWatchEvent{
		Pattern: map[string]interface{}{
			"source": []string{"aws.ec2"},
		},
	}

	resources, err := event.ToCloudFormation("TestFunc", "TestEvent")
	if err != nil {
		t.Fatalf("ToCloudFormation failed: %v", err)
	}

	permission := resources["TestFuncTestEventPermission"].(map[string]interface{})
	properties := permission["Properties"].(map[string]interface{})

	sourceArn, ok := properties["SourceArn"].(map[string]interface{})
	if !ok {
		t.Fatal("SourceArn should be a map (Fn::GetAtt)")
	}

	getAtt, ok := sourceArn["Fn::GetAtt"].([]string)
	if !ok || len(getAtt) != 2 {
		t.Fatal("SourceArn should use Fn::GetAtt with 2 elements")
	}

	if getAtt[0] != "TestFuncTestEvent" || getAtt[1] != "Arn" {
		t.Errorf("Expected Fn::GetAtt ['TestFuncTestEvent', 'Arn'], got %v", getAtt)
	}
}
