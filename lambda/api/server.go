package api

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func New(dynamodbClient *dynamodb.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Type", "text/plain")

		w.Header().Add("Access-Allow-Control-Origin", "*")
		w.Header().Add("Accesss-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		result, err := dynamodbClient.ListTables(r.Context(), nil)
		if err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for _, name := range result.TableNames {
			log.Println(name)
		}
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
	return mux
}
