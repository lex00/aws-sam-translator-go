package pull

import (
	"testing"
)

func TestNewDocumentDBEventProperties(t *testing.T) {
	props := NewDocumentDBEventProperties(
		"arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		"my-database",
	)

	if props.Cluster != "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster" {
		t.Errorf("expected Cluster ARN, got %v", props.Cluster)
	}
	if props.DatabaseName != "my-database" {
		t.Errorf("expected DatabaseName 'my-database', got %s", props.DatabaseName)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_Minimal(t *testing.T) {
	props := NewDocumentDBEventProperties(
		"arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		"my-database",
	)

	esm := props.ToEventSourceMapping("my-function")

	if esm.FunctionName != "my-function" {
		t.Errorf("expected FunctionName 'my-function', got %v", esm.FunctionName)
	}
	if esm.EventSourceArn != "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster" {
		t.Errorf("expected EventSourceArn, got %v", esm.EventSourceArn)
	}
	if esm.DocumentDBEventSourceConfig == nil {
		t.Fatal("expected DocumentDBEventSourceConfig to be set")
	}
	if esm.DocumentDBEventSourceConfig.DatabaseName != "my-database" {
		t.Errorf("expected DatabaseName 'my-database', got %s", esm.DocumentDBEventSourceConfig.DatabaseName)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithCollection(t *testing.T) {
	props := &DocumentDBEventProperties{
		Cluster:        "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName:   "my-database",
		CollectionName: "my-collection",
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.DocumentDBEventSourceConfig.CollectionName != "my-collection" {
		t.Errorf("expected CollectionName 'my-collection', got %s", esm.DocumentDBEventSourceConfig.CollectionName)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithFullDocument(t *testing.T) {
	props := &DocumentDBEventProperties{
		Cluster:      "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName: "my-database",
		FullDocument: "UpdateLookup",
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.DocumentDBEventSourceConfig.FullDocument != "UpdateLookup" {
		t.Errorf("expected FullDocument 'UpdateLookup', got %s", esm.DocumentDBEventSourceConfig.FullDocument)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithStartingPosition(t *testing.T) {
	props := &DocumentDBEventProperties{
		Cluster:          "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName:     "my-database",
		StartingPosition: "LATEST",
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.StartingPosition != "LATEST" {
		t.Errorf("expected StartingPosition 'LATEST', got %s", esm.StartingPosition)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithBatchSize(t *testing.T) {
	batchSize := 200
	props := &DocumentDBEventProperties{
		Cluster:      "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName: "my-database",
		BatchSize:    &batchSize,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.BatchSize == nil || *esm.BatchSize != 200 {
		t.Errorf("expected BatchSize 200, got %v", esm.BatchSize)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithBatchingWindow(t *testing.T) {
	window := 30
	props := &DocumentDBEventProperties{
		Cluster:                        "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName:                   "my-database",
		MaximumBatchingWindowInSeconds: &window,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.MaximumBatchingWindowInSeconds == nil || *esm.MaximumBatchingWindowInSeconds != 30 {
		t.Errorf("expected MaximumBatchingWindowInSeconds 30, got %v", esm.MaximumBatchingWindowInSeconds)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithEnabled(t *testing.T) {
	enabled := false
	props := &DocumentDBEventProperties{
		Cluster:      "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName: "my-database",
		Enabled:      &enabled,
	}

	esm := props.ToEventSourceMapping("my-function")

	if esm.Enabled == nil || *esm.Enabled != false {
		t.Errorf("expected Enabled false, got %v", esm.Enabled)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithSourceAccessConfigurations(t *testing.T) {
	props := &DocumentDBEventProperties{
		Cluster:      "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName: "my-database",
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	if len(esm.SourceAccessConfigurations) != 1 {
		t.Errorf("expected 1 SourceAccessConfiguration, got %d", len(esm.SourceAccessConfigurations))
	}
	if esm.SourceAccessConfigurations[0].Type != "BASIC_AUTH" {
		t.Errorf("unexpected Type: %s", esm.SourceAccessConfigurations[0].Type)
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_WithFilterCriteria(t *testing.T) {
	props := &DocumentDBEventProperties{
		Cluster:      "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName: "my-database",
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"operationType": ["insert", "update"]}`},
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

func TestDocumentDBEventProperties_ToEventSourceMapping_FullConfig(t *testing.T) {
	batchSize := 100
	window := 15
	enabled := true

	props := &DocumentDBEventProperties{
		Cluster:                        "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName:                   "my-database",
		CollectionName:                 "my-collection",
		StartingPosition:               "TRIM_HORIZON",
		FullDocument:                   "UpdateLookup",
		BatchSize:                      &batchSize,
		MaximumBatchingWindowInSeconds: &window,
		Enabled:                        &enabled,
		SourceAccessConfigurations: []SourceAccessConfiguration{
			{
				Type: "BASIC_AUTH",
				URI:  "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret",
			},
		},
		FilterCriteria: &FilterCriteria{
			Filters: []Filter{
				{Pattern: `{"operationType": ["insert"]}`},
			},
		},
	}

	esm := props.ToEventSourceMapping("my-function")

	// Verify all properties are set
	if esm.FunctionName != "my-function" {
		t.Errorf("unexpected FunctionName: %v", esm.FunctionName)
	}
	if esm.DocumentDBEventSourceConfig.DatabaseName != "my-database" {
		t.Errorf("unexpected DatabaseName: %s", esm.DocumentDBEventSourceConfig.DatabaseName)
	}
	if esm.DocumentDBEventSourceConfig.CollectionName != "my-collection" {
		t.Errorf("unexpected CollectionName: %s", esm.DocumentDBEventSourceConfig.CollectionName)
	}
	if esm.DocumentDBEventSourceConfig.FullDocument != "UpdateLookup" {
		t.Errorf("unexpected FullDocument: %s", esm.DocumentDBEventSourceConfig.FullDocument)
	}
	if esm.StartingPosition != "TRIM_HORIZON" {
		t.Errorf("unexpected StartingPosition: %s", esm.StartingPosition)
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
	if len(esm.SourceAccessConfigurations) != 1 {
		t.Errorf("unexpected SourceAccessConfigurations length: %d", len(esm.SourceAccessConfigurations))
	}
	if len(esm.FilterCriteria.Filters) != 1 {
		t.Errorf("unexpected FilterCriteria.Filters length: %d", len(esm.FilterCriteria.Filters))
	}
}

func TestDocumentDBEventProperties_ToEventSourceMapping_ToCloudFormation(t *testing.T) {
	batchSize := 50
	props := &DocumentDBEventProperties{
		Cluster:        "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster",
		DatabaseName:   "my-database",
		CollectionName: "my-collection",
		FullDocument:   "UpdateLookup",
		BatchSize:      &batchSize,
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
	if cfnProps["EventSourceArn"] != "arn:aws:rds:us-east-1:123456789012:cluster:my-docdb-cluster" {
		t.Errorf("expected EventSourceArn in CFN properties")
	}

	docDBConfig := cfnProps["DocumentDBEventSourceConfig"].(map[string]interface{})
	if docDBConfig["DatabaseName"] != "my-database" {
		t.Errorf("unexpected DatabaseName in CFN properties: %v", docDBConfig["DatabaseName"])
	}
	if docDBConfig["CollectionName"] != "my-collection" {
		t.Errorf("unexpected CollectionName in CFN properties: %v", docDBConfig["CollectionName"])
	}
	if docDBConfig["FullDocument"] != "UpdateLookup" {
		t.Errorf("unexpected FullDocument in CFN properties: %v", docDBConfig["FullDocument"])
	}
}
