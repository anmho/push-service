package main

import (
	"aws-lambda-go/api"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"flag"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)
const (
	port = 8080
)


func main() {
	key := os.Getenv("STRIPE_KEY")
	log.Println("stripe key", key)
	isLocal := flag.Bool("local", false, "-local specifies whether to run in prod mode (local=false) (AWS Lambda) or dev mode (localhost)")
	flag.Parse()
	mux := api.New()
	
	if *isLocal {
		slog.Info("server listening", slog.Int("port", port))
		http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	} else {
		lambda.Start(httpadapter.New(mux).ProxyWithContext)
	}
}
