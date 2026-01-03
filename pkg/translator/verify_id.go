// Package translator provides the main SAM to CloudFormation transformation orchestrator.
package translator

import (
	"fmt"
	"regexp"
	"strings"
)

// IDVerificationError represents an error during ID verification.
type IDVerificationError struct {
	ID      string
	Message string
}

func (e *IDVerificationError) Error() string {
	return fmt.Sprintf("invalid ID '%s': %s", e.ID, e.Message)
}

// IDVerifier provides validation and verification for CloudFormation logical IDs.
type IDVerifier struct {
	// knownIDs tracks all known logical IDs to detect duplicates.
	knownIDs map[string]bool
	// reservedPrefixes contains prefixes that shouldn't be used for user IDs.
	reservedPrefixes []string
}

// NewIDVerifier creates a new IDVerifier.
func NewIDVerifier() *IDVerifier {
	return &IDVerifier{
		knownIDs: make(map[string]bool),
		reservedPrefixes: []string{
			"AWS",
			"Custom",
		},
	}
}

// VerifyLogicalID validates a CloudFormation logical ID and tracks it for duplicate detection.
func (v *IDVerifier) VerifyLogicalID(id string) error {
	// Check for empty ID
	if id == "" {
		return &IDVerificationError{ID: id, Message: "logical ID cannot be empty"}
	}

	// Check length
	if len(id) > LogicalIDMaxLength {
		return &IDVerificationError{
			ID:      id,
			Message: fmt.Sprintf("logical ID exceeds maximum length of %d characters", LogicalIDMaxLength),
		}
	}

	// Check format (must be alphanumeric, starting with a letter)
	if !IsValidLogicalID(id) {
		return &IDVerificationError{
			ID:      id,
			Message: "logical ID must be alphanumeric and start with a letter",
		}
	}

	// Check for reserved prefixes
	for _, prefix := range v.reservedPrefixes {
		if strings.HasPrefix(id, prefix) {
			return &IDVerificationError{
				ID:      id,
				Message: fmt.Sprintf("logical ID cannot start with reserved prefix '%s'", prefix),
			}
		}
	}

	// Check for duplicates
	if v.knownIDs[id] {
		return &IDVerificationError{
			ID:      id,
			Message: "duplicate logical ID detected",
		}
	}

	// Track this ID
	v.knownIDs[id] = true

	return nil
}

// VerifyLogicalIDWithoutTracking validates a logical ID without adding it to the known IDs.
// Useful for validation without side effects.
func (v *IDVerifier) VerifyLogicalIDWithoutTracking(id string) error {
	// Check for empty ID
	if id == "" {
		return &IDVerificationError{ID: id, Message: "logical ID cannot be empty"}
	}

	// Check length
	if len(id) > LogicalIDMaxLength {
		return &IDVerificationError{
			ID:      id,
			Message: fmt.Sprintf("logical ID exceeds maximum length of %d characters", LogicalIDMaxLength),
		}
	}

	// Check format
	if !IsValidLogicalID(id) {
		return &IDVerificationError{
			ID:      id,
			Message: "logical ID must be alphanumeric and start with a letter",
		}
	}

	return nil
}

// RegisterID adds an ID to the known IDs set without full verification.
// Useful for tracking existing IDs from a template.
func (v *IDVerifier) RegisterID(id string) {
	v.knownIDs[id] = true
}

// UnregisterID removes an ID from the known IDs set.
func (v *IDVerifier) UnregisterID(id string) {
	delete(v.knownIDs, id)
}

// IsKnown checks if an ID has been registered.
func (v *IDVerifier) IsKnown(id string) bool {
	return v.knownIDs[id]
}

// GetKnownIDs returns a copy of all known IDs.
func (v *IDVerifier) GetKnownIDs() []string {
	ids := make([]string, 0, len(v.knownIDs))
	for id := range v.knownIDs {
		ids = append(ids, id)
	}
	return ids
}

// Clear removes all known IDs.
func (v *IDVerifier) Clear() {
	v.knownIDs = make(map[string]bool)
}

// AddReservedPrefix adds a custom reserved prefix.
func (v *IDVerifier) AddReservedPrefix(prefix string) {
	v.reservedPrefixes = append(v.reservedPrefixes, prefix)
}

// VerifyARN validates an ARN string.
func VerifyARN(arnStr string) error {
	if arnStr == "" {
		return fmt.Errorf("ARN cannot be empty")
	}

	arn, err := ParseARN(arnStr)
	if err != nil {
		return err
	}

	// Verify partition
	if !isValidPartition(arn.Partition) {
		return fmt.Errorf("invalid partition '%s' in ARN", arn.Partition)
	}

	// Verify service (must not be empty)
	if arn.Service == "" {
		return fmt.Errorf("service component cannot be empty in ARN")
	}

	// Verify resource (must not be empty)
	if arn.Resource == "" {
		return fmt.Errorf("resource component cannot be empty in ARN")
	}

	return nil
}

// VerifyARNPartition checks if an ARN uses a specific partition.
func VerifyARNPartition(arnStr, expectedPartition string) error {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return err
	}

	if arn.Partition != expectedPartition {
		return fmt.Errorf("ARN partition '%s' does not match expected partition '%s'", arn.Partition, expectedPartition)
	}

	return nil
}

