package events

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRule_JSONSerialization(t *testing.T) {
	rule := Rule{
		Name:               "MyRule",
		Description:        "A test EventBridge rule",
		EventBusName:       "default",
		ScheduleExpression: "rate(5 minutes)",
		State:              "ENABLED",
		Targets: []Target{
			{
				Id:  "Target1",
				Arn: "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
			},
		},
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Failed to marshal rule to JSON: %v", err)
	}

	var unmarshaled Rule
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal rule from JSON: %v", err)
	}

	if unmarshaled.Name != rule.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, rule.Name)
	}

	if unmarshaled.ScheduleExpression != rule.ScheduleExpression {
		t.Errorf("ScheduleExpression mismatch: got %v, want %v", unmarshaled.ScheduleExpression, rule.ScheduleExpression)
	}

	if len(unmarshaled.Targets) != len(rule.Targets) {
		t.Errorf("Targets length mismatch: got %d, want %d", len(unmarshaled.Targets), len(rule.Targets))
	}
}

func TestRule_YAMLSerialization(t *testing.T) {
	rule := Rule{
		Name: "MyRule",
		EventPattern: map[string]interface{}{
			"source":      []string{"aws.ec2"},
			"detail-type": []string{"EC2 Instance State-change Notification"},
		},
		State: "ENABLED",
	}

	data, err := yaml.Marshal(rule)
	if err != nil {
		t.Fatalf("Failed to marshal rule to YAML: %v", err)
	}

	var unmarshaled Rule
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal rule from YAML: %v", err)
	}

	if unmarshaled.Name != rule.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, rule.Name)
	}
}

func TestRule_WithIntrinsicFunctions(t *testing.T) {
	rule := Rule{
		Name: map[string]interface{}{
			"Fn::Sub": "${AWS::StackName}-rule",
		},
		RoleArn: map[string]interface{}{
			"Fn::GetAtt": []string{"MyRole", "Arn"},
		},
		Targets: []Target{
			{
				Id: "Target1",
				Arn: map[string]interface{}{
					"Fn::GetAtt": []string{"MyFunction", "Arn"},
				},
			},
		},
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Failed to marshal rule with intrinsics: %v", err)
	}

	var unmarshaled Rule
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal rule with intrinsics: %v", err)
	}
}

func TestRule_WithEventPattern(t *testing.T) {
	rule := Rule{
		Name: "PatternRule",
		EventPattern: map[string]interface{}{
			"source":      []string{"custom.myapp"},
			"detail-type": []string{"UserCreated"},
			"detail": map[string]interface{}{
				"userType": []string{"admin", "user"},
			},
		},
		State: "ENABLED",
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Failed to marshal rule with event pattern: %v", err)
	}

	var unmarshaled Rule
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal rule with event pattern: %v", err)
	}
}

func TestTarget_WithInputTransformer(t *testing.T) {
	target := Target{
		Id:  "TransformedTarget",
		Arn: "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
		InputTransformer: &InputTransformer{
			InputPathsMap: map[string]interface{}{
				"instance": "$.detail.instance-id",
				"state":    "$.detail.state",
			},
			InputTemplate: `{"instanceId": <instance>, "newState": <state>}`,
		},
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("Failed to marshal target with input transformer: %v", err)
	}

	var unmarshaled Target
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal target with input transformer: %v", err)
	}

	if unmarshaled.InputTransformer == nil {
		t.Error("InputTransformer should not be nil")
	}
}

func TestTarget_WithEcsParameters(t *testing.T) {
	target := Target{
		Id:  "EcsTarget",
		Arn: "arn:aws:ecs:us-east-1:123456789012:cluster/MyCluster",
		EcsParameters: &EcsParameters{
			TaskDefinitionArn: "arn:aws:ecs:us-east-1:123456789012:task-definition/MyTask:1",
			TaskCount:         1,
			LaunchType:        "FARGATE",
			NetworkConfiguration: &NetworkConfiguration{
				AwsVpcConfiguration: &AwsVpcConfiguration{
					Subnets:        []interface{}{"subnet-12345678"},
					SecurityGroups: []interface{}{"sg-12345678"},
					AssignPublicIp: "ENABLED",
				},
			},
		},
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("Failed to marshal target with ECS parameters: %v", err)
	}

	var unmarshaled Target
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal target with ECS parameters: %v", err)
	}

	if unmarshaled.EcsParameters == nil {
		t.Error("EcsParameters should not be nil")
	}

	if unmarshaled.EcsParameters.NetworkConfiguration == nil {
		t.Error("NetworkConfiguration should not be nil")
	}
}

