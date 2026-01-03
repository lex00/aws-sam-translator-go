// Package dynamodb provides CloudFormation resource models for AWS DynamoDB.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dynamodb-table.html
package dynamodb

// Table represents an AWS::DynamoDB::Table resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dynamodb-table.html
type Table struct {
	// TableName is the name of the table. If not specified, CloudFormation generates a unique name.
	TableName interface{} `json:"TableName,omitempty" yaml:"TableName,omitempty"`

	// AttributeDefinitions describes the attributes that define the key schema.
	AttributeDefinitions []AttributeDefinition `json:"AttributeDefinitions,omitempty" yaml:"AttributeDefinitions,omitempty"`

	// KeySchema specifies the attributes that make up the primary key for the table.
	KeySchema []KeySchemaElement `json:"KeySchema,omitempty" yaml:"KeySchema,omitempty"`

	// BillingMode specifies how you are charged for read/write throughput.
	// Valid values: PROVISIONED | PAY_PER_REQUEST
	BillingMode string `json:"BillingMode,omitempty" yaml:"BillingMode,omitempty"`

	// ProvisionedThroughput specifies the provisioned throughput for the table.
	ProvisionedThroughput *ProvisionedThroughput `json:"ProvisionedThroughput,omitempty" yaml:"ProvisionedThroughput,omitempty"`

	// GlobalSecondaryIndexes specifies one or more global secondary indexes.
	GlobalSecondaryIndexes []GlobalSecondaryIndex `json:"GlobalSecondaryIndexes,omitempty" yaml:"GlobalSecondaryIndexes,omitempty"`

	// LocalSecondaryIndexes specifies one or more local secondary indexes.
	LocalSecondaryIndexes []LocalSecondaryIndex `json:"LocalSecondaryIndexes,omitempty" yaml:"LocalSecondaryIndexes,omitempty"`

	// StreamSpecification specifies the settings for DynamoDB Streams.
	StreamSpecification *StreamSpecification `json:"StreamSpecification,omitempty" yaml:"StreamSpecification,omitempty"`

	// TableClass specifies the table class. Valid values: STANDARD | STANDARD_INFREQUENT_ACCESS
	TableClass string `json:"TableClass,omitempty" yaml:"TableClass,omitempty"`

	// DeletionProtectionEnabled indicates whether deletion protection is enabled.
	DeletionProtectionEnabled interface{} `json:"DeletionProtectionEnabled,omitempty" yaml:"DeletionProtectionEnabled,omitempty"`

	// ContributorInsightsSpecification specifies the settings for CloudWatch contributor insights.
	ContributorInsightsSpecification *ContributorInsightsSpecification `json:"ContributorInsightsSpecification,omitempty" yaml:"ContributorInsightsSpecification,omitempty"`

	// KinesisStreamSpecification specifies the Kinesis Data Streams configuration.
	KinesisStreamSpecification *KinesisStreamSpecification `json:"KinesisStreamSpecification,omitempty" yaml:"KinesisStreamSpecification,omitempty"`

	// PointInTimeRecoverySpecification specifies the settings for point-in-time recovery.
	PointInTimeRecoverySpecification *PointInTimeRecoverySpecification `json:"PointInTimeRecoverySpecification,omitempty" yaml:"PointInTimeRecoverySpecification,omitempty"`

	// SSESpecification specifies the settings for server-side encryption.
	SSESpecification *SSESpecification `json:"SSESpecification,omitempty" yaml:"SSESpecification,omitempty"`

	// TimeToLiveSpecification specifies the Time to Live (TTL) settings.
	TimeToLiveSpecification *TimeToLiveSpecification `json:"TimeToLiveSpecification,omitempty" yaml:"TimeToLiveSpecification,omitempty"`

	// Tags is a list of key-value pairs to apply to the table.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// ImportSourceSpecification specifies the properties for importing data from S3.
	ImportSourceSpecification *ImportSourceSpecification `json:"ImportSourceSpecification,omitempty" yaml:"ImportSourceSpecification,omitempty"`

	// ResourcePolicy specifies a resource-based policy document for the table.
	ResourcePolicy *ResourcePolicy `json:"ResourcePolicy,omitempty" yaml:"ResourcePolicy,omitempty"`
}

