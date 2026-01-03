package cognito

import (
	"testing"
)

func TestResourceTypeUserPool(t *testing.T) {
	if ResourceTypeUserPool != "AWS::Cognito::UserPool" {
		t.Errorf("expected ResourceTypeUserPool to be 'AWS::Cognito::UserPool', got %s", ResourceTypeUserPool)
	}
}

func TestCognitoTriggerConstants(t *testing.T) {
	tests := []struct {
		trigger  CognitoTrigger
		expected string
	}{
		{TriggerPreSignUp, "PreSignUp"},
		{TriggerPostConfirmation, "PostConfirmation"},
		{TriggerPreAuthentication, "PreAuthentication"},
		{TriggerPostAuthentication, "PostAuthentication"},
		{TriggerPreTokenGeneration, "PreTokenGeneration"},
		{TriggerCustomMessage, "CustomMessage"},
		{TriggerUserMigration, "UserMigration"},
		{TriggerDefineAuthChallenge, "DefineAuthChallenge"},
		{TriggerCreateAuthChallenge, "CreateAuthChallenge"},
		{TriggerVerifyAuthChallengeResponse, "VerifyAuthChallengeResponse"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.trigger) != tt.expected {
				t.Errorf("expected trigger %q, got %q", tt.expected, string(tt.trigger))
			}
		})
	}
}

func TestGetLambdaConfigProperty(t *testing.T) {
	tests := []struct {
		trigger  CognitoTrigger
		expected string
	}{
		{TriggerPreSignUp, "PreSignUp"},
		{TriggerPostConfirmation, "PostConfirmation"},
		{TriggerPreAuthentication, "PreAuthentication"},
		{TriggerPostAuthentication, "PostAuthentication"},
		{TriggerPreTokenGeneration, "PreTokenGeneration"},
		{TriggerCustomMessage, "CustomMessage"},
		{TriggerUserMigration, "UserMigration"},
		{TriggerDefineAuthChallenge, "DefineAuthChallenge"},
		{TriggerCreateAuthChallenge, "CreateAuthChallenge"},
		{TriggerVerifyAuthChallengeResponse, "VerifyAuthChallengeResponse"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GetLambdaConfigProperty(tt.trigger)
			if result != tt.expected {
				t.Errorf("GetLambdaConfigProperty(%q) = %q, expected %q", tt.trigger, result, tt.expected)
			}
		})
	}
}

func TestLambdaConfig_AllFields(t *testing.T) {
	config := &LambdaConfig{
		CreateAuthChallenge:         "arn:aws:lambda:us-east-1:123:function:CreateChallenge",
		CustomMessage:               "arn:aws:lambda:us-east-1:123:function:CustomMsg",
		DefineAuthChallenge:         "arn:aws:lambda:us-east-1:123:function:DefineChallenge",
		PostAuthentication:          "arn:aws:lambda:us-east-1:123:function:PostAuth",
		PostConfirmation:            "arn:aws:lambda:us-east-1:123:function:PostConfirm",
		PreAuthentication:           "arn:aws:lambda:us-east-1:123:function:PreAuth",
		PreSignUp:                   "arn:aws:lambda:us-east-1:123:function:PreSignUp",
		PreTokenGeneration:          "arn:aws:lambda:us-east-1:123:function:PreToken",
		UserMigration:               "arn:aws:lambda:us-east-1:123:function:Migration",
		VerifyAuthChallengeResponse: "arn:aws:lambda:us-east-1:123:function:Verify",
		KMSKeyID:                    "arn:aws:kms:us-east-1:123:key/abc",
	}

	if config.CreateAuthChallenge != "arn:aws:lambda:us-east-1:123:function:CreateChallenge" {
		t.Error("CreateAuthChallenge field not set correctly")
	}
	if config.CustomMessage != "arn:aws:lambda:us-east-1:123:function:CustomMsg" {
		t.Error("CustomMessage field not set correctly")
	}
	if config.DefineAuthChallenge != "arn:aws:lambda:us-east-1:123:function:DefineChallenge" {
		t.Error("DefineAuthChallenge field not set correctly")
	}
	if config.PostAuthentication != "arn:aws:lambda:us-east-1:123:function:PostAuth" {
		t.Error("PostAuthentication field not set correctly")
	}
	if config.PostConfirmation != "arn:aws:lambda:us-east-1:123:function:PostConfirm" {
		t.Error("PostConfirmation field not set correctly")
	}
	if config.PreAuthentication != "arn:aws:lambda:us-east-1:123:function:PreAuth" {
		t.Error("PreAuthentication field not set correctly")
	}
	if config.PreSignUp != "arn:aws:lambda:us-east-1:123:function:PreSignUp" {
		t.Error("PreSignUp field not set correctly")
	}
	if config.PreTokenGeneration != "arn:aws:lambda:us-east-1:123:function:PreToken" {
		t.Error("PreTokenGeneration field not set correctly")
	}
	if config.UserMigration != "arn:aws:lambda:us-east-1:123:function:Migration" {
		t.Error("UserMigration field not set correctly")
	}
	if config.VerifyAuthChallengeResponse != "arn:aws:lambda:us-east-1:123:function:Verify" {
		t.Error("VerifyAuthChallengeResponse field not set correctly")
	}
	if config.KMSKeyID != "arn:aws:kms:us-east-1:123:key/abc" {
		t.Error("KMSKeyID field not set correctly")
	}
}

