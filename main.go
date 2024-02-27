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

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	lambdasdk "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/tthvo/lambda-auto-deployer/utils"
)

func HandleRequest(ctx context.Context, event *utils.Event) error {
	if event == nil {
		return fmt.Errorf("event must not be nil")
	}

	log.Printf("received event: %+v\n", event)

	// Get a mapping from lambda name to ECR repository, vice versa
	mapping := utils.GetLambdaMapping()
	if mapping == nil {
		return fmt.Errorf("lambda mapping is not available")
	}

	// Find AWS configurations
	config, err := utils.GetAWSConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load SDK config: %s", err.Error())
	}

	// Create a Lambda client to roll out new lambda version
	client := utils.GetLambdaClient(*config)

	lambda := mapping[event.Detail.RepositoryName]
	if len(lambda) == 0 {
		log.Printf("no lambda to roll out for repository: %s. Skipped", event.Detail.RepositoryName)
		return nil
	}

	imageUri := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", event.Account, event.Region, event.Detail.RepositoryName, event.Detail.ImageTag)

	updateOpts := &lambdasdk.UpdateFunctionCodeInput{
		FunctionName:  &lambda,
		Architectures: []types.Architecture{types.ArchitectureX8664},
		ImageUri:      &imageUri,
	}
	out, err := client.UpdateFunctionCode(ctx, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update function code: %s", err.Error())
	}

	log.Printf("sucessfully roll out %s(%s)", *out.FunctionName, *out.FunctionArn)
	return nil

}

func main() {
	lambda.Start(HandleRequest)
}
