// Package events provides CloudFormation resource models for Amazon EventBridge.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-events-rule.html
package events

// Rule represents an AWS::Events::Rule resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-events-rule.html
type Rule struct {
	// Name is the name of the rule. If not specified, CloudFormation generates a unique name.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Description is the description of the rule.
	Description interface{} `json:"Description,omitempty" yaml:"Description,omitempty"`

	// EventBusName is the name or ARN of the event bus to associate with this rule.
	EventBusName interface{} `json:"EventBusName,omitempty" yaml:"EventBusName,omitempty"`

	// EventPattern is the event pattern that describes the event structure.
	EventPattern interface{} `json:"EventPattern,omitempty" yaml:"EventPattern,omitempty"`

	// ScheduleExpression is the scheduling expression (cron or rate).
	ScheduleExpression interface{} `json:"ScheduleExpression,omitempty" yaml:"ScheduleExpression,omitempty"`

	// State indicates whether the rule is enabled. Valid values: DISABLED | ENABLED | ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS
	State string `json:"State,omitempty" yaml:"State,omitempty"`

	// Targets is the list of targets for the rule.
	Targets []Target `json:"Targets,omitempty" yaml:"Targets,omitempty"`

	// RoleArn is the ARN of the IAM role associated with the rule.
	RoleArn interface{} `json:"RoleArn,omitempty" yaml:"RoleArn,omitempty"`
}

// Target represents a target for an EventBridge rule.
type Target struct {
	// Id is a unique identifier for the target.
	Id interface{} `json:"Id" yaml:"Id"`

	// Arn is the ARN of the target resource.
	Arn interface{} `json:"Arn" yaml:"Arn"`

	// RoleArn is the ARN of the IAM role for cross-account event delivery.
	RoleArn interface{} `json:"RoleArn,omitempty" yaml:"RoleArn,omitempty"`

	// Input is the JSON text to pass to the target.
	Input interface{} `json:"Input,omitempty" yaml:"Input,omitempty"`

	// InputPath is the JSONPath to extract from the event to pass to the target.
	InputPath interface{} `json:"InputPath,omitempty" yaml:"InputPath,omitempty"`

	// InputTransformer specifies settings for transforming input before passing to the target.
	InputTransformer *InputTransformer `json:"InputTransformer,omitempty" yaml:"InputTransformer,omitempty"`

	// BatchParameters specifies the job definition, job name, and other parameters.
	BatchParameters *BatchParameters `json:"BatchParameters,omitempty" yaml:"BatchParameters,omitempty"`

	// DeadLetterConfig specifies the dead-letter queue configuration.
	DeadLetterConfig *DeadLetterConfig `json:"DeadLetterConfig,omitempty" yaml:"DeadLetterConfig,omitempty"`

	// EcsParameters specifies ECS task parameters.
	EcsParameters *EcsParameters `json:"EcsParameters,omitempty" yaml:"EcsParameters,omitempty"`

	// HttpParameters specifies HTTP parameters for API destination targets.
	HttpParameters *HttpParameters `json:"HttpParameters,omitempty" yaml:"HttpParameters,omitempty"`

	// KinesisParameters specifies Kinesis stream parameters.
	KinesisParameters *KinesisParameters `json:"KinesisParameters,omitempty" yaml:"KinesisParameters,omitempty"`

	// RedshiftDataParameters specifies Amazon Redshift Data API parameters.
	RedshiftDataParameters *RedshiftDataParameters `json:"RedshiftDataParameters,omitempty" yaml:"RedshiftDataParameters,omitempty"`

	// RetryPolicy specifies the retry policy settings.
	RetryPolicy *RetryPolicy `json:"RetryPolicy,omitempty" yaml:"RetryPolicy,omitempty"`

	// RunCommandParameters specifies parameters for Systems Manager run command.
	RunCommandParameters *RunCommandParameters `json:"RunCommandParameters,omitempty" yaml:"RunCommandParameters,omitempty"`

	// SageMakerPipelineParameters specifies SageMaker Model Building Pipeline parameters.
	SageMakerPipelineParameters *SageMakerPipelineParameters `json:"SageMakerPipelineParameters,omitempty" yaml:"SageMakerPipelineParameters,omitempty"`

	// SqsParameters specifies SQS queue parameters.
	SqsParameters *SqsParameters `json:"SqsParameters,omitempty" yaml:"SqsParameters,omitempty"`

	// AppSyncParameters specifies AWS AppSync parameters.
	AppSyncParameters *AppSyncParameters `json:"AppSyncParameters,omitempty" yaml:"AppSyncParameters,omitempty"`
}

