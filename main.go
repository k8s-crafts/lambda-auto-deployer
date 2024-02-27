package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tthvo/lambda-auto-deployer/utils"
)

func HandleRequest(ctx context.Context, event utils.Event) error {
	if event == nil {
		return fmt.Errorf("event must not be nil")
	}

	log.Printf("received event: %+v\n", event)

	// Find AWS configurations
	config, err := utils.GetAWSConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load SDK config: %s", err.Error())
	}

	// Create a Lambda client to roll out new lambda version
	client := utils.GetLambdaClient(*config)
	log.Printf("using lamdba client: %+v\n", client)

	// log.Println("using context: ")

	// updateOpts := &lambdasdk.UpdateFunctionCodeInput{
	// 	FunctionName:  &[]string{""}[0],
	// 	Architectures: []types.Architecture{},
	// 	ImageUri:      &[]string{""}[0],
	// }
	// out, err := client.UpdateFunctionCode(ctx, updateOpts)
	// if err != nil {
	// 	return fmt.Errorf("failed to update function code: %s", err.Error())
	// }

	// log.Printf("sucessfully roll out %s(%s)", *out.FunctionName, *out.FunctionArn)
	return nil

}

func main() {
	lambda.Start(HandleRequest)
}
