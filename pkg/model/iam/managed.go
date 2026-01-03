package iam

import (
	"fmt"
	"strings"
)

// ManagedPolicyARN represents a reference to an AWS managed policy or custom managed policy.
type ManagedPolicyARN struct {
	// ARN is the full ARN of the managed policy.
	ARN interface{}

	// IsAWSManaged indicates if this is an AWS-managed policy.
	IsAWSManaged bool
}

// NewManagedPolicyARN creates a ManagedPolicyARN from a full ARN.
func NewManagedPolicyARN(arn interface{}) *ManagedPolicyARN {
	isAWSManaged := false
	if arnStr, ok := arn.(string); ok {
		isAWSManaged = strings.Contains(arnStr, ":iam::aws:policy/")
	}
	return &ManagedPolicyARN{
		ARN:          arn,
		IsAWSManaged: isAWSManaged,
	}
}

// NewAWSManagedPolicyARN creates a ManagedPolicyARN for an AWS-managed policy by name.
func NewAWSManagedPolicyARN(policyName string, partition string) *ManagedPolicyARN {
	if partition == "" {
		partition = "aws"
	}
	arn := fmt.Sprintf("arn:%s:iam::aws:policy/%s", partition, policyName)
	return &ManagedPolicyARN{
		ARN:          arn,
		IsAWSManaged: true,
	}
}

// NewCustomManagedPolicyARN creates a ManagedPolicyARN for a customer-managed policy.
func NewCustomManagedPolicyARN(policyName, accountID, partition string) *ManagedPolicyARN {
	if partition == "" {
		partition = "aws"
	}
	arn := fmt.Sprintf("arn:%s:iam::%s:policy/%s", partition, accountID, policyName)
	return &ManagedPolicyARN{
		ARN:          arn,
		IsAWSManaged: false,
	}
}

// NewCustomManagedPolicyARNWithPath creates a ManagedPolicyARN for a customer-managed policy with path.
func NewCustomManagedPolicyARNWithPath(path, policyName, accountID, partition string) *ManagedPolicyARN {
	if partition == "" {
		partition = "aws"
	}
	// Ensure path starts with / and ends with /
	if path != "" {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}
	} else {
		path = "/"
	}
	arn := fmt.Sprintf("arn:%s:iam::%s:policy%s%s", partition, accountID, path, policyName)
	return &ManagedPolicyARN{
		ARN:          arn,
		IsAWSManaged: false,
	}
}

// String returns the ARN as a string.
func (m *ManagedPolicyARN) String() string {
	if s, ok := m.ARN.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", m.ARN)
}

// Value returns the ARN value (could be string or intrinsic function).
func (m *ManagedPolicyARN) Value() interface{} {
	return m.ARN
}

// ManagedPolicyResolver resolves managed policy references to ARNs.
type ManagedPolicyResolver struct {
	partition string
	accountID string
}

// NewManagedPolicyResolver creates a new ManagedPolicyResolver.
func NewManagedPolicyResolver(partition, accountID string) *ManagedPolicyResolver {
	if partition == "" {
		partition = "aws"
	}
	return &ManagedPolicyResolver{
		partition: partition,
		accountID: accountID,
	}
}

// Resolve resolves a policy reference to a ManagedPolicyARN.
// The input can be:
// - A full ARN string (returned as-is)
// - An AWS managed policy name (e.g., "AWSLambdaBasicExecutionRole")
// - A map with CloudFormation intrinsic function (e.g., {"Ref": "PolicyArn"})
func (r *ManagedPolicyResolver) Resolve(policy interface{}) *ManagedPolicyARN {
	switch p := policy.(type) {
	case string:
		// Check if it's already a full ARN
		if strings.HasPrefix(p, "arn:") {
			return NewManagedPolicyARN(p)
		}
		// Check if it looks like an AWS managed policy name
		if isAWSManagedPolicyName(p) {
			return NewAWSManagedPolicyARN(p, r.partition)
		}
		// Assume it's a custom policy name
		return NewCustomManagedPolicyARN(p, r.accountID, r.partition)

	case map[string]interface{}:
		// CloudFormation intrinsic function - return as-is
		return NewManagedPolicyARN(p)

	default:
		// Return as-is for any other type
		return NewManagedPolicyARN(policy)
	}
}

// ResolveMany resolves multiple policy references to ManagedPolicyARNs.
func (r *ManagedPolicyResolver) ResolveMany(policies []interface{}) []*ManagedPolicyARN {
	result := make([]*ManagedPolicyARN, len(policies))
	for i, p := range policies {
		result[i] = r.Resolve(p)
	}
	return result
}

