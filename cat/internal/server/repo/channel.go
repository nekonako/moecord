package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Channel struct {
	ID          ulid.ULID `db:"id"`
	ServerID    ulid.ULID `db:"server_id"`
	CategoryID  ulid.ULID `db:"category_id"`
	Name        string    `db:"name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	ChannelType string    `db:"channel_type"`
}

type ChannelMember struct {
	ID        ulid.ULID `db:"id"`
	ChannelID ulid.ULID `db:"channel_id"`
	UserID    ulid.ULID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Repository) CreateChannel(ctx context.Context, tx *sqlx.Tx, channel []Channel) error {

	span := tracer.SpanFromContext(ctx, "repo.CreateChannel")
	defer tracer.Finish(span)

	query := `
	INSERT INTO channel (
		id,
        server_id,
        category_id,
		name,
		channel_type,
		created_at,
		updated_at
    ) VALUES (:id, :server_id, :category_id, :name, :channel_type, :created_at, :updated_at);
	`

	_, err := tx.NamedExecContext(ctx, query, channel)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	return nil

}

func (r *Repository) CreateChannelMember(ctx context.Context, tx *sqlx.Tx, member []ChannelMember) error {

	span := tracer.SpanFromContext(ctx, "repo.CreateChannelMember")
	defer tracer.Finish(span)

	query := `
		INSERT INTO channel_member (
			id,
			channel_id,
			user_id,
			created_at
		) VALUES (:id, :channel_id, :user_id, :created_at)
	`

	_, err := tx.NamedExecContext(ctx, query, member)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	return nil

}

func (r *Repository) GetPublicChannel(ctx context.Context, serverID ulid.ULID) ([]Channel, error) {

	span := tracer.SpanFromContext(ctx, "repo.GetPublicChannel")
	defer tracer.Finish(span)

	query := `
	SELECT
		id,
		server_id,
        category_id,
		name,
        channel_type,
		created_at,
		updated_at
	FROM channel WHERE server_id = $1 AND is_private = false
	`

	result := []Channel{}
	err := r.postgres.SelectContext(ctx, &result, query, serverID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return result, err
	}

	return result, nil

}
