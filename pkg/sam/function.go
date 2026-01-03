// Package sam provides SAM resource transformers.
package sam

import (
	"fmt"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/model/iam"
	"github.com/lex00/aws-sam-translator-go/pkg/model/lambda"
)

// Function represents an AWS::Serverless::Function resource.
// See: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-resource-function.html
type Function struct {
	// Handler is the function entry point (required for Zip package type).
	Handler string `json:"Handler,omitempty" yaml:"Handler,omitempty"`

	// Runtime is the runtime environment for the function (required for Zip package type).
	Runtime string `json:"Runtime,omitempty" yaml:"Runtime,omitempty"`

	// CodeUri specifies the function code location.
	// Can be a string (s3://bucket/key) or an object with Bucket, Key, Version properties.
	CodeUri interface{} `json:"CodeUri,omitempty" yaml:"CodeUri,omitempty"`

	// ImageUri is the URI of a container image in Amazon ECR.
	ImageUri interface{} `json:"ImageUri,omitempty" yaml:"ImageUri,omitempty"`

	// PackageType is the deployment package type: Zip or Image.
	PackageType string `json:"PackageType,omitempty" yaml:"PackageType,omitempty"`

	// Description is a description of the function.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// MemorySize is the amount of memory available to the function (MB).
	MemorySize int `json:"MemorySize,omitempty" yaml:"MemorySize,omitempty"`

	// Timeout is the amount of time Lambda allows a function to run (seconds).
	Timeout int `json:"Timeout,omitempty" yaml:"Timeout,omitempty"`

	// Role is the ARN of the function's execution role.
	// If not specified, SAM creates a role automatically.
	Role interface{} `json:"Role,omitempty" yaml:"Role,omitempty"`

	// Policies specifies policies to attach to the function's execution role.
	// Can be a string (managed policy ARN), array of ARNs, or inline policy document.
	Policies interface{} `json:"Policies,omitempty" yaml:"Policies,omitempty"`

	// Environment contains environment variables for the function.
	Environment map[string]interface{} `json:"Environment,omitempty" yaml:"Environment,omitempty"`

	// Events defines the event sources that trigger the function.
	Events map[string]interface{} `json:"Events,omitempty" yaml:"Events,omitempty"`

	// Tags is a map of key-value pairs to apply to the function.
	Tags map[string]string `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// Layers is a list of Lambda layer ARNs to attach to the function.
	Layers []interface{} `json:"Layers,omitempty" yaml:"Layers,omitempty"`

	// VpcConfig connects the function to a VPC.
	VpcConfig map[string]interface{} `json:"VpcConfig,omitempty" yaml:"VpcConfig,omitempty"`

	// FunctionName is the name of the Lambda function.
	FunctionName interface{} `json:"FunctionName,omitempty" yaml:"FunctionName,omitempty"`

	// Architectures is the instruction set architecture (x86_64 or arm64).
	Architectures []string `json:"Architectures,omitempty" yaml:"Architectures,omitempty"`

	// AutoPublishAlias creates a Lambda version and alias with the specified name.
	AutoPublishAlias string `json:"AutoPublishAlias,omitempty" yaml:"AutoPublishAlias,omitempty"`

	// AutoPublishCodeSha256 is the SHA256 hash to trigger auto-publishing.
	AutoPublishCodeSha256 string `json:"AutoPublishCodeSha256,omitempty" yaml:"AutoPublishCodeSha256,omitempty"`

	// DeploymentPreference configures CodeDeploy deployment for gradual rollouts.
	DeploymentPreference map[string]interface{} `json:"DeploymentPreference,omitempty" yaml:"DeploymentPreference,omitempty"`

	// ProvisionedConcurrencyConfig specifies provisioned concurrency settings.
	ProvisionedConcurrencyConfig map[string]interface{} `json:"ProvisionedConcurrencyConfig,omitempty" yaml:"ProvisionedConcurrencyConfig,omitempty"`

	// ReservedConcurrentExecutions is the number of reserved concurrent executions.
	ReservedConcurrentExecutions *int `json:"ReservedConcurrentExecutions,omitempty" yaml:"ReservedConcurrentExecutions,omitempty"`

	// Tracing configures AWS X-Ray tracing. Valid values: Active, PassThrough.
	Tracing string `json:"Tracing,omitempty" yaml:"Tracing,omitempty"`

	// DeadLetterQueue configures the dead letter queue for failed invocations.
	DeadLetterQueue map[string]interface{} `json:"DeadLetterQueue,omitempty" yaml:"DeadLetterQueue,omitempty"`

	// KmsKeyArn is the ARN of the KMS key used to encrypt environment variables.
	KmsKeyArn interface{} `json:"KmsKeyArn,omitempty" yaml:"KmsKeyArn,omitempty"`

	// EphemeralStorage configures the size of the function's /tmp directory.
	EphemeralStorage map[string]interface{} `json:"EphemeralStorage,omitempty" yaml:"EphemeralStorage,omitempty"`

	// SnapStart enables SnapStart for Java functions.
	SnapStart map[string]interface{} `json:"SnapStart,omitempty" yaml:"SnapStart,omitempty"`

	// FileSystemConfigs connects the function to an Amazon EFS file system.
	FileSystemConfigs []map[string]interface{} `json:"FileSystemConfigs,omitempty" yaml:"FileSystemConfigs,omitempty"`

	// ImageConfig overrides the container image settings.
	ImageConfig map[string]interface{} `json:"ImageConfig,omitempty" yaml:"ImageConfig,omitempty"`

	// CodeSigningConfigArn is the ARN of a code-signing configuration.
	CodeSigningConfigArn interface{} `json:"CodeSigningConfigArn,omitempty" yaml:"CodeSigningConfigArn,omitempty"`

	// RuntimeManagementConfig configures runtime management settings.
	RuntimeManagementConfig map[string]interface{} `json:"RuntimeManagementConfig,omitempty" yaml:"RuntimeManagementConfig,omitempty"`

	// PermissionsBoundary is the ARN of a permissions boundary policy.
	PermissionsBoundary interface{} `json:"PermissionsBoundary,omitempty" yaml:"PermissionsBoundary,omitempty"`

	// FunctionUrlConfig configures a Lambda function URL.
	FunctionUrlConfig map[string]interface{} `json:"FunctionUrlConfig,omitempty" yaml:"FunctionUrlConfig,omitempty"`

	// LoggingConfig configures CloudWatch logging settings.
	LoggingConfig map[string]interface{} `json:"LoggingConfig,omitempty" yaml:"LoggingConfig,omitempty"`

	// RecursiveLoop sets loop detection behavior for recursive invocations.
	RecursiveLoop string `json:"RecursiveLoop,omitempty" yaml:"RecursiveLoop,omitempty"`

	// Condition is a CloudFormation condition name for the function.
	Condition string `json:"Condition,omitempty" yaml:"Condition,omitempty"`

	// DependsOn specifies resource dependencies.
	DependsOn interface{} `json:"DependsOn,omitempty" yaml:"DependsOn,omitempty"`

	// Metadata is custom metadata for the resource.
	Metadata map[string]interface{} `json:"Metadata,omitempty" yaml:"Metadata,omitempty"`
}

// TransformContext provides context information for the transformation.
type TransformContext struct {
	// Region is the AWS region.
	Region string

	// AccountID is the AWS account ID.
	AccountID string

	// StackName is the CloudFormation stack name.
	StackName string

	// Partition is the AWS partition (aws, aws-cn, aws-us-gov).
	Partition string
}

// FunctionTransformer transforms AWS::Serverless::Function to CloudFormation.
type FunctionTransformer struct{}

// NewFunctionTransformer creates a new FunctionTransformer.
func NewFunctionTransformer() *FunctionTransformer {
	return &FunctionTransformer{}
}

// Transform converts a SAM Function to CloudFormation resources.
func (t *FunctionTransformer) Transform(logicalID string, f *Function, ctx *TransformContext) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Build the Lambda function properties
	functionProps, err := t.buildFunctionProperties(logicalID, f)
	if err != nil {
		return nil, fmt.Errorf("failed to build function properties: %w", err)
	}

	// Determine role configuration
	roleRef, roleResource, err := t.buildRole(logicalID, f)
	if err != nil {
		return nil, fmt.Errorf("failed to build role: %w", err)
	}

	functionProps["Role"] = roleRef
	if roleResource != nil {
		resources[logicalID+"Role"] = roleResource
	}

	// Build the function resource
	functionResource := map[string]interface{}{
		"Type":       "AWS::Lambda::Function",
		"Properties": functionProps,
	}

	if f.Condition != "" {
		functionResource["Condition"] = f.Condition
	}
	if f.DependsOn != nil {
		functionResource["DependsOn"] = f.DependsOn
	}
	if f.Metadata != nil {
		functionResource["Metadata"] = f.Metadata
	}

	resources[logicalID] = functionResource

	// Handle AutoPublishAlias (versioning)
	if f.AutoPublishAlias != "" {
		versionResources, err := t.buildVersionAndAlias(logicalID, f)
		if err != nil {
			return nil, fmt.Errorf("failed to build version and alias: %w", err)
		}
		for k, v := range versionResources {
			resources[k] = v
		}
	}

	// Handle events
	if len(f.Events) > 0 {
		eventResources, err := t.buildEventResources(logicalID, f)
		if err != nil {
			return nil, fmt.Errorf("failed to build event resources: %w", err)
		}
		for k, v := range eventResources {
			resources[k] = v
		}
	}

	// Handle DeploymentPreference
	if f.DeploymentPreference != nil && f.AutoPublishAlias != "" {
		deployResources, err := t.buildDeploymentPreference(logicalID, f)
		if err != nil {
			return nil, fmt.Errorf("failed to build deployment preference: %w", err)
		}
		for k, v := range deployResources {
			resources[k] = v
		}
	}

	return resources, nil
}

// buildFunctionProperties builds the Lambda function properties.
func (t *FunctionTransformer) buildFunctionProperties(logicalID string, f *Function) (map[string]interface{}, error) {
	props := make(map[string]interface{})

	// Handler (required for Zip)
	if f.Handler != "" {
		props["Handler"] = f.Handler
	}

	// Runtime (required for Zip)
	if f.Runtime != "" {
		props["Runtime"] = f.Runtime
	}

	// Code configuration
	code, err := t.buildCodeConfig(f)
	if err != nil {
		return nil, err
	}
	props["Code"] = code

	// Optional properties
	if f.FunctionName != nil {
		props["FunctionName"] = f.FunctionName
	}

	if f.Description != "" {
		props["Description"] = f.Description
	}

	if f.MemorySize > 0 {
		props["MemorySize"] = f.MemorySize
	}

	if f.Timeout > 0 {
		props["Timeout"] = f.Timeout
	}

	if f.PackageType != "" {
		props["PackageType"] = f.PackageType
	}

	if len(f.Architectures) > 0 {
		props["Architectures"] = f.Architectures
	}

	if f.Environment != nil {
		props["Environment"] = f.Environment
	}

	if len(f.Tags) > 0 {
		tags := make([]interface{}, 0, len(f.Tags))
		for k, v := range f.Tags {
			tags = append(tags, map[string]interface{}{
				"Key":   k,
				"Value": v,
			})
		}
		props["Tags"] = tags
	}

	if len(f.Layers) > 0 {
		props["Layers"] = f.Layers
	}

	if f.VpcConfig != nil {
		props["VpcConfig"] = f.VpcConfig
	}

	if f.ReservedConcurrentExecutions != nil {
		props["ReservedConcurrentExecutions"] = *f.ReservedConcurrentExecutions
	}

	if f.Tracing != "" {
		props["TracingConfig"] = map[string]interface{}{
			"Mode": f.Tracing,
		}
	}

	if f.DeadLetterQueue != nil {
		props["DeadLetterConfig"] = f.DeadLetterQueue
	}

	if f.KmsKeyArn != nil {
		props["KmsKeyArn"] = f.KmsKeyArn
	}

	if f.EphemeralStorage != nil {
		props["EphemeralStorage"] = f.EphemeralStorage
	}

	if f.SnapStart != nil {
		props["SnapStart"] = f.SnapStart
	}

	if len(f.FileSystemConfigs) > 0 {
		props["FileSystemConfigs"] = f.FileSystemConfigs
	}

	if f.ImageConfig != nil {
		props["ImageConfig"] = f.ImageConfig
	}

	if f.CodeSigningConfigArn != nil {
		props["CodeSigningConfigArn"] = f.CodeSigningConfigArn
	}

	if f.RuntimeManagementConfig != nil {
		props["RuntimeManagementConfig"] = f.RuntimeManagementConfig
	}

	if f.LoggingConfig != nil {
		props["LoggingConfig"] = f.LoggingConfig
	}

	if f.RecursiveLoop != "" {
		props["RecursiveLoop"] = f.RecursiveLoop
	}

	return props, nil
}

// buildCodeConfig builds the Code property from CodeUri or ImageUri.
func (t *FunctionTransformer) buildCodeConfig(f *Function) (map[string]interface{}, error) {
	code := make(map[string]interface{})

	if f.ImageUri != nil {
		code["ImageUri"] = f.ImageUri
		return code, nil
	}

	if f.CodeUri == nil {
		return nil, fmt.Errorf("CodeUri or ImageUri is required")
	}

	switch v := f.CodeUri.(type) {
	case string:
		// Parse s3://bucket/key format
		if strings.HasPrefix(v, "s3://") {
			parsed, err := parseS3Uri(v)
			if err != nil {
				return nil, err
			}
			code["S3Bucket"] = parsed["S3Bucket"]
			code["S3Key"] = parsed["S3Key"]
			if version, ok := parsed["S3ObjectVersion"]; ok {
				code["S3ObjectVersion"] = version
			}
		} else {
			// Local path - keep as-is for packaging
			code["S3Bucket"] = v
		}
	case map[string]interface{}:
		if bucket, ok := v["Bucket"]; ok {
			code["S3Bucket"] = bucket
		}
		if key, ok := v["Key"]; ok {
			code["S3Key"] = key
		}
		if version, ok := v["Version"]; ok {
			code["S3ObjectVersion"] = version
		}
	case map[interface{}]interface{}:
		// Convert YAML map
		if bucket, ok := v["Bucket"]; ok {
			code["S3Bucket"] = bucket
		}
		if key, ok := v["Key"]; ok {
			code["S3Key"] = key
		}
		if version, ok := v["Version"]; ok {
			code["S3ObjectVersion"] = version
		}
	default:
		code["S3Bucket"] = v
	}

	return code, nil
}

// parseS3Uri parses an S3 URI string (s3://bucket/key) into components.
func parseS3Uri(uri string) (map[string]interface{}, error) {
	if !strings.HasPrefix(uri, "s3://") {
		return nil, fmt.Errorf("invalid S3 URI: must start with s3://")
	}

	path := strings.TrimPrefix(uri, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid S3 URI: must have bucket and key")
	}

	return map[string]interface{}{
		"S3Bucket": parts[0],
		"S3Key":    parts[1],
	}, nil
}

// buildRole builds the IAM role for the function.
func (t *FunctionTransformer) buildRole(logicalID string, f *Function) (interface{}, map[string]interface{}, error) {
	// If Role is explicitly provided, use it
	if f.Role != nil {
		return f.Role, nil, nil
	}

	// Build an execution role
	trustPolicy := iam.NewAssumeRolePolicyForService(iam.ServiceLambda)
	role := iam.NewRole(trustPolicy)

	// Add basic execution role policy
	managedPolicies := []interface{}{
		"arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
	}

	// Add VPC access policy if VPC is configured
	if f.VpcConfig != nil {
		managedPolicies = append(managedPolicies,
			"arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole")
	}

	// Add X-Ray policy if tracing is enabled
	if f.Tracing == "Active" {
		managedPolicies = append(managedPolicies,
			"arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess")
	}

	// Process Policies property
	if f.Policies != nil {
		additionalPolicies, inlinePolicies, err := t.processPolicies(logicalID, f.Policies)
		if err != nil {
			return nil, nil, err
		}
		managedPolicies = append(managedPolicies, additionalPolicies...)
		role.Policies = append(role.Policies, inlinePolicies...)
	}

	// Set permissions boundary
	if f.PermissionsBoundary != nil {
		role.PermissionsBoundary = f.PermissionsBoundary
	}

	// Build role properties
	roleProps := role.ToCloudFormation()
	roleProps["ManagedPolicyArns"] = managedPolicies

	roleResource := map[string]interface{}{
		"Type":       "AWS::IAM::Role",
		"Properties": roleProps,
	}

	// Return reference to the role
	roleRef := map[string]interface{}{
		"Fn::GetAtt": []string{logicalID + "Role", "Arn"},
	}

	return roleRef, roleResource, nil
}

// processPolicies processes the Policies property and returns managed policy ARNs and inline policies.
func (t *FunctionTransformer) processPolicies(logicalID string, policies interface{}) ([]interface{}, []iam.InlinePolicy, error) {
	var managedPolicies []interface{}
	var inlinePolicies []iam.InlinePolicy

	switch p := policies.(type) {
	case string:
		// Single managed policy ARN
		managedPolicies = append(managedPolicies, p)

	case []interface{}:
		// Array of policies
		for _, item := range p {
			switch v := item.(type) {
			case string:
				managedPolicies = append(managedPolicies, v)
			case map[string]interface{}:
				// Could be an inline policy or a SAM policy template
				if _, hasStatement := v["Statement"]; hasStatement {
					// Inline policy document
					doc := iam.NewPolicyDocument()
					if statements, ok := v["Statement"].([]interface{}); ok {
						for _, stmt := range statements {
							if stmtMap, ok := stmt.(map[string]interface{}); ok {
								statement := t.mapToStatement(stmtMap)
								doc.AddStatement(statement)
							}
						}
					}
					inlinePolicies = append(inlinePolicies, iam.InlinePolicy{
						PolicyName:     logicalID + "Policy",
						PolicyDocument: doc,
					})
				} else {
					// Treat as SAM policy template (already expanded by plugin)
					// or unknown format - add as-is
					managedPolicies = append(managedPolicies, v)
				}
			}
		}

	case map[string]interface{}:
		// Single inline policy document
		if _, hasStatement := p["Statement"]; hasStatement {
			doc := iam.NewPolicyDocument()
			if statements, ok := p["Statement"].([]interface{}); ok {
				for _, stmt := range statements {
					if stmtMap, ok := stmt.(map[string]interface{}); ok {
						statement := t.mapToStatement(stmtMap)
						doc.AddStatement(statement)
					}
				}
			}
			inlinePolicies = append(inlinePolicies, iam.InlinePolicy{
				PolicyName:     logicalID + "Policy",
				PolicyDocument: doc,
			})
		}
	}

	return managedPolicies, inlinePolicies, nil
}

// mapToStatement converts a map to an IAM Statement.
func (t *FunctionTransformer) mapToStatement(m map[string]interface{}) *iam.Statement {
	stmt := iam.NewStatement(iam.EffectAllow)

	if effect, ok := m["Effect"].(string); ok {
		stmt.Effect = effect
	}
	if action, ok := m["Action"]; ok {
		stmt.Action = action
	}
	if resource, ok := m["Resource"]; ok {
		stmt.Resource = resource
	}
	if condition, ok := m["Condition"].(map[string]interface{}); ok {
		stmt.Condition = condition
	}
	if sid, ok := m["Sid"].(string); ok {
		stmt.Sid = sid
	}

	return stmt
}

// buildVersionAndAlias creates Lambda Version and Alias resources.
func (t *FunctionTransformer) buildVersionAndAlias(logicalID string, f *Function) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Version
	versionID := logicalID + "Version"
	versionProps := map[string]interface{}{
		"FunctionName": map[string]interface{}{"Ref": logicalID},
	}
	if f.AutoPublishCodeSha256 != "" {
		versionProps["CodeSha256"] = f.AutoPublishCodeSha256
	}

	versionResource := map[string]interface{}{
		"Type":           lambda.ResourceTypeVersion,
		"DeletionPolicy": "Retain",
		"Properties":     versionProps,
	}
	resources[versionID] = versionResource

	// Create Alias
	aliasID := logicalID + "Alias" + f.AutoPublishAlias
	aliasProps := map[string]interface{}{
		"FunctionName":    map[string]interface{}{"Ref": logicalID},
		"FunctionVersion": map[string]interface{}{"Fn::GetAtt": []string{versionID, "Version"}},
		"Name":            f.AutoPublishAlias,
	}

	// Add provisioned concurrency if specified
	if f.ProvisionedConcurrencyConfig != nil {
		aliasProps["ProvisionedConcurrencyConfig"] = f.ProvisionedConcurrencyConfig
	}

	aliasResource := map[string]interface{}{
		"Type":       lambda.ResourceTypeAlias,
		"Properties": aliasProps,
	}
	resources[aliasID] = aliasResource

	return resources, nil
}

// buildEventResources creates resources for function event sources.
func (t *FunctionTransformer) buildEventResources(logicalID string, f *Function) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Reference to the function (or alias if AutoPublishAlias is set)
	var functionRef interface{}
	if f.AutoPublishAlias != "" {
		functionRef = map[string]interface{}{
			"Ref": logicalID + "Alias" + f.AutoPublishAlias,
		}
	} else {
		functionRef = map[string]interface{}{
			"Fn::GetAtt": []string{logicalID, "Arn"},
		}
	}

	for eventName, eventConfig := range f.Events {
		eventMap, ok := eventConfig.(map[string]interface{})
		if !ok {
			continue
		}

		eventType, _ := eventMap["Type"].(string)
		eventProps, _ := eventMap["Properties"].(map[string]interface{})

		eventResources, err := t.buildEventSource(logicalID, eventName, eventType, eventProps, functionRef)
		if err != nil {
			return nil, fmt.Errorf("failed to build event %s: %w", eventName, err)
		}

		for k, v := range eventResources {
			resources[k] = v
		}
	}

	return resources, nil
}

// buildEventSource creates resources for a single event source.
func (t *FunctionTransformer) buildEventSource(logicalID, eventName, eventType string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	switch eventType {
	case "S3":
		return t.buildS3Event(logicalID, eventName, props, functionRef)
	case "SQS":
		return t.buildSQSEvent(logicalID, eventName, props, functionRef)
	case "Kinesis":
		return t.buildKinesisEvent(logicalID, eventName, props, functionRef)
	case "DynamoDB":
		return t.buildDynamoDBEvent(logicalID, eventName, props, functionRef)
	case "Api":
		return t.buildApiEvent(logicalID, eventName, props, functionRef)
	case "HttpApi":
		return t.buildHttpApiEvent(logicalID, eventName, props, functionRef)
	case "Schedule":
		return t.buildScheduleEvent(logicalID, eventName, props, functionRef)
	case "CloudWatchEvent", "EventBridgeRule":
		return t.buildCloudWatchEvent(logicalID, eventName, props, functionRef)
	case "SNS":
		return t.buildSNSEvent(logicalID, eventName, props, functionRef)
	case "IoTRule":
		return t.buildIoTRuleEvent(logicalID, eventName, props, functionRef)
	case "Cognito":
		return t.buildCognitoEvent(logicalID, eventName, props, functionRef)
	case "MSK":
		return t.buildMSKEvent(logicalID, eventName, props, functionRef)
	case "MQ":
		return t.buildMQEvent(logicalID, eventName, props, functionRef)
	case "SelfManagedKafka":
		return t.buildSelfManagedKafkaEvent(logicalID, eventName, props, functionRef)
	case "CloudWatchLogs":
		return t.buildCloudWatchLogsEvent(logicalID, eventName, props, functionRef)
	case "AlexaSkill":
		return t.buildAlexaSkillEvent(logicalID, eventName, props, functionRef)
	default:
		// Unknown event type - skip
		return resources, nil
	}
}

// buildS3Event creates resources for an S3 event source.
func (t *FunctionTransformer) buildS3Event(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Lambda permission
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":        "lambda:InvokeFunction",
		"FunctionName":  functionRef,
		"Principal":     "s3.amazonaws.com",
		"SourceAccount": map[string]interface{}{"Ref": "AWS::AccountId"},
	}
	if bucket, ok := props["Bucket"]; ok {
		permissionProps["SourceArn"] = t.buildS3BucketArn(bucket)
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildS3BucketArn creates an S3 bucket ARN from various input formats.
func (t *FunctionTransformer) buildS3BucketArn(bucket interface{}) interface{} {
	switch v := bucket.(type) {
	case string:
		if strings.HasPrefix(v, "arn:") {
			return v
		}
		return fmt.Sprintf("arn:aws:s3:::%s", v)
	case map[string]interface{}:
		if _, hasRef := v["Ref"]; hasRef {
			return map[string]interface{}{
				"Fn::Sub": []interface{}{
					"arn:aws:s3:::${Bucket}",
					map[string]interface{}{"Bucket": v},
				},
			}
		}
		if getAtt, hasGetAtt := v["Fn::GetAtt"]; hasGetAtt {
			return map[string]interface{}{
				"Fn::Sub": []interface{}{
					"arn:aws:s3:::${Bucket}",
					map[string]interface{}{"Bucket": map[string]interface{}{"Fn::GetAtt": getAtt}},
				},
			}
		}
	}
	return bucket
}

// buildSQSEvent creates resources for an SQS event source.
func (t *FunctionTransformer) buildSQSEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	esmID := logicalID + eventName
	esmProps := map[string]interface{}{
		"FunctionName": functionRef,
	}

	if queue, ok := props["Queue"]; ok {
		esmProps["EventSourceArn"] = queue
	}
	if batchSize, ok := props["BatchSize"]; ok {
		esmProps["BatchSize"] = batchSize
	}
	if enabled, ok := props["Enabled"]; ok {
		esmProps["Enabled"] = enabled
	}
	if batchingWindow, ok := props["MaximumBatchingWindowInSeconds"]; ok {
		esmProps["MaximumBatchingWindowInSeconds"] = batchingWindow
	}
	if filterCriteria, ok := props["FilterCriteria"]; ok {
		esmProps["FilterCriteria"] = filterCriteria
	}
	if scalingConfig, ok := props["ScalingConfig"]; ok {
		esmProps["ScalingConfig"] = scalingConfig
	}

	resources[esmID] = map[string]interface{}{
		"Type":       lambda.ResourceTypeEventSourceMapping,
		"Properties": esmProps,
	}

	return resources, nil
}

// buildKinesisEvent creates resources for a Kinesis event source.
func (t *FunctionTransformer) buildKinesisEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	esmID := logicalID + eventName
	esmProps := map[string]interface{}{
		"FunctionName": functionRef,
	}

	if stream, ok := props["Stream"]; ok {
		esmProps["EventSourceArn"] = stream
	}
	if startingPosition, ok := props["StartingPosition"]; ok {
		esmProps["StartingPosition"] = startingPosition
	}
	if batchSize, ok := props["BatchSize"]; ok {
		esmProps["BatchSize"] = batchSize
	}
	if enabled, ok := props["Enabled"]; ok {
		esmProps["Enabled"] = enabled
	}
	if batchingWindow, ok := props["MaximumBatchingWindowInSeconds"]; ok {
		esmProps["MaximumBatchingWindowInSeconds"] = batchingWindow
	}
	if bisect, ok := props["BisectBatchOnFunctionError"]; ok {
		esmProps["BisectBatchOnFunctionError"] = bisect
	}
	if maxRetry, ok := props["MaximumRetryAttempts"]; ok {
		esmProps["MaximumRetryAttempts"] = maxRetry
	}
	if maxAge, ok := props["MaximumRecordAgeInSeconds"]; ok {
		esmProps["MaximumRecordAgeInSeconds"] = maxAge
	}
	if parallelization, ok := props["ParallelizationFactor"]; ok {
		esmProps["ParallelizationFactor"] = parallelization
	}
	if destConfig, ok := props["DestinationConfig"]; ok {
		esmProps["DestinationConfig"] = destConfig
	}
	if tumbling, ok := props["TumblingWindowInSeconds"]; ok {
		esmProps["TumblingWindowInSeconds"] = tumbling
	}

	resources[esmID] = map[string]interface{}{
		"Type":       lambda.ResourceTypeEventSourceMapping,
		"Properties": esmProps,
	}

	return resources, nil
}

// buildDynamoDBEvent creates resources for a DynamoDB Streams event source.
func (t *FunctionTransformer) buildDynamoDBEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	esmID := logicalID + eventName
	esmProps := map[string]interface{}{
		"FunctionName": functionRef,
	}

	if stream, ok := props["Stream"]; ok {
		esmProps["EventSourceArn"] = stream
	}
	if startingPosition, ok := props["StartingPosition"]; ok {
		esmProps["StartingPosition"] = startingPosition
	}
	if batchSize, ok := props["BatchSize"]; ok {
		esmProps["BatchSize"] = batchSize
	}
	if enabled, ok := props["Enabled"]; ok {
		esmProps["Enabled"] = enabled
	}
	if batchingWindow, ok := props["MaximumBatchingWindowInSeconds"]; ok {
		esmProps["MaximumBatchingWindowInSeconds"] = batchingWindow
	}
	if bisect, ok := props["BisectBatchOnFunctionError"]; ok {
		esmProps["BisectBatchOnFunctionError"] = bisect
	}
	if maxRetry, ok := props["MaximumRetryAttempts"]; ok {
		esmProps["MaximumRetryAttempts"] = maxRetry
	}
	if maxAge, ok := props["MaximumRecordAgeInSeconds"]; ok {
		esmProps["MaximumRecordAgeInSeconds"] = maxAge
	}
	if parallelization, ok := props["ParallelizationFactor"]; ok {
		esmProps["ParallelizationFactor"] = parallelization
	}
	if destConfig, ok := props["DestinationConfig"]; ok {
		esmProps["DestinationConfig"] = destConfig
	}
	if tumbling, ok := props["TumblingWindowInSeconds"]; ok {
		esmProps["TumblingWindowInSeconds"] = tumbling
	}
	if filterCriteria, ok := props["FilterCriteria"]; ok {
		esmProps["FilterCriteria"] = filterCriteria
	}

	resources[esmID] = map[string]interface{}{
		"Type":       lambda.ResourceTypeEventSourceMapping,
		"Properties": esmProps,
	}

	return resources, nil
}

// buildApiEvent creates resources for an API Gateway event source.
func (t *FunctionTransformer) buildApiEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Lambda permission for API Gateway
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "apigateway.amazonaws.com",
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildHttpApiEvent creates resources for an HTTP API (API Gateway V2) event source.
func (t *FunctionTransformer) buildHttpApiEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Lambda permission for API Gateway V2
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "apigateway.amazonaws.com",
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildScheduleEvent creates resources for a scheduled event (EventBridge).
func (t *FunctionTransformer) buildScheduleEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create EventBridge Rule
	ruleID := logicalID + eventName
	ruleProps := map[string]interface{}{
		"Targets": []interface{}{
			map[string]interface{}{
				"Arn": functionRef,
				"Id":  logicalID,
			},
		},
	}

	if schedule, ok := props["Schedule"]; ok {
		ruleProps["ScheduleExpression"] = schedule
	}
	if name, ok := props["Name"]; ok {
		ruleProps["Name"] = name
	}
	if desc, ok := props["Description"]; ok {
		ruleProps["Description"] = desc
	}
	if state, ok := props["State"]; ok {
		ruleProps["State"] = state
	} else if enabled, ok := props["Enabled"]; ok {
		if enabled == false {
			ruleProps["State"] = "DISABLED"
		}
	}
	if input, ok := props["Input"]; ok {
		targets := ruleProps["Targets"].([]interface{})
		targets[0].(map[string]interface{})["Input"] = input
	}
	if inputPath, ok := props["InputPath"]; ok {
		targets := ruleProps["Targets"].([]interface{})
		targets[0].(map[string]interface{})["InputPath"] = inputPath
	}

	resources[ruleID] = map[string]interface{}{
		"Type":       "AWS::Events::Rule",
		"Properties": ruleProps,
	}

	// Create Lambda permission for EventBridge
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "events.amazonaws.com",
		"SourceArn":    map[string]interface{}{"Fn::GetAtt": []string{ruleID, "Arn"}},
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildCloudWatchEvent creates resources for a CloudWatch/EventBridge event pattern.
func (t *FunctionTransformer) buildCloudWatchEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create EventBridge Rule
	ruleID := logicalID + eventName
	ruleProps := map[string]interface{}{
		"Targets": []interface{}{
			map[string]interface{}{
				"Arn": functionRef,
				"Id":  logicalID,
			},
		},
	}

	if pattern, ok := props["Pattern"]; ok {
		ruleProps["EventPattern"] = pattern
	}
	if name, ok := props["Name"]; ok {
		ruleProps["Name"] = name
	}
	if desc, ok := props["Description"]; ok {
		ruleProps["Description"] = desc
	}
	if state, ok := props["State"]; ok {
		ruleProps["State"] = state
	}
	if input, ok := props["Input"]; ok {
		targets := ruleProps["Targets"].([]interface{})
		targets[0].(map[string]interface{})["Input"] = input
	}
	if inputPath, ok := props["InputPath"]; ok {
		targets := ruleProps["Targets"].([]interface{})
		targets[0].(map[string]interface{})["InputPath"] = inputPath
	}

	resources[ruleID] = map[string]interface{}{
		"Type":       "AWS::Events::Rule",
		"Properties": ruleProps,
	}

	// Create Lambda permission for EventBridge
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "events.amazonaws.com",
		"SourceArn":    map[string]interface{}{"Fn::GetAtt": []string{ruleID, "Arn"}},
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildSNSEvent creates resources for an SNS event source.
func (t *FunctionTransformer) buildSNSEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Lambda permission for SNS
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "sns.amazonaws.com",
	}

	if topic, ok := props["Topic"]; ok {
		permissionProps["SourceArn"] = topic
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	// Create SNS Subscription
	subscriptionID := logicalID + eventName + "Subscription"
	subscriptionProps := map[string]interface{}{
		"Protocol": "lambda",
		"Endpoint": functionRef,
	}

	if topic, ok := props["Topic"]; ok {
		subscriptionProps["TopicArn"] = topic
	}
	if filterPolicy, ok := props["FilterPolicy"]; ok {
		subscriptionProps["FilterPolicy"] = filterPolicy
	}
	if filterPolicyScope, ok := props["FilterPolicyScope"]; ok {
		subscriptionProps["FilterPolicyScope"] = filterPolicyScope
	}

	resources[subscriptionID] = map[string]interface{}{
		"Type":       "AWS::SNS::Subscription",
		"Properties": subscriptionProps,
	}

	return resources, nil
}

// buildIoTRuleEvent creates resources for an IoT Rule event source.
func (t *FunctionTransformer) buildIoTRuleEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create IoT TopicRule
	ruleID := logicalID + eventName
	ruleProps := map[string]interface{}{
		"TopicRulePayload": map[string]interface{}{
			"Actions": []interface{}{
				map[string]interface{}{
					"Lambda": map[string]interface{}{
						"FunctionArn": functionRef,
					},
				},
			},
			"RuleDisabled": false,
		},
	}

	if sql, ok := props["Sql"]; ok {
		ruleProps["TopicRulePayload"].(map[string]interface{})["Sql"] = sql
	}
	if sqlVersion, ok := props["AwsIotSqlVersion"]; ok {
		ruleProps["TopicRulePayload"].(map[string]interface{})["AwsIotSqlVersion"] = sqlVersion
	}
	if desc, ok := props["Description"]; ok {
		ruleProps["TopicRulePayload"].(map[string]interface{})["Description"] = desc
	}

	resources[ruleID] = map[string]interface{}{
		"Type":       "AWS::IoT::TopicRule",
		"Properties": ruleProps,
	}

	// Create Lambda permission for IoT
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "iot.amazonaws.com",
		"SourceArn":    map[string]interface{}{"Fn::GetAtt": []string{ruleID, "Arn"}},
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildCognitoEvent creates resources for a Cognito event source.
func (t *FunctionTransformer) buildCognitoEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Lambda permission for Cognito
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "cognito-idp.amazonaws.com",
	}

	if userPool, ok := props["UserPool"]; ok {
		permissionProps["SourceArn"] = t.buildCognitoUserPoolArn(userPool)
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildCognitoUserPoolArn creates a Cognito User Pool ARN.
func (t *FunctionTransformer) buildCognitoUserPoolArn(userPool interface{}) interface{} {
	switch v := userPool.(type) {
	case string:
		if strings.HasPrefix(v, "arn:") {
			return v
		}
		return map[string]interface{}{
			"Fn::Sub": fmt.Sprintf("arn:aws:cognito-idp:${AWS::Region}:${AWS::AccountId}:userpool/%s", v),
		}
	case map[string]interface{}:
		if _, hasRef := v["Ref"]; hasRef {
			return map[string]interface{}{
				"Fn::GetAtt": []interface{}{v["Ref"], "Arn"},
			}
		}
	}
	return userPool
}

// buildMSKEvent creates resources for an MSK event source.
func (t *FunctionTransformer) buildMSKEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	esmID := logicalID + eventName
	esmProps := map[string]interface{}{
		"FunctionName": functionRef,
	}

	if stream, ok := props["Stream"]; ok {
		esmProps["EventSourceArn"] = stream
	}
	if topics, ok := props["Topics"]; ok {
		esmProps["Topics"] = topics
	}
	if startingPosition, ok := props["StartingPosition"]; ok {
		esmProps["StartingPosition"] = startingPosition
	}
	if batchSize, ok := props["BatchSize"]; ok {
		esmProps["BatchSize"] = batchSize
	}
	if consumerGroupId, ok := props["ConsumerGroupId"]; ok {
		esmProps["AmazonManagedKafkaEventSourceConfig"] = map[string]interface{}{
			"ConsumerGroupId": consumerGroupId,
		}
	}

	resources[esmID] = map[string]interface{}{
		"Type":       lambda.ResourceTypeEventSourceMapping,
		"Properties": esmProps,
	}

	return resources, nil
}

// buildMQEvent creates resources for an MQ event source.
func (t *FunctionTransformer) buildMQEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	esmID := logicalID + eventName
	esmProps := map[string]interface{}{
		"FunctionName": functionRef,
	}

	if broker, ok := props["Broker"]; ok {
		esmProps["EventSourceArn"] = broker
	}
	if queues, ok := props["Queues"]; ok {
		esmProps["Queues"] = queues
	}
	if batchSize, ok := props["BatchSize"]; ok {
		esmProps["BatchSize"] = batchSize
	}
	if sourceAccessConfigs, ok := props["SourceAccessConfigurations"]; ok {
		esmProps["SourceAccessConfigurations"] = sourceAccessConfigs
	}

	resources[esmID] = map[string]interface{}{
		"Type":       lambda.ResourceTypeEventSourceMapping,
		"Properties": esmProps,
	}

	return resources, nil
}

// buildSelfManagedKafkaEvent creates resources for a self-managed Kafka event source.
func (t *FunctionTransformer) buildSelfManagedKafkaEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	esmID := logicalID + eventName
	esmProps := map[string]interface{}{
		"FunctionName": functionRef,
	}

	if brokers, ok := props["KafkaBootstrapServers"]; ok {
		esmProps["SelfManagedEventSource"] = map[string]interface{}{
			"Endpoints": map[string]interface{}{
				"KafkaBootstrapServers": brokers,
			},
		}
	}
	if topics, ok := props["Topics"]; ok {
		esmProps["Topics"] = topics
	}
	if startingPosition, ok := props["StartingPosition"]; ok {
		esmProps["StartingPosition"] = startingPosition
	}
	if batchSize, ok := props["BatchSize"]; ok {
		esmProps["BatchSize"] = batchSize
	}
	if sourceAccessConfigs, ok := props["SourceAccessConfigurations"]; ok {
		esmProps["SourceAccessConfigurations"] = sourceAccessConfigs
	}

	resources[esmID] = map[string]interface{}{
		"Type":       lambda.ResourceTypeEventSourceMapping,
		"Properties": esmProps,
	}

	return resources, nil
}

// buildCloudWatchLogsEvent creates resources for a CloudWatch Logs event source.
func (t *FunctionTransformer) buildCloudWatchLogsEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Lambda permission for CloudWatch Logs
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "logs.amazonaws.com",
	}

	if logGroup, ok := props["LogGroupName"]; ok {
		permissionProps["SourceArn"] = map[string]interface{}{
			"Fn::Sub": fmt.Sprintf("arn:aws:logs:${AWS::Region}:${AWS::AccountId}:log-group:%v:*", logGroup),
		}
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	// Create Subscription Filter
	filterID := logicalID + eventName + "SubscriptionFilter"
	filterProps := map[string]interface{}{
		"DestinationArn": functionRef,
	}

	if logGroup, ok := props["LogGroupName"]; ok {
		filterProps["LogGroupName"] = logGroup
	}
	if filterPattern, ok := props["FilterPattern"]; ok {
		filterProps["FilterPattern"] = filterPattern
	}

	resources[filterID] = map[string]interface{}{
		"Type":       "AWS::Logs::SubscriptionFilter",
		"Properties": filterProps,
		"DependsOn":  permissionID,
	}

	return resources, nil
}

// buildAlexaSkillEvent creates resources for an Alexa Skill event source.
func (t *FunctionTransformer) buildAlexaSkillEvent(logicalID, eventName string, props map[string]interface{}, functionRef interface{}) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	// Create Lambda permission for Alexa
	permissionID := logicalID + eventName + "Permission"
	permissionProps := map[string]interface{}{
		"Action":       "lambda:InvokeFunction",
		"FunctionName": functionRef,
		"Principal":    "alexa-appkit.amazon.com",
	}

	if skillId, ok := props["SkillId"]; ok {
		permissionProps["EventSourceToken"] = skillId
	}

	resources[permissionID] = map[string]interface{}{
		"Type":       lambda.ResourceTypePermission,
		"Properties": permissionProps,
	}

	return resources, nil
}

// buildDeploymentPreference creates CodeDeploy resources for gradual deployments.
func (t *FunctionTransformer) buildDeploymentPreference(logicalID string, f *Function) (map[string]interface{}, error) {
	resources := make(map[string]interface{})

	deployPref := f.DeploymentPreference
	deployType, _ := deployPref["Type"].(string)

	// Skip if deployment is disabled
	if deployType == "" || deployType == "AllAtOnce" {
		return resources, nil
	}

	// Create CodeDeploy Application (if not exists)
	appID := "ServerlessDeploymentApplication"
	resources[appID] = map[string]interface{}{
		"Type": "AWS::CodeDeploy::Application",
		"Properties": map[string]interface{}{
			"ComputePlatform": "Lambda",
		},
	}

	// Create CodeDeploy Deployment Group
	groupID := logicalID + "DeploymentGroup"
	groupProps := map[string]interface{}{
		"ApplicationName": map[string]interface{}{"Ref": appID},
		"DeploymentConfigName": map[string]interface{}{
			"Fn::Sub": fmt.Sprintf("CodeDeployDefault.Lambda%s", deployType),
		},
		"DeploymentStyle": map[string]interface{}{
			"DeploymentOption": "WITH_TRAFFIC_CONTROL",
			"DeploymentType":   "BLUE_GREEN",
		},
	}

	// Add alarms if specified
	if alarms, ok := deployPref["Alarms"]; ok {
		groupProps["AlarmConfiguration"] = map[string]interface{}{
			"Enabled": true,
			"Alarms":  alarms,
		}
	}

	// Add hooks if specified
	if hooks, ok := deployPref["Hooks"]; ok {
		if hooksMap, ok := hooks.(map[string]interface{}); ok {
			deploymentGroupDeploymentStyle := groupProps["DeploymentStyle"].(map[string]interface{})
			if preTraffic, ok := hooksMap["PreTraffic"]; ok {
				deploymentGroupDeploymentStyle["PreTrafficHook"] = preTraffic
			}
			if postTraffic, ok := hooksMap["PostTraffic"]; ok {
				deploymentGroupDeploymentStyle["PostTrafficHook"] = postTraffic
			}
		}
	}

	resources[groupID] = map[string]interface{}{
		"Type":       "AWS::CodeDeploy::DeploymentGroup",
		"Properties": groupProps,
	}

	return resources, nil
}