// ResolveManyToValues resolves multiple policy references and returns their values.
func (r *ManagedPolicyResolver) ResolveManyToValues(policies []interface{}) []interface{} {
	result := make([]interface{}, len(policies))
	for i, p := range policies {
		result[i] = r.Resolve(p).Value()
	}
	return result
}

// Common AWS managed policy names and their mapping.
var awsManagedPolicyPrefixes = []string{
	"AWS",
	"Amazon",
	"IAM",
	"CloudWatch",
	"Alexa",
	"Lex",
	"RDS",
	"EC2",
	"ECS",
	"EKS",
	"ElasticLoadBalancing",
	"CloudFormation",
	"CodeDeploy",
	"CodeBuild",
	"CodePipeline",
	"SecretsManager",
	"Systems",
	"Service",
	"Simple",
	"Translate",
	"Comprehend",
	"Rekognition",
	"Polly",
	"Textract",
	"Transcribe",
}

// isAWSManagedPolicyName checks if a policy name looks like an AWS managed policy.
func isAWSManagedPolicyName(name string) bool {
	for _, prefix := range awsManagedPolicyPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	// Also check for common patterns
	commonAWSPolicies := []string{
		"AdministratorAccess",
		"PowerUserAccess",
		"ReadOnlyAccess",
		"ViewOnlyAccess",
		"SecurityAudit",
		"SupportUser",
		"Billing",
	}
	for _, policy := range commonAWSPolicies {
		if name == policy {
			return true
		}
	}
	return false
}

// Common AWS managed policy ARN generators.

// AWSLambdaBasicExecutionRole returns the ARN for AWSLambdaBasicExecutionRole.
func AWSLambdaBasicExecutionRole(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole", partition)
}

// AWSLambdaVPCAccessExecutionRole returns the ARN for AWSLambdaVPCAccessExecutionRole.
func AWSLambdaVPCAccessExecutionRole(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole", partition)
}

// AWSLambdaDynamoDBExecutionRole returns the ARN for AWSLambdaDynamoDBExecutionRole.
func AWSLambdaDynamoDBExecutionRole(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/service-role/AWSLambdaDynamoDBExecutionRole", partition)
}

// AWSLambdaKinesisExecutionRole returns the ARN for AWSLambdaKinesisExecutionRole.
func AWSLambdaKinesisExecutionRole(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/service-role/AWSLambdaKinesisExecutionRole", partition)
}

// AWSLambdaSQSQueueExecutionRole returns the ARN for AWSLambdaSQSQueueExecutionRole.
func AWSLambdaSQSQueueExecutionRole(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole", partition)
}

// AWSXrayWriteOnlyAccess returns the ARN for AWSXrayWriteOnlyAccess.
func AWSXrayWriteOnlyAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AWSXrayWriteOnlyAccess", partition)
}

// AWSStepFunctionsFullAccess returns the ARN for AWSStepFunctionsFullAccess.
func AWSStepFunctionsFullAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AWSStepFunctionsFullAccess", partition)
}

// CloudWatchLogsFullAccess returns the ARN for CloudWatchLogsFullAccess.
func CloudWatchLogsFullAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/CloudWatchLogsFullAccess", partition)
}

// AmazonDynamoDBFullAccess returns the ARN for AmazonDynamoDBFullAccess.
func AmazonDynamoDBFullAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AmazonDynamoDBFullAccess", partition)
}

// AmazonDynamoDBReadOnlyAccess returns the ARN for AmazonDynamoDBReadOnlyAccess.
func AmazonDynamoDBReadOnlyAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AmazonDynamoDBReadOnlyAccess", partition)
}

// AmazonS3FullAccess returns the ARN for AmazonS3FullAccess.
func AmazonS3FullAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AmazonS3FullAccess", partition)
}

// AmazonS3ReadOnlyAccess returns the ARN for AmazonS3ReadOnlyAccess.
func AmazonS3ReadOnlyAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AmazonS3ReadOnlyAccess", partition)
}

// AmazonSNSFullAccess returns the ARN for AmazonSNSFullAccess.
func AmazonSNSFullAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AmazonSNSFullAccess", partition)
}

// AmazonSQSFullAccess returns the ARN for AmazonSQSFullAccess.
func AmazonSQSFullAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AmazonSQSFullAccess", partition)
}

// AmazonVPCFullAccess returns the ARN for AmazonVPCFullAccess.
func AmazonVPCFullAccess(partition string) string {
	if partition == "" {
		partition = "aws"
	}
	return fmt.Sprintf("arn:%s:iam::aws:policy/AmazonVPCFullAccess", partition)
}
