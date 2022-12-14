package lambda

type AWSLambdaConfig struct {
	Region        string `env:"AWS_REGION,required"`
	FunctionName  string `env:"AWS_LAMBDA_FUNCTION_NAME,required"`
	RuntimeApi    string `env:"AWS_LAMBDA_RUNTIME_API,required"`
	IsSAMLocal    bool   `env:"AWS_SAM_LOCAL" envDefault:"false"`
	EnvironmentID string
	ExtensionName string
}
