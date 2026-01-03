package sam

import (
	"encoding/json"
	"testing"
)

func TestStateMachineTransformer_Transform_Minimal(t *testing.T) {
	transformer := NewStateMachineTransformer()

	// Minimal state machine with just a Definition
	sm := &StateMachine{
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States": map[string]interface{}{
				"Hello": map[string]interface{}{
					"Type": "Pass",
					"End":  true,
				},
			},
		},
	}

	resources, err := transformer.Transform("MyStateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have 2 resources: StateMachine and Role
	if len(resources) != 2 {
		t.Errorf("expected 2 resources (StateMachine + Role), got %d", len(resources))
	}

	// Check StateMachine resource exists
	smResource, ok := resources["MyStateMachine"].(map[string]interface{})
	if !ok {
		t.Fatal("MyStateMachine resource not found")
	}

	// Check type
	if smResource["Type"] != "AWS::StepFunctions::StateMachine" {
		t.Errorf("expected Type 'AWS::StepFunctions::StateMachine', got %v", smResource["Type"])
	}

	// Check properties
	props := smResource["Properties"].(map[string]interface{})

	// Should have DefinitionString with Fn::Join
	defStr, defStrOk := props["DefinitionString"].(map[string]interface{})
	if !defStrOk {
		t.Fatal("DefinitionString not found or not a map")
	}
	if _, hasJoin := defStr["Fn::Join"]; !hasJoin {
		t.Error("DefinitionString should contain Fn::Join")
	}

	// Should have RoleArn referencing generated role
	roleArn, roleArnOk := props["RoleArn"].(map[string]interface{})
	if !roleArnOk {
		t.Fatal("RoleArn not found or not a map")
	}
	if getAttr, hasGetAtt := roleArn["Fn::GetAtt"].([]string); hasGetAtt {
		if getAttr[0] != "MyStateMachineRole" || getAttr[1] != "Arn" {
			t.Errorf("expected RoleArn Fn::GetAtt [MyStateMachineRole, Arn], got %v", getAttr)
		}
	} else {
		t.Error("RoleArn should contain Fn::GetAtt")
	}

	// Should have SAM tag
	tags := props["Tags"].([]map[string]interface{})
	if len(tags) < 1 {
		t.Fatal("expected at least 1 tag")
	}
	if tags[0]["Key"] != "stateMachine:createdBy" || tags[0]["Value"] != "SAM" {
		t.Errorf("expected first tag to be stateMachine:createdBy=SAM, got %v", tags[0])
	}

	// Check Role resource exists
	roleResource, ok := resources["MyStateMachineRole"].(map[string]interface{})
	if !ok {
		t.Fatal("MyStateMachineRole resource not found")
	}

	if roleResource["Type"] != "AWS::IAM::Role" {
		t.Errorf("expected Role Type 'AWS::IAM::Role', got %v", roleResource["Type"])
	}

	roleProps := roleResource["Properties"].(map[string]interface{})
	assumeRolePolicy := roleProps["AssumeRolePolicyDocument"].(map[string]interface{})
	statements := assumeRolePolicy["Statement"].([]interface{})
	if len(statements) != 1 {
		t.Fatalf("expected 1 statement in AssumeRolePolicyDocument, got %d", len(statements))
	}
}

