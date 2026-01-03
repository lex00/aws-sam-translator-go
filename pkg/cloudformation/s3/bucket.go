// Package s3 provides CloudFormation resource models for Amazon S3.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-s3-bucket.html
package s3

// Bucket represents an AWS::S3::Bucket resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-s3-bucket.html
type Bucket struct {
	// BucketName is the name of the bucket.
	BucketName interface{} `json:"BucketName,omitempty" yaml:"BucketName,omitempty"`

	// AccelerateConfiguration specifies the transfer acceleration configuration.
	AccelerateConfiguration *AccelerateConfiguration `json:"AccelerateConfiguration,omitempty" yaml:"AccelerateConfiguration,omitempty"`

	// AccessControl is the canned ACL for the bucket (deprecated).
	AccessControl string `json:"AccessControl,omitempty" yaml:"AccessControl,omitempty"`

	// AnalyticsConfigurations specifies the analytics configurations.
	AnalyticsConfigurations []AnalyticsConfiguration `json:"AnalyticsConfigurations,omitempty" yaml:"AnalyticsConfigurations,omitempty"`

	// BucketEncryption specifies the encryption configuration.
	BucketEncryption *BucketEncryption `json:"BucketEncryption,omitempty" yaml:"BucketEncryption,omitempty"`

	// CorsConfiguration specifies the CORS configuration.
	CorsConfiguration *CorsConfiguration `json:"CorsConfiguration,omitempty" yaml:"CorsConfiguration,omitempty"`

	// IntelligentTieringConfigurations specifies intelligent-tiering configurations.
	IntelligentTieringConfigurations []IntelligentTieringConfiguration `json:"IntelligentTieringConfigurations,omitempty" yaml:"IntelligentTieringConfigurations,omitempty"`

	// InventoryConfigurations specifies inventory configurations.
	InventoryConfigurations []InventoryConfiguration `json:"InventoryConfigurations,omitempty" yaml:"InventoryConfigurations,omitempty"`

	// LifecycleConfiguration specifies the lifecycle configuration.
	LifecycleConfiguration *LifecycleConfiguration `json:"LifecycleConfiguration,omitempty" yaml:"LifecycleConfiguration,omitempty"`

	// LoggingConfiguration specifies the logging configuration.
	LoggingConfiguration *LoggingConfiguration `json:"LoggingConfiguration,omitempty" yaml:"LoggingConfiguration,omitempty"`

	// MetricsConfigurations specifies metrics configurations.
	MetricsConfigurations []MetricsConfiguration `json:"MetricsConfigurations,omitempty" yaml:"MetricsConfigurations,omitempty"`

	// NotificationConfiguration specifies the notification configuration for events.
	NotificationConfiguration *NotificationConfiguration `json:"NotificationConfiguration,omitempty" yaml:"NotificationConfiguration,omitempty"`

	// ObjectLockConfiguration specifies the Object Lock configuration.
	ObjectLockConfiguration *ObjectLockConfiguration `json:"ObjectLockConfiguration,omitempty" yaml:"ObjectLockConfiguration,omitempty"`

	// ObjectLockEnabled indicates whether Object Lock is enabled.
	ObjectLockEnabled interface{} `json:"ObjectLockEnabled,omitempty" yaml:"ObjectLockEnabled,omitempty"`

	// OwnershipControls specifies the ownership controls.
	OwnershipControls *OwnershipControls `json:"OwnershipControls,omitempty" yaml:"OwnershipControls,omitempty"`

	// PublicAccessBlockConfiguration specifies public access block configuration.
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" yaml:"PublicAccessBlockConfiguration,omitempty"`

	// ReplicationConfiguration specifies the replication configuration.
	ReplicationConfiguration *ReplicationConfiguration `json:"ReplicationConfiguration,omitempty" yaml:"ReplicationConfiguration,omitempty"`

	// Tags is a list of key-value pairs to apply to the bucket.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// VersioningConfiguration specifies the versioning configuration.
	VersioningConfiguration *VersioningConfiguration `json:"VersioningConfiguration,omitempty" yaml:"VersioningConfiguration,omitempty"`

	// WebsiteConfiguration specifies the static website configuration.
	WebsiteConfiguration *WebsiteConfiguration `json:"WebsiteConfiguration,omitempty" yaml:"WebsiteConfiguration,omitempty"`
}

// AccelerateConfiguration specifies transfer acceleration configuration.
type AccelerateConfiguration struct {
	// AccelerationStatus indicates whether acceleration is enabled.
	// Valid values: Enabled | Suspended
	AccelerationStatus string `json:"AccelerationStatus" yaml:"AccelerationStatus"`
}

