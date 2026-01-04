// Package main provides a tool for comparing Go translator output with Python aws-sam-translator.
//
// Usage:
//
//	go run tools/compare_output.go -input testdata/input/function_basic.yaml -partition aws
//	go run tools/compare_output.go -all                                   # Compare all fixtures
//	go run tools/compare_output.go -all -partition aws-cn                 # Compare with China partition
//
// This tool requires the Python aws-sam-translator to be installed:
//
//	pip install aws-sam-translator
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/translator"
)

type ComparisonResult struct {
	InputFile   string
	Partition   string
	GoOutput    map[string]interface{}
	PyOutput    map[string]interface{}
	Matches     bool
	Diffs       []string
	GoError     error
	PyError     error
	IsErrorCase bool
}

type ComparisonStats struct {
	Total    int
	Passed   int
	Failed   int
	GoErrors int
	PyErrors int
	Results  []ComparisonResult
}

func main() {
	inputFile := flag.String("input", "", "Single input file to compare")
	all := flag.Bool("all", false, "Compare all fixtures in testdata/input")
	partition := flag.String("partition", "aws", "AWS partition (aws, aws-cn, aws-us-gov)")
	verbose := flag.Bool("v", false, "Verbose output")
	showDiff := flag.Bool("diff", false, "Show detailed diffs for failures")
	flag.Parse()

	if *inputFile == "" && !*all {
		fmt.Println("Usage: compare_output -input <file> | -all [-partition aws|aws-cn|aws-us-gov]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check if Python aws-sam-translator is available
	if !checkPythonTranslator() {
		fmt.Println("Error: Python aws-sam-translator not found")
		fmt.Println("Install with: pip install aws-sam-translator")
		os.Exit(1)
	}

	var files []string
	if *all {
		var err error
		files, err = filepath.Glob("testdata/input/*.yaml")
		if err != nil {
			fmt.Printf("Error finding input files: %v\n", err)
			os.Exit(1)
		}
		sort.Strings(files)
	} else {
		files = []string{*inputFile}
	}

	stats := ComparisonStats{}
	for _, file := range files {
		result := compareFile(file, *partition, *verbose)
		stats.Total++
		stats.Results = append(stats.Results, result)

		if result.GoError != nil {
			stats.GoErrors++
		}
		if result.PyError != nil {
			stats.PyErrors++
		}
		if result.Matches {
			stats.Passed++
		} else {
			stats.Failed++
		}

		// Print result
		if *verbose || !result.Matches {
			printResult(result, *showDiff)
		} else if *all {
			if result.Matches {
				fmt.Printf(".")
			} else {
				fmt.Printf("F")
			}
		}
	}

	if *all {
		fmt.Println()
		printSummary(stats)
	}

	if stats.Failed > 0 {
		os.Exit(1)
	}
}

func checkPythonTranslator() bool {
	cmd := exec.Command("python3", "-c", "import samtranslator")
	return cmd.Run() == nil
}

func compareFile(inputFile, partition string, verbose bool) ComparisonResult {
	result := ComparisonResult{
		InputFile:   inputFile,
		Partition:   partition,
		IsErrorCase: strings.Contains(filepath.Base(inputFile), "error"),
	}

	// Read input template
	input, err := os.ReadFile(inputFile)
	if err != nil {
		result.GoError = fmt.Errorf("failed to read input: %v", err)
		return result
	}

	// Transform with Go translator
	goOutput, goErr := transformWithGo(input, partition)
	result.GoError = goErr
	if goErr == nil {
		result.GoOutput = goOutput
	}

	// Transform with Python translator
	pyOutput, pyErr := transformWithPython(inputFile, partition)
	result.PyError = pyErr
	if pyErr == nil {
		result.PyOutput = pyOutput
	}

	// For error cases, both should error
	if result.IsErrorCase {
		result.Matches = (goErr != nil && pyErr != nil) || compareErrorMessages(goErr, pyErr)
		if !result.Matches {
			result.Diffs = []string{fmt.Sprintf("Go error: %v, Python error: %v", goErr, pyErr)}
		}
		return result
	}

	// For success cases, compare outputs
	if goErr != nil || pyErr != nil {
		result.Matches = false
		if goErr != nil {
			result.Diffs = append(result.Diffs, fmt.Sprintf("Go error: %v", goErr))
		}
		if pyErr != nil {
			result.Diffs = append(result.Diffs, fmt.Sprintf("Python error: %v", pyErr))
		}
		return result
	}

	// Compare the JSON outputs
	result.Matches, result.Diffs = compareJSON(goOutput, pyOutput, "")
	return result
}

func transformWithGo(input []byte, partition string) (map[string]interface{}, error) {
	t := translator.NewWithOptions(translator.Options{
		Region:    getRegionForPartition(partition),
		AccountID: "123456789012",
		StackName: "sam-app",
		Partition: partition,
	})

	output, err := t.TransformBytes(input)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Go output: %v", err)
	}

	return result, nil
}

func transformWithPython(inputFile, partition string) (map[string]interface{}, error) {
	// Create a Python script to transform the template
	script := fmt.Sprintf(`
import json
import sys
from samtranslator.translator.translator import Translator
from samtranslator.parser import parser
from samtranslator.plugins import LifeCycleEvents
from samtranslator.policy_template_processor.processor import PolicyTemplatesProcessor

try:
    with open('%s', 'r') as f:
        manifest = parser.Parser().parse(f.read(), '%s')

    feature_toggle = None
    translator = Translator(manifest, manifest, feature_toggle)

    # Provide managed policy map and intrinsics resolver
    intrinsics_resolver = None
    policy_templates = PolicyTemplatesProcessor('.').process()

    output = translator.translate(
        intrinsics_resolver=intrinsics_resolver,
        policy_template_processor=policy_templates,
    )

    print(json.dumps(output, indent=2, sort_keys=True))
except Exception as e:
    print(json.dumps({"error": str(e)}), file=sys.stderr)
    sys.exit(1)
`, inputFile, partition)

	cmd := exec.Command("python3", "-c", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("python error: %s", stderr.String())
		}
		return nil, fmt.Errorf("python execution failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Python output: %v", err)
	}

	return result, nil
}

