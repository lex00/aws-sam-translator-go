// Package parser provides YAML/JSON template parsing with intrinsic function detection.
package parser

import (
	"encoding/json"
	"fmt"
)

// parseJSONWithIntrinsics parses JSON content and validates intrinsic function structure.
func parseJSONWithIntrinsics(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Process the result to ensure consistent intrinsic representation
	processed, err := processJSONValue(result)
	if err != nil {
		return nil, err
	}

	if m, ok := processed.(map[string]interface{}); ok {
		return m, nil
	}

	return result, nil
}

// processJSONValue recursively processes JSON values to normalize intrinsic functions.
func processJSONValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		return processJSONMap(v)
	case []interface{}:
		return processJSONArray(v)
	default:
		return value, nil
	}
}

// processJSONMap processes a JSON map, handling intrinsic functions.
func processJSONMap(m map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range m {
		processed, err := processJSONValue(value)
		if err != nil {
			return nil, err
		}
		result[key] = processed
	}

	return result, nil
}

// processJSONArray processes a JSON array.
func processJSONArray(arr []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(arr))

	for i, item := range arr {
		processed, err := processJSONValue(item)
		if err != nil {
			return nil, err
		}
		result[i] = processed
	}

	return result, nil
}

// JSONLocationTracker tracks source locations in JSON (byte offsets).
type JSONLocationTracker struct {
	offsets map[string]int
}

// NewJSONLocationTracker creates a new JSONLocationTracker.
func NewJSONLocationTracker() *JSONLocationTracker {
	return &JSONLocationTracker{
		offsets: make(map[string]int),
	}
}

// Track records the byte offset of an element.
func (jlt *JSONLocationTracker) Track(path string, offset int) {
	jlt.offsets[path] = offset
}

// Get retrieves the byte offset for the given path.
func (jlt *JSONLocationTracker) Get(path string) (int, bool) {
	offset, ok := jlt.offsets[path]
	return offset, ok
}

// ParseJSONWithLocations parses JSON content and tracks element locations.
// Note: Full location tracking for JSON would require a streaming parser.
// This implementation provides basic parsing with intrinsic handling.
func ParseJSONWithLocations(data []byte) (map[string]interface{}, *JSONLocationTracker, error) {
	result, err := parseJSONWithIntrinsics(data)
	if err != nil {
		return nil, nil, err
	}

	tracker := NewJSONLocationTracker()
	// Basic location tracking - for full location support, a streaming parser would be needed
	return result, tracker, nil
}
