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
	Name        string    `db:"name"`
	ChannelType string    `db:"channel_type"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type ChannelMember struct {
	ID        ulid.ULID `db:"id"`
	ChannelID ulid.ULID `db:"channel_id"`
	UserID    ulid.ULID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type ChannelCategory struct {
	ID        ulid.ULID `db:"id"`
	ServerID  ulid.ULID `db:"server_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Repository) SaveChannel(ctx context.Context, tx *sqlx.Tx, channel []Channel) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveChannel")
	defer tracer.Finish(span)

	query := `
	INSERT INTO channel (
		id,
        server_id,
		name,
		channel_type,
		created_at,
		updated_at
	) VALUES (:id, :server_id, :name, :channel_type, :created_at, :updated_at);
	`

	_, err := tx.NamedExecContext(ctx, query, channel)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed create channel")
		return err
	}

	return nil

}

func (r *Repository) SaveChannelMember(ctx context.Context, tx *sqlx.Tx, member []ChannelMember) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveChannelMember")
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
		log.Error().Err(err).Msg("failed create channel member")
		return err
	}

	return nil

}

func (r *Repository) ListChannel(ctx context.Context, userID, categoryID ulid.ULID) ([]Channel, error) {
	span := tracer.SpanFromContext(ctx, "repo.ListChannel")
	defer tracer.Finish(span)

	query := `
	SELECT
		c.id,
		c.server_id,
		c.name,
		c.channel_type,
		c.created_at,
		c.updated_at
	FROM channel AS c
	INNER JOIN channel_member AS cm ON cm.channel_id = c.id
	WHERE cm.user_id = $1 AND c.category_id = $2
	`

	result := []Channel{}
	err := r.postgres.SelectContext(ctx, &result, query, userID, categoryID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed select channel")
		return result, err
	}

	return result, nil

}

func (r *Repository) SaveChannelCategory(ctx context.Context, c ChannelCategory) error {
	span := tracer.SpanFromContext(ctx, "repo.SaveChannelCategory")
	defer tracer.Finish(span)

	query := `
		INSERT INTO channel_category (
			id,
			server_id,
			name,
			created_at
		) VALUES (:id, :server_id, :name, :created_at)
	`

	_, err := r.postgres.NamedExecContext(ctx, query, c)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed create channel category")
		return err
	}

	return nil

}

func (r *Repository) ListChannelCategory(ctx context.Context, serverID ulid.ULID) ([]ChannelCategory, error) {
	span := tracer.SpanFromContext(ctx, "repo.ListChannelCategory")
	defer tracer.Finish(span)

	query := `
	SELECT
	    id,
        server_id,
        name,
        created_at
    FROM channel_category WHERE server_id = $1
    `

	result := []ChannelCategory{}
	err := r.postgres.SelectContext(ctx, &result, query, serverID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed select channel category")
		return result, err
	}

	return result, nil

}
