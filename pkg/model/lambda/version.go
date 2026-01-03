package lambda

// Version represents an AWS::Lambda::Version CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-version.html
type Version struct {
	// CodeSha256 is the SHA256 hash of the deployment package.
	// Used to validate that the function code hasn't changed.
	CodeSha256 string `json:"CodeSha256,omitempty" yaml:"CodeSha256,omitempty"`

	// Description is a description for the version.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// FunctionName is the name or ARN of the Lambda function (required).
	FunctionName interface{} `json:"FunctionName" yaml:"FunctionName"`

	// Policy is the resource-based policy for the version.
	Policy interface{} `json:"Policy,omitempty" yaml:"Policy,omitempty"`

	// ProvisionedConcurrencyConfig specifies provisioned concurrency.
	ProvisionedConcurrencyConfig *ProvisionedConcurrencyConfig `json:"ProvisionedConcurrencyConfig,omitempty" yaml:"ProvisionedConcurrencyConfig,omitempty"`

	// RuntimePolicy configures runtime management for the version.
	RuntimePolicy *RuntimePolicy `json:"RuntimePolicy,omitempty" yaml:"RuntimePolicy,omitempty"`
}

// ProvisionedConcurrencyConfig specifies provisioned concurrency for a version.
type ProvisionedConcurrencyConfig struct {
	// ProvisionedConcurrentExecutions is the amount of provisioned concurrency.
	ProvisionedConcurrentExecutions int `json:"ProvisionedConcurrentExecutions" yaml:"ProvisionedConcurrentExecutions"`
}

// RuntimePolicy configures runtime management for a version.
type RuntimePolicy struct {
	// RuntimeVersionArn is the ARN of the runtime version.
	RuntimeVersionArn interface{} `json:"RuntimeVersionArn,omitempty" yaml:"RuntimeVersionArn,omitempty"`

	// UpdateRuntimeOn specifies when to update the runtime.
	// Valid values: Auto, FunctionUpdate, Manual
	UpdateRuntimeOn string `json:"UpdateRuntimeOn" yaml:"UpdateRuntimeOn"`
}

// NewVersion creates a new Version with the required function name.
func NewVersion(functionName interface{}) *Version {
	return &Version{
		FunctionName: functionName,
	}
}

// NewVersionWithDescription creates a new Version with a description.
func NewVersionWithDescription(functionName interface{}, description string) *Version {
	return &Version{
		FunctionName: functionName,
		Description:  description,
	}
}

// WithCodeSha256 sets the code SHA256 hash for validation.
func (v *Version) WithCodeSha256(sha256 string) *Version {
	v.CodeSha256 = sha256
	return v
}

// WithProvisionedConcurrency sets provisioned concurrency for the version.
func (v *Version) WithProvisionedConcurrency(count int) *Version {
	v.ProvisionedConcurrencyConfig = &ProvisionedConcurrencyConfig{
		ProvisionedConcurrentExecutions: count,
	}
	return v
}

// WithRuntimePolicy sets the runtime policy for the version.
func (v *Version) WithRuntimePolicy(updateRuntimeOn string, runtimeVersionArn interface{}) *Version {
	v.RuntimePolicy = &RuntimePolicy{
		UpdateRuntimeOn:   updateRuntimeOn,
		RuntimeVersionArn: runtimeVersionArn,
	}
	return v
}

// ToCloudFormation converts the Version to a CloudFormation resource.
func (v *Version) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	properties["FunctionName"] = v.FunctionName

	if v.CodeSha256 != "" {
		properties["CodeSha256"] = v.CodeSha256
	}
	if v.Description != "" {
		properties["Description"] = v.Description
	}
	if v.Policy != nil {
		properties["Policy"] = v.Policy
	}
	if v.ProvisionedConcurrencyConfig != nil {
		properties["ProvisionedConcurrencyConfig"] = v.ProvisionedConcurrencyConfig.toMap()
	}
	if v.RuntimePolicy != nil {
		properties["RuntimePolicy"] = v.RuntimePolicy.toMap()
	}

	return map[string]interface{}{
		"Type":       ResourceTypeVersion,
		"Properties": properties,
	}
}

func (p *ProvisionedConcurrencyConfig) toMap() map[string]interface{} {
	return map[string]interface{}{
		"ProvisionedConcurrentExecutions": p.ProvisionedConcurrentExecutions,
	}
}

func (r *RuntimePolicy) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["UpdateRuntimeOn"] = r.UpdateRuntimeOn
	if r.RuntimeVersionArn != nil {
		m["RuntimeVersionArn"] = r.RuntimeVersionArn
	}
	return m
}
