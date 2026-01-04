package sam

import (
	"testing"
)

func TestApplicationTransformer_BasicApplication(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	if resources == nil {
		t.Fatal("resources should not be nil")
	}

	// Should create CloudFormation Stack
	stackResource, ok := resources["MyApp"].(map[string]interface{})
	if !ok {
		t.Fatal("should have CloudFormation Stack resource")
	}
	if stackResource["Type"] != TypeCloudFormationStack {
		t.Errorf("expected Type %q, got %v", TypeCloudFormationStack, stackResource["Type"])
	}

	props := stackResource["Properties"].(map[string]interface{})
	if props["TemplateURL"] == nil {
		t.Error("TemplateURL should not be nil")
	}
}

func TestApplicationTransformer_WithSemanticVersion(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: map[string]interface{}{
			"ApplicationId":   "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
			"SemanticVersion": "1.0.0",
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})
	templateURL := props["TemplateURL"].(map[string]interface{})

	// Should have Fn::Transform with parameters
	transform, ok := templateURL["Fn::Transform"].(map[string]interface{})
	if !ok {
		t.Fatal("TemplateURL should contain Fn::Transform")
	}

	params := transform["Parameters"].(map[string]interface{})
	if params["ApplicationId"] != "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app" {
		t.Errorf("expected ApplicationId, got %v", params["ApplicationId"])
	}
	if params["SemanticVersion"] != "1.0.0" {
		t.Errorf("expected SemanticVersion '1.0.0', got %v", params["SemanticVersion"])
	}
}

func TestApplicationTransformer_WithParameters(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
		Parameters: map[string]interface{}{
			"TableName": "MyTable",
			"BucketArn": map[string]interface{}{"Fn::GetAtt": []string{"MyBucket", "Arn"}},
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	params := props["Parameters"].(map[string]interface{})
	if params["TableName"] != "MyTable" {
		t.Errorf("expected TableName 'MyTable', got %v", params["TableName"])
	}
	if params["BucketArn"] == nil {
		t.Error("BucketArn should not be nil")
	}
}

func TestApplicationTransformer_WithNotificationArns(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
		NotificationArns: []interface{}{
			"arn:aws:sns:us-east-1:123456789012:my-topic",
			map[string]interface{}{"Ref": "NotificationTopic"},
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	notificationArns := props["NotificationARNs"].([]interface{})
	if len(notificationArns) != 2 {
		t.Errorf("expected 2 notification ARNs, got %d", len(notificationArns))
	}
}

func TestApplicationTransformer_WithTags(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
		Tags: map[string]string{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	tags := props["Tags"].([]interface{})
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	// Verify tag structure
	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})
		if tagMap["Key"] == nil || tagMap["Value"] == nil {
			t.Error("each tag should have Key and Value")
		}
	}
}

func TestApplicationTransformer_WithTimeoutInMinutes(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location:         "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
		TimeoutInMinutes: 30,
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	if props["TimeoutInMinutes"] != 30 {
		t.Errorf("expected TimeoutInMinutes 30, got %v", props["TimeoutInMinutes"])
	}
}

func TestApplicationTransformer_WithCondition(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location:  "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
		Condition: "IsProduction",
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	if stackResource["Condition"] != "IsProduction" {
		t.Errorf("expected Condition 'IsProduction', got %v", stackResource["Condition"])
	}
}

func TestApplicationTransformer_WithDependsOn(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location:  "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
		DependsOn: []string{"MyBucket", "MyTable"},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	dependsOn := stackResource["DependsOn"].([]string)
	if len(dependsOn) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(dependsOn))
	}
}

