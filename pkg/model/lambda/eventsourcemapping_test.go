package lambda

import (
	"testing"
)

func TestNewEventSourceMapping(t *testing.T) {
	esm := NewEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
}

func TestNewKinesisEventSourceMapping(t *testing.T) {
	esm := NewKinesisEventSourceMapping(
		"my-function",
		"arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		"LATEST",
	)

	if esm.StartingPosition != "LATEST" {
		t.Errorf("expected StartingPosition 'LATEST', got %s", esm.StartingPosition)
	}
	if esm.EventSourceArn != "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream" {
		t.Errorf("unexpected EventSourceArn: %v", esm.EventSourceArn)
	}
}

func TestNewDynamoDBEventSourceMapping(t *testing.T) {
	esm := NewDynamoDBEventSourceMapping(
		"my-function",
		"arn:aws:dynamodb:us-east-1:123456789012:table/my-table/stream/2023-01-01T00:00:00.000",
		"TRIM_HORIZON",
	)

	if esm.StartingPosition != "TRIM_HORIZON" {
		t.Errorf("expected StartingPosition 'TRIM_HORIZON', got %s", esm.StartingPosition)
	}
}

func TestNewSQSEventSourceMapping(t *testing.T) {
	esm := NewSQSEventSourceMapping(
		"my-function",
		"arn:aws:sqs:us-east-1:123456789012:my-queue",
	)

	if esm.EventSourceArn != "arn:aws:sqs:us-east-1:123456789012:my-queue" {
		t.Errorf("unexpected EventSourceArn: %v", esm.EventSourceArn)
	}
	if esm.StartingPosition != "" {
		t.Errorf("expected no StartingPosition for SQS, got %s", esm.StartingPosition)
	}
}

func TestNewMSKEventSourceMapping(t *testing.T) {
	esm := NewMSKEventSourceMapping(
		"my-function",
		"arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		[]string{"my-topic"},
		"LATEST",
	)

	if len(esm.Topics) != 1 || esm.Topics[0] != "my-topic" {
		t.Errorf("unexpected Topics: %v", esm.Topics)
	}
}

func TestNewSelfManagedKafkaEventSourceMapping(t *testing.T) {
	esm := NewSelfManagedKafkaEventSourceMapping(
		"my-function",
		[]string{"broker1:9092", "broker2:9092"},
		[]string{"my-topic"},
		"LATEST",
	)

	if esm.SelfManagedEventSource == nil {
		t.Fatal("expected SelfManagedEventSource to be set")
	}
	if len(esm.SelfManagedEventSource.Endpoints.KafkaBootstrapServers) != 2 {
		t.Errorf("expected 2 bootstrap servers, got %d",
			len(esm.SelfManagedEventSource.Endpoints.KafkaBootstrapServers))
	}
}

func TestEventSourceMappingWithBatchSize(t *testing.T) {
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue").
		WithBatchSize(100)

	if esm.BatchSize == nil || *esm.BatchSize != 100 {
		t.Errorf("expected BatchSize 100, got %v", esm.BatchSize)
	}
}

func TestEventSourceMappingWithBatchingWindow(t *testing.T) {
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue").
		WithBatchingWindow(30)

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 30 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 30, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestEventSourceMappingWithEnabled(t *testing.T) {
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue").
		WithEnabled(false)

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestEventSourceMappingWithBisectOnError(t *testing.T) {
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "LATEST").
		WithBisectOnError(true)

	if !esm.BisectBatchOnFunctionError {
		t.Error("expected BisectBatchOnFunctionError to be true")
	}
}

func TestEventSourceMappingWithMaximumRetryAttempts(t *testing.T) {
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "LATEST").
		WithMaximumRetryAttempts(5)

	if esm.MaximumRetryAttempts == nil || *esm.MaximumRetryAttempts != 5 {
		t.Errorf("expected MaximumRetryAttempts 5, got %v", esm.MaximumRetryAttempts)
	}
}

func TestEventSourceMappingWithMaximumRecordAge(t *testing.T) {
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "LATEST").
		WithMaximumRecordAge(3600)

	if esm.MaximumRecordAgeInSeconds == nil || *esm.MaximumRecordAgeInSeconds != 3600 {
		t.Errorf("expected MaximumRecordAgeInSeconds 3600, got %v", esm.MaximumRecordAgeInSeconds)
	}
}

func TestEventSourceMappingWithParallelizationFactor(t *testing.T) {
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "LATEST").
		WithParallelizationFactor(5)

	if esm.ParallelizationFactor == nil || *esm.ParallelizationFactor != 5 {
		t.Errorf("expected ParallelizationFactor 5, got %v", esm.ParallelizationFactor)
	}
}

func TestEventSourceMappingWithTumblingWindow(t *testing.T) {
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "LATEST").
		WithTumblingWindow(60)

	if esm.TumblingWindowInSeconds == nil || *esm.TumblingWindowInSeconds != 60 {
		t.Errorf("expected TumblingWindowInSeconds 60, got %v", esm.TumblingWindowInSeconds)
	}
}

