package s3

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBucket_JSONSerialization(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
		Tags: []Tag{
			{Key: "Environment", Value: "Production"},
		},
		VersioningConfiguration: &VersioningConfiguration{
			Status: "Enabled",
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket to JSON: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket from JSON: %v", err)
	}

	if unmarshaled.BucketName != bucket.BucketName {
		t.Errorf("BucketName mismatch: got %v, want %v", unmarshaled.BucketName, bucket.BucketName)
	}
}

func TestBucket_YAMLSerialization(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
	}

	data, err := yaml.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket to YAML: %v", err)
	}

	var unmarshaled Bucket
	if err := yaml.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket from YAML: %v", err)
	}

	if unmarshaled.BucketName != bucket.BucketName {
		t.Errorf("BucketName mismatch: got %v, want %v", unmarshaled.BucketName, bucket.BucketName)
	}
}

func TestBucket_WithIntrinsicFunctions(t *testing.T) {
	bucket := Bucket{
		BucketName: map[string]interface{}{
			"Fn::Sub": "${AWS::StackName}-bucket",
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with intrinsics: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with intrinsics: %v", err)
	}
}

func TestBucket_WithNotificationConfiguration(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
		NotificationConfiguration: &NotificationConfiguration{
			LambdaConfigurations: []LambdaConfiguration{
				{
					Event:    "s3:ObjectCreated:*",
					Function: "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
					Filter: &NotificationFilter{
						S3Key: &S3KeyFilter{
							Rules: []FilterRule{
								{Name: "prefix", Value: "uploads/"},
								{Name: "suffix", Value: ".json"},
							},
						},
					},
				},
			},
			QueueConfigurations: []QueueConfiguration{
				{
					Event: "s3:ObjectRemoved:*",
					Queue: "arn:aws:sqs:us-east-1:123456789012:MyQueue",
				},
			},
			TopicConfigurations: []TopicConfiguration{
				{
					Event: "s3:ObjectCreated:Put",
					Topic: "arn:aws:sns:us-east-1:123456789012:MyTopic",
				},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with notifications: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with notifications: %v", err)
	}

	if unmarshaled.NotificationConfiguration == nil {
		t.Error("NotificationConfiguration should not be nil")
	}

	if len(unmarshaled.NotificationConfiguration.LambdaConfigurations) != 1 {
		t.Errorf("LambdaConfigurations length mismatch: got %d, want 1",
			len(unmarshaled.NotificationConfiguration.LambdaConfigurations))
	}
}

func TestBucket_WithEncryption(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-encrypted-bucket",
		BucketEncryption: &BucketEncryption{
			ServerSideEncryptionConfiguration: []ServerSideEncryptionRule{
				{
					BucketKeyEnabled: true,
					ServerSideEncryptionByDefault: &ServerSideEncryptionByDefault{
						SSEAlgorithm:   "aws:kms",
						KMSMasterKeyID: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
					},
				},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal encrypted bucket: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal encrypted bucket: %v", err)
	}

	if unmarshaled.BucketEncryption == nil {
		t.Error("BucketEncryption should not be nil")
	}
}

func TestBucket_WithLifecycleConfiguration(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					Id:     "MoveToGlacier",
					Status: "Enabled",
					Filter: &LifecycleRuleFilter{
						Prefix: "archive/",
					},
					Transitions: []Transition{
						{
							StorageClass:     "GLACIER",
							TransitionInDays: 90,
						},
					},
					ExpirationInDays: 365,
				},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with lifecycle: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with lifecycle: %v", err)
	}

	if unmarshaled.LifecycleConfiguration == nil {
		t.Error("LifecycleConfiguration should not be nil")
	}
}

func TestBucket_WithCorsConfiguration(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
		CorsConfiguration: &CorsConfiguration{
			CorsRules: []CorsRule{
				{
					AllowedHeaders: []interface{}{"*"},
					AllowedMethods: []string{"GET", "PUT", "POST"},
					AllowedOrigins: []interface{}{"https://example.com"},
					ExposedHeaders: []interface{}{"ETag"},
					MaxAge:         3600,
				},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with CORS: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with CORS: %v", err)
	}

	if unmarshaled.CorsConfiguration == nil {
		t.Error("CorsConfiguration should not be nil")
	}
}

func TestBucket_WithPublicAccessBlock(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-secure-bucket",
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAcls:       true,
			BlockPublicPolicy:     true,
			IgnorePublicAcls:      true,
			RestrictPublicBuckets: true,
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with public access block: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with public access block: %v", err)
	}

	if unmarshaled.PublicAccessBlockConfiguration == nil {
		t.Error("PublicAccessBlockConfiguration should not be nil")
	}
}

func TestBucket_WithWebsiteConfiguration(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-website-bucket",
		WebsiteConfiguration: &WebsiteConfiguration{
			IndexDocument: &IndexDocument{Suffix: "index.html"},
			ErrorDocument: &ErrorDocument{Key: "error.html"},
			RoutingRules: []RoutingRule{
				{
					Condition: &RoutingRuleCondition{
						KeyPrefixEquals: "docs/",
					},
					Redirect: Redirect{
						ReplaceKeyPrefixWith: "documents/",
					},
				},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with website config: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with website config: %v", err)
	}

	if unmarshaled.WebsiteConfiguration == nil {
		t.Error("WebsiteConfiguration should not be nil")
	}
}

func TestBucket_WithReplicationConfiguration(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-source-bucket",
		ReplicationConfiguration: &ReplicationConfiguration{
			Role: "arn:aws:iam::123456789012:role/ReplicationRole",
			Rules: []ReplicationRule{
				{
					Id:       "ReplicateAll",
					Status:   "Enabled",
					Priority: 1,
					Destination: ReplicationDestination{
						Bucket:       "arn:aws:s3:::my-destination-bucket",
						StorageClass: "STANDARD",
					},
					Filter: &ReplicationRuleFilter{
						Prefix: "",
					},
				},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with replication: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with replication: %v", err)
	}

	if unmarshaled.ReplicationConfiguration == nil {
		t.Error("ReplicationConfiguration should not be nil")
	}
}

func TestBucket_WithOwnershipControls(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
		OwnershipControls: &OwnershipControls{
			Rules: []OwnershipControlsRule{
				{ObjectOwnership: "BucketOwnerEnforced"},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with ownership controls: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with ownership controls: %v", err)
	}

	if unmarshaled.OwnershipControls == nil {
		t.Error("OwnershipControls should not be nil")
	}
}

func TestBucket_WithObjectLock(t *testing.T) {
	bucket := Bucket{
		BucketName:        "my-locked-bucket",
		ObjectLockEnabled: true,
		ObjectLockConfiguration: &ObjectLockConfiguration{
			ObjectLockEnabled: "Enabled",
			Rule: &ObjectLockRule{
				DefaultRetention: &DefaultRetention{
					Mode: "GOVERNANCE",
					Days: 30,
				},
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with object lock: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with object lock: %v", err)
	}

	if unmarshaled.ObjectLockConfiguration == nil {
		t.Error("ObjectLockConfiguration should not be nil")
	}
}

func TestBucket_OmitEmpty(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal to raw map: %v", err)
	}

	omittedFields := []string{
		"AccelerateConfiguration",
		"AccessControl",
		"BucketEncryption",
		"CorsConfiguration",
		"LifecycleConfiguration",
		"LoggingConfiguration",
		"NotificationConfiguration",
		"PublicAccessBlockConfiguration",
		"ReplicationConfiguration",
		"Tags",
		"VersioningConfiguration",
		"WebsiteConfiguration",
	}

	for _, field := range omittedFields {
		if _, exists := raw[field]; exists {
			t.Errorf("Expected field %q to be omitted when empty", field)
		}
	}
}

func TestNotificationFilter_JSONSerialization(t *testing.T) {
	filter := NotificationFilter{
		S3Key: &S3KeyFilter{
			Rules: []FilterRule{
				{Name: "prefix", Value: "images/"},
				{Name: "suffix", Value: ".jpg"},
			},
		},
	}

	data, err := json.Marshal(filter)
	if err != nil {
		t.Fatalf("Failed to marshal notification filter: %v", err)
	}

	var unmarshaled NotificationFilter
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal notification filter: %v", err)
	}

	if len(unmarshaled.S3Key.Rules) != 2 {
		t.Errorf("Filter rules length mismatch: got %d, want 2", len(unmarshaled.S3Key.Rules))
	}
}

func TestLambdaConfiguration_JSONSerialization(t *testing.T) {
	config := LambdaConfiguration{
		Event: "s3:ObjectCreated:*",
		Function: map[string]interface{}{
			"Fn::GetAtt": []string{"MyFunction", "Arn"},
		},
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal Lambda configuration: %v", err)
	}

	var unmarshaled LambdaConfiguration
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal Lambda configuration: %v", err)
	}

	if unmarshaled.Event != config.Event {
		t.Errorf("Event mismatch: got %v, want %v", unmarshaled.Event, config.Event)
	}
}

func TestBucket_WithEventBridgeNotification(t *testing.T) {
	bucket := Bucket{
		BucketName: "my-bucket",
		NotificationConfiguration: &NotificationConfiguration{
			EventBridgeConfiguration: &EventBridgeConfiguration{
				EventBridgeEnabled: true,
			},
		},
	}

	data, err := json.Marshal(bucket)
	if err != nil {
		t.Fatalf("Failed to marshal bucket with EventBridge: %v", err)
	}

	var unmarshaled Bucket
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal bucket with EventBridge: %v", err)
	}

	if unmarshaled.NotificationConfiguration.EventBridgeConfiguration == nil {
		t.Error("EventBridgeConfiguration should not be nil")
	}
}
