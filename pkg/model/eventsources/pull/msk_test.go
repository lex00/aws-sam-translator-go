package pull

import (
	"testing"
)

func TestNewMSKEventProperties(t *testing.T) {
	props := NewMSKEventProperties(
		"arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		[]string{"my-topic"},
		"LATEST",
	)

	if props.Stream != "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1" {
		t.Errorf("expected Stream ARN, got %v", props.Stream)
	}
	if len(props.Topics) != 1 || props.Topics[0] != "my-topic" {
		t.Errorf("expected Topics [my-topic], got %v", props.Topics)
	}
	if props.StartingPosition != "LATEST" {
		t.Errorf("expected StartingPosition 'LATEST', got %s", props.StartingPosition)
	}
}

func TestMSKEventProperties_ToEventSourceMapping_Minimal(t *testing.T) {
	props := NewMSKEventProperties(
		"arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		[]string{"my-topic"},
		"TRIM_HORIZON",
	)

	esm := props.ToEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1" {
		t.Errorf("expected EventSourceArn, got %v", esm.EventSourceArn)
	}
	if esm.StartingPosition != "TRIM_HORIZON" {
		t.Errorf("expected StartingPosition 'TRIM_HORIZON', got %s", esm.StartingPosition)
	}
	if len(esm.Topics) != 1 || esm.Topics[0] != "my-topic" {
		t.Errorf("expected Topics [my-topic], got %v", esm.Topics)
	}
}

func TestMSKEventProperties_ToEventSourceMapping_MultipleTopics(t *testing.T) {
	props := NewMSKEventProperties(
		"arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		[]string{"topic-1", "topic-2", "topic-3"},
		"LATEST",
	)

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.Topics) != 3 {
		t.Errorf("expected 3 topics, got %d", len(esm.Topics))
	}
	expectedTopics := []string{"topic-1", "topic-2", "topic-3"}
	for i, topic := range esm.Topics {
		if topic != expectedTopics[i] {
			t.Errorf("expected topic %s at index %d, got %s", expectedTopics[i], i, topic)
		}
	}
}

