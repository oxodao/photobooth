package orm

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oxodao/photobooth/models"
)

type Events struct {
	db *sqlx.DB
}

func (e *Events) GetEvents() ([]models.Event, error) {
	events := []models.Event{}

	// @TODO created_at && order by created_at desc
	rows, err := e.db.Queryx(`
		SELECT id, name, date, author, location
		FROM event
		ORDER BY name
	`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		evt := models.Event{}
		err = rows.StructScan(&evt)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	return events, nil
}

func (e *Events) GetEvent(id int64) (*models.Event, error) {
	row := e.db.QueryRowx(`
		SELECT id, name, date, author, location
		FROM event
		WHERE id = ?
	`, id)

	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}

		return nil, row.Err()
	}

	evt := models.Event{}
	err := row.StructScan(&evt)

	return &evt, err
}

func (e *Events) GetImage(id int64) (*models.Image, error) {
	row := e.db.QueryRowx(`
		SELECT id, created_at, unattended, event_id
		FROM image
		WHERE id = ?
	`, id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	pct := models.Image{}
	err := row.StructScan(&pct)

	return &pct, err
}

func (e *Events) InsertImage(eventId int64, unattended bool) (*models.Image, error) {
	currTime := time.Now().Unix()

	row := e.db.QueryRowx(`
		INSERT INTO image(event_id, unattended, created_at)
		VALUES (?, ?, ?)
		RETURNING *
	`, eventId, unattended, currTime)
	if row.Err() != nil {
		return nil, row.Err()
	}

	img := models.Image{}
	err := row.StructScan(&img)
	if err != nil {
		return nil, err
	}

	return &img, nil
}