// Package parser provides YAML/JSON template parsing with intrinsic function detection.
package parser

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// intrinsicTagMapping maps YAML short-form tags to their CloudFormation intrinsic function names.
var intrinsicTagMapping = map[string]string{
	"!Ref":         "Ref",
	"!Sub":         "Fn::Sub",
	"!GetAtt":      "Fn::GetAtt",
	"!Join":        "Fn::Join",
	"!If":          "Fn::If",
	"!Select":      "Fn::Select",
	"!FindInMap":   "Fn::FindInMap",
	"!Base64":      "Fn::Base64",
	"!Cidr":        "Fn::Cidr",
	"!GetAZs":      "Fn::GetAZs",
	"!ImportValue": "Fn::ImportValue",
	"!Split":       "Fn::Split",
	"!Transform":   "Fn::Transform",
	"!And":         "Fn::And",
	"!Equals":      "Fn::Equals",
	"!Not":         "Fn::Not",
	"!Or":          "Fn::Or",
	"!Condition":   "Condition",
}

// NodeWithLocation wraps a parsed value with its source location.
type NodeWithLocation struct {
	Value  interface{}
	Line   int
	Column int
}

// unmarshalYAMLNode recursively unmarshals a yaml.Node, handling CloudFormation intrinsic tags.
// It returns the unmarshaled value and tracks source locations for error reporting.
func unmarshalYAMLNode(node *yaml.Node) (interface{}, error) {
	// Handle intrinsic function short-form tags (e.g., !Ref, !Sub)
	if intrinsicName, ok := intrinsicTagMapping[node.Tag]; ok {
		return handleIntrinsicTag(node, intrinsicName)
	}

	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) > 0 {
			return unmarshalYAMLNode(node.Content[0])
		}
		return nil, nil

	case yaml.MappingNode:
		return unmarshalMappingNode(node)

	case yaml.SequenceNode:
		return unmarshalSequenceNode(node)

	case yaml.ScalarNode:
		return unmarshalScalarNode(node)

	case yaml.AliasNode:
		if node.Alias != nil {
			return unmarshalYAMLNode(node.Alias)
		}
		return nil, nil

	default:
		return nil, nil
	}
}

// handleIntrinsicTag converts a YAML node with an intrinsic tag to the long-form map representation.
func handleIntrinsicTag(node *yaml.Node, intrinsicName string) (interface{}, error) {
	var value interface{}
	var err error

	// Special handling for !GetAtt with string value containing "."
	if intrinsicName == "Fn::GetAtt" && node.Kind == yaml.ScalarNode {
		// Convert "Resource.Attribute" to ["Resource", "Attribute"]
		parts := strings.SplitN(node.Value, ".", 2)
		if len(parts) == 2 {
			value = parts
		} else {
			value = node.Value
		}
	} else {
		// For other intrinsics, unmarshal the value normally
		switch node.Kind {
		case yaml.ScalarNode:
			value, err = unmarshalScalarNode(node)
		case yaml.SequenceNode:
			value, err = unmarshalSequenceNode(node)
		case yaml.MappingNode:
			value, err = unmarshalMappingNode(node)
		default:
			value = node.Value
		}
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		intrinsicName: value,
	}, nil
}

// unmarshalMappingNode unmarshals a YAML mapping node to a map[string]interface{}.
func unmarshalMappingNode(node *yaml.Node) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Mapping nodes have key-value pairs in Content
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		// Get the key as a string
		key := keyNode.Value

		// Recursively unmarshal the value
		value, err := unmarshalYAMLNode(valueNode)
		if err != nil {
			return nil, err
		}

		result[key] = value
	}

	return result, nil
}