func TestEventSourceMappingWithOnFailureDestination(t *testing.T) {
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "LATEST").
		WithOnFailureDestination("arn:aws:sqs:us-east-1:123456789012:dlq")

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

func TestEventSourceMappingWithReportBatchItemFailures(t *testing.T) {
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "LATEST").
		WithReportBatchItemFailures()

	if len(esm.FunctionResponseTypes) != 1 || esm.FunctionResponseTypes[0] != "ReportBatchItemFailures" {
		t.Errorf("expected FunctionResponseTypes [ReportBatchItemFailures], got %v", esm.FunctionResponseTypes)
	}
}

func TestEventSourceMappingWithFilter(t *testing.T) {
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue").
		WithFilter(`{"body": {"status": ["PENDING"]}}`).
		WithFilter(`{"body": {"priority": ["HIGH"]}}`)

	if esm.FilterCriteria == nil {
		t.Fatal("expected FilterCriteria to be set")
	}
	if len(esm.FilterCriteria.Filters) != 2 {
		t.Errorf("expected 2 filters, got %d", len(esm.FilterCriteria.Filters))
	}
}

func TestEventSourceMappingWithScalingConfig(t *testing.T) {
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue").
		WithScalingConfig(10)

	if esm.ScalingConfig == nil {
		t.Fatal("expected ScalingConfig to be set")
	}
	if esm.ScalingConfig.MaximumConcurrency == nil || *esm.ScalingConfig.MaximumConcurrency != 10 {
		t.Errorf("expected MaximumConcurrency 10, got %v", esm.ScalingConfig.MaximumConcurrency)
	}
}

func TestEventSourceMappingAddSourceAccessConfiguration(t *testing.T) {
	esm := NewMSKEventSourceMapping(
		"my-function",
		"arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		[]string{"my-topic"},
		"LATEST",
	).AddSourceAccessConfiguration("SASL_SCRAM_512_AUTH", "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret")

	if len(esm.SourceAccessConfigurations) != 1 {
		t.Errorf("expected 1 SourceAccessConfiguration, got %d", len(esm.SourceAccessConfigurations))
	}
	if esm.SourceAccessConfigurations[0].Type != "SASL_SCRAM_512_AUTH" {
		t.Errorf("unexpected Type: %s", esm.SourceAccessConfigurations[0].Type)
	}
}

func TestEventSourceMappingToCloudFormation_Minimal(t *testing.T) {
	esm := NewEventSourceMapping("my-function")

	result := esm.ToCloudFormation()

	if result["Type"] != ResourceTypeEventSourceMapping {
		t.Errorf("expected Type %s, got %v", ResourceTypeEventSourceMapping, result["Type"])
	}

	props := result["Properties"].(map[string]interface{})
	if props["FunctionName"] != "my-function" {
		t.Errorf("expected FunctionName in properties")
	}
}

func TestEventSourceMappingToCloudFormation_SQS(t *testing.T) {
	batchSize := 10
	enabled := true
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue").
		WithBatchSize(batchSize).
		WithEnabled(enabled).
		WithFilter(`{"body": {"status": ["PENDING"]}}`)

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["BatchSize"] != 10 {
		t.Errorf("expected BatchSize 10, got %v", props["BatchSize"])
	}
	if props["Enabled"] != true {
		t.Errorf("expected Enabled true, got %v", props["Enabled"])
	}
	if props["EventSourceArn"] != "arn:aws:sqs:us-east-1:123456789012:my-queue" {
		t.Errorf("expected EventSourceArn in properties")
	}

	filterCriteria := props["FilterCriteria"].(map[string]interface{})
	filters := filterCriteria["Filters"].([]map[string]interface{})
	if len(filters) != 1 {
		t.Errorf("expected 1 filter, got %d", len(filters))
	}
}

func TestEventSourceMappingToCloudFormation_Kinesis(t *testing.T) {
	esm := NewKinesisEventSourceMapping(
		"my-function",
		"arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
		"LATEST",
	).
		WithBisectOnError(true).
		WithMaximumRetryAttempts(3).
		WithMaximumRecordAge(3600).
		WithParallelizationFactor(5).
		WithTumblingWindow(60).
		WithOnFailureDestination("arn:aws:sqs:us-east-1:123456789012:dlq").
		WithReportBatchItemFailures()

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["StartingPosition"] != "LATEST" {
		t.Errorf("expected StartingPosition 'LATEST', got %v", props["StartingPosition"])
	}
	if props["BisectBatchOnFunctionError"] != true {
		t.Errorf("expected BisectBatchOnFunctionError true, got %v", props["BisectBatchOnFunctionError"])
	}
	if props["MaximumRetryAttempts"] != 3 {
		t.Errorf("expected MaximumRetryAttempts 3, got %v", props["MaximumRetryAttempts"])
	}
	if props["MaximumRecordAgeInSeconds"] != 3600 {
		t.Errorf("expected MaximumRecordAgeInSeconds 3600, got %v", props["MaximumRecordAgeInSeconds"])
	}
	if props["ParallelizationFactor"] != 5 {
		t.Errorf("expected ParallelizationFactor 5, got %v", props["ParallelizationFactor"])
	}
	if props["TumblingWindowInSeconds"] != 60 {
		t.Errorf("expected TumblingWindowInSeconds 60, got %v", props["TumblingWindowInSeconds"])
	}

	destConfig := props["DestinationConfig"].(map[string]interface{})
	onFailure := destConfig["OnFailure"].(map[string]interface{})
	if onFailure["Destination"] != "arn:aws:sqs:us-east-1:123456789012:dlq" {
		t.Errorf("unexpected Destination: %v", onFailure["Destination"])
	}

	responseTypes := props["FunctionResponseTypes"].([]string)
	if len(responseTypes) != 1 || responseTypes[0] != "ReportBatchItemFailures" {
		t.Errorf("unexpected FunctionResponseTypes: %v", responseTypes)
	}
}

