package iam

import (
	"strings"
	"testing"
)

func TestNewManagedPolicyARN(t *testing.T) {
	tests := []struct {
		name         string
		arn          interface{}
		isAWSManaged bool
	}{
		{
			name:         "AWS managed policy",
			arn:          "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
			isAWSManaged: true,
		},
		{
			name:         "Customer managed policy",
			arn:          "arn:aws:iam::123456789012:policy/my-policy",
			isAWSManaged: false,
		},
		{
			name:         "Intrinsic function",
			arn:          map[string]interface{}{"Ref": "PolicyArn"},
			isAWSManaged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arn := NewManagedPolicyARN(tt.arn)

			if arn.IsAWSManaged != tt.isAWSManaged {
				t.Errorf("expected IsAWSManaged %v, got %v", tt.isAWSManaged, arn.IsAWSManaged)
			}
		})
	}
}

func TestNewAWSManagedPolicyARN(t *testing.T) {
	tests := []struct {
		name       string
		policyName string
		partition  string
		expected   string
	}{
		{
			name:       "default partition",
			policyName: "AWSLambdaBasicExecutionRole",
			partition:  "",
			expected:   "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
		},
		{
			name:       "aws partition",
			policyName: "AWSLambdaBasicExecutionRole",
			partition:  "aws",
			expected:   "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
		},
		{
			name:       "aws-cn partition",
			policyName: "AWSLambdaBasicExecutionRole",
			partition:  "aws-cn",
			expected:   "arn:aws-cn:iam::aws:policy/AWSLambdaBasicExecutionRole",
		},
		{
			name:       "aws-us-gov partition",
			policyName: "AWSLambdaBasicExecutionRole",
			partition:  "aws-us-gov",
			expected:   "arn:aws-us-gov:iam::aws:policy/AWSLambdaBasicExecutionRole",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arn := NewAWSManagedPolicyARN(tt.policyName, tt.partition)

			if arn.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, arn.String())
			}

			if !arn.IsAWSManaged {
				t.Error("expected IsAWSManaged to be true")
			}
		})
	}
}

func TestNewCustomManagedPolicyARN(t *testing.T) {
	tests := []struct {
		name       string
		policyName string
		accountID  string
		partition  string
		expected   string
	}{
		{
			name:       "default partition",
			policyName: "my-policy",
			accountID:  "123456789012",
			partition:  "",
			expected:   "arn:aws:iam::123456789012:policy/my-policy",
		},
		{
			name:       "aws-cn partition",
			policyName: "my-policy",
			accountID:  "123456789012",
			partition:  "aws-cn",
			expected:   "arn:aws-cn:iam::123456789012:policy/my-policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arn := NewCustomManagedPolicyARN(tt.policyName, tt.accountID, tt.partition)

			if arn.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, arn.String())
			}

			if arn.IsAWSManaged {
				t.Error("expected IsAWSManaged to be false")
			}
		})
	}
}

func TestNewCustomManagedPolicyARNWithPath(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		policyName string
		accountID  string
		partition  string
		expected   string
	}{
		{
			name:       "with path",
			path:       "/service-role/",
			policyName: "my-policy",
			accountID:  "123456789012",
			partition:  "aws",
			expected:   "arn:aws:iam::123456789012:policy/service-role/my-policy",
		},
		{
			name:       "path without leading slash",
			path:       "service-role",
			policyName: "my-policy",
			accountID:  "123456789012",
			partition:  "aws",
			expected:   "arn:aws:iam::123456789012:policy/service-role/my-policy",
		},
		{
			name:       "path without trailing slash",
			path:       "/service-role",
			policyName: "my-policy",
			accountID:  "123456789012",
			partition:  "aws",
			expected:   "arn:aws:iam::123456789012:policy/service-role/my-policy",
		},
		{
			name:       "empty path",
			path:       "",
			policyName: "my-policy",
			accountID:  "123456789012",
			partition:  "aws",
			expected:   "arn:aws:iam::123456789012:policy/my-policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arn := NewCustomManagedPolicyARNWithPath(tt.path, tt.policyName, tt.accountID, tt.partition)

			if arn.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, arn.String())
			}
		})
	}
}

func TestManagedPolicyARNValue(t *testing.T) {
	// Test with string
	arn := NewManagedPolicyARN("arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole")
	if arn.Value() != "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole" {
		t.Errorf("unexpected value: %v", arn.Value())
	}

	// Test with map
	intrinsic := map[string]interface{}{"Ref": "PolicyArn"}
	arn = NewManagedPolicyARN(intrinsic)
	value, ok := arn.Value().(map[string]interface{})
	if !ok {
		t.Fatalf("expected value to be map")
	}
	if _, hasRef := value["Ref"]; !hasRef {
		t.Errorf("expected intrinsic function to be preserved")
	}
}

