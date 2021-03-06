package events

import (
	"sync"
)

type eventModel struct {
	Name        string   `json:"name"`
	Id          string   `json:"id"`
	Description string   `json:"description:"`
	CreatedBy   string   `json:"createdBy"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
	Location    location `json:"location"`
	Capacity    int      `json:"capacity"`
	Invited     []person `json:"invited"`
	Attending   []person `json:"attending"`
	Declined    []person `json:"declined"`
	SoftDeleted bool     `json:"softDeleted"`
}

type location struct {
	Name       string `json:"name"`
	AddressOne string `json:"addressOne"`
	AddressTwo string `json:"addressTwo"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zip        string `json:"zip"`
	Lat        int    `json:"lat"`
	Lon        int    `json:"lon"`
}

type person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"LastName"`
}

type eventHandlers struct {
	sync.Mutex
	store map[string]eventModel
}

func newEventHandlers() *eventHandlers {
	// Another option is to grab data from DB storage here
	return &eventHandlers{
		store: map[string]eventModel{
			"id1": eventModel{
				Name: "Ultimate Frisbee",
				Id:   "id1",
			},
		},
	}
}