// AnalyticsConfiguration specifies analytics configuration.
type AnalyticsConfiguration struct {
	// Id is the ID of the analytics configuration.
	Id interface{} `json:"Id" yaml:"Id"`

	// Prefix is the prefix for filtering objects.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// StorageClassAnalysis specifies storage class analysis.
	StorageClassAnalysis StorageClassAnalysis `json:"StorageClassAnalysis" yaml:"StorageClassAnalysis"`

	// TagFilters specifies tag filters.
	TagFilters []TagFilter `json:"TagFilters,omitempty" yaml:"TagFilters,omitempty"`
}

// StorageClassAnalysis specifies storage class analysis configuration.
type StorageClassAnalysis struct {
	// DataExport specifies data export configuration.
	DataExport *DataExport `json:"DataExport,omitempty" yaml:"DataExport,omitempty"`
}

// DataExport specifies data export configuration.
type DataExport struct {
	// Destination specifies the export destination.
	Destination Destination `json:"Destination" yaml:"Destination"`

	// OutputSchemaVersion specifies the output schema version. Valid values: V_1
	OutputSchemaVersion string `json:"OutputSchemaVersion" yaml:"OutputSchemaVersion"`
}

// Destination specifies the analytics export destination.
type Destination struct {
	// BucketAccountId is the account ID of the destination bucket owner.
	BucketAccountId interface{} `json:"BucketAccountId,omitempty" yaml:"BucketAccountId,omitempty"`

	// BucketArn is the ARN of the destination bucket.
	BucketArn interface{} `json:"BucketArn" yaml:"BucketArn"`

	// Format is the output format. Valid values: CSV
	Format string `json:"Format" yaml:"Format"`

	// Prefix is the prefix for the destination.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`
}

// TagFilter specifies a tag filter.
type TagFilter struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}

// BucketEncryption specifies bucket encryption configuration.
type BucketEncryption struct {
	// ServerSideEncryptionConfiguration specifies server-side encryption rules.
	ServerSideEncryptionConfiguration []ServerSideEncryptionRule `json:"ServerSideEncryptionConfiguration" yaml:"ServerSideEncryptionConfiguration"`
}

// ServerSideEncryptionRule specifies a server-side encryption rule.
type ServerSideEncryptionRule struct {
	// BucketKeyEnabled indicates whether bucket key is enabled.
	BucketKeyEnabled interface{} `json:"BucketKeyEnabled,omitempty" yaml:"BucketKeyEnabled,omitempty"`

	// ServerSideEncryptionByDefault specifies default encryption.
	ServerSideEncryptionByDefault *ServerSideEncryptionByDefault `json:"ServerSideEncryptionByDefault,omitempty" yaml:"ServerSideEncryptionByDefault,omitempty"`
}

// ServerSideEncryptionByDefault specifies default server-side encryption.
type ServerSideEncryptionByDefault struct {
	// SSEAlgorithm is the encryption algorithm. Valid values: aws:kms | AES256 | aws:kms:dsse
	SSEAlgorithm string `json:"SSEAlgorithm" yaml:"SSEAlgorithm"`

	// KMSMasterKeyID is the AWS KMS key ID.
	KMSMasterKeyID interface{} `json:"KMSMasterKeyID,omitempty" yaml:"KMSMasterKeyID,omitempty"`
}

// CorsConfiguration specifies CORS configuration.
type CorsConfiguration struct {
	// CorsRules specifies the CORS rules.
	CorsRules []CorsRule `json:"CorsRules" yaml:"CorsRules"`
}

// CorsRule specifies a CORS rule.
type CorsRule struct {
	// Id is the ID of the CORS rule.
	Id interface{} `json:"Id,omitempty" yaml:"Id,omitempty"`

	// AllowedHeaders specifies allowed headers.
	AllowedHeaders []interface{} `json:"AllowedHeaders,omitempty" yaml:"AllowedHeaders,omitempty"`

	// AllowedMethods specifies allowed methods.
	AllowedMethods []string `json:"AllowedMethods" yaml:"AllowedMethods"`

	// AllowedOrigins specifies allowed origins.
	AllowedOrigins []interface{} `json:"AllowedOrigins" yaml:"AllowedOrigins"`

	// ExposedHeaders specifies exposed headers.
	ExposedHeaders []interface{} `json:"ExposedHeaders,omitempty" yaml:"ExposedHeaders,omitempty"`

	// MaxAge is the max age in seconds.
	MaxAge interface{} `json:"MaxAge,omitempty" yaml:"MaxAge,omitempty"`
}

