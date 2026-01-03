// Package sam provides SAM resource transformers.
package sam

// Function represents an AWS::Serverless::Function resource.
type Function struct {
	Handler       string                 `json:"Handler,omitempty" yaml:"Handler,omitempty"`
	Runtime       string                 `json:"Runtime,omitempty" yaml:"Runtime,omitempty"`
	CodeUri       interface{}            `json:"CodeUri,omitempty" yaml:"CodeUri,omitempty"`
	Description   string                 `json:"Description,omitempty" yaml:"Description,omitempty"`
	MemorySize    int                    `json:"MemorySize,omitempty" yaml:"MemorySize,omitempty"`
	Timeout       int                    `json:"Timeout,omitempty" yaml:"Timeout,omitempty"`
	Role          interface{}            `json:"Role,omitempty" yaml:"Role,omitempty"`
	Policies      interface{}            `json:"Policies,omitempty" yaml:"Policies,omitempty"`
	Environment   map[string]interface{} `json:"Environment,omitempty" yaml:"Environment,omitempty"`
	Events        map[string]interface{} `json:"Events,omitempty" yaml:"Events,omitempty"`
	Tags          map[string]string      `json:"Tags,omitempty" yaml:"Tags,omitempty"`
	Layers        []interface{}          `json:"Layers,omitempty" yaml:"Layers,omitempty"`
	VpcConfig     map[string]interface{} `json:"VpcConfig,omitempty" yaml:"VpcConfig,omitempty"`
	FunctionName  interface{}            `json:"FunctionName,omitempty" yaml:"FunctionName,omitempty"`
	Architectures []string               `json:"Architectures,omitempty" yaml:"Architectures,omitempty"`
}

// FunctionTransformer transforms AWS::Serverless::Function to CloudFormation.
type FunctionTransformer struct{}

// NewFunctionTransformer creates a new FunctionTransformer.
func NewFunctionTransformer() *FunctionTransformer {
	return &FunctionTransformer{}
}

// Transform converts a SAM Function to CloudFormation resources.
func (t *FunctionTransformer) Transform(logicalID string, f *Function) (map[string]interface{}, error) {
	// TODO: Implement function transformation
	return nil, nil
}
