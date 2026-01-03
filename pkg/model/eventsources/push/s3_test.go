package push

import (
	"testing"
)

func TestNewS3EventSource(t *testing.T) {
	bucket := "my-bucket"
	events := []string{"s3:ObjectCreated:*", "s3:ObjectRemoved:*"}

	source := NewS3EventSource(bucket, events)

	if source == nil {
		t.Fatal("NewS3EventSource returned nil")
	}
	if source.Bucket != bucket {
		t.Errorf("expected bucket %v, got %v", bucket, source.Bucket)
	}
	if len(source.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(source.Events))
	}
	if source.Events[0] != "s3:ObjectCreated:*" {
		t.Errorf("expected first event s3:ObjectCreated:*, got %s", source.Events[0])
	}
}

func TestS3EventSource_WithPrefixFilter(t *testing.T) {
	bucket := "my-bucket"
	events := []string{"s3:ObjectCreated:*"}

	source := NewS3EventSource(bucket, events).
		WithPrefixFilter("uploads/")

	if source.Filter == nil {
		t.Fatal("expected filter to be set")
	}
	if source.Filter.S3Key == nil {
		t.Fatal("expected S3Key filter to be set")
	}
	if len(source.Filter.S3Key.Rules) != 1 {
		t.Fatalf("expected 1 filter rule, got %d", len(source.Filter.S3Key.Rules))
	}
	if source.Filter.S3Key.Rules[0].Name != "prefix" {
		t.Errorf("expected rule name 'prefix', got %s", source.Filter.S3Key.Rules[0].Name)
	}
	if source.Filter.S3Key.Rules[0].Value != "uploads/" {
		t.Errorf("expected rule value 'uploads/', got %v", source.Filter.S3Key.Rules[0].Value)
	}
}

func TestS3EventSource_WithSuffixFilter(t *testing.T) {
	bucket := "my-bucket"
	events := []string{"s3:ObjectCreated:*"}

	source := NewS3EventSource(bucket, events).
		WithSuffixFilter(".jpg")

	if source.Filter == nil {
		t.Fatal("expected filter to be set")
	}
	if source.Filter.S3Key == nil {
		t.Fatal("expected S3Key filter to be set")
	}
	if len(source.Filter.S3Key.Rules) != 1 {
		t.Fatalf("expected 1 filter rule, got %d", len(source.Filter.S3Key.Rules))
	}
	if source.Filter.S3Key.Rules[0].Name != "suffix" {
		t.Errorf("expected rule name 'suffix', got %s", source.Filter.S3Key.Rules[0].Name)
	}
	if source.Filter.S3Key.Rules[0].Value != ".jpg" {
		t.Errorf("expected rule value '.jpg', got %v", source.Filter.S3Key.Rules[0].Value)
	}
}

func TestS3EventSource_WithPrefixAndSuffixFilter(t *testing.T) {
	bucket := "my-bucket"
	events := []string{"s3:ObjectCreated:Put"}

	source := NewS3EventSource(bucket, events).
		WithPrefixFilter("images/").
		WithSuffixFilter(".png")

	if source.Filter == nil {
		t.Fatal("expected filter to be set")
	}
	if source.Filter.S3Key == nil {
		t.Fatal("expected S3Key filter to be set")
	}
	if len(source.Filter.S3Key.Rules) != 2 {
		t.Fatalf("expected 2 filter rules, got %d", len(source.Filter.S3Key.Rules))
	}

	// Check prefix rule
	if source.Filter.S3Key.Rules[0].Name != "prefix" {
		t.Errorf("expected first rule name 'prefix', got %s", source.Filter.S3Key.Rules[0].Name)
	}
	if source.Filter.S3Key.Rules[0].Value != "images/" {
		t.Errorf("expected first rule value 'images/', got %v", source.Filter.S3Key.Rules[0].Value)
	}

	// Check suffix rule
	if source.Filter.S3Key.Rules[1].Name != "suffix" {
		t.Errorf("expected second rule name 'suffix', got %s", source.Filter.S3Key.Rules[1].Name)
	}
	if source.Filter.S3Key.Rules[1].Value != ".png" {
		t.Errorf("expected second rule value '.png', got %v", source.Filter.S3Key.Rules[1].Value)
	}
}

