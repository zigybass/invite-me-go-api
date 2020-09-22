package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zigybass/invite-me-go-api/pkg/events"
)

func testThis(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET request")
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/events", events.Db.GetEvents).Methods("GET")
	router.HandleFunc("/events/{id}", events.Db.GetEvent).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", router))
}
