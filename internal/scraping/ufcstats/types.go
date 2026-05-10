package ufcstats

import "time"

type Event struct {
	ID        *string    `json:"id,omitempty"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	Date      string     `json:"date"`
	Location  string     `json:"location"`
	EventSync *bool      `json:"event_sync,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
