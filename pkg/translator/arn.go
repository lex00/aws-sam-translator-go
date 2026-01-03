// Package translator provides the main SAM to CloudFormation transformation orchestrator.
package translator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/region"
)

// ARN represents a parsed AWS ARN.
type ARN struct {
	Partition string
	Service   string
	Region    string
	AccountID string
	Resource  string
}

// String returns the ARN as a string.
func (a ARN) String() string {
	return fmt.Sprintf("arn:%s:%s:%s:%s:%s",
		a.Partition,
		a.Service,
		a.Region,
		a.AccountID,
		a.Resource,
	)
}

// arnPattern is the regex pattern for parsing ARNs.
var arnPattern = regexp.MustCompile(`^arn:([^:]+):([^:]+):([^:]*):([^:]*):(.+)$`)

// ArnGenerator generates AWS ARNs for various resource types.
type ArnGenerator struct {
	partition string
	region    string
	accountID string
}

// NewArnGenerator creates a new ArnGenerator with default partition (aws).
func NewArnGenerator(regionStr, accountID string) *ArnGenerator {
	partition := region.GetArnPartition(regionStr)
	return &ArnGenerator{
		partition: partition,
		region:    regionStr,
		accountID: accountID,
	}
}

// NewArnGeneratorWithPartition creates an ArnGenerator with a specific partition.
func NewArnGeneratorWithPartition(partition, regionStr, accountID string) *ArnGenerator {
	return &ArnGenerator{
		partition: partition,
		region:    regionStr,
		accountID: accountID,
	}
}

// Lambda generates an ARN for a Lambda function.
func (g *ArnGenerator) Lambda(functionName string) string {
	return g.build("lambda", g.region, g.accountID, "function:"+functionName)
}

// LambdaAlias generates an ARN for a Lambda function alias.
func (g *ArnGenerator) LambdaAlias(functionName, alias string) string {
	return g.build("lambda", g.region, g.accountID, "function:"+functionName+":"+alias)
}

// LambdaVersion generates an ARN for a Lambda function version.
func (g *ArnGenerator) LambdaVersion(functionName, version string) string {
	return g.build("lambda", g.region, g.accountID, "function:"+functionName+":"+version)
}

// LambdaLayer generates an ARN for a Lambda layer.
func (g *ArnGenerator) LambdaLayer(layerName string) string {
	return g.build("lambda", g.region, g.accountID, "layer:"+layerName)
}

// LambdaLayerVersion generates an ARN for a Lambda layer version.
func (g *ArnGenerator) LambdaLayerVersion(layerName string, version int) string {
	return g.build("lambda", g.region, g.accountID, fmt.Sprintf("layer:%s:%d", layerName, version))
}

// APIGateway generates an ARN for an API Gateway REST API.
func (g *ArnGenerator) APIGateway(apiID string) string {
	return g.build("apigateway", g.region, "", "/restapis/"+apiID)
}

// APIGatewayStage generates an ARN for an API Gateway stage.
func (g *ArnGenerator) APIGatewayStage(apiID, stageName string) string {
	return g.build("apigateway", g.region, "", fmt.Sprintf("/restapis/%s/stages/%s", apiID, stageName))
}

// APIGatewayV2 generates an ARN for an API Gateway V2 (HTTP) API.
func (g *ArnGenerator) APIGatewayV2(apiID string) string {
	return g.build("apigateway", g.region, "", "/apis/"+apiID)
}

// APIGatewayExecute generates an ARN for API Gateway execution.
func (g *ArnGenerator) APIGatewayExecute(apiID, stageName, method, path string) string {
	resource := fmt.Sprintf("%s/%s/%s%s", apiID, stageName, method, path)
	return g.buildExecuteAPI("execute-api", g.region, g.accountID, resource)
}

// IAMRole generates an ARN for an IAM role.
func (g *ArnGenerator) IAMRole(roleName string) string {
	return g.build("iam", "", g.accountID, "role/"+roleName)
}

