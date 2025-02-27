# Using envgen with AWS SSM Parameter Store

This document provides a step-by-step guide on how to use `envgen` to generate `.env` files from AWS SSM Parameter Store.

## Prerequisites

1. AWS CLI installed and configured with appropriate credentials
2. AWS SSM Parameter Store set up with your parameters
3. `envgen` installed or built from source

## Setting Up AWS SSM Parameters

Before using `envgen`, you need to store your sensitive parameters in AWS SSM Parameter Store. Here's how to do it:

```bash
# Store a simple string parameter
aws ssm put-parameter --name "/myapp/db_host" --value "db.example.com" --type String

# Store a secure string parameter (encrypted)
aws ssm put-parameter --name "/myapp/db_password" --value "your-secure-password" --type SecureString

# Update an existing parameter
aws ssm put-parameter --name "/myapp/db_host" --value "new-db.example.com" --type String --overwrite
```

## Creating a Template File

Create a `.env.template` file in your project root with placeholders for SSM parameters:

```ini
# Application configuration
APP_NAME=myapp
APP_ENV=development
LOG_LEVEL=debug

# Database configuration
DB_HOST={{SSM:/myapp/db_host}}
DB_PORT=5432
DB_NAME={{SSM:/myapp/db_name}}
DB_USER={{SSM:/myapp/db_user}}
DB_PASSWORD={{SSM:/myapp/db_password}}
```

## Generating the .env File

Run `envgen` to generate your `.env` file:

```bash
# Basic usage (uses default .env.template and outputs to .env)
./envgen

# Specify a custom template and output file
./envgen -t ./config/.env.template -o ./config/.env

# Specify AWS region
./envgen -r us-west-2

# Enable verbose output
./envgen -V
```

## Additional Commands

### Validating Templates

You can validate your template without generating a `.env` file:

```bash
# Validate the default template
./envgen validate

# Validate a specific template
./envgen validate -t ./config/.env.template
```

### Listing SSM Parameters

You can list all SSM parameters referenced in your template:

```bash
# List parameters in the default template
./envgen list

# List parameters in a specific template
./envgen list -t ./config/.env.template
```

This is useful for auditing which parameters your application depends on.

## Verifying the Generated .env File

After running `envgen`, your `.env` file should contain the actual values from AWS SSM Parameter Store:

```ini
# Application configuration
APP_NAME=myapp
APP_ENV=development
LOG_LEVEL=debug

# Database configuration
DB_HOST=db.example.com
DB_PORT=5432
DB_NAME=myapp_db
DB_USER=dbuser
DB_PASSWORD=your-secure-password
```

## Using with CI/CD Pipelines

You can integrate `envgen` into your CI/CD pipelines to generate environment files during deployment:

```yaml
# Example GitHub Actions workflow
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v1
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2
      
      - name: Download envgen
        run: |
          curl -L -o envgen https://github.com/bonyuta0204/envgen/releases/download/v1.0.0/envgen-linux-amd64
          chmod +x envgen
      
      - name: Generate .env file
        run: ./envgen -V
      
      # Continue with your deployment steps
```

## Troubleshooting

### AWS Credentials Issues

If you encounter AWS credentials issues:

1. Ensure your AWS credentials are properly configured
2. Check that your IAM user/role has the necessary permissions
3. Try setting the AWS region explicitly with the `-r` flag

### Parameter Not Found

If a parameter is not found:

1. Verify that the parameter exists in AWS SSM Parameter Store
2. Check the parameter name and path in your template
3. Ensure you have permission to access the parameter

### Other Issues

For other issues, run with verbose output enabled:

```bash
./envgen -V
```

This will provide more information about what's happening during execution.