// InputTransformer specifies settings for input transformation.
type InputTransformer struct {
	// InputPathsMap is a map of JSON paths to extract from the event.
	InputPathsMap map[string]interface{} `json:"InputPathsMap,omitempty" yaml:"InputPathsMap,omitempty"`

	// InputTemplate is the template to use for the transformed input.
	InputTemplate interface{} `json:"InputTemplate" yaml:"InputTemplate"`
}

// BatchParameters specifies AWS Batch job parameters.
type BatchParameters struct {
	// JobDefinition is the ARN or name of the job definition.
	JobDefinition interface{} `json:"JobDefinition" yaml:"JobDefinition"`

	// JobName is the name of the job.
	JobName interface{} `json:"JobName" yaml:"JobName"`

	// ArrayProperties specifies the array properties of the job.
	ArrayProperties *BatchArrayProperties `json:"ArrayProperties,omitempty" yaml:"ArrayProperties,omitempty"`

	// RetryStrategy specifies the retry strategy.
	RetryStrategy *BatchRetryStrategy `json:"RetryStrategy,omitempty" yaml:"RetryStrategy,omitempty"`
}

// BatchArrayProperties specifies array properties for a batch job.
type BatchArrayProperties struct {
	// Size is the size of the array job.
	Size interface{} `json:"Size,omitempty" yaml:"Size,omitempty"`
}

// BatchRetryStrategy specifies the retry strategy for a batch job.
type BatchRetryStrategy struct {
	// Attempts is the number of times to attempt to retry the job.
	Attempts interface{} `json:"Attempts,omitempty" yaml:"Attempts,omitempty"`
}

// DeadLetterConfig specifies dead-letter queue configuration.
type DeadLetterConfig struct {
	// Arn is the ARN of the SQS queue to use as the dead-letter queue.
	Arn interface{} `json:"Arn,omitempty" yaml:"Arn,omitempty"`
}

// EcsParameters specifies ECS task parameters.
type EcsParameters struct {
	// TaskDefinitionArn is the ARN of the task definition.
	TaskDefinitionArn interface{} `json:"TaskDefinitionArn" yaml:"TaskDefinitionArn"`

	// TaskCount is the number of tasks to create.
	TaskCount interface{} `json:"TaskCount,omitempty" yaml:"TaskCount,omitempty"`

	// LaunchType is the launch type. Valid values: EC2 | FARGATE | EXTERNAL
	LaunchType string `json:"LaunchType,omitempty" yaml:"LaunchType,omitempty"`

	// NetworkConfiguration specifies the network configuration.
	NetworkConfiguration *NetworkConfiguration `json:"NetworkConfiguration,omitempty" yaml:"NetworkConfiguration,omitempty"`

	// PlatformVersion is the platform version for the task.
	PlatformVersion interface{} `json:"PlatformVersion,omitempty" yaml:"PlatformVersion,omitempty"`

	// Group is the task group.
	Group interface{} `json:"Group,omitempty" yaml:"Group,omitempty"`

	// CapacityProviderStrategy specifies the capacity provider strategy.
	CapacityProviderStrategy []CapacityProviderStrategyItem `json:"CapacityProviderStrategy,omitempty" yaml:"CapacityProviderStrategy,omitempty"`

	// EnableECSManagedTags indicates whether to enable ECS managed tags.
	EnableECSManagedTags interface{} `json:"EnableECSManagedTags,omitempty" yaml:"EnableECSManagedTags,omitempty"`

	// EnableExecuteCommand indicates whether to enable execute command.
	EnableExecuteCommand interface{} `json:"EnableExecuteCommand,omitempty" yaml:"EnableExecuteCommand,omitempty"`

	// PlacementConstraints specifies placement constraints.
	PlacementConstraints []PlacementConstraint `json:"PlacementConstraints,omitempty" yaml:"PlacementConstraints,omitempty"`

	// PlacementStrategies specifies placement strategies.
	PlacementStrategies []PlacementStrategy `json:"PlacementStrategies,omitempty" yaml:"PlacementStrategies,omitempty"`

	// PropagateTags specifies tag propagation. Valid values: TASK_DEFINITION
	PropagateTags string `json:"PropagateTags,omitempty" yaml:"PropagateTags,omitempty"`

	// ReferenceId is a reference ID for the task.
	ReferenceId interface{} `json:"ReferenceId,omitempty" yaml:"ReferenceId,omitempty"`

	// Tags is a list of tags to apply to the task.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`
}

