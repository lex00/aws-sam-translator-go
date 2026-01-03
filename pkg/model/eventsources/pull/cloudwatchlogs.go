package pull

// CloudWatchLogsEventProperties represents SAM properties for a CloudWatch Logs event source.
// This event type creates an AWS::Logs::SubscriptionFilter resource, not an EventSourceMapping.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-cloudwatchlogs.html
type CloudWatchLogsEventProperties struct {
	// LogGroupName is the name of the CloudWatch Logs log group (required).
	LogGroupName interface{} `json:"LogGroupName" yaml:"LogGroupName"`

	// FilterPattern is the filtering expression for the subscription (required).
	// An empty string matches all log events.
	FilterPattern string `json:"FilterPattern" yaml:"FilterPattern"`
}

// NewCloudWatchLogsEventProperties creates a new CloudWatchLogsEventProperties with required fields.
func NewCloudWatchLogsEventProperties(logGroupName interface{}, filterPattern string) *CloudWatchLogsEventProperties {
	return &CloudWatchLogsEventProperties{
		LogGroupName:  logGroupName,
		FilterPattern: filterPattern,
	}
}

// SubscriptionFilter represents an AWS::Logs::SubscriptionFilter CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-logs-subscriptionfilter.html
type SubscriptionFilter struct {
	// DestinationArn is the ARN of the destination resource (Lambda function, Kinesis stream, etc.).
	DestinationArn interface{} `json:"DestinationArn" yaml:"DestinationArn"`

	// FilterPattern is the filtering expression.
	FilterPattern string `json:"FilterPattern" yaml:"FilterPattern"`

	// LogGroupName is the name of the log group to associate with the subscription filter.
	LogGroupName interface{} `json:"LogGroupName" yaml:"LogGroupName"`

	// FilterName is the name of the subscription filter.
	FilterName interface{} `json:"FilterName,omitempty" yaml:"FilterName,omitempty"`

	// RoleArn is the ARN of the IAM role for cross-account log delivery.
	RoleArn interface{} `json:"RoleArn,omitempty" yaml:"RoleArn,omitempty"`

	// Distribution is the method for distributing log data to the destination.
	// Valid values: Random, ByLogStream
	Distribution string `json:"Distribution,omitempty" yaml:"Distribution,omitempty"`
}

// ResourceTypeSubscriptionFilter is the CloudFormation resource type for subscription filters.
const ResourceTypeSubscriptionFilter = "AWS::Logs::SubscriptionFilter"

// ToSubscriptionFilter converts CloudWatchLogsEventProperties to a SubscriptionFilter.
func (c *CloudWatchLogsEventProperties) ToSubscriptionFilter(destinationArn interface{}) *SubscriptionFilter {
	return &SubscriptionFilter{
		DestinationArn: destinationArn,
		FilterPattern:  c.FilterPattern,
		LogGroupName:   c.LogGroupName,
	}
}

// WithFilterName sets the filter name.
func (s *SubscriptionFilter) WithFilterName(name interface{}) *SubscriptionFilter {
	s.FilterName = name
	return s
}

// WithRoleArn sets the IAM role ARN for cross-account delivery.
func (s *SubscriptionFilter) WithRoleArn(roleArn interface{}) *SubscriptionFilter {
	s.RoleArn = roleArn
	return s
}

// WithDistribution sets the distribution method.
func (s *SubscriptionFilter) WithDistribution(distribution string) *SubscriptionFilter {
	s.Distribution = distribution
	return s
}

// ToCloudFormation converts the SubscriptionFilter to a CloudFormation resource.
func (s *SubscriptionFilter) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	properties["DestinationArn"] = s.DestinationArn
	properties["FilterPattern"] = s.FilterPattern
	properties["LogGroupName"] = s.LogGroupName

	if s.FilterName != nil {
		properties["FilterName"] = s.FilterName
	}
	if s.RoleArn != nil {
		properties["RoleArn"] = s.RoleArn
	}
	if s.Distribution != "" {
		properties["Distribution"] = s.Distribution
	}

	return map[string]interface{}{
		"Type":       ResourceTypeSubscriptionFilter,
		"Properties": properties,
	}
}
