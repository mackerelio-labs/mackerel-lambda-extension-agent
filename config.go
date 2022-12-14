package main

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/caarlos0/env/v6"
	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/host/mackerel"
	"github.com/mackerelio-labs/mackerel-lambda-extension-agent/lambda"
)

type Config struct {
	MackerelConfig  mackerel.MackerelConfig
	AWSLambdaConfig lambda.AWSLambdaConfig
}

func GetConfig() (*Config, error) {
	conf, err := parseEnv()
	if err != nil {
		return nil, err
	}

	if conf.MackerelConfig.ApiKey != "" && conf.MackerelConfig.ApiKeySSMParamName != "" {
		return nil, errors.New("either EXT_MACKEREL_API_KEY or EXT_MACKEREL_API_KEY_SSM can be specified")
	}

	if conf.MackerelConfig.ApiKeySSMParamName != "" {
		apiKey, err := fetchMackerelApiKeyFromSSM(conf.AWSLambdaConfig.Region, conf.MackerelConfig.ApiKeySSMParamName)
		if err != nil {
			return nil, err
		}
		conf.MackerelConfig.ApiKey = apiKey
	}

	environmentID, err := getEnvironmentID()
	if err != nil {
		return nil, err
	}
	conf.AWSLambdaConfig.EnvironmentID = environmentID

	conf.AWSLambdaConfig.ExtensionName = getExtensionName()

	return conf, nil
}

func parseEnv() (*Config, error) {
	conf := &Config{}
	if err := env.Parse(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func fetchMackerelApiKeyFromSSM(region string, name string) (string, error) {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(region)
	client := ssm.New(sess)
	res, err := client.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return aws.StringValue(res.Parameter.Value), nil
}

func getEnvironmentID() (string, error) {
	bytes, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
	if err != nil {
		return "", err
	}
	environmentID := strings.TrimRight(string(bytes), "\r\n")
	return environmentID, nil
}

func getExtensionName() string {
	return path.Base(os.Args[0])
}
