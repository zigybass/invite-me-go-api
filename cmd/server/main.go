package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	// r.HandleFunc("/events", testing)
	// r.HandleFunc("/events/", testing)

	log.Fatal(http.ListenAndServe(":8081", r))
}
