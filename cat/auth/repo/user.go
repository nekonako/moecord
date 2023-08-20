package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type User struct {
	ID       ulid.ULID `db:"id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
}

func (r *Repository) SaveUser(ctx context.Context, tx *sqlx.Tx, user User) error {

	span := trace.SpanFromContext(ctx)
	defer tracer.Finish(span)

	query := `
		INSERT INTO "user" (
			id,
			username,
			email,
			password
		) VALUES (:id, :username, :email, :password)
	`

	_, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert user")
		return err
	}

	return nil
}
