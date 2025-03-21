package types

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CalendarEventSummaryResponse struct {
	CalendarId   string      `json:"calendarId"`
	EventId      string      `json:"eventId"`
	EventType    string      `json:"eventType"`
	StartTime    string      `json:"startTime"`
	EndTime      string      `json:"endTime"`
	Notes        pgtype.Text `json:"notes"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	LocationName string      `json:"locationName"`
}

type NewContent struct {
	FragmentType         string `json:"fragmentType"`
	Name                 string `json:"name"`
	Body                 JSONB  `json:"body"`
	ExtendedAttributes   JSONB  `json:"extendedAttributes"`
	PublishedStatus      string `json:"publishedStatus"`
	PublishStartDateTime string `json:"publishStartDateTime"`
	PublishEndDateTime   string `json:"publishEndDateTime"`
}

type UpdateContent struct {
	Id                   string `json:"id"`
	FragmentType         string `json:"fragmentType"`
	Name                 string `json:"name"`
	Body                 JSONB  `json:"body"`
	ExtendedAttributes   JSONB  `json:"extendedAttributes"`
	PublishedStatus      string `json:"publishedStatus"`
	PublishStartDateTime string `json:"publishStartDateTime"`
	PublishEndDateTime   string `json:"publishEndDateTime"`
}

type PublishedStatus struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

// TODO: We may actually care about the shape of these columns
// if so, we need to create structs to represent the fields
type JSONB map[string]interface{}