func TestStateMachineTransformer_Transform_WithExplicitRole(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Name:          "MyStateMachineWithRole",
		Type:          "STANDARD",
		Role:          "arn:aws:iam::123456123456:role/service-role/SampleRole",
		DefinitionUri: "s3://sam-demo-bucket/my-state-machine.asl.json",
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should only have 1 resource (StateMachine, no Role)
	if len(resources) != 1 {
		t.Errorf("expected 1 resource (StateMachine only), got %d", len(resources))
	}

	smResource := resources["StateMachine"].(map[string]interface{})
	props := smResource["Properties"].(map[string]interface{})

	// Check StateMachineName
	if props["StateMachineName"] != "MyStateMachineWithRole" {
		t.Errorf("expected StateMachineName 'MyStateMachineWithRole', got %v", props["StateMachineName"])
	}

	// Check StateMachineType
	if props["StateMachineType"] != "STANDARD" {
		t.Errorf("expected StateMachineType 'STANDARD', got %v", props["StateMachineType"])
	}

	// Check RoleArn is the explicit ARN
	if props["RoleArn"] != "arn:aws:iam::123456123456:role/service-role/SampleRole" {
		t.Errorf("expected RoleArn to be explicit ARN, got %v", props["RoleArn"])
	}

	// Check DefinitionS3Location
	s3Loc := props["DefinitionS3Location"].(map[string]interface{})
	if s3Loc["Bucket"] != "sam-demo-bucket" {
		t.Errorf("expected Bucket 'sam-demo-bucket', got %v", s3Loc["Bucket"])
	}
	if s3Loc["Key"] != "my-state-machine.asl.json" {
		t.Errorf("expected Key 'my-state-machine.asl.json', got %v", s3Loc["Key"])
	}
}

func TestStateMachineTransformer_Transform_WithTags(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States":  map[string]interface{}{},
		},
		Tags: map[string]interface{}{
			"TagOne": "ValueOne",
			"TagTwo": "ValueTwo",
		},
		Tracing: &TracingConfig{Enabled: true},
	}

	resources, err := transformer.Transform("MyStateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	smResource := resources["MyStateMachine"].(map[string]interface{})
	props := smResource["Properties"].(map[string]interface{})

	// Check Tags
	tags := props["Tags"].([]map[string]interface{})
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}

	// First tag should be SAM tag
	if tags[0]["Key"] != "stateMachine:createdBy" || tags[0]["Value"] != "SAM" {
		t.Errorf("expected first tag to be stateMachine:createdBy=SAM, got %v", tags[0])
	}

	// Check TracingConfiguration
	tracing := props["TracingConfiguration"].(map[string]interface{})
	if tracing["Enabled"] != true {
		t.Errorf("expected TracingConfiguration.Enabled to be true, got %v", tracing["Enabled"])
	}

	// Check that role has X-Ray policy
	roleResource := resources["MyStateMachineRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})
	managedPolicies := roleProps["ManagedPolicyArns"].([]interface{})
	found := false
	for _, arn := range managedPolicies {
		if arn == "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected X-Ray managed policy when tracing is enabled")
	}
}

func TestStateMachineTransformer_Transform_WithInlinePolicies(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States":  map[string]interface{}{},
		},
		Policies: []interface{}{
			map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []interface{}{
					map[string]interface{}{
						"Effect":   "Deny",
						"Action":   "*",
						"Resource": "*",
					},
				},
			},
		},
	}

	resources, err := transformer.Transform("MyStateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	roleResource := resources["MyStateMachineRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})

	// Check inline policies
	policies := roleProps["Policies"].([]map[string]interface{})
	if len(policies) != 1 {
		t.Fatalf("expected 1 inline policy, got %d", len(policies))
	}

	if policies[0]["PolicyName"] != "MyStateMachineRolePolicy0" {
		t.Errorf("expected PolicyName 'MyStateMachineRolePolicy0', got %v", policies[0]["PolicyName"])
	}
}

