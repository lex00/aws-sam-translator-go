package pull

import (
	"testing"
)

func TestNewSQSEventProperties(t *testing.T) {
	props := NewSQSEventProperties("arn:aws:sqs:us-east-1:123456789012:my-queue")

	if props.Queue != "arn:aws:sqs:us-east-1:123456789012:my-queue" {
		t.Errorf("expected Queue ARN, got %v", props.Queue)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_Minimal(t *testing.T) {
	props := NewSQSEventProperties("arn:aws:sqs:us-east-1:123456789012:my-queue")

	esm := props.ToEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:sqs:us-east-1:123456789012:my-queue" {
		t.Errorf("expected EventSourceArn, got %v", esm.EventSourceArn)
	}
	// SQS does not use StartingPosition
	if esm.StartingPosition != "" {
		t.Errorf("expected empty StartingPosition for SQS, got %s", esm.StartingPosition)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_WithBatchSize(t *testing.T) {
	batchSize := 100
	props := &SQSEventProperties{
		Queue:     "arn:aws:sqs:us-east-1:123456789012:my-queue",
		BatchSize: &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.BatchSize == nil || *esm.BatchSize != 100 {
		t.Errorf("expected BatchSize 100, got %v", esm.BatchSize)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_WithBatchingWindow(t *testing.T) {
	window := 30
	props := &SQSEventProperties{
		Queue:                          "arn:aws:sqs:us-east-1:123456789012:my-queue",
		MaximumBatchingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 30 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 30, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_WithEnabled(t *testing.T) {
	enabled := false
	props := &SQSEventProperties{
		Queue:   "arn:aws:sqs:us-east-1:123456789012:my-queue",
		Enabled: &enabled,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_WithFunctionResponseTypes(t *testing.T) {
	props := &SQSEventProperties{
		Queue:                 "arn:aws:sqs:us-east-1:123456789012:my-queue",
		FunctionResponseTypes: []string{"ReportBatchItemFailures"},
	}

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.FunctionResponseTypes) != 1 || esm.FunctionResponseTypes[0] != "ReportBatchItemFailures" {
		t.Errorf("expected FunctionResponseTypes [ReportBatchItemFailures], got %v", esm.FunctionResponseTypes)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_WithFilterCriteria(t *testing.T) {
	props := &SQSEventProperties{
		Queue: "arn:aws:sqs:us-east-1:123456789012:my-queue",
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"body": {"status": ["PENDING"]}}`},
				{Pattern: `{"body": {"priority": ["HIGH"]}}`},
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.FilterCriteria == nil {
		t.Fatal("expected FilterCriteria to be set")
	}
	if len(esm.FilterCriteria.Filters) != 2 {
		t.Errorf("expected 2 filters, got %d", len(esm.FilterCriteria.Filters))
	}
}

func TestSQSEventProperties_ToEventSourceMapping_WithScalingConfig(t *testing.T) {
	maxConcurrency := 10
	props := &SQSEventProperties{
		Queue: "arn:aws:sqs:us-east-1:123456789012:my-queue",
		ScalingConfig: &ScalingConfig{
			MaximumConcurrency: &maxConcurrency,
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.ScalingConfig == nil {
		t.Fatal("expected ScalingConfig to be set")
	}
	if esm.ScalingConfig.MaximumConcurrency == nil || *esm.ScalingConfig.MaximumConcurrency != 10 {
		t.Errorf("expected MaximumConcurrency 10, got %v", esm.ScalingConfig.MaximumConcurrency)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_FullConfig(t *testing.T) {
	batchSize := 50
	window := 15
	enabled := true
	maxConcurrency := 20

	props := &SQSEventProperties{
		Queue:                          "arn:aws:sqs:us-east-1:123456789012:my-queue",
		BatchSize:                      &batchSize,
		MaximumBatchingWindowInSeconds: &window,
		Enabled:                        &enabled,
		FunctionResponseTypes:          []string{"ReportBatchItemFailures"},
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"body": {"type": ["ORDER"]}}`},
			},
		},
		ScalingConfig: &ScalingConfig{
			MaximumConcurrency: &maxConcurrency,
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	// Verify all properties are set
	if esm.FunctionName != "my-function" {
		t.Errorf("unexpected FunctionName: %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:sqs:us-east-1:123456789012:my-queue" {
		t.Errorf("unexpected EventSourceArn: %v", esm.EventSourceArn)
	}
	if *esm.BatchSize != 50 {
		t.Errorf("unexpected BatchSize: %d", *esm.BatchSize)
	}
	if *esm.MaximumBatchingWindowInSeconds != 15 {
		t.Errorf("unexpected MaximumBatchingWindowInSeconds: %d", *esm.MaximumBatchingWindowInSeconds)
	}
	if !*esm.Enabled {
		t.Error("unexpected Enabled value")
	}
	if len(esm.FunctionResponseTypes) != 1 {
		t.Errorf("unexpected FunctionResponseTypes length: %d", len(esm.FunctionResponseTypes))
	}
	if len(esm.FilterCriteria.Filters) != 1 {
		t.Errorf("unexpected FilterCriteria.Filters length: %d", len(esm.FilterCriteria.Filters))
	}
	if *esm.ScalingConfig.MaximumConcurrency != 20 {
		t.Errorf("unexpected MaximumConcurrency: %d", *esm.ScalingConfig.MaximumConcurrency)
	}
}

func TestSQSEventProperties_ToEventSourceMapping_ToCloudFormation(t *testing.T) {
	batchSize := 10
	props := &SQSEventProperties{
		Queue:     "arn:aws:sqs:us-east-1:123456789012:my-queue",
		BatchSize: &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")
	cfn := esm.ToCloudFormation()

	if cfn["Type"] != "AWS::Lambda::EventSourceMapping" {
		t.Errorf("expected Type 'AWS::Lambda::EventSourceMapping', got %v", cfn["Type"])
	}

	cfnProps := cfn["Properties"].(map[string]interface{})
	if cfnProps["FunctionName"] != "my-function" {
		t.Errorf("expected FunctionName in CFN properties")
	}
	if cfnProps["EventSourceArn"] != "arn:aws:sqs:us-east-1:123456789012:my-queue" {
		t.Errorf("expected EventSourceArn in CFN properties")
	}
	if cfnProps["BatchSize"] != 10 {
		t.Errorf("expected BatchSize in CFN properties")
	}
	// Verify StartingPosition is not set for SQS
	if _, ok := cfnProps["StartingPosition"]; ok {
		t.Error("StartingPosition should not be set for SQS")
	}
}

func TestSQSEventProperties_ToEventSourceMapping_FIFOQueue(t *testing.T) {
	batchSize := 10
	props := &SQSEventProperties{
		Queue:     "arn:aws:sqs:us-east-1:123456789012:my-queue.fifo",
		BatchSize: &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.EventSourceArn != "arn:aws:sqs:us-east-1:123456789012:my-queue.fifo" {
		t.Errorf("expected FIFO queue ARN, got %v", esm.EventSourceArn)
	}
}
