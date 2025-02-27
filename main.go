package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bonyuta0204/envgen/ssm"
	"github.com/bonyuta0204/envgen/template"
	"github.com/urfave/cli/v2"
)

const (
	defaultTemplateFile = ".env.template"
	defaultOutputFile   = ".env"
)

func main() {
	app := &cli.App{
		Name:        "envgen",
		Usage:       "Generate .env files from AWS SSM Parameter Store",
		Description: "A tool that generates .env files by retrieving values from AWS SSM Parameter Store based on a template file",
		Version:     "1.0.0",
		Compiled:    time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Yuta Nakamura",
				Email: "nakamurayuta0204@github.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Value:   defaultTemplateFile,
				Usage:   "Path to the template file",
				EnvVars: []string{"ENVGEN_TEMPLATE"},
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   defaultOutputFile,
				Usage:   "Path to the output file",
				EnvVars: []string{"ENVGEN_OUTPUT"},
			},
			&cli.StringFlag{
				Name:    "region",
				Aliases: []string{"r"},
				Usage:   "AWS region (overrides AWS_REGION environment variable)",
				EnvVars: []string{"AWS_REGION"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "Enable verbose output",
				EnvVars: []string{"ENVGEN_VERBOSE"},
			},
		},
		Action: func(c *cli.Context) error {
			return generateEnv(c)
		},
		Commands: []*cli.Command{
			{
				Name:  "validate",
				Usage: "Validate the template file without generating the .env file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "template",
						Aliases: []string{"t"},
						Value:   defaultTemplateFile,
						Usage:   "Path to the template file",
					},
					&cli.StringFlag{
						Name:    "region",
						Aliases: []string{"r"},
						Usage:   "AWS region (overrides AWS_REGION environment variable)",
					},
				},
				Action: func(c *cli.Context) error {
					return validateTemplate(c)
				},
			},
			{
				Name:  "list",
				Usage: "List all SSM parameters referenced in the template file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "template",
						Aliases: []string{"t"},
						Value:   defaultTemplateFile,
						Usage:   "Path to the template file",
					},
				},
				Action: func(c *cli.Context) error {
					return listParameters(c)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateEnv(c *cli.Context) error {
	templateFile := c.String("template")
	outputFile := c.String("output")
	region := c.String("region")
	verbose := c.Bool("verbose")

	// Validate template file exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file %s does not exist", templateFile)
	}

	// Create SSM client
	ssmClient, err := ssm.NewClient(region)
	if err != nil {
		return fmt.Errorf("error initializing AWS SSM client: %w", err)
	}

	// Read template file
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Process template
	processor := template.NewProcessor(ssmClient)
	result, err := processor.Process(string(templateContent))
	if err != nil {
		return fmt.Errorf("error processing template: %w", err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}
	}

	// Write output file
	if err := os.WriteFile(outputFile, []byte(result), 0600); err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	if verbose {
		fmt.Printf("Successfully generated %s from %s\n", outputFile, templateFile)
	}

	return nil
}

func validateTemplate(c *cli.Context) error {
	templateFile := c.String("template")
	region := c.String("region")

	// Validate template file exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file %s does not exist", templateFile)
	}

	// Create SSM client
	ssmClient, err := ssm.NewClient(region)
	if err != nil {
		return fmt.Errorf("error initializing AWS SSM client: %w", err)
	}

	// Read template file
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Process template
	processor := template.NewProcessor(ssmClient)
	_, err = processor.Process(string(templateContent))
	if err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	fmt.Printf("Template %s is valid\n", templateFile)
	return nil
}

func listParameters(c *cli.Context) error {
	templateFile := c.String("template")

	// Validate template file exists
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("template file %s does not exist", templateFile)
	}

	// Read template file
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	// Extract parameters from template
	parameters, err := template.ExtractParameters(string(templateContent))
	if err != nil {
		return fmt.Errorf("error extracting parameters: %w", err)
	}

	if len(parameters) == 0 {
		fmt.Println("No SSM parameters found in template")
		return nil
	}

	fmt.Println("SSM parameters referenced in template:")
	for _, param := range parameters {
		fmt.Printf("  - %s\n", param)
	}

	return nil
}
