package pull

import (
	"testing"
)

func TestNewScheduleV2EventProperties(t *testing.T) {
	props := NewScheduleV2EventProperties("rate(1 hour)")

	if props.ScheduleExpression != "rate(1 hour)" {
		t.Errorf("expected ScheduleExpression 'rate(1 hour)', got %s", props.ScheduleExpression)
	}
}

func TestScheduleV2EventProperties_ToSchedule_Minimal(t *testing.T) {
	props := NewScheduleV2EventProperties("rate(5 minutes)")

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.ScheduleExpression != "rate(5 minutes)" {
		t.Errorf("expected ScheduleExpression, got %s", schedule.ScheduleExpression)
	}
	if schedule.Target.Arn != "arn:aws:lambda:us-east-1:123456789012:function:my-function" {
		t.Errorf("expected Target.Arn, got %v", schedule.Target.Arn)
	}
	if schedule.Target.RoleArn != "arn:aws:iam::123456789012:role/SchedulerRole" {
		t.Errorf("expected Target.RoleArn, got %v", schedule.Target.RoleArn)
	}
	// Default FlexibleTimeWindow should be OFF
	if schedule.FlexibleTimeWindow.Mode != "OFF" {
		t.Errorf("expected FlexibleTimeWindow.Mode 'OFF', got %s", schedule.FlexibleTimeWindow.Mode)
	}
}

func TestScheduleV2EventProperties_CronExpression(t *testing.T) {
	props := NewScheduleV2EventProperties("cron(0 12 * * ? *)")

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.ScheduleExpression != "cron(0 12 * * ? *)" {
		t.Errorf("expected cron expression, got %s", schedule.ScheduleExpression)
	}
}

func TestScheduleV2EventProperties_WithTimezone(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression:         "cron(0 9 * * ? *)",
		ScheduleExpressionTimezone: "America/New_York",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.ScheduleExpressionTimezone != "America/New_York" {
		t.Errorf("expected timezone 'America/New_York', got %s", schedule.ScheduleExpressionTimezone)
	}
}

func TestScheduleV2EventProperties_WithFlexibleTimeWindow(t *testing.T) {
	maxWindow := 15
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		FlexibleTimeWindow: &FlexibleTimeWindow{
			Mode:                   "FLEXIBLE",
			MaximumWindowInMinutes: &maxWindow,
		},
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.FlexibleTimeWindow.Mode != "FLEXIBLE" {
		t.Errorf("expected Mode 'FLEXIBLE', got %s", schedule.FlexibleTimeWindow.Mode)
	}
	if *schedule.FlexibleTimeWindow.MaximumWindowInMinutes != 15 {
		t.Errorf("expected MaximumWindowInMinutes 15, got %d", *schedule.FlexibleTimeWindow.MaximumWindowInMinutes)
	}
}

func TestScheduleV2EventProperties_WithName(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		Name:               "my-schedule",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.Name != "my-schedule" {
		t.Errorf("expected Name 'my-schedule', got %v", schedule.Name)
	}
}

func TestScheduleV2EventProperties_WithDescription(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		Description:        "My hourly schedule",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.Description != "My hourly schedule" {
		t.Errorf("expected Description, got %s", schedule.Description)
	}
}

func TestScheduleV2EventProperties_WithState(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		State:              "DISABLED",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.State != "DISABLED" {
		t.Errorf("expected State 'DISABLED', got %s", schedule.State)
	}
}

func TestScheduleV2EventProperties_WithGroupName(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		GroupName:          "my-schedule-group",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.GroupName != "my-schedule-group" {
		t.Errorf("expected GroupName, got %v", schedule.GroupName)
	}
}

func TestScheduleV2EventProperties_WithDateRange(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		StartDate:          "2024-01-01T00:00:00Z",
		EndDate:            "2024-12-31T23:59:59Z",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.StartDate != "2024-01-01T00:00:00Z" {
		t.Errorf("expected StartDate, got %s", schedule.StartDate)
	}
	if schedule.EndDate != "2024-12-31T23:59:59Z" {
		t.Errorf("expected EndDate, got %s", schedule.EndDate)
	}
}

func TestScheduleV2EventProperties_WithInput(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		Input:              `{"key": "value"}`,
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.Target.Input != `{"key": "value"}` {
		t.Errorf("expected Input, got %v", schedule.Target.Input)
	}
}

