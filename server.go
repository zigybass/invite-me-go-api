package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Event struct {
	Name    string `json:"name"`
	ID      string `json:"ID"`
	OnGoing bool   `json:"onGoing"`
}

type eventHandlers struct {
	sync.Mutex
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

	h.Lock()
	i := 0
	for _, event := range h.store {
		events[i] = event
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(events)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func (h *eventHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json' but got '%s'", contentType)))
		return
	}

	var event Event
	err = json.Unmarshal(bodyBytes, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	event.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[event.ID] = event
	defer h.Unlock()
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
