package utils

import (
	"reflect"
	"testing"
)

func TestDeepMerge(t *testing.T) {
	tests := []struct {
		name     string
		dst      map[string]interface{}
		src      map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "empty maps",
			dst:      map[string]interface{}{},
			src:      map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: "src overwrites dst",
			dst:  map[string]interface{}{"a": 1},
			src:  map[string]interface{}{"a": 2},
			expected: map[string]interface{}{
				"a": 2,
			},
		},
		{
			name: "merge adds new keys",
			dst:  map[string]interface{}{"a": 1},
			src:  map[string]interface{}{"b": 2},
			expected: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		},
		{
			name: "deep merge nested maps",
			dst: map[string]interface{}{
				"outer": map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			src: map[string]interface{}{
				"outer": map[string]interface{}{
					"b": 3,
					"c": 4,
				},
			},
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"a": 1,
					"b": 3,
					"c": 4,
				},
			},
		},
		{
			name: "non-map overwrites map",
			dst: map[string]interface{}{
				"key": map[string]interface{}{"nested": 1},
			},
			src: map[string]interface{}{
				"key": "string value",
			},
			expected: map[string]interface{}{
				"key": "string value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeepMerge(tt.dst, tt.src)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeepMerge() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeepCopy(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name:  "empty map",
			input: map[string]interface{}{},
		},
		{
			name: "simple map",
			input: map[string]interface{}{
				"string": "value",
				"number": 42,
				"bool":   true,
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
		},
		{
			name: "map with slice",
			input: map[string]interface{}{
				"list": []interface{}{"a", "b", "c"},
			},
		},
		{
			name: "complex nested structure",
			input: map[string]interface{}{
				"resources": map[string]interface{}{
					"function": map[string]interface{}{
						"type": "AWS::Lambda::Function",
						"tags": []interface{}{
							map[string]interface{}{"Key": "Name", "Value": "Test"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeepCopy(tt.input)

			// Check equality
			if !reflect.DeepEqual(result, tt.input) {
				t.Errorf("DeepCopy() = %v, want %v", result, tt.input)
			}

			// Check independence (modifying copy shouldn't affect original)
			if len(result) > 0 {
				for k := range result {
					result[k] = "modified"
					break
				}
				if reflect.DeepEqual(result, tt.input) && len(tt.input) > 0 {
					t.Error("DeepCopy() did not create independent copy")
				}
			}
		})
	}
}

func TestSortedKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected []string
	}{
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: []string{},
		},
		{
			name: "single key",
			input: map[string]interface{}{
				"alpha": 1,
			},
			expected: []string{"alpha"},
		},
		{
			name: "multiple keys",
			input: map[string]interface{}{
				"charlie": 3,
				"alpha":   1,
				"bravo":   2,
			},
			expected: []string{"alpha", "bravo", "charlie"},
		},
		{
			name: "AWS resource names",
			input: map[string]interface{}{
				"MyFunction":       nil,
				"ApiGateway":       nil,
				"DynamoDBTable":    nil,
				"LambdaPermission": nil,
			},
			expected: []string{"ApiGateway", "DynamoDBTable", "LambdaPermission", "MyFunction"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedKeys(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortedKeys() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedStringKeys(t *testing.T) {
	input := map[string]string{
		"z": "last",
		"a": "first",
		"m": "middle",
	}
	expected := []string{"a", "m", "z"}

	result := SortedStringKeys(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortedStringKeys() = %v, want %v", result, expected)
	}
}

func TestSortedResourceKeys(t *testing.T) {
	type Resource struct {
		Type string
	}

	input := map[string]Resource{
		"Zebra":    {Type: "AWS::S3::Bucket"},
		"Apple":    {Type: "AWS::Lambda::Function"},
		"Mushroom": {Type: "AWS::DynamoDB::Table"},
	}
	expected := []string{"Apple", "Mushroom", "Zebra"}

	result := SortedResourceKeys(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortedResourceKeys() = %v, want %v", result, expected)
	}
}

func TestDeepEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{
			name:     "equal strings",
			a:        "hello",
			b:        "hello",
			expected: true,
		},
		{
			name:     "different strings",
			a:        "hello",
			b:        "world",
			expected: false,
		},
		{
			name:     "equal maps",
			a:        map[string]interface{}{"key": "value"},
			b:        map[string]interface{}{"key": "value"},
			expected: true,
		},
		{
			name:     "different maps",
			a:        map[string]interface{}{"key": "value1"},
			b:        map[string]interface{}{"key": "value2"},
			expected: false,
		},
		{
			name:     "equal slices",
			a:        []interface{}{1, 2, 3},
			b:        []interface{}{1, 2, 3},
			expected: true,
		},
		{
			name:     "different slices",
			a:        []interface{}{1, 2, 3},
			b:        []interface{}{1, 2, 4},
			expected: false,
		},
		{
			name:     "nil values",
			a:        nil,
			b:        nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeepEqual(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("DeepEqual() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMapContains(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]interface{}
		subset   map[string]interface{}
		expected bool
	}{
		{
			name:     "empty subset",
			m:        map[string]interface{}{"a": 1, "b": 2},
			subset:   map[string]interface{}{},
			expected: true,
		},
		{
			name:     "exact match",
			m:        map[string]interface{}{"a": 1, "b": 2},
			subset:   map[string]interface{}{"a": 1, "b": 2},
			expected: true,
		},
		{
			name:     "partial match",
			m:        map[string]interface{}{"a": 1, "b": 2, "c": 3},
			subset:   map[string]interface{}{"a": 1, "b": 2},
			expected: true,
		},
		{
			name:     "missing key",
			m:        map[string]interface{}{"a": 1},
			subset:   map[string]interface{}{"a": 1, "b": 2},
			expected: false,
		},
		{
			name:     "different value",
			m:        map[string]interface{}{"a": 1, "b": 2},
			subset:   map[string]interface{}{"a": 1, "b": 3},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapContains(tt.m, tt.subset)
			if result != tt.expected {
				t.Errorf("MapContains() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStringSliceContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		value    string
		expected bool
	}{
		{
			name:     "empty slice",
			slice:    []string{},
			value:    "test",
			expected: false,
		},
		{
			name:     "contains value",
			slice:    []string{"a", "b", "c"},
			value:    "b",
			expected: true,
		},
		{
			name:     "does not contain value",
			slice:    []string{"a", "b", "c"},
			value:    "d",
			expected: false,
		},
		{
			name:     "first element",
			slice:    []string{"first", "second", "third"},
			value:    "first",
			expected: true,
		},
		{
			name:     "last element",
			slice:    []string{"first", "second", "third"},
			value:    "third",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringSliceContains(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("StringSliceContains() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "already sorted",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "reverse order",
			input:    []string{"c", "b", "a"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "random order",
			input:    []string{"banana", "apple", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := make([]string, len(tt.input))
			copy(original, tt.input)

			result := SortStringSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortStringSlice() = %v, want %v", result, tt.expected)
			}

			// Verify original slice is not modified
			if !reflect.DeepEqual(tt.input, original) {
				t.Error("SortStringSlice() modified the original slice")
			}
		})
	}
}

func TestUniqueStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "all same",
			input:    []string{"x", "x", "x"},
			expected: []string{"x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqueStrings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqueStrings() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	tests := []struct {
		name     string
		maps     []map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "no maps",
			maps:     []map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: "single map",
			maps: []map[string]interface{}{
				{"a": 1, "b": 2},
			},
			expected: map[string]interface{}{"a": 1, "b": 2},
		},
		{
			name: "multiple maps",
			maps: []map[string]interface{}{
				{"a": 1},
				{"b": 2},
				{"c": 3},
			},
			expected: map[string]interface{}{"a": 1, "b": 2, "c": 3},
		},
		{
			name: "later maps override",
			maps: []map[string]interface{}{
				{"a": 1, "b": 2},
				{"b": 3, "c": 4},
			},
			expected: map[string]interface{}{"a": 1, "b": 3, "c": 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeMaps(tt.maps...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeMaps() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name         string
		m            map[string]interface{}
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "key exists with string value",
			m:            map[string]interface{}{"key": "value"},
			key:          "key",
			defaultValue: "default",
			expected:     "value",
		},
		{
			name:         "key does not exist",
			m:            map[string]interface{}{"other": "value"},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "key exists but not string",
			m:            map[string]interface{}{"key": 123},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "empty map",
			m:            map[string]interface{}{},
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStringValue(tt.m, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetStringValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetMapValue(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]interface{}
		key      string
		expected map[string]interface{}
	}{
		{
			name:     "key exists with map value",
			m:        map[string]interface{}{"key": map[string]interface{}{"nested": "value"}},
			key:      "key",
			expected: map[string]interface{}{"nested": "value"},
		},
		{
			name:     "key does not exist",
			m:        map[string]interface{}{"other": "value"},
			key:      "key",
			expected: nil,
		},
		{
			name:     "key exists but not map",
			m:        map[string]interface{}{"key": "string"},
			key:      "key",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMapValue(tt.m, tt.key)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetMapValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetSliceValue(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]interface{}
		key      string
		expected []interface{}
	}{
		{
			name:     "key exists with slice value",
			m:        map[string]interface{}{"key": []interface{}{1, 2, 3}},
			key:      "key",
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "key does not exist",
			m:        map[string]interface{}{"other": "value"},
			key:      "key",
			expected: nil,
		},
		{
			name:     "key exists but not slice",
			m:        map[string]interface{}{"key": "string"},
			key:      "key",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSliceValue(tt.m, tt.key)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetSliceValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}