func TestScheduleV2EventProperties_WithRetryPolicy(t *testing.T) {
	maxAge := 3600
	maxRetries := 3
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		RetryPolicy: &ScheduleRetryPolicy{
			MaximumEventAgeInSeconds: &maxAge,
			MaximumRetryAttempts:     &maxRetries,
		},
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.Target.RetryPolicy == nil {
		t.Fatal("expected RetryPolicy to be set")
	}
	if *schedule.Target.RetryPolicy.MaximumEventAgeInSeconds != 3600 {
		t.Errorf("expected MaximumEventAgeInSeconds 3600, got %d", *schedule.Target.RetryPolicy.MaximumEventAgeInSeconds)
	}
	if *schedule.Target.RetryPolicy.MaximumRetryAttempts != 3 {
		t.Errorf("expected MaximumRetryAttempts 3, got %d", *schedule.Target.RetryPolicy.MaximumRetryAttempts)
	}
}

func TestScheduleV2EventProperties_WithDeadLetterConfig(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		DeadLetterConfig: &ScheduleDeadLetterConfig{
			Arn: "arn:aws:sqs:us-east-1:123456789012:dlq",
		},
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.Target.DeadLetterConfig == nil {
		t.Fatal("expected DeadLetterConfig to be set")
	}
	if schedule.Target.DeadLetterConfig.Arn != "arn:aws:sqs:us-east-1:123456789012:dlq" {
		t.Errorf("expected DLQ ARN, got %v", schedule.Target.DeadLetterConfig.Arn)
	}
}

func TestScheduleV2EventProperties_WithKmsKey(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		KmsKeyArn:          "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	if schedule.KmsKeyArn != "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012" {
		t.Errorf("expected KmsKeyArn, got %v", schedule.KmsKeyArn)
	}
}

func TestScheduleV2EventProperties_FullConfig(t *testing.T) {
	maxWindow := 30
	maxAge := 7200
	maxRetries := 5

	props := &ScheduleV2EventProperties{
		ScheduleExpression:         "cron(0 8 * * ? *)",
		ScheduleExpressionTimezone: "Europe/London",
		FlexibleTimeWindow: &FlexibleTimeWindow{
			Mode:                   "FLEXIBLE",
			MaximumWindowInMinutes: &maxWindow,
		},
		Name:        "daily-report-schedule",
		Description: "Runs daily report at 8 AM London time",
		State:       "ENABLED",
		GroupName:   "reports",
		StartDate:   "2024-01-01T00:00:00Z",
		EndDate:     "2024-12-31T23:59:59Z",
		Input:       `{"report": "daily"}`,
		RetryPolicy: &ScheduleRetryPolicy{
			MaximumEventAgeInSeconds: &maxAge,
			MaximumRetryAttempts:     &maxRetries,
		},
		DeadLetterConfig: &ScheduleDeadLetterConfig{
			Arn: "arn:aws:sqs:us-east-1:123456789012:dlq",
		},
		KmsKeyArn: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)

	// Verify all properties are set
	if schedule.ScheduleExpression != "cron(0 8 * * ? *)" {
		t.Errorf("unexpected ScheduleExpression: %s", schedule.ScheduleExpression)
	}
	if schedule.ScheduleExpressionTimezone != "Europe/London" {
		t.Errorf("unexpected ScheduleExpressionTimezone: %s", schedule.ScheduleExpressionTimezone)
	}
	if schedule.FlexibleTimeWindow.Mode != "FLEXIBLE" {
		t.Errorf("unexpected FlexibleTimeWindow.Mode: %s", schedule.FlexibleTimeWindow.Mode)
	}
	if schedule.Name != "daily-report-schedule" {
		t.Errorf("unexpected Name: %v", schedule.Name)
	}
	if schedule.Description != "Runs daily report at 8 AM London time" {
		t.Errorf("unexpected Description: %s", schedule.Description)
	}
	if schedule.State != "ENABLED" {
		t.Errorf("unexpected State: %s", schedule.State)
	}
	if schedule.GroupName != "reports" {
		t.Errorf("unexpected GroupName: %v", schedule.GroupName)
	}
}

func TestSchedule_ToCloudFormation_Minimal(t *testing.T) {
	props := NewScheduleV2EventProperties("rate(1 hour)")

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)
	cfn := schedule.ToCloudFormation()

	if cfn["Type"] != "AWS::Scheduler::Schedule" {
		t.Errorf("expected Type 'AWS::Scheduler::Schedule', got %v", cfn["Type"])
	}

	cfnProps := cfn["Properties"].(map[string]interface{})
	if cfnProps["ScheduleExpression"] != "rate(1 hour)" {
		t.Errorf("expected ScheduleExpression in CFN properties")
	}

	flexWindow := cfnProps["FlexibleTimeWindow"].(map[string]interface{})
	if flexWindow["Mode"] != "OFF" {
		t.Errorf("expected FlexibleTimeWindow.Mode 'OFF' in CFN properties")
	}

	target := cfnProps["Target"].(map[string]interface{})
	if target["Arn"] != "arn:aws:lambda:us-east-1:123456789012:function:my-function" {
		t.Errorf("expected Target.Arn in CFN properties")
	}
	if target["RoleArn"] != "arn:aws:iam::123456789012:role/SchedulerRole" {
		t.Errorf("expected Target.RoleArn in CFN properties")
	}
}

