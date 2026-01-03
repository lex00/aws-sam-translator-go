package translator

import (
	"testing"
)

func TestLogicalIDGenerator_Generate(t *testing.T) {
	gen := NewLogicalIDGenerator()

	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "single part",
			parts:    []string{"MyFunction"},
			expected: "MyFunction",
		},
		{
			name:     "multiple parts",
			parts:    []string{"My", "Function"},
			expected: "MyFunction",
		},
		{
			name:     "with special characters",
			parts:    []string{"My-Function_Name"},
			expected: "MyFunctionName",
		},
		{
			name:     "with numbers",
			parts:    []string{"Function123"},
			expected: "Function123",
		},
		{
			name:     "starts with number",
			parts:    []string{"123Function"},
			expected: "R123Function",
		},
		{
			name:     "empty input",
			parts:    []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.Generate(tt.parts...)
			if result != tt.expected {
				t.Errorf("Generate(%v) = %q, want %q", tt.parts, result, tt.expected)
			}
		})
	}
}

func TestLogicalIDGenerator_GenerateWithPrefix(t *testing.T) {
	gen := NewLogicalIDGeneratorWithPrefix("SAM")

	result := gen.Generate("MyFunction")
	if result != "SAMMyFunction" {
		t.Errorf("Generate with prefix = %q, want %q", result, "SAMMyFunction")
	}
}

func TestLogicalIDGenerator_GenerateHashed(t *testing.T) {
	gen := NewLogicalIDGenerator()

	// Same data should produce same hash
	id1 := gen.GenerateHashed("test-data", "Function")
	id2 := gen.GenerateHashed("test-data", "Function")
	if id1 != id2 {
		t.Errorf("Hashed IDs should be stable: %q != %q", id1, id2)
	}

	// Different data should produce different hash
	id3 := gen.GenerateHashed("different-data", "Function")
	if id1 == id3 {
		t.Errorf("Different data should produce different IDs: %q == %q", id1, id3)
	}

	// ID should contain the base part
	if len(id1) < len("Function")+LogicalIDHashLength {
		t.Errorf("Hashed ID too short: %q", id1)
	}
}

func TestLogicalIDGenerator_GenerateAPIDeploymentID(t *testing.T) {
	gen := NewLogicalIDGenerator()

	// Same spec should produce same ID
	id1 := gen.GenerateAPIDeploymentID("MyAPI", `{"openapi": "3.0"}`)
	id2 := gen.GenerateAPIDeploymentID("MyAPI", `{"openapi": "3.0"}`)
	if id1 != id2 {
		t.Errorf("API deployment IDs should be stable: %q != %q", id1, id2)
	}

	// Different spec should produce different ID
	id3 := gen.GenerateAPIDeploymentID("MyAPI", `{"openapi": "3.1"}`)
	if id1 == id3 {
		t.Errorf("Different specs should produce different IDs: %q == %q", id1, id3)
	}

	// ID should contain "Deployment"
	if id1[:len("MyAPIDeployment")] != "MyAPIDeployment" {
		t.Errorf("API deployment ID should start with 'MyAPIDeployment': %q", id1)
	}
}

func TestLogicalIDGenerator_GenerateFromMap(t *testing.T) {
	gen := NewLogicalIDGenerator()

	// Same map should produce same ID (deterministic)
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	id1 := gen.GenerateFromMap("Prefix", data)
	id2 := gen.GenerateFromMap("Prefix", data)
	if id1 != id2 {
		t.Errorf("Map-based IDs should be stable: %q != %q", id1, id2)
	}

	// Order shouldn't matter (deterministic serialization)
	data2 := map[string]interface{}{
		"key2": 123,
		"key1": "value1",
	}
	id3 := gen.GenerateFromMap("Prefix", data2)
	if id1 != id3 {
		t.Errorf("Map order should not affect ID: %q != %q", id1, id3)
	}
}

func TestIsValidLogicalID(t *testing.T) {
	tests := []struct {
		id    string
		valid bool
	}{
		{"MyFunction", true},
		{"Function123", true},
		{"a", true},
		{"A123B456", true},
		{"123Function", false}, // starts with number
		{"My-Function", false}, // contains hyphen
		{"My_Function", false}, // contains underscore
		{"", false},            // empty
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			result := IsValidLogicalID(tt.id)
			if result != tt.valid {
				t.Errorf("IsValidLogicalID(%q) = %v, want %v", tt.id, result, tt.valid)
			}
		})
	}
}

func TestMakeLogicalIDSafe(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"MyFunction", "MyFunction"},
		{"My-Function", "MyFunction"},
		{"My_Function", "MyFunction"},
		{"123Function", "R123Function"},
		{"my.function.name", "myfunctionname"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MakeLogicalIDSafe(tt.input)
			if result != tt.expected {
				t.Errorf("MakeLogicalIDSafe(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeLogicalID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"MyFunction", "MyFunction"},
		{"My-Function", "MyFunction"},
		{"123Start", "R123Start"},
		{"special!@#chars", "specialchars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeLogicalID(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeLogicalID(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalculateHash(t *testing.T) {
	// Hash should be deterministic
	hash1 := calculateHash("test data")
	hash2 := calculateHash("test data")
	if hash1 != hash2 {
		t.Errorf("Hash should be deterministic: %q != %q", hash1, hash2)
	}

	// Hash should have correct length
	if len(hash1) != LogicalIDHashLength {
		t.Errorf("Hash length = %d, want %d", len(hash1), LogicalIDHashLength)
	}

	// Different data should produce different hash
	hash3 := calculateHash("different data")
	if hash1 == hash3 {
		t.Errorf("Different data should produce different hash: %q == %q", hash1, hash3)
	}
}

func TestSerializeMapDeterministic(t *testing.T) {
	// Maps with same content should produce same serialization
	map1 := map[string]interface{}{
		"b": 2,
		"a": 1,
		"c": 3,
	}
	map2 := map[string]interface{}{
		"a": 1,
		"c": 3,
		"b": 2,
	}

	result1 := serializeMapDeterministic(map1)
	result2 := serializeMapDeterministic(map2)

	if result1 != result2 {
		t.Errorf("Serialization should be deterministic regardless of insertion order: %q != %q", result1, result2)
	}
}

func TestSerializeValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"nil", nil, "nil"},
		{"string", "test", "test"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"slice", []interface{}{"a", "b"}, "[a,b]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := serializeValue(tt.value)
			if result != tt.expected {
				t.Errorf("serializeValue(%v) = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}
