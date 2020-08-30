package main

import (
	"net/http"
)

func eventsHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {

	http.HandleFunc("/events", eventsHandler)
	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		panic(err)
	}
}