func TestSchedule_ToCloudFormation_FullConfig(t *testing.T) {
	maxWindow := 15
	maxAge := 3600
	maxRetries := 3

	props := &ScheduleV2EventProperties{
		ScheduleExpression:         "cron(0 12 * * ? *)",
		ScheduleExpressionTimezone: "UTC",
		FlexibleTimeWindow: &FlexibleTimeWindow{
			Mode:                   "FLEXIBLE",
			MaximumWindowInMinutes: &maxWindow,
		},
		Name:        "my-schedule",
		Description: "My schedule",
		State:       "ENABLED",
		GroupName:   "my-group",
		StartDate:   "2024-01-01T00:00:00Z",
		EndDate:     "2024-12-31T23:59:59Z",
		Input:       `{"key": "value"}`,
		RetryPolicy: &ScheduleRetryPolicy{
			MaximumEventAgeInSeconds: &maxAge,
			MaximumRetryAttempts:     &maxRetries,
		},
		DeadLetterConfig: &ScheduleDeadLetterConfig{
			Arn: "arn:aws:sqs:us-east-1:123456789012:dlq",
		},
		KmsKeyArn: "arn:aws:kms:us-east-1:123456789012:key/my-key",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/SchedulerRole",
	)
	cfn := schedule.ToCloudFormation()

	cfnProps := cfn["Properties"].(map[string]interface{})

	if cfnProps["ScheduleExpressionTimezone"] != "UTC" {
		t.Errorf("expected ScheduleExpressionTimezone in CFN properties")
	}
	if cfnProps["Name"] != "my-schedule" {
		t.Errorf("expected Name in CFN properties")
	}
	if cfnProps["Description"] != "My schedule" {
		t.Errorf("expected Description in CFN properties")
	}
	if cfnProps["State"] != "ENABLED" {
		t.Errorf("expected State in CFN properties")
	}
	if cfnProps["GroupName"] != "my-group" {
		t.Errorf("expected GroupName in CFN properties")
	}
	if cfnProps["StartDate"] != "2024-01-01T00:00:00Z" {
		t.Errorf("expected StartDate in CFN properties")
	}
	if cfnProps["EndDate"] != "2024-12-31T23:59:59Z" {
		t.Errorf("expected EndDate in CFN properties")
	}
	if cfnProps["KmsKeyArn"] != "arn:aws:kms:us-east-1:123456789012:key/my-key" {
		t.Errorf("expected KmsKeyArn in CFN properties")
	}

	flexWindow := cfnProps["FlexibleTimeWindow"].(map[string]interface{})
	if flexWindow["Mode"] != "FLEXIBLE" {
		t.Errorf("expected FlexibleTimeWindow.Mode 'FLEXIBLE' in CFN properties")
	}
	if flexWindow["MaximumWindowInMinutes"] != 15 {
		t.Errorf("expected FlexibleTimeWindow.MaximumWindowInMinutes in CFN properties")
	}

	target := cfnProps["Target"].(map[string]interface{})
	if target["Input"] != `{"key": "value"}` {
		t.Errorf("expected Target.Input in CFN properties")
	}

	retryPolicy := target["RetryPolicy"].(map[string]interface{})
	if retryPolicy["MaximumEventAgeInSeconds"] != 3600 {
		t.Errorf("expected RetryPolicy.MaximumEventAgeInSeconds in CFN properties")
	}
	if retryPolicy["MaximumRetryAttempts"] != 3 {
		t.Errorf("expected RetryPolicy.MaximumRetryAttempts in CFN properties")
	}

	dlq := target["DeadLetterConfig"].(map[string]interface{})
	if dlq["Arn"] != "arn:aws:sqs:us-east-1:123456789012:dlq" {
		t.Errorf("expected DeadLetterConfig.Arn in CFN properties")
	}
}

func TestScheduleV2EventProperties_WithRoleArnOverride(t *testing.T) {
	props := &ScheduleV2EventProperties{
		ScheduleExpression: "rate(1 hour)",
		RoleArn:            "arn:aws:iam::123456789012:role/CustomRole",
	}

	schedule := props.ToSchedule(
		"arn:aws:lambda:us-east-1:123456789012:function:my-function",
		"arn:aws:iam::123456789012:role/DefaultRole",
	)

	// RoleArn from properties should override the parameter
	if schedule.Target.RoleArn != "arn:aws:iam::123456789012:role/CustomRole" {
		t.Errorf("expected custom RoleArn, got %v", schedule.Target.RoleArn)
	}
}