// VerifyARNService checks if an ARN is for a specific service.
func VerifyARNService(arnStr, expectedService string) error {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return err
	}

	if arn.Service != expectedService {
		return fmt.Errorf("ARN service '%s' does not match expected service '%s'", arn.Service, expectedService)
	}

	return nil
}

// VerifyARNRegion checks if an ARN uses a specific region.
func VerifyARNRegion(arnStr, expectedRegion string) error {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return err
	}

	if arn.Region != expectedRegion {
		return fmt.Errorf("ARN region '%s' does not match expected region '%s'", arn.Region, expectedRegion)
	}

	return nil
}

// VerifyARNAccountID checks if an ARN uses a specific account ID.
func VerifyARNAccountID(arnStr, expectedAccountID string) error {
	arn, err := ParseARN(arnStr)
	if err != nil {
		return err
	}

	if arn.AccountID != expectedAccountID {
		return fmt.Errorf("ARN account ID '%s' does not match expected account ID '%s'", arn.AccountID, expectedAccountID)
	}

	return nil
}

// isValidPartition checks if a partition string is valid.
func isValidPartition(partition string) bool {
	validPartitions := map[string]bool{
		"aws":        true,
		"aws-cn":     true,
		"aws-us-gov": true,
	}
	return validPartitions[partition]
}

// IDStabilityChecker verifies that generated IDs are deterministic.
type IDStabilityChecker struct {
	generator *LogicalIDGenerator
	cache     map[string]string
}

// NewIDStabilityChecker creates a new IDStabilityChecker.
func NewIDStabilityChecker() *IDStabilityChecker {
	return &IDStabilityChecker{
		generator: NewLogicalIDGenerator(),
		cache:     make(map[string]string),
	}
}

// CheckStability verifies that the same inputs produce the same ID.
func (c *IDStabilityChecker) CheckStability(key string, parts ...string) (string, error) {
	id := c.generator.Generate(parts...)

	if cachedID, exists := c.cache[key]; exists {
		if cachedID != id {
			return "", fmt.Errorf("ID stability violation: key '%s' generated '%s' but previously generated '%s'", key, id, cachedID)
		}
	} else {
		c.cache[key] = id
	}

	return id, nil
}

// CheckHashedStability verifies that hashed IDs are stable.
func (c *IDStabilityChecker) CheckHashedStability(key, data string, parts ...string) (string, error) {
	id := c.generator.GenerateHashed(data, parts...)

	cacheKey := key + ":" + data
	if cachedID, exists := c.cache[cacheKey]; exists {
		if cachedID != id {
			return "", fmt.Errorf("hashed ID stability violation: key '%s' generated '%s' but previously generated '%s'", cacheKey, id, cachedID)
		}
	} else {
		c.cache[cacheKey] = id
	}

	return id, nil
}

// CheckAPIDeploymentIDStability verifies API deployment ID stability.
func (c *IDStabilityChecker) CheckAPIDeploymentIDStability(apiLogicalID, openAPISpec string) (string, error) {
	id := c.generator.GenerateAPIDeploymentID(apiLogicalID, openAPISpec)

	cacheKey := "api:" + apiLogicalID + ":" + calculateHash(openAPISpec)
	if cachedID, exists := c.cache[cacheKey]; exists {
		if cachedID != id {
			return "", fmt.Errorf("API deployment ID stability violation: generated '%s' but previously generated '%s'", id, cachedID)
		}
	} else {
		c.cache[cacheKey] = id
	}

	return id, nil
}

// VerifyAPIDeploymentIDChanges verifies that API deployment ID changes when spec changes.
func (c *IDStabilityChecker) VerifyAPIDeploymentIDChanges(apiLogicalID, oldSpec, newSpec string) error {
	oldID := c.generator.GenerateAPIDeploymentID(apiLogicalID, oldSpec)
	newID := c.generator.GenerateAPIDeploymentID(apiLogicalID, newSpec)

	if oldSpec != newSpec && oldID == newID {
		return fmt.Errorf("API deployment ID did not change when spec changed")
	}

	if oldSpec == newSpec && oldID != newID {
		return fmt.Errorf("API deployment ID changed when spec did not change")
	}

	return nil
}

// Clear resets the stability cache.
func (c *IDStabilityChecker) Clear() {
	c.cache = make(map[string]string)
}

// ValidateResourceName validates a resource name according to CloudFormation rules.
var resourceNamePattern = regexp.MustCompile(`^[a-zA-Z][-a-zA-Z0-9]*$`)

// ValidateResourceName checks if a name is valid for CloudFormation resources.
func ValidateResourceName(name string) error {
	if name == "" {
		return fmt.Errorf("resource name cannot be empty")
	}

	if len(name) > 128 {
		return fmt.Errorf("resource name exceeds maximum length of 128 characters")
	}

	if !resourceNamePattern.MatchString(name) {
		return fmt.Errorf("resource name must start with a letter and contain only alphanumeric characters and hyphens")
	}

	return nil
}

// ValidateStackName validates a CloudFormation stack name.
func ValidateStackName(name string) error {
	if name == "" {
		return fmt.Errorf("stack name cannot be empty")
	}

	if len(name) > 128 {
		return fmt.Errorf("stack name exceeds maximum length of 128 characters")
	}

	stackNamePattern := regexp.MustCompile(`^[a-zA-Z][-a-zA-Z0-9]*$`)
	if !stackNamePattern.MatchString(name) {
		return fmt.Errorf("stack name must start with a letter and contain only alphanumeric characters and hyphens")
	}

	return nil
}
