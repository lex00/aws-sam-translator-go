package pull

import (
	"testing"
)

func TestNewKinesisEventProperties(t *testing.T) {
	props := NewKinesisEventProperties(
		"arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		"LATEST",
	)

	if props.Stream != "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream" {
		t.Errorf("expected Stream ARN, got %v", props.Stream)
	}
	if props.StartingPosition != "LATEST" {
		t.Errorf("expected StartingPosition 'LATEST', got %s", props.StartingPosition)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_Minimal(t *testing.T) {
	props := NewKinesisEventProperties(
		"arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		"TRIM_HORIZON",
	)

	esm := props.ToEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream" {
		t.Errorf("expected EventSourceArn, got %v", esm.EventSourceArn)
	}
	if esm.StartingPosition != "TRIM_HORIZON" {
		t.Errorf("expected StartingPosition 'TRIM_HORIZON', got %s", esm.StartingPosition)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithBatchSize(t *testing.T) {
	batchSize := 500
	props := &KinesisEventProperties{
		Stream:           "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition: "LATEST",
		BatchSize:        &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.BatchSize == nil || *esm.BatchSize != 500 {
		t.Errorf("expected BatchSize 500, got %v", esm.BatchSize)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithBatchingWindow(t *testing.T) {
	window := 30
	props := &KinesisEventProperties{
		Stream:                         "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:               "LATEST",
		MaximumBatchingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 30 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 30, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithEnabled(t *testing.T) {
	enabled := false
	props := &KinesisEventProperties{
		Stream:           "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition: "LATEST",
		Enabled:          &enabled,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithBisectBatchOnError(t *testing.T) {
	bisect := true
	props := &KinesisEventProperties{
		Stream:                     "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:           "LATEST",
		BisectBatchOnFunctionError: &bisect,
	}

	esm := props.ToEventSourceMapping("my-function")

	if !esm.BisectBatchOnFunctionError {
		t.Error("expected BisectBatchOnFunctionError to be true")
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithRetryAttempts(t *testing.T) {
	retries := 3
	props := &KinesisEventProperties{
		Stream:               "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:     "LATEST",
		MaximumRetryAttempts: &retries,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumRetryAttempts == nil || *esm.MaximumRetryAttempts != 3 {
		t.Errorf("expected MaximumRetryAttempts 3, got %v", esm.MaximumRetryAttempts)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithRecordAge(t *testing.T) {
	age := 3600
	props := &KinesisEventProperties{
		Stream:                    "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:          "LATEST",
		MaximumRecordAgeInSeconds: &age,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumRecordAgeInSeconds == nil || *esm.MaximumRecordAgeInSeconds != 3600 {
		t.Errorf("expected MaximumRecordAgeInSeconds 3600, got %v", esm.MaximumRecordAgeInSeconds)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithParallelizationFactor(t *testing.T) {
	factor := 5
	props := &KinesisEventProperties{
		Stream:                "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:      "LATEST",
		ParallelizationFactor: &factor,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.ParallelizationFactor == nil || *esm.ParallelizationFactor != 5 {
		t.Errorf("expected ParallelizationFactor 5, got %v", esm.ParallelizationFactor)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithTumblingWindow(t *testing.T) {
	window := 60
	props := &KinesisEventProperties{
		Stream:                  "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:        "LATEST",
		TumblingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.TumblingWindowInSeconds == nil || *esm.TumblingWindowInSeconds != 60 {
		t.Errorf("expected TumblingWindowInSeconds 60, got %v", esm.TumblingWindowInSeconds)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithDestinationConfig(t *testing.T) {
	props := &KinesisEventProperties{
		Stream:           "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition: "LATEST",
		DestinationConfig: &DestinationConfig{
			OnFailure: &OnFailure{
				Destination: "arn:aws:sqs:us-east-1:123456789012:dlq",
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.DestinationConfig == nil {
		t.Fatal("expected DestinationConfig to be set")
	}
	if esm.DestinationConfig.OnFailure == nil {
		t.Fatal("expected OnFailure to be set")
	}
	if esm.DestinationConfig.OnFailure.Destination != "arn:aws:sqs:us-east-1:123456789012:dlq" {
		t.Errorf("unexpected Destination: %v", esm.DestinationConfig.OnFailure.Destination)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithFunctionResponseTypes(t *testing.T) {
	props := &KinesisEventProperties{
		Stream:                "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:      "LATEST",
		FunctionResponseTypes: []string{"ReportBatchItemFailures"},
	}

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.FunctionResponseTypes) != 1 || esm.FunctionResponseTypes[0] != "ReportBatchItemFailures" {
		t.Errorf("expected FunctionResponseTypes [ReportBatchItemFailures], got %v", esm.FunctionResponseTypes)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_WithFilterCriteria(t *testing.T) {
	props := &KinesisEventProperties{
		Stream:           "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition: "LATEST",
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"data": {"status": ["PENDING"]}}`},
				{Pattern: `{"data": {"priority": ["HIGH"]}}`},
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

func TestKinesisEventProperties_ToEventSourceMapping_WithStartingPositionTimestamp(t *testing.T) {
	timestamp := 1609459200.0
	props := &KinesisEventProperties{
		Stream:                    "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:          "AT_TIMESTAMP",
		StartingPositionTimestamp: &timestamp,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.StartingPositionTimestamp == nil || *esm.StartingPositionTimestamp != 1609459200.0 {
		t.Errorf("expected StartingPositionTimestamp 1609459200.0, got %v", esm.StartingPositionTimestamp)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_FullConfig(t *testing.T) {
	batchSize := 100
	window := 5
	enabled := true
	bisect := true
	retries := 3
	age := 86400
	parallel := 2
	tumbling := 30
	timestamp := 1609459200.0

	props := &KinesisEventProperties{
		Stream:                         "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition:               "AT_TIMESTAMP",
		StartingPositionTimestamp:      &timestamp,
		BatchSize:                      &batchSize,
		MaximumBatchingWindowInSeconds: &window,
		Enabled:                        &enabled,
		BisectBatchOnFunctionError:     &bisect,
		MaximumRetryAttempts:           &retries,
		MaximumRecordAgeInSeconds:      &age,
		ParallelizationFactor:          &parallel,
		TumblingWindowInSeconds:        &tumbling,
		DestinationConfig: &DestinationConfig{
			OnFailure: &OnFailure{
				Destination: "arn:aws:sqs:us-east-1:123456789012:dlq",
			},
		},
		FunctionResponseTypes: []string{"ReportBatchItemFailures"},
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"data": {"status": ["PENDING"]}}`},
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	// Verify all properties are set
	if esm.FunctionName != "my-function" {
		t.Errorf("unexpected FunctionName: %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream" {
		t.Errorf("unexpected EventSourceArn: %v", esm.EventSourceArn)
	}
	if esm.StartingPosition != "AT_TIMESTAMP" {
		t.Errorf("unexpected StartingPosition: %s", esm.StartingPosition)
	}
	if *esm.BatchSize != 100 {
		t.Errorf("unexpected BatchSize: %d", *esm.BatchSize)
	}
	if *esm.MaximumBatchingWindowInSeconds != 5 {
		t.Errorf("unexpected MaximumBatchingWindowInSeconds: %d", *esm.MaximumBatchingWindowInSeconds)
	}
	if !*esm.Enabled {
		t.Error("unexpected Enabled value")
	}
	if !esm.BisectBatchOnFunctionError {
		t.Error("unexpected BisectBatchOnFunctionError value")
	}
	if *esm.MaximumRetryAttempts != 3 {
		t.Errorf("unexpected MaximumRetryAttempts: %d", *esm.MaximumRetryAttempts)
	}
	if *esm.MaximumRecordAgeInSeconds != 86400 {
		t.Errorf("unexpected MaximumRecordAgeInSeconds: %d", *esm.MaximumRecordAgeInSeconds)
	}
	if *esm.ParallelizationFactor != 2 {
		t.Errorf("unexpected ParallelizationFactor: %d", *esm.ParallelizationFactor)
	}
	if *esm.TumblingWindowInSeconds != 30 {
		t.Errorf("unexpected TumblingWindowInSeconds: %d", *esm.TumblingWindowInSeconds)
	}
}

func TestKinesisEventProperties_ToEventSourceMapping_ToCloudFormation(t *testing.T) {
	batchSize := 100
	props := &KinesisEventProperties{
		Stream:           "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		StartingPosition: "LATEST",
		BatchSize:        &batchSize,
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
	if cfnProps["EventSourceArn"] != "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream" {
		t.Errorf("expected EventSourceArn in CFN properties")
	}
	if cfnProps["StartingPosition"] != "LATEST" {
		t.Errorf("expected StartingPosition in CFN properties")
	}
	if cfnProps["BatchSize"] != 100 {
		t.Errorf("expected BatchSize in CFN properties")
	}
}
