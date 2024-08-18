package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/login", LoginHandler).Methods("POST")
	router.HandleFunc("/api/write", AuthMiddleware(WriteHandler)).Methods("POST")
	router.HandleFunc("/api/read", AuthMiddleware(ReadHandler)).Methods("POST")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
