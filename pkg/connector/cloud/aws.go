package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

var awsCfg *AWSConfig

type AWSConfig struct {
	Cfg aws.Config
}

func GetAWSConfig() (*AWSConfig, error) {
	if awsCfg != nil {
		return awsCfg, nil
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion("ap-northeast-2"))
	if err != nil {
		return nil, err
	}
	awsCfg = &AWSConfig{cfg}

	return awsCfg, nil
}
