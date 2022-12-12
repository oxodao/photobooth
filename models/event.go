package models

type Event struct {
	Id       int64      `json:"id" db:"id"`
	Name     string     `json:"name" db:"name"`
	Date     *Timestamp `json:"date" db:"date"`
	Author   *string    `json:"author" db:"author"`
	Location *string    `json:"location" db:"location"`
}