// AttributeDefinition represents an attribute definition for the table.
type AttributeDefinition struct {
	// AttributeName is the name of the attribute.
	AttributeName interface{} `json:"AttributeName" yaml:"AttributeName"`

	// AttributeType is the data type for the attribute. Valid values: S | N | B
	AttributeType string `json:"AttributeType" yaml:"AttributeType"`
}

// KeySchemaElement represents a single element of a key schema.
type KeySchemaElement struct {
	// AttributeName is the name of a key attribute.
	AttributeName interface{} `json:"AttributeName" yaml:"AttributeName"`

	// KeyType is the role of the key attribute. Valid values: HASH | RANGE
	KeyType string `json:"KeyType" yaml:"KeyType"`
}

// ProvisionedThroughput represents the provisioned throughput settings.
type ProvisionedThroughput struct {
	// ReadCapacityUnits is the maximum number of strongly consistent reads per second.
	ReadCapacityUnits interface{} `json:"ReadCapacityUnits" yaml:"ReadCapacityUnits"`

	// WriteCapacityUnits is the maximum number of writes per second.
	WriteCapacityUnits interface{} `json:"WriteCapacityUnits" yaml:"WriteCapacityUnits"`
}

// GlobalSecondaryIndex represents a global secondary index.
type GlobalSecondaryIndex struct {
	// IndexName is the name of the global secondary index.
	IndexName interface{} `json:"IndexName" yaml:"IndexName"`

	// KeySchema specifies the key schema for the global secondary index.
	KeySchema []KeySchemaElement `json:"KeySchema" yaml:"KeySchema"`

	// Projection specifies attributes that are copied from the table to the index.
	Projection Projection `json:"Projection" yaml:"Projection"`

	// ProvisionedThroughput specifies the provisioned throughput for the index.
	ProvisionedThroughput *ProvisionedThroughput `json:"ProvisionedThroughput,omitempty" yaml:"ProvisionedThroughput,omitempty"`

	// ContributorInsightsSpecification specifies CloudWatch contributor insights settings.
	ContributorInsightsSpecification *ContributorInsightsSpecification `json:"ContributorInsightsSpecification,omitempty" yaml:"ContributorInsightsSpecification,omitempty"`
}

// LocalSecondaryIndex represents a local secondary index.
type LocalSecondaryIndex struct {
	// IndexName is the name of the local secondary index.
	IndexName interface{} `json:"IndexName" yaml:"IndexName"`

	// KeySchema specifies the key schema for the local secondary index.
	KeySchema []KeySchemaElement `json:"KeySchema" yaml:"KeySchema"`

	// Projection specifies attributes that are copied from the table to the index.
	Projection Projection `json:"Projection" yaml:"Projection"`
}

// Projection represents the attributes to project into the index.
type Projection struct {
	// ProjectionType specifies the set of attributes to project.
	// Valid values: ALL | KEYS_ONLY | INCLUDE
	ProjectionType string `json:"ProjectionType,omitempty" yaml:"ProjectionType,omitempty"`

	// NonKeyAttributes specifies the non-key attributes to project (for INCLUDE projection).
	NonKeyAttributes []interface{} `json:"NonKeyAttributes,omitempty" yaml:"NonKeyAttributes,omitempty"`
}

// StreamSpecification represents DynamoDB Streams settings.
type StreamSpecification struct {
	// StreamViewType specifies what information is written to the stream.
	// Valid values: KEYS_ONLY | NEW_IMAGE | OLD_IMAGE | NEW_AND_OLD_IMAGES
	StreamViewType string `json:"StreamViewType" yaml:"StreamViewType"`

	// ResourcePolicy specifies a resource-based policy document for the stream.
	ResourcePolicy *ResourcePolicy `json:"ResourcePolicy,omitempty" yaml:"ResourcePolicy,omitempty"`
}

// ContributorInsightsSpecification represents CloudWatch Contributor Insights settings.
type ContributorInsightsSpecification struct {
	// Enabled indicates whether CloudWatch Contributor Insights are enabled.
	Enabled interface{} `json:"Enabled" yaml:"Enabled"`
}

