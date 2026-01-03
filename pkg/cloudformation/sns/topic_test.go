package sns

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTopic_JSONSerialization(t *testing.T) {
	topic := Topic{
		TopicName:   "MyTopic",
		DisplayName: "My Topic Display Name",
		Subscription: []TopicSubscription{
			{Endpoint: "user@example.com", Protocol: "email"},
			{Endpoint: "arn:aws:sqs:us-east-1:123456789012:MyQueue", Protocol: "sqs"},
		},
		Tags: []Tag{
			{Key: "Environment", Value: "Production"},
		},
	}

	data, err := json.Marshal(topic)
	if err != nil {
		t.Fatalf("Failed to marshal topic to JSON: %v", err)
	}

	var unmarshaled Topic
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal topic from JSON: %v", err)
	}

	if unmarshaled.TopicName != topic.TopicName {
		t.Errorf("TopicName mismatch: got %v, want %v", unmarshaled.TopicName, topic.TopicName)
	}

	if len(unmarshaled.Subscription) != len(topic.Subscription) {
		t.Errorf("Subscription length mismatch: got %d, want %d", len(unmarshaled.Subscription), len(topic.Subscription))
	}
}

func TestTopic_YAMLSerialization(t *testing.T) {
	topic := Topic{
		TopicName:   "MyTopic",
		DisplayName: "My Topic",
	}

	data, err := yaml.Marshal(topic)
	if err != nil {
		t.Fatalf("Failed to marshal topic to YAML: %v", err)
	}

	var unmarshaled Topic
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal topic from YAML: %v", err)
	}

	if unmarshaled.TopicName != topic.TopicName {
		t.Errorf("TopicName mismatch: got %v, want %v", unmarshaled.TopicName, topic.TopicName)
	}
}

func TestTopic_WithIntrinsicFunctions(t *testing.T) {
	topic := Topic{
		TopicName: map[string]interface{}{
			"Fn::Sub": "${AWS::StackName}-topic",
		},
		KmsMasterKeyId: map[string]interface{}{
			"Fn::GetAtt": []string{"MyKey", "Arn"},
		},
	}

	data, err := json.Marshal(topic)
	if err != nil {
		t.Fatalf("Failed to marshal topic with intrinsics: %v", err)
	}

	var unmarshaled Topic
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal topic with intrinsics: %v", err)
	}
}

func TestTopic_FifoTopic(t *testing.T) {
	topic := Topic{
		TopicName:                 "MyTopic.fifo",
		FifoTopic:                 true,
		ContentBasedDeduplication: true,
	}

	data, err := json.Marshal(topic)
	if err != nil {
		t.Fatalf("Failed to marshal FIFO topic: %v", err)
	}

	var unmarshaled Topic
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal FIFO topic: %v", err)
	}

	if unmarshaled.FifoTopic != true {
		t.Errorf("FifoTopic mismatch: got %v, want true", unmarshaled.FifoTopic)
	}
}

func TestTopic_WithDeliveryStatusLogging(t *testing.T) {
	topic := Topic{
		TopicName: "MyTopic",
		DeliveryStatusLogging: []LoggingConfig{
			{
				Protocol:                  "sqs",
				SuccessFeedbackRoleArn:    "arn:aws:iam::123456789012:role/LoggingRole",
				SuccessFeedbackSampleRate: "100",
				FailureFeedbackRoleArn:    "arn:aws:iam::123456789012:role/LoggingRole",
			},
		},
	}

	data, err := json.Marshal(topic)
	if err != nil {
		t.Fatalf("Failed to marshal topic with logging: %v", err)
	}

	var unmarshaled Topic
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal topic with logging: %v", err)
	}

	if len(unmarshaled.DeliveryStatusLogging) != 1 {
		t.Errorf("DeliveryStatusLogging length mismatch: got %d, want 1", len(unmarshaled.DeliveryStatusLogging))
	}
}

