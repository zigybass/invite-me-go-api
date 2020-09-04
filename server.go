package main

import (
	"encoding/json"
	"net/http"
)

type Event struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	OnGoing bool   `json:"onGoing"`
}

type eventHandlers struct {
	store map[string]Event
}

func (h *eventHandlers) events(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

}

func (h *eventHandlers) get(w http.ResponseWriter, r *http.Request) {
	events := make([]Event, len(h.store))

	i := 0
	for _, event := range h.store {
		events[i] = event
		i++
	}

	jsonBytes, err := json.Marshal(events)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func (h *eventHandlers) post(w http.ResponseWriter, r *http.Request) {
}

func newEventHandlers() *eventHandlers {
	// Another option is to grab data from DB storage here
	return &eventHandlers{
		store: map[string]Event{
			"id1": Event{
				Name:    "Ultimate Frisbee",
				ID:      "id1",
				OnGoing: true,
			},
			"id2": Event{
				Name:    "Soccer",
				ID:      "id2",
				OnGoing: false,
			},
		},
	}
}

func main() {
	eventHandlers := newEventHandlers()
	http.HandleFunc("/events", eventHandlers.events)
	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		panic(err)
	}
}
