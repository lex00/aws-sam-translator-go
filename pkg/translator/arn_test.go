package translator

import (
	"testing"
)

func TestArnGenerator_Lambda(t *testing.T) {
	tests := []struct {
		name         string
		region       string
		accountID    string
		functionName string
		expected     string
	}{
		{
			name:         "standard AWS partition",
			region:       "us-east-1",
			accountID:    "123456789012",
			functionName: "MyFunction",
			expected:     "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
		},
		{
			name:         "China partition",
			region:       "cn-north-1",
			accountID:    "123456789012",
			functionName: "MyFunction",
			expected:     "arn:aws-cn:lambda:cn-north-1:123456789012:function:MyFunction",
		},
		{
			name:         "GovCloud partition",
			region:       "us-gov-west-1",
			accountID:    "123456789012",
			functionName: "MyFunction",
			expected:     "arn:aws-us-gov:lambda:us-gov-west-1:123456789012:function:MyFunction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewArnGenerator(tt.region, tt.accountID)
			result := gen.Lambda(tt.functionName)
			if result != tt.expected {
				t.Errorf("Lambda() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestArnGenerator_LambdaAlias(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.LambdaAlias("MyFunction", "prod")
	expected := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction:prod"
	if result != expected {
		t.Errorf("LambdaAlias() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_LambdaVersion(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.LambdaVersion("MyFunction", "1")
	expected := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction:1"
	if result != expected {
		t.Errorf("LambdaVersion() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_LambdaLayer(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.LambdaLayer("MyLayer")
	expected := "arn:aws:lambda:us-east-1:123456789012:layer:MyLayer"
	if result != expected {
		t.Errorf("LambdaLayer() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_LambdaLayerVersion(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.LambdaLayerVersion("MyLayer", 5)
	expected := "arn:aws:lambda:us-east-1:123456789012:layer:MyLayer:5"
	if result != expected {
		t.Errorf("LambdaLayerVersion() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_APIGateway(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.APIGateway("abc123")
	expected := "arn:aws:apigateway:us-east-1::/restapis/abc123"
	if result != expected {
		t.Errorf("APIGateway() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_APIGatewayStage(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.APIGatewayStage("abc123", "prod")
	expected := "arn:aws:apigateway:us-east-1::/restapis/abc123/stages/prod"
	if result != expected {
		t.Errorf("APIGatewayStage() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_APIGatewayV2(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.APIGatewayV2("xyz789")
	expected := "arn:aws:apigateway:us-east-1::/apis/xyz789"
	if result != expected {
		t.Errorf("APIGatewayV2() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_APIGatewayExecute(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.APIGatewayExecute("abc123", "prod", "GET", "/users")
	expected := "arn:aws:execute-api:us-east-1:123456789012:abc123/prod/GET/users"
	if result != expected {
		t.Errorf("APIGatewayExecute() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_IAMRole(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.IAMRole("MyRole")
	expected := "arn:aws:iam::123456789012:role/MyRole"
	if result != expected {
		t.Errorf("IAMRole() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_IAMRoleWithPath(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")

	tests := []struct {
		path     string
		roleName string
		expected string
	}{
		{"/service-role/", "MyRole", "arn:aws:iam::123456789012:role/service-role/MyRole"},
		{"service-role", "MyRole", "arn:aws:iam::123456789012:role/service-role/MyRole"},
		{"/", "MyRole", "arn:aws:iam::123456789012:role/MyRole"},
		{"", "MyRole", "arn:aws:iam::123456789012:role/MyRole"},
	}

	for _, tt := range tests {
		t.Run(tt.path+"/"+tt.roleName, func(t *testing.T) {
			result := gen.IAMRoleWithPath(tt.path, tt.roleName)
			if result != tt.expected {
				t.Errorf("IAMRoleWithPath(%q, %q) = %q, want %q", tt.path, tt.roleName, result, tt.expected)
			}
		})
	}
}

func TestArnGenerator_IAMPolicy(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.IAMPolicy("MyPolicy")
	expected := "arn:aws:iam::123456789012:policy/MyPolicy"
	if result != expected {
		t.Errorf("IAMPolicy() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_IAMManagedPolicy(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.IAMManagedPolicy("AmazonS3ReadOnlyAccess")
	expected := "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
	if result != expected {
		t.Errorf("IAMManagedPolicy() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_S3Bucket(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.S3Bucket("my-bucket")
	expected := "arn:aws:s3:::my-bucket"
	if result != expected {
		t.Errorf("S3Bucket() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_S3Object(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.S3Object("my-bucket", "path/to/object")
	expected := "arn:aws:s3:::my-bucket/path/to/object"
	if result != expected {
		t.Errorf("S3Object() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_DynamoDBTable(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.DynamoDBTable("MyTable")
	expected := "arn:aws:dynamodb:us-east-1:123456789012:table/MyTable"
	if result != expected {
		t.Errorf("DynamoDBTable() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_DynamoDBIndex(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.DynamoDBIndex("MyTable", "GSI1")
	expected := "arn:aws:dynamodb:us-east-1:123456789012:table/MyTable/index/GSI1"
	if result != expected {
		t.Errorf("DynamoDBIndex() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_DynamoDBStream(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.DynamoDBStream("MyTable", "2023-01-01T00:00:00.000")
	expected := "arn:aws:dynamodb:us-east-1:123456789012:table/MyTable/stream/2023-01-01T00:00:00.000"
	if result != expected {
		t.Errorf("DynamoDBStream() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_SNSTopic(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.SNSTopic("MyTopic")
	expected := "arn:aws:sns:us-east-1:123456789012:MyTopic"
	if result != expected {
		t.Errorf("SNSTopic() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_SQSQueue(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.SQSQueue("MyQueue")
	expected := "arn:aws:sqs:us-east-1:123456789012:MyQueue"
	if result != expected {
		t.Errorf("SQSQueue() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_KinesisStream(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.KinesisStream("MyStream")
	expected := "arn:aws:kinesis:us-east-1:123456789012:stream/MyStream"
	if result != expected {
		t.Errorf("KinesisStream() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_StepFunctionsStateMachine(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.StepFunctionsStateMachine("MyStateMachine")
	expected := "arn:aws:states:us-east-1:123456789012:stateMachine:MyStateMachine"
	if result != expected {
		t.Errorf("StepFunctionsStateMachine() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_EventsRule(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.EventsRule("MyRule")
	expected := "arn:aws:events:us-east-1:123456789012:rule/MyRule"
	if result != expected {
		t.Errorf("EventsRule() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_EventsEventBus(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.EventsEventBus("MyEventBus")
	expected := "arn:aws:events:us-east-1:123456789012:event-bus/MyEventBus"
	if result != expected {
		t.Errorf("EventsEventBus() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_CloudWatchLogGroup(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.CloudWatchLogGroup("/aws/lambda/MyFunction")
	expected := "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/MyFunction"
	if result != expected {
		t.Errorf("CloudWatchLogGroup() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_CloudWatchAlarm(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.CloudWatchAlarm("MyAlarm")
	expected := "arn:aws:cloudwatch:us-east-1:123456789012:alarm:MyAlarm"
	if result != expected {
		t.Errorf("CloudWatchAlarm() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_SecretsManagerSecret(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.SecretsManagerSecret("MySecret")
	expected := "arn:aws:secretsmanager:us-east-1:123456789012:secret:MySecret"
	if result != expected {
		t.Errorf("SecretsManagerSecret() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_KMSKey(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.KMSKey("12345678-1234-1234-1234-123456789012")
	expected := "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
	if result != expected {
		t.Errorf("KMSKey() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_KMSAlias(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")

	// Without prefix
	result := gen.KMSAlias("MyKey")
	expected := "arn:aws:kms:us-east-1:123456789012:alias/MyKey"
	if result != expected {
		t.Errorf("KMSAlias() = %q, want %q", result, expected)
	}

	// With prefix
	result = gen.KMSAlias("alias/MyKey")
	if result != expected {
		t.Errorf("KMSAlias() with alias prefix = %q, want %q", result, expected)
	}
}

func TestArnGenerator_CognitoUserPool(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.CognitoUserPool("us-east-1_abcd1234")
	expected := "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_abcd1234"
	if result != expected {
		t.Errorf("CognitoUserPool() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_CodeDeployApplication(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.CodeDeployApplication("MyApp")
	expected := "arn:aws:codedeploy:us-east-1:123456789012:application:MyApp"
	if result != expected {
		t.Errorf("CodeDeployApplication() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_CodeDeployDeploymentGroup(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.CodeDeployDeploymentGroup("MyApp", "MyDeploymentGroup")
	expected := "arn:aws:codedeploy:us-east-1:123456789012:deploymentgroup:MyApp/MyDeploymentGroup"
	if result != expected {
		t.Errorf("CodeDeployDeploymentGroup() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_Generic(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.Generic("custom", "resource/name")
	expected := "arn:aws:custom:us-east-1:123456789012:resource/name"
	if result != expected {
		t.Errorf("Generic() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_GenericGlobal(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.GenericGlobal("iam", "role/MyRole")
	expected := "arn:aws:iam::123456789012:role/MyRole"
	if result != expected {
		t.Errorf("GenericGlobal() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_GenericNoAccount(t *testing.T) {
	gen := NewArnGenerator("us-east-1", "123456789012")
	result := gen.GenericNoAccount("s3", "my-bucket")
	expected := "arn:aws:s3:us-east-1::my-bucket"
	if result != expected {
		t.Errorf("GenericNoAccount() = %q, want %q", result, expected)
	}
}

func TestArnGenerator_AllPartitions(t *testing.T) {
	partitions := []struct {
		region    string
		partition string
	}{
		{"us-east-1", "aws"},
		{"eu-west-1", "aws"},
		{"cn-north-1", "aws-cn"},
		{"cn-northwest-1", "aws-cn"},
		{"us-gov-west-1", "aws-us-gov"},
		{"us-gov-east-1", "aws-us-gov"},
	}

	for _, p := range partitions {
		t.Run(p.region, func(t *testing.T) {
			gen := NewArnGenerator(p.region, "123456789012")
			arn := gen.Lambda("MyFunction")
			parsed, err := ParseARN(arn)
			if err != nil {
				t.Fatalf("Failed to parse ARN: %v", err)
			}
			if parsed.Partition != p.partition {
				t.Errorf("Expected partition %q for region %q, got %q", p.partition, p.region, parsed.Partition)
			}
		})
	}
}

func TestParseARN(t *testing.T) {
	tests := []struct {
		name      string
		arnStr    string
		wantErr   bool
		partition string
		service   string
		region    string
		accountID string
		resource  string
	}{
		{
			name:      "valid Lambda ARN",
			arnStr:    "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
			partition: "aws",
			service:   "lambda",
			region:    "us-east-1",
			accountID: "123456789012",
			resource:  "function:MyFunction",
		},
		{
			name:      "valid S3 ARN",
			arnStr:    "arn:aws:s3:::my-bucket",
			partition: "aws",
			service:   "s3",
			region:    "",
			accountID: "",
			resource:  "my-bucket",
		},
		{
			name:      "valid IAM ARN",
			arnStr:    "arn:aws:iam::123456789012:role/MyRole",
			partition: "aws",
			service:   "iam",
			region:    "",
			accountID: "123456789012",
			resource:  "role/MyRole",
		},
		{
			name:    "empty ARN",
			arnStr:  "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			arnStr:  "not-an-arn",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arn, err := ParseARN(tt.arnStr)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseARN() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseARN() unexpected error: %v", err)
				return
			}
			if arn.Partition != tt.partition {
				t.Errorf("Partition = %q, want %q", arn.Partition, tt.partition)
			}
			if arn.Service != tt.service {
				t.Errorf("Service = %q, want %q", arn.Service, tt.service)
			}
			if arn.Region != tt.region {
				t.Errorf("Region = %q, want %q", arn.Region, tt.region)
			}
			if arn.AccountID != tt.accountID {
				t.Errorf("AccountID = %q, want %q", arn.AccountID, tt.accountID)
			}
			if arn.Resource != tt.resource {
				t.Errorf("Resource = %q, want %q", arn.Resource, tt.resource)
			}
		})
	}
}

func TestARN_String(t *testing.T) {
	arn := ARN{
		Partition: "aws",
		Service:   "lambda",
		Region:    "us-east-1",
		AccountID: "123456789012",
		Resource:  "function:MyFunction",
	}
	expected := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction"
	if arn.String() != expected {
		t.Errorf("ARN.String() = %q, want %q", arn.String(), expected)
	}
}

func TestIsValidARN(t *testing.T) {
	tests := []struct {
		arn   string
		valid bool
	}{
		{"arn:aws:lambda:us-east-1:123456789012:function:MyFunction", true},
		{"arn:aws:s3:::my-bucket", true},
		{"arn:aws-cn:lambda:cn-north-1:123456789012:function:MyFunction", true},
		{"arn:aws-us-gov:lambda:us-gov-west-1:123456789012:function:MyFunction", true},
		{"not-an-arn", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.arn, func(t *testing.T) {
			if IsValidARN(tt.arn) != tt.valid {
				t.Errorf("IsValidARN(%q) = %v, want %v", tt.arn, !tt.valid, tt.valid)
			}
		})
	}
}

func TestGetPartition(t *testing.T) {
	partition, err := GetPartition("arn:aws:lambda:us-east-1:123456789012:function:MyFunction")
	if err != nil {
		t.Fatalf("GetPartition() unexpected error: %v", err)
	}
	if partition != "aws" {
		t.Errorf("GetPartition() = %q, want %q", partition, "aws")
	}
}

func TestGetService(t *testing.T) {
	service, err := GetService("arn:aws:lambda:us-east-1:123456789012:function:MyFunction")
	if err != nil {
		t.Fatalf("GetService() unexpected error: %v", err)
	}
	if service != "lambda" {
		t.Errorf("GetService() = %q, want %q", service, "lambda")
	}
}

func TestGetResource(t *testing.T) {
	resource, err := GetResource("arn:aws:lambda:us-east-1:123456789012:function:MyFunction")
	if err != nil {
		t.Fatalf("GetResource() unexpected error: %v", err)
	}
	if resource != "function:MyFunction" {
		t.Errorf("GetResource() = %q, want %q", resource, "function:MyFunction")
	}
}

func TestReplacePartition(t *testing.T) {
	result, err := ReplacePartition("arn:aws:lambda:us-east-1:123456789012:function:MyFunction", "aws-cn")
	if err != nil {
		t.Fatalf("ReplacePartition() unexpected error: %v", err)
	}
	expected := "arn:aws-cn:lambda:us-east-1:123456789012:function:MyFunction"
	if result != expected {
		t.Errorf("ReplacePartition() = %q, want %q", result, expected)
	}
}

func TestReplaceRegion(t *testing.T) {
	result, err := ReplaceRegion("arn:aws:lambda:us-east-1:123456789012:function:MyFunction", "eu-west-1")
	if err != nil {
		t.Fatalf("ReplaceRegion() unexpected error: %v", err)
	}
	expected := "arn:aws:lambda:eu-west-1:123456789012:function:MyFunction"
	if result != expected {
		t.Errorf("ReplaceRegion() = %q, want %q", result, expected)
	}
}

func TestReplaceAccountID(t *testing.T) {
	result, err := ReplaceAccountID("arn:aws:lambda:us-east-1:123456789012:function:MyFunction", "999999999999")
	if err != nil {
		t.Fatalf("ReplaceAccountID() unexpected error: %v", err)
	}
	expected := "arn:aws:lambda:us-east-1:999999999999:function:MyFunction"
	if result != expected {
		t.Errorf("ReplaceAccountID() = %q, want %q", result, expected)
	}
}

func TestNewArnGeneratorWithPartition(t *testing.T) {
	gen := NewArnGeneratorWithPartition("aws-cn", "cn-north-1", "123456789012")
	result := gen.Lambda("MyFunction")
	expected := "arn:aws-cn:lambda:cn-north-1:123456789012:function:MyFunction"
	if result != expected {
		t.Errorf("Lambda() with explicit partition = %q, want %q", result, expected)
	}
}
