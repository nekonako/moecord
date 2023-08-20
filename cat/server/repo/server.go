package repo

import (
	"context"
	"time"

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

func (r *Repository) ListServerUser(ctx context.Context, userID ulid.ULID) ([]Server, error) {

	span := tracer.SpanFromContext(ctx, "repo.ListServer")
	defer tracer.Finish(span)

	query := `
	SELECT 
		s.id,
		s.owner_id,
		s.name,
		s.direct_message,
		s.created_at,
		s.updated_at
	FROM server_member AS sm
	INNER JOIN servers AS s ON s.id = sm.server_id
	WHERE sm.user_id = $1 ORDER BY sm.id
	`

	result := []Server{}
	err := r.postgres.SelectContext(ctx, &result, query, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed get server")
		return result, err
	}

	return result, nil

}
