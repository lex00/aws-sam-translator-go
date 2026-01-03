package pull

import (
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// SelfManagedKafkaEventProperties represents SAM properties for a self-managed Kafka event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-selfmanagedkafka.html
type SelfManagedKafkaEventProperties struct {
	// KafkaBootstrapServers is a list of Kafka bootstrap servers (required).
	KafkaBootstrapServers []string `json:"KafkaBootstrapServers" yaml:"KafkaBootstrapServers"`

	// Topics is a list of Kafka topics (required).
	Topics []string `json:"Topics" yaml:"Topics"`

	// StartingPosition specifies the position in the stream to start reading (required).
	// Valid values: TRIM_HORIZON, LATEST
	StartingPosition string `json:"StartingPosition" yaml:"StartingPosition"`

	// BatchSize is the maximum number of records in each batch (1-10000, default 100).
	BatchSize *int `json:"BatchSize,omitempty" yaml:"BatchSize,omitempty"`

	// MaximumBatchingWindowInSeconds is the maximum batching window in seconds (0-300).
	MaximumBatchingWindowInSeconds *int `json:"MaximumBatchingWindowInSeconds,omitempty" yaml:"MaximumBatchingWindowInSeconds,omitempty"`

	// Enabled indicates whether the event source mapping is active.
	Enabled *bool `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`

	// ConsumerGroupId is the identifier for the Kafka consumer group.
	ConsumerGroupId string `json:"ConsumerGroupId,omitempty" yaml:"ConsumerGroupId,omitempty"`

	// SourceAccessConfigurations specifies authentication and VPC configuration.
	SourceAccessConfigurations []SourceAccessConfiguration `json:"SourceAccessConfigurations,omitempty" yaml:"SourceAccessConfigurations,omitempty"`

	// FilterCriteria specifies event filtering for the mapping.
	FilterCriteria *FilterCriteria `json:"FilterCriteria,omitempty" yaml:"FilterCriteria,omitempty"`
}

// NewSelfManagedKafkaEventProperties creates a new SelfManagedKafkaEventProperties with required fields.
func NewSelfManagedKafkaEventProperties(bootstrapServers []string, topics []string, startingPosition string) *SelfManagedKafkaEventProperties {
	return &SelfManagedKafkaEventProperties{
		KafkaBootstrapServers: bootstrapServers,
		Topics:                topics,
		StartingPosition:      startingPosition,
	}
}

// ToEventSourceMapping converts SelfManagedKafkaEventProperties to a Lambda EventSourceMapping.
func (s *SelfManagedKafkaEventProperties) ToEventSourceMapping(functionName interface{}) *lambda.EventSourceMapping {
	esm := lambda.NewSelfManagedKafkaEventSourceMapping(
		functionName,
		s.KafkaBootstrapServers,
		s.Topics,
		s.StartingPosition,
	)

	if s.BatchSize != nil {
		esm.WithBatchSize(*s.BatchSize)
	}
	if s.MaximumBatchingWindowInSeconds != nil {
		esm.WithBatchingWindow(*s.MaximumBatchingWindowInSeconds)
	}
	if s.Enabled != nil {
		esm.WithEnabled(*s.Enabled)
	}
	if s.ConsumerGroupId != "" {
		esm.SelfManagedKafkaEventSourceConfig = &lambda.SelfManagedKafkaEventSourceConfig{
			ConsumerGroupId: s.ConsumerGroupId,
		}
	}
	for _, sac := range s.SourceAccessConfigurations {
		esm.AddSourceAccessConfiguration(sac.Type, sac.URI)
	}
	if s.FilterCriteria != nil && len(s.FilterCriteria.Filters) > 0 {
		for _, filter := range s.FilterCriteria.Filters {
			esm.WithFilter(filter.Pattern)
		}
	}

	return esm
}
