// Package lambda provides CloudFormation resource models for AWS Lambda.
package lambda

// ResourceType constants for Lambda resources.
const (
	ResourceTypeFunction           = "AWS::Lambda::Function"
	ResourceTypeVersion            = "AWS::Lambda::Version"
	ResourceTypeAlias              = "AWS::Lambda::Alias"
	ResourceTypePermission         = "AWS::Lambda::Permission"
	ResourceTypeEventSourceMapping = "AWS::Lambda::EventSourceMapping"
	ResourceTypeLayerVersion       = "AWS::Lambda::LayerVersion"
)

// Function represents an AWS::Lambda::Function CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-function.html
type Function struct {
	// Architectures is the instruction set architecture for the function.
	// Valid values: x86_64, arm64
	Architectures []string `json:"Architectures,omitempty" yaml:"Architectures,omitempty"`

	// Code is the code for the function (required).
	Code *Code `json:"Code" yaml:"Code"`

	// CodeSigningConfigArn is the ARN of a code-signing configuration.
	CodeSigningConfigArn interface{} `json:"CodeSigningConfigArn,omitempty" yaml:"CodeSigningConfigArn,omitempty"`

	// DeadLetterConfig configures error handling for asynchronous invocation.
	DeadLetterConfig *DeadLetterConfig `json:"DeadLetterConfig,omitempty" yaml:"DeadLetterConfig,omitempty"`

	// Description is a description of the function.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// Environment contains environment variables for the function.
	Environment *Environment `json:"Environment,omitempty" yaml:"Environment,omitempty"`

	// EphemeralStorage configures the size of the function's /tmp directory.
	EphemeralStorage *EphemeralStorage `json:"EphemeralStorage,omitempty" yaml:"EphemeralStorage,omitempty"`

	// FileSystemConfigs connects the function to an Amazon EFS file system.
	FileSystemConfigs []FileSystemConfig `json:"FileSystemConfigs,omitempty" yaml:"FileSystemConfigs,omitempty"`

	// FunctionName is the name of the Lambda function.
	FunctionName interface{} `json:"FunctionName,omitempty" yaml:"FunctionName,omitempty"`

	// Handler is the name of the method within your code that Lambda calls.
	Handler string `json:"Handler,omitempty" yaml:"Handler,omitempty"`

	// ImageConfig overrides the container image settings.
	ImageConfig *ImageConfig `json:"ImageConfig,omitempty" yaml:"ImageConfig,omitempty"`

	// KmsKeyArn is the ARN of the KMS key used to encrypt environment variables.
	KmsKeyArn interface{} `json:"KmsKeyArn,omitempty" yaml:"KmsKeyArn,omitempty"`

	// Layers is a list of function layer ARNs to add to the function.
	Layers []interface{} `json:"Layers,omitempty" yaml:"Layers,omitempty"`

	// LoggingConfig configures Amazon CloudWatch logging for the function.
	LoggingConfig *LoggingConfig `json:"LoggingConfig,omitempty" yaml:"LoggingConfig,omitempty"`

	// MemorySize is the amount of memory available to the function (MB).
	MemorySize int `json:"MemorySize,omitempty" yaml:"MemorySize,omitempty"`

	// PackageType is the type of deployment package.
	// Valid values: Zip, Image
	PackageType string `json:"PackageType,omitempty" yaml:"PackageType,omitempty"`

	// RecursiveLoop sets the loop detection behavior for recursive function invocations.
	RecursiveLoop string `json:"RecursiveLoop,omitempty" yaml:"RecursiveLoop,omitempty"`

	// ReservedConcurrentExecutions is the number of concurrent executions reserved.
	ReservedConcurrentExecutions *int `json:"ReservedConcurrentExecutions,omitempty" yaml:"ReservedConcurrentExecutions,omitempty"`

	// Role is the ARN of the function's execution role (required).
	Role interface{} `json:"Role" yaml:"Role"`

	// Runtime is the identifier of the function's runtime.
	Runtime string `json:"Runtime,omitempty" yaml:"Runtime,omitempty"`

	// RuntimeManagementConfig configures runtime management.
	RuntimeManagementConfig *RuntimeManagementConfig `json:"RuntimeManagementConfig,omitempty" yaml:"RuntimeManagementConfig,omitempty"`

	// SnapStart enables SnapStart for Java functions.
	SnapStart *SnapStart `json:"SnapStart,omitempty" yaml:"SnapStart,omitempty"`

	// Tags is a list of tags to apply to the function.
	Tags []Tag `json:"Tags,omitempty" yaml:"Tags,omitempty"`

	// Timeout is the amount of time Lambda allows a function to run (seconds).
	Timeout int `json:"Timeout,omitempty" yaml:"Timeout,omitempty"`

	// TracingConfig configures AWS X-Ray tracing.
	TracingConfig *TracingConfig `json:"TracingConfig,omitempty" yaml:"TracingConfig,omitempty"`

	// VpcConfig connects the function to a VPC.
	VpcConfig *VpcConfig `json:"VpcConfig,omitempty" yaml:"VpcConfig,omitempty"`
}

