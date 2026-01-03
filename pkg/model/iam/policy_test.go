package iam

import (
	"encoding/json"
	"testing"
)

func TestNewPolicy(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewPolicy("my-policy", doc)

	if policy.PolicyName != "my-policy" {
		t.Errorf("expected PolicyName 'my-policy', got %v", policy.PolicyName)
	}

	if policy.PolicyDocument == nil {
		t.Error("expected PolicyDocument to be set")
	}
}

func TestPolicyAttachToRole(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewPolicy("my-policy", doc).
		AttachToRole("my-role").
		AttachToRole("another-role")

	if len(policy.Roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(policy.Roles))
	}
}

func TestPolicyAttachToUser(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewPolicy("my-policy", doc).AttachToUser("my-user")

	if len(policy.Users) != 1 {
		t.Errorf("expected 1 user, got %d", len(policy.Users))
	}
}

func TestPolicyAttachToGroup(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewPolicy("my-policy", doc).AttachToGroup("my-group")

	if len(policy.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(policy.Groups))
	}
}

func TestPolicyToCloudFormation(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().
			WithActions("s3:GetObject", "s3:PutObject").
			WithResource("arn:aws:s3:::my-bucket/*"),
	)
	policy := NewPolicy("my-policy", doc).
		AttachToRole("my-role").
		AttachToUser("my-user")

	props := policy.ToCloudFormation()

	if props["PolicyName"] != "my-policy" {
		t.Errorf("expected PolicyName in props")
	}

	policyDoc, ok := props["PolicyDocument"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected PolicyDocument to be map")
	}

	if policyDoc["Version"] != PolicyDocumentVersion {
		t.Errorf("expected Version in PolicyDocument")
	}

	roles, ok := props["Roles"].([]interface{})
	if !ok {
		t.Fatalf("expected Roles to be []interface{}")
	}

	if len(roles) != 1 {
		t.Errorf("expected 1 role, got %d", len(roles))
	}

	users, ok := props["Users"].([]interface{})
	if !ok {
		t.Fatalf("expected Users to be []interface{}")
	}

	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestPolicyToResource(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewPolicy("my-policy", doc).AttachToRole("my-role")

	resource := policy.ToResource()

	if resource["Type"] != "AWS::IAM::Policy" {
		t.Errorf("expected Type 'AWS::IAM::Policy', got %v", resource["Type"])
	}

	if _, ok := resource["Properties"].(map[string]interface{}); !ok {
		t.Error("expected Properties to be map")
	}
}