// IAMRoleWithPath generates an ARN for an IAM role with a path.
func (g *ArnGenerator) IAMRoleWithPath(path, roleName string) string {
	if path == "" || path == "/" {
		return g.IAMRole(roleName)
	}
	// Ensure path starts and ends with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return g.build("iam", "", g.accountID, "role"+path+roleName)
}

// IAMPolicy generates an ARN for an IAM policy.
func (g *ArnGenerator) IAMPolicy(policyName string) string {
	return g.build("iam", "", g.accountID, "policy/"+policyName)
}

// IAMManagedPolicy generates an ARN for an AWS managed policy.
func (g *ArnGenerator) IAMManagedPolicy(policyName string) string {
	return fmt.Sprintf("arn:%s:iam::aws:policy/%s", g.partition, policyName)
}

// S3Bucket generates an ARN for an S3 bucket.
func (g *ArnGenerator) S3Bucket(bucketName string) string {
	return g.build("s3", "", "", bucketName)
}

// S3Object generates an ARN for an S3 object.
func (g *ArnGenerator) S3Object(bucketName, key string) string {
	return g.build("s3", "", "", bucketName+"/"+key)
}

// DynamoDBTable generates an ARN for a DynamoDB table.
func (g *ArnGenerator) DynamoDBTable(tableName string) string {
	return g.build("dynamodb", g.region, g.accountID, "table/"+tableName)
}

// DynamoDBIndex generates an ARN for a DynamoDB index.
func (g *ArnGenerator) DynamoDBIndex(tableName, indexName string) string {
	return g.build("dynamodb", g.region, g.accountID, "table/"+tableName+"/index/"+indexName)
}

// DynamoDBStream generates an ARN for a DynamoDB stream.
func (g *ArnGenerator) DynamoDBStream(tableName, streamLabel string) string {
	return g.build("dynamodb", g.region, g.accountID, "table/"+tableName+"/stream/"+streamLabel)
}

// SNSTopic generates an ARN for an SNS topic.
func (g *ArnGenerator) SNSTopic(topicName string) string {
	return g.build("sns", g.region, g.accountID, topicName)
}

// SQSQueue generates an ARN for an SQS queue.
func (g *ArnGenerator) SQSQueue(queueName string) string {
	return g.build("sqs", g.region, g.accountID, queueName)
}

// KinesisStream generates an ARN for a Kinesis stream.
func (g *ArnGenerator) KinesisStream(streamName string) string {
	return g.build("kinesis", g.region, g.accountID, "stream/"+streamName)
}

// StepFunctionsStateMachine generates an ARN for a Step Functions state machine.
func (g *ArnGenerator) StepFunctionsStateMachine(stateMachineName string) string {
	return g.build("states", g.region, g.accountID, "stateMachine:"+stateMachineName)
}

// EventsRule generates an ARN for an EventBridge rule.
func (g *ArnGenerator) EventsRule(ruleName string) string {
	return g.build("events", g.region, g.accountID, "rule/"+ruleName)
}

// EventsEventBus generates an ARN for an EventBridge event bus.
func (g *ArnGenerator) EventsEventBus(eventBusName string) string {
	return g.build("events", g.region, g.accountID, "event-bus/"+eventBusName)
}

// CloudWatchLogGroup generates an ARN for a CloudWatch log group.
func (g *ArnGenerator) CloudWatchLogGroup(logGroupName string) string {
	return g.build("logs", g.region, g.accountID, "log-group:"+logGroupName)
}

// CloudWatchAlarm generates an ARN for a CloudWatch alarm.
func (g *ArnGenerator) CloudWatchAlarm(alarmName string) string {
	return g.build("cloudwatch", g.region, g.accountID, "alarm:"+alarmName)
}

// SecretsManagerSecret generates an ARN for a Secrets Manager secret.
func (g *ArnGenerator) SecretsManagerSecret(secretName string) string {
	return g.build("secretsmanager", g.region, g.accountID, "secret:"+secretName)
}

// KMSKey generates an ARN for a KMS key.
func (g *ArnGenerator) KMSKey(keyID string) string {
	return g.build("kms", g.region, g.accountID, "key/"+keyID)
}

