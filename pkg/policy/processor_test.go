package policy

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestProcessor_Version(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if p.Version() == "" {
		t.Error("Version() is empty")
	}
}

func TestProcessor_TemplateNames(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	names := p.TemplateNames()
	if len(names) == 0 {
		t.Error("TemplateNames() returned empty list")
	}

	// Should have 81 templates based on the JSON file
	if len(names) != 81 {
		t.Errorf("TemplateNames() returned %d templates, expected 81", len(names))
	}
}

func TestProcessor_HasTemplate(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name     string
		template string
		expected bool
	}{
		{"existing template", "DynamoDBCrudPolicy", true},
		{"another existing", "S3ReadPolicy", true},
		{"non-existing", "NonExistentPolicy", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.HasTemplate(tt.template); got != tt.expected {
				t.Errorf("HasTemplate(%q) = %v, want %v", tt.template, got, tt.expected)
			}
		})
	}
}

func TestProcessor_GetTemplate(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test getting an existing template
	tmpl, ok := p.GetTemplate("DynamoDBCrudPolicy")
	if !ok {
		t.Fatal("GetTemplate(DynamoDBCrudPolicy) returned false")
	}
	if tmpl.Description == "" {
		t.Error("Template description is empty")
	}
	if tmpl.Definition == nil {
		t.Error("Template definition is nil")
	}

	// Test getting non-existing template
	_, ok = p.GetTemplate("NonExistent")
	if ok {
		t.Error("GetTemplate(NonExistent) should return false")
	}
}

func TestProcessor_Expand_SimpleTemplate(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// AMIDescribePolicy has no parameters
	result, err := p.Expand("AMIDescribePolicy", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	// Verify structure
	statements, ok := result["Statement"].([]interface{})
	if !ok {
		t.Fatal("Result should have Statement array")
	}
	if len(statements) == 0 {
		t.Fatal("Statement array should not be empty")
	}

	// Verify first statement
	stmt, ok := statements[0].(map[string]interface{})
	if !ok {
		t.Fatal("Statement should be a map")
	}
	if stmt["Effect"] != "Allow" {
		t.Errorf("Effect = %v, want Allow", stmt["Effect"])
	}
}

func TestProcessor_Expand_WithParameters(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// DynamoDBCrudPolicy requires TableName parameter
	params := map[string]interface{}{
		"TableName": "MyTable",
	}
	result, err := p.Expand("DynamoDBCrudPolicy", params)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	// Verify structure
	statements, ok := result["Statement"].([]interface{})
	if !ok {
		t.Fatal("Result should have Statement array")
	}
	if len(statements) == 0 {
		t.Fatal("Statement array should not be empty")
	}
}

func TestProcessor_Expand_MissingParameter(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// DynamoDBCrudPolicy requires TableName parameter
	_, err = p.Expand("DynamoDBCrudPolicy", map[string]interface{}{})
	if err == nil {
		t.Error("Expand() should error when required parameter is missing")
	}
}

func TestProcessor_Expand_UnknownTemplate(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = p.Expand("NonExistentPolicy", map[string]interface{}{})
	if err == nil {
		t.Error("Expand() should error for unknown template")
	}
}

func TestProcessor_ExpandStatements(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	statements, err := p.ExpandStatements("AMIDescribePolicy", map[string]interface{}{})
	if err != nil {
		t.Fatalf("ExpandStatements() error = %v", err)
	}

	if len(statements) == 0 {
		t.Error("ExpandStatements() returned empty slice")
	}
}

func TestProcessor_Expand_ParameterSubstitution(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test with AWSSecretsManagerGetSecretValuePolicy which uses Fn::Sub with Ref
	params := map[string]interface{}{
		"SecretArn": "arn:aws:secretsmanager:us-east-1:123456789:secret:mysecret",
	}

	result, err := p.Expand("AWSSecretsManagerGetSecretValuePolicy", params)
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	// Convert to JSON for easier inspection
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	// The result should contain our parameter value substituted
	jsonStr := string(jsonBytes)
	if jsonStr == "" {
		t.Error("JSON output is empty")
	}

	// Verify the Fn::Sub structure is preserved with substituted values
	statements := result["Statement"].([]interface{})
	stmt := statements[0].(map[string]interface{})
	resource := stmt["Resource"].(map[string]interface{})
	fnSub := resource["Fn::Sub"].([]interface{})

	// The second element should have our substituted value
	varMap := fnSub[1].(map[string]interface{})
	if varMap["secretArn"] != params["SecretArn"] {
		t.Errorf("Parameter not substituted correctly, got %v", varMap["secretArn"])
	}
}

func TestNewFromBytes_InvalidJSON(t *testing.T) {
	_, err := NewFromBytes([]byte("invalid json"))
	if err == nil {
		t.Error("NewFromBytes() should error on invalid JSON")
	}
}

func TestNewFromBytes_ValidJSON(t *testing.T) {
	jsonData := `{
		"Templates": {
			"TestPolicy": {
				"Definition": {
					"Statement": [{"Effect": "Allow", "Action": ["s3:GetObject"], "Resource": "*"}]
				},
				"Description": "Test policy",
				"Parameters": {}
			}
		},
		"Version": "1.0.0"
	}`

	p, err := NewFromBytes([]byte(jsonData))
	if err != nil {
		t.Fatalf("NewFromBytes() error = %v", err)
	}

	if !p.HasTemplate("TestPolicy") {
		t.Error("TestPolicy should exist")
	}

	if p.Version() != "1.0.0" {
		t.Errorf("Version() = %v, want 1.0.0", p.Version())
	}
}
