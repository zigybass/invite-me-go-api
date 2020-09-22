package events

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/zigybass/invite-me-go-api/pkg/cors"
)

type Event struct {
	mux     sync.Mutex
	Name    string `json:"name"`
	Id      string `json:"id"`
	OnGoing bool   `json:"onGoing"`
}

type eventHandlers struct {
	mux   sync.Mutex
	store map[string]Event
}

var events = []Event{
	{
		Name:    "Ultimate Frisbee",
		Id:      "id1",
		OnGoing: true,
	},
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	cors.SetupCORS(&w, r)

	jsonBytes, err := json.Marshal(events)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func (h *eventHandlers) AddEvent(w http.ResponseWriter, r *http.Request) {
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
	// h.Lock()
	h.store[event.Id] = event
	// h.Unlock()

	jsonBytes, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func (h *eventHandlers) GetEvent(w http.ResponseWriter, r *http.Request) {
	cors.SetupCORS(&w, r)
	parts := strings.Split(r.URL.String(), "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// h.Lock()
	event, ok := h.store[parts[2]]
	// h.Unlock()

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

func NewEventHandlers() *eventHandlers {
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