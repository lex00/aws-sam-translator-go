package sam

import (
	"encoding/json"
	"testing"
)

func TestFunctionTransformer_BasicFunction(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:    "index.handler",
		Runtime:    "nodejs18.x",
		CodeUri:    "s3://bucket/code.zip",
		Timeout:    30,
		MemorySize: 256,
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	if resources == nil {
		t.Fatal("resources should not be nil")
	}

	// Should create Lambda function
	fnResource, ok := resources["MyFunction"].(map[string]interface{})
	if !ok {
		t.Fatal("should have Lambda function resource")
	}
	if fnResource["Type"] != "AWS::Lambda::Function" {
		t.Errorf("expected Type 'AWS::Lambda::Function', got %v", fnResource["Type"])
	}

	props := fnResource["Properties"].(map[string]interface{})
	if props["Handler"] != "index.handler" {
		t.Errorf("expected Handler 'index.handler', got %v", props["Handler"])
	}
	if props["Runtime"] != "nodejs18.x" {
		t.Errorf("expected Runtime 'nodejs18.x', got %v", props["Runtime"])
	}
	if props["Timeout"] != 30 {
		t.Errorf("expected Timeout 30, got %v", props["Timeout"])
	}
	if props["MemorySize"] != 256 {
		t.Errorf("expected MemorySize 256, got %v", props["MemorySize"])
	}

	// Should create IAM role (auto-generated)
	roleResource, ok := resources["MyFunctionRole"].(map[string]interface{})
	if !ok {
		t.Fatal("should have IAM role resource")
	}
	if roleResource["Type"] != "AWS::IAM::Role" {
		t.Errorf("expected Type 'AWS::IAM::Role', got %v", roleResource["Type"])
	}
}

func TestFunctionTransformer_WithExplicitRole(t *testing.T) {
	transformer := NewFunctionTransformer()

	// When Role is explicitly specified, no role should be generated
	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Role:    "arn:aws:iam::123456789012:role/MyExistingRole",
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should NOT create IAM role
	if _, hasRole := resources["MyFunctionRole"]; hasRole {
		t.Error("should not create role when explicitly provided")
	}

	// Function should reference the provided role
	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	if props["Role"] != "arn:aws:iam::123456789012:role/MyExistingRole" {
		t.Errorf("expected Role 'arn:aws:iam::123456789012:role/MyExistingRole', got %v", props["Role"])
	}
}