// NetworkConfiguration specifies network configuration for ECS tasks.
type NetworkConfiguration struct {
	// AwsVpcConfiguration specifies the VPC configuration.
	AwsVpcConfiguration *AwsVpcConfiguration `json:"AwsVpcConfiguration,omitempty" yaml:"AwsVpcConfiguration,omitempty"`
}

// AwsVpcConfiguration specifies VPC configuration.
type AwsVpcConfiguration struct {
	// Subnets is the list of subnet IDs.
	Subnets []interface{} `json:"Subnets" yaml:"Subnets"`

	// SecurityGroups is the list of security group IDs.
	SecurityGroups []interface{} `json:"SecurityGroups,omitempty" yaml:"SecurityGroups,omitempty"`

	// AssignPublicIp indicates whether to assign a public IP. Valid values: DISABLED | ENABLED
	AssignPublicIp string `json:"AssignPublicIp,omitempty" yaml:"AssignPublicIp,omitempty"`
}

// CapacityProviderStrategyItem specifies a capacity provider strategy item.
type CapacityProviderStrategyItem struct {
	// CapacityProvider is the short name of the capacity provider.
	CapacityProvider interface{} `json:"CapacityProvider" yaml:"CapacityProvider"`

	// Weight is the relative percentage of total tasks.
	Weight interface{} `json:"Weight,omitempty" yaml:"Weight,omitempty"`

	// Base is the base number of tasks.
	Base interface{} `json:"Base,omitempty" yaml:"Base,omitempty"`
}

// PlacementConstraint specifies a placement constraint.
type PlacementConstraint struct {
	// Type is the type of constraint. Valid values: distinctInstance | memberOf
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`

	// Expression is the cluster query language expression.
	Expression interface{} `json:"Expression,omitempty" yaml:"Expression,omitempty"`
}

// PlacementStrategy specifies a placement strategy.
type PlacementStrategy struct {
	// Type is the type of strategy. Valid values: random | spread | binpack
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`

	// Field is the field to apply the strategy against.
	Field interface{} `json:"Field,omitempty" yaml:"Field,omitempty"`
}

