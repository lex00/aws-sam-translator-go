package pull

import (
	"testing"
)

func TestNewDynamoDBEventProperties(t *testing.T) {
	props := NewDynamoDBEventProperties(
		"arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		"TRIM_HORIZON",
	)

	if props.Stream != "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000" {
		t.Errorf("expected Stream ARN, got %v", props.Stream)
	}
	if props.StartingPosition != "TRIM_HORIZON" {
		t.Errorf("expected StartingPosition 'TRIM_HORIZON', got %s", props.StartingPosition)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_Minimal(t *testing.T) {
	props := NewDynamoDBEventProperties(
		"arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		"LATEST",
	)

	esm := props.ToEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000" {
		t.Errorf("expected EventSourceArn, got %v", esm.EventSourceArn)
	}
	if esm.StartingPosition != "LATEST" {
		t.Errorf("expected StartingPosition 'LATEST', got %s", esm.StartingPosition)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithBatchSize(t *testing.T) {
	batchSize := 200
	props := &DynamoDBEventProperties{
		Stream:           "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition: "TRIM_HORIZON",
		BatchSize:        &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.BatchSize == nil || *esm.BatchSize != 200 {
		t.Errorf("expected BatchSize 200, got %v", esm.BatchSize)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithBatchingWindow(t *testing.T) {
	window := 60
	props := &DynamoDBEventProperties{
		Stream:                         "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:               "TRIM_HORIZON",
		MaximumBatchingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 60 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 60, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithEnabled(t *testing.T) {
	enabled := false
	props := &DynamoDBEventProperties{
		Stream:           "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition: "TRIM_HORIZON",
		Enabled:          &enabled,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithBisectBatchOnError(t *testing.T) {
	bisect := true
	props := &DynamoDBEventProperties{
		Stream:                     "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:           "TRIM_HORIZON",
		BisectBatchOnFunctionError: &bisect,
	}

	esm := props.ToEventSourceMapping("my-function")

	if !esm.BisectBatchOnFunctionError {
		t.Error("expected BisectBatchOnFunctionError to be true")
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithRetryAttempts(t *testing.T) {
	retries := 5
	props := &DynamoDBEventProperties{
		Stream:               "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:     "TRIM_HORIZON",
		MaximumRetryAttempts: &retries,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumRetryAttempts == nil || *esm.MaximumRetryAttempts != 5 {
		t.Errorf("expected MaximumRetryAttempts 5, got %v", esm.MaximumRetryAttempts)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithRecordAge(t *testing.T) {
	age := 7200
	props := &DynamoDBEventProperties{
		Stream:                    "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:          "TRIM_HORIZON",
		MaximumRecordAgeInSeconds: &age,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumRecordAgeInSeconds == nil || *esm.MaximumRecordAgeInSeconds != 7200 {
		t.Errorf("expected MaximumRecordAgeInSeconds 7200, got %v", esm.MaximumRecordAgeInSeconds)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithParallelizationFactor(t *testing.T) {
	factor := 10
	props := &DynamoDBEventProperties{
		Stream:                "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:      "TRIM_HORIZON",
		ParallelizationFactor: &factor,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.ParallelizationFactor == nil || *esm.ParallelizationFactor != 10 {
		t.Errorf("expected ParallelizationFactor 10, got %v", esm.ParallelizationFactor)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithTumblingWindow(t *testing.T) {
	window := 120
	props := &DynamoDBEventProperties{
		Stream:                  "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:        "TRIM_HORIZON",
		TumblingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.TumblingWindowInSeconds == nil || *esm.TumblingWindowInSeconds != 120 {
		t.Errorf("expected TumblingWindowInSeconds 120, got %v", esm.TumblingWindowInSeconds)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithDestinationConfig(t *testing.T) {
	props := &DynamoDBEventProperties{
		Stream:           "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition: "TRIM_HORIZON",
		DestinationConfig: &DestinationConfig{
			OnFailure: &OnFailure{
				Destination: "arn:aws:sns:us-east-1:123456789012:my-topic",
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
	if esm.DestinationConfig.OnFailure.Destination != "arn:aws:sns:us-east-1:123456789012:my-topic" {
		t.Errorf("unexpected Destination: %v", esm.DestinationConfig.OnFailure.Destination)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithFunctionResponseTypes(t *testing.T) {
	props := &DynamoDBEventProperties{
		Stream:                "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:      "TRIM_HORIZON",
		FunctionResponseTypes: []string{"ReportBatchItemFailures"},
	}

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.FunctionResponseTypes) != 1 || esm.FunctionResponseTypes[0] != "ReportBatchItemFailures" {
		t.Errorf("expected FunctionResponseTypes [ReportBatchItemFailures], got %v", esm.FunctionResponseTypes)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_WithFilterCriteria(t *testing.T) {
	props := &DynamoDBEventProperties{
		Stream:           "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition: "TRIM_HORIZON",
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"eventName": ["INSERT"]}`},
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.FilterCriteria == nil {
		t.Fatal("expected FilterCriteria to be set")
	}
	if len(esm.FilterCriteria.Filters) != 1 {
		t.Errorf("expected 1 filter, got %d", len(esm.FilterCriteria.Filters))
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_FullConfig(t *testing.T) {
	batchSize := 100
	window := 10
	enabled := true
	bisect := true
	retries := 2
	age := 3600
	parallel := 5
	tumbling := 60

	props := &DynamoDBEventProperties{
		Stream:                         "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition:               "TRIM_HORIZON",
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
				{Pattern: `{"eventName": ["INSERT", "MODIFY"]}`},
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	// Verify all properties are set
	if esm.FunctionName != "my-function" {
		t.Errorf("unexpected FunctionName: %v", esm.FunctionName)
	}
	if *esm.BatchSize != 100 {
		t.Errorf("unexpected BatchSize: %d", *esm.BatchSize)
	}
	if *esm.MaximumBatchingWindowInSeconds != 10 {
		t.Errorf("unexpected MaximumBatchingWindowInSeconds: %d", *esm.MaximumBatchingWindowInSeconds)
	}
	if !*esm.Enabled {
		t.Error("unexpected Enabled value")
	}
	if !esm.BisectBatchOnFunctionError {
		t.Error("unexpected BisectBatchOnFunctionError value")
	}
	if *esm.MaximumRetryAttempts != 2 {
		t.Errorf("unexpected MaximumRetryAttempts: %d", *esm.MaximumRetryAttempts)
	}
	if *esm.MaximumRecordAgeInSeconds != 3600 {
		t.Errorf("unexpected MaximumRecordAgeInSeconds: %d", *esm.MaximumRecordAgeInSeconds)
	}
	if *esm.ParallelizationFactor != 5 {
		t.Errorf("unexpected ParallelizationFactor: %d", *esm.ParallelizationFactor)
	}
	if *esm.TumblingWindowInSeconds != 60 {
		t.Errorf("unexpected TumblingWindowInSeconds: %d", *esm.TumblingWindowInSeconds)
	}
}

func TestDynamoDBEventProperties_ToEventSourceMapping_ToCloudFormation(t *testing.T) {
	batchSize := 50
	props := &DynamoDBEventProperties{
		Stream:           "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		StartingPosition: "TRIM_HORIZON",
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
	if cfnProps["EventSourceArn"] != "arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000" {
		t.Errorf("expected EventSourceArn in CFN properties")
	}
	if cfnProps["StartingPosition"] != "TRIM_HORIZON" {
		t.Errorf("expected StartingPosition in CFN properties")
	}
	if cfnProps["BatchSize"] != 50 {
		t.Errorf("expected BatchSize in CFN properties")
	}
}
