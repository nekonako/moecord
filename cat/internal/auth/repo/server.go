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
	Avatar        string    `db:"avatar"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type ServerMember struct {
	ID        ulid.ULID `db:"id"`
	ServerID  ulid.ULID `db:"server_id"`
	UserID    ulid.ULID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
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
            avatar,
			created_at,
			updated_at
        ) VALUES (:id, :owner_id, :name, :direct_message, :avatar, :created_at, :updated_at)
	`

	_, err := tx.NamedExecContext(ctx, query, server)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert user")
		return err
	}

	return nil

}

func (r *Repository) SaveServerMember(ctx context.Context, tx *sqlx.Tx, member ServerMember) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveServer")
	defer tracer.Finish(span)

	query := `
		INSERT INTO server_member (
			id,
			server_id,
			user_id,
			created_at
		) VALUES (:id, :server_id, :user_id, :created_at)
	`

	_, err := tx.NamedExecContext(ctx, query, member)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert server member")
		return err
	}

	return nil

}
