package main

import (
	"aws-lambda-go/api"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"log/slog"
	"net/http"
	"os"

	"flag"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/joho/godotenv"
)

const (
	port = 8080
)

func WithEndpoint(endpoint string) func(o *dynamodb.Options) {
	return func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	}
}

func MakeDynamoClient(ctx context.Context, region string, isLocal bool) (*dynamodb.Client, error) {
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

func main() {
	isLocal := flag.Bool("local", false, "-local specifies whether to run in prod mode (local=false) (AWS Lambda) or dev mode (localhost)")
	flag.Parse()

	ctx := context.Background()

	var dynamoClient *dynamodb.Client
	var err error
	if *isLocal {
		dynamoClient, err = MakeDynamoClient(ctx, "localhost", true)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		dynamoClient, err = MakeDynamoClient(ctx, "us-west-2", false)
		if err != nil {
			log.Fatalln(err)
		}
	}

	mux := api.New(dynamoClient)

	if *isLocal {
		_ = godotenv.Load(".env.development")
		slog.Info("server listening", slog.Int("port", port))
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("shutting down server")
				os.Exit(0)
			} else {
				log.Fatalln(err)
			}
		}
	} else {
		lambda.Start(httpadapter.New(mux).ProxyWithContext)
	}
}