// Code represents the deployment package for a Lambda function.
type Code struct {
	// ImageUri is the URI of a container image in Amazon ECR.
	ImageUri interface{} `json:"ImageUri,omitempty" yaml:"ImageUri,omitempty"`

	// S3Bucket is the S3 bucket name.
	S3Bucket interface{} `json:"S3Bucket,omitempty" yaml:"S3Bucket,omitempty"`

	// S3Key is the S3 object key.
	S3Key interface{} `json:"S3Key,omitempty" yaml:"S3Key,omitempty"`

	// S3ObjectVersion is the S3 object version.
	S3ObjectVersion interface{} `json:"S3ObjectVersion,omitempty" yaml:"S3ObjectVersion,omitempty"`

	// ZipFile is the inline code for the function.
	ZipFile interface{} `json:"ZipFile,omitempty" yaml:"ZipFile,omitempty"`
}

// DeadLetterConfig configures how Lambda handles events that fail.
type DeadLetterConfig struct {
	// TargetArn is the ARN of an SNS topic or SQS queue.
	TargetArn interface{} `json:"TargetArn,omitempty" yaml:"TargetArn,omitempty"`
}

// Environment contains environment variable configuration.
type Environment struct {
	// Variables is a map of environment variables.
	Variables map[string]interface{} `json:"Variables,omitempty" yaml:"Variables,omitempty"`
}

// EphemeralStorage configures the function's /tmp directory.
type EphemeralStorage struct {
	// Size is the size of the /tmp directory in MB (512-10240).
	Size int `json:"Size" yaml:"Size"`
}

// FileSystemConfig connects the function to an EFS file system.
type FileSystemConfig struct {
	// Arn is the ARN of the access point.
	Arn interface{} `json:"Arn" yaml:"Arn"`

	// LocalMountPath is the path where the function can access the file system.
	LocalMountPath string `json:"LocalMountPath" yaml:"LocalMountPath"`
}

// ImageConfig overrides container image configuration.
type ImageConfig struct {
	// Command overrides the CMD in the container image.
	Command []string `json:"Command,omitempty" yaml:"Command,omitempty"`

	// EntryPoint overrides the ENTRYPOINT in the container image.
	EntryPoint []string `json:"EntryPoint,omitempty" yaml:"EntryPoint,omitempty"`

	// WorkingDirectory sets the working directory.
	WorkingDirectory string `json:"WorkingDirectory,omitempty" yaml:"WorkingDirectory,omitempty"`
}

// LoggingConfig configures CloudWatch logging.
type LoggingConfig struct {
	// ApplicationLogLevel is the log level for application logs.
	// Valid values: TRACE, DEBUG, INFO, WARN, ERROR, FATAL
	ApplicationLogLevel string `json:"ApplicationLogLevel,omitempty" yaml:"ApplicationLogLevel,omitempty"`

	// LogFormat is the format for logging.
	// Valid values: Text, JSON
	LogFormat string `json:"LogFormat,omitempty" yaml:"LogFormat,omitempty"`

	// LogGroup is the CloudWatch log group name.
	LogGroup interface{} `json:"LogGroup,omitempty" yaml:"LogGroup,omitempty"`

	// SystemLogLevel is the log level for system logs.
	// Valid values: DEBUG, INFO, WARN
	SystemLogLevel string `json:"SystemLogLevel,omitempty" yaml:"SystemLogLevel,omitempty"`
}

// RuntimeManagementConfig configures runtime management settings.
type RuntimeManagementConfig struct {
	// RuntimeVersionArn is the ARN of a specific runtime version.
	RuntimeVersionArn interface{} `json:"RuntimeVersionArn,omitempty" yaml:"RuntimeVersionArn,omitempty"`

	// UpdateRuntimeOn specifies when to update the runtime.
	// Valid values: Auto, FunctionUpdate, Manual
	UpdateRuntimeOn string `json:"UpdateRuntimeOn" yaml:"UpdateRuntimeOn"`
}