func getRegionForPartition(partition string) string {
	switch partition {
	case "aws-cn":
		return "cn-north-1"
	case "aws-us-gov":
		return "us-gov-west-1"
	default:
		return "us-east-1"
	}
}

func compareErrorMessages(goErr, pyErr error) bool {
	// If both are nil or both are non-nil, consider it a match for error cases
	if (goErr == nil) == (pyErr == nil) {
		return true
	}
	return false
}

func compareJSON(a, b map[string]interface{}, path string) (bool, []string) {
	var diffs []string

	// Get all keys from both maps
	keys := make(map[string]bool)
	for k := range a {
		keys[k] = true
	}
	for k := range b {
		keys[k] = true
	}

	var sortedKeys []string
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}

		aVal, aOk := a[key]
		bVal, bOk := b[key]

		if !aOk {
			diffs = append(diffs, fmt.Sprintf("Missing in Go: %s", currentPath))
			continue
		}
		if !bOk {
			diffs = append(diffs, fmt.Sprintf("Extra in Go: %s", currentPath))
			continue
		}

		// Compare values recursively
		matches, subDiffs := compareValues(aVal, bVal, currentPath)
		if !matches {
			diffs = append(diffs, subDiffs...)
		}
	}

	return len(diffs) == 0, diffs
}

