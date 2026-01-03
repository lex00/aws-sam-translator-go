// Package pull provides SAM event source handlers for pull-based event sources.
// These handlers convert SAM event properties to CloudFormation resources.
package pull

import (
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// KinesisEventProperties represents SAM properties for a Kinesis stream event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-kinesis.html
type KinesisEventProperties struct {
	// Stream is the ARN of the Kinesis stream (required).
	Stream interface{} `json:"Stream" yaml:"Stream"`

	// StartingPosition specifies the position in the stream to start reading (required).
	// Valid values: TRIM_HORIZON, LATEST, AT_TIMESTAMP
	StartingPosition string `json:"StartingPosition" yaml:"StartingPosition"`

	// StartingPositionTimestamp is the time from which to start reading (for AT_TIMESTAMP).
	StartingPositionTimestamp *float64 `json:"StartingPositionTimestamp,omitempty" yaml:"StartingPositionTimestamp,omitempty"`

	// BatchSize is the maximum number of records in each batch (1-10000, default 100).
	BatchSize *int `json:"BatchSize,omitempty" yaml:"BatchSize,omitempty"`

	// MaximumBatchingWindowInSeconds is the maximum batching window in seconds (0-300).
	MaximumBatchingWindowInSeconds *int `json:"MaximumBatchingWindowInSeconds,omitempty" yaml:"MaximumBatchingWindowInSeconds,omitempty"`

	// Enabled indicates whether the event source mapping is active.
	Enabled *bool `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`

	// BisectBatchOnFunctionError splits a batch when a function returns an error.
	BisectBatchOnFunctionError *bool `json:"BisectBatchOnFunctionError,omitempty" yaml:"BisectBatchOnFunctionError,omitempty"`

	// MaximumRetryAttempts is the maximum number of retry attempts (0-10000, -1 for infinite).
	MaximumRetryAttempts *int `json:"MaximumRetryAttempts,omitempty" yaml:"MaximumRetryAttempts,omitempty"`

	// MaximumRecordAgeInSeconds discards records older than the specified age (60-604800, -1 for infinite).
	MaximumRecordAgeInSeconds *int `json:"MaximumRecordAgeInSeconds,omitempty" yaml:"MaximumRecordAgeInSeconds,omitempty"`

	// ParallelizationFactor is the number of batches to process concurrently (1-10).
	ParallelizationFactor *int `json:"ParallelizationFactor,omitempty" yaml:"ParallelizationFactor,omitempty"`

	// TumblingWindowInSeconds is the duration of a processing window in seconds (0-900).
	TumblingWindowInSeconds *int `json:"TumblingWindowInSeconds,omitempty" yaml:"TumblingWindowInSeconds,omitempty"`

	// DestinationConfig configures destinations for events that fail processing.
	DestinationConfig *DestinationConfig `json:"DestinationConfig,omitempty" yaml:"DestinationConfig,omitempty"`

	// FunctionResponseTypes specifies the response type for stream sources.
	// Valid values: ReportBatchItemFailures
	FunctionResponseTypes []string `json:"FunctionResponseTypes,omitempty" yaml:"FunctionResponseTypes,omitempty"`

	// FilterCriteria specifies event filtering for the mapping.
	FilterCriteria *FilterCriteria `json:"FilterCriteria,omitempty" yaml:"FilterCriteria,omitempty"`
}

// DestinationConfig configures destinations for failed events.
type DestinationConfig struct {
	// OnFailure configures the failure destination.
	OnFailure *OnFailure `json:"OnFailure,omitempty" yaml:"OnFailure,omitempty"`
}

// OnFailure specifies the destination for records that fail processing.
type OnFailure struct {
	// Destination is the ARN of the destination resource.
	Destination interface{} `json:"Destination,omitempty" yaml:"Destination,omitempty"`
}

// FilterCriteria specifies event filtering configuration.
type FilterCriteria struct {
	// Filters is a list of filter patterns.
	Filters []Filter `json:"Filters,omitempty" yaml:"Filters,omitempty"`
}

// Filter represents a single filter pattern.
type Filter struct {
	// Pattern is the filter pattern in JSON format.
	Pattern string `json:"Pattern,omitempty" yaml:"Pattern,omitempty"`
}

// NewKinesisEventProperties creates a new KinesisEventProperties with required fields.
func NewKinesisEventProperties(stream interface{}, startingPosition string) *KinesisEventProperties {
	return &KinesisEventProperties{
		Stream:           stream,
		StartingPosition: startingPosition,
	}
}

// ToEventSourceMapping converts KinesisEventProperties to a Lambda EventSourceMapping.
func (k *KinesisEventProperties) ToEventSourceMapping(functionName interface{}) *lambda.EventSourceMapping {
	esm := lambda.NewKinesisEventSourceMapping(functionName, k.Stream, k.StartingPosition)

	if k.StartingPositionTimestamp != nil {
		esm.StartingPositionTimestamp = k.StartingPositionTimestamp
	}
	if k.BatchSize != nil {
		esm.WithBatchSize(*k.BatchSize)
	}
	if k.MaximumBatchingWindowInSeconds != nil {
		esm.WithBatchingWindow(*k.MaximumBatchingWindowInSeconds)
	}
	if k.Enabled != nil {
		esm.WithEnabled(*k.Enabled)
	}
	if k.BisectBatchOnFunctionError != nil {
		esm.WithBisectOnError(*k.BisectBatchOnFunctionError)
	}
	if k.MaximumRetryAttempts != nil {
		esm.WithMaximumRetryAttempts(*k.MaximumRetryAttempts)
	}
	if k.MaximumRecordAgeInSeconds != nil {
		esm.WithMaximumRecordAge(*k.MaximumRecordAgeInSeconds)
	}
	if k.ParallelizationFactor != nil {
		esm.WithParallelizationFactor(*k.ParallelizationFactor)
	}
	if k.TumblingWindowInSeconds != nil {
		esm.WithTumblingWindow(*k.TumblingWindowInSeconds)
	}
	if k.DestinationConfig != nil && k.DestinationConfig.OnFailure != nil {
		esm.WithOnFailureDestination(k.DestinationConfig.OnFailure.Destination)
	}
	if len(k.FunctionResponseTypes) > 0 {
		esm.FunctionResponseTypes = k.FunctionResponseTypes
	}
	if k.FilterCriteria != nil && len(k.FilterCriteria.Filters) > 0 {
		for _, filter := range k.FilterCriteria.Filters {
			esm.WithFilter(filter.Pattern)
		}
	}

	return esm
}
