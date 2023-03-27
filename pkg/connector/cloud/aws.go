package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/haandol/devops-academy-eda-demo/pkg/config"
)

var awsCfg *AWSConfig

type AWSConfig struct {
	Cfg aws.Config
}

func GetAWSConfig(appCfg *config.AWS) (*AWSConfig, error) {
	if awsCfg != nil {
		return awsCfg, nil
	}

	optFns := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(appCfg.Region),
	}
	if appCfg.UseLocal {
		optFns = append(optFns, awsconfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "0",
					URL:           "http://dynamodb:8000",
					SigningRegion: appCfg.Region,
				}, nil
			}),
		))
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), optFns...)
	if err != nil {
		return nil, err
	}
	awsCfg = &AWSConfig{cfg}

	return awsCfg, nil
}
