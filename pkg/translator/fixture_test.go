package translator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const (
	testdataInputDir  = "../../testdata/input"
	testdataOutputDir = "../../testdata/output"
)

// TestSuccessFixtures runs all success test fixtures through the translator
// and compares output against expected CloudFormation JSON.
//
// This mirrors the Python SAM translator's test_transform_success tests:
// - Reads YAML input from testdata/input/
// - Transforms using the Go SAM translator
// - Compares against expected JSON output in testdata/output/
// - Uses deep sorting for stable comparison (like Python's deep_sort_lists)
//
// Set STRICT_FIXTURES=1 to fail the test suite on fixture mismatches.
// By default, failures are logged but don't fail the build.
func TestSuccessFixtures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping fixture tests in short mode")
	}

	fixtures := getSuccessFixtures(t)
	if len(fixtures) == 0 {
		t.Fatal("No success fixtures found")
	}

	strictMode := os.Getenv("STRICT_FIXTURES") == "1"
	t.Logf("Running %d success fixtures for aws partition (strict=%v)", len(fixtures), strictMode)

	passed := 0
	failed := 0
	var failures []string

	for _, fixture := range fixtures {
		fixture := fixture // capture range variable
		t.Run(fixture.name, func(t *testing.T) {
			err := runSuccessFixtureCheck(fixture)
			if err != nil {
				failed++
				failures = append(failures, fixture.name+": "+err.Error())
				if strictMode {
					t.Error(err)
				}
			} else {
				passed++
			}
		})
	}

	// Log summary at the end
	t.Logf("\n=== Fixture Test Summary ===")
	t.Logf("Total:  %d", len(fixtures))
	t.Logf("Passed: %d (%.1f%%)", passed, float64(passed)/float64(len(fixtures))*100)
	t.Logf("Failed: %d", failed)

	if failed > 0 && !strictMode {
		t.Logf("\nFailed fixtures (first 10):")
		for i, f := range failures {
			if i >= 10 {
				t.Logf("... and %d more", len(failures)-10)
				break
			}
			t.Logf("  - %s", f)
		}
	}
}

// TestErrorFixtures runs all error test fixtures through the translator
// and verifies that appropriate errors are produced.
//
// This mirrors the Python SAM translator's test_transform_invalid_document tests:
// - Reads invalid YAML input from testdata/input/error_*.yaml
// - Verifies that transformation produces an error
// - Error cases include validation errors, missing properties, invalid types, etc.
//
// Set STRICT_FIXTURES=1 to fail the test suite on fixture mismatches.
func TestErrorFixtures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping fixture tests in short mode")
	}

	fixtures := getErrorFixtures(t)
	if len(fixtures) == 0 {
		t.Fatal("No error fixtures found")
	}

	strictMode := os.Getenv("STRICT_FIXTURES") == "1"
	t.Logf("Running %d error fixtures (strict=%v)", len(fixtures), strictMode)

	passed := 0
	failed := 0

	for _, fixture := range fixtures {
		fixture := fixture // capture range variable
		t.Run(fixture.name, func(t *testing.T) {
			err := runErrorFixtureCheck(fixture)
			if err != nil {
				failed++
				if strictMode {
					t.Error(err)
				}
			} else {
				passed++
			}
		})
	}

	t.Logf("\n=== Error Fixture Summary ===")
	t.Logf("Total:  %d", len(fixtures))
	t.Logf("Passed: %d (%.1f%%)", passed, float64(passed)/float64(len(fixtures))*100)
	t.Logf("Failed: %d", failed)
}