// IntelligentTieringConfiguration specifies intelligent-tiering configuration.
type IntelligentTieringConfiguration struct {
	// Id is the ID of the configuration.
	Id interface{} `json:"Id" yaml:"Id"`

	// Prefix is the prefix filter.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// Status indicates whether the configuration is enabled. Valid values: Enabled | Disabled
	Status string `json:"Status" yaml:"Status"`

	// TagFilters specifies tag filters.
	TagFilters []TagFilter `json:"TagFilters,omitempty" yaml:"TagFilters,omitempty"`

	// Tierings specifies the tiering configuration.
	Tierings []Tiering `json:"Tierings" yaml:"Tierings"`
}

// Tiering specifies a tiering configuration.
type Tiering struct {
	// AccessTier is the access tier. Valid values: ARCHIVE_ACCESS | DEEP_ARCHIVE_ACCESS
	AccessTier string `json:"AccessTier" yaml:"AccessTier"`

	// Days is the number of days before transitioning.
	Days interface{} `json:"Days" yaml:"Days"`
}

// InventoryConfiguration specifies inventory configuration.
type InventoryConfiguration struct {
	// Id is the ID of the inventory configuration.
	Id interface{} `json:"Id" yaml:"Id"`

	// Destination specifies the inventory destination.
	Destination InventoryDestination `json:"Destination" yaml:"Destination"`

	// Enabled indicates whether the inventory is enabled.
	Enabled interface{} `json:"Enabled" yaml:"Enabled"`

	// IncludedObjectVersions specifies which object versions to include.
	// Valid values: All | Current
	IncludedObjectVersions string `json:"IncludedObjectVersions" yaml:"IncludedObjectVersions"`

	// OptionalFields specifies optional fields to include.
	OptionalFields []string `json:"OptionalFields,omitempty" yaml:"OptionalFields,omitempty"`

	// Prefix is the prefix filter.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// ScheduleFrequency is the frequency. Valid values: Daily | Weekly
	ScheduleFrequency string `json:"ScheduleFrequency" yaml:"ScheduleFrequency"`
}

// InventoryDestination specifies the inventory destination.
type InventoryDestination struct {
	// BucketAccountId is the destination bucket owner account ID.
	BucketAccountId interface{} `json:"BucketAccountId,omitempty" yaml:"BucketAccountId,omitempty"`

	// BucketArn is the destination bucket ARN.
	BucketArn interface{} `json:"BucketArn" yaml:"BucketArn"`

	// Format is the output format. Valid values: CSV | ORC | Parquet
	Format string `json:"Format" yaml:"Format"`

	// Prefix is the destination prefix.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`
}

// LifecycleConfiguration specifies lifecycle configuration.
type LifecycleConfiguration struct {
	// Rules specifies lifecycle rules.
	Rules []LifecycleRule `json:"Rules" yaml:"Rules"`
}

// LifecycleRule specifies a lifecycle rule.
type LifecycleRule struct {
	// Id is the rule ID.
	Id interface{} `json:"Id,omitempty" yaml:"Id,omitempty"`

	// Status indicates whether the rule is enabled. Valid values: Enabled | Disabled
	Status string `json:"Status" yaml:"Status"`

	// Prefix is the prefix filter (deprecated, use Filter).
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// Filter specifies the filter.
	Filter *LifecycleRuleFilter `json:"Filter,omitempty" yaml:"Filter,omitempty"`

	// AbortIncompleteMultipartUpload specifies abort settings.
	AbortIncompleteMultipartUpload *AbortIncompleteMultipartUpload `json:"AbortIncompleteMultipartUpload,omitempty" yaml:"AbortIncompleteMultipartUpload,omitempty"`

	// ExpirationDate specifies when objects expire.
	ExpirationDate interface{} `json:"ExpirationDate,omitempty" yaml:"ExpirationDate,omitempty"`

	// ExpirationInDays specifies the expiration in days.
	ExpirationInDays interface{} `json:"ExpirationInDays,omitempty" yaml:"ExpirationInDays,omitempty"`

	// ExpiredObjectDeleteMarker indicates whether to delete expired object delete markers.
	ExpiredObjectDeleteMarker interface{} `json:"ExpiredObjectDeleteMarker,omitempty" yaml:"ExpiredObjectDeleteMarker,omitempty"`

	// NoncurrentVersionExpiration specifies noncurrent version expiration.
	NoncurrentVersionExpiration *NoncurrentVersionExpiration `json:"NoncurrentVersionExpiration,omitempty" yaml:"NoncurrentVersionExpiration,omitempty"`

	// NoncurrentVersionExpirationInDays specifies noncurrent expiration days (deprecated).
	NoncurrentVersionExpirationInDays interface{} `json:"NoncurrentVersionExpirationInDays,omitempty" yaml:"NoncurrentVersionExpirationInDays,omitempty"`

	// NoncurrentVersionTransition specifies noncurrent version transition (deprecated).
	NoncurrentVersionTransition *NoncurrentVersionTransition `json:"NoncurrentVersionTransition,omitempty" yaml:"NoncurrentVersionTransition,omitempty"`

	// NoncurrentVersionTransitions specifies noncurrent version transitions.
	NoncurrentVersionTransitions []NoncurrentVersionTransition `json:"NoncurrentVersionTransitions,omitempty" yaml:"NoncurrentVersionTransitions,omitempty"`

	// ObjectSizeGreaterThan specifies the minimum object size.
	ObjectSizeGreaterThan interface{} `json:"ObjectSizeGreaterThan,omitempty" yaml:"ObjectSizeGreaterThan,omitempty"`

	// ObjectSizeLessThan specifies the maximum object size.
	ObjectSizeLessThan interface{} `json:"ObjectSizeLessThan,omitempty" yaml:"ObjectSizeLessThan,omitempty"`

	// Transition specifies transition (deprecated).
	Transition *Transition `json:"Transition,omitempty" yaml:"Transition,omitempty"`

	// Transitions specifies transitions.
	Transitions []Transition `json:"Transitions,omitempty" yaml:"Transitions,omitempty"`
}