// SnapStart enables SnapStart for the function.
type SnapStart struct {
	// ApplyOn specifies when to apply SnapStart.
	// Valid values: PublishedVersions, None
	ApplyOn string `json:"ApplyOn" yaml:"ApplyOn"`
}

// Tag represents a key-value tag.
type Tag struct {
	// Key is the tag key.
	Key string `json:"Key" yaml:"Key"`

	// Value is the tag value.
	Value string `json:"Value" yaml:"Value"`
}

// TracingConfig configures AWS X-Ray tracing.
type TracingConfig struct {
	// Mode is the tracing mode.
	// Valid values: Active, PassThrough
	Mode string `json:"Mode" yaml:"Mode"`
}

// VpcConfig connects the function to a VPC.
type VpcConfig struct {
	// Ipv6AllowedForDualStack allows outbound IPv6 traffic.
	Ipv6AllowedForDualStack bool `json:"Ipv6AllowedForDualStack,omitempty" yaml:"Ipv6AllowedForDualStack,omitempty"`

	// SecurityGroupIds is a list of VPC security group IDs.
	SecurityGroupIds []interface{} `json:"SecurityGroupIds,omitempty" yaml:"SecurityGroupIds,omitempty"`

	// SubnetIds is a list of VPC subnet IDs.
	SubnetIds []interface{} `json:"SubnetIds,omitempty" yaml:"SubnetIds,omitempty"`
}

// NewFunction creates a new Function with the required properties.
func NewFunction(code *Code, role interface{}) *Function {
	return &Function{
		Code: code,
		Role: role,
	}
}

// NewCodeFromS3 creates a Code configuration from an S3 location.
func NewCodeFromS3(bucket, key interface{}) *Code {
	return &Code{
		S3Bucket: bucket,
		S3Key:    key,
	}
}

// NewCodeFromS3WithVersion creates a Code configuration from an S3 location with version.
func NewCodeFromS3WithVersion(bucket, key, version interface{}) *Code {
	return &Code{
		S3Bucket:        bucket,
		S3Key:           key,
		S3ObjectVersion: version,
	}
}

// NewCodeFromImage creates a Code configuration from a container image.
func NewCodeFromImage(imageUri interface{}) *Code {
	return &Code{
		ImageUri: imageUri,
	}
}

// NewCodeFromZip creates a Code configuration from inline code.
func NewCodeFromZip(zipFile interface{}) *Code {
	return &Code{
		ZipFile: zipFile,
	}
}

// ToCloudFormation converts the Function to a CloudFormation resource.
func (f *Function) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	if len(f.Architectures) > 0 {
		properties["Architectures"] = f.Architectures
	}
	if f.Code != nil {
		properties["Code"] = f.Code.toMap()
	}
	if f.CodeSigningConfigArn != nil {
		properties["CodeSigningConfigArn"] = f.CodeSigningConfigArn
	}
	if f.DeadLetterConfig != nil {
		properties["DeadLetterConfig"] = f.DeadLetterConfig.toMap()
	}
	if f.Description != "" {
		properties["Description"] = f.Description
	}
	if f.Environment != nil {
		properties["Environment"] = f.Environment.toMap()
	}
	if f.EphemeralStorage != nil {
		properties["EphemeralStorage"] = f.EphemeralStorage.toMap()
	}
	if len(f.FileSystemConfigs) > 0 {
		configs := make([]map[string]interface{}, len(f.FileSystemConfigs))
		for i, fsc := range f.FileSystemConfigs {
			configs[i] = fsc.toMap()
		}
		properties["FileSystemConfigs"] = configs
	}
	if f.FunctionName != nil {
		properties["FunctionName"] = f.FunctionName
	}
	if f.Handler != "" {
		properties["Handler"] = f.Handler
	}
	if f.ImageConfig != nil {
		properties["ImageConfig"] = f.ImageConfig.toMap()
	}
	if f.KmsKeyArn != nil {
		properties["KmsKeyArn"] = f.KmsKeyArn
	}
	if len(f.Layers) > 0 {
		properties["Layers"] = f.Layers
	}
	if f.LoggingConfig != nil {
		properties["LoggingConfig"] = f.LoggingConfig.toMap()
	}
	if f.MemorySize > 0 {
		properties["MemorySize"] = f.MemorySize
	}
	if f.PackageType != "" {
		properties["PackageType"] = f.PackageType
	}
	if f.RecursiveLoop != "" {
		properties["RecursiveLoop"] = f.RecursiveLoop
	}
	if f.ReservedConcurrentExecutions != nil {
		properties["ReservedConcurrentExecutions"] = *f.ReservedConcurrentExecutions
	}
	if f.Role != nil {
		properties["Role"] = f.Role
	}
	if f.Runtime != "" {
		properties["Runtime"] = f.Runtime
	}
	if f.RuntimeManagementConfig != nil {
		properties["RuntimeManagementConfig"] = f.RuntimeManagementConfig.toMap()
	}
	if f.SnapStart != nil {
		properties["SnapStart"] = f.SnapStart.toMap()
	}
	if len(f.Tags) > 0 {
		tags := make([]map[string]interface{}, len(f.Tags))
		for i, t := range f.Tags {
			tags[i] = t.toMap()
		}
		properties["Tags"] = tags
	}
	if f.Timeout > 0 {
		properties["Timeout"] = f.Timeout
	}
	if f.TracingConfig != nil {
		properties["TracingConfig"] = f.TracingConfig.toMap()
	}
	if f.VpcConfig != nil {
		properties["VpcConfig"] = f.VpcConfig.toMap()
	}

	return map[string]interface{}{
		"Type":       ResourceTypeFunction,
		"Properties": properties,
	}
}

