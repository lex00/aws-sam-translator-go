package sqs

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestQueue_JSONSerialization(t *testing.T) {
	queue := Queue{
		QueueName:              "MyQueue",
		VisibilityTimeout:      30,
		MessageRetentionPeriod: 345600, // 4 days
		DelaySeconds:           0,
		MaximumMessageSize:     262144,
		Tags: []Tag{
			{Key: "Environment", Value: "Production"},
		},
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal queue to JSON: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal queue from JSON: %v", err)
	}

	if unmarshaled.QueueName != queue.QueueName {
		t.Errorf("QueueName mismatch: got %v, want %v", unmarshaled.QueueName, queue.QueueName)
	}
}

func TestQueue_YAMLSerialization(t *testing.T) {
	queue := Queue{
		QueueName:         "MyQueue",
		VisibilityTimeout: 30,
	}

	data, err := yaml.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal queue to YAML: %v", err)
	}

	var unmarshaled Queue
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal queue from YAML: %v", err)
	}

	if unmarshaled.QueueName != queue.QueueName {
		t.Errorf("QueueName mismatch: got %v, want %v", unmarshaled.QueueName, queue.QueueName)
	}
}

func TestQueue_WithIntrinsicFunctions(t *testing.T) {
	queue := Queue{
		QueueName: map[string]interface{}{
			"Fn::Sub": "${AWS::StackName}-queue",
		},
		KmsMasterKeyId: map[string]interface{}{
			"Fn::GetAtt": []string{"MyKey", "Arn"},
		},
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal queue with intrinsics: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal queue with intrinsics: %v", err)
	}
}

func TestQueue_FifoQueue(t *testing.T) {
	queue := Queue{
		QueueName:                 "MyQueue.fifo",
		FifoQueue:                 true,
		ContentBasedDeduplication: true,
		DeduplicationScope:        "messageGroup",
		FifoThroughputLimit:       "perMessageGroupId",
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal FIFO queue: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal FIFO queue: %v", err)
	}

	if unmarshaled.FifoQueue != true {
		t.Errorf("FifoQueue mismatch: got %v, want true", unmarshaled.FifoQueue)
	}

	if unmarshaled.DeduplicationScope != queue.DeduplicationScope {
		t.Errorf("DeduplicationScope mismatch: got %v, want %v", unmarshaled.DeduplicationScope, queue.DeduplicationScope)
	}
}

func TestQueue_WithRedrivePolicy(t *testing.T) {
	queue := Queue{
		QueueName: "MyQueue",
		RedrivePolicy: map[string]interface{}{
			"deadLetterTargetArn": "arn:aws:sqs:us-east-1:123456789012:MyDLQ",
			"maxReceiveCount":     5,
		},
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal queue with redrive policy: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal queue with redrive policy: %v", err)
	}

	if unmarshaled.RedrivePolicy == nil {
		t.Error("RedrivePolicy should not be nil")
	}
}

func TestQueue_WithRedriveAllowPolicy(t *testing.T) {
	queue := Queue{
		QueueName: "MyDLQ",
		RedriveAllowPolicy: map[string]interface{}{
			"redrivePermission": "byQueue",
			"sourceQueueArns": []string{
				"arn:aws:sqs:us-east-1:123456789012:Queue1",
				"arn:aws:sqs:us-east-1:123456789012:Queue2",
			},
		},
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal queue with redrive allow policy: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal queue with redrive allow policy: %v", err)
	}

	if unmarshaled.RedriveAllowPolicy == nil {
		t.Error("RedriveAllowPolicy should not be nil")
	}
}

func TestQueue_WithEncryption(t *testing.T) {
	queue := Queue{
		QueueName:                    "MyEncryptedQueue",
		KmsMasterKeyId:               "alias/my-key",
		KmsDataKeyReusePeriodSeconds: 300,
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal encrypted queue: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal encrypted queue: %v", err)
	}

	if unmarshaled.KmsMasterKeyId != queue.KmsMasterKeyId {
		t.Errorf("KmsMasterKeyId mismatch: got %v, want %v", unmarshaled.KmsMasterKeyId, queue.KmsMasterKeyId)
	}
}

func TestQueue_WithSqsManagedEncryption(t *testing.T) {
	queue := Queue{
		QueueName:            "MySseManagedQueue",
		SqsManagedSseEnabled: true,
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal SQS managed encrypted queue: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SQS managed encrypted queue: %v", err)
	}

	if unmarshaled.SqsManagedSseEnabled != true {
		t.Errorf("SqsManagedSseEnabled mismatch: got %v, want true", unmarshaled.SqsManagedSseEnabled)
	}
}

func TestQueue_WithLongPolling(t *testing.T) {
	queue := Queue{
		QueueName:                     "MyLongPollingQueue",
		ReceiveMessageWaitTimeSeconds: 20,
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal long polling queue: %v", err)
	}

	var unmarshaled Queue
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal long polling queue: %v", err)
	}

	// Compare as strings since JSON unmarshaling converts numbers to float64
	if fmt.Sprintf("%v", unmarshaled.ReceiveMessageWaitTimeSeconds) != fmt.Sprintf("%v", queue.ReceiveMessageWaitTimeSeconds) {
		t.Errorf("ReceiveMessageWaitTimeSeconds mismatch: got %v, want %v", unmarshaled.ReceiveMessageWaitTimeSeconds, queue.ReceiveMessageWaitTimeSeconds)
	}
}

func TestQueue_OmitEmpty(t *testing.T) {
	queue := Queue{
		QueueName: "MyQueue",
	}

	data, err := json.Marshal(queue)
	if err != nil {
		t.Fatalf("Failed to marshal queue: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	omittedFields := []string{
		"ContentBasedDeduplication",
		"DeduplicationScope",
		"DelaySeconds",
		"FifoQueue",
		"FifoThroughputLimit",
		"KmsMasterKeyId",
		"MaximumMessageSize",
		"MessageRetentionPeriod",
		"ReceiveMessageWaitTimeSeconds",
		"RedriveAllowPolicy",
		"RedrivePolicy",
		"SqsManagedSseEnabled",
		"Tags",
		"VisibilityTimeout",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestRedrivePolicy_JSONSerialization(t *testing.T) {
	policy := RedrivePolicy{
		DeadLetterTargetArn: "arn:aws:sqs:us-east-1:123456789012:MyDLQ",
		MaxReceiveCount:     5,
	}

	data, err := json.Marshal(policy)
	if err != nil {
		t.Fatalf("Failed to marshal RedrivePolicy: %v", err)
	}

	var unmarshaled RedrivePolicy
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RedrivePolicy: %v", err)
	}

	if unmarshaled.DeadLetterTargetArn != policy.DeadLetterTargetArn {
		t.Errorf("DeadLetterTargetArn mismatch: got %v, want %v", unmarshaled.DeadLetterTargetArn, policy.DeadLetterTargetArn)
	}
}

func TestRedriveAllowPolicy_JSONSerialization(t *testing.T) {
	policy := RedriveAllowPolicy{
		RedrivePermission: "byQueue",
		SourceQueueArns:   []interface{}{"arn:aws:sqs:us-east-1:123456789012:Queue1"},
	}

	data, err := json.Marshal(policy)
	if err != nil {
		t.Fatalf("Failed to marshal RedriveAllowPolicy: %v", err)
	}

	var unmarshaled RedriveAllowPolicy
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal RedriveAllowPolicy: %v", err)
	}

	if unmarshaled.RedrivePermission != policy.RedrivePermission {
		t.Errorf("RedrivePermission mismatch: got %v, want %v", unmarshaled.RedrivePermission, policy.RedrivePermission)
	}
}