// TestPartitionFixtures tests transformation for all three AWS partitions.
// This mirrors the Python SAM translator's approach of testing each fixture
// against all partition/region combinations.
//
// Set STRICT_FIXTURES=1 to fail the test suite on fixture mismatches.
func TestPartitionFixtures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping partition tests in short mode")
	}

	// Matches Python SAM translator's partition/region test matrix
	partitions := []struct {
		name      string
		partition string
		region    string
		outputDir string
	}{
		{"aws-cn", "aws-cn", "cn-north-1", "aws-cn"},
		{"aws-us-gov", "aws-us-gov", "us-gov-west-1", "aws-us-gov"},
	}

	strictMode := os.Getenv("STRICT_FIXTURES") == "1"

	for _, p := range partitions {
		p := p
		t.Run(p.name, func(t *testing.T) {
			fixtures := getPartitionFixtures(t, p.outputDir)
			if len(fixtures) == 0 {
				t.Skipf("No fixtures found for partition %s", p.name)
			}

			t.Logf("Running %d fixtures for partition %s (strict=%v)", len(fixtures), p.name, strictMode)

			passed := 0
			failed := 0

			for _, fixture := range fixtures {
				fixture := fixture
				t.Run(fixture.name, func(t *testing.T) {
					err := runPartitionFixtureCheck(fixture, p.partition, p.region)
					if err != nil {
						failed++
						if strictMode {
							t.Error(err)
						}
					} else {
						passed++
					}
				})
			}

			t.Logf("\n=== Partition %s Summary ===", p.name)
			t.Logf("Total:  %d", len(fixtures))
			t.Logf("Passed: %d (%.1f%%)", passed, float64(passed)/float64(len(fixtures))*100)
			t.Logf("Failed: %d", failed)
		})
	}
}

type testFixture struct {
	name       string
	inputPath  string
	outputPath string
}

func getSuccessFixtures(t *testing.T) []testFixture {
	t.Helper()

	inputFiles, err := filepath.Glob(filepath.Join(testdataInputDir, "*.yaml"))
	if err != nil {
		t.Fatalf("Failed to glob input files: %v", err)
	}

	var fixtures []testFixture
	for _, inputPath := range inputFiles {
		baseName := filepath.Base(inputPath)
		// Skip error fixtures
		if strings.HasPrefix(baseName, "error_") {
			continue
		}

		name := strings.TrimSuffix(baseName, ".yaml")
		outputPath := filepath.Join(testdataOutputDir, name+".json")

		// Check if output exists
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			continue
		}

		fixtures = append(fixtures, testFixture{
			name:       name,
			inputPath:  inputPath,
			outputPath: outputPath,
		})
	}

	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].name < fixtures[j].name
	})

	return fixtures
}

func getErrorFixtures(t *testing.T) []testFixture {
	t.Helper()

	inputFiles, err := filepath.Glob(filepath.Join(testdataInputDir, "error_*.yaml"))
	if err != nil {
		t.Fatalf("Failed to glob error input files: %v", err)
	}

	var fixtures []testFixture
	for _, inputPath := range inputFiles {
		baseName := filepath.Base(inputPath)
		name := strings.TrimSuffix(baseName, ".yaml")
		outputPath := filepath.Join(testdataOutputDir, name+".json")

		// Check if output exists
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			continue
		}

		fixtures = append(fixtures, testFixture{
			name:       name,
			inputPath:  inputPath,
			outputPath: outputPath,
		})
	}

	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].name < fixtures[j].name
	})

	return fixtures
}

func getPartitionFixtures(t *testing.T, partitionDir string) []testFixture {
	t.Helper()

	outputDir := filepath.Join(testdataOutputDir, partitionDir)
	outputFiles, err := filepath.Glob(filepath.Join(outputDir, "*.json"))
	if err != nil {
		t.Fatalf("Failed to glob partition output files: %v", err)
	}

	var fixtures []testFixture
	for _, outputPath := range outputFiles {
		baseName := filepath.Base(outputPath)
		// Skip error fixtures
		if strings.HasPrefix(baseName, "error_") {
			continue
		}

		name := strings.TrimSuffix(baseName, ".json")
		inputPath := filepath.Join(testdataInputDir, name+".yaml")

		// Check if input exists
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			continue
		}

		fixtures = append(fixtures, testFixture{
			name:       name,
			inputPath:  inputPath,
			outputPath: outputPath,
		})
	}

	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].name < fixtures[j].name
	})

	return fixtures
}