func TestStateMachineTransformer_Transform_WithScheduleEvent(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Role:          "arn:aws:iam::123456123456:role/service-role/SampleRole",
		DefinitionUri: "s3://sam-demo-bucket/my_state_machine.asl.json",
		Events: map[string]interface{}{
			"ScheduleEvent": map[string]interface{}{
				"Type": "Schedule",
				"Properties": map[string]interface{}{
					"Schedule":    "rate(1 minute)",
					"Name":        "TestSchedule",
					"Description": "test schedule",
					"Enabled":     false,
				},
			},
		},
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have 3 resources: StateMachine, Events::Rule, and invocation Role
	if len(resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(resources))
	}

	// Check Events::Rule
	ruleResource, ok := resources["StateMachineScheduleEvent"].(map[string]interface{})
	if !ok {
		t.Fatal("StateMachineScheduleEvent resource not found")
	}

	if ruleResource["Type"] != "AWS::Events::Rule" {
		t.Errorf("expected Type 'AWS::Events::Rule', got %v", ruleResource["Type"])
	}

	ruleProps := ruleResource["Properties"].(map[string]interface{})
	if ruleProps["ScheduleExpression"] != "rate(1 minute)" {
		t.Errorf("expected ScheduleExpression 'rate(1 minute)', got %v", ruleProps["ScheduleExpression"])
	}
	if ruleProps["Name"] != "TestSchedule" {
		t.Errorf("expected Name 'TestSchedule', got %v", ruleProps["Name"])
	}
	if ruleProps["Description"] != "test schedule" {
		t.Errorf("expected Description 'test schedule', got %v", ruleProps["Description"])
	}
	if ruleProps["State"] != "DISABLED" {
		t.Errorf("expected State 'DISABLED', got %v", ruleProps["State"])
	}

	// Check Targets
	targets := ruleProps["Targets"].([]interface{})
	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}

	target := targets[0].(map[string]interface{})
	if target["Id"] != "StateMachineScheduleEventStepFunctionsTarget" {
		t.Errorf("expected Target Id 'StateMachineScheduleEventStepFunctionsTarget', got %v", target["Id"])
	}

	// Check invocation role
	roleResource, ok := resources["StateMachineScheduleEventRole"].(map[string]interface{})
	if !ok {
		t.Fatal("StateMachineScheduleEventRole resource not found")
	}

	if roleResource["Type"] != "AWS::IAM::Role" {
		t.Errorf("expected Type 'AWS::IAM::Role', got %v", roleResource["Type"])
	}
}

func TestStateMachineTransformer_Transform_WithCloudWatchEvent(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Role:          "arn:aws:iam::123456123456:role/service-role/SampleRole",
		DefinitionUri: "s3://sam-demo-bucket/my_state_machine.asl.json",
		Events: map[string]interface{}{
			"CWEvent": map[string]interface{}{
				"Type": "CloudWatchEvent",
				"Properties": map[string]interface{}{
					"RuleName": "MyRule",
					"State":    "ENABLED",
					"Pattern": map[string]interface{}{
						"detail": map[string]interface{}{
							"state": []interface{}{"terminated"},
						},
					},
				},
			},
		},
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have 3 resources: StateMachine, Events::Rule, and invocation Role
	if len(resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(resources))
	}

	// Check Events::Rule
	ruleResource := resources["StateMachineCWEvent"].(map[string]interface{})
	ruleProps := ruleResource["Properties"].(map[string]interface{})

	if ruleProps["Name"] != "MyRule" {
		t.Errorf("expected Name 'MyRule', got %v", ruleProps["Name"])
	}
	if ruleProps["State"] != "ENABLED" {
		t.Errorf("expected State 'ENABLED', got %v", ruleProps["State"])
	}

	// Check EventPattern
	pattern := ruleProps["EventPattern"].(map[string]interface{})
	detail := pattern["detail"].(map[string]interface{})
	states := detail["state"].([]interface{})
	if len(states) != 1 || states[0] != "terminated" {
		t.Errorf("expected EventPattern detail.state to be ['terminated'], got %v", states)
	}
}

