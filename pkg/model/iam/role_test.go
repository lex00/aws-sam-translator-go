package iam

import (
	"encoding/json"
	"testing"
)

func TestNewRole(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRole(trustPolicy)

	if role.AssumeRolePolicyDocument == nil {
		t.Error("expected AssumeRolePolicyDocument to be set")
	}
}

func TestNewRoleWithName(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRoleWithName("my-lambda-role", trustPolicy)

	if role.RoleName != "my-lambda-role" {
		t.Errorf("expected RoleName 'my-lambda-role', got %v", role.RoleName)
	}
}

func TestRoleWithPath(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRole(trustPolicy).WithPath("/service-role/")

	if role.Path != "/service-role/" {
		t.Errorf("expected Path '/service-role/', got %s", role.Path)
	}
}

func TestRoleWithDescription(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRole(trustPolicy).WithDescription("My test role")

	if role.Description != "My test role" {
		t.Errorf("expected Description 'My test role', got %s", role.Description)
	}
}

func TestRoleAddManagedPolicyArn(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRole(trustPolicy).
		AddManagedPolicyArn("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole").
		AddManagedPolicyArn("arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess")

	if len(role.ManagedPolicyArns) != 2 {
		t.Errorf("expected 2 managed policy ARNs, got %d", len(role.ManagedPolicyArns))
	}
}

func TestRoleAddInlinePolicy(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	inlinePolicy := DynamoDBCrudPolicy("arn:aws:dynamodb:us-east-1:123456789012:table/my-table")

	role := NewRole(trustPolicy).AddInlinePolicy("DynamoDBAccess", inlinePolicy)

	if len(role.Policies) != 1 {
		t.Errorf("expected 1 inline policy, got %d", len(role.Policies))
	}

	if role.Policies[0].PolicyName != "DynamoDBAccess" {
		t.Errorf("expected PolicyName 'DynamoDBAccess', got %s", role.Policies[0].PolicyName)
	}
}

func TestRoleAddTag(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRole(trustPolicy).
		AddTag("Environment", "production").
		AddTag("Team", "platform")

	if len(role.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(role.Tags))
	}

	if role.Tags[0].Key != "Environment" || role.Tags[0].Value != "production" {
		t.Errorf("unexpected first tag: %+v", role.Tags[0])
	}
}

func TestRoleToCloudFormation(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRoleWithName("my-role", trustPolicy).
		WithPath("/service-role/").
		WithDescription("Test role").
		WithMaxSessionDuration(7200).
		AddManagedPolicyArn("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole").
		AddTag("Environment", "test")

	props := role.ToCloudFormation()

	if props["RoleName"] != "my-role" {
		t.Errorf("expected RoleName in props")
	}

	if props["Path"] != "/service-role/" {
		t.Errorf("expected Path in props")
	}

	if props["Description"] != "Test role" {
		t.Errorf("expected Description in props")
	}

	if props["MaxSessionDuration"] != 7200 {
		t.Errorf("expected MaxSessionDuration 7200, got %v", props["MaxSessionDuration"])
	}

	assumeRolePolicy, ok := props["AssumeRolePolicyDocument"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected AssumeRolePolicyDocument to be map")
	}

	if assumeRolePolicy["Version"] != PolicyDocumentVersion {
		t.Errorf("expected Version in AssumeRolePolicyDocument")
	}

	managedPolicies, ok := props["ManagedPolicyArns"].([]interface{})
	if !ok {
		t.Fatalf("expected ManagedPolicyArns to be []interface{}")
	}

	if len(managedPolicies) != 1 {
		t.Errorf("expected 1 managed policy, got %d", len(managedPolicies))
	}

	tags, ok := props["Tags"].([]map[string]interface{})
	if !ok {
		t.Fatalf("expected Tags to be []map[string]interface{}")
	}

	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
}

func TestRoleToResource(t *testing.T) {
	trustPolicy := NewAssumeRolePolicyForService(ServiceLambda)
	role := NewRole(trustPolicy)

	resource := role.ToResource()

	if resource["Type"] != "AWS::IAM::Role" {
		t.Errorf("expected Type 'AWS::IAM::Role', got %v", resource["Type"])
	}

	if _, ok := resource["Properties"].(map[string]interface{}); !ok {
		t.Error("expected Properties to be map")
	}
}

