package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type App struct {
	id string
}

func newApp(id string) *App {
	return &App{
		id: id,
	}
}

func (app *App) Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	responseBody := map[string]string{
		"message": "Hello World",
	}

	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "internal server error"}`,
		}, nil
	}
	response := events.APIGatewayProxyResponse{
		Body:       string(responseJSON),
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                     "text/plain",
			"Access-Allow-Control-Origin":      "*",
			"Accesss-Control-Allow-Headers":    "Content-Type",
			"Access-Control-Allow-Methods":     "OPTIONS, POST, GET, PUT, DELETE",
			"Access-Control-Allow-Credentials": "true",
		},
	}
	return response, nil
}

func main() {

	// id := "someRandomString"
	// app := newApp(id)

	// lambda.Start(app.Handler)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Type", "text/plain")

		w.Header().Add("Access-Allow-Control-Origin", "*")
		w.Header().Add("Accesss-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		w.Write([]byte(`{"message": "root"}`))

	})
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Type", "text/plain")

		w.Header().Add("Access-Allow-Control-Origin", "*")
		w.Header().Add("Accesss-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		w.Write([]byte(`{"message": "Hello World"}`))

	})

	lambda.Start(httpadapter.New(mux).ProxyWithContext)
}
