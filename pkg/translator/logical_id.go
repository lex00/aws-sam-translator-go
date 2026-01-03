// Package translator provides the main SAM to CloudFormation transformation orchestrator.
package translator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// LogicalIDMaxLength is the maximum length for CloudFormation logical IDs.
const LogicalIDMaxLength = 255

// LogicalIDHashLength is the length of the hash suffix used for generated IDs.
const LogicalIDHashLength = 8

// logicalIDPattern is the valid pattern for CloudFormation logical IDs.
var logicalIDPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)

// LogicalIDGenerator generates deterministic CloudFormation logical IDs.
type LogicalIDGenerator struct {
	// prefix is an optional prefix to prepend to generated IDs.
	prefix string
}

// NewLogicalIDGenerator creates a new LogicalIDGenerator.
func NewLogicalIDGenerator() *LogicalIDGenerator {
	return &LogicalIDGenerator{}
}

// NewLogicalIDGeneratorWithPrefix creates a LogicalIDGenerator with a prefix.
func NewLogicalIDGeneratorWithPrefix(prefix string) *LogicalIDGenerator {
	return &LogicalIDGenerator{prefix: prefix}
}

// Generate creates a logical ID from the given parts.
// The generated ID is deterministic - the same inputs always produce the same output.
func (g *LogicalIDGenerator) Generate(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	// Join parts and sanitize
	combined := strings.Join(parts, "")
	sanitized := sanitizeLogicalID(combined)

	// Add prefix if present
	if g.prefix != "" {
		sanitized = sanitizeLogicalID(g.prefix) + sanitized
	}

	// Truncate if necessary to leave room for hash
	maxBase := LogicalIDMaxLength - LogicalIDHashLength
	if len(sanitized) > maxBase {
		sanitized = sanitized[:maxBase]
	}

	return sanitized
}

// GenerateHashed creates a logical ID with a deterministic hash suffix.
// This is useful when collisions might occur or when the ID needs to change
// when the underlying data changes.
func (g *LogicalIDGenerator) GenerateHashed(data string, parts ...string) string {
	// Generate the base ID
	base := g.Generate(parts...)
	if base == "" && data == "" {
		return ""
	}

	// Calculate hash from data
	hash := calculateHash(data)

	// Combine base and hash
	if base == "" {
		return hash
	}

	// Ensure we don't exceed max length
	maxBase := LogicalIDMaxLength - LogicalIDHashLength - 1 // -1 for potential separator handling
	if len(base)+len(hash) > LogicalIDMaxLength {
		base = base[:maxBase]
	}

	return base + hash
}

// GenerateAPIDeploymentID creates a logical ID for API Gateway deployments.
// The ID changes when the API specification changes, ensuring deployments
// are recreated when needed.
func (g *LogicalIDGenerator) GenerateAPIDeploymentID(apiLogicalID string, openAPISpec string) string {
	// Hash the OpenAPI spec to create a unique deployment ID
	specHash := calculateHash(openAPISpec)

	// Create the deployment ID
	base := apiLogicalID + "Deployment"
	return g.Generate(base) + specHash
}

// GenerateFromMap creates a deterministic logical ID from a map.
// The map is serialized in a deterministic order to ensure stable hashing.
func (g *LogicalIDGenerator) GenerateFromMap(prefix string, data map[string]interface{}) string {
	// Serialize map deterministically
	serialized := serializeMapDeterministic(data)

	// Generate hashed ID
	return g.GenerateHashed(serialized, prefix)
}

// sanitizeLogicalID removes invalid characters from a logical ID.
// CloudFormation logical IDs must be alphanumeric and start with a letter.
func sanitizeLogicalID(id string) string {
	var result strings.Builder
	result.Grow(len(id))

	startsWithLetter := false
	for _, r := range id {
		if r >= 'A' && r <= 'Z' {
			if result.Len() == 0 {
				startsWithLetter = true
			}
			result.WriteRune(r)
		} else if r >= 'a' && r <= 'z' {
			if result.Len() == 0 {
				startsWithLetter = true
			}
			result.WriteRune(r)
		} else if r >= '0' && r <= '9' {
			// Always include numbers - we'll add prefix if needed
			result.WriteRune(r)
		}
		// All other characters are dropped
	}

	// Ensure it starts with a letter
	str := result.String()
	if len(str) > 0 && !startsWithLetter {
		str = "R" + str
	}

	return str
}

// calculateHash computes a deterministic hash for the given data.
func calculateHash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	hash := hex.EncodeToString(h.Sum(nil))

	// Return first N characters (alphanumeric)
	return hash[:LogicalIDHashLength]
}

// serializeMapDeterministic converts a map to a deterministic string representation.
func serializeMapDeterministic(m map[string]interface{}) string {
	if m == nil {
		return ""
	}

	// Get sorted keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build deterministic string
	var result strings.Builder
	for _, k := range keys {
		result.WriteString(k)
		result.WriteString(":")
		result.WriteString(serializeValue(m[k]))
		result.WriteString(";")
	}

	return result.String()
}

// serializeValue converts a value to a deterministic string representation.
func serializeValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "nil"
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", val)
	case []interface{}:
		var parts []string
		for _, item := range val {
			parts = append(parts, serializeValue(item))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case map[string]interface{}:
		return "{" + serializeMapDeterministic(val) + "}"
	default:
		return fmt.Sprintf("%v", val)
	}
}

// IsValidLogicalID checks if a string is a valid CloudFormation logical ID.
func IsValidLogicalID(id string) bool {
	if id == "" || len(id) > LogicalIDMaxLength {
		return false
	}
	return logicalIDPattern.MatchString(id)
}

// MakeLogicalIDSafe ensures a string is a valid CloudFormation logical ID.
// If the input is already valid, it's returned unchanged.
// Otherwise, it's sanitized to make it valid.
func MakeLogicalIDSafe(id string) string {
	if IsValidLogicalID(id) {
		return id
	}
	return sanitizeLogicalID(id)
}
