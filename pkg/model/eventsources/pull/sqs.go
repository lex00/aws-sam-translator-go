package pull

import (
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// SQSEventProperties represents SAM properties for an SQS queue event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-sqs.html
type SQSEventProperties struct {
	// Queue is the ARN of the SQS queue (required).
	Queue interface{} `json:"Queue" yaml:"Queue"`

	// BatchSize is the maximum number of records in each batch (1-10000, default 10).
	BatchSize *int `json:"BatchSize,omitempty" yaml:"BatchSize,omitempty"`

	// MaximumBatchingWindowInSeconds is the maximum batching window in seconds (0-300).
	MaximumBatchingWindowInSeconds *int `json:"MaximumBatchingWindowInSeconds,omitempty" yaml:"MaximumBatchingWindowInSeconds,omitempty"`

	// Enabled indicates whether the event source mapping is active.
	Enabled *bool `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`

	// FunctionResponseTypes specifies the response type for the event source.
	// Valid values: ReportBatchItemFailures
	FunctionResponseTypes []string `json:"FunctionResponseTypes,omitempty" yaml:"FunctionResponseTypes,omitempty"`

	// FilterCriteria specifies event filtering for the mapping.
	FilterCriteria *FilterCriteria `json:"FilterCriteria,omitempty" yaml:"FilterCriteria,omitempty"`

	// ScalingConfig configures scaling for the event source mapping.
	ScalingConfig *ScalingConfig `json:"ScalingConfig,omitempty" yaml:"ScalingConfig,omitempty"`
}

// ScalingConfig configures scaling for the event source mapping.
type ScalingConfig struct {
	// MaximumConcurrency is the maximum number of concurrent functions (2-1000).
	MaximumConcurrency *int `json:"MaximumConcurrency,omitempty" yaml:"MaximumConcurrency,omitempty"`
}

// NewSQSEventProperties creates a new SQSEventProperties with required fields.
func NewSQSEventProperties(queue interface{}) *SQSEventProperties {
	return &SQSEventProperties{
		Queue: queue,
	}
}

// ToEventSourceMapping converts SQSEventProperties to a Lambda EventSourceMapping.
func (s *SQSEventProperties) ToEventSourceMapping(functionName interface{}) *lambda.EventSourceMapping {
	esm := lambda.NewSQSEventSourceMapping(functionName, s.Queue)

	if s.BatchSize != nil {
		esm.WithBatchSize(*s.BatchSize)
	}
	if s.MaximumBatchingWindowInSeconds != nil {
		esm.WithBatchingWindow(*s.MaximumBatchingWindowInSeconds)
	}
	if s.Enabled != nil {
		esm.WithEnabled(*s.Enabled)
	}
	if len(s.FunctionResponseTypes) > 0 {
		esm.FunctionResponseTypes = s.FunctionResponseTypes
	}
	if s.FilterCriteria != nil && len(s.FilterCriteria.Filters) > 0 {
		for _, filter := range s.FilterCriteria.Filters {
			esm.WithFilter(filter.Pattern)
		}
	}
	if s.ScalingConfig != nil && s.ScalingConfig.MaximumConcurrency != nil {
		esm.WithScalingConfig(*s.ScalingConfig.MaximumConcurrency)
	}

	return esm
}
