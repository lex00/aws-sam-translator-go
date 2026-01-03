package lambda

import (
	"testing"
)

func TestNewPermission(t *testing.T) {
	perm := NewPermission("lambda:InvokeFunction", "my-function", "s3.amazonaws.com")

	if perm.Action != "lambda:InvokeFunction" {
		t.Errorf("expected Action 'lambda:InvokeFunction', got %s", perm.Action)
	}
	if perm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", perm.FunctionName)
	}
	if perm.Principal != "s3.amazonaws.com" {
		t.Errorf("expected Principal 's3.amazonaws.com', got %s", perm.Principal)
	}
}

func TestNewInvokePermission(t *testing.T) {
	perm := NewInvokePermission("my-function", "events.amazonaws.com")

	if perm.Action != "lambda:InvokeFunction" {
		t.Errorf("expected Action 'lambda:InvokeFunction', got %s", perm.Action)
	}
}

func TestNewAPIGatewayPermission(t *testing.T) {
	sourceArn := "arn:aws:execute-api:us-east-1:123456789012:abc123/*/*/*"
	perm := NewAPIGatewayPermission("my-function", sourceArn)

	if perm.Principal != "apigateway.amazonaws.com" {
		t.Errorf("expected Principal 'apigateway.amazonaws.com', got %s", perm.Principal)
	}
	if perm.SourceArn != sourceArn {
		t.Errorf("expected SourceArn %s, got %v", sourceArn, perm.SourceArn)
	}
}

func TestNewS3Permission(t *testing.T) {
	sourceArn := "arn:aws:s3:::my-bucket"
	sourceAccount := "123456789012"
	perm := NewS3Permission("my-function", sourceArn, sourceAccount)

	if perm.Principal != "s3.amazonaws.com" {
		t.Errorf("expected Principal 's3.amazonaws.com', got %s", perm.Principal)
	}
	if perm.SourceArn != sourceArn {
		t.Errorf("expected SourceArn %s, got %v", sourceArn, perm.SourceArn)
	}
	if perm.SourceAccount != sourceAccount {
		t.Errorf("expected SourceAccount %s, got %v", sourceAccount, perm.SourceAccount)
	}
}

func TestNewSNSPermission(t *testing.T) {
	sourceArn := "arn:aws:sns:us-east-1:123456789012:my-topic"
	perm := NewSNSPermission("my-function", sourceArn)

	if perm.Principal != "sns.amazonaws.com" {
		t.Errorf("expected Principal 'sns.amazonaws.com', got %s", perm.Principal)
	}
}

func TestNewEventsPermission(t *testing.T) {
	sourceArn := "arn:aws:events:us-east-1:123456789012:rule/my-rule"
	perm := NewEventsPermission("my-function", sourceArn)

	if perm.Principal != "events.amazonaws.com" {
		t.Errorf("expected Principal 'events.amazonaws.com', got %s", perm.Principal)
	}
}

func TestNewCloudWatchLogsPermission(t *testing.T) {
	sourceArn := "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/*"
	perm := NewCloudWatchLogsPermission("my-function", sourceArn)

	if perm.Principal != "logs.amazonaws.com" {
		t.Errorf("expected Principal 'logs.amazonaws.com', got %s", perm.Principal)
	}
}

func TestNewCognitoPermission(t *testing.T) {
	sourceArn := "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_abc123"
	perm := NewCognitoPermission("my-function", sourceArn)

	if perm.Principal != "cognito-idp.amazonaws.com" {
		t.Errorf("expected Principal 'cognito-idp.amazonaws.com', got %s", perm.Principal)
	}
}

func TestNewIoTPermission(t *testing.T) {
	sourceArn := "arn:aws:iot:us-east-1:123456789012:rule/my-rule"
	perm := NewIoTPermission("my-function", sourceArn)

	if perm.Principal != "iot.amazonaws.com" {
		t.Errorf("expected Principal 'iot.amazonaws.com', got %s", perm.Principal)
	}
}

func TestNewAlexaPermission(t *testing.T) {
	perm := NewAlexaPermission("my-function", "amzn1.ask.skill.12345678-1234-1234-1234-123456789012")

	if perm.Principal != "alexa-appkit.amazon.com" {
		t.Errorf("expected Principal 'alexa-appkit.amazon.com', got %s", perm.Principal)
	}
	if perm.EventSourceToken != "amzn1.ask.skill.12345678-1234-1234-1234-123456789012" {
		t.Errorf("unexpected EventSourceToken: %s", perm.EventSourceToken)
	}
}

