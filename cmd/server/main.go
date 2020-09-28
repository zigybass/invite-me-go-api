package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zigybass/invite-me-go-api/pkg/events"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/events", events.Db.GetEvents).Methods("GET")
	router.HandleFunc("/events/{id}", events.Db.GetEvent).Methods("GET")

	router.HandleFunc("/events", events.Db.AddEvent).Methods("POST", "OPTIONS")

	router.HandleFunc("/events/{id}", events.Db.DeleteEvent).Methods("DELETE", "OPTIONS")

	router.Use(mux.CORSMethodMiddleware(router))

	log.Fatal(http.ListenAndServe(":8081", router))
}