func TestStateMachineTransformer_Transform_WithScheduleV2Event(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Role:          "arn:aws:iam::123456123456:role/service-role/SampleRole",
		DefinitionUri: "s3://sam-demo-bucket/my_state_machine.asl.json",
		Events: map[string]interface{}{
			"ScheduleV2": map[string]interface{}{
				"Type": "ScheduleV2",
				"Properties": map[string]interface{}{
					"ScheduleExpression": "rate(1 minute)",
				},
			},
		},
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have 3 resources: StateMachine, Scheduler::Schedule, and invocation Role
	if len(resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(resources))
	}

	// Check Scheduler::Schedule
	scheduleResource := resources["StateMachineScheduleV2"].(map[string]interface{})
	if scheduleResource["Type"] != "AWS::Scheduler::Schedule" {
		t.Errorf("expected Type 'AWS::Scheduler::Schedule', got %v", scheduleResource["Type"])
	}

	scheduleProps := scheduleResource["Properties"].(map[string]interface{})
	if scheduleProps["ScheduleExpression"] != "rate(1 minute)" {
		t.Errorf("expected ScheduleExpression 'rate(1 minute)', got %v", scheduleProps["ScheduleExpression"])
	}

	// Should have default FlexibleTimeWindow
	ftw := scheduleProps["FlexibleTimeWindow"].(map[string]interface{})
	if ftw["Mode"] != "OFF" {
		t.Errorf("expected FlexibleTimeWindow.Mode 'OFF', got %v", ftw["Mode"])
	}

	// Check invocation role for scheduler.amazonaws.com
	roleResource := resources["StateMachineScheduleV2Role"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})
	assumeRole := roleProps["AssumeRolePolicyDocument"].(map[string]interface{})
	statements := assumeRole["Statement"].([]interface{})
	stmt := statements[0].(map[string]interface{})
	principal := stmt["Principal"].(map[string]interface{})
	if principal["Service"] != "scheduler.amazonaws.com" {
		t.Errorf("expected scheduler.amazonaws.com trust, got %v", principal["Service"])
	}
}

func TestStateMachineTransformer_Transform_DefinitionUriObject(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Role: "arn:aws:iam::123456123456:role/SampleRole",
		DefinitionUri: map[string]interface{}{
			"Bucket":  "my-bucket",
			"Key":     "path/to/definition.asl.json",
			"Version": "123",
		},
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	smResource := resources["StateMachine"].(map[string]interface{})
	props := smResource["Properties"].(map[string]interface{})

	s3Loc := props["DefinitionS3Location"].(map[string]interface{})
	if s3Loc["Bucket"] != "my-bucket" {
		t.Errorf("expected Bucket 'my-bucket', got %v", s3Loc["Bucket"])
	}
	if s3Loc["Key"] != "path/to/definition.asl.json" {
		t.Errorf("expected Key 'path/to/definition.asl.json', got %v", s3Loc["Key"])
	}
	if s3Loc["Version"] != "123" {
		t.Errorf("expected Version '123', got %v", s3Loc["Version"])
	}
}

func TestStateMachineTransformer_Transform_WithLogging(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Type: "EXPRESS",
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States":  map[string]interface{}{},
		},
		Logging: &LoggingConfig{
			Level:                "ALL",
			IncludeExecutionData: true,
			Destinations: []interface{}{
				map[string]interface{}{
					"CloudWatchLogsLogGroup": map[string]interface{}{
						"LogGroupArn": map[string]interface{}{
							"Fn::GetAtt": []string{"LogGroup", "Arn"},
						},
					},
				},
			},
		},
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	smResource := resources["StateMachine"].(map[string]interface{})
	props := smResource["Properties"].(map[string]interface{})

	if props["StateMachineType"] != "EXPRESS" {
		t.Errorf("expected StateMachineType 'EXPRESS', got %v", props["StateMachineType"])
	}

	loggingConfig := props["LoggingConfiguration"].(map[string]interface{})
	if loggingConfig["Level"] != "ALL" {
		t.Errorf("expected Level 'ALL', got %v", loggingConfig["Level"])
	}
	if loggingConfig["IncludeExecutionData"] != true {
		t.Errorf("expected IncludeExecutionData true, got %v", loggingConfig["IncludeExecutionData"])
	}
}

