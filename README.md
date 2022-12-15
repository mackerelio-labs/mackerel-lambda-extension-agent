# mackerel-lambda-extension-agent

!! This is an **experimental** project. Please be careful when using in production. !!

mackerel-lambda-extension is a extension layer of AWS Lambda to monitor each Lambda runtime environment with Mackerel.

Each Lambda runtime environment is individually registered as a host on Mackerel, and the metrics derived from [Lambda Telemetry API](https://docs.aws.amazon.com/lambda/latest/dg/telemetry-api.html) are posted.

## Usage

### Registration

You can download the Lamdba layer zip file on [Release](https://github.com/mackerelio-labs/mackerel-lambda-extension-agent/releases) page.

Add mackerel-lambda-extension-agent.zip file as a layer to the Lambda function you wish to monitor.

If you want to build it yourself, use the following command:

```sh
cd /path/to/mackerel-lambda-extension-agent
make build-extension
```

### Settings

Specify the environment variables for the target Lambda function.

Either `EXT_MACKEREL_API_KEY` or `EXT_MACKEREL_API_KEY_SSM` must be specified.

| Name | Description |
| :-- | :-- |
| `EXT_MACKEREL_API_KEY` | Mackerel API key. Read and write permission is required |
| `EXT_MACKEREL_API_KEY_SSM` | Name of SSM parameter store where Mackerel API key is stored with encryption |
| `EXT_MACKEREL_ROLE_FULL_NAMES` | Service and role to which hosts belong. The format is `<service>:<role>,...,<service>:<role>`.  |
| `EXT_LOG_LEVEL` | Select a log level from the following options: `DEBUG`, `INFO`, `WARNING`, `ERROR`. Default is `WARNING` |

### Example: Configuration by Terraform

```hcl
locals {
   extension_file = "./mackerel-lambda-extension-agent.zip"
}

resource "aws_lambda_layer_version" "agent" {
  filename         = local.extension_file
  layer_name       = "mackerel-lambda-extension-agent"
  source_code_hash = filebase64sha256(local.extension_file)
}

resource "aws_lambda_function" "this" {
  function_name = "sample-function"

  ...

  layers = [aws_lambda_layer_version.agent.arn]

  environment {
    variables = {
      EXT_MACKEREL_API_KEY         = var.mackerel_api_key
      EXT_MACKEREL_ROLE_FULL_NAMES = "lambda:sample-function"
      EXT_LOG_LEVEL                = "DEBUG"
    }
  }
}
```

## License

The source code is licensed MIT.
