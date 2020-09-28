package events

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/zigybass/invite-me-go-api/pkg/cors"
)

var Db = newEventHandlers()

func (h *eventHandlers) GetEvents(w http.ResponseWriter, r *http.Request) {
	cors.SetupCORS(&w, r)

	var events []eventModel

	h.Lock()
	i := 0
	for _, event := range h.store {
		if event.SoftDeleted != true {
			events = append(events, event)
		}
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

func (h *eventHandlers) AddEvent(w http.ResponseWriter, r *http.Request) {
	cors.SetupCORS(&w, r)
	if r.Method == http.MethodOptions {
		return
	}

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

	var event eventModel
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

func (h *eventHandlers) GetEvent(w http.ResponseWriter, r *http.Request) {
	cors.SetupCORS(&w, r)
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

func (h *eventHandlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	cors.SetupCORS(&w, r)

	parts := strings.Split(r.URL.String(), "/")
	i := parts[2]

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	event, ok := h.store[i]
	event.SoftDeleted = true
	h.store[event.Id] = event
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
