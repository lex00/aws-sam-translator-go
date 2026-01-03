package iot

import (
	"testing"
)

func TestResourceTypeTopicRule(t *testing.T) {
	if ResourceTypeTopicRule != "AWS::IoT::TopicRule" {
		t.Errorf("expected ResourceTypeTopicRule to be 'AWS::IoT::TopicRule', got %s", ResourceTypeTopicRule)
	}
}

func TestNewTopicRule(t *testing.T) {
	sql := "SELECT * FROM 'my/topic'"
	functionArn := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction"

	rule := NewTopicRule(sql, functionArn)

	if rule == nil {
		t.Fatal("NewTopicRule returned nil")
	}
	if rule.TopicRulePayload == nil {
		t.Fatal("TopicRulePayload is nil")
	}
	if rule.TopicRulePayload.Sql != sql {
		t.Errorf("expected Sql %q, got %q", sql, rule.TopicRulePayload.Sql)
	}
	if len(rule.TopicRulePayload.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(rule.TopicRulePayload.Actions))
	}
	if rule.TopicRulePayload.Actions[0].Lambda == nil {
		t.Fatal("Lambda action is nil")
	}
	if rule.TopicRulePayload.Actions[0].Lambda.FunctionArn != functionArn {
		t.Errorf("expected FunctionArn %q, got %q", functionArn, rule.TopicRulePayload.Actions[0].Lambda.FunctionArn)
	}
}

func TestNewTopicRule_WithIntrinsicFunction(t *testing.T) {
	sql := "SELECT * FROM 'sensors/+'"
	functionArn := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyFunction", "Arn"},
	}

	rule := NewTopicRule(sql, functionArn)

	if rule.TopicRulePayload.Actions[0].Lambda.FunctionArn == nil {
		t.Fatal("FunctionArn is nil")
	}
	arnMap, ok := rule.TopicRulePayload.Actions[0].Lambda.FunctionArn.(map[string]interface{})
	if !ok {
		t.Fatal("FunctionArn is not a map")
	}
	if _, hasGetAtt := arnMap["Fn::GetAtt"]; !hasGetAtt {
		t.Error("expected Fn::GetAtt in FunctionArn")
	}
}

func TestTopicRule_WithRuleName(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn")
	rule.WithRuleName("MyRule")

	if rule.RuleName != "MyRule" {
		t.Errorf("expected RuleName 'MyRule', got %v", rule.RuleName)
	}
}

func TestTopicRule_WithAwsIotSqlVersion(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn")
	rule.WithAwsIotSqlVersion("2016-03-23")

	if rule.TopicRulePayload.AwsIotSqlVersion != "2016-03-23" {
		t.Errorf("expected AwsIotSqlVersion '2016-03-23', got %v", rule.TopicRulePayload.AwsIotSqlVersion)
	}
}

func TestTopicRule_WithDescription(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn")
	rule.WithDescription("My IoT rule description")

	if rule.TopicRulePayload.Description != "My IoT rule description" {
		t.Errorf("expected Description 'My IoT rule description', got %v", rule.TopicRulePayload.Description)
	}
}

func TestTopicRule_BuilderChaining(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn").
		WithRuleName("ChainedRule").
		WithAwsIotSqlVersion("2016-03-23").
		WithDescription("Chained description")

	if rule.RuleName != "ChainedRule" {
		t.Errorf("expected RuleName 'ChainedRule', got %v", rule.RuleName)
	}
	if rule.TopicRulePayload.AwsIotSqlVersion != "2016-03-23" {
		t.Errorf("expected AwsIotSqlVersion '2016-03-23', got %v", rule.TopicRulePayload.AwsIotSqlVersion)
	}
	if rule.TopicRulePayload.Description != "Chained description" {
		t.Errorf("expected Description 'Chained description', got %v", rule.TopicRulePayload.Description)
	}
}

func TestTopicRule_ToCloudFormation_Basic(t *testing.T) {
	rule := NewTopicRule("SELECT * FROM 'test'", "arn:aws:lambda:us-east-1:123:function:Fn")

	cf := rule.ToCloudFormation()

	// Check Type
	if cf["Type"] != ResourceTypeTopicRule {
		t.Errorf("expected Type %q, got %v", ResourceTypeTopicRule, cf["Type"])
	}

	// Check Properties
	props, ok := cf["Properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Properties is not a map")
	}

	// Check TopicRulePayload
	payload, ok := props["TopicRulePayload"].(map[string]interface{})
	if !ok {
		t.Fatal("TopicRulePayload is not a map")
	}

	if payload["Sql"] != "SELECT * FROM 'test'" {
		t.Errorf("expected Sql 'SELECT * FROM 'test'', got %v", payload["Sql"])
	}

	actions, ok := payload["Actions"].([]interface{})
	if !ok {
		t.Fatal("Actions is not a slice")
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

	if lambdaAction["FunctionArn"] != "arn:aws:lambda:us-east-1:123:function:Fn" {
		t.Errorf("expected FunctionArn, got %v", lambdaAction["FunctionArn"])
	}
}

