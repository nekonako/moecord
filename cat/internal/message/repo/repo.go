package repo

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	postgres *sqlx.DB
}

func New(postgres *sqlx.DB) *Repository {
	return &Repository{
		postgres: postgres,
	}
}
