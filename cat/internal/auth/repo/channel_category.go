package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ChannelCategory struct {
	ID        ulid.ULID `db:"id"`
	ServerID  ulid.ULID `db:"server_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Repository) CreateChannelCategory(ctx context.Context, tx *sqlx.Tx, c []ChannelCategory) error {

	ctx, span := tracer.Start(ctx, "repo.CreateChannelCategory")
	defer tracer.Finish(span)

	query := `INSERT INTO channel_category (id, server_id, name, created_at) VALUES (:id, :server_id, :name, :created_at)`
	_, err := tx.NamedExecContext(ctx, query, c)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	return nil

}