func TestFunctionTransformer_WithEnvironment(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Environment: map[string]interface{}{
			"Variables": map[string]interface{}{
				"TABLE_NAME": "MyTable",
				"DEBUG":      "true",
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	env := props["Environment"].(map[string]interface{})
	vars := env["Variables"].(map[string]interface{})
	if vars["TABLE_NAME"] != "MyTable" {
		t.Errorf("expected TABLE_NAME 'MyTable', got %v", vars["TABLE_NAME"])
	}
	if vars["DEBUG"] != "true" {
		t.Errorf("expected DEBUG 'true', got %v", vars["DEBUG"])
	}
}

func TestFunctionTransformer_WithTags(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Tags: map[string]string{
			"Environment": "production",
			"Team":        "backend",
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	tags := props["Tags"].([]interface{})
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}

func TestFunctionTransformer_WithLayers(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Layers: []interface{}{
			"arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1",
			map[string]interface{}{"Ref": "MyLayerVersion"},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	layers := props["Layers"].([]interface{})
	if len(layers) != 2 {
		t.Errorf("expected 2 layers, got %d", len(layers))
	}
}

func TestFunctionTransformer_WithVpcConfig(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		VpcConfig: map[string]interface{}{
			"SecurityGroupIds": []interface{}{"sg-12345678"},
			"SubnetIds":        []interface{}{"subnet-12345678", "subnet-87654321"},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	vpcConfig := props["VpcConfig"].(map[string]interface{})
	if vpcConfig["SecurityGroupIds"] == nil {
		t.Error("VpcConfig should have SecurityGroupIds")
	}
	if vpcConfig["SubnetIds"] == nil {
		t.Error("VpcConfig should have SubnetIds")
	}

	// Should add VPC access policy to role
	roleResource := resources["MyFunctionRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})
	managedPolicies := roleProps["ManagedPolicyArns"].([]interface{})
	found := false
	for _, p := range managedPolicies {
		if pStr, ok := p.(string); ok && pStr == "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole" {
			found = true
			break
		}
	}
	if !found {
		t.Error("should include VPC access execution role policy")
	}
}

func TestFunctionTransformer_WithArchitectures(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:       "index.handler",
		Runtime:       "nodejs18.x",
		CodeUri:       "s3://bucket/code.zip",
		Architectures: []string{"arm64"},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	archs := props["Architectures"].([]string)
	if len(archs) != 1 || archs[0] != "arm64" {
		t.Errorf("expected Architectures [arm64], got %v", archs)
	}
}

func TestFunctionTransformer_WithS3Event(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"S3Event": map[string]interface{}{
				"Type": "S3",
				"Properties": map[string]interface{}{
					"Bucket": map[string]interface{}{"Ref": "MyBucket"},
					"Events": "s3:ObjectCreated:*",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create Lambda permission for S3
	if _, hasPermission := resources["MyFunctionS3EventPermission"]; !hasPermission {
		t.Error("should create Lambda permission for S3 event")
	}
}

func TestFunctionTransformer_WithSQSEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"SQSEvent": map[string]interface{}{
				"Type": "SQS",
				"Properties": map[string]interface{}{
					"Queue":     "arn:aws:sqs:us-east-1:123456789012:MyQueue",
					"BatchSize": 10,
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create EventSourceMapping for SQS
	esm, hasESM := resources["MyFunctionSQSEvent"].(map[string]interface{})
	if !hasESM {
		t.Fatal("should create EventSourceMapping for SQS event")
	}
	if esm["Type"] != "AWS::Lambda::EventSourceMapping" {
		t.Errorf("expected Type 'AWS::Lambda::EventSourceMapping', got %v", esm["Type"])
	}
}

func TestFunctionTransformer_WithApiEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"ApiEvent": map[string]interface{}{
				"Type": "Api",
				"Properties": map[string]interface{}{
					"Path":   "/hello",
					"Method": "GET",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create Lambda permission for API Gateway
	if _, hasPermission := resources["MyFunctionApiEventPermission"]; !hasPermission {
		t.Error("should create Lambda permission for API event")
	}
}

func TestFunctionTransformer_WithManagedPolicies(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Policies: []interface{}{
			"arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	roleResource := resources["MyFunctionRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})
	managedPolicies := roleProps["ManagedPolicyArns"].([]interface{})

	// Should include both policies plus basic execution role
	if len(managedPolicies) < 2 {
		t.Errorf("expected at least 2 managed policies, got %d", len(managedPolicies))
	}
}

func TestFunctionTransformer_WithInlinePolicies(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Policies: map[string]interface{}{
			"Statement": []interface{}{
				map[string]interface{}{
					"Effect":   "Allow",
					"Action":   []interface{}{"dynamodb:GetItem"},
					"Resource": "*",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	roleResource := resources["MyFunctionRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})
	policies := roleProps["Policies"].([]map[string]interface{})
	if len(policies) == 0 {
		t.Error("expected inline policies to be added")
	}
}

func TestFunctionTransformer_WithFunctionName(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:      "index.handler",
		Runtime:      "nodejs18.x",
		CodeUri:      "s3://bucket/code.zip",
		FunctionName: "my-custom-function-name",
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	if props["FunctionName"] != "my-custom-function-name" {
		t.Errorf("expected FunctionName 'my-custom-function-name', got %v", props["FunctionName"])
	}
}

func TestFunctionTransformer_WithCodeUriObject(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: map[string]interface{}{
			"Bucket":  "my-bucket",
			"Key":     "code/package.zip",
			"Version": "abc123",
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	code := props["Code"].(map[string]interface{})
	if code["S3Bucket"] != "my-bucket" {
		t.Errorf("expected S3Bucket 'my-bucket', got %v", code["S3Bucket"])
	}
	if code["S3Key"] != "code/package.zip" {
		t.Errorf("expected S3Key 'code/package.zip', got %v", code["S3Key"])
	}
	if code["S3ObjectVersion"] != "abc123" {
		t.Errorf("expected S3ObjectVersion 'abc123', got %v", code["S3ObjectVersion"])
	}
}

func TestFunctionTransformer_WithAutoPublishAlias(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:          "index.handler",
		Runtime:          "nodejs18.x",
		CodeUri:          "s3://bucket/code.zip",
		AutoPublishAlias: "live",
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create Version
	if _, hasVersion := resources["MyFunctionVersion"]; !hasVersion {
		t.Error("should create Lambda Version")
	}

	// Should create Alias pointing to version
	aliasResource, hasAlias := resources["MyFunctionAliaslive"].(map[string]interface{})
	if !hasAlias {
		t.Fatal("should create Lambda Alias")
	}
	if aliasResource["Type"] != "AWS::Lambda::Alias" {
		t.Errorf("expected Type 'AWS::Lambda::Alias', got %v", aliasResource["Type"])
	}

	aliasProps := aliasResource["Properties"].(map[string]interface{})
	if aliasProps["Name"] != "live" {
		t.Errorf("expected Name 'live', got %v", aliasProps["Name"])
	}
}

func TestFunctionTransformer_WithDeploymentPreference(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:          "index.handler",
		Runtime:          "nodejs18.x",
		CodeUri:          "s3://bucket/code.zip",
		AutoPublishAlias: "live",
		DeploymentPreference: map[string]interface{}{
			"Type": "Linear10PercentEvery1Minute",
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create CodeDeploy resources
	if _, hasDeploymentGroup := resources["MyFunctionDeploymentGroup"]; !hasDeploymentGroup {
		t.Error("should create CodeDeploy deployment group")
	}
}

func TestFunctionTransformer_WithSnapStart(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "com.example.Handler::handleRequest",
		Runtime: "java11",
		CodeUri: "s3://bucket/code.jar",
		SnapStart: map[string]interface{}{
			"ApplyOn": "PublishedVersions",
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	snapStart := props["SnapStart"].(map[string]interface{})
	if snapStart["ApplyOn"] != "PublishedVersions" {
		t.Errorf("expected ApplyOn 'PublishedVersions', got %v", snapStart["ApplyOn"])
	}
}

func TestFunctionTransformer_WithScheduleEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"ScheduleEvent": map[string]interface{}{
				"Type": "Schedule",
				"Properties": map[string]interface{}{
					"Schedule": "rate(1 hour)",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create EventBridge rule
	if _, hasRule := resources["MyFunctionScheduleEvent"]; !hasRule {
		t.Error("should create EventBridge rule for Schedule event")
	}
}

func TestFunctionTransformer_WithKinesisEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"KinesisEvent": map[string]interface{}{
				"Type": "Kinesis",
				"Properties": map[string]interface{}{
					"Stream":           "arn:aws:kinesis:us-east-1:123456789012:stream/MyStream",
					"StartingPosition": "LATEST",
					"BatchSize":        100,
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create EventSourceMapping
	esm, hasESM := resources["MyFunctionKinesisEvent"].(map[string]interface{})
	if !hasESM {
		t.Fatal("should create EventSourceMapping for Kinesis event")
	}
	if esm["Type"] != "AWS::Lambda::EventSourceMapping" {
		t.Errorf("expected Type 'AWS::Lambda::EventSourceMapping', got %v", esm["Type"])
	}
}

func TestFunctionTransformer_WithDynamoDBEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"DDBEvent": map[string]interface{}{
				"Type": "DynamoDB",
				"Properties": map[string]interface{}{
					"Stream":           "arn:aws:dynamodb:us-east-1:123456789012:table/MyTable/stream/2021-01-01T00:00:00.000",
					"StartingPosition": "TRIM_HORIZON",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create EventSourceMapping
	if _, hasESM := resources["MyFunctionDDBEvent"]; !hasESM {
		t.Error("should create EventSourceMapping for DynamoDB event")
	}
}

func TestFunctionTransformer_WithSNSEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"SNSEvent": map[string]interface{}{
				"Type": "SNS",
				"Properties": map[string]interface{}{
					"Topic": "arn:aws:sns:us-east-1:123456789012:MyTopic",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create Lambda permission and SNS subscription
	if _, hasPermission := resources["MyFunctionSNSEventPermission"]; !hasPermission {
		t.Error("should create Lambda permission for SNS event")
	}
}

func TestFunctionTransformer_WithCloudWatchEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"CloudWatchEvent": map[string]interface{}{
				"Type": "CloudWatchEvent",
				"Properties": map[string]interface{}{
					"Pattern": map[string]interface{}{
						"source":      []interface{}{"aws.ec2"},
						"detail-type": []interface{}{"EC2 Instance State-change Notification"},
					},
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create EventBridge rule
	if _, hasRule := resources["MyFunctionCloudWatchEvent"]; !hasRule {
		t.Error("should create EventBridge rule for CloudWatch event")
	}
}

func TestFunctionTransformer_MultipleEvents(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"ApiEvent": map[string]interface{}{
				"Type": "Api",
				"Properties": map[string]interface{}{
					"Path":   "/hello",
					"Method": "GET",
				},
			},
			"ScheduleEvent": map[string]interface{}{
				"Type": "Schedule",
				"Properties": map[string]interface{}{
					"Schedule": "rate(5 minutes)",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should have resources for both events
	if _, ok := resources["MyFunction"]; !ok {
		t.Error("should have MyFunction resource")
	}
	if _, ok := resources["MyFunctionRole"]; !ok {
		t.Error("should have MyFunctionRole resource")
	}
}

func TestFunctionTransformer_WithDescription(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:     "index.handler",
		Runtime:     "nodejs18.x",
		CodeUri:     "s3://bucket/code.zip",
		Description: "My awesome Lambda function",
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	if props["Description"] != "My awesome Lambda function" {
		t.Errorf("expected Description 'My awesome Lambda function', got %v", props["Description"])
	}
}

func TestFunctionTransformer_WithReservedConcurrentExecutions(t *testing.T) {
	transformer := NewFunctionTransformer()

	reserved := 100
	fn := &Function{
		Handler:                      "index.handler",
		Runtime:                      "nodejs18.x",
		CodeUri:                      "s3://bucket/code.zip",
		ReservedConcurrentExecutions: &reserved,
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	if props["ReservedConcurrentExecutions"] != 100 {
		t.Errorf("expected ReservedConcurrentExecutions 100, got %v", props["ReservedConcurrentExecutions"])
	}
}

func TestFunctionTransformer_WithProvisionedConcurrencyConfig(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:          "index.handler",
		Runtime:          "nodejs18.x",
		CodeUri:          "s3://bucket/code.zip",
		AutoPublishAlias: "live",
		ProvisionedConcurrencyConfig: map[string]interface{}{
			"ProvisionedConcurrentExecutions": 10,
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Provisioned concurrency is set on the alias
	aliasResource := resources["MyFunctionAliaslive"].(map[string]interface{})
	aliasProps := aliasResource["Properties"].(map[string]interface{})
	pcConfig := aliasProps["ProvisionedConcurrencyConfig"].(map[string]interface{})
	if pcConfig["ProvisionedConcurrentExecutions"] != 10 {
		t.Errorf("expected ProvisionedConcurrentExecutions 10, got %v", pcConfig["ProvisionedConcurrentExecutions"])
	}
}

func TestFunctionTransformer_WithTracing(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Tracing: "Active",
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	tracingConfig := props["TracingConfig"].(map[string]interface{})
	if tracingConfig["Mode"] != "Active" {
		t.Errorf("expected Mode 'Active', got %v", tracingConfig["Mode"])
	}

	// Should add X-Ray write policy to role
	roleResource := resources["MyFunctionRole"].(map[string]interface{})
	roleProps := roleResource["Properties"].(map[string]interface{})
	managedPolicies := roleProps["ManagedPolicyArns"].([]interface{})
	found := false
	for _, p := range managedPolicies {
		if pStr, ok := p.(string); ok && pStr == "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess" {
			found = true
			break
		}
	}
	if !found {
		t.Error("should include X-Ray write access policy")
	}
}

func TestFunctionTransformer_WithDeadLetterConfig(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		DeadLetterQueue: map[string]interface{}{
			"TargetArn": "arn:aws:sqs:us-east-1:123456789012:dlq",
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	dlc := props["DeadLetterConfig"].(map[string]interface{})
	if dlc["TargetArn"] != "arn:aws:sqs:us-east-1:123456789012:dlq" {
		t.Errorf("expected TargetArn 'arn:aws:sqs:us-east-1:123456789012:dlq', got %v", dlc["TargetArn"])
	}
}

func TestFunctionTransformer_WithKmsKeyArn(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:   "index.handler",
		Runtime:   "nodejs18.x",
		CodeUri:   "s3://bucket/code.zip",
		KmsKeyArn: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	if props["KmsKeyArn"] != "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012" {
		t.Errorf("expected KmsKeyArn, got %v", props["KmsKeyArn"])
	}
}

func TestFunctionTransformer_WithEphemeralStorage(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:          "index.handler",
		Runtime:          "nodejs18.x",
		CodeUri:          "s3://bucket/code.zip",
		EphemeralStorage: map[string]interface{}{"Size": 1024},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	fnResource := resources["MyFunction"].(map[string]interface{})
	props := fnResource["Properties"].(map[string]interface{})
	ephemeralStorage := props["EphemeralStorage"].(map[string]interface{})
	if ephemeralStorage["Size"] != 1024 {
		t.Errorf("expected Size 1024, got %v", ephemeralStorage["Size"])
	}
}

func TestFunctionTransformer_WithHttpApiEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"HttpApiEvent": map[string]interface{}{
				"Type": "HttpApi",
				"Properties": map[string]interface{}{
					"Path":   "/hello",
					"Method": "GET",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create Lambda permission for API Gateway V2
	if _, hasPermission := resources["MyFunctionHttpApiEventPermission"]; !hasPermission {
		t.Error("should create Lambda permission for HttpApi event")
	}
}

func TestFunctionTransformer_ToJSON(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler:    "index.handler",
		Runtime:    "nodejs18.x",
		CodeUri:    "s3://bucket/code.zip",
		MemorySize: 256,
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should be serializable to JSON
	if _, err := json.Marshal(resources); err != nil {
		t.Fatalf("Failed to marshal resources to JSON: %v", err)
	}
}

func TestFunctionTransformer_WithIoTRuleEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"IoTRule": map[string]interface{}{
				"Type": "IoTRule",
				"Properties": map[string]interface{}{
					"Sql":              "SELECT * FROM 'topic/test'",
					"AwsIotSqlVersion": "2016-03-23",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create IoT TopicRule
	if _, hasRule := resources["MyFunctionIoTRule"]; !hasRule {
		t.Error("should create IoT TopicRule")
	}
}

func TestFunctionTransformer_WithCognitoEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"CognitoEvent": map[string]interface{}{
				"Type": "Cognito",
				"Properties": map[string]interface{}{
					"UserPool": map[string]interface{}{"Ref": "MyUserPool"},
					"Trigger":  "PreSignUp",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create Lambda permission for Cognito
	if _, hasPermission := resources["MyFunctionCognitoEventPermission"]; !hasPermission {
		t.Error("should create Lambda permission for Cognito event")
	}
}

func TestFunctionTransformer_WithMSKEvent(t *testing.T) {
	transformer := NewFunctionTransformer()

	fn := &Function{
		Handler: "index.handler",
		Runtime: "nodejs18.x",
		CodeUri: "s3://bucket/code.zip",
		Events: map[string]interface{}{
			"MSKEvent": map[string]interface{}{
				"Type": "MSK",
				"Properties": map[string]interface{}{
					"Stream":           "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abc123",
					"Topics":           []interface{}{"my-topic"},
					"StartingPosition": "LATEST",
				},
			},
		},
	}

	resources, err := transformer.Transform("MyFunction", fn, nil)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Should create EventSourceMapping for MSK
	if _, hasESM := resources["MyFunctionMSKEvent"]; !hasESM {
		t.Error("should create EventSourceMapping for MSK event")
	}
}
