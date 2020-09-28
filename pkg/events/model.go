package events

import "time"

type EventModel struct {
	Name        string    `json:"name"`
	Id          string    `json:"id"`
	Description string    `json:"description:"`
	CreatedBy   string    `json:"createdBy"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Location    location  `json:"location"`
	Capacity    int       `json:"capacity"`
	Invited     []person  `json:"invited"`
	Attending   []person  `json:"attending"`
	Declined    []person  `json:"declined"`
	OnGoing     bool      `json:"onGoing"`
	SoftDeleted bool      `json:"softDeleted"`
}

type location struct {
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