// KMSAlias generates an ARN for a KMS alias.
func (g *ArnGenerator) KMSAlias(aliasName string) string {
	alias := aliasName
	if !strings.HasPrefix(alias, "alias/") {
		alias = "alias/" + alias
	}
	return g.build("kms", g.region, g.accountID, alias)
}

// CognitoUserPool generates an ARN for a Cognito user pool.
func (g *ArnGenerator) CognitoUserPool(userPoolID string) string {
	return g.build("cognito-idp", g.region, g.accountID, "userpool/"+userPoolID)
}

// CodeDeployApplication generates an ARN for a CodeDeploy application.
func (g *ArnGenerator) CodeDeployApplication(applicationName string) string {
	return g.build("codedeploy", g.region, g.accountID, "application:"+applicationName)
}

// CodeDeployDeploymentGroup generates an ARN for a CodeDeploy deployment group.
func (g *ArnGenerator) CodeDeployDeploymentGroup(applicationName, deploymentGroupName string) string {
	return g.build("codedeploy", g.region, g.accountID, "deploymentgroup:"+applicationName+"/"+deploymentGroupName)
}

// Generic generates an ARN for any service with custom resource specification.
func (g *ArnGenerator) Generic(service, resource string) string {
	return g.build(service, g.region, g.accountID, resource)
}

// GenericGlobal generates an ARN for a global (non-regional) service.
func (g *ArnGenerator) GenericGlobal(service, resource string) string {
	return g.build(service, "", g.accountID, resource)
}

// GenericNoAccount generates an ARN for a service that doesn't use account ID.
func (g *ArnGenerator) GenericNoAccount(service, resource string) string {
	return g.build(service, g.region, "", resource)
}

// build constructs an ARN string.
func (g *ArnGenerator) build(service, arnRegion, accountID, resource string) string {
	return fmt.Sprintf("arn:%s:%s:%s:%s:%s",
		g.partition,
		service,
		arnRegion,
		accountID,
		resource,
	)
}

// buildExecuteAPI constructs an ARN for execute-api (API Gateway invocations).
func (g *ArnGenerator) buildExecuteAPI(service, arnRegion, accountID, resource string) string {
	return fmt.Sprintf("arn:%s:%s:%s:%s:%s",
		g.partition,
		service,
		arnRegion,
		accountID,
		resource,
	)
}

// ParseARN parses an ARN string into its components.
func ParseARN(arnStr string) (*ARN, error) {
	if arnStr == "" {
		return nil, fmt.Errorf("ARN cannot be empty")
	}

	matches := arnPattern.FindStringSubmatch(arnStr)
	if matches == nil {
		return nil, fmt.Errorf("invalid ARN format: %s", arnStr)
	}

	return &ARN{
		Partition: matches[1],
		Service:   matches[2],
		Region:    matches[3],
		AccountID: matches[4],
		Resource:  matches[5],
	}, nil
}

// IsValidARN checks if a string is a valid ARN format.
func IsValidARN(arnStr string) bool {
	_, err := ParseARN(arnStr)
	return err == nil
}

// GetPartition returns the partition from an ARN string.
func GetPartition(arnStr string) (string, error) {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return "", err
	}
	return arn.Partition, nil
}

// GetService returns the service from an ARN string.
func GetService(arnStr string) (string, error) {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return "", err
	}
	return arn.Service, nil
}

// GetResource returns the resource from an ARN string.
func GetResource(arnStr string) (string, error) {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return "", err
	}
	return arn.Resource, nil
}

// ReplacePartition creates a new ARN with a different partition.
func ReplacePartition(arnStr string, newPartition string) (string, error) {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return "", err
	}
	arn.Partition = newPartition
	return arn.String(), nil
}

// ReplaceRegion creates a new ARN with a different region.
func ReplaceRegion(arnStr string, newRegion string) (string, error) {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return "", err
	}
	arn.Region = newRegion
	return arn.String(), nil
}

// ReplaceAccountID creates a new ARN with a different account ID.
func ReplaceAccountID(arnStr string, newAccountID string) (string, error) {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return "", err
	}
	arn.AccountID = newAccountID
	return arn.String(), nil
}