func TestPermissionWithSourceArn(t *testing.T) {
	perm := NewInvokePermission("my-function", "events.amazonaws.com").
		WithSourceArn("arn:aws:events:us-east-1:123456789012:rule/my-rule")

	if perm.SourceArn != "arn:aws:events:us-east-1:123456789012:rule/my-rule" {
		t.Errorf("unexpected SourceArn: %v", perm.SourceArn)
	}
}

func TestPermissionWithSourceAccount(t *testing.T) {
	perm := NewInvokePermission("my-function", "s3.amazonaws.com").
		WithSourceAccount("123456789012")

	if perm.SourceAccount != "123456789012" {
		t.Errorf("unexpected SourceAccount: %v", perm.SourceAccount)
	}
}

func TestPermissionWithEventSourceToken(t *testing.T) {
	perm := NewInvokePermission("my-function", "alexa-appkit.amazon.com").
		WithEventSourceToken("token123")

	if perm.EventSourceToken != "token123" {
		t.Errorf("unexpected EventSourceToken: %s", perm.EventSourceToken)
	}
}

func TestPermissionWithPrincipalOrgID(t *testing.T) {
	perm := NewInvokePermission("my-function", "*").
		WithPrincipalOrgID("o-1234567890")

	if perm.PrincipalOrgID != "o-1234567890" {
		t.Errorf("unexpected PrincipalOrgID: %s", perm.PrincipalOrgID)
	}
}

func TestPermissionWithFunctionUrlAuthType(t *testing.T) {
	perm := NewInvokePermission("my-function", "*").
		WithFunctionUrlAuthType("NONE")

	if perm.FunctionUrlAuthType != "NONE" {
		t.Errorf("unexpected FunctionUrlAuthType: %s", perm.FunctionUrlAuthType)
	}
}

func TestPermissionToCloudFormation_Minimal(t *testing.T) {
	perm := NewInvokePermission("my-function", "events.amazonaws.com")

	result := perm.ToCloudFormation()

	if result["Type"] != ResourceTypePermission {
		t.Errorf("expected Type %s, got %v", ResourceTypePermission, result["Type"])
	}

	props := result["Properties"].(map[string]interface{})
	if props["Action"] != "lambda:InvokeFunction" {
		t.Errorf("expected Action in properties")
	}
	if props["FunctionName"] != "my-function" {
		t.Errorf("expected FunctionName in properties")
	}
	if props["Principal"] != "events.amazonaws.com" {
		t.Errorf("expected Principal in properties")
	}
}

func TestPermissionToCloudFormation_Full(t *testing.T) {
	perm := NewPermission("lambda:InvokeFunction", "my-function", "*").
		WithSourceArn("arn:aws:s3:::my-bucket").
		WithSourceAccount("123456789012").
		WithEventSourceToken("token").
		WithPrincipalOrgID("o-123").
		WithFunctionUrlAuthType("AWS_IAM")

	result := perm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["SourceArn"] != "arn:aws:s3:::my-bucket" {
		t.Errorf("expected SourceArn in properties")
	}
	if props["SourceAccount"] != "123456789012" {
		t.Errorf("expected SourceAccount in properties")
	}
	if props["EventSourceToken"] != "token" {
		t.Errorf("expected EventSourceToken in properties")
	}
	if props["PrincipalOrgID"] != "o-123" {
		t.Errorf("expected PrincipalOrgID in properties")
	}
	if props["FunctionUrlAuthType"] != "AWS_IAM" {
		t.Errorf("expected FunctionUrlAuthType in properties")
	}
}

func TestPermissionWithIntrinsicFunctionName(t *testing.T) {
	fnRef := map[string]interface{}{"Ref": "MyLambdaFunction"}
	perm := NewInvokePermission(fnRef, "events.amazonaws.com")

	result := perm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	fnName := props["FunctionName"].(map[string]interface{})
	if fnName["Ref"] != "MyLambdaFunction" {
		t.Errorf("expected Ref to MyLambdaFunction, got %v", fnName)
	}
}
