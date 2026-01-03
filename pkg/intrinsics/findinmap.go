package intrinsics

import "fmt"

// FindInMapAction handles the Fn::FindInMap intrinsic function.
// Fn::FindInMap returns the value corresponding to keys in a two-level map.
type FindInMapAction struct{}

// Name returns the intrinsic function name.
func (a *FindInMapAction) Name() string {
	return "Fn::FindInMap"
}

// Resolve resolves a Fn::FindInMap intrinsic function.
// Fn::FindInMap takes the form: [MapName, TopLevelKey, SecondLevelKey]
// It looks up template Mappings to find the corresponding value.
func (a *FindInMapAction) Resolve(ctx *ResolveContext, value interface{}) (interface{}, error) {
	arr, ok := value.([]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("expected array, got %T", value))
	}

	if len(arr) != 3 {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("expected 3 elements, got %d", len(arr)))
	}

	// Extract map name, top level key, and second level key
	mapName, err := a.extractKey(arr[0], "MapName")
	if err != nil {
		return nil, err
	}

	topLevelKey, err := a.extractKey(arr[1], "TopLevelKey")
	if err != nil {
		return nil, err
	}

	secondLevelKey, err := a.extractKey(arr[2], "SecondLevelKey")
	if err != nil {
		return nil, err
	}

	// If any key is still an intrinsic (not a string), preserve for CloudFormation
	if mapName == "" || topLevelKey == "" || secondLevelKey == "" {
		return map[string]interface{}{"Fn::FindInMap": value}, nil
	}

	// Look up the mapping in the template
	if ctx.Template == nil || ctx.Template.Mappings == nil {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("mapping '%s' not found", mapName))
	}

	mapping, ok := ctx.Template.Mappings[mapName]
	if !ok {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("mapping '%s' not found", mapName))
	}

	mappingMap, ok := mapping.(map[string]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("mapping '%s' is not a map", mapName))
	}

	topLevel, ok := mappingMap[topLevelKey]
	if !ok {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("key '%s' not found in mapping '%s'", topLevelKey, mapName))
	}

	topLevelMap, ok := topLevel.(map[string]interface{})
	if !ok {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("value at '%s.%s' is not a map", mapName, topLevelKey))
	}

	result, ok := topLevelMap[secondLevelKey]
	if !ok {
		return nil, NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("key '%s' not found in mapping '%s.%s'", secondLevelKey, mapName, topLevelKey))
	}

	return result, nil
}

// extractKey extracts a string key from a value, handling intrinsics.
func (a *FindInMapAction) extractKey(value interface{}, keyName string) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case map[string]interface{}:
		// This is an unresolved intrinsic - return empty to signal pass-through
		return "", nil
	default:
		return "", NewIntrinsicError("Fn::FindInMap", fmt.Sprintf("%s must be string or intrinsic, got %T", keyName, value))
	}
}