// unmarshalSequenceNode unmarshals a YAML sequence node to a []interface{}.
func unmarshalSequenceNode(node *yaml.Node) ([]interface{}, error) {
	result := make([]interface{}, 0, len(node.Content))

	for _, item := range node.Content {
		value, err := unmarshalYAMLNode(item)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}

// unmarshalScalarNode unmarshals a YAML scalar node, preserving type information.
func unmarshalScalarNode(node *yaml.Node) (interface{}, error) {
	var value interface{}

	// Use yaml.Unmarshal to properly convert scalar types (int, bool, float, etc.)
	if err := node.Decode(&value); err != nil {
		return nil, err
	}

	return value, nil
}

// parseYAMLWithIntrinsics parses YAML content with proper handling of CloudFormation intrinsic tags.
func parseYAMLWithIntrinsics(data []byte) (map[string]interface{}, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	result, err := unmarshalYAMLNode(&root)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	// If the result is not a map, wrap it
	return map[string]interface{}{
		"value": result,
	}, nil
}

// LocationTracker tracks source locations for template elements.
type LocationTracker struct {
	locations map[string]SourceLocation
}

// SourceLocation represents a position in the source template.
type SourceLocation struct {
	Line   int
	Column int
}

// NewLocationTracker creates a new LocationTracker.
func NewLocationTracker() *LocationTracker {
	return &LocationTracker{
		locations: make(map[string]SourceLocation),
	}
}

// Track records the location of an element with the given path.
func (lt *LocationTracker) Track(path string, line, column int) {
	lt.locations[path] = SourceLocation{Line: line, Column: column}
}

// Get retrieves the location for the given path.
func (lt *LocationTracker) Get(path string) (SourceLocation, bool) {
	loc, ok := lt.locations[path]
	return loc, ok
}

// parseYAMLWithLocations parses YAML content and tracks source locations.
func parseYAMLWithLocations(data []byte) (map[string]interface{}, *LocationTracker, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, nil, err
	}

	tracker := NewLocationTracker()
	result, err := unmarshalYAMLNodeWithLocations(&root, "", tracker)
	if err != nil {
		return nil, nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, tracker, nil
	}

	return map[string]interface{}{
		"value": result,
	}, tracker, nil
}

// unmarshalYAMLNodeWithLocations unmarshals a YAML node while tracking source locations.
func unmarshalYAMLNodeWithLocations(node *yaml.Node, path string, tracker *LocationTracker) (interface{}, error) {
	// Track the location of this node
	if path != "" {
		tracker.Track(path, node.Line, node.Column)
	}

	// Handle intrinsic function short-form tags
	if intrinsicName, ok := intrinsicTagMapping[node.Tag]; ok {
		return handleIntrinsicTagWithLocations(node, intrinsicName, path, tracker)
	}

	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) > 0 {
			return unmarshalYAMLNodeWithLocations(node.Content[0], path, tracker)
		}
		return nil, nil

	case yaml.MappingNode:
		return unmarshalMappingNodeWithLocations(node, path, tracker)

	case yaml.SequenceNode:
		return unmarshalSequenceNodeWithLocations(node, path, tracker)

	case yaml.ScalarNode:
		return unmarshalScalarNode(node)

	case yaml.AliasNode:
		if node.Alias != nil {
			return unmarshalYAMLNodeWithLocations(node.Alias, path, tracker)
		}
		return nil, nil

	default:
		return nil, nil
	}
}

// handleIntrinsicTagWithLocations handles intrinsic tags while tracking locations.
func handleIntrinsicTagWithLocations(node *yaml.Node, intrinsicName, path string, tracker *LocationTracker) (interface{}, error) {
	var value interface{}
	var err error

	intrinsicPath := path
	if intrinsicPath != "" {
		intrinsicPath += "." + intrinsicName
	} else {
		intrinsicPath = intrinsicName
	}
	tracker.Track(intrinsicPath, node.Line, node.Column)

	if intrinsicName == "Fn::GetAtt" && node.Kind == yaml.ScalarNode {
		parts := strings.SplitN(node.Value, ".", 2)
		if len(parts) == 2 {
			value = parts
		} else {
			value = node.Value
		}
	} else {
		switch node.Kind {
		case yaml.ScalarNode:
			value, err = unmarshalScalarNode(node)
		case yaml.SequenceNode:
			value, err = unmarshalSequenceNodeWithLocations(node, intrinsicPath, tracker)
		case yaml.MappingNode:
			value, err = unmarshalMappingNodeWithLocations(node, intrinsicPath, tracker)
		default:
			value = node.Value
		}
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		intrinsicName: value,
	}, nil
}

// unmarshalMappingNodeWithLocations unmarshals a mapping node while tracking locations.
func unmarshalMappingNodeWithLocations(node *yaml.Node, path string, tracker *LocationTracker) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		key := keyNode.Value
		keyPath := path
		if keyPath != "" {
			keyPath += "." + key
		} else {
			keyPath = key
		}

		value, err := unmarshalYAMLNodeWithLocations(valueNode, keyPath, tracker)
		if err != nil {
			return nil, err
		}

		result[key] = value
	}

	return result, nil
}

// unmarshalSequenceNodeWithLocations unmarshals a sequence node while tracking locations.
func unmarshalSequenceNodeWithLocations(node *yaml.Node, path string, tracker *LocationTracker) ([]interface{}, error) {
	result := make([]interface{}, 0, len(node.Content))

	for i, item := range node.Content {
		itemPath := path + "[" + string(rune('0'+i)) + "]"
		if i >= 10 {
			// For indices >= 10, use proper string conversion
			itemPath = path + "[" + intToString(i) + "]"
		}

		value, err := unmarshalYAMLNodeWithLocations(item, itemPath, tracker)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}

// intToString converts an integer to its string representation.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