func TestPolicyValidate(t *testing.T) {
	tests := []struct {
		name    string
		policy  *Policy
		wantErr bool
	}{
		{
			name: "valid policy with role",
			policy: func() *Policy {
				doc := NewPolicyDocument().AddStatement(
					NewAllowStatement().WithAction("s3:*").WithResource("*"),
				)
				return NewPolicy("my-policy", doc).AttachToRole("my-role")
			}(),
			wantErr: false,
		},
		{
			name: "valid policy with user",
			policy: func() *Policy {
				doc := NewPolicyDocument().AddStatement(
					NewAllowStatement().WithAction("s3:*").WithResource("*"),
				)
				return NewPolicy("my-policy", doc).AttachToUser("my-user")
			}(),
			wantErr: false,
		},
		{
			name: "valid policy with group",
			policy: func() *Policy {
				doc := NewPolicyDocument().AddStatement(
					NewAllowStatement().WithAction("s3:*").WithResource("*"),
				)
				return NewPolicy("my-policy", doc).AttachToGroup("my-group")
			}(),
			wantErr: false,
		},
		{
			name: "missing PolicyName",
			policy: func() *Policy {
				doc := NewPolicyDocument().AddStatement(
					NewAllowStatement().WithAction("s3:*").WithResource("*"),
				)
				return &Policy{PolicyDocument: doc, Roles: []interface{}{"my-role"}}
			}(),
			wantErr: true,
		},
		{
			name: "missing PolicyDocument",
			policy: &Policy{
				PolicyName: "my-policy",
				Roles:      []interface{}{"my-role"},
			},
			wantErr: true,
		},
		{
			name: "no attachments",
			policy: func() *Policy {
				doc := NewPolicyDocument().AddStatement(
					NewAllowStatement().WithAction("s3:*").WithResource("*"),
				)
				return NewPolicy("my-policy", doc) // No attachments
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewManagedPolicy(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewManagedPolicy(doc)

	if policy.PolicyDocument == nil {
		t.Error("expected PolicyDocument to be set")
	}
}

func TestNewManagedPolicyWithName(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewManagedPolicyWithName("my-managed-policy", doc)

	if policy.ManagedPolicyName != "my-managed-policy" {
		t.Errorf("expected ManagedPolicyName 'my-managed-policy', got %v", policy.ManagedPolicyName)
	}
}

func TestManagedPolicyWithDescription(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewManagedPolicy(doc).WithDescription("My managed policy")

	if policy.Description != "My managed policy" {
		t.Errorf("expected Description 'My managed policy', got %s", policy.Description)
	}
}

func TestManagedPolicyWithPath(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewManagedPolicy(doc).WithPath("/my-path/")

	if policy.Path != "/my-path/" {
		t.Errorf("expected Path '/my-path/', got %s", policy.Path)
	}
}

func TestManagedPolicyToCloudFormation(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewManagedPolicyWithName("my-managed-policy", doc).
		WithDescription("Test managed policy").
		WithPath("/service-role/").
		AttachToRole("my-role")

	props := policy.ToCloudFormation()

	if props["ManagedPolicyName"] != "my-managed-policy" {
		t.Errorf("expected ManagedPolicyName in props")
	}

	if props["Description"] != "Test managed policy" {
		t.Errorf("expected Description in props")
	}

	if props["Path"] != "/service-role/" {
		t.Errorf("expected Path in props")
	}

	roles, ok := props["Roles"].([]interface{})
	if !ok {
		t.Fatalf("expected Roles to be []interface{}")
	}

	if len(roles) != 1 {
		t.Errorf("expected 1 role, got %d", len(roles))
	}
}

func TestManagedPolicyToResource(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)
	policy := NewManagedPolicy(doc)

	resource := policy.ToResource()

	if resource["Type"] != "AWS::IAM::ManagedPolicy" {
		t.Errorf("expected Type 'AWS::IAM::ManagedPolicy', got %v", resource["Type"])
	}
}

func TestManagedPolicyValidate(t *testing.T) {
	tests := []struct {
		name    string
		policy  *ManagedPolicy
		wantErr bool
	}{
		{
			name: "valid managed policy",
			policy: func() *ManagedPolicy {
				doc := NewPolicyDocument().AddStatement(
					NewAllowStatement().WithAction("s3:*").WithResource("*"),
				)
				return NewManagedPolicy(doc)
			}(),
			wantErr: false,
		},
		{
			name:    "missing PolicyDocument",
			policy:  &ManagedPolicy{},
			wantErr: true,
		},
		{
			name: "invalid PolicyDocument",
			policy: &ManagedPolicy{
				PolicyDocument: &PolicyDocument{Version: PolicyDocumentVersion}, // No statements
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicyJSONSerialization(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().
			WithSid("AllowS3Access").
			WithActions("s3:GetObject", "s3:PutObject").
			WithResource("arn:aws:s3:::my-bucket/*"),
	)
	policy := NewPolicy("my-policy", doc).AttachToRole("my-role")

	resource := policy.ToResource()

	data, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	// Verify it can be unmarshaled back
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result["Type"] != "AWS::IAM::Policy" {
		t.Errorf("expected Type 'AWS::IAM::Policy', got %v", result["Type"])
	}
}

func TestPolicyWithIntrinsicFunctions(t *testing.T) {
	// Test policy with CloudFormation intrinsic functions
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().
			WithAction("s3:GetObject").
			WithResource(map[string]interface{}{
				"Fn::Sub": "arn:aws:s3:::${BucketName}/*",
			}),
	)

	policy := NewPolicy(
		map[string]interface{}{"Fn::Sub": "${AWS::StackName}-policy"},
		doc,
	).AttachToRole(map[string]interface{}{"Ref": "LambdaRole"})

	props := policy.ToCloudFormation()

	// Verify intrinsic functions are preserved
	policyName, ok := props["PolicyName"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected PolicyName to be map (intrinsic function)")
	}

	if _, hasFnSub := policyName["Fn::Sub"]; !hasFnSub {
		t.Error("expected Fn::Sub in PolicyName")
	}

	roles, rolesOk := props["Roles"].([]interface{})
	if !rolesOk {
		t.Fatalf("expected Roles to be []interface{}")
	}

	roleRef, roleOk := roles[0].(map[string]interface{})
	if !roleOk {
		t.Fatalf("expected role to be map (intrinsic function)")
	}

	if _, ok := roleRef["Ref"]; !ok {
		t.Error("expected Ref in role")
	}
}

func TestManagedPolicyAttachments(t *testing.T) {
	doc := NewPolicyDocument().AddStatement(
		NewAllowStatement().WithAction("s3:*").WithResource("*"),
	)

	policy := NewManagedPolicy(doc).
		AttachToRole("role1").
		AttachToRoles([]interface{}{"role2", "role3"}).
		AttachToUser("user1").
		AttachToUsers([]interface{}{"user2"}).
		AttachToGroup("group1").
		AttachToGroups([]interface{}{"group2", "group3"})

	if len(policy.Roles) != 3 {
		t.Errorf("expected 3 roles, got %d", len(policy.Roles))
	}

	if len(policy.Users) != 2 {
		t.Errorf("expected 2 users, got %d", len(policy.Users))
	}

	if len(policy.Groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(policy.Groups))
	}
}