func TestS3EventSource_ToCloudFormationResources(t *testing.T) {
	bucket := map[string]interface{}{"Ref": "MyBucket"}
	events := []string{"s3:ObjectCreated:*"}
	functionRef := map[string]interface{}{"Ref": "MyFunction"}
	functionName := "MyFunction"

	source := NewS3EventSource(bucket, events)

	resources, notification, err := source.ToCloudFormationResources(functionRef, functionName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that permission resource was created
	if len(resources) == 0 {
		t.Fatal("expected at least one resource")
	}

	permissionID := "MyFunctionS3Permission"
	if _, ok := resources[permissionID]; !ok {
		t.Errorf("expected resource %s to be created", permissionID)
	}

	// Check permission properties
	permission := resources[permissionID].(map[string]interface{})
	if permission["Type"] != "AWS::Lambda::Permission" {
		t.Errorf("expected Type AWS::Lambda::Permission, got %v", permission["Type"])
	}

	properties := permission["Properties"].(map[string]interface{})
	if properties["Action"] != "lambda:InvokeFunction" {
		t.Errorf("expected Action lambda:InvokeFunction, got %v", properties["Action"])
	}
	if properties["Principal"] != "s3.amazonaws.com" {
		t.Errorf("expected Principal s3.amazonaws.com, got %v", properties["Principal"])
	}

	// Check notification configuration
	if notification == nil {
		t.Fatal("expected notification configuration")
	}
	if notification.Event != "s3:ObjectCreated:*" {
		t.Errorf("expected Event s3:ObjectCreated:*, got %s", notification.Event)
	}
	if notification.Function == nil {
		t.Error("expected Function to be set")
	}
}

func TestS3EventSource_ToCloudFormationResourcesWithFilter(t *testing.T) {
	bucket := "my-bucket"
	events := []string{"s3:ObjectCreated:Put"}
	functionRef := map[string]interface{}{"GetAtt": []interface{}{"MyFunction", "Arn"}}
	functionName := "MyFunction"

	source := NewS3EventSource(bucket, events).
		WithPrefixFilter("uploads/").
		WithSuffixFilter(".jpg")

	resources, notification, err := source.ToCloudFormationResources(functionRef, functionName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resources) == 0 {
		t.Fatal("expected at least one resource")
	}

	// Check notification has filter
	if notification == nil {
		t.Fatal("expected notification configuration")
	}
	if notification.Filter == nil {
		t.Fatal("expected notification filter to be set")
	}
	if notification.Filter.S3Key == nil {
		t.Fatal("expected S3Key filter to be set")
	}
	if len(notification.Filter.S3Key.Rules) != 2 {
		t.Fatalf("expected 2 filter rules, got %d", len(notification.Filter.S3Key.Rules))
	}

	// Verify filter rules
	if notification.Filter.S3Key.Rules[0].Name != "prefix" {
		t.Errorf("expected first rule name 'prefix', got %s", notification.Filter.S3Key.Rules[0].Name)
	}
	if notification.Filter.S3Key.Rules[0].Value != "uploads/" {
		t.Errorf("expected first rule value 'uploads/', got %v", notification.Filter.S3Key.Rules[0].Value)
	}
	if notification.Filter.S3Key.Rules[1].Name != "suffix" {
		t.Errorf("expected second rule name 'suffix', got %s", notification.Filter.S3Key.Rules[1].Name)
	}
	if notification.Filter.S3Key.Rules[1].Value != ".jpg" {
		t.Errorf("expected second rule value '.jpg', got %v", notification.Filter.S3Key.Rules[1].Value)
	}
}

func TestS3EventSource_ToCloudFormationResourcesMultiEvent(t *testing.T) {
	bucket := "my-bucket"
	events := []string{"s3:ObjectCreated:*", "s3:ObjectRemoved:*", "s3:ObjectRestore:*"}
	functionRef := map[string]interface{}{"Ref": "MyFunction"}
	functionName := "MyFunction"

	source := NewS3EventSource(bucket, events)

	resources, notifications, err := source.ToCloudFormationResourcesMultiEvent(functionRef, functionName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that permission resource was created
	if len(resources) == 0 {
		t.Fatal("expected at least one resource")
	}

	// Check notifications for each event
	if len(notifications) != 3 {
		t.Fatalf("expected 3 notifications, got %d", len(notifications))
	}

	expectedEvents := []string{"s3:ObjectCreated:*", "s3:ObjectRemoved:*", "s3:ObjectRestore:*"}
	for i, notification := range notifications {
		if notification.Event != expectedEvents[i] {
			t.Errorf("notification %d: expected Event %s, got %s", i, expectedEvents[i], notification.Event)
		}
		if notification.Function == nil {
			t.Errorf("notification %d: expected Function to be set", i)
		}
	}
}

func TestS3EventSource_ToCloudFormationResourcesMultiEventWithFilter(t *testing.T) {
	bucket := map[string]interface{}{"Fn::GetAtt": []interface{}{"MyBucket", "Arn"}}
	events := []string{"s3:ObjectCreated:Put", "s3:ObjectCreated:Post"}
	functionRef := map[string]interface{}{"Ref": "MyFunction"}
	functionName := "MyFunction"

	source := NewS3EventSource(bucket, events).
		WithPrefixFilter("data/")

	resources, notifications, err := source.ToCloudFormationResourcesMultiEvent(functionRef, functionName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resources) == 0 {
		t.Fatal("expected at least one resource")
	}

	if len(notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(notifications))
	}

	// Both notifications should have the same filter
	for i, notification := range notifications {
		if notification.Filter == nil {
			t.Errorf("notification %d: expected filter to be set", i)
			continue
		}
		if notification.Filter.S3Key == nil {
			t.Errorf("notification %d: expected S3Key filter to be set", i)
			continue
		}
		if len(notification.Filter.S3Key.Rules) != 1 {
			t.Errorf("notification %d: expected 1 filter rule, got %d", i, len(notification.Filter.S3Key.Rules))
			continue
		}
		if notification.Filter.S3Key.Rules[0].Name != "prefix" {
			t.Errorf("notification %d: expected rule name 'prefix', got %s", i, notification.Filter.S3Key.Rules[0].Name)
		}
		if notification.Filter.S3Key.Rules[0].Value != "data/" {
			t.Errorf("notification %d: expected rule value 'data/', got %v", i, notification.Filter.S3Key.Rules[0].Value)
		}
	}
}

func TestS3EventSource_Validate(t *testing.T) {
	tests := []struct {
		name      string
		source    *S3EventSource
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid source",
			source: &S3EventSource{
				Bucket: "my-bucket",
				Events: []string{"s3:ObjectCreated:*"},
			},
			expectErr: false,
		},
		{
			name: "missing bucket",
			source: &S3EventSource{
				Events: []string{"s3:ObjectCreated:*"},
			},
			expectErr: true,
			errMsg:    "bucket is required",
		},
		{
			name: "missing events",
			source: &S3EventSource{
				Bucket: "my-bucket",
				Events: []string{},
			},
			expectErr: true,
			errMsg:    "must specify at least one event",
		},
		{
			name: "empty event name",
			source: &S3EventSource{
				Bucket: "my-bucket",
				Events: []string{"s3:ObjectCreated:*", ""},
			},
			expectErr: true,
			errMsg:    "event name cannot be empty",
		},
		{
			name: "invalid filter rule name",
			source: &S3EventSource{
				Bucket: "my-bucket",
				Events: []string{"s3:ObjectCreated:*"},
				Filter: &S3NotificationFilter{
					S3Key: &S3KeyFilter{
						Rules: []S3FilterRule{
							{Name: "invalid", Value: "test"},
						},
					},
				},
			},
			expectErr: true,
			errMsg:    "invalid S3 filter rule name",
		},
		{
			name: "nil filter rule value",
			source: &S3EventSource{
				Bucket: "my-bucket",
				Events: []string{"s3:ObjectCreated:*"},
				Filter: &S3NotificationFilter{
					S3Key: &S3KeyFilter{
						Rules: []S3FilterRule{
							{Name: "prefix", Value: nil},
						},
					},
				},
			},
			expectErr: true,
			errMsg:    "filter rule value cannot be nil",
		},
		{
			name: "valid with prefix filter",
			source: &S3EventSource{
				Bucket: "my-bucket",
				Events: []string{"s3:ObjectCreated:Put"},
				Filter: &S3NotificationFilter{
					S3Key: &S3KeyFilter{
						Rules: []S3FilterRule{
							{Name: "prefix", Value: "uploads/"},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "valid with suffix filter",
			source: &S3EventSource{
				Bucket: "my-bucket",
				Events: []string{"s3:ObjectRemoved:*"},
				Filter: &S3NotificationFilter{
					S3Key: &S3KeyFilter{
						Rules: []S3FilterRule{
							{Name: "suffix", Value: ".log"},
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.Validate()
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errMsg)
				} else if tt.errMsg != "" {
					// Just check if error message contains the expected substring
					// (not exact match to allow for variations)
					errStr := err.Error()
					found := false
					if len(tt.errMsg) > 0 {
						for i := 0; i <= len(errStr)-len(tt.errMsg); i++ {
							if errStr[i:i+len(tt.errMsg)] == tt.errMsg {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, errStr)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestS3EventSource_Getters(t *testing.T) {
	bucket := "my-bucket"
	events := []string{"s3:ObjectCreated:*"}
	filter := &S3NotificationFilter{
		S3Key: &S3KeyFilter{
			Rules: []S3FilterRule{
				{Name: "prefix", Value: "uploads/"},
			},
		},
	}

	source := &S3EventSource{
		Bucket: bucket,
		Events: events,
		Filter: filter,
	}

	if source.GetBucket() != bucket {
		t.Errorf("GetBucket() = %v, want %v", source.GetBucket(), bucket)
	}

	returnedEvents := source.GetEvents()
	if len(returnedEvents) != len(events) {
		t.Errorf("GetEvents() length = %d, want %d", len(returnedEvents), len(events))
	}
	if len(returnedEvents) > 0 && returnedEvents[0] != events[0] {
		t.Errorf("GetEvents()[0] = %s, want %s", returnedEvents[0], events[0])
	}

	returnedFilter := source.GetFilter()
	if returnedFilter != filter {
		t.Errorf("GetFilter() = %v, want %v", returnedFilter, filter)
	}
}

func TestS3EventSource_IntrinsicFunctions(t *testing.T) {
	// Test with CloudFormation intrinsic functions
	bucket := map[string]interface{}{
		"Fn::GetAtt": []interface{}{"MyBucket", "Arn"},
	}
	events := []string{"s3:ObjectCreated:*"}
	functionRef := map[string]interface{}{
		"Ref": "MyFunction",
	}
	functionName := "MyFunction"

	source := NewS3EventSource(bucket, events).
		WithPrefixFilter(map[string]interface{}{
			"Ref": "UploadPrefix",
		})

	resources, notification, err := source.ToCloudFormationResources(functionRef, functionName)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resources) == 0 {
		t.Fatal("expected at least one resource")
	}

	// Verify permission has the intrinsic function for SourceArn
	permissionID := "MyFunctionS3Permission"
	permission := resources[permissionID].(map[string]interface{})
	properties := permission["Properties"].(map[string]interface{})

	sourceArn := properties["SourceArn"].(map[string]interface{})
	if sourceArn["Fn::GetAtt"] == nil {
		t.Error("expected SourceArn to contain Fn::GetAtt intrinsic function")
	}

	// Verify notification filter has intrinsic function
	if notification.Filter == nil || notification.Filter.S3Key == nil {
		t.Fatal("expected notification filter to be set")
	}
	if len(notification.Filter.S3Key.Rules) == 0 {
		t.Fatal("expected at least one filter rule")
	}

	filterValue, ok := notification.Filter.S3Key.Rules[0].Value.(map[string]interface{})
	if !ok {
		t.Fatal("expected filter value to be a map (intrinsic function)")
	}
	if filterValue["Ref"] != "UploadPrefix" {
		t.Errorf("expected filter value Ref to be 'UploadPrefix', got %v", filterValue["Ref"])
	}
}

func TestS3EventSource_VariousEventTypes(t *testing.T) {
	eventTypes := []string{
		"s3:ObjectCreated:*",
		"s3:ObjectCreated:Put",
		"s3:ObjectCreated:Post",
		"s3:ObjectCreated:Copy",
		"s3:ObjectCreated:CompleteMultipartUpload",
		"s3:ObjectRemoved:*",
		"s3:ObjectRemoved:Delete",
		"s3:ObjectRemoved:DeleteMarkerCreated",
		"s3:ObjectRestore:Post",
		"s3:ObjectRestore:Completed",
		"s3:ReducedRedundancyLostObject",
		"s3:Replication:OperationFailedReplication",
		"s3:LifecycleExpiration:Delete",
		"s3:LifecycleTransition",
		"s3:IntelligentTiering",
		"s3:ObjectTagging:Put",
		"s3:ObjectAcl:Put",
	}

	bucket := "test-bucket"
	functionRef := map[string]interface{}{"Ref": "TestFunction"}
	functionName := "TestFunction"

	for _, eventType := range eventTypes {
		t.Run(eventType, func(t *testing.T) {
			source := NewS3EventSource(bucket, []string{eventType})

			resources, notification, err := source.ToCloudFormationResources(functionRef, functionName)

			if err != nil {
				t.Fatalf("unexpected error for event type %s: %v", eventType, err)
			}

			if len(resources) == 0 {
				t.Fatalf("expected resources for event type %s", eventType)
			}

			if notification == nil {
				t.Fatalf("expected notification for event type %s", eventType)
			}

			if notification.Event != eventType {
				t.Errorf("expected notification event %s, got %s", eventType, notification.Event)
			}
		})
	}
}
