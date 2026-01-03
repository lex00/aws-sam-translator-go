package push

import (
	"encoding/json"
	"testing"
)

func TestScheduleEvent_BasicSchedule(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Fn::GetAtt": []string{"MyFunction", "Arn"},
	}, "Schedule1")

	// Should generate 2 resources: Rule and Permission
	if len(resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resources))
	}

	// Check Rule resource exists
	ruleLogicalId := "MyFunctionSchedule1"
	if _, ok := resources[ruleLogicalId]; !ok {
		t.Errorf("Expected rule resource with ID %s", ruleLogicalId)
	}

	// Check Permission resource exists
	permissionLogicalId := "MyFunctionSchedule1Permission"
	if _, ok := resources[permissionLogicalId]; !ok {
		t.Errorf("Expected permission resource with ID %s", permissionLogicalId)
	}

	// Verify Rule properties
	rule := resources[ruleLogicalId].(map[string]interface{})
	if rule["Type"] != "AWS::Events::Rule" {
		t.Errorf("Expected type AWS::Events::Rule, got %v", rule["Type"])
	}

	props := rule["Properties"].(map[string]interface{})
	if props["ScheduleExpression"] != "rate(1 minute)" {
		t.Errorf("Expected ScheduleExpression 'rate(1 minute)', got %v", props["ScheduleExpression"])
	}

	// Verify targets
	targets := props["Targets"].([]map[string]interface{})
	if len(targets) != 1 {
		t.Errorf("Expected 1 target, got %d", len(targets))
	}

	target := targets[0]
	expectedTargetId := "MyFunctionSchedule1LambdaTarget"
	if target["Id"] != expectedTargetId {
		t.Errorf("Expected target ID %s, got %v", expectedTargetId, target["Id"])
	}

	// Verify Permission properties
	permission := resources[permissionLogicalId].(map[string]interface{})
	if permission["Type"] != "AWS::Lambda::Permission" {
		t.Errorf("Expected type AWS::Lambda::Permission, got %v", permission["Type"])
	}

	permProps := permission["Properties"].(map[string]interface{})
	if permProps["Action"] != "lambda:InvokeFunction" {
		t.Errorf("Expected Action lambda:InvokeFunction, got %v", permProps["Action"])
	}

	if permProps["Principal"] != "events.amazonaws.com" {
		t.Errorf("Expected Principal events.amazonaws.com, got %v", permProps["Principal"])
	}
}

func TestScheduleEvent_WithAllProperties(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule:    "cron(0 12 * * ? *)",
		Name:        "my-schedule-rule",
		Description: "Test Schedule",
		State:       "ENABLED",
		Input:       `{"key": "value"}`,
	}

	resources := schedule.ToCloudFormationResources("TestFunction", map[string]interface{}{
		"Ref": "TestFunction",
	}, "MySchedule")

	ruleLogicalId := "TestFunctionMySchedule"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	// Verify all properties
	if props["ScheduleExpression"] != "cron(0 12 * * ? *)" {
		t.Errorf("Expected cron expression, got %v", props["ScheduleExpression"])
	}

	if props["Name"] != "my-schedule-rule" {
		t.Errorf("Expected Name 'my-schedule-rule', got %v", props["Name"])
	}

	if props["Description"] != "Test Schedule" {
		t.Errorf("Expected Description 'Test Schedule', got %v", props["Description"])
	}

	if props["State"] != "ENABLED" {
		t.Errorf("Expected State 'ENABLED', got %v", props["State"])
	}

	targets := props["Targets"].([]map[string]interface{})
	if targets[0]["Input"] != `{"key": "value"}` {
		t.Errorf("Expected Input JSON, got %v", targets[0]["Input"])
	}
}

func TestScheduleEvent_StateWithIntrinsicFunction(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
		State: map[string]interface{}{
			"Ref": "ScheduleState",
		},
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	// Verify State is preserved as intrinsic function
	state := props["State"].(map[string]interface{})
	if state["Ref"] != "ScheduleState" {
		t.Errorf("Expected State with Ref, got %v", state)
	}
}

