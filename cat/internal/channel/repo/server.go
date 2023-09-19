package repo

import (
	"context"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Server struct {
	ID   ulid.ULID `db:"id"`
	Name string    `db:"name"`
}

func (r *Repository) GetFistServer(ctx context.Context) (Server, error) {
	span := tracer.SpanFromContext(ctx, "repo.GetFistServer")
	defer tracer.Finish(span)

	query := `
	SELECT
	    id,
        name
    FROM servers ORDER BY id ASC LIMIT 1
    `

	result := Server{}
	err := r.postgres.GetContext(ctx, &result, query)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed get first server")
		return result, err
	}

	return result, nil

}
