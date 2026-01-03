package push

import (
	"testing"
)

func TestIoTRuleEventSourceHandler_GenerateResources(t *testing.T) {
	handler := NewIoTRuleEventSourceHandler()

	tests := []struct {
		name              string
		functionLogicalID string
		eventLogicalID    string
		event             *IoTRuleEvent
		wantErr           bool
		validate          func(t *testing.T, resources map[string]interface{})
	}{
		{
			name:              "basic IoTRule event",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "MyIoTRule",
			event: &IoTRuleEvent{
				Sql: "SELECT * FROM 'my/topic'",
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				// Check that two resources are created
				if len(resources) != 2 {
					t.Errorf("expected 2 resources, got %d", len(resources))
				}

				// Check topic rule resource
				ruleResource, ok := resources["MyFunctionMyIoTRule"]
				if !ok {
					t.Fatal("topic rule resource not found")
				}

				ruleMap, ok := ruleResource.(map[string]interface{})
				if !ok {
					t.Fatal("topic rule resource is not a map")
				}

				if ruleMap["Type"] != "AWS::IoT::TopicRule" {
					t.Errorf("expected Type AWS::IoT::TopicRule, got %v", ruleMap["Type"])
				}

				props, ok := ruleMap["Properties"].(map[string]interface{})
				if !ok {
					t.Fatal("topic rule properties is not a map")
				}

				payload, ok := props["TopicRulePayload"].(map[string]interface{})
				if !ok {
					t.Fatal("TopicRulePayload is not a map")
				}

				if payload["Sql"] != "SELECT * FROM 'my/topic'" {
					t.Errorf("expected Sql SELECT * FROM 'my/topic', got %v", payload["Sql"])
				}

				actions, ok := payload["Actions"].([]interface{})
				if !ok {
					t.Fatal("Actions is not an array")
				}

				if len(actions) != 1 {
					t.Fatalf("expected 1 action, got %d", len(actions))
				}

				action, ok := actions[0].(map[string]interface{})
				if !ok {
					t.Fatal("action is not a map")
				}

				lambdaAction, ok := action["Lambda"].(map[string]interface{})
				if !ok {
					t.Fatal("Lambda action is not a map")
				}

				functionArn, ok := lambdaAction["FunctionArn"].(map[string]interface{})
				if !ok {
					t.Fatal("FunctionArn is not a map")
				}

				getAtt, ok := functionArn["Fn::GetAtt"].([]interface{})
				if !ok {
					t.Fatal("Fn::GetAtt is not an array")
				}

				if len(getAtt) != 2 || getAtt[0] != "MyFunction" || getAtt[1] != "Arn" {
					t.Errorf("expected Fn::GetAtt [MyFunction, Arn], got %v", getAtt)
				}

				// Check permission resource
				permResource, ok := resources["MyFunctionMyIoTRulePermission"]
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

				if permProps["Principal"] != "iot.amazonaws.com" {
					t.Errorf("expected Principal iot.amazonaws.com, got %v", permProps["Principal"])
				}
			},
		},
		{
			name:              "IoTRule event with SQL version",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "MyIoTRule",
			event: &IoTRuleEvent{
				Sql:              "SELECT * FROM 'sensors/temperature'",
				AwsIotSqlVersion: "2016-03-23",
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				ruleResource := resources["MyFunctionMyIoTRule"].(map[string]interface{})
				props := ruleResource["Properties"].(map[string]interface{})
				payload := props["TopicRulePayload"].(map[string]interface{})

				if payload["AwsIotSqlVersion"] != "2016-03-23" {
					t.Errorf("expected AwsIotSqlVersion 2016-03-23, got %v", payload["AwsIotSqlVersion"])
				}
			},
		},
		{
			name:              "missing Sql property",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "MyIoTRule",
			event:             &IoTRuleEvent{},
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

func TestIoTRuleEventSourceHandler_Validate(t *testing.T) {
	handler := NewIoTRuleEventSourceHandler()

	tests := []struct {
		name    string
		event   *IoTRuleEvent
		wantErr bool
	}{
		{
			name: "valid event",
			event: &IoTRuleEvent{
				Sql: "SELECT * FROM 'topic'",
			},
			wantErr: false,
		},
		{
			name:    "missing Sql",
			event:   &IoTRuleEvent{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Validate(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
