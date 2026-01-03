package pull

import (
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// DocumentDBEventProperties represents SAM properties for a DocumentDB change stream event source.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-documentdb.html
type DocumentDBEventProperties struct {
	// Cluster is the ARN of the Amazon DocumentDB cluster (required).
	Cluster interface{} `json:"Cluster" yaml:"Cluster"`

	// DatabaseName is the name of the database to consume (required).
	DatabaseName string `json:"DatabaseName" yaml:"DatabaseName"`

	// CollectionName is the name of the collection to consume.
	CollectionName string `json:"CollectionName,omitempty" yaml:"CollectionName,omitempty"`

	// StartingPosition specifies the position in the stream to start reading.
	// Valid values: TRIM_HORIZON, LATEST
	StartingPosition string `json:"StartingPosition,omitempty" yaml:"StartingPosition,omitempty"`

	// FullDocument specifies what to include in the change stream document.
	// Valid values: Default, UpdateLookup
	FullDocument string `json:"FullDocument,omitempty" yaml:"FullDocument,omitempty"`

	// BatchSize is the maximum number of records in each batch (1-10000, default 100).
	BatchSize *int `json:"BatchSize,omitempty" yaml:"BatchSize,omitempty"`

	// MaximumBatchingWindowInSeconds is the maximum batching window in seconds (0-300).
	MaximumBatchingWindowInSeconds *int `json:"MaximumBatchingWindowInSeconds,omitempty" yaml:"MaximumBatchingWindowInSeconds,omitempty"`

	// Enabled indicates whether the event source mapping is active.
	Enabled *bool `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`

	// SourceAccessConfigurations specifies authentication and VPC configuration.
	SourceAccessConfigurations []SourceAccessConfiguration `json:"SourceAccessConfigurations,omitempty" yaml:"SourceAccessConfigurations,omitempty"`

	// FilterCriteria specifies event filtering for the mapping.
	FilterCriteria *FilterCriteria `json:"FilterCriteria,omitempty" yaml:"FilterCriteria,omitempty"`
}

// NewDocumentDBEventProperties creates a new DocumentDBEventProperties with required fields.
func NewDocumentDBEventProperties(cluster interface{}, databaseName string) *DocumentDBEventProperties {
	return &DocumentDBEventProperties{
		Cluster:      cluster,
		DatabaseName: databaseName,
	}
}

// ToEventSourceMapping converts DocumentDBEventProperties to a Lambda EventSourceMapping.
func (d *DocumentDBEventProperties) ToEventSourceMapping(functionName interface{}) *lambda.EventSourceMapping {
	esm := lambda.NewEventSourceMapping(functionName)
	esm.EventSourceArn = d.Cluster

	// Configure DocumentDB-specific settings
	esm.DocumentDBEventSourceConfig = &lambda.DocumentDBEventSourceConfig{
		DatabaseName: d.DatabaseName,
	}
	if d.CollectionName != "" {
		esm.DocumentDBEventSourceConfig.CollectionName = d.CollectionName
	}
	if d.FullDocument != "" {
		esm.DocumentDBEventSourceConfig.FullDocument = d.FullDocument
	}

	if d.StartingPosition != "" {
		esm.StartingPosition = d.StartingPosition
	}
	if d.BatchSize != nil {
		esm.WithBatchSize(*d.BatchSize)
	}
	if d.MaximumBatchingWindowInSeconds != nil {
		esm.WithBatchingWindow(*d.MaximumBatchingWindowInSeconds)
	}
	if d.Enabled != nil {
		esm.WithEnabled(*d.Enabled)
	}
	for _, sac := range d.SourceAccessConfigurations {
		esm.AddSourceAccessConfiguration(sac.Type, sac.URI)
	}
	if d.FilterCriteria != nil && len(d.FilterCriteria.Filters) > 0 {
		for _, filter := range d.FilterCriteria.Filters {
			esm.WithFilter(filter.Pattern)
		}
	}

	return esm
}
