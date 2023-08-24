package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/scylladb/gocqlx/v2"
)

type Repository struct {
	postgres *sqlx.DB
	scylla   *gocqlx.Session
}

func New(postgres *sqlx.DB, scylla *gocqlx.Session) *Repository {
	return &Repository{
		postgres: postgres,
		scylla:   scylla,
	}
}