func TestRoleValidate(t *testing.T) {
	tests := []struct {
		name    string
		role    *Role
		wantErr bool
	}{
		{
			name:    "valid role",
			role:    NewRole(NewAssumeRolePolicyForService(ServiceLambda)),
			wantErr: false,
		},
		{
			name:    "missing AssumeRolePolicyDocument",
			role:    &Role{},
			wantErr: true,
		},
		{
			name: "invalid MaxSessionDuration too low",
			role: func() *Role {
				r := NewRole(NewAssumeRolePolicyForService(ServiceLambda))
				r.MaxSessionDuration = 1800 // Less than 3600
				return r
			}(),
			wantErr: true,
		},
		{
			name: "invalid MaxSessionDuration too high",
			role: func() *Role {
				r := NewRole(NewAssumeRolePolicyForService(ServiceLambda))
				r.MaxSessionDuration = 50000 // More than 43200
				return r
			}(),
			wantErr: true,
		},
		{
			name: "valid MaxSessionDuration",
			role: func() *Role {
				r := NewRole(NewAssumeRolePolicyForService(ServiceLambda))
				r.MaxSessionDuration = 7200
				return r
			}(),
			wantErr: false,
		},
		{
			name: "inline policy missing PolicyName",
			role: func() *Role {
				r := NewRole(NewAssumeRolePolicyForService(ServiceLambda))
				r.Policies = []InlinePolicy{{PolicyDocument: NewPolicyDocument().AddStatement(NewAllowStatement().WithAction("s3:*").WithResource("*"))}}
				return r
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.role.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTrustRelationship(t *testing.T) {
	trust := NewServiceTrustRelationship("lambda.amazonaws.com")
	doc := trust.ToPolicyDocument()

	if err := doc.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	m := doc.ToMap()

	statements, ok := m["Statement"].([]interface{})
	if !ok {
		t.Fatalf("expected Statement to be []interface{}")
	}

	if len(statements) != 1 {
		t.Errorf("expected 1 statement, got %d", len(statements))
	}
}

func TestTrustRelationshipWithCondition(t *testing.T) {
	trust := NewServiceTrustRelationship("lambda.amazonaws.com").
		WithCondition("StringEquals", map[string]interface{}{
			"aws:SourceAccount": "123456789012",
		})

	doc := trust.ToPolicyDocument()

	if len(doc.Statement[0].Condition) == 0 {
		t.Error("expected Condition to be set")
	}
}

func TestNewLambdaExecutionRole(t *testing.T) {
	role := NewLambdaExecutionRole()

	if err := role.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	// Verify the trust policy is for Lambda
	m := role.AssumeRolePolicyDocument.ToMap()
	statements := m["Statement"].([]interface{})
	stmt := statements[0].(map[string]interface{})
	principal := stmt["Principal"].(map[string]interface{})

	if principal["Service"] != ServiceLambda {
		t.Errorf("expected Lambda service principal, got %v", principal["Service"])
	}
}

func TestNewStepFunctionsExecutionRole(t *testing.T) {
	role := NewStepFunctionsExecutionRole()

	if err := role.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	m := role.AssumeRolePolicyDocument.ToMap()
	statements := m["Statement"].([]interface{})
	stmt := statements[0].(map[string]interface{})
	principal := stmt["Principal"].(map[string]interface{})

	if principal["Service"] != ServiceStepFunctions {
		t.Errorf("expected Step Functions service principal, got %v", principal["Service"])
	}
}

func TestRoleJSONSerialization(t *testing.T) {
	role := NewLambdaExecutionRoleWithName("my-lambda-role").
		WithPath("/service-role/").
		WithDescription("Lambda execution role").
		AddManagedPolicyArn(AWSLambdaBasicExecutionRole("aws")).
		AddInlinePolicy("S3Access", S3ReadPolicy("arn:aws:s3:::my-bucket/*")).
		AddTag("Environment", "production")

	resource := role.ToResource()

	data, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	// Verify it can be unmarshaled back
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result["Type"] != "AWS::IAM::Role" {
		t.Errorf("expected Type 'AWS::IAM::Role', got %v", result["Type"])
	}
}

func TestFederatedTrustRelationship(t *testing.T) {
	trust := NewFederatedTrustRelationship("arn:aws:iam::123456789012:saml-provider/MyProvider").
		WithCondition("StringEquals", map[string]interface{}{
			"SAML:aud": "https://signin.aws.amazon.com/saml",
		})

	doc := trust.ToPolicyDocument()

	if err := doc.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	m := doc.ToMap()
	statements := m["Statement"].([]interface{})
	stmt := statements[0].(map[string]interface{})
	principal := stmt["Principal"].(map[string]interface{})

	if principal["Federated"] != "arn:aws:iam::123456789012:saml-provider/MyProvider" {
		t.Errorf("expected Federated principal, got %v", principal["Federated"])
	}
}

func TestAWSTrustRelationship(t *testing.T) {
	trust := NewAWSTrustRelationship("arn:aws:iam::123456789012:root")
	doc := trust.ToPolicyDocument()

	if err := doc.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	m := doc.ToMap()
	statements := m["Statement"].([]interface{})
	stmt := statements[0].(map[string]interface{})
	principal := stmt["Principal"].(map[string]interface{})

	if principal["AWS"] != "arn:aws:iam::123456789012:root" {
		t.Errorf("expected AWS principal, got %v", principal["AWS"])
	}
}
