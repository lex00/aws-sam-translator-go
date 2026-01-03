package lambda

// EventSourceMapping represents an AWS::Lambda::EventSourceMapping CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-eventsourcemapping.html
type EventSourceMapping struct {
	// AmazonManagedKafkaEventSourceConfig configures Amazon MSK event source.
	AmazonManagedKafkaEventSourceConfig *AmazonManagedKafkaEventSourceConfig `json:"AmazonManagedKafkaEventSourceConfig,omitempty" yaml:"AmazonManagedKafkaEventSourceConfig,omitempty"`

	// BatchSize is the maximum number of records in each batch.
	BatchSize *int `json:"BatchSize,omitempty" yaml:"BatchSize,omitempty"`

	// BisectBatchOnFunctionError splits a batch when a function returns an error.
	BisectBatchOnFunctionError bool `json:"BisectBatchOnFunctionError,omitempty" yaml:"BisectBatchOnFunctionError,omitempty"`

	// DestinationConfig configures destinations for events that fail processing.
	DestinationConfig *EventSourceDestinationConfig `json:"DestinationConfig,omitempty" yaml:"DestinationConfig,omitempty"`

	// DocumentDBEventSourceConfig configures DocumentDB event source.
	DocumentDBEventSourceConfig *DocumentDBEventSourceConfig `json:"DocumentDBEventSourceConfig,omitempty" yaml:"DocumentDBEventSourceConfig,omitempty"`

	// Enabled indicates whether the event source mapping is active.
	Enabled *bool `json:"Enabled,omitempty" yaml:"Enabled,omitempty"`

	// EventSourceArn is the ARN of the event source.
	EventSourceArn interface{} `json:"EventSourceArn,omitempty" yaml:"EventSourceArn,omitempty"`

	// FilterCriteria specifies event filtering for the mapping.
	FilterCriteria *FilterCriteria `json:"FilterCriteria,omitempty" yaml:"FilterCriteria,omitempty"`

	// FunctionName is the name or ARN of the Lambda function (required).
	FunctionName interface{} `json:"FunctionName" yaml:"FunctionName"`

	// FunctionResponseTypes specifies the response type for stream sources.
	// Valid values: ReportBatchItemFailures
	FunctionResponseTypes []string `json:"FunctionResponseTypes,omitempty" yaml:"FunctionResponseTypes,omitempty"`

	// KmsKeyArn is the ARN of a KMS key to encrypt the event source mapping.
	KmsKeyArn interface{} `json:"KmsKeyArn,omitempty" yaml:"KmsKeyArn,omitempty"`

	// MaximumBatchingWindowInSeconds is the maximum batching window in seconds.
	MaximumBatchingWindowInSeconds *int `json:"MaximumBatchingWindowInSeconds,omitempty" yaml:"MaximumBatchingWindowInSeconds,omitempty"`

	// MaximumRecordAgeInSeconds discards records older than the specified age.
	MaximumRecordAgeInSeconds *int `json:"MaximumRecordAgeInSeconds,omitempty" yaml:"MaximumRecordAgeInSeconds,omitempty"`

	// MaximumRetryAttempts is the maximum number of retry attempts.
	MaximumRetryAttempts *int `json:"MaximumRetryAttempts,omitempty" yaml:"MaximumRetryAttempts,omitempty"`

	// ParallelizationFactor is the number of batches to process concurrently.
	ParallelizationFactor *int `json:"ParallelizationFactor,omitempty" yaml:"ParallelizationFactor,omitempty"`

	// Queues is a list of Amazon MQ queues.
	Queues []string `json:"Queues,omitempty" yaml:"Queues,omitempty"`

	// ScalingConfig configures scaling for the event source mapping.
	ScalingConfig *ScalingConfig `json:"ScalingConfig,omitempty" yaml:"ScalingConfig,omitempty"`

	// SelfManagedEventSource configures a self-managed Apache Kafka source.
	SelfManagedEventSource *SelfManagedEventSource `json:"SelfManagedEventSource,omitempty" yaml:"SelfManagedEventSource,omitempty"`

	// SelfManagedKafkaEventSourceConfig configures self-managed Kafka source.
	SelfManagedKafkaEventSourceConfig *SelfManagedKafkaEventSourceConfig `json:"SelfManagedKafkaEventSourceConfig,omitempty" yaml:"SelfManagedKafkaEventSourceConfig,omitempty"`

	// SourceAccessConfigurations specifies authentication protocols.
	SourceAccessConfigurations []SourceAccessConfiguration `json:"SourceAccessConfigurations,omitempty" yaml:"SourceAccessConfigurations,omitempty"`

	// StartingPosition is the position in the stream to start reading.
	// Valid values: AT_TIMESTAMP, LATEST, TRIM_HORIZON
	StartingPosition string `json:"StartingPosition,omitempty" yaml:"StartingPosition,omitempty"`

	// StartingPositionTimestamp is the timestamp to start reading (for AT_TIMESTAMP).
	StartingPositionTimestamp *float64 `json:"StartingPositionTimestamp,omitempty" yaml:"StartingPositionTimestamp,omitempty"`

	// Tags is a list of tags for the event source mapping.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// Topics is a list of Kafka topics.
	Topics []string `json:"Topics,omitempty" yaml:"Topics,omitempty"`

	// TumblingWindowInSeconds is the duration of a processing window in seconds.
	TumblingWindowInSeconds *int `json:"TumblingWindowInSeconds,omitempty" yaml:"TumblingWindowInSeconds,omitempty"`
}

