// Package main provides the sam-translate CLI tool.
package main

import (
	"fmt"
	"os"

	"github.com/lex00/aws-sam-translator-go/pkg/translator"
)

// Version is set at build time.
var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("sam-translate version %s\n", Version)
			return
		case "--help", "-h":
			printHelp()
			return
		}
	}

	printHelp()
}

func printHelp() {
	fmt.Println(`sam-translate - Transform SAM templates to CloudFormation

Usage:
  sam-translate [flags]

Flags:
  --template-file, -t    Path to SAM template file
  --output-template, -o  Path to output CloudFormation template
  --stdout               Write output to stdout
  --verbose              Enable verbose logging
  --version, -v          Print version information
  --help, -h             Show this help message

Example:
  sam-translate --template-file template.yaml --output-template output.yaml`)
	fmt.Println()
	fmt.Printf("Translator version: %s\n", translator.Version)
}
