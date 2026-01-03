package pull

import (
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// MSKEventProperties represents SAM properties for an Amazon MSK event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-msk.html
type MSKEventProperties struct {
	// Stream is the ARN of the Amazon MSK cluster (required).
	Stream interface{} `json:"Stream" yaml:"Stream"`

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

	// SourceAccessConfigurations specifies authentication configuration.
	SourceAccessConfigurations []SourceAccessConfiguration `json:"SourceAccessConfigurations,omitempty" yaml:"SourceAccessConfigurations,omitempty"`

	// FilterCriteria specifies event filtering for the mapping.
	FilterCriteria *FilterCriteria `json:"FilterCriteria,omitempty" yaml:"FilterCriteria,omitempty"`
}

// SourceAccessConfiguration specifies authentication configuration.
type SourceAccessConfiguration struct {
	// Type is the authentication protocol type.
	// Valid values: BASIC_AUTH, CLIENT_CERTIFICATE_TLS_AUTH, SASL_SCRAM_256_AUTH,
	// SASL_SCRAM_512_AUTH, SERVER_ROOT_CA_CERTIFICATE, VIRTUAL_HOST, VPC_SECURITY_GROUP, VPC_SUBNET
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`

	// URI is the value for the authentication configuration.
	URI interface{} `json:"URI,omitempty" yaml:"URI,omitempty"`
}

// NewMSKEventProperties creates a new MSKEventProperties with required fields.
func NewMSKEventProperties(stream interface{}, topics []string, startingPosition string) *MSKEventProperties {
	return &MSKEventProperties{
		Stream:           stream,
		Topics:           topics,
		StartingPosition: startingPosition,
	}
}

// ToEventSourceMapping converts MSKEventProperties to a Lambda EventSourceMapping.
func (m *MSKEventProperties) ToEventSourceMapping(functionName interface{}) *lambda.EventSourceMapping {
	esm := lambda.NewMSKEventSourceMapping(functionName, m.Stream, m.Topics, m.StartingPosition)

	if m.BatchSize != nil {
		esm.WithBatchSize(*m.BatchSize)
	}
	if m.MaximumBatchingWindowInSeconds != nil {
		esm.WithBatchingWindow(*m.MaximumBatchingWindowInSeconds)
	}
	if m.Enabled != nil {
		esm.WithEnabled(*m.Enabled)
	}
	if m.ConsumerGroupId != "" {
		esm.AmazonManagedKafkaEventSourceConfig = &lambda.AmazonManagedKafkaEventSourceConfig{
			ConsumerGroupId: m.ConsumerGroupId,
		}
	}
	for _, sac := range m.SourceAccessConfigurations {
		esm.AddSourceAccessConfiguration(sac.Type, sac.URI)
	}
	if m.FilterCriteria != nil && len(m.FilterCriteria.Filters) > 0 {
		for _, filter := range m.FilterCriteria.Filters {
			esm.WithFilter(filter.Pattern)
		}
	}

	return esm
}
