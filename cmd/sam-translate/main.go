// Package main provides the sam-translate CLI tool for transforming SAM templates to CloudFormation.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lex00/aws-sam-translator-go/pkg/region"
	"github.com/lex00/aws-sam-translator-go/pkg/translator"
	"github.com/spf13/cobra"
)

// Exit codes matching Python sam-translator.
const (
	ExitSuccess        = 0
	ExitTransformError = 1
	ExitInvalidArgs    = 2
)

// Options holds the CLI configuration.
type Options struct {
	TemplateFile   string
	OutputTemplate string
	Stdout         bool
	Verbose        bool
	Region         string
}

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(ExitInvalidArgs)
	}
}

// newRootCmd creates the root cobra command.
func newRootCmd() *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:     "sam-translate",
		Short:   "Transform SAM templates to CloudFormation",
		Long:    `sam-translate transforms AWS Serverless Application Model (SAM) templates to standard CloudFormation templates.`,
		Version: fmt.Sprintf("%s (translator: %s)", getVersion(), translator.Version),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate that template file is provided
			if opts.TemplateFile == "" {
				return fmt.Errorf("required flag \"template-file\" not set")
			}

			// Validate that either output file or stdout is specified
			if opts.OutputTemplate == "" && !opts.Stdout {
				return fmt.Errorf("either --output-template or --stdout must be specified")
			}

			// Run the transform
			exitCode := runTransform(&opts, cmd.OutOrStdout(), cmd.ErrOrStderr())
			if exitCode != ExitSuccess {
				os.Exit(exitCode)
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Define flags
	cmd.Flags().StringVarP(&opts.TemplateFile, "template-file", "t", "", "Path to SAM template file (required)")
	cmd.Flags().StringVarP(&opts.OutputTemplate, "output-template", "o", "", "Path to output CloudFormation template")
	cmd.Flags().BoolVar(&opts.Stdout, "stdout", false, "Write output to stdout")
	cmd.Flags().BoolVar(&opts.Verbose, "verbose", false, "Enable verbose logging")
	cmd.Flags().StringVar(&opts.Region, "region", "", "AWS region for partition detection (default: us-east-1)")

	// Mark template-file as required
	_ = cmd.MarkFlagRequired("template-file")

	return cmd
}

// runTransform performs the actual SAM to CloudFormation transformation.
// It returns an exit code to facilitate testing.
func runTransform(opts *Options, stdout io.Writer, stderr io.Writer) int {
	// Default writers if not provided
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	// Log verbose info
	if opts.Verbose {
		fmt.Fprintf(stderr, "Reading template from: %s\n", opts.TemplateFile)
		if opts.Region != "" {
			fmt.Fprintf(stderr, "Using region: %s\n", opts.Region)
		}
	}

	// Read the input template
	input, err := os.ReadFile(opts.TemplateFile)
	if err != nil {
		fmt.Fprintf(stderr, "Error: failed to read template file: %v\n", err)
		return ExitTransformError
	}

	if opts.Verbose {
		fmt.Fprintf(stderr, "Template size: %d bytes\n", len(input))
	}

	// Create translator with options
	translatorOpts := translator.Options{
		Region:    region.RegionOrDefault(opts.Region),
		Partition: getPartitionForRegion(opts.Region),
	}

	if opts.Verbose {
		fmt.Fprintf(stderr, "Using partition: %s\n", translatorOpts.Partition)
	}

	tr := translator.NewWithOptions(translatorOpts)

	// Perform the transformation
	if opts.Verbose {
		fmt.Fprintf(stderr, "Transforming template...\n")
	}

	output, err := tr.TransformBytes(input)
	if err != nil {
		// Format the error message
		errMsg := formatError(err)
		fmt.Fprintf(stderr, "Error: %s\n", errMsg)
		return ExitTransformError
	}

	if opts.Verbose {
		fmt.Fprintf(stderr, "Transformation successful. Output size: %d bytes\n", len(output))
	}

	// Write output
	if opts.Stdout {
		_, err = stdout.Write(output)
		if err != nil {
			fmt.Fprintf(stderr, "Error: failed to write to stdout: %v\n", err)
			return ExitTransformError
		}
		// Add newline for better terminal output
		fmt.Fprintln(stdout)
	}

	if opts.OutputTemplate != "" {
		if opts.Verbose {
			fmt.Fprintf(stderr, "Writing output to: %s\n", opts.OutputTemplate)
		}
		err = os.WriteFile(opts.OutputTemplate, output, 0644)
		if err != nil {
			fmt.Fprintf(stderr, "Error: failed to write output file: %v\n", err)
			return ExitTransformError
		}
	}

	return ExitSuccess
}

// getPartitionForRegion returns the AWS partition for the given region.
func getPartitionForRegion(regionStr string) string {
	if regionStr == "" {
		return "aws"
	}

	// Use the region package for partition detection
	partition := region.GetPartitionForRegion(regionStr)
	return string(partition)
}

// formatError formats a transformation error for user-friendly output.
func formatError(err error) string {
	if err == nil {
		return ""
	}

	// Handle TransformError with multiple errors
	if te, ok := err.(*translator.TransformError); ok {
		var msgs []string
		for _, e := range te.Errors {
			msgs = append(msgs, formatSingleError(e))
		}
		if len(msgs) == 1 {
			return msgs[0]
		}
		return "multiple errors:\n  - " + strings.Join(msgs, "\n  - ")
	}

	return formatSingleError(err)
}

// formatSingleError formats a single error with source location if available.
func formatSingleError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	// Check if the error already has line/column info
	if strings.Contains(errStr, "line") && strings.Contains(errStr, "column") {
		return errStr
	}

	return errStr
}
