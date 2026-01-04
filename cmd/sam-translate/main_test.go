package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestParseArgs tests command-line argument parsing with cobra.
func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOpts Options
		wantErr  bool
	}{
		{
			name: "template file short flag",
			args: []string{"-t", "template.yaml"},
			wantOpts: Options{
				TemplateFile: "template.yaml",
			},
		},
		{
			name: "template file long flag",
			args: []string{"--template-file", "template.yaml"},
			wantOpts: Options{
				TemplateFile: "template.yaml",
			},
		},
		{
			name: "output template short flag",
			args: []string{"-t", "input.yaml", "-o", "output.yaml"},
			wantOpts: Options{
				TemplateFile:   "input.yaml",
				OutputTemplate: "output.yaml",
			},
		},
		{
			name: "output template long flag",
			args: []string{"-t", "input.yaml", "--output-template", "output.yaml"},
			wantOpts: Options{
				TemplateFile:   "input.yaml",
				OutputTemplate: "output.yaml",
			},
		},
		{
			name: "stdout flag",
			args: []string{"-t", "template.yaml", "--stdout"},
			wantOpts: Options{
				TemplateFile: "template.yaml",
				Stdout:       true,
			},
		},
		{
			name: "verbose flag",
			args: []string{"-t", "template.yaml", "--verbose"},
			wantOpts: Options{
				TemplateFile: "template.yaml",
				Verbose:      true,
			},
		},
		{
			name: "region flag",
			args: []string{"-t", "template.yaml", "--region", "us-west-2"},
			wantOpts: Options{
				TemplateFile: "template.yaml",
				Region:       "us-west-2",
			},
		},
		{
			name: "all flags combined",
			args: []string{"-t", "input.yaml", "-o", "output.yaml", "--region", "eu-west-1", "--verbose"},
			wantOpts: Options{
				TemplateFile:   "input.yaml",
				OutputTemplate: "output.yaml",
				Region:         "eu-west-1",
				Verbose:        true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRootCmd()
			cmd.SetArgs(tt.args)

			// Parse flags without executing
			if err := cmd.ParseFlags(tt.args); err != nil {
				if !tt.wantErr {
					t.Errorf("ParseFlags() error = %v", err)
				}
				return
			}

			// Get parsed values using cobra's flag methods
			opts := getOptionsFromCmd(cmd)

			if opts.TemplateFile != tt.wantOpts.TemplateFile {
				t.Errorf("TemplateFile = %q, want %q", opts.TemplateFile, tt.wantOpts.TemplateFile)
			}
			if opts.OutputTemplate != tt.wantOpts.OutputTemplate {
				t.Errorf("OutputTemplate = %q, want %q", opts.OutputTemplate, tt.wantOpts.OutputTemplate)
			}
			if opts.Stdout != tt.wantOpts.Stdout {
				t.Errorf("Stdout = %v, want %v", opts.Stdout, tt.wantOpts.Stdout)
			}
			if opts.Verbose != tt.wantOpts.Verbose {
				t.Errorf("Verbose = %v, want %v", opts.Verbose, tt.wantOpts.Verbose)
			}
			if opts.Region != tt.wantOpts.Region {
				t.Errorf("Region = %q, want %q", opts.Region, tt.wantOpts.Region)
			}
		})
	}
}

func getOptionsFromCmd(cmd *cobra.Command) Options {
	templateFile, _ := cmd.Flags().GetString("template-file")
	outputTemplate, _ := cmd.Flags().GetString("output-template")
	stdout, _ := cmd.Flags().GetBool("stdout")
	verbose, _ := cmd.Flags().GetBool("verbose")
	region, _ := cmd.Flags().GetString("region")

	return Options{
		TemplateFile:   templateFile,
		OutputTemplate: outputTemplate,
		Stdout:         stdout,
		Verbose:        verbose,
		Region:         region,
	}
}

