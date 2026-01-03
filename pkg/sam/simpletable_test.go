package sam

import (
	"testing"

	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/dynamodb"
)

func TestSimpleTableTransformer_Transform_Minimal(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	// Minimal table with no properties (uses defaults)
	st := &SimpleTable{}

	resources, err := transformer.Transform("MinimalTable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Check resource exists
	resource, ok := resources["MinimalTable"].(map[string]interface{})
	if !ok {
		t.Fatal("MinimalTable resource not found")
	}

	// Check type
	if resource["Type"] != "AWS::DynamoDB::Table" {
		t.Errorf("expected Type 'AWS::DynamoDB::Table', got %v", resource["Type"])
	}

	// Check properties
	props := resource["Properties"].(map[string]interface{})

	// Should have BillingMode PAY_PER_REQUEST by default
	if props["BillingMode"] != "PAY_PER_REQUEST" {
		t.Errorf("expected BillingMode 'PAY_PER_REQUEST', got %v", props["BillingMode"])
	}

	// Should have default primary key
	attrDefs := props["AttributeDefinitions"].([]dynamodb.AttributeDefinition)
	if len(attrDefs) != 1 {
		t.Errorf("expected 1 attribute definition, got %d", len(attrDefs))
	}
	if attrDefs[0].AttributeName != "id" {
		t.Errorf("expected AttributeName 'id', got %v", attrDefs[0].AttributeName)
	}
	if attrDefs[0].AttributeType != "S" {
		t.Errorf("expected AttributeType 'S', got %v", attrDefs[0].AttributeType)
	}

	// Check KeySchema
	keySchema := props["KeySchema"].([]dynamodb.KeySchemaElement)
	if len(keySchema) != 1 {
		t.Errorf("expected 1 key schema element, got %d", len(keySchema))
	}
	if keySchema[0].AttributeName != "id" {
		t.Errorf("expected KeySchema AttributeName 'id', got %v", keySchema[0].AttributeName)
	}
	if keySchema[0].KeyType != "HASH" {
		t.Errorf("expected KeyType 'HASH', got %v", keySchema[0].KeyType)
	}
}

func TestSimpleTableTransformer_Transform_WithPrimaryKey(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	st := &SimpleTable{
		PrimaryKey: &PrimaryKey{
			Name: "member-number",
			Type: "Number",
		},
	}

	resources, err := transformer.Transform("MemberTable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources["MemberTable"].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	attrDefs := props["AttributeDefinitions"].([]dynamodb.AttributeDefinition)
	if attrDefs[0].AttributeName != "member-number" {
		t.Errorf("expected AttributeName 'member-number', got %v", attrDefs[0].AttributeName)
	}
	if attrDefs[0].AttributeType != "N" {
		t.Errorf("expected AttributeType 'N' for Number, got %v", attrDefs[0].AttributeType)
	}

	keySchema := props["KeySchema"].([]dynamodb.KeySchemaElement)
	if keySchema[0].AttributeName != "member-number" {
		t.Errorf("expected KeySchema AttributeName 'member-number', got %v", keySchema[0].AttributeName)
	}
}

func TestSimpleTableTransformer_Transform_WithProvisionedThroughput(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	st := &SimpleTable{
		PrimaryKey: &PrimaryKey{
			Name: "id",
			Type: "String",
		},
		ProvisionedThroughput: &ProvisionedThroughput{
			ReadCapacityUnits:  20,
			WriteCapacityUnits: 10,
		},
	}

	resources, err := transformer.Transform("ProvisionedTable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources["ProvisionedTable"].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	// Should NOT have BillingMode when ProvisionedThroughput is set
	if _, ok := props["BillingMode"]; ok {
		t.Error("BillingMode should not be set when ProvisionedThroughput is specified")
	}

	// Check ProvisionedThroughput
	pt := props["ProvisionedThroughput"].(map[string]interface{})
	if pt["ReadCapacityUnits"] != 20 {
		t.Errorf("expected ReadCapacityUnits 20, got %v", pt["ReadCapacityUnits"])
	}
	if pt["WriteCapacityUnits"] != 10 {
		t.Errorf("expected WriteCapacityUnits 10, got %v", pt["WriteCapacityUnits"])
	}
}

func TestSimpleTableTransformer_Transform_WithSSE(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	st := &SimpleTable{
		SSESpecification: &SSESpecification{
			SSEEnabled: true,
		},
	}

	resources, err := transformer.Transform("SSETable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources["SSETable"].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	sse := props["SSESpecification"].(map[string]interface{})
	if sse["SSEEnabled"] != true {
		t.Errorf("expected SSEEnabled true, got %v", sse["SSEEnabled"])
	}
}

func TestSimpleTableTransformer_Transform_WithPITR(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	st := &SimpleTable{
		PointInTimeRecoverySpecification: &PointInTimeRecoverySpecification{
			PointInTimeRecoveryEnabled: true,
		},
	}

	resources, err := transformer.Transform("PITRTable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources["PITRTable"].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	pitr := props["PointInTimeRecoverySpecification"].(map[string]interface{})
	if pitr["PointInTimeRecoveryEnabled"] != true {
		t.Errorf("expected PointInTimeRecoveryEnabled true, got %v", pitr["PointInTimeRecoveryEnabled"])
	}
}

func TestSimpleTableTransformer_Transform_WithTableName(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	st := &SimpleTable{
		TableName: "MyCustomTableName",
	}

	resources, err := transformer.Transform("NamedTable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources["NamedTable"].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	if props["TableName"] != "MyCustomTableName" {
		t.Errorf("expected TableName 'MyCustomTableName', got %v", props["TableName"])
	}
}

func TestSimpleTableTransformer_Transform_WithTags(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	st := &SimpleTable{
		Tags: map[string]string{
			"Environment": "Production",
			"Team":        "Backend",
		},
	}

	resources, err := transformer.Transform("TaggedTable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources["TaggedTable"].(map[string]interface{})
	props := resource["Properties"].(map[string]interface{})

	tags := props["Tags"].([]dynamodb.Tag)
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	// Tags might be in any order, so check both exist
	tagMap := make(map[interface{}]interface{})
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}

	if tagMap["Environment"] != "Production" {
		t.Errorf("expected Environment tag 'Production', got %v", tagMap["Environment"])
	}
	if tagMap["Team"] != "Backend" {
		t.Errorf("expected Team tag 'Backend', got %v", tagMap["Team"])
	}
}

func TestMapAttributeType(t *testing.T) {
	tests := []struct {
		samType  string
		expected string
	}{
		{"String", "S"},
		{"Number", "N"},
		{"Binary", "B"},
		{"Unknown", "S"}, // Default to S for unknown types
		{"", "S"},        // Default to S for empty
	}

	for _, tt := range tests {
		t.Run(tt.samType, func(t *testing.T) {
			result := mapAttributeType(tt.samType)
			if result != tt.expected {
				t.Errorf("mapAttributeType(%q) = %q, want %q", tt.samType, result, tt.expected)
			}
		})
	}
}

func TestSimpleTableTransformer_Transform_Complete(t *testing.T) {
	transformer := NewSimpleTableTransformer()

	// Complete table matching testdata/input/simpletable.yaml CompleteTable
	st := &SimpleTable{
		PrimaryKey: &PrimaryKey{
			Name: "member-number",
			Type: "Number",
		},
		ProvisionedThroughput: &ProvisionedThroughput{
			ReadCapacityUnits:  20,
			WriteCapacityUnits: 10,
		},
	}

	resources, err := transformer.Transform("CompleteTable", st)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	resource := resources["CompleteTable"].(map[string]interface{})

	// Verify type
	if resource["Type"] != "AWS::DynamoDB::Table" {
		t.Errorf("expected Type 'AWS::DynamoDB::Table', got %v", resource["Type"])
	}

	props := resource["Properties"].(map[string]interface{})

	// Verify AttributeDefinitions
	attrDefs := props["AttributeDefinitions"].([]dynamodb.AttributeDefinition)
	if len(attrDefs) != 1 {
		t.Fatalf("expected 1 attribute definition, got %d", len(attrDefs))
	}
	if attrDefs[0].AttributeName != "member-number" {
		t.Errorf("expected AttributeName 'member-number', got %v", attrDefs[0].AttributeName)
	}
	if attrDefs[0].AttributeType != "N" {
		t.Errorf("expected AttributeType 'N', got %v", attrDefs[0].AttributeType)
	}

	// Verify KeySchema
	keySchema := props["KeySchema"].([]dynamodb.KeySchemaElement)
	if len(keySchema) != 1 {
		t.Fatalf("expected 1 key schema element, got %d", len(keySchema))
	}
	if keySchema[0].AttributeName != "member-number" {
		t.Errorf("expected KeySchema AttributeName 'member-number', got %v", keySchema[0].AttributeName)
	}
	if keySchema[0].KeyType != "HASH" {
		t.Errorf("expected KeyType 'HASH', got %v", keySchema[0].KeyType)
	}

	// Verify ProvisionedThroughput
	pt := props["ProvisionedThroughput"].(map[string]interface{})
	if pt["ReadCapacityUnits"] != 20 {
		t.Errorf("expected ReadCapacityUnits 20, got %v", pt["ReadCapacityUnits"])
	}
	if pt["WriteCapacityUnits"] != 10 {
		t.Errorf("expected WriteCapacityUnits 10, got %v", pt["WriteCapacityUnits"])
	}

	// Should not have BillingMode
	if _, ok := props["BillingMode"]; ok {
		t.Error("BillingMode should not be present when ProvisionedThroughput is set")
	}
}
