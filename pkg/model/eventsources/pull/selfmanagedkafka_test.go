package pull

import (
	"testing"
)

func TestNewSelfManagedKafkaEventProperties(t *testing.T) {
	props := NewSelfManagedKafkaEventProperties(
		[]string{"broker1:9092", "broker2:9092"},
		[]string{"my-topic"},
		"LATEST",
	)

	if len(props.KafkaBootstrapServers) != 2 {
		t.Errorf("expected 2 bootstrap servers, got %d", len(props.KafkaBootstrapServers))
	}
	if len(props.Topics) != 1 || props.Topics[0] != "my-topic" {
		t.Errorf("expected Topics [my-topic], got %v", props.Topics)
	}
	if props.StartingPosition != "LATEST" {
		t.Errorf("expected StartingPosition 'LATEST', got %s", props.StartingPosition)
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_Minimal(t *testing.T) {
	props := NewSelfManagedKafkaEventProperties(
		[]string{"broker1:9092"},
		[]string{"my-topic"},
		"TRIM_HORIZON",
	)

	esm := props.ToEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
	if esm.SelfManagedEventSource == nil {
		t.Fatal("expected SelfManagedEventSource to be set")
	}
	if len(esm.SelfManagedEventSource.Endpoints.KafkaBootstrapServers) != 1 {
		t.Errorf("expected 1 bootstrap server, got %d",
			len(esm.SelfManagedEventSource.Endpoints.KafkaBootstrapServers))
	}
	if esm.StartingPosition != "TRIM_HORIZON" {
		t.Errorf("expected StartingPosition 'TRIM_HORIZON', got %s", esm.StartingPosition)
	}
	if len(esm.Topics) != 1 || esm.Topics[0] != "my-topic" {
		t.Errorf("expected Topics [my-topic], got %v", esm.Topics)
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_MultipleServers(t *testing.T) {
	props := NewSelfManagedKafkaEventProperties(
		[]string{"broker1:9092", "broker2:9092", "broker3:9092"},
		[]string{"my-topic"},
		"LATEST",
	)

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.SelfManagedEventSource.Endpoints.KafkaBootstrapServers) != 3 {
		t.Errorf("expected 3 bootstrap servers, got %d",
			len(esm.SelfManagedEventSource.Endpoints.KafkaBootstrapServers))
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_WithBatchSize(t *testing.T) {
	batchSize := 500
	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers: []string{"broker1:9092"},
		Topics:                []string{"my-topic"},
		StartingPosition:      "LATEST",
		BatchSize:             &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.BatchSize == nil || *esm.BatchSize != 500 {
		t.Errorf("expected BatchSize 500, got %v", esm.BatchSize)
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_WithBatchingWindow(t *testing.T) {
	window := 30
	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers:          []string{"broker1:9092"},
		Topics:                         []string{"my-topic"},
		StartingPosition:               "LATEST",
		MaximumBatchingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 30 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 30, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_WithEnabled(t *testing.T) {
	enabled := false
	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers: []string{"broker1:9092"},
		Topics:                []string{"my-topic"},
		StartingPosition:      "LATEST",
		Enabled:               &enabled,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_WithConsumerGroupId(t *testing.T) {
	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers: []string{"broker1:9092"},
		Topics:                []string{"my-topic"},
		StartingPosition:      "LATEST",
		ConsumerGroupId:       "my-consumer-group",
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.SelfManagedKafkaEventSourceConfig == nil {
		t.Fatal("expected SelfManagedKafkaEventSourceConfig to be set")
	}
	if esm.SelfManagedKafkaEventSourceConfig.ConsumerGroupId != "my-consumer-group" {
		t.Errorf("expected ConsumerGroupId 'my-consumer-group', got %s",
			esm.SelfManagedKafkaEventSourceConfig.ConsumerGroupId)
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_WithSourceAccessConfigurations(t *testing.T) {
	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers: []string{"broker1:9092"},
		Topics:                []string{"my-topic"},
		StartingPosition:      "LATEST",
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
	if esm.SourceAccessConfigurations[0].Type != "SASL_SCRAM_512_AUTH" {
		t.Errorf("unexpected Type: %s", esm.SourceAccessConfigurations[0].Type)
	}
	if esm.SourceAccessConfigurations[1].Type != "VPC_SUBNET" {
		t.Errorf("unexpected Type: %s", esm.SourceAccessConfigurations[1].Type)
	}
	if esm.SourceAccessConfigurations[2].Type != "VPC_SECURITY_GROUP" {
		t.Errorf("unexpected Type: %s", esm.SourceAccessConfigurations[2].Type)
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_WithFilterCriteria(t *testing.T) {
	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers: []string{"broker1:9092"},
		Topics:                []string{"my-topic"},
		StartingPosition:      "LATEST",
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

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_FullConfig(t *testing.T) {
	batchSize := 100
	window := 30
	enabled := true

	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers:          []string{"broker1:9092", "broker2:9092"},
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
			{
				Type: "VPC_SUBNET",
				URI:  "subnet-12345678",
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
	if esm.SelfManagedKafkaEventSourceConfig.ConsumerGroupId != "my-consumer-group" {
		t.Errorf("unexpected ConsumerGroupId: %s", esm.SelfManagedKafkaEventSourceConfig.ConsumerGroupId)
	}
	if len(esm.SourceAccessConfigurations) != 2 {
		t.Errorf("unexpected SourceAccessConfigurations length: %d", len(esm.SourceAccessConfigurations))
	}
	if len(esm.FilterCriteria.Filters) != 1 {
		t.Errorf("unexpected FilterCriteria.Filters length: %d", len(esm.FilterCriteria.Filters))
	}
}

func TestSelfManagedKafkaEventProperties_ToEventSourceMapping_ToCloudFormation(t *testing.T) {
	batchSize := 100
	props := &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers: []string{"broker1:9092", "broker2:9092"},
		Topics:                []string{"my-topic"},
		StartingPosition:      "LATEST",
		BatchSize:             &batchSize,
		ConsumerGroupId:       "my-consumer-group",
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
	if cfnProps["StartingPosition"] != "LATEST" {
		t.Errorf("expected StartingPosition in CFN properties")
	}

	selfManaged := cfnProps["SelfManagedEventSource"].(map[string]interface{})
	endpoints := selfManaged["Endpoints"].(map[string]interface{})
	servers := endpoints["KafkaBootstrapServers"].([]string)
	if len(servers) != 2 {
		t.Errorf("expected 2 bootstrap servers in CFN, got %d", len(servers))
	}

	topics := cfnProps["Topics"].([]string)
	if len(topics) != 1 || topics[0] != "my-topic" {
		t.Errorf("unexpected Topics in CFN properties: %v", topics)
	}

	kafkaConfig := cfnProps["SelfManagedKafkaEventSourceConfig"].(map[string]interface{})
	if kafkaConfig["ConsumerGroupId"] != "my-consumer-group" {
		t.Errorf("unexpected ConsumerGroupId in CFN properties: %v", kafkaConfig["ConsumerGroupId"])
	}
}
