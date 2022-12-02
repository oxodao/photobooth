package orm

import (
	"github.com/jmoiron/sqlx"
	"github.com/oxodao/photomaton/models"
)

type AppState struct {
	db *sqlx.DB
}

func (x *AppState) GetState() (models.AppState, error) {
	as := models.AppState{}

	row := x.db.QueryRowx(`
		SELECT hwid, token, current_event
		FROM app_state
		WHERE id = 1
	`)
	if row.Err() != nil {
		return as, row.Err()
	}

	err := row.StructScan(&as)
	if err != nil {
		return as, err
	}

	return as, nil
}

func (x *AppState) CreateState(state models.AppState) error {
	_, err := x.db.Exec(`
		INSERT INTO app_state(id, hwid, token, current_event)
		VALUES (1, ?, NULL, NULL)
	`, state.HardwareID)

	return err
}

func (x *AppState) SetState(state models.AppState) error {

	return nil
}