func TestScheduleEvent_EnabledProperty(t *testing.T) {
	tests := []struct {
		name          string
		enabled       interface{}
		expectedState string
	}{
		{
			name:          "Enabled true",
			enabled:       true,
			expectedState: "ENABLED",
		},
		{
			name:          "Enabled false",
			enabled:       false,
			expectedState: "DISABLED",
		},
		{
			name:          "Enabled string true",
			enabled:       "true",
			expectedState: "ENABLED",
		},
		{
			name:          "Enabled string false",
			enabled:       "false",
			expectedState: "DISABLED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := &ScheduleEvent{
				Schedule: "rate(1 minute)",
				Enabled:  tt.enabled,
			}

			resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
				"Ref": "MyFunction",
			}, "Schedule1")

			ruleLogicalId := "MyFunctionSchedule1"
			rule := resources[ruleLogicalId].(map[string]interface{})
			props := rule["Properties"].(map[string]interface{})

			if props["State"] != tt.expectedState {
				t.Errorf("Expected State %s, got %v", tt.expectedState, props["State"])
			}
		})
	}
}

func TestScheduleEvent_EnabledWithIntrinsicFunction(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
		Enabled: map[string]interface{}{
			"Fn::Sub": "Enabled",
		},
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	// Verify Enabled is preserved as intrinsic function
	state := props["State"].(map[string]interface{})
	if state["Fn::Sub"] != "Enabled" {
		t.Errorf("Expected State with Fn::Sub, got %v", state)
	}
}

func TestScheduleEvent_StateTakesPrecedenceOverEnabled(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
		State:    "DISABLED",
		Enabled:  true, // Should be ignored
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	if props["State"] != "DISABLED" {
		t.Errorf("Expected State 'DISABLED', got %v", props["State"])
	}
}

func TestScheduleEvent_WithDeadLetterConfig_ARN(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
		DeadLetterConfig: &DeadLetterConfig{
			Type:      "ARN",
			TargetArn: "arn:aws:sqs:us-east-1:123456789012:MyDLQ",
		},
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	targets := props["Targets"].([]map[string]interface{})
	target := targets[0]

	dlc := target["DeadLetterConfig"].(map[string]interface{})
	if dlc["Arn"] != "arn:aws:sqs:us-east-1:123456789012:MyDLQ" {
		t.Errorf("Expected DLQ ARN, got %v", dlc["Arn"])
	}
}

func TestScheduleEvent_WithDeadLetterConfig_SQS(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
		DeadLetterConfig: &DeadLetterConfig{
			Type:           "SQS",
			QueueLogicalId: "MyDLQ",
		},
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	targets := props["Targets"].([]map[string]interface{})
	target := targets[0]

	dlc := target["DeadLetterConfig"].(map[string]interface{})
	arn := dlc["Arn"].(map[string]interface{})

	// Verify Fn::GetAtt is used for SQS queue
	getAtt := arn["Fn::GetAtt"].([]string)
	if len(getAtt) != 2 || getAtt[0] != "MyDLQ" || getAtt[1] != "Arn" {
		t.Errorf("Expected Fn::GetAtt for MyDLQ, got %v", arn)
	}
}

func TestScheduleEvent_WithRetryPolicy(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
		RetryPolicy: &RetryPolicy{
			MaximumEventAgeInSeconds: 3600,
			MaximumRetryAttempts:     2,
		},
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	targets := props["Targets"].([]map[string]interface{})
	target := targets[0]

	rp := target["RetryPolicy"].(map[string]interface{})
	if rp["MaximumEventAgeInSeconds"] != 3600 {
		t.Errorf("Expected MaximumEventAgeInSeconds 3600, got %v", rp["MaximumEventAgeInSeconds"])
	}

	if rp["MaximumRetryAttempts"] != 2 {
		t.Errorf("Expected MaximumRetryAttempts 2, got %v", rp["MaximumRetryAttempts"])
	}
}

func TestScheduleEvent_PermissionReferencesRule(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	permissionLogicalId := "MyFunctionSchedule1Permission"
	permission := resources[permissionLogicalId].(map[string]interface{})
	props := permission["Properties"].(map[string]interface{})

	// Verify SourceArn references the rule
	sourceArn := props["SourceArn"].(map[string]interface{})
	getAtt := sourceArn["Fn::GetAtt"].([]string)

	if len(getAtt) != 2 || getAtt[0] != "MyFunctionSchedule1" || getAtt[1] != "Arn" {
		t.Errorf("Expected SourceArn to reference MyFunctionSchedule1, got %v", getAtt)
	}

	// Verify FunctionName references the function
	functionName := props["FunctionName"].(map[string]interface{})
	if functionName["Ref"] != "MyFunction" {
		t.Errorf("Expected FunctionName Ref to MyFunction, got %v", functionName)
	}
}

