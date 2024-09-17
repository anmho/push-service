package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func WithEndpoint(endpoint string) func(o *dynamodb.Options) {
	return func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	}
}

func MakeClient(ctx context.Context, region string, isLocal bool) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	var opts []func(*dynamodb.Options)
	if isLocal {
		opts = append(opts, WithEndpoint("http://localhost:8000"))
	}

	client := dynamodb.NewFromConfig(cfg, opts...)
	return client, nil
}