func TestEventSourceMappingToCloudFormation_SelfManagedKafka(t *testing.T) {
	esm := NewSelfManagedKafkaEventSourceMapping(
		"my-function",
		[]string{"broker1:9092", "broker2:9092"},
		[]string{"my-topic"},
		"LATEST",
	).AddSourceAccessConfiguration("SASL_SCRAM_512_AUTH", "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret")

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	selfManaged := props["SelfManagedEventSource"].(map[string]interface{})
	endpoints := selfManaged["Endpoints"].(map[string]interface{})
	servers := endpoints["KafkaBootstrapServers"].([]string)
	if len(servers) != 2 {
		t.Errorf("expected 2 bootstrap servers, got %d", len(servers))
	}

	topics := props["Topics"].([]string)
	if len(topics) != 1 || topics[0] != "my-topic" {
		t.Errorf("unexpected Topics: %v", topics)
	}

	configs := props["SourceAccessConfigurations"].([]map[string]interface{})
	if len(configs) != 1 {
		t.Errorf("expected 1 SourceAccessConfiguration, got %d", len(configs))
	}
}

func TestEventSourceMappingWithQueues(t *testing.T) {
	esm := NewEventSourceMapping("my-function")
	esm.Queues = []string{"my-queue"}
	esm.EventSourceArn = "arn:aws:mq:us-east-1:123456789012:broker:my-broker:12345678-1234-1234-1234-123456789012"

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	queues := props["Queues"].([]string)
	if len(queues) != 1 || queues[0] != "my-queue" {
		t.Errorf("unexpected Queues: %v", queues)
	}
}

func TestEventSourceMappingWithAmazonManagedKafkaConfig(t *testing.T) {
	esm := NewEventSourceMapping("my-function")
	esm.AmazonManagedKafkaEventSourceConfig = &AmazonManagedKafkaEventSourceConfig{
		ConsumerGroupId: "my-consumer-group",
	}

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	config := props["AmazonManagedKafkaEventSourceConfig"].(map[string]interface{})
	if config["ConsumerGroupId"] != "my-consumer-group" {
		t.Errorf("unexpected ConsumerGroupId: %v", config["ConsumerGroupId"])
	}
}

func TestEventSourceMappingWithDocumentDBConfig(t *testing.T) {
	esm := NewEventSourceMapping("my-function")
	esm.DocumentDBEventSourceConfig = &DocumentDBEventSourceConfig{
		CollectionName: "my-collection",
		DatabaseName:   "my-database",
		FullDocument:   "UpdateLookup",
	}

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	config := props["DocumentDBEventSourceConfig"].(map[string]interface{})
	if config["CollectionName"] != "my-collection" {
		t.Errorf("unexpected CollectionName: %v", config["CollectionName"])
	}
	if config["DatabaseName"] != "my-database" {
		t.Errorf("unexpected DatabaseName: %v", config["DatabaseName"])
	}
	if config["FullDocument"] != "UpdateLookup" {
		t.Errorf("unexpected FullDocument: %v", config["FullDocument"])
	}
}

func TestEventSourceMappingWithTags(t *testing.T) {
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue")
	esm.Tags = []Tag{
		{Key: "Environment", Value: "production"},
	}

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	tags := props["Tags"].([]map[string]interface{})
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
	if tags[0]["Key"] != "Environment" || tags[0]["Value"] != "production" {
		t.Errorf("unexpected tag: %v", tags[0])
	}
}

func TestEventSourceMappingWithKmsKeyArn(t *testing.T) {
	esm := NewSQSEventSourceMapping("my-function", "arn:aws:sqs:us-east-1:123456789012:my-queue")
	esm.KmsKeyArn = "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["KmsKeyArn"] != "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012" {
		t.Errorf("unexpected KmsKeyArn: %v", props["KmsKeyArn"])
	}
}

func TestEventSourceMappingWithStartingPositionTimestamp(t *testing.T) {
	timestamp := 1609459200.0
	esm := NewKinesisEventSourceMapping("my-function", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", "AT_TIMESTAMP")
	esm.StartingPositionTimestamp = &timestamp

	result := esm.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	if props["StartingPositionTimestamp"] != 1609459200.0 {
		t.Errorf("unexpected StartingPositionTimestamp: %v", props["StartingPositionTimestamp"])
	}
}