// LifecycleRuleFilter specifies a lifecycle rule filter.
type LifecycleRuleFilter struct {
	// And specifies an AND filter.
	And *LifecycleRuleAndOperator `json:"And,omitempty" yaml:"And,omitempty"`

	// ObjectSizeGreaterThan specifies minimum size.
	ObjectSizeGreaterThan interface{} `json:"ObjectSizeGreaterThan,omitempty" yaml:"ObjectSizeGreaterThan,omitempty"`

	// ObjectSizeLessThan specifies maximum size.
	ObjectSizeLessThan interface{} `json:"ObjectSizeLessThan,omitempty" yaml:"ObjectSizeLessThan,omitempty"`

	// Prefix specifies the prefix.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// Tag specifies the tag filter.
	Tag *TagFilter `json:"Tag,omitempty" yaml:"Tag,omitempty"`
}

// LifecycleRuleAndOperator specifies an AND operator for filters.
type LifecycleRuleAndOperator struct {
	// ObjectSizeGreaterThan specifies minimum size.
	ObjectSizeGreaterThan interface{} `json:"ObjectSizeGreaterThan,omitempty" yaml:"ObjectSizeGreaterThan,omitempty"`

	// ObjectSizeLessThan specifies maximum size.
	ObjectSizeLessThan interface{} `json:"ObjectSizeLessThan,omitempty" yaml:"ObjectSizeLessThan,omitempty"`

	// Prefix specifies the prefix.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// Tags specifies the tags.
	Tags []TagFilter `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// AbortIncompleteMultipartUpload specifies when to abort multipart uploads.
type AbortIncompleteMultipartUpload struct {
	// DaysAfterInitiation is the number of days after initiation.
	DaysAfterInitiation interface{} `json:"DaysAfterInitiation" yaml:"DaysAfterInitiation"`
}

// NoncurrentVersionExpiration specifies noncurrent version expiration.
type NoncurrentVersionExpiration struct {
	// NewerNoncurrentVersions specifies the number of newer versions to retain.
	NewerNoncurrentVersions interface{} `json:"NewerNoncurrentVersions,omitempty" yaml:"NewerNoncurrentVersions,omitempty"`

	// NoncurrentDays specifies the number of days.
	NoncurrentDays interface{} `json:"NoncurrentDays" yaml:"NoncurrentDays"`
}

// NoncurrentVersionTransition specifies noncurrent version transition.
type NoncurrentVersionTransition struct {
	// NewerNoncurrentVersions specifies the number of newer versions to retain.
	NewerNoncurrentVersions interface{} `json:"NewerNoncurrentVersions,omitempty" yaml:"NewerNoncurrentVersions,omitempty"`

	// StorageClass specifies the storage class.
	StorageClass string `json:"StorageClass" yaml:"StorageClass"`

	// TransitionInDays specifies the transition days.
	TransitionInDays interface{} `json:"TransitionInDays" yaml:"TransitionInDays"`
}

// Transition specifies a transition.
type Transition struct {
	// StorageClass specifies the storage class.
	StorageClass string `json:"StorageClass" yaml:"StorageClass"`

	// TransitionDate specifies the transition date.
	TransitionDate interface{} `json:"TransitionDate,omitempty" yaml:"TransitionDate,omitempty"`

	// TransitionInDays specifies the transition days.
	TransitionInDays interface{} `json:"TransitionInDays,omitempty" yaml:"TransitionInDays,omitempty"`
}

// LoggingConfiguration specifies logging configuration.
type LoggingConfiguration struct {
	// DestinationBucketName is the destination bucket name.
	DestinationBucketName interface{} `json:"DestinationBucketName,omitempty" yaml:"DestinationBucketName,omitempty"`

	// LogFilePrefix is the log file prefix.
	LogFilePrefix interface{} `json:"LogFilePrefix,omitempty" yaml:"LogFilePrefix,omitempty"`

	// TargetObjectKeyFormat specifies the target key format.
	TargetObjectKeyFormat *TargetObjectKeyFormat `json:"TargetObjectKeyFormat,omitempty" yaml:"TargetObjectKeyFormat,omitempty"`
}

// TargetObjectKeyFormat specifies target object key format.
type TargetObjectKeyFormat struct {
	// PartitionedPrefix specifies partitioned prefix settings.
	PartitionedPrefix *PartitionedPrefix `json:"PartitionedPrefix,omitempty" yaml:"PartitionedPrefix,omitempty"`

	// SimplePrefix specifies simple prefix settings.
	SimplePrefix *SimplePrefix `json:"SimplePrefix,omitempty" yaml:"SimplePrefix,omitempty"`
}

// PartitionedPrefix specifies partitioned prefix settings.
type PartitionedPrefix struct {
	// PartitionDateSource specifies the date source. Valid values: EventTime | DeliveryTime
	PartitionDateSource string `json:"PartitionDateSource,omitempty" yaml:"PartitionDateSource,omitempty"`
}

// SimplePrefix is an empty struct for simple prefix format.
type SimplePrefix struct{}

// MetricsConfiguration specifies metrics configuration.
type MetricsConfiguration struct {
	// Id is the metrics configuration ID.
	Id interface{} `json:"Id" yaml:"Id"`

	// AccessPointArn specifies the access point ARN.
	AccessPointArn interface{} `json:"AccessPointArn,omitempty" yaml:"AccessPointArn,omitempty"`

	// Prefix is the prefix filter.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// TagFilters specifies tag filters.
	TagFilters []TagFilter `json:"TagFilters,omitempty" yaml:"TagFilters,omitempty"`
}

// NotificationConfiguration specifies notification configuration for bucket events.
// This is a key configuration for SAM event sources.
type NotificationConfiguration struct {
	// EventBridgeConfiguration specifies EventBridge configuration.
	EventBridgeConfiguration *EventBridgeConfiguration `json:"EventBridgeConfiguration,omitempty" yaml:"EventBridgeConfiguration,omitempty"`

	// LambdaConfigurations specifies Lambda function configurations.
	LambdaConfigurations []LambdaConfiguration `json:"LambdaConfigurations,omitempty" yaml:"LambdaConfigurations,omitempty"`

	// QueueConfigurations specifies SQS queue configurations.
	QueueConfigurations []QueueConfiguration `json:"QueueConfigurations,omitempty" yaml:"QueueConfigurations,omitempty"`

	// TopicConfigurations specifies SNS topic configurations.
	TopicConfigurations []TopicConfiguration `json:"TopicConfigurations,omitempty" yaml:"TopicConfigurations,omitempty"`
}

// EventBridgeConfiguration specifies EventBridge configuration.
type EventBridgeConfiguration struct {
	// EventBridgeEnabled indicates whether EventBridge is enabled.
	EventBridgeEnabled interface{} `json:"EventBridgeEnabled,omitempty" yaml:"EventBridgeEnabled,omitempty"`
}

// LambdaConfiguration specifies Lambda function notification configuration.
type LambdaConfiguration struct {
	// Event is the bucket event. Valid values include: s3:ObjectCreated:*, s3:ObjectRemoved:*, etc.
	Event string `json:"Event" yaml:"Event"`

	// Function is the ARN of the Lambda function.
	Function interface{} `json:"Function" yaml:"Function"`

	// Filter specifies the notification filter.
	Filter *NotificationFilter `json:"Filter,omitempty" yaml:"Filter,omitempty"`
}

// QueueConfiguration specifies SQS queue notification configuration.
type QueueConfiguration struct {
	// Event is the bucket event.
	Event string `json:"Event" yaml:"Event"`

	// Queue is the ARN of the SQS queue.
	Queue interface{} `json:"Queue" yaml:"Queue"`

	// Filter specifies the notification filter.
	Filter *NotificationFilter `json:"Filter,omitempty" yaml:"Filter,omitempty"`
}

// TopicConfiguration specifies SNS topic notification configuration.
type TopicConfiguration struct {
	// Event is the bucket event.
	Event string `json:"Event" yaml:"Event"`

	// Topic is the ARN of the SNS topic.
	Topic interface{} `json:"Topic" yaml:"Topic"`

	// Filter specifies the notification filter.
	Filter *NotificationFilter `json:"Filter,omitempty" yaml:"Filter,omitempty"`
}

// NotificationFilter specifies a notification filter.
type NotificationFilter struct {
	// S3Key specifies the S3 key filter.
	S3Key *S3KeyFilter `json:"S3Key,omitempty" yaml:"S3Key,omitempty"`
}

// S3KeyFilter specifies the S3 key filter rules.
type S3KeyFilter struct {
	// Rules specifies the filter rules.
	Rules []FilterRule `json:"Rules" yaml:"Rules"`
}

// FilterRule specifies a filter rule.
type FilterRule struct {
	// Name is the filter name. Valid values: prefix | suffix
	Name string `json:"Name" yaml:"Name"`

	// Value is the filter value.
	Value interface{} `json:"Value" yaml:"Value"`
}

// ObjectLockConfiguration specifies Object Lock configuration.
type ObjectLockConfiguration struct {
	// ObjectLockEnabled indicates whether Object Lock is enabled.
	ObjectLockEnabled string `json:"ObjectLockEnabled,omitempty" yaml:"ObjectLockEnabled,omitempty"`

	// Rule specifies the Object Lock rule.
	Rule *ObjectLockRule `json:"Rule,omitempty" yaml:"Rule,omitempty"`
}

// ObjectLockRule specifies an Object Lock rule.
type ObjectLockRule struct {
	// DefaultRetention specifies the default retention.
	DefaultRetention *DefaultRetention `json:"DefaultRetention,omitempty" yaml:"DefaultRetention,omitempty"`
}

// DefaultRetention specifies default retention settings.
type DefaultRetention struct {
	// Days is the number of days.
	Days interface{} `json:"Days,omitempty" yaml:"Days,omitempty"`

	// Mode is the retention mode. Valid values: GOVERNANCE | COMPLIANCE
	Mode string `json:"Mode,omitempty" yaml:"Mode,omitempty"`

	// Years is the number of years.
	Years interface{} `json:"Years,omitempty" yaml:"Years,omitempty"`
}

// OwnershipControls specifies ownership controls.
type OwnershipControls struct {
	// Rules specifies the ownership control rules.
	Rules []OwnershipControlsRule `json:"Rules" yaml:"Rules"`
}

// OwnershipControlsRule specifies an ownership control rule.
type OwnershipControlsRule struct {
	// ObjectOwnership specifies object ownership.
	// Valid values: BucketOwnerEnforced | ObjectWriter | BucketOwnerPreferred
	ObjectOwnership string `json:"ObjectOwnership" yaml:"ObjectOwnership"`
}

// PublicAccessBlockConfiguration specifies public access block configuration.
type PublicAccessBlockConfiguration struct {
	// BlockPublicAcls indicates whether to block public ACLs.
	BlockPublicAcls interface{} `json:"BlockPublicAcls,omitempty" yaml:"BlockPublicAcls,omitempty"`

	// BlockPublicPolicy indicates whether to block public bucket policies.
	BlockPublicPolicy interface{} `json:"BlockPublicPolicy,omitempty" yaml:"BlockPublicPolicy,omitempty"`

	// IgnorePublicAcls indicates whether to ignore public ACLs.
	IgnorePublicAcls interface{} `json:"IgnorePublicAcls,omitempty" yaml:"IgnorePublicAcls,omitempty"`

	// RestrictPublicBuckets indicates whether to restrict public buckets.
	RestrictPublicBuckets interface{} `json:"RestrictPublicBuckets,omitempty" yaml:"RestrictPublicBuckets,omitempty"`
}

// ReplicationConfiguration specifies replication configuration.
type ReplicationConfiguration struct {
	// Role is the IAM role ARN.
	Role interface{} `json:"Role" yaml:"Role"`

	// Rules specifies the replication rules.
	Rules []ReplicationRule `json:"Rules" yaml:"Rules"`
}

// ReplicationRule specifies a replication rule.
type ReplicationRule struct {
	// Id is the rule ID.
	Id interface{} `json:"Id,omitempty" yaml:"Id,omitempty"`

	// Destination specifies the destination.
	Destination ReplicationDestination `json:"Destination" yaml:"Destination"`

	// Filter specifies the filter.
	Filter *ReplicationRuleFilter `json:"Filter,omitempty" yaml:"Filter,omitempty"`

	// Prefix is the prefix filter (deprecated).
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// Priority is the rule priority.
	Priority interface{} `json:"Priority,omitempty" yaml:"Priority,omitempty"`

	// Status is the rule status. Valid values: Enabled | Disabled
	Status string `json:"Status" yaml:"Status"`

	// DeleteMarkerReplication specifies delete marker replication.
	DeleteMarkerReplication *DeleteMarkerReplication `json:"DeleteMarkerReplication,omitempty" yaml:"DeleteMarkerReplication,omitempty"`

	// SourceSelectionCriteria specifies source selection criteria.
	SourceSelectionCriteria *SourceSelectionCriteria `json:"SourceSelectionCriteria,omitempty" yaml:"SourceSelectionCriteria,omitempty"`
}

// ReplicationDestination specifies a replication destination.
type ReplicationDestination struct {
	// Bucket is the destination bucket ARN.
	Bucket interface{} `json:"Bucket" yaml:"Bucket"`

	// Account is the destination account ID.
	Account interface{} `json:"Account,omitempty" yaml:"Account,omitempty"`

	// AccessControlTranslation specifies ACL translation.
	AccessControlTranslation *AccessControlTranslation `json:"AccessControlTranslation,omitempty" yaml:"AccessControlTranslation,omitempty"`

	// EncryptionConfiguration specifies encryption configuration.
	EncryptionConfiguration *ReplicationEncryptionConfiguration `json:"EncryptionConfiguration,omitempty" yaml:"EncryptionConfiguration,omitempty"`

	// Metrics specifies metrics.
	Metrics *Metrics `json:"Metrics,omitempty" yaml:"Metrics,omitempty"`

	// ReplicationTime specifies replication time.
	ReplicationTime *ReplicationTime `json:"ReplicationTime,omitempty" yaml:"ReplicationTime,omitempty"`

	// StorageClass specifies the storage class.
	StorageClass string `json:"StorageClass,omitempty" yaml:"StorageClass,omitempty"`
}

// AccessControlTranslation specifies ACL translation.
type AccessControlTranslation struct {
	// Owner specifies the owner. Valid values: Destination
	Owner string `json:"Owner" yaml:"Owner"`
}

// ReplicationEncryptionConfiguration specifies replication encryption configuration.
type ReplicationEncryptionConfiguration struct {
	// ReplicaKmsKeyID specifies the replica KMS key ID.
	ReplicaKmsKeyID interface{} `json:"ReplicaKmsKeyID" yaml:"ReplicaKmsKeyID"`
}

// Metrics specifies replication metrics.
type Metrics struct {
	// EventThreshold specifies the event threshold.
	EventThreshold *ReplicationTimeValue `json:"EventThreshold,omitempty" yaml:"EventThreshold,omitempty"`

	// Status is the metrics status. Valid values: Enabled | Disabled
	Status string `json:"Status" yaml:"Status"`
}

// ReplicationTime specifies replication time control.
type ReplicationTime struct {
	// Status is the replication time status. Valid values: Enabled | Disabled
	Status string `json:"Status" yaml:"Status"`

	// Time specifies the time.
	Time ReplicationTimeValue `json:"Time" yaml:"Time"`
}

// ReplicationTimeValue specifies a replication time value.
type ReplicationTimeValue struct {
	// Minutes is the number of minutes.
	Minutes interface{} `json:"Minutes" yaml:"Minutes"`
}

// ReplicationRuleFilter specifies a replication rule filter.
type ReplicationRuleFilter struct {
	// And specifies an AND filter.
	And *ReplicationRuleAndOperator `json:"And,omitempty" yaml:"And,omitempty"`

	// Prefix specifies the prefix.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// Tag specifies the tag filter.
	Tag *TagFilter `json:"Tag,omitempty" yaml:"Tag,omitempty"`
}

// ReplicationRuleAndOperator specifies an AND operator.
type ReplicationRuleAndOperator struct {
	// Prefix specifies the prefix.
	Prefix interface{} `json:"Prefix,omitempty" yaml:"Prefix,omitempty"`

	// Tags specifies the tags.
	Tags []TagFilter `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// DeleteMarkerReplication specifies delete marker replication.
type DeleteMarkerReplication struct {
	// Status is the status. Valid values: Enabled | Disabled
	Status string `json:"Status,omitempty" yaml:"Status,omitempty"`
}

// SourceSelectionCriteria specifies source selection criteria.
type SourceSelectionCriteria struct {
	// ReplicaModifications specifies replica modifications.
	ReplicaModifications *ReplicaModifications `json:"ReplicaModifications,omitempty" yaml:"ReplicaModifications,omitempty"`

	// SseKmsEncryptedObjects specifies SSE-KMS encrypted objects.
	SseKmsEncryptedObjects *SseKmsEncryptedObjects `json:"SseKmsEncryptedObjects,omitempty" yaml:"SseKmsEncryptedObjects,omitempty"`
}

// ReplicaModifications specifies replica modifications.
type ReplicaModifications struct {
	// Status is the status. Valid values: Enabled | Disabled
	Status string `json:"Status" yaml:"Status"`
}

// SseKmsEncryptedObjects specifies SSE-KMS encrypted objects.
type SseKmsEncryptedObjects struct {
	// Status is the status. Valid values: Enabled | Disabled
	Status string `json:"Status" yaml:"Status"`
}

// VersioningConfiguration specifies versioning configuration.
type VersioningConfiguration struct {
	// Status is the versioning status. Valid values: Enabled | Suspended
	Status string `json:"Status" yaml:"Status"`
}

// WebsiteConfiguration specifies static website configuration.
type WebsiteConfiguration struct {
	// ErrorDocument specifies the error document.
	ErrorDocument *ErrorDocument `json:"ErrorDocument,omitempty" yaml:"ErrorDocument,omitempty"`

	// IndexDocument specifies the index document.
	IndexDocument *IndexDocument `json:"IndexDocument,omitempty" yaml:"IndexDocument,omitempty"`

	// RedirectAllRequestsTo specifies redirect configuration.
	RedirectAllRequestsTo *RedirectAllRequestsTo `json:"RedirectAllRequestsTo,omitempty" yaml:"RedirectAllRequestsTo,omitempty"`

	// RoutingRules specifies routing rules.
	RoutingRules []RoutingRule `json:"RoutingRules,omitempty" yaml:"RoutingRules,omitempty"`
}

// ErrorDocument specifies the error document.
type ErrorDocument struct {
	// Key is the error document key.
	Key interface{} `json:"Key" yaml:"Key"`
}

// IndexDocument specifies the index document.
type IndexDocument struct {
	// Suffix is the index document suffix.
	Suffix interface{} `json:"Suffix" yaml:"Suffix"`
}

// RedirectAllRequestsTo specifies redirect configuration.
type RedirectAllRequestsTo struct {
	// HostName is the redirect hostname.
	HostName interface{} `json:"HostName" yaml:"HostName"`

	// Protocol is the redirect protocol.
	Protocol string `json:"Protocol,omitempty" yaml:"Protocol,omitempty"`
}

// RoutingRule specifies a routing rule.
type RoutingRule struct {
	// Condition specifies the condition.
	Condition *RoutingRuleCondition `json:"Condition,omitempty" yaml:"Condition,omitempty"`

	// Redirect specifies the redirect.
	Redirect Redirect `json:"Redirect" yaml:"Redirect"`
}

// RoutingRuleCondition specifies a routing rule condition.
type RoutingRuleCondition struct {
	// HttpErrorCodeReturnedEquals specifies the HTTP error code.
	HttpErrorCodeReturnedEquals interface{} `json:"HttpErrorCodeReturnedEquals,omitempty" yaml:"HttpErrorCodeReturnedEquals,omitempty"`

	// KeyPrefixEquals specifies the key prefix.
	KeyPrefixEquals interface{} `json:"KeyPrefixEquals,omitempty" yaml:"KeyPrefixEquals,omitempty"`
}

// Redirect specifies a redirect.
type Redirect struct {
	// HostName is the redirect hostname.
	HostName interface{} `json:"HostName,omitempty" yaml:"HostName,omitempty"`

	// HttpRedirectCode is the HTTP redirect code.
	HttpRedirectCode interface{} `json:"HttpRedirectCode,omitempty" yaml:"HttpRedirectCode,omitempty"`

	// Protocol is the redirect protocol.
	Protocol string `json:"Protocol,omitempty" yaml:"Protocol,omitempty"`

	// ReplaceKeyPrefixWith specifies key prefix replacement.
	ReplaceKeyPrefixWith interface{} `json:"ReplaceKeyPrefixWith,omitempty" yaml:"ReplaceKeyPrefixWith,omitempty"`

	// ReplaceKeyWith specifies key replacement.
	ReplaceKeyWith interface{} `json:"ReplaceKeyWith,omitempty" yaml:"ReplaceKeyWith,omitempty"`
}

// Tag represents a key-value pair tag.
type Tag struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}
