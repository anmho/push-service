package api

import "net/http"



func New() http.Handler {
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
	return mux
}