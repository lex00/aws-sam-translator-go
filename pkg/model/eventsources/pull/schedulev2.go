package pull

// ScheduleV2EventProperties represents SAM properties for an EventBridge Scheduler event source.
// This event type creates an AWS::Scheduler::Schedule resource, not an EventSourceMapping.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-schedulev2.html
type ScheduleV2EventProperties struct {
	// ScheduleExpression is the scheduling expression (required).
	// Can be a cron expression or rate expression.
	ScheduleExpression string `json:"ScheduleExpression" yaml:"ScheduleExpression"`

	// ScheduleExpressionTimezone is the timezone for the schedule expression.
	// Defaults to UTC.
	ScheduleExpressionTimezone string `json:"ScheduleExpressionTimezone,omitempty" yaml:"ScheduleExpressionTimezone,omitempty"`

	// FlexibleTimeWindow configures the flexible time window for the schedule.
	FlexibleTimeWindow *FlexibleTimeWindow `json:"FlexibleTimeWindow,omitempty" yaml:"FlexibleTimeWindow,omitempty"`

	// Name is the name of the schedule.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Description is the description of the schedule.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// State indicates whether the schedule is enabled.
	// Valid values: ENABLED, DISABLED
	State string `json:"State,omitempty" yaml:"State,omitempty"`

	// GroupName is the name of the schedule group.
	GroupName interface{} `json:"GroupName,omitempty" yaml:"GroupName,omitempty"`

	// StartDate is the date when the schedule starts.
	StartDate string `json:"StartDate,omitempty" yaml:"StartDate,omitempty"`

	// EndDate is the date when the schedule ends.
	EndDate string `json:"EndDate,omitempty" yaml:"EndDate,omitempty"`

	// Input is the JSON text to pass to the target.
	Input interface{} `json:"Input,omitempty" yaml:"Input,omitempty"`

	// RetryPolicy configures retry behavior.
	RetryPolicy *ScheduleRetryPolicy `json:"RetryPolicy,omitempty" yaml:"RetryPolicy,omitempty"`

	// DeadLetterConfig configures the dead-letter queue.
	DeadLetterConfig *ScheduleDeadLetterConfig `json:"DeadLetterConfig,omitempty" yaml:"DeadLetterConfig,omitempty"`

	// KmsKeyArn is the ARN of the KMS key to encrypt the schedule.
	KmsKeyArn interface{} `json:"KmsKeyArn,omitempty" yaml:"KmsKeyArn,omitempty"`

	// RoleArn is the ARN of the IAM role for the schedule.
	RoleArn interface{} `json:"RoleArn,omitempty" yaml:"RoleArn,omitempty"`
}

// FlexibleTimeWindow configures the flexible time window for a schedule.
type FlexibleTimeWindow struct {
	// Mode is the flexible time window mode.
	// Valid values: OFF, FLEXIBLE
	Mode string `json:"Mode" yaml:"Mode"`

	// MaximumWindowInMinutes is the maximum time window (1-1440).
	// Only valid when Mode is FLEXIBLE.
	MaximumWindowInMinutes *int `json:"MaximumWindowInMinutes,omitempty" yaml:"MaximumWindowInMinutes,omitempty"`
}

// ScheduleRetryPolicy configures retry behavior for a schedule.
type ScheduleRetryPolicy struct {
	// MaximumEventAgeInSeconds is the maximum age of an event before it's discarded (60-86400).
	MaximumEventAgeInSeconds *int `json:"MaximumEventAgeInSeconds,omitempty" yaml:"MaximumEventAgeInSeconds,omitempty"`

	// MaximumRetryAttempts is the maximum number of retry attempts (0-185).
	MaximumRetryAttempts *int `json:"MaximumRetryAttempts,omitempty" yaml:"MaximumRetryAttempts,omitempty"`
}

// ScheduleDeadLetterConfig configures the dead-letter queue for a schedule.
type ScheduleDeadLetterConfig struct {
	// Arn is the ARN of the SQS queue to use as the dead-letter queue.
	Arn interface{} `json:"Arn,omitempty" yaml:"Arn,omitempty"`
}

// NewScheduleV2EventProperties creates a new ScheduleV2EventProperties with required fields.
func NewScheduleV2EventProperties(scheduleExpression string) *ScheduleV2EventProperties {
	return &ScheduleV2EventProperties{
		ScheduleExpression: scheduleExpression,
	}
}

// Schedule represents an AWS::Scheduler::Schedule CloudFormation resource.
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-scheduler-schedule.html
type Schedule struct {
	// ScheduleExpression is the scheduling expression (required).
	ScheduleExpression string `json:"ScheduleExpression" yaml:"ScheduleExpression"`

	// ScheduleExpressionTimezone is the timezone for the schedule expression.
	ScheduleExpressionTimezone string `json:"ScheduleExpressionTimezone,omitempty" yaml:"ScheduleExpressionTimezone,omitempty"`

	// FlexibleTimeWindow configures the flexible time window (required).
	FlexibleTimeWindow *FlexibleTimeWindow `json:"FlexibleTimeWindow" yaml:"FlexibleTimeWindow"`

	// Target specifies the target for the schedule (required).
	Target *ScheduleTarget `json:"Target" yaml:"Target"`

	// Name is the name of the schedule.
	Name interface{} `json:"Name,omitempty" yaml:"Name,omitempty"`

	// Description is the description of the schedule.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// State indicates whether the schedule is enabled.
	State string `json:"State,omitempty" yaml:"State,omitempty"`

	// GroupName is the name of the schedule group.
	GroupName interface{} `json:"GroupName,omitempty" yaml:"GroupName,omitempty"`

	// StartDate is the date when the schedule starts.
	StartDate string `json:"StartDate,omitempty" yaml:"StartDate,omitempty"`

	// EndDate is the date when the schedule ends.
	EndDate string `json:"EndDate,omitempty" yaml:"EndDate,omitempty"`

	// KmsKeyArn is the ARN of the KMS key to encrypt the schedule.
	KmsKeyArn interface{} `json:"KmsKeyArn,omitempty" yaml:"KmsKeyArn,omitempty"`
}

