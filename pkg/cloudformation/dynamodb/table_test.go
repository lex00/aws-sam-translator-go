package dynamodb

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTable_JSONSerialization(t *testing.T) {
	table := Table{
		TableName: "MyTable",
		AttributeDefinitions: []AttributeDefinition{
			{AttributeName: "PK", AttributeType: "S"},
			{AttributeName: "SK", AttributeType: "S"},
			{AttributeName: "GSI1PK", AttributeType: "S"},
		},
		KeySchema: []KeySchemaElement{
			{AttributeName: "PK", KeyType: "HASH"},
			{AttributeName: "SK", KeyType: "RANGE"},
		},
		BillingMode: "PAY_PER_REQUEST",
		GlobalSecondaryIndexes: []GlobalSecondaryIndex{
			{
				IndexName: "GSI1",
				KeySchema: []KeySchemaElement{
					{AttributeName: "GSI1PK", KeyType: "HASH"},
				},
				Projection: Projection{
					ProjectionType: "ALL",
				},
			},
		},
		StreamSpecification: &StreamSpecification{
			StreamViewType: "NEW_AND_OLD_IMAGES",
		},
		Tags: []Tag{
			{Key: "Environment", Value: "Production"},
		},
	}

	data, err := json.Marshal(table)
	if err != nil {
		t.Fatalf("Failed to marshal table to JSON: %v", err)
	}

	var unmarshaled Table
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal table from JSON: %v", err)
	}

	if unmarshaled.TableName != table.TableName {
		t.Errorf("TableName mismatch: got %v, want %v", unmarshaled.TableName, table.TableName)
	}

	if len(unmarshaled.AttributeDefinitions) != len(table.AttributeDefinitions) {
		t.Errorf("AttributeDefinitions length mismatch: got %d, want %d",
			len(unmarshaled.AttributeDefinitions), len(table.AttributeDefinitions))
	}

	if len(unmarshaled.GlobalSecondaryIndexes) != len(table.GlobalSecondaryIndexes) {
		t.Errorf("GlobalSecondaryIndexes length mismatch: got %d, want %d",
			len(unmarshaled.GlobalSecondaryIndexes), len(table.GlobalSecondaryIndexes))
	}

	if unmarshaled.BillingMode != table.BillingMode {
		t.Errorf("BillingMode mismatch: got %v, want %v", unmarshaled.BillingMode, table.BillingMode)
	}
}

func TestTable_YAMLSerialization(t *testing.T) {
	table := Table{
		TableName: "MyTable",
		AttributeDefinitions: []AttributeDefinition{
			{AttributeName: "id", AttributeType: "S"},
		},
		KeySchema: []KeySchemaElement{
			{AttributeName: "id", KeyType: "HASH"},
		},
		ProvisionedThroughput: &ProvisionedThroughput{
			ReadCapacityUnits:  5,
			WriteCapacityUnits: 5,
		},
	}

	data, err := yaml.Marshal(table)
	if err != nil {
		t.Fatalf("Failed to marshal table to YAML: %v", err)
	}

	var unmarshaled Table
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal table from YAML: %v", err)
	}

	if unmarshaled.TableName != table.TableName {
		t.Errorf("TableName mismatch: got %v, want %v", unmarshaled.TableName, table.TableName)
	}

	if unmarshaled.ProvisionedThroughput == nil {
		t.Error("ProvisionedThroughput should not be nil")
	}
}

func TestTable_WithIntrinsicFunctions(t *testing.T) {
	table := Table{
		TableName: map[string]interface{}{
			"Fn::Sub": "${AWS::StackName}-table",
		},
		AttributeDefinitions: []AttributeDefinition{
			{
				AttributeName: map[string]interface{}{"Ref": "PrimaryKeyName"},
				AttributeType: "S",
			},
		},
		KeySchema: []KeySchemaElement{
			{
				AttributeName: map[string]interface{}{"Ref": "PrimaryKeyName"},
				KeyType:       "HASH",
			},
		},
		SSESpecification: &SSESpecification{
			SSEEnabled: true,
			KMSMasterKeyId: map[string]interface{}{
				"Fn::GetAtt": []string{"MyKey", "Arn"},
			},
		},
	}

	data, err := json.Marshal(table)
	if err != nil {
		t.Fatalf("Failed to marshal table with intrinsics: %v", err)
	}

	var unmarshaled Table
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal table with intrinsics: %v", err)
	}

	if unmarshaled.SSESpecification == nil {
		t.Error("SSESpecification should not be nil")
	}
}

