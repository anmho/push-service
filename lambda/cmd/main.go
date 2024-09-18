package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"log"
	"log/slog"
	"net/http"
	"os"
	"send-to/api"
	"send-to/dynamo"
	"send-to/expo"
	"send-to/push"

	"flag"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/joho/godotenv"
)

const (
	port = 8080
)

func main() {
	isLocal := flag.Bool("local", false, "-local specifies whether to run in prod mode (local=false) (AWS Lambda) or dev mode (localhost)")
	flag.Parse()

	ctx := context.Background()

	var dynamoClient *dynamodb.Client
	var err error
	//if *isLocal {
	//	dynamoClient, err = dynamo.MakeService(ctx, "localhost", true)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//} else {
	var region = "us-west-2"
	dynamoClient, err = dynamo.MakeClient(ctx, region, false)
	if err != nil {
		log.Fatalln(err)
	}
	//}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	schedulerClient := scheduler.NewFromConfig(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	expoClient := expo.MakeClient()
	pushService := push.MakeService(dynamoClient, schedulerClient, expoClient)
	mux := api.New(pushService)

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