func TestStateMachineTransformer_Transform_WithRolePath(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States":  map[string]interface{}{},
		},
		RolePath: "/my/custom/path/",
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	roleResource := resources["StateMachineRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})

	if roleProps["Path"] != "/my/custom/path/" {
		t.Errorf("expected Path '/my/custom/path/', got %v", roleProps["Path"])
	}
}

func TestStateMachineTransformer_Transform_WithPermissionsBoundary(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States":  map[string]interface{}{},
		},
		PermissionsBoundary: "arn:aws:iam::123456789:policy/MyBoundary",
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	roleResource := resources["StateMachineRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})

	if roleProps["PermissionsBoundary"] != "arn:aws:iam::123456789:policy/MyBoundary" {
		t.Errorf("expected PermissionsBoundary 'arn:aws:iam::123456789:policy/MyBoundary', got %v", roleProps["PermissionsBoundary"])
	}
}

func TestStateMachineTransformer_Transform_WithDefinitionSubstitutions(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States":  map[string]interface{}{},
		},
		DefinitionSubstitutions: map[string]interface{}{
			"TableName": map[string]interface{}{"Ref": "MyTable"},
			"BucketArn": map[string]interface{}{"Fn::GetAtt": []string{"MyBucket", "Arn"}},
		},
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	smResource := resources["StateMachine"].(map[string]interface{})
	props := smResource["Properties"].(map[string]interface{})

	subs := props["DefinitionSubstitutions"].(map[string]interface{})
	if subs["TableName"] == nil || subs["BucketArn"] == nil {
		t.Error("expected DefinitionSubstitutions to be present")
	}
}

func TestStateMachineTransformer_parseDefinitionUri_InvalidFormat(t *testing.T) {
	transformer := NewStateMachineTransformer()

	// Test non-s3 URI
	_, err := transformer.parseDefinitionUri("https://example.com/definition.json")
	if err == nil {
		t.Error("expected error for non-s3 URI")
	}

	// Test invalid s3 URI (no key)
	_, err = transformer.parseDefinitionUri("s3://bucket-only")
	if err == nil {
		t.Error("expected error for s3 URI without key")
	}
}

func TestStateMachineTransformer_convertDefinitionToString(t *testing.T) {
	transformer := NewStateMachineTransformer()

	definition := map[string]interface{}{
		"StartAt": "Hello",
		"States": map[string]interface{}{
			"Hello": map[string]interface{}{
				"Type": "Pass",
				"End":  true,
			},
		},
	}

	result, err := transformer.convertDefinitionToString(definition)
	if err != nil {
		t.Fatalf("convertDefinitionToString failed: %v", err)
	}

	fnJoin := result.(map[string]interface{})["Fn::Join"].([]interface{})
	if fnJoin[0] != "\n" {
		t.Errorf("expected Fn::Join separator to be newline, got %v", fnJoin[0])
	}

	lines := fnJoin[1].([]interface{})
	if len(lines) == 0 {
		t.Error("expected at least one line in Fn::Join array")
	}

	// Verify it can be parsed back as JSON when joined
	joined := ""
	for _, line := range lines {
		if joined != "" {
			joined += "\n"
		}
		joined += line.(string)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(joined), &parsed); err != nil {
		t.Errorf("joined definition should be valid JSON: %v", err)
	}
}

func TestStateMachineTransformer_Transform_ManagedPolicyString(t *testing.T) {
	transformer := NewStateMachineTransformer()

	sm := &StateMachine{
		Definition: map[string]interface{}{
			"StartAt": "Hello",
			"States":  map[string]interface{}{},
		},
		Policies: "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
	}

	resources, err := transformer.Transform("StateMachine", sm)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	roleResource := resources["StateMachineRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})

	managedPolicies := roleProps["ManagedPolicyArns"].([]interface{})
	found := false
	for _, arn := range managedPolicies {
		if arn == "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected managed policy ARN to be added")
	}
}
