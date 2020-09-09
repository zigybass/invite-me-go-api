package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zigybass/invite-me-go-api/cors"
)

type Event struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	OnGoing bool   `json:"onGoing"`
}

type eventHandlers struct {
	sync.Mutex
	store map[string]Event
}

func (h *eventHandlers) events(w http.ResponseWriter, r *http.Request) {
	cors.SetupCORS(&w, r)
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
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
		return
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

	event.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[event.Id] = event
	h.Unlock()

	jsonBytes, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func (h *eventHandlers) getEvent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	event, ok := h.store[parts[2]]
	h.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(event)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func newEventHandlers() *eventHandlers {
	// Another option is to grab data from DB storage here
	return &eventHandlers{
		store: map[string]Event{
			"id1": Event{
				Name:    "Ultimate Frisbee",
				Id:      "id1",
				OnGoing: true,
			},
		},
	}
}

func main() {
	eventHandlers := newEventHandlers()
	http.HandleFunc("/events", eventHandlers.events)
	http.HandleFunc("/events/", eventHandlers.getEvent)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