func (c *Code) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if c.ImageUri != nil {
		m["ImageUri"] = c.ImageUri
	}
	if c.S3Bucket != nil {
		m["S3Bucket"] = c.S3Bucket
	}
	if c.S3Key != nil {
		m["S3Key"] = c.S3Key
	}
	if c.S3ObjectVersion != nil {
		m["S3ObjectVersion"] = c.S3ObjectVersion
	}
	if c.ZipFile != nil {
		m["ZipFile"] = c.ZipFile
	}
	return m
}

func (d *DeadLetterConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if d.TargetArn != nil {
		m["TargetArn"] = d.TargetArn
	}
	return m
}

func (e *Environment) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if e.Variables != nil {
		m["Variables"] = e.Variables
	}
	return m
}

func (e *EphemeralStorage) toMap() map[string]interface{} {
	return map[string]interface{}{
		"Size": e.Size,
	}
}

func (f *FileSystemConfig) toMap() map[string]interface{} {
	return map[string]interface{}{
		"Arn":            f.Arn,
		"LocalMountPath": f.LocalMountPath,
	}
}

func (i *ImageConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if len(i.Command) > 0 {
		m["Command"] = i.Command
	}
	if len(i.EntryPoint) > 0 {
		m["EntryPoint"] = i.EntryPoint
	}
	if i.WorkingDirectory != "" {
		m["WorkingDirectory"] = i.WorkingDirectory
	}
	return m
}

func (l *LoggingConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if l.ApplicationLogLevel != "" {
		m["ApplicationLogLevel"] = l.ApplicationLogLevel
	}
	if l.LogFormat != "" {
		m["LogFormat"] = l.LogFormat
	}
	if l.LogGroup != nil {
		m["LogGroup"] = l.LogGroup
	}
	if l.SystemLogLevel != "" {
		m["SystemLogLevel"] = l.SystemLogLevel
	}
	return m
}

func (r *RuntimeManagementConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if r.RuntimeVersionArn != nil {
		m["RuntimeVersionArn"] = r.RuntimeVersionArn
	}
	m["UpdateRuntimeOn"] = r.UpdateRuntimeOn
	return m
}

func (s *SnapStart) toMap() map[string]interface{} {
	return map[string]interface{}{
		"ApplyOn": s.ApplyOn,
	}
}

func (t *Tag) toMap() map[string]interface{} {
	return map[string]interface{}{
		"Key":   t.Key,
		"Value": t.Value,
	}
}

func (t *TracingConfig) toMap() map[string]interface{} {
	return map[string]interface{}{
		"Mode": t.Mode,
	}
}

func (v *VpcConfig) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	if v.Ipv6AllowedForDualStack {
		m["Ipv6AllowedForDualStack"] = v.Ipv6AllowedForDualStack
	}
	if len(v.SecurityGroupIds) > 0 {
		m["SecurityGroupIds"] = v.SecurityGroupIds
	}
	if len(v.SubnetIds) > 0 {
		m["SubnetIds"] = v.SubnetIds
	}
	return m
}
