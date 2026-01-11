package main

import "runtime/debug"

// Version is set at build time via ldflags.
// When not set (default "dev"), the version will be inferred from
// Go module info when installed via `go install`.
var Version = "dev"

// getVersion returns the CLI version, attempting to read it from
// Go build info if the compile-time Version is "dev".
func getVersion() string {
	// If Version was set at build time via ldflags, use it
	if Version != "dev" {
		return Version
	}

	// Try to get version from Go module info (works with `go install pkg@version`)
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	// Fallback to "dev"
	return "dev"
}
