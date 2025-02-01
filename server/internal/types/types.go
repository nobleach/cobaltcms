package types

import "github.com/jackc/pgx/v5/pgtype"

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

type PublishedStatus struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
