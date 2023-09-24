package repo

import (
	"context"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID       ulid.ULID `db:"id"`
	Username string    `db:"username"`
	Avatar   string    `db:"avatar"`
}

func (r *Repository) GetUser(ctx context.Context, id ulid.ULID) (User, error) {
	ctx, span := tracer.Start(ctx, "repo.SaveMessage")
	defer tracer.Finish(span)

	result := User{}
	query := `SELECT id, username, avatar FROM users WHERE id = $1`
	err := r.postgres.GetContext(ctx, &result, query, id)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed save message")
		return result, err
	}

	return result, err

}
