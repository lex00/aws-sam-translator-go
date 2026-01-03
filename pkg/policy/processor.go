// Package policy provides SAM policy template expansion.
package policy

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed templates.json
var templatesFS embed.FS

// PolicyTemplatesFile represents the structure of policy_templates.json.
type PolicyTemplatesFile struct {
	Templates map[string]Template `json:"Templates"`
	Version   string              `json:"Version"`
}

// Template represents a SAM policy template.
type Template struct {
	Definition  map[string]interface{}   `json:"Definition"`
	Description string                   `json:"Description"`
	Parameters  map[string]TemplateParam `json:"Parameters"`
}

// TemplateParam represents a template parameter definition.
type TemplateParam struct {
	Description string `json:"Description"`
}

// Processor handles policy template expansion.
type Processor struct {
	templates map[string]Template
	version   string
}

// New creates a new policy Processor with embedded templates.
func New() (*Processor, error) {
	data, err := templatesFS.ReadFile("templates.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded templates: %w", err)
	}

	return NewFromBytes(data)
}

// NewFromBytes creates a new policy Processor from JSON bytes.
func NewFromBytes(data []byte) (*Processor, error) {
	var file PolicyTemplatesFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("failed to parse policy templates: %w", err)
	}

	return &Processor{
		templates: file.Templates,
		version:   file.Version,
	}, nil
}

// Version returns the version of the loaded policy templates.
func (p *Processor) Version() string {
	return p.version
}

// TemplateNames returns a list of all available template names.
func (p *Processor) TemplateNames() []string {
	names := make([]string, 0, len(p.templates))
	for name := range p.templates {
		names = append(names, name)
	}
	return names
}

// HasTemplate checks if a template with the given name exists.
func (p *Processor) HasTemplate(name string) bool {
	_, ok := p.templates[name]
	return ok
}

// GetTemplate returns the template definition for inspection.
func (p *Processor) GetTemplate(name string) (Template, bool) {
	t, ok := p.templates[name]
	return t, ok
}

// Expand expands a policy template with the given parameters.
// It returns the IAM policy statements with parameters substituted.
func (p *Processor) Expand(templateName string, params map[string]interface{}) (map[string]interface{}, error) {
	template, ok := p.templates[templateName]
	if !ok {
		return nil, fmt.Errorf("unknown policy template: %s", templateName)
	}

	// Validate that all required parameters are provided
	for paramName := range template.Parameters {
		if _, exists := params[paramName]; !exists {
			return nil, fmt.Errorf("missing required parameter '%s' for template '%s'", paramName, templateName)
		}
	}

	// Deep copy and substitute parameters in the definition
	expanded := deepCopyAndSubstitute(template.Definition, params)

	result, ok := expanded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("template definition is not a map")
	}

	return result, nil
}

// ExpandStatements expands a policy template and returns just the Statement array.
func (p *Processor) ExpandStatements(templateName string, params map[string]interface{}) ([]interface{}, error) {
	definition, err := p.Expand(templateName, params)
	if err != nil {
		return nil, err
	}

	statements, ok := definition["Statement"]
	if !ok {
		return nil, fmt.Errorf("template '%s' has no Statement field", templateName)
	}

	stmtArray, ok := statements.([]interface{})
	if !ok {
		return nil, fmt.Errorf("template '%s' Statement is not an array", templateName)
	}

	return stmtArray, nil
}

// deepCopyAndSubstitute performs a deep copy of the value and substitutes
// parameter references with their actual values.
func deepCopyAndSubstitute(value interface{}, params map[string]interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		// Check if this is a Ref to a parameter
		if ref, ok := v["Ref"].(string); ok && len(v) == 1 {
			if paramValue, ok := params[ref]; ok {
				return paramValue
			}
		}

		// Check if this is an Fn::Sub expression
		if sub, ok := v["Fn::Sub"]; ok && len(v) == 1 {
			return processFnSub(sub, params)
		}

		// Regular map - deep copy all entries
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			result[key] = deepCopyAndSubstitute(val, params)
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = deepCopyAndSubstitute(item, params)
		}
		return result

	default:
		// Primitive types (string, number, bool, nil) - return as-is
		return v
	}
}

// processFnSub handles Fn::Sub intrinsic function parameter substitution.
// Fn::Sub can be either:
// - A simple string: "arn:aws:s3:::${BucketName}"
// - An array: ["${param}", {"param": {"Ref": "ParamName"}}]
func processFnSub(sub interface{}, params map[string]interface{}) interface{} {
	switch s := sub.(type) {
	case string:
		// Simple Fn::Sub with just a string - leave as-is for CloudFormation
		return map[string]interface{}{"Fn::Sub": s}

	case []interface{}:
		if len(s) != 2 {
			return map[string]interface{}{"Fn::Sub": sub}
		}

		template, ok := s[0].(string)
		if !ok {
			return map[string]interface{}{"Fn::Sub": sub}
		}

		varMap, ok := s[1].(map[string]interface{})
		if !ok {
			return map[string]interface{}{"Fn::Sub": sub}
		}

		// Process the variable map, substituting Refs with actual parameter values
		newVarMap := make(map[string]interface{}, len(varMap))
		for varName, varValue := range varMap {
			newVarMap[varName] = deepCopyAndSubstitute(varValue, params)
		}

		return map[string]interface{}{
			"Fn::Sub": []interface{}{template, newVarMap},
		}

	default:
		return map[string]interface{}{"Fn::Sub": sub}
	}
}
