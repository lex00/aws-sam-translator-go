package translator

import (
	"strings"
	"testing"
)

func TestIDVerifier_VerifyLogicalID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid ID",
			id:      "MyFunction",
			wantErr: false,
		},
		{
			name:    "empty ID",
			id:      "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "starts with number",
			id:      "123Function",
			wantErr: true,
			errMsg:  "must be alphanumeric and start with a letter",
		},
		{
			name:    "contains hyphen",
			id:      "My-Function",
			wantErr: true,
			errMsg:  "must be alphanumeric and start with a letter",
		},
		{
			name:    "reserved prefix AWS",
			id:      "AWSMyFunction",
			wantErr: true,
			errMsg:  "reserved prefix",
		},
		{
			name:    "reserved prefix Custom",
			id:      "CustomMyFunction",
			wantErr: true,
			errMsg:  "reserved prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh verifier for each test to avoid duplicate detection
			verifier := NewIDVerifier()
			err := verifier.VerifyLogicalID(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Errorf("VerifyLogicalID(%q) expected error, got nil", tt.id)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("VerifyLogicalID(%q) error = %q, want error containing %q", tt.id, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("VerifyLogicalID(%q) unexpected error: %v", tt.id, err)
				}
			}
		})
	}
}

func TestIDVerifier_DuplicateDetection(t *testing.T) {
	v := NewIDVerifier()

	// First verification should succeed
	err := v.VerifyLogicalID("MyFunction")
	if err != nil {
		t.Errorf("First verification should succeed: %v", err)
	}

	// Second verification of same ID should fail
	err = v.VerifyLogicalID("MyFunction")
	if err == nil {
		t.Error("Duplicate verification should fail")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("Error should mention duplicate: %v", err)
	}
}

func TestIDVerifier_VerifyLogicalIDWithoutTracking(t *testing.T) {
	v := NewIDVerifier()

	// Verification without tracking should succeed multiple times
	err := v.VerifyLogicalIDWithoutTracking("MyFunction")
	if err != nil {
		t.Errorf("First verification should succeed: %v", err)
	}

	err = v.VerifyLogicalIDWithoutTracking("MyFunction")
	if err != nil {
		t.Errorf("Second verification without tracking should also succeed: %v", err)
	}

	// ID should not be in known IDs
	if v.IsKnown("MyFunction") {
		t.Error("ID should not be tracked")
	}
}

func TestIDVerifier_RegisterAndUnregister(t *testing.T) {
	v := NewIDVerifier()

	// Register an ID
	v.RegisterID("MyFunction")
	if !v.IsKnown("MyFunction") {
		t.Error("ID should be known after registration")
	}

	// Unregister the ID
	v.UnregisterID("MyFunction")
	if v.IsKnown("MyFunction") {
		t.Error("ID should not be known after unregistration")
	}
}

func TestIDVerifier_GetKnownIDs(t *testing.T) {
	v := NewIDVerifier()

	v.RegisterID("Function1")
	v.RegisterID("Function2")
	v.RegisterID("Function3")

	ids := v.GetKnownIDs()
	if len(ids) != 3 {
		t.Errorf("Expected 3 known IDs, got %d", len(ids))
	}

	// Check all IDs are present
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}
	for _, expected := range []string{"Function1", "Function2", "Function3"} {
		if !idMap[expected] {
			t.Errorf("Expected ID %q not found in known IDs", expected)
		}
	}
}

func TestIDVerifier_Clear(t *testing.T) {
	v := NewIDVerifier()

	v.RegisterID("Function1")
	v.RegisterID("Function2")

	v.Clear()

	if len(v.GetKnownIDs()) != 0 {
		t.Error("Known IDs should be empty after Clear()")
	}
}

func TestIDVerifier_AddReservedPrefix(t *testing.T) {
	v := NewIDVerifier()

	// Add custom reserved prefix
	v.AddReservedPrefix("Internal")

	err := v.VerifyLogicalID("InternalFunction")
	if err == nil {
		t.Error("ID with custom reserved prefix should fail verification")
	}
	if !strings.Contains(err.Error(), "reserved prefix") {
		t.Errorf("Error should mention reserved prefix: %v", err)
	}
}

