// Package types provides core type definitions for SAM and CloudFormation resources.
package types

// Template represents a SAM or CloudFormation template.
type Template struct {
	AWSTemplateFormatVersion string                 `json:"AWSTemplateFormatVersion,omitempty" yaml:"AWSTemplateFormatVersion,omitempty"`
	Transform                interface{}            `json:"Transform,omitempty" yaml:"Transform,omitempty"`
	Description              string                 `json:"Description,omitempty" yaml:"Description,omitempty"`
	Metadata                 map[string]interface{} `json:"Metadata,omitempty" yaml:"Metadata,omitempty"`
	Parameters               map[string]Parameter   `json:"Parameters,omitempty" yaml:"Parameters,omitempty"`
	Mappings                 map[string]interface{} `json:"Mappings,omitempty" yaml:"Mappings,omitempty"`
	Conditions               map[string]interface{} `json:"Conditions,omitempty" yaml:"Conditions,omitempty"`
	Resources                map[string]Resource    `json:"Resources,omitempty" yaml:"Resources,omitempty"`
	Outputs                  map[string]Output      `json:"Outputs,omitempty" yaml:"Outputs,omitempty"`
	Globals                  map[string]interface{} `json:"Globals,omitempty" yaml:"Globals,omitempty"`
}

// Parameter represents a CloudFormation parameter.
type Parameter struct {
	Type                  string      `json:"Type" yaml:"Type"`
	Default               interface{} `json:"Default,omitempty" yaml:"Default,omitempty"`
	Description           string      `json:"Description,omitempty" yaml:"Description,omitempty"`
	AllowedValues         []string    `json:"AllowedValues,omitempty" yaml:"AllowedValues,omitempty"`
	AllowedPattern        string      `json:"AllowedPattern,omitempty" yaml:"AllowedPattern,omitempty"`
	ConstraintDescription string      `json:"ConstraintDescription,omitempty" yaml:"ConstraintDescription,omitempty"`
	MaxLength             int         `json:"MaxLength,omitempty" yaml:"MaxLength,omitempty"`
	MinLength             int         `json:"MinLength,omitempty" yaml:"MinLength,omitempty"`
	MaxValue              float64     `json:"MaxValue,omitempty" yaml:"MaxValue,omitempty"`
	MinValue              float64     `json:"MinValue,omitempty" yaml:"MinValue,omitempty"`
	NoEcho                bool        `json:"NoEcho,omitempty" yaml:"NoEcho,omitempty"`
}

// Resource represents a CloudFormation or SAM resource.
type Resource struct {
	Type           string                 `json:"Type" yaml:"Type"`
	Properties     map[string]interface{} `json:"Properties,omitempty" yaml:"Properties,omitempty"`
	Metadata       map[string]interface{} `json:"Metadata,omitempty" yaml:"Metadata,omitempty"`
	DependsOn      interface{}            `json:"DependsOn,omitempty" yaml:"DependsOn,omitempty"`
	Condition      string                 `json:"Condition,omitempty" yaml:"Condition,omitempty"`
	DeletionPolicy string                 `json:"DeletionPolicy,omitempty" yaml:"DeletionPolicy,omitempty"`
	UpdatePolicy   map[string]interface{} `json:"UpdatePolicy,omitempty" yaml:"UpdatePolicy,omitempty"`
}

// Output represents a CloudFormation output.
type Output struct {
	Description string      `json:"Description,omitempty" yaml:"Description,omitempty"`
	Value       interface{} `json:"Value" yaml:"Value"`
	Export      *Export     `json:"Export,omitempty" yaml:"Export,omitempty"`
	Condition   string      `json:"Condition,omitempty" yaml:"Condition,omitempty"`
}

// Export represents an output export configuration.
type Export struct {
	Name interface{} `json:"Name" yaml:"Name"`
}