func TestTarget_WithDeadLetterConfig(t *testing.T) {
	target := Target{
		Id:  "DLQTarget",
		Arn: "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
		DeadLetterConfig: &DeadLetterConfig{
			Arn: "arn:aws:sqs:us-east-1:123456789012:MyDLQ",
		},
		RetryPolicy: &RetryPolicy{
			MaximumEventAgeInSeconds: 3600,
			MaximumRetryAttempts:     3,
		},
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("Failed to marshal target with DLQ config: %v", err)
	}

	var unmarshaled Target
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal target with DLQ config: %v", err)
	}

	if unmarshaled.DeadLetterConfig == nil {
		t.Error("DeadLetterConfig should not be nil")
	}

	if unmarshaled.RetryPolicy == nil {
		t.Error("RetryPolicy should not be nil")
	}
}

func TestTarget_WithBatchParameters(t *testing.T) {
	target := Target{
		Id:  "BatchTarget",
		Arn: "arn:aws:batch:us-east-1:123456789012:job-queue/MyQueue",
		BatchParameters: &BatchParameters{
			JobDefinition: "arn:aws:batch:us-east-1:123456789012:job-definition/MyJob:1",
			JobName:       "MyBatchJob",
			ArrayProperties: &BatchArrayProperties{
				Size: 10,
			},
			RetryStrategy: &BatchRetryStrategy{
				Attempts: 3,
			},
		},
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("Failed to marshal target with Batch parameters: %v", err)
	}

	var unmarshaled Target
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal target with Batch parameters: %v", err)
	}

	if unmarshaled.BatchParameters == nil {
		t.Error("BatchParameters should not be nil")
	}
}

func TestTarget_WithHttpParameters(t *testing.T) {
	target := Target{
		Id:  "ApiTarget",
		Arn: "arn:aws:execute-api:us-east-1:123456789012:api-id/stage/method/path",
		HttpParameters: &HttpParameters{
			PathParameterValues: []interface{}{"value1", "value2"},
			HeaderParameters: map[string]interface{}{
				"Content-Type": "application/json",
			},
			QueryStringParameters: map[string]interface{}{
				"key": "value",
			},
		},
	}

	data, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("Failed to marshal target with HTTP parameters: %v", err)
	}

	var unmarshaled Target
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal target with HTTP parameters: %v", err)
	}

	if unmarshaled.HttpParameters == nil {
		t.Error("HttpParameters should not be nil")
	}
}

func TestRule_OmitEmpty(t *testing.T) {
	rule := Rule{
		Name:  "MinimalRule",
		State: "ENABLED",
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Failed to marshal rule: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	omittedFields := []string{
		"Description",
		"EventBusName",
		"EventPattern",
		"ScheduleExpression",
		"Targets",
		"RoleArn",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestSqsParameters_JSONSerialization(t *testing.T) {
	params := SqsParameters{
		MessageGroupId: "my-group-id",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal SqsParameters: %v", err)
	}

	var unmarshaled SqsParameters
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SqsParameters: %v", err)
	}

	if unmarshaled.MessageGroupId != params.MessageGroupId {
		t.Errorf("MessageGroupId mismatch: got %v, want %v", unmarshaled.MessageGroupId, params.MessageGroupId)
	}
}

func TestKinesisParameters_JSONSerialization(t *testing.T) {
	params := KinesisParameters{
		PartitionKeyPath: "$.detail.id",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal KinesisParameters: %v", err)
	}

	var unmarshaled KinesisParameters
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal KinesisParameters: %v", err)
	}

	if unmarshaled.PartitionKeyPath != params.PartitionKeyPath {
		t.Errorf("PartitionKeyPath mismatch: got %v, want %v", unmarshaled.PartitionKeyPath, params.PartitionKeyPath)
	}
}