// KinesisStreamSpecification represents Kinesis Data Streams settings.
type KinesisStreamSpecification struct {
	// StreamArn is the ARN for a specific Kinesis data stream.
	StreamArn interface{} `json:"StreamArn" yaml:"StreamArn"`

	// ApproximateCreationDateTimePrecision specifies timestamp precision.
	// Valid values: MILLISECOND | MICROSECOND
	ApproximateCreationDateTimePrecision string `json:"ApproximateCreationDateTimePrecision,omitempty" yaml:"ApproximateCreationDateTimePrecision,omitempty"`
}

// PointInTimeRecoverySpecification represents point-in-time recovery settings.
type PointInTimeRecoverySpecification struct {
	// PointInTimeRecoveryEnabled indicates whether point-in-time recovery is enabled.
	PointInTimeRecoveryEnabled interface{} `json:"PointInTimeRecoveryEnabled,omitempty" yaml:"PointInTimeRecoveryEnabled,omitempty"`
}

// SSESpecification represents server-side encryption settings.
type SSESpecification struct {
	// SSEEnabled indicates whether server-side encryption is enabled.
	SSEEnabled interface{} `json:"SSEEnabled,omitempty" yaml:"SSEEnabled,omitempty"`

	// SSEType specifies the server-side encryption type. Valid values: KMS
	SSEType string `json:"SSEType,omitempty" yaml:"SSEType,omitempty"`

	// KMSMasterKeyId specifies the AWS KMS key to use for encryption.
	KMSMasterKeyId interface{} `json:"KMSMasterKeyId,omitempty" yaml:"KMSMasterKeyId,omitempty"`
}

// TimeToLiveSpecification represents Time to Live settings.
type TimeToLiveSpecification struct {
	// AttributeName is the name of the TTL attribute.
	AttributeName interface{} `json:"AttributeName" yaml:"AttributeName"`

	// Enabled indicates whether TTL is enabled.
	Enabled interface{} `json:"Enabled" yaml:"Enabled"`
}

// Tag represents a key-value pair tag.
type Tag struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}

// ImportSourceSpecification represents settings for importing data from S3.
type ImportSourceSpecification struct {
	// S3BucketSource specifies the S3 bucket containing data to import.
	S3BucketSource S3BucketSource `json:"S3BucketSource" yaml:"S3BucketSource"`

	// InputFormat specifies the format of the source data. Valid values: CSV | DYNAMODB_JSON | ION
	InputFormat string `json:"InputFormat" yaml:"InputFormat"`

	// InputFormatOptions specifies additional options for the input format.
	InputFormatOptions *InputFormatOptions `json:"InputFormatOptions,omitempty" yaml:"InputFormatOptions,omitempty"`

	// InputCompressionType specifies the compression type. Valid values: GZIP | ZSTD | NONE
	InputCompressionType string `json:"InputCompressionType,omitempty" yaml:"InputCompressionType,omitempty"`
}

// S3BucketSource represents the S3 bucket source for import.
type S3BucketSource struct {
	// S3Bucket is the name of the S3 bucket containing the source data.
	S3Bucket interface{} `json:"S3Bucket" yaml:"S3Bucket"`

	// S3BucketOwner is the account number of the S3 bucket owner.
	S3BucketOwner interface{} `json:"S3BucketOwner,omitempty" yaml:"S3BucketOwner,omitempty"`

	// S3KeyPrefix is the key prefix of the S3 bucket.
	S3KeyPrefix interface{} `json:"S3KeyPrefix,omitempty" yaml:"S3KeyPrefix,omitempty"`
}

// InputFormatOptions represents format options for import source.
type InputFormatOptions struct {
	// Csv specifies options for CSV format.
	Csv *CsvOptions `json:"Csv,omitempty" yaml:"Csv,omitempty"`
}

// CsvOptions represents CSV format options.
type CsvOptions struct {
	// Delimiter is the delimiter character.
	Delimiter string `json:"Delimiter,omitempty" yaml:"Delimiter,omitempty"`

	// HeaderList is the list of header names.
	HeaderList []interface{} `json:"HeaderList,omitempty" yaml:"HeaderList,omitempty"`
}

// ResourcePolicy represents a resource-based policy.
type ResourcePolicy struct {
	// PolicyDocument is the resource-based policy document.
	PolicyDocument interface{} `json:"PolicyDocument" yaml:"PolicyDocument"`
}