func TestMSKEventProperties_ToEventSourceMapping_WithBatchSize(t *testing.T) {
	batchSize := 500
	props := &MSKEventProperties{
		Stream:           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:           []string{"my-topic"},
		StartingPosition: "LATEST",
		BatchSize:        &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.BatchSize == nil || *esm.BatchSize != 500 {
		t.Errorf("expected BatchSize 500, got %v", esm.BatchSize)
	}
}

func TestMSKEventProperties_ToEventSourceMapping_WithBatchingWindow(t *testing.T) {
	window := 60
	props := &MSKEventProperties{
		Stream:                         "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:                         []string{"my-topic"},
		StartingPosition:               "LATEST",
		MaximumBatchingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 60 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 60, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestMSKEventProperties_ToEventSourceMapping_WithEnabled(t *testing.T) {
	enabled := false
	props := &MSKEventProperties{
		Stream:           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:           []string{"my-topic"},
		StartingPosition: "LATEST",
		Enabled:          &enabled,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestMSKEventProperties_ToEventSourceMapping_WithConsumerGroupId(t *testing.T) {
	props := &MSKEventProperties{
		Stream:           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:           []string{"my-topic"},
		StartingPosition: "LATEST",
		ConsumerGroupId:  "my-consumer-group",
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.AmazonManagedKafkaEventSourceConfig == nil {
		t.Fatal("expected AmazonManagedKafkaEventSourceConfig to be set")
	}
	if esm.AmazonManagedKafkaEventSourceConfig.ConsumerGroupId != "my-consumer-group" {
		t.Errorf("expected ConsumerGroupId 'my-consumer-group', got %s",
			esm.AmazonManagedKafkaEventSourceConfig.ConsumerGroupId)
	}
}

func TestMSKEventProperties_ToEventSourceMapping_WithSourceAccessConfigurations(t *testing.T) {
	props := &MSKEventProperties{
		Stream:           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:           []string{"my-topic"},
		StartingPosition: "LATEST",
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "SASL_SCRAM_512_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.SourceAccessConfigurations) != 1 {
		t.Errorf("expected 1 SourceAccessConfiguration, got %d", len(esm.SourceAccessConfigurations))
	}
	if esm.SourceAccessConfigurations[0].Type != "SASL_SCRAM_512_AUTH" {
		t.Errorf("unexpected Type: %s", esm.SourceAccessConfigurations[0].Type)
	}
	if esm.SourceAccessConfigurations[0].URI != "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret" {
		t.Errorf("unexpected URI: %v", esm.SourceAccessConfigurations[0].URI)
	}
}

func TestMSKEventProperties_ToEventSourceMapping_WithMultipleSourceAccessConfigs(t *testing.T) {
	props := &MSKEventProperties{
		Stream:           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:           []string{"my-topic"},
		StartingPosition: "LATEST",
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "SASL_SCRAM_512_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
			{
				Type: "VPC_SUBNET",
				URI:  "subnet-12345678",
			},
			{
				Type: "VPC_SECURITY_GROUP",
				URI:  "sg-12345678",
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.SourceAccessConfigurations) != 3 {
		t.Errorf("expected 3 SourceAccessConfigurations, got %d", len(esm.SourceAccessConfigurations))
	}
}

func TestMSKEventProperties_ToEventSourceMapping_WithFilterCriteria(t *testing.T) {
	props := &MSKEventProperties{
		Stream:           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:           []string{"my-topic"},
		StartingPosition: "LATEST",
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"value": {"type": ["ORDER"]}}`},
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

func TestMSKEventProperties_ToEventSourceMapping_FullConfig(t *testing.T) {
	batchSize := 100
	window := 30
	enabled := true

	props := &MSKEventProperties{
		Stream:                         "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:                         []string{"topic-1", "topic-2"},
		StartingPosition:               "LATEST",
		BatchSize:                      &batchSize,
		MaximumBatchingWindowInSeconds: &window,
		Enabled:                        &enabled,
		ConsumerGroupId:                "my-consumer-group",
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "SASL_SCRAM_512_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"value": {"priority": ["HIGH"]}}`},
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
	if *esm.MaximumBatchingWindowInSeconds != 30 {
		t.Errorf("unexpected MaximumBatchingWindowInSeconds: %d", *esm.MaximumBatchingWindowInSeconds)
	}
	if !*esm.Enabled {
		t.Error("unexpected Enabled value")
	}
	if esm.AmazonManagedKafkaEventSourceConfig.ConsumerGroupId != "my-consumer-group" {
		t.Errorf("unexpected ConsumerGroupId: %s", esm.AmazonManagedKafkaEventSourceConfig.ConsumerGroupId)
	}
	if len(esm.SourceAccessConfigurations) != 1 {
		t.Errorf("unexpected SourceAccessConfigurations length: %d", len(esm.SourceAccessConfigurations))
	}
	if len(esm.FilterCriteria.Filters) != 1 {
		t.Errorf("unexpected FilterCriteria.Filters length: %d", len(esm.FilterCriteria.Filters))
	}
}

func TestMSKEventProperties_ToEventSourceMapping_ToCloudFormation(t *testing.T) {
	batchSize := 100
	props := &MSKEventProperties{
		Stream:           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1",
		Topics:           []string{"my-topic"},
		StartingPosition: "LATEST",
		BatchSize:        &batchSize,
		ConsumerGroupId:  "my-consumer-group",
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
	if cfnProps["EventSourceArn"] != "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/12345678-1234-1234-1234-123456789012-1" {
		t.Errorf("expected EventSourceArn in CFN properties")
	}
	if cfnProps["StartingPosition"] != "LATEST" {
		t.Errorf("expected StartingPosition in CFN properties")
	}
	topics := cfnProps["Topics"].([]string)
	if len(topics) != 1 || topics[0] != "my-topic" {
		t.Errorf("unexpected Topics in CFN properties: %v", topics)
	}
	mskConfig := cfnProps["AmazonManagedKafkaEventSourceConfig"].(map[string]interface{})
	if mskConfig["ConsumerGroupId"] != "my-consumer-group" {
		t.Errorf("unexpected ConsumerGroupId in CFN properties: %v", mskConfig["ConsumerGroupId"])
	}
}