func TestScheduleEvent_JSONSerialization(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule:    "rate(5 minutes)",
		Name:        "TestSchedule",
		Description: "Test description",
		State:       "ENABLED",
		Input:       `{"test": true}`,
		RetryPolicy: &RetryPolicy{
			MaximumEventAgeInSeconds: 7200,
			MaximumRetryAttempts:     3,
		},
	}

	data, err := json.Marshal(schedule)
	if err != nil {
		t.Fatalf("Failed to marshal schedule: %v", err)
	}

	var unmarshaled ScheduleEvent
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal schedule: %v", err)
	}

	if unmarshaled.Schedule != schedule.Schedule {
		t.Errorf("Schedule mismatch: got %v, want %v", unmarshaled.Schedule, schedule.Schedule)
	}

	if unmarshaled.Name != schedule.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, schedule.Name)
	}

	// JSON unmarshaling converts numbers to float64
	attempts, ok := unmarshaled.RetryPolicy.MaximumRetryAttempts.(float64)
	if !ok || attempts != 3.0 {
		t.Errorf("MaximumRetryAttempts mismatch: got %v (type %T), want 3", unmarshaled.RetryPolicy.MaximumRetryAttempts, unmarshaled.RetryPolicy.MaximumRetryAttempts)
	}
}

func TestScheduleEvent_MultipleSchedulesForSameFunction(t *testing.T) {
	// Test that different event logical IDs generate different resource IDs
	schedule1 := &ScheduleEvent{
		Schedule: "rate(1 minute)",
	}

	schedule2 := &ScheduleEvent{
		Schedule: "rate(5 minutes)",
	}

	resources1 := schedule1.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	resources2 := schedule2.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule2")

	// Verify they generate different resource IDs
	if _, ok := resources1["MyFunctionSchedule1"]; !ok {
		t.Error("Expected MyFunctionSchedule1 resource")
	}

	if _, ok := resources2["MyFunctionSchedule2"]; !ok {
		t.Error("Expected MyFunctionSchedule2 resource")
	}

	// Verify they have different target IDs
	rule1 := resources1["MyFunctionSchedule1"].(map[string]interface{})
	props1 := rule1["Properties"].(map[string]interface{})
	targets1 := props1["Targets"].([]map[string]interface{})

	rule2 := resources2["MyFunctionSchedule2"].(map[string]interface{})
	props2 := rule2["Properties"].(map[string]interface{})
	targets2 := props2["Targets"].([]map[string]interface{})

	if targets1[0]["Id"] == targets2[0]["Id"] {
		t.Error("Expected different target IDs for different schedules")
	}
}

func TestScheduleEvent_NoStateOrEnabled(t *testing.T) {
	// When neither State nor Enabled is set, State should not be in the output
	schedule := &ScheduleEvent{
		Schedule: "rate(1 minute)",
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Ref": "MyFunction",
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	// State should not be present
	if _, exists := props["State"]; exists {
		t.Error("Expected State to not be present when neither State nor Enabled is set")
	}
}

func TestScheduleEvent_ComplexIntrinsicFunctions(t *testing.T) {
	schedule := &ScheduleEvent{
		Schedule: map[string]interface{}{
			"Fn::Sub": "rate(${ScheduleRate} minutes)",
		},
		Name: map[string]interface{}{
			"Fn::Join": []interface{}{
				"-",
				[]interface{}{"schedule", map[string]interface{}{"Ref": "Environment"}},
			},
		},
		Description: map[string]interface{}{
			"Fn::Sub": "Schedule for ${Environment}",
		},
	}

	resources := schedule.ToCloudFormationResources("MyFunction", map[string]interface{}{
		"Fn::GetAtt": []string{"MyFunction", "Arn"},
	}, "Schedule1")

	ruleLogicalId := "MyFunctionSchedule1"
	rule := resources[ruleLogicalId].(map[string]interface{})
	props := rule["Properties"].(map[string]interface{})

	// Verify intrinsic functions are preserved
	scheduleExpr := props["ScheduleExpression"].(map[string]interface{})
	if scheduleExpr["Fn::Sub"] != "rate(${ScheduleRate} minutes)" {
		t.Errorf("Expected Fn::Sub in ScheduleExpression, got %v", scheduleExpr)
	}

	name := props["Name"].(map[string]interface{})
	if name["Fn::Join"] == nil {
		t.Error("Expected Fn::Join in Name")
	}

	description := props["Description"].(map[string]interface{})
	if description["Fn::Sub"] != "Schedule for ${Environment}" {
		t.Errorf("Expected Fn::Sub in Description, got %v", description)
	}
}