// AmazonManagedKafkaEventSourceConfig configures Amazon MSK event source.
type AmazonManagedKafkaEventSourceConfig struct {
	// ConsumerGroupId is the identifier for the Kafka consumer group.
	ConsumerGroupId string `json:"ConsumerGroupId,omitempty" yaml:"ConsumerGroupId,omitempty"`
}

// EventSourceDestinationConfig configures destinations for failed events.
type EventSourceDestinationConfig struct {
	// OnFailure configures the failure destination.
	OnFailure *OnFailure `json:"OnFailure,omitempty" yaml:"OnFailure,omitempty"`
}

// OnFailure specifies the destination for records that fail processing.
type OnFailure struct {
	// Destination is the ARN of the destination resource.
	Destination interface{} `json:"Destination,omitempty" yaml:"Destination,omitempty"`
}

// DocumentDBEventSourceConfig configures DocumentDB event source.
type DocumentDBEventSourceConfig struct {
	// CollectionName is the name of the collection to consume.
	CollectionName string `json:"CollectionName,omitempty" yaml:"CollectionName,omitempty"`

	// DatabaseName is the name of the database to consume.
	DatabaseName string `json:"DatabaseName,omitempty" yaml:"DatabaseName,omitempty"`

	// FullDocument specifies what to include in the change stream document.
	// Valid values: Default, UpdateLookup
	FullDocument string `json:"FullDocument,omitempty" yaml:"FullDocument,omitempty"`
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

// ScalingConfig configures scaling for the event source mapping.
type ScalingConfig struct {
	// MaximumConcurrency is the maximum number of concurrent functions.
	MaximumConcurrency *int `json:"MaximumConcurrency,omitempty" yaml:"MaximumConcurrency,omitempty"`
}

// SelfManagedEventSource configures a self-managed Apache Kafka source.
type SelfManagedEventSource struct {
	// Endpoints specifies the list of bootstrap servers.
	Endpoints *Endpoints `json:"Endpoints,omitempty" yaml:"Endpoints,omitempty"`
}

// Endpoints contains the list of Kafka bootstrap servers.
type Endpoints struct {
	// KafkaBootstrapServers is a list of Kafka bootstrap servers.
	KafkaBootstrapServers []string `json:"KafkaBootstrapServers,omitempty" yaml:"KafkaBootstrapServers,omitempty"`
}

// SelfManagedKafkaEventSourceConfig configures self-managed Kafka source.
type SelfManagedKafkaEventSourceConfig struct {
	// ConsumerGroupId is the identifier for the Kafka consumer group.
	ConsumerGroupId string `json:"ConsumerGroupId,omitempty" yaml:"ConsumerGroupId,omitempty"`
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

// NewEventSourceMapping creates a new EventSourceMapping with the required function name.
func NewEventSourceMapping(functionName interface{}) *EventSourceMapping {
	return &EventSourceMapping{
		FunctionName: functionName,
	}
}

// NewKinesisEventSourceMapping creates an event source mapping for Kinesis.
func NewKinesisEventSourceMapping(functionName interface{}, eventSourceArn interface{}, startingPosition string) *EventSourceMapping {
	return &EventSourceMapping{
		FunctionName:     functionName,
		EventSourceArn:   eventSourceArn,
		StartingPosition: startingPosition,
	}
}

// NewDynamoDBEventSourceMapping creates an event source mapping for DynamoDB Streams.
func NewDynamoDBEventSourceMapping(functionName interface{}, eventSourceArn interface{}, startingPosition string) *EventSourceMapping {
	return &EventSourceMapping{
		FunctionName:     functionName,
		EventSourceArn:   eventSourceArn,
		StartingPosition: startingPosition,
	}
}

// NewSQSEventSourceMapping creates an event source mapping for SQS.
func NewSQSEventSourceMapping(functionName interface{}, eventSourceArn interface{}) *EventSourceMapping {
	return &EventSourceMapping{
		FunctionName:   functionName,
		EventSourceArn: eventSourceArn,
	}
}

// NewMSKEventSourceMapping creates an event source mapping for Amazon MSK.
func NewMSKEventSourceMapping(functionName interface{}, eventSourceArn interface{}, topics []string, startingPosition string) *EventSourceMapping {
	return &EventSourceMapping{
		FunctionName:     functionName,
		EventSourceArn:   eventSourceArn,
		Topics:           topics,
		StartingPosition: startingPosition,
	}
}

// NewSelfManagedKafkaEventSourceMapping creates an event source mapping for self-managed Kafka.
func NewSelfManagedKafkaEventSourceMapping(functionName interface{}, bootstrapServers []string, topics []string, startingPosition string) *EventSourceMapping {
	return &EventSourceMapping{
		FunctionName: functionName,
		SelfManagedEventSource: &SelfManagedEventSource{
			Endpoints: &Endpoints{
				KafkaBootstrapServers: bootstrapServers,
			},
		},
		Topics:           topics,
		StartingPosition: startingPosition,
	}
}

// WithBatchSize sets the batch size for the mapping.
func (e *EventSourceMapping) WithBatchSize(size int) *EventSourceMapping {
	e.BatchSize = &size
	return e
}

// WithBatchingWindow sets the maximum batching window.
func (e *EventSourceMapping) WithBatchingWindow(seconds int) *EventSourceMapping {
	e.MaximumBatchingWindowInSeconds = &seconds
	return e
}

// WithEnabled sets whether the mapping is enabled.
func (e *EventSourceMapping) WithEnabled(enabled bool) *EventSourceMapping {
	e.Enabled = &enabled
	return e
}

// WithBisectOnError enables batch bisection on function error.
func (e *EventSourceMapping) WithBisectOnError(bisect bool) *EventSourceMapping {
	e.BisectBatchOnFunctionError = bisect
	return e
}

// WithMaximumRetryAttempts sets the maximum retry attempts.
func (e *EventSourceMapping) WithMaximumRetryAttempts(attempts int) *EventSourceMapping {
	e.MaximumRetryAttempts = &attempts
	return e
}

// WithMaximumRecordAge sets the maximum record age in seconds.
func (e *EventSourceMapping) WithMaximumRecordAge(seconds int) *EventSourceMapping {
	e.MaximumRecordAgeInSeconds = &seconds
	return e
}

// WithParallelizationFactor sets the parallelization factor.
func (e *EventSourceMapping) WithParallelizationFactor(factor int) *EventSourceMapping {
	e.ParallelizationFactor = &factor
	return e
}

// WithTumblingWindow sets the tumbling window duration.
func (e *EventSourceMapping) WithTumblingWindow(seconds int) *EventSourceMapping {
	e.TumblingWindowInSeconds = &seconds
	return e
}

// WithOnFailureDestination sets the failure destination.
func (e *EventSourceMapping) WithOnFailureDestination(destination interface{}) *EventSourceMapping {
	e.DestinationConfig = &EventSourceDestinationConfig{
		OnFailure: &OnFailure{
			Destination: destination,
		},
	}
	return e
}

// WithReportBatchItemFailures enables reporting of batch item failures.
func (e *EventSourceMapping) WithReportBatchItemFailures() *EventSourceMapping {
	e.FunctionResponseTypes = []string{"ReportBatchItemFailures"}
	return e
}

// WithFilter adds a filter pattern to the mapping.
func (e *EventSourceMapping) WithFilter(pattern string) *EventSourceMapping {
	if e.FilterCriteria == nil {
		e.FilterCriteria = &FilterCriteria{}
	}
	e.FilterCriteria.Filters = append(e.FilterCriteria.Filters, Filter{Pattern: pattern})
	return e
}

// WithScalingConfig sets the scaling configuration.
func (e *EventSourceMapping) WithScalingConfig(maxConcurrency int) *EventSourceMapping {
	e.ScalingConfig = &ScalingConfig{
		MaximumConcurrency: &maxConcurrency,
	}
	return e
}

// AddSourceAccessConfiguration adds a source access configuration.
func (e *EventSourceMapping) AddSourceAccessConfiguration(configType string, uri interface{}) *EventSourceMapping {
	e.SourceAccessConfigurations = append(e.SourceAccessConfigurations, SourceAccessConfiguration{
		Type: configType,
		URI:  uri,
	})
	return e
}

// ToCloudFormation converts the EventSourceMapping to a CloudFormation resource.
func (e *EventSourceMapping) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	properties["FunctionName"] = e.FunctionName

	if e.AmazonManagedKafkaEventSourceConfig != nil {
		properties["AmazonManagedKafkaEventSourceConfig"] = e.AmazonManagedKafkaEventSourceConfig.toMap()
	}
	if e.BatchSize != nil {
		properties["BatchSize"] = *e.BatchSize
	}
	if e.BisectBatchOnFunctionError {
		properties["BisectBatchOnFunctionError"] = e.BisectBatchOnFunctionError
	}
	if e.DestinationConfig != nil {
		properties["DestinationConfig"] = e.DestinationConfig.toMap()
	}
	if e.DocumentDBEventSourceConfig != nil {
		properties["DocumentDBEventSourceConfig"] = e.DocumentDBEventSourceConfig.toMap()
	}
	if e.Enabled != nil {
		properties["Enabled"] = *e.Enabled
	}
	if e.EventSourceArn != nil {
		properties["EventSourceArn"] = e.EventSourceArn
	}
	if e.FilterCriteria != nil && len(e.FilterCriteria.Filters) > 0 {
		properties["FilterCriteria"] = e.FilterCriteria.toMap()
	}
	if len(e.FunctionResponseTypes) > 0 {
		properties["FunctionResponseTypes"] = e.FunctionResponseTypes
	}
	if e.KmsKeyArn != nil {
		properties["KmsKeyArn"] = e.KmsKeyArn
	}
	if e.MaximumBatchingWindowInSeconds != nil {
		properties["MaximumBatchingWindowInSeconds"] = *e.MaximumBatchingWindowInSeconds
	}
	if e.MaximumRecordAgeInSeconds != nil {
		properties["MaximumRecordAgeInSeconds"] = *e.MaximumRecordAgeInSeconds
	}
	if e.MaximumRetryAttempts != nil {
		properties["MaximumRetryAttempts"] = *e.MaximumRetryAttempts
	}
	if e.ParallelizationFactor != nil {
		properties["ParallelizationFactor"] = *e.ParallelizationFactor
	}
	if len(e.Queues) > 0 {
		properties["Queues"] = e.Queues
	}
	if e.ScalingConfig != nil {
		properties["ScalingConfig"] = e.ScalingConfig.toMap()
	}
	if e.SelfManagedEventSource != nil {
		properties["SelfManagedEventSource"] = e.SelfManagedEventSource.toMap()
	}
	if e.SelfManagedKafkaEventSourceConfig != nil {
		properties["SelfManagedKafkaEventSourceConfig"] = e.SelfManagedKafkaEventSourceConfig.toMap()
	}
	if len(e.SourceAccessConfigurations) > 0 {
		configs := make([]map[string]interface{}, len(e.SourceAccessConfigurations))
		for i, sac := range e.SourceAccessConfigurations {
			configs[i] = sac.toMap()
		}
		properties["SourceAccessConfigurations"] = configs
	}
	if e.StartingPosition != "" {
		properties["StartingPosition"] = e.StartingPosition
	}
	if e.StartingPositionTimestamp != nil {
		properties["StartingPositionTimestamp"] = *e.StartingPositionTimestamp
	}
	if len(e.Tags) > 0 {
		tags := make([]map[string]interface{}, len(e.Tags))
		for i, t := range e.Tags {
			tags[i] = t.toMap()
		}
		properties["Tags"] = tags
	}
	if len(e.Topics) > 0 {
		properties["Topics"] = e.Topics
	}
	if e.TumblingWindowInSeconds != nil {
		properties["TumblingWindowInSeconds"] = *e.TumblingWindowInSeconds
	}

	return map[string]interface{}{
		"Type":       ResourceTypeEventSourceMapping,
		"Properties": properties,
	}
}

func (a *AmazonManagedKafkaEventSourceConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if a.ConsumerGroupId != "" {
		m["ConsumerGroupId"] = a.ConsumerGroupId
	}
	return m
}

func (d *EventSourceDestinationConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if d.OnFailure != nil {
		m["OnFailure"] = d.OnFailure.toMap()
	}
	return m
}

func (o *OnFailure) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if o.Destination != nil {
		m["Destination"] = o.Destination
	}
	return m
}