// Tag represents a key-value pair tag.
type Tag struct {
	// Key is the tag key.
	Key interface{} `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value interface{} `json:"Value" yaml:"Value"`
}

// HttpParameters specifies HTTP parameters for API destination targets.
type HttpParameters struct {
	// PathParameterValues is a list of path parameter values.
	PathParameterValues []interface{} `json:"PathParameterValues,omitempty" yaml:"PathParameterValues,omitempty"`

	// HeaderParameters is a map of header parameters.
	HeaderParameters map[string]interface{} `json:"HeaderParameters,omitempty" yaml:"HeaderParameters,omitempty"`

	// QueryStringParameters is a map of query string parameters.
	QueryStringParameters map[string]interface{} `json:"QueryStringParameters,omitempty" yaml:"QueryStringParameters,omitempty"`
}

// KinesisParameters specifies Kinesis stream parameters.
type KinesisParameters struct {
	// PartitionKeyPath is the JSON path to the partition key.
	PartitionKeyPath interface{} `json:"PartitionKeyPath" yaml:"PartitionKeyPath"`
}

// RedshiftDataParameters specifies Amazon Redshift Data API parameters.
type RedshiftDataParameters struct {
	// Database is the name of the database.
	Database interface{} `json:"Database" yaml:"Database"`

	// DbUser is the database user name.
	DbUser interface{} `json:"DbUser,omitempty" yaml:"DbUser,omitempty"`

	// SecretManagerArn is the ARN of the secret containing credentials.
	SecretManagerArn interface{} `json:"SecretManagerArn,omitempty" yaml:"SecretManagerArn,omitempty"`

	// Sql is the SQL statement text.
	Sql interface{} `json:"Sql,omitempty" yaml:"Sql,omitempty"`

	// Sqls is a list of SQL statements.
	Sqls []interface{} `json:"Sqls,omitempty" yaml:"Sqls,omitempty"`

	// StatementName is the name of the SQL statement.
	StatementName interface{} `json:"StatementName,omitempty" yaml:"StatementName,omitempty"`

	// WithEvent indicates whether to send an event back.
	WithEvent interface{} `json:"WithEvent,omitempty" yaml:"WithEvent,omitempty"`
}

// RetryPolicy specifies retry policy settings.
type RetryPolicy struct {
	// MaximumEventAgeInSeconds is the maximum age of an event.
	MaximumEventAgeInSeconds interface{} `json:"MaximumEventAgeInSeconds,omitempty" yaml:"MaximumEventAgeInSeconds,omitempty"`

	// MaximumRetryAttempts is the maximum number of retry attempts.
	MaximumRetryAttempts interface{} `json:"MaximumRetryAttempts,omitempty" yaml:"MaximumRetryAttempts,omitempty"`
}

// RunCommandParameters specifies parameters for Systems Manager run command.
type RunCommandParameters struct {
	// RunCommandTargets is a list of run command targets.
	RunCommandTargets []RunCommandTarget `json:"RunCommandTargets" yaml:"RunCommandTargets"`
}

// RunCommandTarget specifies a run command target.
type RunCommandTarget struct {
	// Key is the tag key or InstanceIds.
	Key interface{} `json:"Key" yaml:"Key"`

	// Values is a list of tag values or instance IDs.
	Values []interface{} `json:"Values" yaml:"Values"`
}

// SageMakerPipelineParameters specifies SageMaker Model Building Pipeline parameters.
type SageMakerPipelineParameters struct {
	// PipelineParameterList is a list of pipeline parameters.
	PipelineParameterList []SageMakerPipelineParameter `json:"PipelineParameterList,omitempty" yaml:"PipelineParameterList,omitempty"`
}

// SageMakerPipelineParameter specifies a SageMaker pipeline parameter.
type SageMakerPipelineParameter struct {
	// Name is the name of the parameter.
	Name interface{} `json:"Name" yaml:"Name"`

	// Value is the value of the parameter.
	Value interface{} `json:"Value" yaml:"Value"`
}

// SqsParameters specifies SQS queue parameters.
type SqsParameters struct {
	// MessageGroupId is the FIFO message group ID.
	MessageGroupId interface{} `json:"MessageGroupId,omitempty" yaml:"MessageGroupId,omitempty"`
}

// AppSyncParameters specifies AWS AppSync parameters.
type AppSyncParameters struct {
	// GraphQLOperation is the GraphQL operation.
	GraphQLOperation interface{} `json:"GraphQLOperation,omitempty" yaml:"GraphQLOperation,omitempty"`
}
