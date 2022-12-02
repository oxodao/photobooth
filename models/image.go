package models

type Image struct {
	Id         int64     `json:"id" db:"id"`
	Date       Timestamp `json:"date" db:"created_at"`
	Unattended bool      `json:"unattended" db:"unattended"`
	EventId    int64     `json:"-" db:"event_id"`
}
