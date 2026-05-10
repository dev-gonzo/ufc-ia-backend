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

type Fighter struct {
	ID        *string    `json:"id,omitempty"`
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	Record    string     `json:"record"`
	Nickname  *string    `json:"nickname,omitempty"`
	Height    *string    `json:"height,omitempty"`
	Weight    *string    `json:"weight,omitempty"`
	Reach     *string    `json:"reach,omitempty"`
	Stance    *string    `json:"stance,omitempty"`
	DOB       *string    `json:"dob,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Fight struct {
	ID            *string    `json:"id,omitempty"`
	EventID       *string    `json:"event_id,omitempty"`
	URL           string     `json:"url"`
	WeightClass   string     `json:"weight_class"`
	Method        string     `json:"method"`
	Round         int        `json:"round"`
	Time          string     `json:"time"`
	Winner        string     `json:"winner"`
	RedFighterID  *string    `json:"red_fighter_id,omitempty"`
	BlueFighterID *string    `json:"blue_fighter_id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
