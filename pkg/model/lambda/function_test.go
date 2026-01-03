package lambda

import (
	"testing"
)

func TestNewFunction(t *testing.T) {
	code := NewCodeFromS3("my-bucket", "my-key")
	role := "arn:aws:iam::123456789012:role/MyRole"
	fn := NewFunction(code, role)

	if fn.Code == nil {
		t.Error("expected Code to be set")
	}
	if fn.Role != role {
		t.Errorf("expected Role %s, got %v", role, fn.Role)
	}
}

func TestNewCodeFromS3(t *testing.T) {
	code := NewCodeFromS3("my-bucket", "my-key")

	if code.S3Bucket != "my-bucket" {
		t.Errorf("expected S3Bucket 'my-bucket', got %v", code.S3Bucket)
	}
	if code.S3Key != "my-key" {
		t.Errorf("expected S3Key 'my-key', got %v", code.S3Key)
	}
}

func TestNewCodeFromS3WithVersion(t *testing.T) {
	code := NewCodeFromS3WithVersion("my-bucket", "my-key", "v1")

	if code.S3Bucket != "my-bucket" {
		t.Errorf("expected S3Bucket 'my-bucket', got %v", code.S3Bucket)
	}
	if code.S3Key != "my-key" {
		t.Errorf("expected S3Key 'my-key', got %v", code.S3Key)
	}
	if code.S3ObjectVersion != "v1" {
		t.Errorf("expected S3ObjectVersion 'v1', got %v", code.S3ObjectVersion)
	}
}

func TestNewCodeFromImage(t *testing.T) {
	code := NewCodeFromImage("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-repo:latest")

	if code.ImageUri != "123456789012.dkr.ecr.us-east-1.amazonaws.com/my-repo:latest" {
		t.Errorf("unexpected ImageUri: %v", code.ImageUri)
	}
}

func TestNewCodeFromZip(t *testing.T) {
	zipContent := "exports.handler = async () => 'hello'"
	code := NewCodeFromZip(zipContent)

	if code.ZipFile != zipContent {
		t.Errorf("expected ZipFile content, got %v", code.ZipFile)
	}
}

func TestFunctionToCloudFormation_Minimal(t *testing.T) {
	code := NewCodeFromS3("my-bucket", "my-key")
	role := "arn:aws:iam::123456789012:role/MyRole"
	fn := NewFunction(code, role)

	result := fn.ToCloudFormation()

	if result["Type"] != ResourceTypeFunction {
		t.Errorf("expected Type %s, got %v", ResourceTypeFunction, result["Type"])
	}

	props := result["Properties"].(map[string]interface{})
	if props["Role"] != role {
		t.Errorf("expected Role in properties")
	}

	codeMap := props["Code"].(map[string]interface{})
	if codeMap["S3Bucket"] != "my-bucket" {
		t.Errorf("expected S3Bucket in Code")
	}
}

func TestFunctionToCloudFormation_Full(t *testing.T) {
	code := NewCodeFromS3("my-bucket", "my-key")
	role := "arn:aws:iam::123456789012:role/MyRole"
	reservedConcurrency := 10

	fn := &Function{
		Architectures:                []string{"arm64"},
		Code:                         code,
		CodeSigningConfigArn:         "arn:aws:lambda:us-east-1:123456789012:code-signing-config:csc-1234",
		DeadLetterConfig:             &DeadLetterConfig{TargetArn: "arn:aws:sqs:us-east-1:123456789012:dlq"},
		Description:                  "My test function",
		Environment:                  &Environment{Variables: map[string]interface{}{"KEY": "value"}},
		EphemeralStorage:             &EphemeralStorage{Size: 1024},
		FunctionName:                 "my-function",
		Handler:                      "index.handler",
		KmsKeyArn:                    "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012",
		Layers:                       []interface{}{"arn:aws:lambda:us-east-1:123456789012:layer:my-layer:1"},
		LoggingConfig:                &LoggingConfig{LogFormat: "JSON", ApplicationLogLevel: "INFO"},
		MemorySize:                   512,
		PackageType:                  "Zip",
		RecursiveLoop:                "Terminate",
		ReservedConcurrentExecutions: &reservedConcurrency,
		Role:                         role,
		Runtime:                      "nodejs18.x",
		RuntimeManagementConfig:      &RuntimeManagementConfig{UpdateRuntimeOn: "Auto"},
		SnapStart:                    &SnapStart{ApplyOn: "PublishedVersions"},
		Tags:                         []Tag{{Key: "Environment", Value: "test"}},
		Timeout:                      30,
		TracingConfig:                &TracingConfig{Mode: "Active"},
		VpcConfig: &VpcConfig{
			SecurityGroupIds: []interface{}{"sg-12345678"},
			SubnetIds:        []interface{}{"subnet-12345678"},
		},
	}

	result := fn.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	// Verify architectures
	archs := props["Architectures"].([]string)
	if len(archs) != 1 || archs[0] != "arm64" {
		t.Errorf("unexpected Architectures: %v", archs)
	}

	// Verify handler
	if props["Handler"] != "index.handler" {
		t.Errorf("expected Handler 'index.handler', got %v", props["Handler"])
	}

	// Verify memory
	if props["MemorySize"] != 512 {
		t.Errorf("expected MemorySize 512, got %v", props["MemorySize"])
	}

	// Verify timeout
	if props["Timeout"] != 30 {
		t.Errorf("expected Timeout 30, got %v", props["Timeout"])
	}

	// Verify runtime
	if props["Runtime"] != "nodejs18.x" {
		t.Errorf("expected Runtime 'nodejs18.x', got %v", props["Runtime"])
	}

	// Verify environment
	env := props["Environment"].(map[string]interface{})
	vars := env["Variables"].(map[string]interface{})
	if vars["KEY"] != "value" {
		t.Errorf("expected environment variable KEY=value")
	}

	// Verify VPC config
	vpc := props["VpcConfig"].(map[string]interface{})
	if vpc["SecurityGroupIds"] == nil {
		t.Error("expected SecurityGroupIds in VpcConfig")
	}

	// Verify tracing config
	tracing := props["TracingConfig"].(map[string]interface{})
	if tracing["Mode"] != "Active" {
		t.Errorf("expected TracingConfig Mode 'Active', got %v", tracing["Mode"])
	}

	// Verify reserved concurrency
	if props["ReservedConcurrentExecutions"] != 10 {
		t.Errorf("expected ReservedConcurrentExecutions 10, got %v", props["ReservedConcurrentExecutions"])
	}
}

