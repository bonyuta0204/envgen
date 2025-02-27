# envgen - AWS SSM Parameter Store to .env File Generator

`envgen` is a command-line tool that generates `.env` files by retrieving values from AWS Systems Manager Parameter Store (SSM) based on a template file. This allows you to keep sensitive information like API keys and database credentials secure while still maintaining version control over your environment configuration.

## Features

- Automatically generates `.env` files from `.env.template`
- Retrieves sensitive values from AWS SSM Parameter Store
- Supports version control for non-sensitive environment variables
- Simple placeholder syntax: `{{SSM:/path/to/param}}`
- Secure handling of sensitive data
- Easy integration with developer workflows

## Installation

### From Source

```bash
go install github.com/bonyuta0204/envgen@latest
```

### Binary Releases

Download the latest binary from the [Releases](https://github.com/bonyuta0204/envgen/releases) page.

## Usage

```bash
envgen [global options] command [command options]
```

### Global Options

- `--template, -t`: Path to the template file (default: `.env.template`)
- `--output, -o`: Path to the output file (default: `.env`)
- `--region, -r`: AWS region (overrides AWS_REGION environment variable)
- `--verbose, -V`: Enable verbose output
- `--help, -h`: Show help
- `--version, -v`: Print the version

### Commands

#### Generate (default)

Generates a `.env` file by retrieving values from AWS SSM Parameter Store.

```bash
# Generate .env file using default template (.env.template)
envgen

# Specify custom template and output files
envgen -t config/.env.template -o .env.production

# Specify AWS region
envgen -r us-west-2
```

#### Validate

Validates the template file without generating the `.env` file.

```bash
# Validate the default template
envgen validate

# Validate a specific template
envgen validate -t config/.env.template
```

#### List

Lists all SSM parameters referenced in the template file.

```bash
# List parameters in the default template
envgen list

# List parameters in a specific template
envgen list -t config/.env.template
```

## Template Format

The template file follows the standard `.env` format but supports placeholders for AWS SSM parameters:

```ini
# Regular environment variables
APP_ENV=development
LOG_LEVEL=debug

# Variables with SSM placeholders
DB_HOST={{SSM:/myapp/db_host}}
DB_USER={{SSM:/myapp/db_user}}
DB_PASSWORD={{SSM:/myapp/db_password}}
```

## AWS Configuration

### Required Permissions

The tool requires the following AWS IAM permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter",
        "ssm:GetParametersByPath"
      ],
      "Resource": "arn:aws:ssm:*:*:parameter/*"
    }
  ]
}
```

### AWS Credentials

The tool uses the AWS SDK's default credential provider chain, which looks for credentials in the following order:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM role for Amazon EC2 or ECS task role
4. AWS SSO

## Best Practices

1. **Never commit `.env` files to your repository**
   - Add `.env` to your `.gitignore` file

2. **Use a consistent naming convention for SSM parameters**
   - Example: `/app_name/environment/parameter_name`

3. **Use AWS KMS to encrypt sensitive parameters**
   - The tool automatically handles decryption

4. **Keep non-sensitive defaults in the template**
   - Only use SSM placeholders for sensitive values

## License

MIT
