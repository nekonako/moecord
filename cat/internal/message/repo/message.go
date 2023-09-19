package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Message struct {
	ID        ulid.ULID    `db:"id" json:"id"`
	ChannelID ulid.ULID    `db:"channel_id" json:"channel_id"`
	SenderID  ulid.ULID    `db:"sender_id" json:"sender_id"`
	Content   string       `db:"content" json:"content"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	Username  string       `db:"username"`
}

func (r *Repository) SaveMessage(ctx context.Context, message Message) error {

	ctx, span := tracer.Start(ctx, "repo.SaveMessage")
	defer tracer.Finish(span)

	query := `
    INSERT INTO message
    (
        id,
        sender_id,
        channel_id,
        content,
        created_at,
        updated_at
    )
    VALUES (:id, :sender_id, :channel_id, :content, :created_at, :updated_at)`

	_, err := r.postgres.NamedExecContext(ctx, query, message)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed save message")
		return err
	}

	return nil

}

func (r *Repository) ListMessages(ctx context.Context, channelID ulid.ULID) ([]Message, error) {

	span := tracer.SpanFromContext(ctx, "repo.ListMessage")
	defer tracer.Finish(span)

	messages := []Message{}
	query := `
        SELECT
            m.id,
            m.sender_id,
            m.channel_id,
            m.content,
            m.created_at,
            m.updated_at,
            u.username
        FROM message AS m
        INNER JOIN users AS u ON u.id = m.sender_id
        WHERE m.channel_id = $1
    `
	err := r.postgres.SelectContext(ctx, &messages, query, channelID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed get messages")
		return messages, err
	}

	return messages, nil

}