func TestTopic_OmitEmpty(t *testing.T) {
	topic := Topic{
		TopicName: "MyTopic",
	}

	data, err := json.Marshal(topic)
	if err != nil {
		t.Fatalf("Failed to marshal topic: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	omittedFields := []string{
		"DisplayName",
		"FifoTopic",
		"ContentBasedDeduplication",
		"KmsMasterKeyId",
		"Subscription",
		"Tags",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestSubscription_JSONSerialization(t *testing.T) {
	sub := Subscription{
		TopicArn:           "arn:aws:sns:us-east-1:123456789012:MyTopic",
		Protocol:           "sqs",
		Endpoint:           "arn:aws:sqs:us-east-1:123456789012:MyQueue",
		RawMessageDelivery: true,
		FilterPolicy: map[string]interface{}{
			"eventType": []string{"order_created", "order_updated"},
		},
		FilterPolicyScope: "MessageAttributes",
	}

	data, err := json.Marshal(sub)
	if err != nil {
		t.Fatalf("Failed to marshal subscription to JSON: %v", err)
	}

	var unmarshaled Subscription
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal subscription from JSON: %v", err)
	}

	if unmarshaled.TopicArn != sub.TopicArn {
		t.Errorf("TopicArn mismatch: got %v, want %v", unmarshaled.TopicArn, sub.TopicArn)
	}

	if unmarshaled.Protocol != sub.Protocol {
		t.Errorf("Protocol mismatch: got %v, want %v", unmarshaled.Protocol, sub.Protocol)
	}
}

func TestSubscription_YAMLSerialization(t *testing.T) {
	sub := Subscription{
		TopicArn: "arn:aws:sns:us-east-1:123456789012:MyTopic",
		Protocol: "lambda",
		Endpoint: "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
	}

	data, err := yaml.Marshal(sub)
	if err != nil {
		t.Fatalf("Failed to marshal subscription to YAML: %v", err)
	}

	var unmarshaled Subscription
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal subscription from YAML: %v", err)
	}

	if unmarshaled.Protocol != sub.Protocol {
		t.Errorf("Protocol mismatch: got %v, want %v", unmarshaled.Protocol, sub.Protocol)
	}
}

func TestSubscription_WithIntrinsicFunctions(t *testing.T) {
	sub := Subscription{
		TopicArn: map[string]interface{}{
			"Ref": "MyTopic",
		},
		Protocol: "sqs",
		Endpoint: map[string]interface{}{
			"Fn::GetAtt": []string{"MyQueue", "Arn"},
		},
	}

	data, err := json.Marshal(sub)
	if err != nil {
		t.Fatalf("Failed to marshal subscription with intrinsics: %v", err)
	}

	var unmarshaled Subscription
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal subscription with intrinsics: %v", err)
	}
}

func TestSubscription_WithDeliveryPolicy(t *testing.T) {
	sub := Subscription{
		TopicArn: "arn:aws:sns:us-east-1:123456789012:MyTopic",
		Protocol: "http",
		Endpoint: "https://example.com/webhook",
		DeliveryPolicy: map[string]interface{}{
			"healthyRetryPolicy": map[string]interface{}{
				"numRetries":         3,
				"minDelayTarget":     20,
				"maxDelayTarget":     20,
				"numMaxDelayRetries": 0,
				"backoffFunction":    "linear",
			},
		},
	}

	data, err := json.Marshal(sub)
	if err != nil {
		t.Fatalf("Failed to marshal subscription with delivery policy: %v", err)
	}

	var unmarshaled Subscription
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal subscription with delivery policy: %v", err)
	}

	if unmarshaled.DeliveryPolicy == nil {
		t.Error("DeliveryPolicy should not be nil")
	}
}

func TestSubscription_WithRedrivePolicy(t *testing.T) {
	sub := Subscription{
		TopicArn: "arn:aws:sns:us-east-1:123456789012:MyTopic",
		Protocol: "sqs",
		Endpoint: "arn:aws:sqs:us-east-1:123456789012:MyQueue",
		RedrivePolicy: map[string]interface{}{
			"deadLetterTargetArn": "arn:aws:sqs:us-east-1:123456789012:MyDLQ",
		},
	}

	data, err := json.Marshal(sub)
	if err != nil {
		t.Fatalf("Failed to marshal subscription with redrive policy: %v", err)
	}

	var unmarshaled Subscription
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal subscription with redrive policy: %v", err)
	}

	if unmarshaled.RedrivePolicy == nil {
		t.Error("RedrivePolicy should not be nil")
	}
}

func TestSubscription_OmitEmpty(t *testing.T) {
	sub := Subscription{
		TopicArn: "arn:aws:sns:us-east-1:123456789012:MyTopic",
		Protocol: "email",
	}

	data, err := json.Marshal(sub)
	if err != nil {
		t.Fatalf("Failed to marshal subscription: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	omittedFields := []string{
		"Endpoint",
		"DeliveryPolicy",
		"FilterPolicy",
		"FilterPolicyScope",
		"RawMessageDelivery",
		"RedrivePolicy",
		"Region",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}
