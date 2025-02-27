package ssm

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Client represents an AWS SSM client with caching capabilities
type Client struct {
	ssmClient *ssm.Client
	cache     map[string]string
	mutex     sync.RWMutex
}

// NewClient creates a new SSM client
func NewClient(region string) (*Client, error) {
	cfg, err := loadConfig(region)
	if err != nil {
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)
	return &Client{
		ssmClient: client,
		cache:     make(map[string]string),
	}, nil
}

// loadConfig loads the AWS configuration
func loadConfig(region string) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}
	return config.LoadDefaultConfig(context.TODO(), opts...)
}

// GetParameter retrieves a parameter from SSM Parameter Store
func (c *Client) GetParameter(name string) (string, error) {
	// Check cache first
	c.mutex.RLock()
	if value, ok := c.cache[name]; ok {
		c.mutex.RUnlock()
		return value, nil
	}
	c.mutex.RUnlock()

	// Parameter not in cache, fetch from AWS
	input := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	}

	result, err := c.ssmClient.GetParameter(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get parameter %s: %w", name, err)
	}

	if result.Parameter.Value == nil {
		return "", fmt.Errorf("parameter %s has nil value", name)
	}

	value := *result.Parameter.Value

	// Cache the result
	c.mutex.Lock()
	c.cache[name] = value
	c.mutex.Unlock()

	return value, nil
}

// GetParametersByPath retrieves all parameters under a path prefix
func (c *Client) GetParametersByPath(path string) (map[string]string, error) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	result := make(map[string]string)
	var nextToken *string

	for {
		input := &ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
			NextToken:      nextToken,
		}

		resp, err := c.ssmClient.GetParametersByPath(context.TODO(), input)
		if err != nil {
			return nil, fmt.Errorf("failed to get parameters by path %s: %w", path, err)
		}

		for _, param := range resp.Parameters {
			paramName := *param.Name
			paramValue := *param.Value

			// Cache the result
			c.mutex.Lock()
			c.cache[paramName] = paramValue
			c.mutex.Unlock()

			// Store in result map
			result[paramName] = paramValue
		}

		nextToken = resp.NextToken
		if nextToken == nil {
			break
		}
	}

	return result, nil
}