func TestTopicRule_ToCloudFormation_WithAllOptions(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn").
		WithRuleName("FullRule").
		WithAwsIotSqlVersion("2016-03-23").
		WithDescription("Full description")

	cf := rule.ToCloudFormation()
	props := cf["Properties"].(map[string]interface{})

	if props["RuleName"] != "FullRule" {
		t.Errorf("expected RuleName 'FullRule', got %v", props["RuleName"])
	}

	payload := props["TopicRulePayload"].(map[string]interface{})

	if payload["AwsIotSqlVersion"] != "2016-03-23" {
		t.Errorf("expected AwsIotSqlVersion '2016-03-23', got %v", payload["AwsIotSqlVersion"])
	}
	if payload["Description"] != "Full description" {
		t.Errorf("expected Description 'Full description', got %v", payload["Description"])
	}
}

func TestTopicRule_ToCloudFormation_WithTags(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn")
	rule.Tags = []Tag{
		{Key: "Environment", Value: "Production"},
		{Key: "Team", Value: "IoT"},
	}

	cf := rule.ToCloudFormation()
	props := cf["Properties"].(map[string]interface{})

	tags, ok := props["Tags"].([]interface{})
	if !ok {
		t.Fatal("Tags is not a slice")
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}

	tag1 := tags[0].(map[string]interface{})
	if tag1["Key"] != "Environment" || tag1["Value"] != "Production" {
		t.Errorf("unexpected first tag: %v", tag1)
	}

	tag2 := tags[1].(map[string]interface{})
	if tag2["Key"] != "Team" || tag2["Value"] != "IoT" {
		t.Errorf("unexpected second tag: %v", tag2)
	}
}

func TestTopicRule_ToCloudFormation_WithRuleDisabled(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn")
	rule.TopicRulePayload.RuleDisabled = true

	cf := rule.ToCloudFormation()
	props := cf["Properties"].(map[string]interface{})
	payload := props["TopicRulePayload"].(map[string]interface{})

	if payload["RuleDisabled"] != true {
		t.Errorf("expected RuleDisabled true, got %v", payload["RuleDisabled"])
	}
}

func TestTopicRule_ToCloudFormation_WithErrorAction(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn")
	rule.TopicRulePayload.ErrorAction = &Action{
		Lambda: &LambdaAction{
			FunctionArn: "arn:aws:lambda:us-east-1:123:function:ErrorHandler",
		},
	}

	cf := rule.ToCloudFormation()
	props := cf["Properties"].(map[string]interface{})
	payload := props["TopicRulePayload"].(map[string]interface{})

	errorAction, ok := payload["ErrorAction"].(map[string]interface{})
	if !ok {
		t.Fatal("ErrorAction is not a map")
	}

	lambdaAction, ok := errorAction["Lambda"].(map[string]interface{})
	if !ok {
		t.Fatal("ErrorAction Lambda is not a map")
	}

	if lambdaAction["FunctionArn"] != "arn:aws:lambda:us-east-1:123:function:ErrorHandler" {
		t.Errorf("expected ErrorAction FunctionArn, got %v", lambdaAction["FunctionArn"])
	}
}

func TestTopicRule_ToCloudFormation_OmitsNilOptionalFields(t *testing.T) {
	rule := NewTopicRule("SELECT *", "arn:aws:lambda:us-east-1:123:function:Fn")

	cf := rule.ToCloudFormation()
	props := cf["Properties"].(map[string]interface{})
	payload := props["TopicRulePayload"].(map[string]interface{})

	// These should not be present when nil
	if _, exists := props["RuleName"]; exists {
		t.Error("RuleName should not be present when nil")
	}
	if _, exists := props["Tags"]; exists {
		t.Error("Tags should not be present when empty")
	}
	if _, exists := payload["AwsIotSqlVersion"]; exists {
		t.Error("AwsIotSqlVersion should not be present when nil")
	}
	if _, exists := payload["Description"]; exists {
		t.Error("Description should not be present when nil")
	}
	if _, exists := payload["RuleDisabled"]; exists {
		t.Error("RuleDisabled should not be present when nil")
	}
	if _, exists := payload["ErrorAction"]; exists {
		t.Error("ErrorAction should not be present when nil")
	}
}
