package template

import (
	"fmt"
	"regexp"
	"strings"
)

// SSMParameterInterface defines the interface for retrieving SSM parameters
type SSMParameterInterface interface {
	GetParameter(name string) (string, error)
}

// Processor handles template processing
type Processor struct {
	ssmClient SSMParameterInterface
}

// NewProcessor creates a new template processor
func NewProcessor(ssmClient SSMParameterInterface) *Processor {
	return &Processor{
		ssmClient: ssmClient,
	}
}

// Process processes the template and replaces SSM placeholders with actual values
func (p *Processor) Process(templateContent string) (string, error) {
	// Regular expression to match SSM parameter placeholders: {{SSM:/path/to/param}}
	ssmRegex := regexp.MustCompile(`{{SSM:([^}]+)}}`)

	// Find all SSM parameter placeholders
	lines := strings.Split(templateContent, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		// Skip empty lines and comments
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			result = append(result, line)
			continue
		}

		// Process line with SSM placeholders
		processedLine, err := p.processLine(line, ssmRegex)
		if err != nil {
			return "", err
		}
		result = append(result, processedLine)
	}

	return strings.Join(result, "\n"), nil
}

// processLine processes a single line of the template
func (p *Processor) processLine(line string, ssmRegex *regexp.Regexp) (string, error) {
	// Find all SSM parameter placeholders in the line
	matches := ssmRegex.FindAllStringSubmatchIndex(line, -1)
	if len(matches) == 0 {
		return line, nil
	}

	// Process each placeholder
	var lastIndex int
	var result strings.Builder

	for _, match := range matches {
		// Append text before the placeholder
		result.WriteString(line[lastIndex:match[0]])

		// Extract parameter name
		paramName := line[match[2]:match[3]]

		// Get parameter value from SSM
		paramValue, err := p.ssmClient.GetParameter(paramName)
		if err != nil {
			return "", fmt.Errorf("failed to get SSM parameter %s: %w", paramName, err)
		}

		// Append parameter value
		result.WriteString(paramValue)

		// Update last index
		lastIndex = match[1]
	}

	// Append remaining text
	result.WriteString(line[lastIndex:])

	return result.String(), nil
}

// ExtractParameters extracts all SSM parameter names from a template
func ExtractParameters(templateContent string) ([]string, error) {
	// Regular expression to match SSM parameter placeholders: {{SSM:/path/to/param}}
	ssmRegex := regexp.MustCompile(`{{SSM:([^}]+)}}`)
	
	// Find all SSM parameter placeholders
	matches := ssmRegex.FindAllStringSubmatch(templateContent, -1)
	
	// Extract parameter names
	paramSet := make(map[string]struct{})
	for _, match := range matches {
		if len(match) >= 2 {
			paramSet[match[1]] = struct{}{}
		}
	}
	
	// Convert to slice
	params := make([]string, 0, len(paramSet))
	for param := range paramSet {
		params = append(params, param)
	}
	
	return params, nil
}
