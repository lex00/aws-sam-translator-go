package iam

import (
	"encoding/json"
	"testing"
)

func TestNewPolicyDocument(t *testing.T) {
	doc := NewPolicyDocument()

	if doc.Version != PolicyDocumentVersion {
		t.Errorf("expected Version %s, got %s", PolicyDocumentVersion, doc.Version)
	}

	if len(doc.Statement) != 0 {
		t.Errorf("expected empty Statement, got %d statements", len(doc.Statement))
	}
}

func TestNewPolicyDocumentWithId(t *testing.T) {
	doc := NewPolicyDocumentWithId("test-policy-id")

	if doc.Id != "test-policy-id" {
		t.Errorf("expected Id 'test-policy-id', got %s", doc.Id)
	}
}

func TestPolicyDocumentAddStatement(t *testing.T) {
	doc := NewPolicyDocument()
	stmt := NewAllowStatement().
		WithAction("s3:GetObject").
		WithResource("arn:aws:s3:::my-bucket/*")

	doc.AddStatement(stmt)

	if len(doc.Statement) != 1 {
		t.Errorf("expected 1 statement, got %d", len(doc.Statement))
	}
}

func TestPolicyDocumentToMap(t *testing.T) {
	doc := NewPolicyDocument()
	stmt := NewAllowStatement().
		WithSid("AllowS3Read").
		WithAction("s3:GetObject").
		WithResource("arn:aws:s3:::my-bucket/*")

	doc.AddStatement(stmt)

	m := doc.ToMap()

	if m["Version"] != PolicyDocumentVersion {
		t.Errorf("expected Version in map")
	}

	statements, ok := m["Statement"].([]interface{})
	if !ok {
		t.Fatalf("expected Statement to be []interface{}")
	}

	if len(statements) != 1 {
		t.Errorf("expected 1 statement in map, got %d", len(statements))
	}
}