func TestLambdaConfig_WithIntrinsicFunctions(t *testing.T) {
	config := &LambdaConfig{
		PreSignUp: map[string]interface{}{
			"Fn::GetAtt": []interface{}{"PreSignUpFunction", "Arn"},
		},
		PostConfirmation: map[string]interface{}{
			"Ref": "PostConfirmationFunctionArn",
		},
	}

	preSignUp, ok := config.PreSignUp.(map[string]interface{})
	if !ok {
		t.Fatal("PreSignUp is not a map")
	}
	if _, hasGetAtt := preSignUp["Fn::GetAtt"]; !hasGetAtt {
		t.Error("PreSignUp should have Fn::GetAtt")
	}

	postConfirm, ok := config.PostConfirmation.(map[string]interface{})
	if !ok {
		t.Fatal("PostConfirmation is not a map")
	}
	if _, hasRef := postConfirm["Ref"]; !hasRef {
		t.Error("PostConfirmation should have Ref")
	}
}

func TestCustomEmailSender(t *testing.T) {
	sender := &CustomEmailSender{
		LambdaArn:     "arn:aws:lambda:us-east-1:123:function:EmailSender",
		LambdaVersion: "V1_0",
	}

	if sender.LambdaArn != "arn:aws:lambda:us-east-1:123:function:EmailSender" {
		t.Errorf("expected LambdaArn, got %v", sender.LambdaArn)
	}
	if sender.LambdaVersion != "V1_0" {
		t.Errorf("expected LambdaVersion 'V1_0', got %v", sender.LambdaVersion)
	}
}

func TestCustomSMSSender(t *testing.T) {
	sender := &CustomSMSSender{
		LambdaArn:     "arn:aws:lambda:us-east-1:123:function:SMSSender",
		LambdaVersion: "V1_0",
	}

	if sender.LambdaArn != "arn:aws:lambda:us-east-1:123:function:SMSSender" {
		t.Errorf("expected LambdaArn, got %v", sender.LambdaArn)
	}
	if sender.LambdaVersion != "V1_0" {
		t.Errorf("expected LambdaVersion 'V1_0', got %v", sender.LambdaVersion)
	}
}

func TestPreTokenGenerationConfig(t *testing.T) {
	config := &PreTokenGenerationConfig{
		LambdaArn:     "arn:aws:lambda:us-east-1:123:function:PreTokenGen",
		LambdaVersion: "V2_0",
	}

	if config.LambdaArn != "arn:aws:lambda:us-east-1:123:function:PreTokenGen" {
		t.Errorf("expected LambdaArn, got %v", config.LambdaArn)
	}
	if config.LambdaVersion != "V2_0" {
		t.Errorf("expected LambdaVersion 'V2_0', got %v", config.LambdaVersion)
	}
}

func TestLambdaConfig_WithCustomSenders(t *testing.T) {
	config := &LambdaConfig{
		CustomEmailSender: &CustomEmailSender{
			LambdaArn:     "arn:aws:lambda:us-east-1:123:function:EmailSender",
			LambdaVersion: "V1_0",
		},
		CustomSMSSender: &CustomSMSSender{
			LambdaArn:     "arn:aws:lambda:us-east-1:123:function:SMSSender",
			LambdaVersion: "V1_0",
		},
		PreTokenGenerationConfig: &PreTokenGenerationConfig{
			LambdaArn:     "arn:aws:lambda:us-east-1:123:function:PreTokenGen",
			LambdaVersion: "V2_0",
		},
		KMSKeyID: "arn:aws:kms:us-east-1:123:key/abc",
	}

	if config.CustomEmailSender == nil {
		t.Fatal("CustomEmailSender is nil")
	}
	if config.CustomSMSSender == nil {
		t.Fatal("CustomSMSSender is nil")
	}
	if config.PreTokenGenerationConfig == nil {
		t.Fatal("PreTokenGenerationConfig is nil")
	}
}

func TestCognitoTrigger_CountAllTriggers(t *testing.T) {
	// Ensure we have exactly 10 trigger types as documented
	triggers := []CognitoTrigger{
		TriggerPreSignUp,
		TriggerPostConfirmation,
		TriggerPreAuthentication,
		TriggerPostAuthentication,
		TriggerPreTokenGeneration,
		TriggerCustomMessage,
		TriggerUserMigration,
		TriggerDefineAuthChallenge,
		TriggerCreateAuthChallenge,
		TriggerVerifyAuthChallengeResponse,
	}

	if len(triggers) != 10 {
		t.Errorf("expected 10 trigger types, got %d", len(triggers))
	}

	// Ensure all are unique
	seen := make(map[CognitoTrigger]bool)
	for _, trigger := range triggers {
		if seen[trigger] {
			t.Errorf("duplicate trigger: %s", trigger)
		}
		seen[trigger] = true
	}
}
