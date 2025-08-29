package provider

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type AwsConfig struct {
	ExternalLabels map[string]interface{}
	cfg            aws.Config
}

func NewAWSCredentialCfg(region, ak, sk string, labels map[string]interface{}) (AwsConfig, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		func(options *config.LoadOptions) error {
			options.Region = region
			options.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     ak,
					SecretAccessKey: sk,
				}, nil
			})

			return nil
		},
	)
	if err != nil {
		return AwsConfig{}, err
	}

	return AwsConfig{
		ExternalLabels: labels,
		cfg:            cfg,
	}, nil
}

func (a AwsConfig) CloudWatchCli() *cloudwatch.Client {
	return cloudwatch.NewFromConfig(a.cfg)
}

func (a AwsConfig) RdsCli() *rds.Client {
	return rds.NewFromConfig(a.cfg)
}

func (a AwsConfig) GetExternalLabels() map[string]interface{} {
	return a.ExternalLabels
}
