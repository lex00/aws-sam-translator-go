package pull

import (
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// MQEventProperties represents SAM properties for an Amazon MQ event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-mq.html
type MQEventProperties struct {
	// Broker is the ARN of the Amazon MQ broker (required).
	Broker interface{} `json:"Broker" yaml:"Broker"`

	// Queues is a list of queue names to consume from (required).
	Queues []string `json:"Queues" yaml:"Queues"`

	// SourceAccessConfigurations specifies authentication configuration (required).
	SourceAccessConfigurations []SourceAccessConfiguration `json:"SourceAccessConfigurations" yaml:"SourceAccessConfigurations"`

	// BatchSize is the maximum number of records in each batch (1-10000, default 100).
	BatchSize *int `json:"BatchSize,omitempty" yaml:"BatchSize,omitempty"`

	// MaximumBatchingWindowInSeconds is the maximum batching window in seconds (0-300).
	MaximumBatchingWindowInSeconds *int `json:"MaximumBatchingWindowInSeconds,omitempty" yaml:"MaximumBatchingWindowInSeconds,omitempty"`

	// Enabled indicates whether the event source mapping is active.
	Enabled *bool `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`

	// FilterCriteria specifies event filtering for the mapping.
	FilterCriteria *FilterCriteria `json:"FilterCriteria,omitempty" yaml:"FilterCriteria,omitempty"`
}

// NewMQEventProperties creates a new MQEventProperties with required fields.
func NewMQEventProperties(broker interface{}, queues []string, sourceAccessConfigs []SourceAccessConfiguration) *MQEventProperties {
	return &MQEventProperties{
		Broker:                     broker,
		Queues:                     queues,
		SourceAccessConfigurations: sourceAccessConfigs,
	}
}

// ToEventSourceMapping converts MQEventProperties to a Lambda EventSourceMapping.
func (m *MQEventProperties) ToEventSourceMapping(functionName interface{}) *lambda.EventSourceMapping {
	esm := lambda.NewEventSourceMapping(functionName)
	esm.EventSourceArn = m.Broker
	esm.Queues = m.Queues

	if m.BatchSize != nil {
		esm.WithBatchSize(*m.BatchSize)
	}
	if m.MaximumBatchingWindowInSeconds != nil {
		esm.WithBatchingWindow(*m.MaximumBatchingWindowInSeconds)
	}
	if m.Enabled != nil {
		esm.WithEnabled(*m.Enabled)
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