func TestFunctionWithFileSystemConfigs(t *testing.T) {
	code := NewCodeFromS3("my-bucket", "my-key")
	fn := NewFunction(code, "arn:aws:iam::123456789012:role/MyRole")
	fn.FileSystemConfigs = []FileSystemConfig{
		{
			Arn:            "arn:aws:elasticfilesystem:us-east-1:123456789012:access-point/fsap-12345678",
			LocalMountPath: "/mnt/efs",
		},
	}

	result := fn.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	configs := props["FileSystemConfigs"].([]map[string]interface{})
	if len(configs) != 1 {
		t.Errorf("expected 1 FileSystemConfig, got %d", len(configs))
	}
	if configs[0]["LocalMountPath"] != "/mnt/efs" {
		t.Errorf("expected LocalMountPath '/mnt/efs', got %v", configs[0]["LocalMountPath"])
	}
}

func TestFunctionWithImageConfig(t *testing.T) {
	code := NewCodeFromImage("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-repo:latest")
	fn := NewFunction(code, "arn:aws:iam::123456789012:role/MyRole")
	fn.PackageType = "Image"
	fn.ImageConfig = &ImageConfig{
		Command:          []string{"app.handler"},
		EntryPoint:       []string{"/lambda-entrypoint.sh"},
		WorkingDirectory: "/var/task",
	}

	result := fn.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	imageConfig := props["ImageConfig"].(map[string]interface{})
	if imageConfig["WorkingDirectory"] != "/var/task" {
		t.Errorf("expected WorkingDirectory '/var/task', got %v", imageConfig["WorkingDirectory"])
	}
}

func TestFunctionWithTags(t *testing.T) {
	code := NewCodeFromS3("my-bucket", "my-key")
	fn := NewFunction(code, "arn:aws:iam::123456789012:role/MyRole")
	fn.Tags = []Tag{
		{Key: "Environment", Value: "production"},
		{Key: "Team", Value: "platform"},
	}

	result := fn.ToCloudFormation()
	props := result["Properties"].(map[string]interface{})

	tags := props["Tags"].([]map[string]interface{})
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
	if tags[0]["Key"] != "Environment" || tags[0]["Value"] != "production" {
		t.Errorf("unexpected first tag: %v", tags[0])
	}
}

func TestCodeToMap(t *testing.T) {
	tests := []struct {
		name     string
		code     *Code
		expected map[string]interface{}
	}{
		{
			name: "S3 code",
			code: &Code{S3Bucket: "bucket", S3Key: "key"},
			expected: map[string]interface{}{
				"S3Bucket": "bucket",
				"S3Key":    "key",
			},
		},
		{
			name: "S3 code with version",
			code: &Code{S3Bucket: "bucket", S3Key: "key", S3ObjectVersion: "v1"},
			expected: map[string]interface{}{
				"S3Bucket":        "bucket",
				"S3Key":           "key",
				"S3ObjectVersion": "v1",
			},
		},
		{
			name: "Image code",
			code: &Code{ImageUri: "uri"},
			expected: map[string]interface{}{
				"ImageUri": "uri",
			},
		},
		{
			name: "Zip code",
			code: &Code{ZipFile: "code"},
			expected: map[string]interface{}{
				"ZipFile": "code",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.code.toMap()
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("expected %s=%v, got %v", k, v, result[k])
				}
			}
		})
	}
}

func TestVpcConfigWithIpv6(t *testing.T) {
	vpc := &VpcConfig{
		Ipv6AllowedForDualStack: true,
		SecurityGroupIds:        []interface{}{"sg-123"},
		SubnetIds:               []interface{}{"subnet-123"},
	}

	result := vpc.toMap()
	if result["Ipv6AllowedForDualStack"] != true {
		t.Error("expected Ipv6AllowedForDualStack to be true")
	}
}
