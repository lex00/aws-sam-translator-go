package push

import (
	"testing"
)

func TestSNSEventSourceHandler_GenerateResources(t *testing.T) {
	handler := NewSNSEventSourceHandler()

	tests := []struct {
		name              string
		functionLogicalID string
		eventLogicalID    string
		event             *SNSEvent
		wantErr           bool
		validate          func(t *testing.T, resources map[string]interface{})
	}{
		{
			name:              "basic SNS event with topic ARN",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "MySNSTopic",
			event: &SNSEvent{
				Topic: "arn:aws:sns:us-east-1:123456789012:my-topic",
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				// Check that two resources are created
				if len(resources) != 2 {
					t.Errorf("expected 2 resources, got %d", len(resources))
				}

				// Check subscription resource
				subResource, ok := resources["MyFunctionMySNSTopic"]
				if !ok {
					t.Fatal("subscription resource not found")
				}

				subMap, ok := subResource.(map[string]interface{})
				if !ok {
					t.Fatal("subscription resource is not a map")
				}

				if subMap["Type"] != "AWS::SNS::Subscription" {
					t.Errorf("expected Type AWS::SNS::Subscription, got %v", subMap["Type"])
				}

				props, ok := subMap["Properties"].(map[string]interface{})
				if !ok {
					t.Fatal("subscription properties is not a map")
				}

				if props["Protocol"] != "lambda" {
					t.Errorf("expected Protocol lambda, got %v", props["Protocol"])
				}

				if props["TopicArn"] != "arn:aws:sns:us-east-1:123456789012:my-topic" {
					t.Errorf("expected TopicArn arn:aws:sns:us-east-1:123456789012:my-topic, got %v", props["TopicArn"])
				}

				endpoint, ok := props["Endpoint"].(map[string]interface{})
				if !ok {
					t.Fatal("endpoint is not a map")
				}

				getAtt, ok := endpoint["Fn::GetAtt"].([]interface{})
				if !ok {
					t.Fatal("Fn::GetAtt is not an array")
				}

				if len(getAtt) != 2 || getAtt[0] != "MyFunction" || getAtt[1] != "Arn" {
					t.Errorf("expected Fn::GetAtt [MyFunction, Arn], got %v", getAtt)
				}

				// Check permission resource
				permResource, ok := resources["MyFunctionMySNSTopicPermission"]
				if !ok {
					t.Fatal("permission resource not found")
				}

				permMap, ok := permResource.(map[string]interface{})
				if !ok {
					t.Fatal("permission resource is not a map")
				}

				if permMap["Type"] != "AWS::Lambda::Permission" {
					t.Errorf("expected Type AWS::Lambda::Permission, got %v", permMap["Type"])
				}

				permProps, ok := permMap["Properties"].(map[string]interface{})
				if !ok {
					t.Fatal("permission properties is not a map")
				}

				if permProps["Action"] != "lambda:InvokeFunction" {
					t.Errorf("expected Action lambda:InvokeFunction, got %v", permProps["Action"])
				}

				if permProps["Principal"] != "sns.amazonaws.com" {
					t.Errorf("expected Principal sns.amazonaws.com, got %v", permProps["Principal"])
				}

				if permProps["SourceArn"] != "arn:aws:sns:us-east-1:123456789012:my-topic" {
					t.Errorf("expected SourceArn arn:aws:sns:us-east-1:123456789012:my-topic, got %v", permProps["SourceArn"])
				}
			},
		},
		{
			name:              "SNS event with topic reference",
			functionLogicalID: "SaveNotificationFunction",
			eventLogicalID:    "NotificationTopic",
			event: &SNSEvent{
				Topic: map[string]interface{}{
					"Ref": "SNSTopicArn",
				},
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				if len(resources) != 2 {
					t.Errorf("expected 2 resources, got %d", len(resources))
				}

				subResource := resources["SaveNotificationFunctionNotificationTopic"].(map[string]interface{})
				props := subResource["Properties"].(map[string]interface{})

				topicArn, ok := props["TopicArn"].(map[string]interface{})
				if !ok {
					t.Fatal("TopicArn is not a map")
				}

				if topicArn["Ref"] != "SNSTopicArn" {
					t.Errorf("expected TopicArn Ref SNSTopicArn, got %v", topicArn["Ref"])
				}
			},
		},
		{
			name:              "SNS event with all parameters",
			functionLogicalID: "MyAwesomeFunction",
			eventLogicalID:    "NotificationTopic",
			event: &SNSEvent{
				Topic:  "arn:aws:sns:us-west-2:987654321098:key/dec86919-7219-4e8d-8871-7f1609df2c7f",
				Region: "region",
				FilterPolicy: map[string]interface{}{
					"store": []interface{}{"example_corp"},
					"event": []interface{}{
						map[string]interface{}{
							"anything-but": "order_cancelled",
						},
					},
					"customer_interests": []interface{}{"rugby", "football", "baseball"},
					"price_usd": []interface{}{
						map[string]interface{}{
							"numeric": []interface{}{">=", 100},
						},
					},
					"before": map[string]interface{}{
						"owner": []interface{}{"0x0"},
					},
				},
				FilterPolicyScope: "MessageAttributes",
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				if len(resources) != 2 {
					t.Errorf("expected 2 resources, got %d", len(resources))
				}

				subResource := resources["MyAwesomeFunctionNotificationTopic"].(map[string]interface{})
				props := subResource["Properties"].(map[string]interface{})

				if props["Region"] != "region" {
					t.Errorf("expected Region region, got %v", props["Region"])
				}

				if props["FilterPolicyScope"] != "MessageAttributes" {
					t.Errorf("expected FilterPolicyScope MessageAttributes, got %v", props["FilterPolicyScope"])
				}

				filterPolicy, ok := props["FilterPolicy"].(map[string]interface{})
				if !ok {
					t.Fatal("FilterPolicy is not a map")
				}

				if filterPolicy["store"] == nil {
					t.Error("expected FilterPolicy to have store key")
				}

				if filterPolicy["event"] == nil {
					t.Error("expected FilterPolicy to have event key")
				}

				if filterPolicy["customer_interests"] == nil {
					t.Error("expected FilterPolicy to have customer_interests key")
				}

				if filterPolicy["price_usd"] == nil {
					t.Error("expected FilterPolicy to have price_usd key")
				}

				if filterPolicy["before"] == nil {
					t.Error("expected FilterPolicy to have before key")
				}
			},
		},
		{
			name:              "SNS event with redrive policy",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "MySNSTopic",
			event: &SNSEvent{
				Topic: "arn:aws:sns:us-east-1:123456789012:my-topic",
				RedrivePolicy: map[string]interface{}{
					"deadLetterTargetArn": "arn:aws:sqs:us-east-1:123456789012:my-dlq",
				},
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				if len(resources) != 2 {
					t.Errorf("expected 2 resources, got %d", len(resources))
				}

				subResource := resources["MyFunctionMySNSTopic"].(map[string]interface{})
				props := subResource["Properties"].(map[string]interface{})

				redrivePolicy, ok := props["RedrivePolicy"].(map[string]interface{})
				if !ok {
					t.Fatal("RedrivePolicy is not a map")
				}

				if redrivePolicy["deadLetterTargetArn"] != "arn:aws:sqs:us-east-1:123456789012:my-dlq" {
					t.Errorf("expected deadLetterTargetArn arn:aws:sqs:us-east-1:123456789012:my-dlq, got %v", redrivePolicy["deadLetterTargetArn"])
				}
			},
		},
		{
			name:              "missing topic should error",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "MySNSTopic",
			event:             &SNSEvent{},
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := handler.GenerateResources(tt.functionLogicalID, tt.eventLogicalID, tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, resources)
			}
		})
	}
}

func TestNewSNSEventSourceHandler(t *testing.T) {
	handler := NewSNSEventSourceHandler()
	if handler == nil {
		t.Error("expected non-nil handler")
	}
}

func TestSNSEvent_Properties(t *testing.T) {
	// Test that SNSEvent struct can be created with all properties
	event := &SNSEvent{
		Topic:             "arn:aws:sns:us-east-1:123456789012:my-topic",
		Region:            "us-east-1",
		FilterPolicy:      map[string]interface{}{"key": "value"},
		FilterPolicyScope: "MessageAttributes",
		SqsSubscription:   false,
		RedrivePolicy:     map[string]interface{}{"deadLetterTargetArn": "arn"},
	}

	if event.Topic != "arn:aws:sns:us-east-1:123456789012:my-topic" {
		t.Errorf("expected Topic to be set")
	}
	if event.Region != "us-east-1" {
		t.Errorf("expected Region to be set")
	}
	if event.FilterPolicyScope != "MessageAttributes" {
		t.Errorf("expected FilterPolicyScope to be set")
	}
}
