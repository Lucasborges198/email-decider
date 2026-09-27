package main

import (
	"log"
	"net/http"
	handler "email-decider/handler"
)

func main() {
	http.HandleFunc("/decision", handler.FetchLayaService)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}