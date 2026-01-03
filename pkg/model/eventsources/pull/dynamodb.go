package pull

import (
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// DynamoDBEventProperties represents SAM properties for a DynamoDB Streams event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-dynamodb.html
type DynamoDBEventProperties struct {
	// Stream is the ARN of the DynamoDB stream (required).
	Stream interface{} `json:"Stream" yaml:"Stream"`

	// StartingPosition specifies the position in the stream to start reading (required).
	// Valid values: TRIM_HORIZON, LATEST
	StartingPosition string `json:"StartingPosition" yaml:"StartingPosition"`

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

// NewDynamoDBEventProperties creates a new DynamoDBEventProperties with required fields.
func NewDynamoDBEventProperties(stream interface{}, startingPosition string) *DynamoDBEventProperties {
	return &DynamoDBEventProperties{
		Stream:           stream,
		StartingPosition: startingPosition,
	}
}

// ToEventSourceMapping converts DynamoDBEventProperties to a Lambda EventSourceMapping.
func (d *DynamoDBEventProperties) ToEventSourceMapping(functionName interface{}) *lambda.EventSourceMapping {
	esm := lambda.NewDynamoDBEventSourceMapping(functionName, d.Stream, d.StartingPosition)

	if d.BatchSize != nil {
		esm.WithBatchSize(*d.BatchSize)
	}
	if d.MaximumBatchingWindowInSeconds != nil {
		esm.WithBatchingWindow(*d.MaximumBatchingWindowInSeconds)
	}
	if d.Enabled != nil {
		esm.WithEnabled(*d.Enabled)
	}
	if d.BisectBatchOnFunctionError != nil {
		esm.WithBisectOnError(*d.BisectBatchOnFunctionError)
	}
	if d.MaximumRetryAttempts != nil {
		esm.WithMaximumRetryAttempts(*d.MaximumRetryAttempts)
	}
	if d.MaximumRecordAgeInSeconds != nil {
		esm.WithMaximumRecordAge(*d.MaximumRecordAgeInSeconds)
	}
	if d.ParallelizationFactor != nil {
		esm.WithParallelizationFactor(*d.ParallelizationFactor)
	}
	if d.TumblingWindowInSeconds != nil {
		esm.WithTumblingWindow(*d.TumblingWindowInSeconds)
	}
	if d.DestinationConfig != nil && d.DestinationConfig.OnFailure != nil {
		esm.WithOnFailureDestination(d.DestinationConfig.OnFailure.Destination)
	}
	if len(d.FunctionResponseTypes) > 0 {
		esm.FunctionResponseTypes = d.FunctionResponseTypes
	}
	if d.FilterCriteria != nil && len(d.FilterCriteria.Filters) > 0 {
		for _, filter := range d.FilterCriteria.Filters {
			esm.WithFilter(filter.Pattern)
		}
	}

	return esm
}
