// Package parser provides YAML/JSON template parsing with intrinsic function detection.
package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/types"
)

// Parser handles SAM/CloudFormation template parsing.
type Parser struct {
	// TrackLocations enables source location tracking for error reporting.
	TrackLocations bool
	// Locations stores the tracked source locations when TrackLocations is enabled.
	Locations *LocationTracker
}

// New creates a new Parser instance.
func New() *Parser {
	return &Parser{
		TrackLocations: false,
		Locations:      nil,
	}
}

// NewWithLocationTracking creates a new Parser with location tracking enabled.
func NewWithLocationTracking() *Parser {
	return &Parser{
		TrackLocations: true,
		Locations:      NewLocationTracker(),
	}
}

// ParseYAML parses a YAML template with full intrinsic function support.
func (p *Parser) ParseYAML(data []byte) (*types.Template, error) {
	var rawData map[string]interface{}
	var err error

	if p.TrackLocations {
		rawData, p.Locations, err = parseYAMLWithLocations(data)
	} else {
		rawData, err = parseYAMLWithIntrinsics(data)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return p.mapToTemplate(rawData)
}

// ParseJSON parses a JSON template.
func (p *Parser) ParseJSON(data []byte) (*types.Template, error) {
	rawData, err := parseJSONWithIntrinsics(data)
	if err != nil {
		return nil, err
	}

	return p.mapToTemplate(rawData)
}

// Parse automatically detects the format (YAML or JSON) and parses the template.
func (p *Parser) Parse(data []byte) (*types.Template, error) {
	// Try to detect format based on content
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) > 0 && trimmed[0] == '{' {
		return p.ParseJSON(data)
	}
	return p.ParseYAML(data)
}

// ParseRawYAML parses YAML and returns the raw map structure without converting to Template.
func (p *Parser) ParseRawYAML(data []byte) (map[string]interface{}, error) {
	if p.TrackLocations {
		result, tracker, err := parseYAMLWithLocations(data)
		p.Locations = tracker
		return result, err
	}
	return parseYAMLWithIntrinsics(data)
}

// ParseRawJSON parses JSON and returns the raw map structure without converting to Template.
func (p *Parser) ParseRawJSON(data []byte) (map[string]interface{}, error) {
	return parseJSONWithIntrinsics(data)
}

// mapToTemplate converts a raw map to a Template struct.
func (p *Parser) mapToTemplate(data map[string]interface{}) (*types.Template, error) {
	template := &types.Template{}

	// AWSTemplateFormatVersion
	if v, ok := data["AWSTemplateFormatVersion"]; ok {
		if s, ok := v.(string); ok {
			template.AWSTemplateFormatVersion = s
		}
	}

	// Transform
	if v, ok := data["Transform"]; ok {
		template.Transform = v
	}

	// Description
	if v, ok := data["Description"]; ok {
		if s, ok := v.(string); ok {
			template.Description = s
		}
	}

	// Metadata
	if v, ok := data["Metadata"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			template.Metadata = m
		}
	}

	// Parameters
	if v, ok := data["Parameters"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			template.Parameters = p.parseParameters(m)
		}
	}

	// Mappings
	if v, ok := data["Mappings"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			template.Mappings = m
		}
	}

	// Conditions
	if v, ok := data["Conditions"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			template.Conditions = m
		}
	}

	// Resources
	if v, ok := data["Resources"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			resources, err := p.parseResources(m)
			if err != nil {
				return nil, err
			}
			template.Resources = resources
		}
	}

	// Outputs
	if v, ok := data["Outputs"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			template.Outputs = p.parseOutputs(m)
		}
	}

	// Globals (SAM-specific)
	if v, ok := data["Globals"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			template.Globals = m
		}
	}

	return template, nil
}

// parseParameters converts a raw map to Parameter structs.
func (p *Parser) parseParameters(data map[string]interface{}) map[string]types.Parameter {
	params := make(map[string]types.Parameter)

	for name, value := range data {
		if m, ok := value.(map[string]interface{}); ok {
			param := types.Parameter{}

			if v, ok := m["Type"].(string); ok {
				param.Type = v
			}
			if v, ok := m["Default"]; ok {
				param.Default = v
			}
			if v, ok := m["Description"].(string); ok {
				param.Description = v
			}
			if v, ok := m["AllowedValues"].([]interface{}); ok {
				param.AllowedValues = toStringSlice(v)
			}
			if v, ok := m["AllowedPattern"].(string); ok {
				param.AllowedPattern = v
			}
			if v, ok := m["ConstraintDescription"].(string); ok {
				param.ConstraintDescription = v
			}
			if v, ok := m["MaxLength"].(int); ok {
				param.MaxLength = v
			}
			if v, ok := m["MinLength"].(int); ok {
				param.MinLength = v
			}
			if v, ok := m["MaxValue"].(float64); ok {
				param.MaxValue = v
			}
			if v, ok := m["MinValue"].(float64); ok {
				param.MinValue = v
			}
			if v, ok := m["NoEcho"].(bool); ok {
				param.NoEcho = v
			}

			params[name] = param
		}
	}

	return params
}