func TestTable_WithLocalSecondaryIndex(t *testing.T) {
	table := Table{
		TableName: "MyTable",
		AttributeDefinitions: []AttributeDefinition{
			{AttributeName: "PK", AttributeType: "S"},
			{AttributeName: "SK", AttributeType: "S"},
			{AttributeName: "LSI1SK", AttributeType: "N"},
		},
		KeySchema: []KeySchemaElement{
			{AttributeName: "PK", KeyType: "HASH"},
			{AttributeName: "SK", KeyType: "RANGE"},
		},
		LocalSecondaryIndexes: []LocalSecondaryIndex{
			{
				IndexName: "LSI1",
				KeySchema: []KeySchemaElement{
					{AttributeName: "PK", KeyType: "HASH"},
					{AttributeName: "LSI1SK", KeyType: "RANGE"},
				},
				Projection: Projection{
					ProjectionType:   "INCLUDE",
					NonKeyAttributes: []interface{}{"attribute1", "attribute2"},
				},
			},
		},
	}

	data, err := json.Marshal(table)
	if err != nil {
		t.Fatalf("Failed to marshal table with LSI: %v", err)
	}

	var unmarshaled Table
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal table with LSI: %v", err)
	}

	if len(unmarshaled.LocalSecondaryIndexes) != 1 {
		t.Errorf("LocalSecondaryIndexes length mismatch: got %d, want 1",
			len(unmarshaled.LocalSecondaryIndexes))
	}
}

func TestTable_WithTTL(t *testing.T) {
	table := Table{
		TableName: "MyTable",
		TimeToLiveSpecification: &TimeToLiveSpecification{
			AttributeName: "expirationTime",
			Enabled:       true,
		},
	}

	data, err := json.Marshal(table)
	if err != nil {
		t.Fatalf("Failed to marshal table with TTL: %v", err)
	}

	var unmarshaled Table
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal table with TTL: %v", err)
	}

	if unmarshaled.TimeToLiveSpecification == nil {
		t.Error("TimeToLiveSpecification should not be nil")
	}
}

func TestTable_WithPointInTimeRecovery(t *testing.T) {
	table := Table{
		TableName: "MyTable",
		PointInTimeRecoverySpecification: &PointInTimeRecoverySpecification{
			PointInTimeRecoveryEnabled: true,
		},
	}

	data, err := json.Marshal(table)
	if err != nil {
		t.Fatalf("Failed to marshal table with PITR: %v", err)
	}

	var unmarshaled Table
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal table with PITR: %v", err)
	}

	if unmarshaled.PointInTimeRecoverySpecification == nil {
		t.Error("PointInTimeRecoverySpecification should not be nil")
	}
}

func TestTable_OmitEmpty(t *testing.T) {
	table := Table{
		TableName: "MyTable",
	}

	data, err := json.Marshal(table)
	if err != nil {
		t.Fatalf("Failed to marshal table: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	// These should be omitted because they're empty
	omittedFields := []string{
		"AttributeDefinitions",
		"KeySchema",
		"BillingMode",
		"ProvisionedThroughput",
		"GlobalSecondaryIndexes",
		"LocalSecondaryIndexes",
		"StreamSpecification",
		"Tags",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestGlobalSecondaryIndex_JSONSerialization(t *testing.T) {
	gsi := GlobalSecondaryIndex{
		IndexName: "GSI1",
		KeySchema: []KeySchemaElement{
			{AttributeName: "GSI1PK", KeyType: "HASH"},
			{AttributeName: "GSI1SK", KeyType: "RANGE"},
		},
		Projection: Projection{
			ProjectionType: "KEYS_ONLY",
		},
		ProvisionedThroughput: &ProvisionedThroughput{
			ReadCapacityUnits:  10,
			WriteCapacityUnits: 5,
		},
	}

	data, err := json.Marshal(gsi)
	if err != nil {
		t.Fatalf("Failed to marshal GSI: %v", err)
	}

	var unmarshaled GlobalSecondaryIndex
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal GSI: %v", err)
	}

	if unmarshaled.IndexName != gsi.IndexName {
		t.Errorf("IndexName mismatch: got %v, want %v", unmarshaled.IndexName, gsi.IndexName)
	}
}

func TestKinesisStreamSpecification_JSONSerialization(t *testing.T) {
	spec := KinesisStreamSpecification{
		StreamArn:                            "arn:aws:kinesis:us-east-1:123456789012:stream/MyStream",
		ApproximateCreationDateTimePrecision: "MICROSECOND",
	}

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Failed to marshal KinesisStreamSpecification: %v", err)
	}

	var unmarshaled KinesisStreamSpecification
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal KinesisStreamSpecification: %v", err)
	}

	if unmarshaled.StreamArn != spec.StreamArn {
		t.Errorf("StreamArn mismatch: got %v, want %v", unmarshaled.StreamArn, spec.StreamArn)
	}
}

func TestImportSourceSpecification_JSONSerialization(t *testing.T) {
	spec := ImportSourceSpecification{
		S3BucketSource: S3BucketSource{
			S3Bucket:    "my-import-bucket",
			S3KeyPrefix: "data/",
		},
		InputFormat:          "CSV",
		InputCompressionType: "GZIP",
		InputFormatOptions: &InputFormatOptions{
			Csv: &CsvOptions{
				Delimiter:  ",",
				HeaderList: []interface{}{"id", "name", "email"},
			},
		},
	}

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Failed to marshal ImportSourceSpecification: %v", err)
	}

	var unmarshaled ImportSourceSpecification
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ImportSourceSpecification: %v", err)
	}

	if unmarshaled.InputFormat != spec.InputFormat {
		t.Errorf("InputFormat mismatch: got %v, want %v", unmarshaled.InputFormat, spec.InputFormat)
	}
}
