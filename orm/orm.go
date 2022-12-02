package orm

import (
	"github.com/jmoiron/sqlx"
	"github.com/oxodao/photomaton/utils"

	_ "github.com/mattn/go-sqlite3"
)

var GET *ORM

type ORM struct {
	db       *sqlx.DB
	AppState AppState
	Events   Events
}

func Load() error {
	db := sqlx.MustConnect("sqlite3", utils.GetPath("photomaton.db"))

	GET = &ORM{
		db:       db,
		AppState: AppState{db},
		Events:   Events{db},
	}

	return nil
}