// parseResources converts a raw map to Resource structs.
func (p *Parser) parseResources(data map[string]interface{}) (map[string]types.Resource, error) {
	resources := make(map[string]types.Resource)

	for name, value := range data {
		if m, ok := value.(map[string]interface{}); ok {
			resource := types.Resource{}

			if v, ok := m["Type"].(string); ok {
				resource.Type = v
			} else {
				return nil, &ParseError{
					Message:  fmt.Sprintf("resource '%s' is missing required 'Type' property", name),
					Location: p.getLocation("Resources." + name),
				}
			}

			if v, ok := m["Properties"].(map[string]interface{}); ok {
				resource.Properties = v
			}
			if v, ok := m["Metadata"].(map[string]interface{}); ok {
				resource.Metadata = v
			}
			if v, ok := m["DependsOn"]; ok {
				resource.DependsOn = v
			}
			if v, ok := m["Condition"].(string); ok {
				resource.Condition = v
			}
			if v, ok := m["DeletionPolicy"].(string); ok {
				resource.DeletionPolicy = v
			}
			if v, ok := m["UpdatePolicy"].(map[string]interface{}); ok {
				resource.UpdatePolicy = v
			}

			resources[name] = resource
		}
	}

	return resources, nil
}

// parseOutputs converts a raw map to Output structs.
func (p *Parser) parseOutputs(data map[string]interface{}) map[string]types.Output {
	outputs := make(map[string]types.Output)

	for name, value := range data {
		if m, ok := value.(map[string]interface{}); ok {
			output := types.Output{}

			if v, ok := m["Description"].(string); ok {
				output.Description = v
			}
			if v, ok := m["Value"]; ok {
				output.Value = v
			}
			if v, ok := m["Condition"].(string); ok {
				output.Condition = v
			}
			if v, ok := m["Export"].(map[string]interface{}); ok {
				export := &types.Export{}
				if name, ok := v["Name"]; ok {
					export.Name = name
				}
				output.Export = export
			}

			outputs[name] = output
		}
	}

	return outputs
}

// getLocation returns the source location for a given path.
func (p *Parser) getLocation(path string) SourceLocation {
	if p.Locations != nil {
		if loc, ok := p.Locations.Get(path); ok {
			return loc
		}
	}
	return SourceLocation{}
}

// toStringSlice converts a []interface{} to []string.
func toStringSlice(arr []interface{}) []string {
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// ParseError represents a parsing error with source location.
type ParseError struct {
	Message  string
	Location SourceLocation
}

func (e *ParseError) Error() string {
	if e.Location.Line > 0 {
		return fmt.Sprintf("parse error at line %d, column %d: %s", e.Location.Line, e.Location.Column, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

// ValidateTemplate performs basic structural validation on a template.
func ValidateTemplate(data map[string]interface{}) error {
	// Check for required Resources section
	if _, ok := data["Resources"]; !ok {
		return &ParseError{
			Message: "template must have a 'Resources' section",
		}
	}

	resources, ok := data["Resources"].(map[string]interface{})
	if !ok {
		return &ParseError{
			Message: "'Resources' must be a mapping",
		}
	}

	// Validate each resource has a Type
	for name, resource := range resources {
		resMap, ok := resource.(map[string]interface{})
		if !ok {
			return &ParseError{
				Message: fmt.Sprintf("resource '%s' must be a mapping", name),
			}
		}

		if _, ok := resMap["Type"]; !ok {
			return &ParseError{
				Message: fmt.Sprintf("resource '%s' is missing required 'Type' property", name),
			}
		}
	}

	return nil
}

// ValidateIntrinsics validates all intrinsic functions in the template.
func ValidateIntrinsics(data map[string]interface{}) []error {
	var errs []error
	validateIntrinsicsRecursive(data, "", &errs)
	return errs
}

func validateIntrinsicsRecursive(value interface{}, path string, errs *[]error) {
	switch v := value.(type) {
	case map[string]interface{}:
		if name := GetIntrinsicName(v); name != "" {
			if err := ValidateIntrinsicStructure(name, v[name]); err != nil {
				*errs = append(*errs, err)
			}
			// Also validate nested intrinsics within the value
			validateIntrinsicsRecursive(v[name], path+"."+name, errs)
		} else {
			for key, val := range v {
				childPath := path
				if childPath != "" {
					childPath += "." + key
				} else {
					childPath = key
				}
				validateIntrinsicsRecursive(val, childPath, errs)
			}
		}
	case []interface{}:
		for i, item := range v {
			itemPath := path + "[" + intToString(i) + "]"
			validateIntrinsicsRecursive(item, itemPath, errs)
		}
	}
}

// MarshalYAML converts a template back to YAML bytes.
func MarshalYAML(template *types.Template) ([]byte, error) {
	return json.Marshal(template)
}

// MarshalJSON converts a template to JSON bytes.
func MarshalJSON(template *types.Template) ([]byte, error) {
	return json.MarshalIndent(template, "", "  ")
}
