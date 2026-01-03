package plugins

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

func TestPolicyTemplatesPlugin_Name(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	if plugin.Name() != "PolicyTemplatesPlugin" {
		t.Errorf("Expected name 'PolicyTemplatesPlugin', got '%s'", plugin.Name())
	}
}

func TestPolicyTemplatesPlugin_Priority(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	if plugin.Priority() != 400 {
		t.Errorf("Expected priority 400, got %d", plugin.Priority())
	}
}

func TestPolicyTemplatesPlugin_ExpandsS3ReadPolicy(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Policies": []interface{}{
						map[string]interface{}{
							"S3ReadPolicy": map[string]interface{}{
								"BucketName": "my-bucket",
							},
						},
					},
				},
			},
		},
	}

	err = plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myFunc := template.Resources["MyFunction"]
	policies, ok := myFunc.Properties["Policies"].([]interface{})
	if !ok {
		t.Fatalf("Expected Policies to be an array")
	}

	if len(policies) != 1 {
		t.Fatalf("Expected 1 policy, got %d", len(policies))
	}

	policy, ok := policies[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected policy to be a map")
	}

	// Should have Statement field (expanded policy)
	_, ok = policy["Statement"]
	if !ok {
		t.Errorf("Expected expanded policy to have Statement field")
	}
}

func TestPolicyTemplatesPlugin_KeepsManagedPolicyArn(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Policies": []interface{}{
						"arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
					},
				},
			},
		},
	}

	err = plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myFunc := template.Resources["MyFunction"]
	policies, ok := myFunc.Properties["Policies"].([]interface{})
	if !ok {
		t.Fatalf("Expected Policies to be an array")
	}

	if len(policies) != 1 {
		t.Fatalf("Expected 1 policy, got %d", len(policies))
	}

	// Should keep as string
	arn, ok := policies[0].(string)
	if !ok {
		t.Fatalf("Expected policy to be a string")
	}

	if arn != "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole" {
		t.Errorf("Expected ARN to be unchanged, got %s", arn)
	}
}

func TestPolicyTemplatesPlugin_KeepsInlinePolicy(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	inlinePolicy := map[string]interface{}{
		"Statement": []interface{}{
			map[string]interface{}{
				"Effect":   "Allow",
				"Action":   []interface{}{"s3:GetObject"},
				"Resource": "arn:aws:s3:::my-bucket/*",
			},
		},
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler":  "index.handler",
					"Runtime":  "python3.9",
					"Policies": inlinePolicy,
				},
			},
		},
	}

	err = plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myFunc := template.Resources["MyFunction"]
	policies, ok := myFunc.Properties["Policies"].([]interface{})
	if !ok {
		t.Fatalf("Expected Policies to be an array")
	}

	if len(policies) != 1 {
		t.Fatalf("Expected 1 policy, got %d", len(policies))
	}

	policy, ok := policies[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected policy to be a map")
	}

	// Should keep inline policy unchanged
	statements, ok := policy["Statement"].([]interface{})
	if !ok {
		t.Fatalf("Expected Statement to be an array")
	}

	if len(statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(statements))
	}
}

func TestPolicyTemplatesPlugin_MixedPolicies(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
					"Policies": []interface{}{
						"arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
						map[string]interface{}{
							"S3ReadPolicy": map[string]interface{}{
								"BucketName": "my-bucket",
							},
						},
					},
				},
			},
		},
	}

	err = plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	myFunc := template.Resources["MyFunction"]
	policies, ok := myFunc.Properties["Policies"].([]interface{})
	if !ok {
		t.Fatalf("Expected Policies to be an array")
	}

	if len(policies) != 2 {
		t.Fatalf("Expected 2 policies, got %d", len(policies))
	}

	// First should be ARN string
	_, ok = policies[0].(string)
	if !ok {
		t.Errorf("Expected first policy to be a string")
	}

	// Second should be expanded policy
	policy, ok := policies[1].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected second policy to be a map")
	}

	_, ok = policy["Statement"]
	if !ok {
		t.Errorf("Expected expanded policy to have Statement field")
	}
}

func TestPolicyTemplatesPlugin_NoPolicies(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	template := &types.Template{
		Resources: map[string]types.Resource{
			"MyFunction": {
				Type: "AWS::Serverless::Function",
				Properties: map[string]interface{}{
					"Handler": "index.handler",
					"Runtime": "python3.9",
				},
			},
		},
	}

	err = plugin.BeforeTransform(template)
	if err != nil {
		t.Fatalf("BeforeTransform failed: %v", err)
	}

	// Should not error or modify function
	myFunc := template.Resources["MyFunction"]
	_, ok := myFunc.Properties["Policies"]
	if ok {
		t.Errorf("Should not add Policies when not present")
	}
}

func TestPolicyTemplatesPlugin_AfterTransform(t *testing.T) {
	plugin, err := NewPolicyTemplatesPlugin()
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	template := &types.Template{}

	err = plugin.AfterTransform(template)
	if err != nil {
		t.Errorf("AfterTransform should not error, got: %v", err)
	}
}
