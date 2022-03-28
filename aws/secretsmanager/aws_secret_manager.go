package awssecretmanager

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
	"os"
)

var Client *secretsmanager.Client

func init() {
	region := os.Getenv("AWS_REGION")

	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	awsv2.AWSV2Instrumentor(&awsConfig.APIOptions)

	Client = secretsmanager.NewFromConfig(awsConfig)
}

func GetSecret(ctx context.Context, secretId string) (*secretsmanager.GetSecretValueOutput, error) {
	return Client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretId,
	})
}