// TestTransformCommand tests the transform command functionality.
func TestTransformCommand(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "sam-translate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test SAM template
	samTemplate := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: s3://bucket/key
`
	inputFile := filepath.Join(tmpDir, "template.yaml")
	if err := os.WriteFile(inputFile, []byte(samTemplate), 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	t.Run("transform to file", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "output.yaml")
		opts := &Options{
			TemplateFile:   inputFile,
			OutputTemplate: outputFile,
		}

		exitCode := runTransform(opts, nil, nil)
		if exitCode != ExitSuccess {
			t.Errorf("runTransform() returned %d, want %d", exitCode, ExitSuccess)
		}

		// Verify output file exists and is valid JSON
		output, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(output, &result); err != nil {
			t.Errorf("output is not valid JSON: %v", err)
		}

		// Verify the function was transformed
		resources, ok := result["Resources"].(map[string]interface{})
		if !ok {
			t.Fatal("Resources not found in output")
		}
		fn, ok := resources["MyFunction"].(map[string]interface{})
		if !ok {
			t.Fatal("MyFunction not found in output")
		}
		if fn["Type"] != "AWS::Lambda::Function" {
			t.Errorf("expected AWS::Lambda::Function, got %v", fn["Type"])
		}
	})

	t.Run("transform to stdout", func(t *testing.T) {
		var stdout bytes.Buffer
		opts := &Options{
			TemplateFile: inputFile,
			Stdout:       true,
		}

		exitCode := runTransform(opts, &stdout, nil)
		if exitCode != ExitSuccess {
			t.Errorf("runTransform() returned %d, want %d", exitCode, ExitSuccess)
		}

		// Verify stdout is valid JSON (trim trailing newline)
		var result map[string]interface{}
		if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &result); err != nil {
			t.Errorf("stdout is not valid JSON: %v", err)
		}
	})

	t.Run("transform with region", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "output-region.yaml")
		opts := &Options{
			TemplateFile:   inputFile,
			OutputTemplate: outputFile,
			Region:         "us-west-2",
		}

		exitCode := runTransform(opts, nil, nil)
		if exitCode != ExitSuccess {
			t.Errorf("runTransform() returned %d, want %d", exitCode, ExitSuccess)
		}
	})
}

// TestExitCodes tests that exit codes match the Python version.
func TestExitCodes(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "sam-translate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("success returns 0", func(t *testing.T) {
		samTemplate := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: s3://bucket/key
`
		inputFile := filepath.Join(tmpDir, "valid.yaml")
		if err := os.WriteFile(inputFile, []byte(samTemplate), 0644); err != nil {
			t.Fatalf("failed to write input file: %v", err)
		}

		var stdout bytes.Buffer
		opts := &Options{
			TemplateFile: inputFile,
			Stdout:       true,
		}
		exitCode := runTransform(opts, &stdout, nil)
		if exitCode != ExitSuccess {
			t.Errorf("exitCode = %d, want %d", exitCode, ExitSuccess)
		}
	})

	t.Run("transform error returns 1", func(t *testing.T) {
		// Template with invalid function (missing required properties)
		invalidTemplate := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  BadFunction:
    Type: AWS::Serverless::Function
    Properties: {}
`
		inputFile := filepath.Join(tmpDir, "invalid.yaml")
		if err := os.WriteFile(inputFile, []byte(invalidTemplate), 0644); err != nil {
			t.Fatalf("failed to write input file: %v", err)
		}

		var stdout, stderr bytes.Buffer
		opts := &Options{
			TemplateFile: inputFile,
			Stdout:       true,
		}
		exitCode := runTransform(opts, &stdout, &stderr)
		if exitCode != ExitTransformError {
			t.Errorf("exitCode = %d, want %d", exitCode, ExitTransformError)
		}
	})

	t.Run("file not found returns 1", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		opts := &Options{
			TemplateFile: "/nonexistent/template.yaml",
			Stdout:       true,
		}
		exitCode := runTransform(opts, &stdout, &stderr)
		if exitCode != ExitTransformError {
			t.Errorf("exitCode = %d, want %d", exitCode, ExitTransformError)
		}
	})
}

// TestErrorReportingWithSourceLocations tests error messages include line/column info.
func TestErrorReportingWithSourceLocations(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "sam-translate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Template with syntax error
	invalidYAML := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties
      Handler: index.handler
`
	inputFile := filepath.Join(tmpDir, "syntax-error.yaml")
	if err := os.WriteFile(inputFile, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	opts := &Options{
		TemplateFile: inputFile,
		Stdout:       true,
	}
	exitCode := runTransform(opts, &stdout, &stderr)
	if exitCode != ExitTransformError {
		t.Errorf("exitCode = %d, want %d", exitCode, ExitTransformError)
	}

	// Error message should be present
	errOutput := stderr.String()
	if errOutput == "" {
		t.Error("expected error output")
	}
}

// TestYAMLAndJSONSupport tests that both YAML and JSON input are supported.
func TestYAMLAndJSONSupport(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "sam-translate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("YAML input", func(t *testing.T) {
		samTemplate := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: s3://bucket/key
`
		inputFile := filepath.Join(tmpDir, "template.yaml")
		if err := os.WriteFile(inputFile, []byte(samTemplate), 0644); err != nil {
			t.Fatalf("failed to write input file: %v", err)
		}

		var stdout bytes.Buffer
		opts := &Options{
			TemplateFile: inputFile,
			Stdout:       true,
		}
		exitCode := runTransform(opts, &stdout, nil)
		if exitCode != ExitSuccess {
			t.Errorf("exitCode = %d, want %d", exitCode, ExitSuccess)
		}
	})

	t.Run("JSON input", func(t *testing.T) {
		samTemplate := `{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Transform": "AWS::Serverless-2016-10-31",
  "Resources": {
    "MyFunction": {
      "Type": "AWS::Serverless::Function",
      "Properties": {
        "Handler": "index.handler",
        "Runtime": "nodejs18.x",
        "CodeUri": "s3://bucket/key"
      }
    }
  }
}`
		inputFile := filepath.Join(tmpDir, "template.json")
		if err := os.WriteFile(inputFile, []byte(samTemplate), 0644); err != nil {
			t.Fatalf("failed to write input file: %v", err)
		}

		var stdout bytes.Buffer
		opts := &Options{
			TemplateFile: inputFile,
			Stdout:       true,
		}
		exitCode := runTransform(opts, &stdout, nil)
		if exitCode != ExitSuccess {
			t.Errorf("exitCode = %d, want %d", exitCode, ExitSuccess)
		}
	})
}

// TestHelpOutput tests the help message output.
func TestHelpOutput(t *testing.T) {
	cmd := newRootCmd()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	helpText := stdout.String()

	// Verify help contains essential information
	requiredStrings := []string{
		"--template-file",
		"-t",
		"--output-template",
		"-o",
		"--stdout",
		"--verbose",
		"--region",
		"--help",
		"--version",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(helpText, s) {
			t.Errorf("help text missing %q", s)
		}
	}
}

// TestVersionOutput tests the version output.
func TestVersionOutput(t *testing.T) {
	cmd := newRootCmd()
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	versionText := stdout.String()

	if !strings.Contains(versionText, "sam-translate") {
		t.Error("version text should contain 'sam-translate'")
	}
}

// TestVerboseOutput tests that verbose mode produces additional output.
func TestVerboseOutput(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "sam-translate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	samTemplate := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: s3://bucket/key
`
	inputFile := filepath.Join(tmpDir, "template.yaml")
	if err := os.WriteFile(inputFile, []byte(samTemplate), 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	opts := &Options{
		TemplateFile: inputFile,
		Stdout:       true,
		Verbose:      true,
	}
	exitCode := runTransform(opts, &stdout, &stderr)
	if exitCode != ExitSuccess {
		t.Errorf("exitCode = %d, want %d", exitCode, ExitSuccess)
	}

	// Verbose mode should produce some stderr output
	if stderr.Len() == 0 {
		t.Error("expected verbose output on stderr")
	}
}

// TestPartitionDetection tests that region flag affects partition detection.
func TestPartitionDetection(t *testing.T) {
	tests := []struct {
		region    string
		partition string
	}{
		{"us-east-1", "aws"},
		{"us-west-2", "aws"},
		{"eu-west-1", "aws"},
		{"cn-north-1", "aws-cn"},
		{"cn-northwest-1", "aws-cn"},
		{"us-gov-west-1", "aws-us-gov"},
	}

	for _, tt := range tests {
		t.Run(tt.region, func(t *testing.T) {
			partition := getPartitionForRegion(tt.region)
			if partition != tt.partition {
				t.Errorf("getPartitionForRegion(%q) = %q, want %q", tt.region, partition, tt.partition)
			}
		})
	}
}

// TestMissingTemplateFile tests that missing template file produces error.
func TestMissingTemplateFile(t *testing.T) {
	cmd := newRootCmd()
	var stderr bytes.Buffer
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	// Should fail due to missing required template file
	if err == nil {
		t.Error("expected error for missing template file")
	}
}

// TestMissingOutputDestination tests that missing output destination produces error.
func TestMissingOutputDestination(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "sam-translate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	samTemplate := `AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: nodejs18.x
      CodeUri: s3://bucket/key
`
	inputFile := filepath.Join(tmpDir, "template.yaml")
	if writeErr := os.WriteFile(inputFile, []byte(samTemplate), 0644); writeErr != nil {
		t.Fatalf("failed to write input file: %v", writeErr)
	}

	cmd := newRootCmd()
	var stderr bytes.Buffer
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"-t", inputFile})

	err = cmd.Execute()
	// Should fail due to missing output destination
	if err == nil {
		t.Error("expected error for missing output destination")
	}
}