func TestPolicyDocumentValidate(t *testing.T) {
	tests := []struct {
		name    string
		doc     *PolicyDocument
		wantErr bool
	}{
		{
			name: "valid document",
			doc: func() *PolicyDocument {
				doc := NewPolicyDocument()
				doc.AddStatement(NewAllowStatement().WithAction("s3:*").WithResource("*"))
				return doc
			}(),
			wantErr: false,
		},
		{
			name:    "empty statements",
			doc:     NewPolicyDocument(),
			wantErr: true,
		},
		{
			name: "empty version",
			doc: &PolicyDocument{
				Version:   "",
				Statement: []*Statement{NewAllowStatement().WithAction("s3:*").WithResource("*")},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStatementToMap(t *testing.T) {
	stmt := NewAllowStatement().
		WithSid("TestStatement").
		WithActions("s3:GetObject", "s3:PutObject").
		WithResource("arn:aws:s3:::my-bucket/*").
		WithCondition("StringEquals", map[string]interface{}{
			"aws:PrincipalAccount": "123456789012",
		})

	m := stmt.ToMap()

	if m["Effect"] != EffectAllow {
		t.Errorf("expected Effect 'Allow', got %v", m["Effect"])
	}

	if m["Sid"] != "TestStatement" {
		t.Errorf("expected Sid 'TestStatement', got %v", m["Sid"])
	}

	actions, ok := m["Action"].([]interface{})
	if !ok {
		t.Fatalf("expected Action to be []interface{}")
	}

	if len(actions) != 2 {
		t.Errorf("expected 2 actions, got %d", len(actions))
	}

	if m["Resource"] != "arn:aws:s3:::my-bucket/*" {
		t.Errorf("expected Resource, got %v", m["Resource"])
	}

	condition, ok := m["Condition"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Condition to be map")
	}

	if _, ok := condition["StringEquals"]; !ok {
		t.Errorf("expected StringEquals condition")
	}
}

func TestStatementValidate(t *testing.T) {
	tests := []struct {
		name    string
		stmt    *Statement
		wantErr bool
	}{
		{
			name:    "valid allow statement",
			stmt:    NewAllowStatement().WithAction("s3:*").WithResource("*"),
			wantErr: false,
		},
		{
			name:    "valid deny statement",
			stmt:    NewDenyStatement().WithAction("s3:*").WithResource("*"),
			wantErr: false,
		},
		{
			name:    "valid trust policy statement",
			stmt:    NewAllowStatement().WithAction("sts:AssumeRole").WithServicePrincipal("lambda.amazonaws.com"),
			wantErr: false,
		},
		{
			name:    "invalid effect",
			stmt:    NewStatement("Invalid").WithAction("s3:*").WithResource("*"),
			wantErr: true,
		},
		{
			name:    "missing action",
			stmt:    NewAllowStatement().WithResource("*"),
			wantErr: true,
		},
		{
			name:    "missing resource for non-trust policy",
			stmt:    NewAllowStatement().WithAction("s3:*"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.stmt.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicyDocumentBuilder(t *testing.T) {
	doc := NewPolicyDocumentBuilder().
		WithId("my-policy").
		AllowActions("s3:GetObject", "s3:ListBucket").
		OnResources("arn:aws:s3:::my-bucket", "arn:aws:s3:::my-bucket/*").
		WithSid("S3Access").
		Done().
		DenyActions("s3:DeleteBucket").
		OnResources("arn:aws:s3:::my-bucket").
		WithSid("DenyDelete").
		Done().
		Build()

	if doc.Id != "my-policy" {
		t.Errorf("expected Id 'my-policy', got %s", doc.Id)
	}

	if len(doc.Statement) != 2 {
		t.Errorf("expected 2 statements, got %d", len(doc.Statement))
	}

	if doc.Statement[0].Effect != EffectAllow {
		t.Errorf("expected first statement to be Allow")
	}

	if doc.Statement[1].Effect != EffectDeny {
		t.Errorf("expected second statement to be Deny")
	}
}

func TestNewAssumeRolePolicyForService(t *testing.T) {
	doc := NewAssumeRolePolicyForService("lambda.amazonaws.com")

	if len(doc.Statement) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Statement))
	}

	stmt := doc.Statement[0]

	if stmt.Effect != EffectAllow {
		t.Errorf("expected Allow effect")
	}

	if stmt.Action != "sts:AssumeRole" {
		t.Errorf("expected sts:AssumeRole action, got %v", stmt.Action)
	}

	principal, ok := stmt.Principal.(map[string]interface{})
	if !ok {
		t.Fatalf("expected Principal to be map")
	}

	if principal["Service"] != "lambda.amazonaws.com" {
		t.Errorf("expected Service principal lambda.amazonaws.com, got %v", principal["Service"])
	}
}

func TestNewAssumeRolePolicyForServices(t *testing.T) {
	doc := NewAssumeRolePolicyForServices([]string{"lambda.amazonaws.com", "apigateway.amazonaws.com"})

	if len(doc.Statement) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(doc.Statement))
	}

	principal, ok := doc.Statement[0].Principal.(map[string]interface{})
	if !ok {
		t.Fatalf("expected Principal to be map")
	}

	services, ok := principal["Service"].([]interface{})
	if !ok {
		t.Fatalf("expected Service to be []interface{}")
	}

	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}
}

func TestLambdaBasicExecutionPolicy(t *testing.T) {
	logGroupArn := "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/my-function:*"
	doc := LambdaBasicExecutionPolicy(logGroupArn)

	if err := doc.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	m := doc.ToMap()

	// Verify it can be serialized to JSON
	_, err := json.Marshal(m)
	if err != nil {
		t.Errorf("failed to marshal to JSON: %v", err)
	}
}

func TestDynamoDBCrudPolicy(t *testing.T) {
	tableArn := "arn:aws:dynamodb:us-east-1:123456789012:table/my-table"
	doc := DynamoDBCrudPolicy(tableArn)

	if err := doc.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	// Should have 8 DynamoDB actions
	actions, ok := doc.Statement[0].Action.([]interface{})
	if !ok {
		t.Fatalf("expected Action to be []interface{}")
	}

	if len(actions) != 8 {
		t.Errorf("expected 8 actions, got %d", len(actions))
	}
}

func TestS3CrudPolicy(t *testing.T) {
	bucketArn := "arn:aws:s3:::my-bucket"
	objectsArn := "arn:aws:s3:::my-bucket/*"
	doc := S3CrudPolicy(bucketArn, objectsArn)

	if err := doc.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	if len(doc.Statement) != 2 {
		t.Errorf("expected 2 statements, got %d", len(doc.Statement))
	}
}

func TestVPCAccessPolicy(t *testing.T) {
	doc := VPCAccessPolicy()

	if err := doc.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	if doc.Statement[0].Resource != "*" {
		t.Errorf("expected Resource to be '*', got %v", doc.Statement[0].Resource)
	}
}

func TestStatementWithMultiplePrincipals(t *testing.T) {
	stmt := NewAllowStatement().
		WithAction("sts:AssumeRole").
		WithPrincipal(map[string]interface{}{
			"Service": []interface{}{
				"lambda.amazonaws.com",
				"apigateway.amazonaws.com",
			},
			"AWS": "arn:aws:iam::123456789012:root",
		})

	m := stmt.ToMap()

	principal, ok := m["Principal"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected Principal to be map")
	}

	services, ok := principal["Service"].([]interface{})
	if !ok {
		t.Fatalf("expected Service to be []interface{}")
	}

	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}

	aws, ok := principal["AWS"].(string)
	if !ok {
		t.Fatalf("expected AWS to be string")
	}

	if aws != "arn:aws:iam::123456789012:root" {
		t.Errorf("unexpected AWS principal: %s", aws)
	}
}