func TestIDVerifier_LongID(t *testing.T) {
	v := NewIDVerifier()

	// Create an ID that's too long
	longID := strings.Repeat("a", LogicalIDMaxLength+1)

	err := v.VerifyLogicalID(longID)
	if err == nil {
		t.Error("Long ID should fail verification")
	}
	if !strings.Contains(err.Error(), "exceeds maximum length") {
		t.Errorf("Error should mention length: %v", err)
	}
}

func TestVerifyARN(t *testing.T) {
	tests := []struct {
		name    string
		arn     string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid Lambda ARN",
			arn:     "arn:aws:lambda:us-east-1:123456789012:function:MyFunction",
			wantErr: false,
		},
		{
			name:    "valid S3 ARN",
			arn:     "arn:aws:s3:::my-bucket",
			wantErr: false,
		},
		{
			name:    "empty ARN",
			arn:     "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "invalid partition",
			arn:     "arn:invalid:lambda:us-east-1:123456789012:function:MyFunction",
			wantErr: true,
			errMsg:  "invalid partition",
		},
		{
			name:    "invalid format",
			arn:     "not-an-arn",
			wantErr: true,
			errMsg:  "invalid ARN format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyARN(tt.arn)
			if tt.wantErr {
				if err == nil {
					t.Errorf("VerifyARN(%q) expected error, got nil", tt.arn)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("VerifyARN(%q) error = %q, want error containing %q", tt.arn, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("VerifyARN(%q) unexpected error: %v", tt.arn, err)
				}
			}
		})
	}
}

func TestVerifyARNPartition(t *testing.T) {
	arn := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction"

	// Should pass with correct partition
	err := VerifyARNPartition(arn, "aws")
	if err != nil {
		t.Errorf("VerifyARNPartition with correct partition should succeed: %v", err)
	}

	// Should fail with incorrect partition
	err = VerifyARNPartition(arn, "aws-cn")
	if err == nil {
		t.Error("VerifyARNPartition with incorrect partition should fail")
	}
}

func TestVerifyARNService(t *testing.T) {
	arn := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction"

	// Should pass with correct service
	err := VerifyARNService(arn, "lambda")
	if err != nil {
		t.Errorf("VerifyARNService with correct service should succeed: %v", err)
	}

	// Should fail with incorrect service
	err = VerifyARNService(arn, "s3")
	if err == nil {
		t.Error("VerifyARNService with incorrect service should fail")
	}
}

func TestVerifyARNRegion(t *testing.T) {
	arn := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction"

	// Should pass with correct region
	err := VerifyARNRegion(arn, "us-east-1")
	if err != nil {
		t.Errorf("VerifyARNRegion with correct region should succeed: %v", err)
	}

	// Should fail with incorrect region
	err = VerifyARNRegion(arn, "eu-west-1")
	if err == nil {
		t.Error("VerifyARNRegion with incorrect region should fail")
	}
}

func TestVerifyARNAccountID(t *testing.T) {
	arn := "arn:aws:lambda:us-east-1:123456789012:function:MyFunction"

	// Should pass with correct account ID
	err := VerifyARNAccountID(arn, "123456789012")
	if err != nil {
		t.Errorf("VerifyARNAccountID with correct account should succeed: %v", err)
	}

	// Should fail with incorrect account ID
	err = VerifyARNAccountID(arn, "999999999999")
	if err == nil {
		t.Error("VerifyARNAccountID with incorrect account should fail")
	}
}

func TestIsValidPartition(t *testing.T) {
	tests := []struct {
		partition string
		valid     bool
	}{
		{"aws", true},
		{"aws-cn", true},
		{"aws-us-gov", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.partition, func(t *testing.T) {
			if isValidPartition(tt.partition) != tt.valid {
				t.Errorf("isValidPartition(%q) = %v, want %v", tt.partition, !tt.valid, tt.valid)
			}
		})
	}
}

