package utils

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

const (
	DefaultRegion = "ca-central-1"
	TimeOut       = 2 * time.Minute

	LAMBDA_NAMES_ENV_VAR = "LAMBDA_NAMES"
	LAMBDA_REPOS_ENV_VAR = "LAMBDA_IMAGE_REPOS"
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

func GetLambdaMapping() map[string]string {
	lambdaStr, found := os.LookupEnv(LAMBDA_NAMES_ENV_VAR)
	if !found {
		return nil
	}
	repoStr, _ := os.LookupEnv(LAMBDA_REPOS_ENV_VAR)
	repositories := strings.Split(repoStr, ",")

	mapping := make(map[string]string)

	for i, lambda := range strings.Split(lambdaStr, ",") {
		repo := ""
		if i < len(repositories) {
			repo = repositories[i]
		}
		mapping[lambda] = repo

		if len(repo) > 0 {
			mapping[repo] = lambda
		}
	}

	return mapping
}

type Event struct {
	Detail  EventDetail `json:"detail"`
	Account string      `json:"account"`
	Region  string      `json:"region"`
}

type EventDetail struct {
	ActionType     string `json:"action-type"`
	ImageDigest    string `json:"image-digest"`
	ImageTag       string `json:"image-tag"`
	RepositoryName string `json:"repository-name"`
	Result         string `json:"result"`
}