func (d *DocumentDBEventSourceConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if d.CollectionName != "" {
		m["CollectionName"] = d.CollectionName
	}
	if d.DatabaseName != "" {
		m["DatabaseName"] = d.DatabaseName
	}
	if d.FullDocument != "" {
		m["FullDocument"] = d.FullDocument
	}
	return m
}

func (f *FilterCriteria) toMap() map[string]interface{} {
	if len(f.Filters) == 0 {
		return nil
	}
	filters := make([]map[string]interface{}, len(f.Filters))
	for i, filter := range f.Filters {
		filters[i] = map[string]interface{}{
			"Pattern": filter.Pattern,
		}
	}
	return map[string]interface{}{
		"Filters": filters,
	}
}

func (s *ScalingConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if s.MaximumConcurrency != nil {
		m["MaximumConcurrency"] = *s.MaximumConcurrency
	}
	return m
}

func (s *SelfManagedEventSource) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if s.Endpoints != nil {
		m["Endpoints"] = s.Endpoints.toMap()
	}
	return m
}

func (e *Endpoints) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if len(e.KafkaBootstrapServers) > 0 {
		m["KafkaBootstrapServers"] = e.KafkaBootstrapServers
	}
	return m
}

func (s *SelfManagedKafkaEventSourceConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if s.ConsumerGroupId != "" {
		m["ConsumerGroupId"] = s.ConsumerGroupId
	}
	return m
}

func (s *SourceAccessConfiguration) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if s.Type != "" {
		m["Type"] = s.Type
	}
	if s.URI != nil {
		m["URI"] = s.URI
	}
	return m
}