func TestIDStabilityChecker_CheckStability(t *testing.T) {
	checker := NewIDStabilityChecker()

	// First check should succeed and cache the ID
	id1, err := checker.CheckStability("test-key", "Part1", "Part2")
	if err != nil {
		t.Errorf("First CheckStability should succeed: %v", err)
	}

	// Second check with same key should return same ID
	id2, err := checker.CheckStability("test-key", "Part1", "Part2")
	if err != nil {
		t.Errorf("Second CheckStability with same inputs should succeed: %v", err)
	}
	if id1 != id2 {
		t.Errorf("IDs should be stable: %q != %q", id1, id2)
	}
}

func TestIDStabilityChecker_CheckHashedStability(t *testing.T) {
	checker := NewIDStabilityChecker()

	// Check hashed stability
	id1, err := checker.CheckHashedStability("test-key", "test-data", "Prefix")
	if err != nil {
		t.Errorf("First CheckHashedStability should succeed: %v", err)
	}

	id2, err := checker.CheckHashedStability("test-key", "test-data", "Prefix")
	if err != nil {
		t.Errorf("Second CheckHashedStability should succeed: %v", err)
	}
	if id1 != id2 {
		t.Errorf("Hashed IDs should be stable: %q != %q", id1, id2)
	}
}

func TestIDStabilityChecker_CheckAPIDeploymentIDStability(t *testing.T) {
	checker := NewIDStabilityChecker()

	spec := `{"openapi": "3.0"}`

	id1, err := checker.CheckAPIDeploymentIDStability("MyAPI", spec)
	if err != nil {
		t.Errorf("First CheckAPIDeploymentIDStability should succeed: %v", err)
	}

	id2, err := checker.CheckAPIDeploymentIDStability("MyAPI", spec)
	if err != nil {
		t.Errorf("Second CheckAPIDeploymentIDStability should succeed: %v", err)
	}
	if id1 != id2 {
		t.Errorf("API deployment IDs should be stable: %q != %q", id1, id2)
	}
}

func TestIDStabilityChecker_VerifyAPIDeploymentIDChanges(t *testing.T) {
	checker := NewIDStabilityChecker()

	oldSpec := `{"openapi": "3.0"}`
	newSpec := `{"openapi": "3.1"}`

	// Different specs should produce different IDs
	err := checker.VerifyAPIDeploymentIDChanges("MyAPI", oldSpec, newSpec)
	if err != nil {
		t.Errorf("Different specs should produce different IDs: %v", err)
	}

	// Same specs should produce same IDs
	err = checker.VerifyAPIDeploymentIDChanges("MyAPI", oldSpec, oldSpec)
	if err != nil {
		t.Errorf("Same specs should produce same IDs: %v", err)
	}
}

func TestIDStabilityChecker_Clear(t *testing.T) {
	checker := NewIDStabilityChecker()

	checker.CheckStability("test-key", "Part1")
	checker.Clear()

	// After clear, cache should be empty, new values should be accepted
	_, err := checker.CheckStability("test-key", "DifferentPart")
	if err != nil {
		t.Errorf("After Clear, new values should be accepted: %v", err)
	}
}

func TestValidateResourceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my-resource", false},
		{"valid alphanumeric", "MyResource123", false},
		{"empty name", "", true},
		{"starts with number", "123resource", true},
		{"starts with hyphen", "-resource", true},
		{"too long", strings.Repeat("a", 129), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResourceName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateResourceName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateStackName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "my-stack", false},
		{"valid alphanumeric", "MyStack123", false},
		{"empty name", "", true},
		{"starts with number", "123stack", true},
		{"starts with hyphen", "-stack", true},
		{"too long", strings.Repeat("a", 129), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStackName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStackName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestIDVerificationError(t *testing.T) {
	err := &IDVerificationError{
		ID:      "TestID",
		Message: "test error message",
	}

	expected := "invalid ID 'TestID': test error message"
	if err.Error() != expected {
		t.Errorf("IDVerificationError.Error() = %q, want %q", err.Error(), expected)
	}
}
