package main

import (
	"log"
	"net/http"

	"serverless-runtime/internal/api"
)

func main() {

	mux := http.NewServeMux() // register routes
	mux.HandleFunc("/invoke/", api.InvokeHandler)
	mux.HandleFunc("/deploy", api.DeployHandler)
	mux.HandleFunc("/functions", api.ListFunctionsHandler)
	mux.HandleFunc("/metrics", api.MetricsHandler)
	log.Println("API server running on :8080")

	err := http.ListenAndServe(":8080", mux) // start http server

	if err != nil {
		log.Fatal(err)
	}

}