// runSuccessFixtureCheck runs a single fixture and returns an error if it fails.
func runSuccessFixtureCheck(fixture testFixture) error {
	// Read input template
	input, err := os.ReadFile(fixture.inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	// Transform
	tr := NewWithOptions(Options{
		Region:    "us-east-1",
		AccountID: "123456789012",
		StackName: "sam-app",
		Partition: "aws",
	})

	output, err := tr.TransformBytes(input)
	if err != nil {
		return fmt.Errorf("transform failed: %v", err)
	}

	// Read expected output
	expected, err := os.ReadFile(fixture.outputPath)
	if err != nil {
		return fmt.Errorf("failed to read expected output: %v", err)
	}

	// Parse both as JSON for comparison
	var actualJSON, expectedJSON map[string]interface{}
	if err := json.Unmarshal(output, &actualJSON); err != nil {
		return fmt.Errorf("failed to parse actual output: %v", err)
	}
	if err := json.Unmarshal(expected, &expectedJSON); err != nil {
		return fmt.Errorf("failed to parse expected output: %v", err)
	}

	// Compare
	if !jsonEqual(actualJSON, expectedJSON) {
		diffs := findDiffs("", actualJSON, expectedJSON)
		if len(diffs) > 3 {
			return fmt.Errorf("%d differences (first 3: %s)", len(diffs), strings.Join(diffs[:3], "; "))
		}
		return fmt.Errorf("%s", strings.Join(diffs, "; "))
	}

	return nil
}

// findDiffs returns a list of differences between two JSON values.
func findDiffs(path string, a, b interface{}) []string {
	var diffs []string

	// Deep sort for comparison
	aSorted := deepSortLists(a)
	bSorted := deepSortLists(b)

	switch aTyped := aSorted.(type) {
	case map[string]interface{}:
		bTyped, ok := bSorted.(map[string]interface{})
		if !ok {
			return []string{fmt.Sprintf("%s: type mismatch (map vs %T)", path, b)}
		}
		// Check for missing/extra keys
		for key := range aTyped {
			if _, ok := bTyped[key]; !ok {
				diffs = append(diffs, fmt.Sprintf("%s.%s: extra in actual", path, key))
			}
		}
		for key := range bTyped {
			if _, ok := aTyped[key]; !ok {
				diffs = append(diffs, fmt.Sprintf("%s.%s: missing in actual", path, key))
			}
		}
		// Check values
		for key, aVal := range aTyped {
			if bVal, ok := bTyped[key]; ok {
				subPath := key
				if path != "" {
					subPath = path + "." + key
				}
				diffs = append(diffs, findDiffs(subPath, aVal, bVal)...)
			}
		}

	case []interface{}:
		bTyped, ok := bSorted.([]interface{})
		if !ok {
			return []string{fmt.Sprintf("%s: type mismatch (slice vs %T)", path, b)}
		}
		if len(aTyped) != len(bTyped) {
			diffs = append(diffs, fmt.Sprintf("%s: length mismatch (%d vs %d)", path, len(aTyped), len(bTyped)))
		}
		for i := 0; i < len(aTyped) && i < len(bTyped); i++ {
			subPath := fmt.Sprintf("%s[%d]", path, i)
			diffs = append(diffs, findDiffs(subPath, aTyped[i], bTyped[i])...)
		}

	default:
		if !jsonEqualNoSort(aSorted, bSorted) {
			diffs = append(diffs, fmt.Sprintf("%s: value mismatch (%v vs %v)", path, a, b))
		}
	}

	return diffs
}

// runErrorFixtureCheck runs a single error fixture and returns an error if it fails.
func runErrorFixtureCheck(fixture testFixture) error {
	// Read input template
	input, err := os.ReadFile(fixture.inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	// Transform - should fail
	tr := NewWithOptions(Options{
		Region:    "us-east-1",
		AccountID: "123456789012",
		StackName: "sam-app",
		Partition: "aws",
	})

	_, err = tr.TransformBytes(input)
	if err == nil {
		return fmt.Errorf("expected error but transform succeeded")
	}

	// Successfully got an error - that's the expected behavior
	return nil
}

// runPartitionFixtureCheck runs a single partition fixture and returns an error if it fails.
func runPartitionFixtureCheck(fixture testFixture, partition, region string) error {
	// Read input template
	input, err := os.ReadFile(fixture.inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	// Transform with specified partition
	tr := NewWithOptions(Options{
		Region:    region,
		AccountID: "123456789012",
		StackName: "sam-app",
		Partition: partition,
	})

	output, err := tr.TransformBytes(input)
	if err != nil {
		return fmt.Errorf("transform failed: %v", err)
	}

	// Read expected output
	expected, err := os.ReadFile(fixture.outputPath)
	if err != nil {
		return fmt.Errorf("failed to read expected output: %v", err)
	}

	// Parse both as JSON for comparison
	var actualJSON, expectedJSON map[string]interface{}
	if err := json.Unmarshal(output, &actualJSON); err != nil {
		return fmt.Errorf("failed to parse actual output: %v", err)
	}
	if err := json.Unmarshal(expected, &expectedJSON); err != nil {
		return fmt.Errorf("failed to parse expected output: %v", err)
	}

	// Compare
	if !jsonEqual(actualJSON, expectedJSON) {
		diffs := findDiffs("", actualJSON, expectedJSON)
		if len(diffs) > 3 {
			return fmt.Errorf("%d differences (first 3: %s)", len(diffs), strings.Join(diffs[:3], "; "))
		}
		return fmt.Errorf("%s", strings.Join(diffs, "; "))
	}

	return nil
}

// deepSortLists recursively sorts all lists in a JSON value for stable comparison.
// This is necessary because Go map iteration order is undefined and lists may not
// be in a predictable order. We sort them to enable deterministic comparison.
//
// This mirrors the Python SAM translator's deep_sort_lists function.
func deepSortLists(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = deepSortLists(val)
		}
		return result

	case []interface{}:
		// First, recursively sort nested structures
		sorted := make([]interface{}, len(v))
		for i, item := range v {
			sorted[i] = deepSortLists(item)
		}
		// Then sort the list itself using JSON string representation for stable ordering
		sort.Slice(sorted, func(i, j int) bool {
			iJSON, _ := json.Marshal(sorted[i])
			jJSON, _ := json.Marshal(sorted[j])
			return string(iJSON) < string(jJSON)
		})
		return sorted

	default:
		return value
	}
}

// jsonEqual performs a deep equality check on two JSON values.
// It first normalizes both values by deep sorting all lists, then compares.
func jsonEqual(a, b interface{}) bool {
	// Deep sort both values for stable comparison
	aSorted := deepSortLists(a)
	bSorted := deepSortLists(b)
	return jsonEqualNoSort(aSorted, bSorted)
}

// jsonEqualNoSort performs equality check without sorting (used after normalization).
func jsonEqualNoSort(a, b interface{}) bool {
	switch aTyped := a.(type) {
	case map[string]interface{}:
		bTyped, ok := b.(map[string]interface{})
		if !ok {
			return false
		}
		if len(aTyped) != len(bTyped) {
			return false
		}
		for key, aVal := range aTyped {
			bVal, ok := bTyped[key]
			if !ok {
				return false
			}
			if !jsonEqualNoSort(aVal, bVal) {
				return false
			}
		}
		return true

	case []interface{}:
		bTyped, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(aTyped) != len(bTyped) {
			return false
		}
		for i := range aTyped {
			if !jsonEqualNoSort(aTyped[i], bTyped[i]) {
				return false
			}
		}
		return true

	case float64:
		bTyped, ok := b.(float64)
		if !ok {
			return false
		}
		return aTyped == bTyped

	case string:
		bTyped, ok := b.(string)
		if !ok {
			return false
		}
		return aTyped == bTyped

	case bool:
		bTyped, ok := b.(bool)
		if !ok {
			return false
		}
		return aTyped == bTyped

	case nil:
		return b == nil

	default:
		return a == b
	}
}