func TestManagedPolicyResolver(t *testing.T) {
	resolver := NewManagedPolicyResolver("aws", "123456789012")

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "full ARN",
			input:    "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
			expected: "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
		},
		{
			name:     "AWS managed policy name",
			input:    "AWSLambdaBasicExecutionRole",
			expected: "arn:aws:iam::aws:policy/AWSLambdaBasicExecutionRole",
		},
		{
			name:     "Amazon managed policy name",
			input:    "AmazonDynamoDBFullAccess",
			expected: "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
		},
		{
			name:     "custom policy name",
			input:    "my-custom-policy",
			expected: "arn:aws:iam::123456789012:policy/my-custom-policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.Resolve(tt.input)

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestManagedPolicyResolverWithIntrinsic(t *testing.T) {
	resolver := NewManagedPolicyResolver("aws", "123456789012")

	intrinsic := map[string]interface{}{"Ref": "PolicyArn"}
	result := resolver.Resolve(intrinsic)

	// Should preserve the intrinsic function
	value, ok := result.Value().(map[string]interface{})
	if !ok {
		t.Fatalf("expected value to be map")
	}

	if _, ok := value["Ref"]; !ok {
		t.Error("expected Ref to be preserved")
	}
}

func TestManagedPolicyResolverResolveMany(t *testing.T) {
	resolver := NewManagedPolicyResolver("aws", "123456789012")

	policies := []interface{}{
		"AWSLambdaBasicExecutionRole",
		"arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess",
		"my-custom-policy",
	}

	results := resolver.ResolveMany(policies)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if !results[0].IsAWSManaged {
		t.Error("expected first result to be AWS managed")
	}

	if !results[1].IsAWSManaged {
		t.Error("expected second result to be AWS managed")
	}

	if results[2].IsAWSManaged {
		t.Error("expected third result to not be AWS managed")
	}
}

func TestManagedPolicyResolverResolveManyToValues(t *testing.T) {
	resolver := NewManagedPolicyResolver("aws", "123456789012")

	policies := []interface{}{
		"AWSLambdaBasicExecutionRole",
		map[string]interface{}{"Ref": "PolicyArn"},
	}

	results := resolver.ResolveManyToValues(policies)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// First should be a string ARN
	arnStr, ok := results[0].(string)
	if !ok {
		t.Fatalf("expected first result to be string")
	}

	if !strings.HasPrefix(arnStr, "arn:aws:iam::aws:policy/") {
		t.Errorf("unexpected ARN: %s", arnStr)
	}

	// Second should be preserved intrinsic
	intrinsic, ok := results[1].(map[string]interface{})
	if !ok {
		t.Fatalf("expected second result to be map")
	}

	if _, ok := intrinsic["Ref"]; !ok {
		t.Error("expected Ref to be preserved")
	}
}

func TestIsAWSManagedPolicyName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"AWSLambdaBasicExecutionRole", true},
		{"AWSLambdaVPCAccessExecutionRole", true},
		{"AmazonDynamoDBFullAccess", true},
		{"AmazonS3ReadOnlyAccess", true},
		{"CloudWatchLogsFullAccess", true},
		{"AdministratorAccess", true},
		{"PowerUserAccess", true},
		{"ReadOnlyAccess", true},
		{"my-custom-policy", false},
		{"custom-policy-name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAWSManagedPolicyName(tt.name)

			if result != tt.expected {
				t.Errorf("isAWSManagedPolicyName(%s) = %v, expected %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestManagedPolicyARNHelpers(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) string
		expected string
	}{
		{
			name:     "AWSLambdaBasicExecutionRole",
			fn:       AWSLambdaBasicExecutionRole,
			expected: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
		},
		{
			name:     "AWSLambdaVPCAccessExecutionRole",
			fn:       AWSLambdaVPCAccessExecutionRole,
			expected: "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
		},
		{
			name:     "AWSLambdaDynamoDBExecutionRole",
			fn:       AWSLambdaDynamoDBExecutionRole,
			expected: "arn:aws:iam::aws:policy/service-role/AWSLambdaDynamoDBExecutionRole",
		},
		{
			name:     "AWSXrayWriteOnlyAccess",
			fn:       AWSXrayWriteOnlyAccess,
			expected: "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess",
		},
		{
			name:     "CloudWatchLogsFullAccess",
			fn:       CloudWatchLogsFullAccess,
			expected: "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess",
		},
		{
			name:     "AmazonDynamoDBFullAccess",
			fn:       AmazonDynamoDBFullAccess,
			expected: "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
		},
		{
			name:     "AmazonS3FullAccess",
			fn:       AmazonS3FullAccess,
			expected: "arn:aws:iam::aws:policy/AmazonS3FullAccess",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn("aws")

			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestManagedPolicyARNHelpersWithPartition(t *testing.T) {
	// Test with aws-cn partition
	arn := AWSLambdaBasicExecutionRole("aws-cn")
	expected := "arn:aws-cn:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

	if arn != expected {
		t.Errorf("expected %s, got %s", expected, arn)
	}

	// Test with aws-us-gov partition
	arn = AmazonS3ReadOnlyAccess("aws-us-gov")
	expected = "arn:aws-us-gov:iam::aws:policy/AmazonS3ReadOnlyAccess"

	if arn != expected {
		t.Errorf("expected %s, got %s", expected, arn)
	}
}

func TestManagedPolicyResolverDefaultPartition(t *testing.T) {
	// Test with empty partition - should default to "aws"
	resolver := NewManagedPolicyResolver("", "123456789012")

	result := resolver.Resolve("AWSLambdaBasicExecutionRole")

	if !strings.HasPrefix(result.String(), "arn:aws:") {
		t.Errorf("expected arn:aws: prefix, got %s", result.String())
	}
}
