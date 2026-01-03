// Package sam provides SAM resource transformers.
package sam

import (
	"github.com/lex00/aws-sam-translator-go/pkg/cloudformation/dynamodb"
)

// SimpleTable represents an AWS::Serverless::SimpleTable resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-simpletable.html
type SimpleTable struct {
	// PrimaryKey specifies the primary key for the table.
	// If not specified, defaults to {Name: "id", Type: "String"}.
	PrimaryKey *PrimaryKey `json:"PrimaryKey,omitempty" yaml:"PrimaryKey,omitempty"`

	// ProvisionedThroughput specifies the read and write capacity units.
	// If not specified, uses BillingMode: PAY_PER_REQUEST.
	ProvisionedThroughput *ProvisionedThroughput `json:"ProvisionedThroughput,omitempty" yaml:"ProvisionedThroughput,omitempty"`

	// TableName is the name of the DynamoDB table.
	TableName interface{} `json:"TableName,omitempty" yaml:"TableName,omitempty"`

	// SSESpecification specifies server-side encryption settings.
	SSESpecification *SSESpecification `json:"SSESpecification,omitempty" yaml:"SSESpecification,omitempty"`

	// Tags is a map of key-value pairs to apply to the table.
	Tags map[string]string `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// PointInTimeRecoverySpecification specifies point-in-time recovery settings.
	PointInTimeRecoverySpecification *PointInTimeRecoverySpecification `json:"PointInTimeRecoverySpecification,omitempty" yaml:"PointInTimeRecoverySpecification,omitempty"`
}

// PrimaryKey represents the primary key definition for a SimpleTable.
type PrimaryKey struct {
	// Name is the attribute name for the primary key.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Type is the data type for the primary key attribute.
	// Valid values: String, Number, Binary
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

// ProvisionedThroughput represents provisioned throughput settings.
type ProvisionedThroughput struct {
	// ReadCapacityUnits is the maximum number of strongly consistent reads per second.
	ReadCapacityUnits interface{} `json:"ReadCapacityUnits" yaml:"ReadCapacityUnits"`

	// WriteCapacityUnits is the maximum number of writes per second.
	WriteCapacityUnits interface{} `json:"WriteCapacityUnits" yaml:"WriteCapacityUnits"`
}

// SSESpecification represents server-side encryption settings.
type SSESpecification struct {
	// SSEEnabled indicates whether server-side encryption is enabled.
	SSEEnabled interface{} `json:"SSEEnabled,omitempty" yaml:"SSEEnabled,omitempty"`
}

// PointInTimeRecoverySpecification represents point-in-time recovery settings.
type PointInTimeRecoverySpecification struct {
	// PointInTimeRecoveryEnabled indicates whether point-in-time recovery is enabled.
	PointInTimeRecoveryEnabled interface{} `json:"PointInTimeRecoveryEnabled,omitempty" yaml:"PointInTimeRecoveryEnabled,omitempty"`
}

// SimpleTableTransformer transforms AWS::Serverless::SimpleTable to CloudFormation.
type SimpleTableTransformer struct{}

// NewSimpleTableTransformer creates a new SimpleTableTransformer.
func NewSimpleTableTransformer() *SimpleTableTransformer {
	return &SimpleTableTransformer{}
}

// Transform converts a SAM SimpleTable to CloudFormation DynamoDB::Table resource.
func (t *SimpleTableTransformer) Transform(logicalID string, st *SimpleTable) (map[string]interface{}, error) {
	// Get primary key name and type, applying defaults
	pkName, pkType := t.getPrimaryKeyConfig(st)

	// Map SAM type to DynamoDB AttributeType
	attrType := mapAttributeType(pkType)

	// Build AttributeDefinitions
	attrDefs := []dynamodb.AttributeDefinition{
		{
			AttributeName: pkName,
			AttributeType: attrType,
		},
	}

	// Build KeySchema
	keySchema := []dynamodb.KeySchemaElement{
		{
			AttributeName: pkName,
			KeyType:       "HASH",
		},
	}

	// Build properties map
	properties := make(map[string]interface{})
	properties["AttributeDefinitions"] = attrDefs
	properties["KeySchema"] = keySchema

	// Handle ProvisionedThroughput or BillingMode
	if st.ProvisionedThroughput != nil {
		properties["ProvisionedThroughput"] = map[string]interface{}{
			"ReadCapacityUnits":  st.ProvisionedThroughput.ReadCapacityUnits,
			"WriteCapacityUnits": st.ProvisionedThroughput.WriteCapacityUnits,
		}
	} else {
		// Default to on-demand billing
		properties["BillingMode"] = "PAY_PER_REQUEST"
	}

	// Handle optional properties
	if st.TableName != nil {
		properties["TableName"] = st.TableName
	}

	if st.SSESpecification != nil {
		properties["SSESpecification"] = map[string]interface{}{
			"SSEEnabled": st.SSESpecification.SSEEnabled,
		}
	}

	if st.PointInTimeRecoverySpecification != nil {
		properties["PointInTimeRecoverySpecification"] = map[string]interface{}{
			"PointInTimeRecoveryEnabled": st.PointInTimeRecoverySpecification.PointInTimeRecoveryEnabled,
		}
	}

	if len(st.Tags) > 0 {
		tags := make([]dynamodb.Tag, 0, len(st.Tags))
		for k, v := range st.Tags {
			tags = append(tags, dynamodb.Tag{Key: k, Value: v})
		}
		properties["Tags"] = tags
	}

	// Build the CloudFormation resource
	resources := map[string]interface{}{
		logicalID: map[string]interface{}{
			"Type":       "AWS::DynamoDB::Table",
			"Properties": properties,
		},
	}

	return resources, nil
}

// getPrimaryKeyConfig returns the primary key name and type from SimpleTable.
// If PrimaryKey is not specified, returns defaults: "id" and "String".
func (t *SimpleTableTransformer) getPrimaryKeyConfig(st *SimpleTable) (interface{}, string) {
	if st.PrimaryKey == nil {
		return "id", "String"
	}

	name := st.PrimaryKey.Name
	if name == nil {
		name = "id"
	}

	pkType := st.PrimaryKey.Type
	if pkType == "" {
		pkType = "String"
	}

	return name, pkType
}

// mapAttributeType converts SAM SimpleTable type strings to DynamoDB attribute types.
// String -> S, Number -> N, Binary -> B
func mapAttributeType(samType string) string {
	switch samType {
	case "String":
		return "S"
	case "Number":
		return "N"
	case "Binary":
		return "B"
	default:
		// Default to String if unknown
		return "S"
	}
}