func TestApplicationTransformer_WithS3Location(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: map[string]interface{}{
			"Bucket": "my-bucket",
			"Key":    "templates/my-template.yaml",
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	templateURL, ok := props["TemplateURL"].(string)
	if !ok {
		t.Fatal("TemplateURL should be a string for S3 location")
	}
	if templateURL != "https://my-bucket.s3.amazonaws.com/templates/my-template.yaml" {
		t.Errorf("expected S3 URL, got %v", templateURL)
	}
}

func TestApplicationTransformer_WithS3LocationAndVersion(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: map[string]interface{}{
			"Bucket":  "my-bucket",
			"Key":     "templates/my-template.yaml",
			"Version": "abc123",
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	templateURL, ok := props["TemplateURL"].(string)
	if !ok {
		t.Fatal("TemplateURL should be a string for S3 location")
	}
	if templateURL != "https://my-bucket.s3.amazonaws.com/templates/my-template.yaml?versionId=abc123" {
		t.Errorf("expected S3 URL with version, got %v", templateURL)
	}
}

func TestApplicationTransformer_WithHTTPSLocation(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: "https://s3.amazonaws.com/my-bucket/templates/my-template.yaml",
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	templateURL := props["TemplateURL"].(string)
	if templateURL != "https://s3.amazonaws.com/my-bucket/templates/my-template.yaml" {
		t.Errorf("expected HTTPS URL, got %v", templateURL)
	}
}

func TestApplicationTransformer_WithIntrinsicS3Location(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: map[string]interface{}{
			"Bucket": map[string]interface{}{"Ref": "TemplateBucket"},
			"Key":    map[string]interface{}{"Ref": "TemplateKey"},
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	props := stackResource["Properties"].(map[string]interface{})

	templateURL, ok := props["TemplateURL"].(map[string]interface{})
	if !ok {
		t.Fatal("TemplateURL should be a map for intrinsic S3 location")
	}

	fnSub, ok := templateURL["Fn::Sub"].([]interface{})
	if !ok {
		t.Fatal("TemplateURL should contain Fn::Sub")
	}
	if len(fnSub) != 2 {
		t.Errorf("Fn::Sub should have 2 elements, got %d", len(fnSub))
	}
}

func TestApplicationTransformer_NilLocation(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{}

	_, err := transformer.Transform("MyApp", app, nil)
	if err == nil {
		t.Error("expected error for nil Location")
	}
}

func TestIsSARApplicationID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app", true},
		{"arn:aws:serverlessrepo:eu-west-1:999999999999:applications/another-app", true},
		{"arn:aws:s3:::my-bucket", false},
		{"https://s3.amazonaws.com/bucket/key", false},
		{"", false},
	}

	for _, tc := range tests {
		result := isSARApplicationID(tc.input)
		if result != tc.expected {
			t.Errorf("isSARApplicationID(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestApplicationTransformer_WithMetadata(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
		Metadata: map[string]interface{}{
			"cfn-lint": map[string]interface{}{
				"config": map[string]interface{}{
					"ignore_checks": []string{"W3002"},
				},
			},
		},
	}

	resources, err := transformer.Transform("MyApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyApp"].(map[string]interface{})
	metadata := stackResource["Metadata"].(map[string]interface{})
	if metadata["cfn-lint"] == nil {
		t.Error("Metadata should contain cfn-lint")
	}
}

func TestApplicationTransformer_AllProperties(t *testing.T) {
	transformer := NewApplicationTransformer()

	app := &Application{
		Location: map[string]interface{}{
			"ApplicationId":   "arn:aws:serverlessrepo:us-east-1:123456789012:applications/my-app",
			"SemanticVersion": "2.0.0",
		},
		Parameters: map[string]interface{}{
			"TableName": "MyTable",
		},
		NotificationArns: []interface{}{
			"arn:aws:sns:us-east-1:123456789012:my-topic",
		},
		Tags: map[string]string{
			"Environment": "test",
		},
		TimeoutInMinutes: 60,
		Condition:        "CreateApp",
		DependsOn:        "MyBucket",
		Metadata: map[string]interface{}{
			"Custom": "metadata",
		},
	}

	resources, err := transformer.Transform("MyCompleteApp", app, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	stackResource := resources["MyCompleteApp"].(map[string]interface{})

	// Check resource type
	if stackResource["Type"] != TypeCloudFormationStack {
		t.Errorf("expected Type %q, got %v", TypeCloudFormationStack, stackResource["Type"])
	}

	// Check optional attributes
	if stackResource["Condition"] != "CreateApp" {
		t.Errorf("expected Condition 'CreateApp', got %v", stackResource["Condition"])
	}
	if stackResource["DependsOn"] != "MyBucket" {
		t.Errorf("expected DependsOn 'MyBucket', got %v", stackResource["DependsOn"])
	}
	if stackResource["Metadata"] == nil {
		t.Error("Metadata should not be nil")
	}

	// Check properties
	props := stackResource["Properties"].(map[string]interface{})
	if props["TemplateURL"] == nil {
		t.Error("TemplateURL should not be nil")
	}
	if props["Parameters"] == nil {
		t.Error("Parameters should not be nil")
	}
	if props["NotificationARNs"] == nil {
		t.Error("NotificationARNs should not be nil")
	}
	if props["Tags"] == nil {
		t.Error("Tags should not be nil")
	}
	if props["TimeoutInMinutes"] != 60 {
		t.Errorf("expected TimeoutInMinutes 60, got %v", props["TimeoutInMinutes"])
	}
}
