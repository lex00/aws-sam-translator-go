package lambda

// Alias represents an AWS::Lambda::Alias CloudFormation resource.
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-alias.html
type Alias struct {
	// Description is a description of the alias.
	Description string `json:"Description,omitempty" yaml:"Description,omitempty"`

	// FunctionName is the name or ARN of the Lambda function (required).
	FunctionName interface{} `json:"FunctionName" yaml:"FunctionName"`

	// FunctionVersion is the function version that the alias invokes (required).
	FunctionVersion interface{} `json:"FunctionVersion" yaml:"FunctionVersion"`

	// Name is the name of the alias (required).
	Name string `json:"Name" yaml:"Name"`

	// ProvisionedConcurrencyConfig specifies provisioned concurrency.
	ProvisionedConcurrencyConfig *AliasProvisionedConcurrencyConfig `json:"ProvisionedConcurrencyConfig,omitempty" yaml:"ProvisionedConcurrencyConfig,omitempty"`

	// RoutingConfig specifies weighted routing for canary deployments.
	RoutingConfig *AliasRoutingConfig `json:"RoutingConfig,omitempty" yaml:"RoutingConfig,omitempty"`
}

// AliasProvisionedConcurrencyConfig specifies provisioned concurrency for an alias.
type AliasProvisionedConcurrencyConfig struct {
	// ProvisionedConcurrentExecutions is the amount of provisioned concurrency.
	ProvisionedConcurrentExecutions int `json:"ProvisionedConcurrentExecutions" yaml:"ProvisionedConcurrentExecutions"`
}

// AliasRoutingConfig specifies weighted routing configuration.
type AliasRoutingConfig struct {
	// AdditionalVersionWeights specifies additional versions with traffic weights.
	AdditionalVersionWeights []VersionWeight `json:"AdditionalVersionWeights,omitempty" yaml:"AdditionalVersionWeights,omitempty"`
}

// VersionWeight specifies a version and its traffic weight.
type VersionWeight struct {
	// FunctionVersion is the version identifier.
	FunctionVersion interface{} `json:"FunctionVersion" yaml:"FunctionVersion"`

	// FunctionWeight is the percentage of traffic to route (0.0 to 1.0).
	FunctionWeight float64 `json:"FunctionWeight" yaml:"FunctionWeight"`
}

// NewAlias creates a new Alias with the required properties.
func NewAlias(functionName interface{}, functionVersion interface{}, name string) *Alias {
	return &Alias{
		FunctionName:    functionName,
		FunctionVersion: functionVersion,
		Name:            name,
	}
}

// NewAliasWithDescription creates a new Alias with a description.
func NewAliasWithDescription(functionName interface{}, functionVersion interface{}, name, description string) *Alias {
	return &Alias{
		FunctionName:    functionName,
		FunctionVersion: functionVersion,
		Name:            name,
		Description:     description,
	}
}

// WithProvisionedConcurrency sets provisioned concurrency for the alias.
func (a *Alias) WithProvisionedConcurrency(count int) *Alias {
	a.ProvisionedConcurrencyConfig = &AliasProvisionedConcurrencyConfig{
		ProvisionedConcurrentExecutions: count,
	}
	return a
}

// WithRoutingConfig sets the routing configuration for canary deployments.
func (a *Alias) WithRoutingConfig(versionWeights []VersionWeight) *Alias {
	a.RoutingConfig = &AliasRoutingConfig{
		AdditionalVersionWeights: versionWeights,
	}
	return a
}

// AddVersionWeight adds a version weight for traffic shifting.
func (a *Alias) AddVersionWeight(version interface{}, weight float64) *Alias {
	if a.RoutingConfig == nil {
		a.RoutingConfig = &AliasRoutingConfig{}
	}
	a.RoutingConfig.AdditionalVersionWeights = append(
		a.RoutingConfig.AdditionalVersionWeights,
		VersionWeight{
			FunctionVersion: version,
			FunctionWeight:  weight,
		},
	)
	return a
}

// ToCloudFormation converts the Alias to a CloudFormation resource.
func (a *Alias) ToCloudFormation() map[string]interface{} {
	properties := make(map[string]interface{})

	properties["FunctionName"] = a.FunctionName
	properties["FunctionVersion"] = a.FunctionVersion
	properties["Name"] = a.Name

	if a.Description != "" {
		properties["Description"] = a.Description
	}
	if a.ProvisionedConcurrencyConfig != nil {
		properties["ProvisionedConcurrencyConfig"] = a.ProvisionedConcurrencyConfig.toMap()
	}
	if a.RoutingConfig != nil && len(a.RoutingConfig.AdditionalVersionWeights) > 0 {
		properties["RoutingConfig"] = a.RoutingConfig.toMap()
	}

	return map[string]interface{}{
		"Type":       ResourceTypeAlias,
		"Properties": properties,
	}
}

func (p *AliasProvisionedConcurrencyConfig) toMap() map[string]interface{} {
	return map[string]interface{}{
		"ProvisionedConcurrentExecutions": p.ProvisionedConcurrentExecutions,
	}
}

func (r *AliasRoutingConfig) toMap() map[string]interface{} {
	if len(r.AdditionalVersionWeights) == 0 {
		return nil
	}

	weights := make([]map[string]interface{}, len(r.AdditionalVersionWeights))
	for i, vw := range r.AdditionalVersionWeights {
		weights[i] = map[string]interface{}{
			"FunctionVersion": vw.FunctionVersion,
			"FunctionWeight":  vw.FunctionWeight,
		}
	}

	return map[string]interface{}{
		"AdditionalVersionWeights": weights,
	}
}
