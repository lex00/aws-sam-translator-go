package pull

import (
	"testing"
)

func TestNewMQEventProperties(t *testing.T) {
	sourceAccessConfigs := []SourceAccessConfiguration{
		{
			Type: "BASIC_AUTH",
			URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
		},
	}
	props := NewMQEventProperties(
		"arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		[]string{"my-queue"},
		sourceAccessConfigs,
	)

	if props.Broker != "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012" {
		t.Errorf("expected Broker ARN, got %v", props.Broker)
	}
	if len(props.Queues) != 1 || props.Queues[0] != "my-queue" {
		t.Errorf("expected Queues [my-queue], got %v", props.Queues)
	}
	if len(props.SourceAccessConfigurations) != 1 {
		t.Errorf("expected 1 SourceAccessConfiguration, got %d", len(props.SourceAccessConfigurations))
	}
}

func TestMQEventProperties_ToEventSourceMapping_Minimal(t *testing.T) {
	sourceAccessConfigs := []SourceAccessConfiguration{
		{
			Type: "BASIC_AUTH",
			URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
		},
	}
	props := NewMQEventProperties(
		"arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		[]string{"my-queue"},
		sourceAccessConfigs,
	)

	esm := props.ToEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012" {
		t.Errorf("expected EventSourceArn, got %v", esm.EventSourceArn)
	}
	if len(esm.Queues) != 1 || esm.Queues[0] != "my-queue" {
		t.Errorf("expected Queues [my-queue], got %v", esm.Queues)
	}
	if len(esm.SourceAccessConfigurations) != 1 {
		t.Errorf("expected 1 SourceAccessConfiguration, got %d", len(esm.SourceAccessConfigurations))
	}
}

func TestMQEventProperties_ToEventSourceMapping_MultipleQueues(t *testing.T) {
	sourceAccessConfigs := []SourceAccessConfiguration{
		{
			Type: "BASIC_AUTH",
			URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
		},
	}
	props := NewMQEventProperties(
		"arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		[]string{"queue-1", "queue-2"},
		sourceAccessConfigs,
	)

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.Queues) != 2 {
		t.Errorf("expected 2 queues, got %d", len(esm.Queues))
	}
}

func TestMQEventProperties_ToEventSourceMapping_WithBatchSize(t *testing.T) {
	batchSize := 50
	props := &MQEventProperties{
		Broker: "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		Queues: []string{"my-queue"},
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
		BatchSize: &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.BatchSize == nil || *esm.BatchSize != 50 {
		t.Errorf("expected BatchSize 50, got %v", esm.BatchSize)
	}
}

func TestMQEventProperties_ToEventSourceMapping_WithBatchingWindow(t *testing.T) {
	window := 30
	props := &MQEventProperties{
		Broker: "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		Queues: []string{"my-queue"},
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
		MaximumBatchingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 30 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 30, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestMQEventProperties_ToEventSourceMapping_WithEnabled(t *testing.T) {
	enabled := false
	props := &MQEventProperties{
		Broker: "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		Queues: []string{"my-queue"},
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
		Enabled: &enabled,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestMQEventProperties_ToEventSourceMapping_WithVPCConfig(t *testing.T) {
	props := &MQEventProperties{
		Broker: "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		Queues: []string{"my-queue"},
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
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

func TestMQEventProperties_ToEventSourceMapping_WithFilterCriteria(t *testing.T) {
	props := &MQEventProperties{
		Broker: "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		Queues: []string{"my-queue"},
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"data": {"type": ["ORDER"]}}`},
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

func TestMQEventProperties_ToEventSourceMapping_FullConfig(t *testing.T) {
	batchSize := 100
	window := 15
	enabled := true

	props := &MQEventProperties{
		Broker: "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		Queues: []string{"queue-1", "queue-2"},
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
		BatchSize:                      &batchSize,
		MaximumBatchingWindowInSeconds: &window,
		Enabled:                        &enabled,
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"data": {"priority": ["HIGH"]}}`},
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
	if *esm.MaximumBatchingWindowInSeconds != 15 {
		t.Errorf("unexpected MaximumBatchingWindowInSeconds: %d", *esm.MaximumBatchingWindowInSeconds)
	}
	if !*esm.Enabled {
		t.Error("unexpected Enabled value")
	}
	if len(esm.Queues) != 2 {
		t.Errorf("unexpected Queues length: %d", len(esm.Queues))
	}
	if len(esm.SourceAccessConfigurations) != 1 {
		t.Errorf("unexpected SourceAccessConfigurations length: %d", len(esm.SourceAccessConfigurations))
	}
	if len(esm.FilterCriteria.Filters) != 1 {
		t.Errorf("unexpected FilterCriteria.Filters length: %d", len(esm.FilterCriteria.Filters))
	}
}

func TestMQEventProperties_ToEventSourceMapping_ToCloudFormation(t *testing.T) {
	batchSize := 10
	props := &MQEventProperties{
		Broker: "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012",
		Queues: []string{"my-queue"},
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
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
	if cfnProps["EventSourceArn"] != "arn:aws:mq:us-east-1:123456789012:broker:my-broker:b-12345678-1234-1234-1234-123456789012" {
		t.Errorf("expected EventSourceArn in CFN properties")
	}

	queues := cfnProps["Queues"].([]string)
	if len(queues) != 1 || queues[0] != "my-queue" {
		t.Errorf("unexpected Queues in CFN properties: %v", queues)
	}

	configs := cfnProps["SourceAccessConfigurations"].([]map[string]interface{})
	if len(configs) != 1 {
		t.Errorf("expected 1 SourceAccessConfiguration in CFN, got %d", len(configs))
	}
	if configs[0]["Type"] != "BASIC_AUTH" {
		t.Errorf("unexpected Type in SourceAccessConfiguration: %v", configs[0]["Type"])
	}
}

func TestMQEventProperties_RabbitMQ(t *testing.T) {
	sourceAccessConfigs := []SourceAccessConfiguration{
		{
			Type: "BASIC_AUTH",
			URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
		},
		{
			Type: "VIRTUAL_HOST",
			URI:  "/my-vhost",
		},
	}
	props := NewMQEventProperties(
		"arn:aws:mq:us-east-1:123456789012:broker:my-rabbit-broker:b-12345678-1234-1234-1234-123456789012",
		[]string{"my-queue"},
		sourceAccessConfigs,
	)

	esm := props.ToEventSourceMapping("my-function")

	// RabbitMQ uses VIRTUAL_HOST for vhost specification
	hasVirtualHost := false
	for _, sac := range esm.SourceAccessConfigurations {
		if sac.Type == "VIRTUAL_HOST" {
			hasVirtualHost = true
			break
		}
	}
	if !hasVirtualHost {
		t.Error("expected VIRTUAL_HOST in SourceAccessConfigurations for RabbitMQ")
	}
}