func compareValues(a, b interface{}, path string) (bool, []string) {
	// Handle nil cases
	if a == nil && b == nil {
		return true, nil
	}
	if a == nil || b == nil {
		return false, []string{fmt.Sprintf("Mismatch at %s: Go=%v, Python=%v", path, a, b)}
	}

	// Compare based on type
	switch aTyped := a.(type) {
	case map[string]interface{}:
		bTyped, ok := b.(map[string]interface{})
		if !ok {
			return false, []string{fmt.Sprintf("Type mismatch at %s: Go=map, Python=%T", path, b)}
		}
		return compareJSON(aTyped, bTyped, path)

	case []interface{}:
		bTyped, ok := b.([]interface{})
		if !ok {
			return false, []string{fmt.Sprintf("Type mismatch at %s: Go=slice, Python=%T", path, b)}
		}
		if len(aTyped) != len(bTyped) {
			return false, []string{fmt.Sprintf("Array length mismatch at %s: Go=%d, Python=%d", path, len(aTyped), len(bTyped))}
		}
		var diffs []string
		for i := range aTyped {
			itemPath := fmt.Sprintf("%s[%d]", path, i)
			matches, subDiffs := compareValues(aTyped[i], bTyped[i], itemPath)
			if !matches {
				diffs = append(diffs, subDiffs...)
			}
		}
		return len(diffs) == 0, diffs

	case string:
		bTyped, ok := b.(string)
		if !ok {
			return false, []string{fmt.Sprintf("Type mismatch at %s: Go=string, Python=%T", path, b)}
		}
		if aTyped != bTyped {
			return false, []string{fmt.Sprintf("String mismatch at %s: Go=%q, Python=%q", path, aTyped, bTyped)}
		}
		return true, nil

	case float64:
		// JSON numbers are always float64
		bTyped, ok := b.(float64)
		if !ok {
			// Check if it's an int that matches
			if bInt, ok := b.(int); ok && float64(bInt) == aTyped {
				return true, nil
			}
			return false, []string{fmt.Sprintf("Type mismatch at %s: Go=number(%v), Python=%T(%v)", path, a, b, b)}
		}
		if aTyped != bTyped {
			return false, []string{fmt.Sprintf("Number mismatch at %s: Go=%v, Python=%v", path, aTyped, bTyped)}
		}
		return true, nil

	case bool:
		bTyped, ok := b.(bool)
		if !ok {
			return false, []string{fmt.Sprintf("Type mismatch at %s: Go=bool, Python=%T", path, b)}
		}
		if aTyped != bTyped {
			return false, []string{fmt.Sprintf("Bool mismatch at %s: Go=%v, Python=%v", path, aTyped, bTyped)}
		}
		return true, nil

	default:
		// For other types, use string comparison
		aStr := fmt.Sprintf("%v", a)
		bStr := fmt.Sprintf("%v", b)
		if aStr != bStr {
			return false, []string{fmt.Sprintf("Value mismatch at %s: Go=%v (%T), Python=%v (%T)", path, a, a, b, b)}
		}
		return true, nil
	}
}

func printResult(result ComparisonResult, showDiff bool) {
	status := "PASS"
	if !result.Matches {
		status = "FAIL"
	}

	fmt.Printf("\n%s: %s [%s]\n", status, result.InputFile, result.Partition)

	if result.GoError != nil {
		fmt.Printf("  Go Error: %v\n", result.GoError)
	}
	if result.PyError != nil {
		fmt.Printf("  Python Error: %v\n", result.PyError)
	}

	if showDiff && len(result.Diffs) > 0 {
		fmt.Println("  Differences:")
		for _, diff := range result.Diffs {
			fmt.Printf("    - %s\n", diff)
		}
	} else if len(result.Diffs) > 0 {
		fmt.Printf("  %d differences found\n", len(result.Diffs))
	}
}

func printSummary(stats ComparisonStats) {
	fmt.Println("\n=== Comparison Summary ===")
	fmt.Printf("Total:    %d\n", stats.Total)
	fmt.Printf("Passed:   %d (%.1f%%)\n", stats.Passed, float64(stats.Passed)/float64(stats.Total)*100)
	fmt.Printf("Failed:   %d\n", stats.Failed)
	fmt.Printf("Go Errors: %d\n", stats.GoErrors)
	fmt.Printf("Py Errors: %d\n", stats.PyErrors)

	if stats.Failed > 0 {
		fmt.Println("\nFailed fixtures:")
		for _, r := range stats.Results {
			if !r.Matches {
				fmt.Printf("  - %s\n", r.InputFile)
			}
		}
	}
}
