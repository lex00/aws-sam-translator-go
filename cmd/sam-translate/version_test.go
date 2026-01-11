package main

import (
	"strings"
	"testing"
)

// TestGetVersion tests the version resolution logic.
func TestGetVersion(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	t.Run("returns ldflags version when set", func(t *testing.T) {
		Version = "v1.2.3"
		got := getVersion()
		if got != "v1.2.3" {
			t.Errorf("getVersion() = %q, want %q", got, "v1.2.3")
		}
	})

	t.Run("returns dev or module version when Version is dev", func(t *testing.T) {
		Version = "dev"
		got := getVersion()
		// Should return either "dev" or a valid version from build info
		// In test environment, build info may return "(devel)" or actual version
		if got == "" {
			t.Error("getVersion() returned empty string")
		}
	})

	t.Run("version format is valid", func(t *testing.T) {
		Version = "dev"
		got := getVersion()
		// Version should be either "dev" or a semver-like string
		if got != "dev" && !strings.HasPrefix(got, "v") && !strings.Contains(got, ".") {
			// Allow module versions which may not start with 'v'
			if got == "(devel)" {
				// This is acceptable in development
				return
			}
			t.Errorf("getVersion() = %q, expected 'dev' or version string", got)
		}
	})
}

// TestVersionNotEmpty ensures version is never empty.
func TestVersionNotEmpty(t *testing.T) {
	// Save original version
	originalVersion := Version
	defer func() { Version = originalVersion }()

	testCases := []string{"dev", "v1.0.0", "1.0.0", "custom-version"}
	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			Version = tc
			got := getVersion()
			if got == "" {
				t.Errorf("getVersion() with Version=%q returned empty string", tc)
			}
		})
	}
}