// ScheduleTarget specifies the target for a schedule.
type ScheduleTarget struct {
	// Arn is the ARN of the target resource (required).
	Arn interface{} `json:"Arn" yaml:"Arn"`

	// RoleArn is the ARN of the IAM role for the target (required).
	RoleArn interface{} `json:"RoleArn" yaml:"RoleArn"`

	// Input is the JSON text to pass to the target.
	Input interface{} `json:"Input,omitempty" yaml:"Input,omitempty"`

	// RetryPolicy configures retry behavior.
	RetryPolicy *ScheduleRetryPolicy `json:"RetryPolicy,omitempty" yaml:"RetryPolicy,omitempty"`

	// DeadLetterConfig configures the dead-letter queue.
	DeadLetterConfig *ScheduleDeadLetterConfig `json:"DeadLetterConfig,omitempty" yaml:"DeadLetterConfig,omitempty"`
}

// ResourceTypeSchedule is the CloudFormation resource type for EventBridge Scheduler schedules.
const ResourceTypeSchedule = "AWS::Scheduler::Schedule"

// ToSchedule converts ScheduleV2EventProperties to a Schedule.
func (s *ScheduleV2EventProperties) ToSchedule(targetArn interface{}, roleArn interface{}) *Schedule {
	schedule := &Schedule{
		ScheduleExpression: s.ScheduleExpression,
		FlexibleTimeWindow: &FlexibleTimeWindow{
			Mode: "OFF",
		},
		Target: &ScheduleTarget{
			Arn:     targetArn,
			RoleArn: roleArn,
		},
	}

	if s.ScheduleExpressionTimezone != "" {
		schedule.ScheduleExpressionTimezone = s.ScheduleExpressionTimezone
	}
	if s.FlexibleTimeWindow != nil {
		schedule.FlexibleTimeWindow = s.FlexibleTimeWindow
	}
	if s.Name != nil {
		schedule.Name = s.Name
	}
	if s.Description != "" {
		schedule.Description = s.Description
	}
	if s.State != "" {
		schedule.State = s.State
	}
	if s.GroupName != nil {
		schedule.GroupName = s.GroupName
	}
	if s.StartDate != "" {
		schedule.StartDate = s.StartDate
	}
	if s.EndDate != "" {
		schedule.EndDate = s.EndDate
	}
	if s.Input != nil {
		schedule.Target.Input = s.Input
	}
	if s.RetryPolicy != nil {
		schedule.Target.RetryPolicy = s.RetryPolicy
	}
	if s.DeadLetterConfig != nil {
		schedule.Target.DeadLetterConfig = s.DeadLetterConfig
	}
	if s.KmsKeyArn != nil {
		schedule.KmsKeyArn = s.KmsKeyArn
	}
	// If RoleArn is specified in properties, use it for the target
	if s.RoleArn != nil {
		schedule.Target.RoleArn = s.RoleArn
	}

	return schedule
}

// ToCloudFormation converts the Schedule to a CloudFormation resource.
func (s *Schedule) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	properties["ScheduleExpression"] = s.ScheduleExpression
	properties["FlexibleTimeWindow"] = s.FlexibleTimeWindow.toMap()
	properties["Target"] = s.Target.toMap()

	if s.ScheduleExpressionTimezone != "" {
		properties["ScheduleExpressionTimezone"] = s.ScheduleExpressionTimezone
	}
	if s.Name != nil {
		properties["Name"] = s.Name
	}
	if s.Description != "" {
		properties["Description"] = s.Description
	}
	if s.State != "" {
		properties["State"] = s.State
	}
	if s.GroupName != nil {
		properties["GroupName"] = s.GroupName
	}
	if s.StartDate != "" {
		properties["StartDate"] = s.StartDate
	}
	if s.EndDate != "" {
		properties["EndDate"] = s.EndDate
	}
	if s.KmsKeyArn != nil {
		properties["KmsKeyArn"] = s.KmsKeyArn
	}

	return map[string]interface{}{
		"Type":       ResourceTypeSchedule,
		"Properties": properties,
	}
}

func (f *FlexibleTimeWindow) toMap() map[string]interface{} {
	m := map[string]interface{}{
		"Mode": f.Mode,
	}
	if f.MaximumWindowInMinutes != nil {
		m["MaximumWindowInMinutes"] = *f.MaximumWindowInMinutes
	}
	return m
}

func (t *ScheduleTarget) toMap() map[string]interface{} {
	m := map[string]interface{}{
		"Arn":     t.Arn,
		"RoleArn": t.RoleArn,
	}
	if t.Input != nil {
		m["Input"] = t.Input
	}
	if t.RetryPolicy != nil {
		m["RetryPolicy"] = t.RetryPolicy.toMap()
	}
	if t.DeadLetterConfig != nil {
		m["DeadLetterConfig"] = t.DeadLetterConfig.toMap()
	}
	return m
}

func (r *ScheduleRetryPolicy) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if r.MaximumEventAgeInSeconds != nil {
		m["MaximumEventAgeInSeconds"] = *r.MaximumEventAgeInSeconds
	}
	if r.MaximumRetryAttempts != nil {
		m["MaximumRetryAttempts"] = *r.MaximumRetryAttempts
	}
	return m
}

func (d *ScheduleDeadLetterConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if d.Arn != nil {
		m["Arn"] = d.Arn
	}
	return m
}
