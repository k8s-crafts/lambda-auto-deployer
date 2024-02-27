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
