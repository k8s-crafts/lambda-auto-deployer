package utils

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

const (
	DefaultRegion = "ca-central-1"
	TimeOut       = 2 * time.Minute
)

func GetAWSRegion() string {
	region, found := os.LookupEnv("AWS_REGION")
	if !found {
		return DefaultRegion
	}
	return region
}

func GetAWSConfig(ctx context.Context) (*aws.Config, error) {
	region := GetAWSRegion()
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetLambdaClient(config aws.Config) *lambda.Client {
	return lambda.NewFromConfig(config)
}

type Event map[string]interface{}

type LambdaSource struct {
	// Lambda source code is packaged as a container image.
	Image Image
}

type Image struct {
	// The name of the container image, for example, `public.ecr.aws/lambda/python:3.12`.
	Name string
	// Platform of the container image, for example, `linnux/amd64` or `linnux/arm64`
	Platform Platform
}

type Platform struct {
	// Architecture field specifies the CPU architecture, for example, `amd64` or `arm64`.
	Architecture string
	// OS specifies the operating system, for example, `linux` or `windows`.
	OS string
}
