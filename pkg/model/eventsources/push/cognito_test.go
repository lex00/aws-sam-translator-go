package push

import (
	"testing"
)

func TestCognitoEventSourceHandler_GenerateResources(t *testing.T) {
	handler := NewCognitoEventSourceHandler()

	tests := []struct {
		name              string
		functionLogicalID string
		eventLogicalID    string
		event             *CognitoEvent
		wantErr           bool
		validate          func(t *testing.T, resources map[string]interface{})
	}{
		{
			name:              "basic Cognito event with Ref",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "CognitoPreSignUp",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
				Trigger:  "PreSignUp",
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				// Check that one resource is created (just the permission)
				if len(resources) != 1 {
					t.Errorf("expected 1 resource, got %d", len(resources))
				}

				// Check permission resource
				permResource, ok := resources["MyFunctionCognitoPreSignUpPermission"]
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

				if permProps["Principal"] != "cognito-idp.amazonaws.com" {
					t.Errorf("expected Principal cognito-idp.amazonaws.com, got %v", permProps["Principal"])
				}

				// Check that SourceArn is a GetAtt for the UserPool ARN
				sourceArn, ok := permProps["SourceArn"].(map[string]interface{})
				if !ok {
					t.Fatal("SourceArn is not a map")
				}

				getAtt, ok := sourceArn["Fn::GetAtt"].([]interface{})
				if !ok {
					t.Fatal("Fn::GetAtt is not an array")
				}

				if len(getAtt) != 2 || getAtt[0] != "MyCognitoUserPool" || getAtt[1] != "Arn" {
					t.Errorf("expected Fn::GetAtt [MyCognitoUserPool, Arn], got %v", getAtt)
				}
			},
		},
		{
			name:              "Cognito event with ARN string",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "CognitoPostAuth",
			event: &CognitoEvent{
				UserPool: "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_abc123",
				Trigger:  "PostAuthentication",
			},
			wantErr: false,
			validate: func(t *testing.T, resources map[string]interface{}) {
				permResource := resources["MyFunctionCognitoPostAuthPermission"].(map[string]interface{})
				permProps := permResource["Properties"].(map[string]interface{})

				if permProps["SourceArn"] != "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_abc123" {
					t.Errorf("expected SourceArn to be the ARN string, got %v", permProps["SourceArn"])
				}
			},
		},
		{
			name:              "missing UserPool property",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "CognitoEvent",
			event: &CognitoEvent{
				Trigger: "PreSignUp",
			},
			wantErr: true,
		},
		{
			name:              "missing Trigger property",
			functionLogicalID: "MyFunction",
			eventLogicalID:    "CognitoEvent",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
			},
			wantErr: true,
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

func TestCognitoEventSourceHandler_Validate(t *testing.T) {
	handler := NewCognitoEventSourceHandler()

	tests := []struct {
		name    string
		event   *CognitoEvent
		wantErr bool
	}{
		{
			name: "valid event",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
				Trigger:  "PreSignUp",
			},
			wantErr: false,
		},
		{
			name: "valid event with PostConfirmation trigger",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
				Trigger:  "PostConfirmation",
			},
			wantErr: false,
		},
		{
			name: "invalid trigger type",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
				Trigger:  "InvalidTrigger",
			},
			wantErr: true,
		},
		{
			name: "missing UserPool",
			event: &CognitoEvent{
				Trigger: "PreSignUp",
			},
			wantErr: true,
		},
		{
			name: "missing Trigger",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
			},
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

func TestCognitoEventSourceHandler_GetLambdaConfigUpdate(t *testing.T) {
	handler := NewCognitoEventSourceHandler()

	tests := []struct {
		name              string
		functionLogicalID string
		event             *CognitoEvent
		wantTrigger       string
		wantErr           bool
	}{
		{
			name:              "PreSignUp trigger",
			functionLogicalID: "MyFunction",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
				Trigger:  "PreSignUp",
			},
			wantTrigger: "PreSignUp",
			wantErr:     false,
		},
		{
			name:              "PostAuthentication trigger",
			functionLogicalID: "MyFunction",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
				Trigger:  "PostAuthentication",
			},
			wantTrigger: "PostAuthentication",
			wantErr:     false,
		},
		{
			name:              "invalid trigger",
			functionLogicalID: "MyFunction",
			event: &CognitoEvent{
				UserPool: map[string]interface{}{"Ref": "MyCognitoUserPool"},
				Trigger:  "InvalidTrigger",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trigger, functionArn, err := handler.GetLambdaConfigUpdate(tt.functionLogicalID, tt.event)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLambdaConfigUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if trigger != tt.wantTrigger {
					t.Errorf("expected trigger %s, got %s", tt.wantTrigger, trigger)
				}

				arnMap, ok := functionArn.(map[string]interface{})
				if !ok {
					t.Fatal("functionArn is not a map")
				}

				getAtt, ok := arnMap["Fn::GetAtt"].([]interface{})
				if !ok {
					t.Fatal("Fn::GetAtt is not an array")
				}

				if len(getAtt) != 2 || getAtt[0] != "MyFunction" || getAtt[1] != "Arn" {
					t.Errorf("expected Fn::GetAtt [MyFunction, Arn], got %v", getAtt)
				}
			}
		})
	}
}
