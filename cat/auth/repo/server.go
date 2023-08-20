package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Server struct {
	ID            ulid.ULID `db:"id"`
	OwnerID       ulid.ULID `db:"owner_id"`
	Name          string    `db:"name"`
	DirectMessage bool      `db:"direct_message"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (r *Repository) SaveServer(ctx context.Context, tx *sqlx.Tx, server Server) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveServer")
	defer tracer.Finish(span)

	query := `
		INSERT INTO servers (
			id,
			owner_id,
			name,
			direct_message,
			created_at,
			updated_at
		) VALUES (:id, :owner_id, :name, :direct_message, :created_at, :updated_at)
	`

	_, err := tx.NamedExecContext(ctx, query, server)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert user")
		return err
	}

	return nil

}
