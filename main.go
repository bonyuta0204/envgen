package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bonyuta0204/envgen/ssm"
	"github.com/bonyuta0204/envgen/template"
)

const (
	defaultTemplateFile = ".env.template"
	defaultOutputFile   = ".env"
)

func main() {
	// Parse command line flags
	templateFile := flag.String("template", defaultTemplateFile, "Path to the template file")
	outputFile := flag.String("output", defaultOutputFile, "Path to the output file")
	region := flag.String("region", "", "AWS region (overrides AWS_REGION environment variable)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	// Validate template file exists
	if _, err := os.Stat(*templateFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Template file %s does not exist\n", *templateFile)
		os.Exit(1)
	}

	// Create SSM client
	ssmClient, err := ssm.NewClient(*region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing AWS SSM client: %v\n", err)
		os.Exit(1)
	}

	// Read template file
	templateContent, err := os.ReadFile(*templateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading template file: %v\n", err)
		os.Exit(1)
	}

	// Process template
	processor := template.NewProcessor(ssmClient)
	result, err := processor.Process(string(templateContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing template: %v\n", err)
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(*outputFile)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Write output file
	if err := os.WriteFile(*outputFile, []byte(result), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Successfully generated %s from %s\n", *outputFile, *templateFile)
	}
}
