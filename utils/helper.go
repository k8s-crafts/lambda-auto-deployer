// Lambda Auto Deployer
// Copyright (C) 2024 Thuan Vo
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

	LAMBDA_NAMES_ENV_VAR      = "LAMBDA_NAMES"
	LAMBDA_REPOS_ENV_VAR      = "LAMBDA_IMAGE_REPOS"
	LAMBDA_IMAGE_TAGS_ENV_VAR = "LAMBDA_IMAGE_TAGS"
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

type ImageTagFilter func(tag string) bool

func GetImageTagFilter() ImageTagFilter {
	tags, found := os.LookupEnv(LAMBDA_IMAGE_TAGS_ENV_VAR)
	return func(tag string) bool {
		if !found {
			return true
		}
		return strings.Contains(tags, tag)
	}
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
