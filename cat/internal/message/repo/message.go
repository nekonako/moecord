package repo

import (
	"context"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

type Message struct {
	ID        []byte
	ChannelID []byte
	SenderID  []byte
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var MessageMetadata = table.Metadata{
	Name:    "messages",
	Columns: []string{"id", "channel_id", "sender_id", "content", "created_at", "updated_at"},
	PartKey: []string{"id"},
}

var MessageByChannelIDMetadata = table.Metadata{
	Name:    "messages_by_channel_id",
	Columns: []string{"id", "channel_id", "sender_id", "content", "created_at", "updated_at"},
	PartKey: []string{"channel_id"},
	SortKey: []string{"id"},
}

var MessageTable = table.New(MessageMetadata)
var MessageByChannelIDTable = table.New(MessageByChannelIDMetadata)

func (r *Repository) SaveMessage(ctx context.Context, scylla *gocqlx.Session, message Message) error {

	_, span := tracer.Start(ctx, "repo.SaveMessage")
	defer tracer.Finish(span)

	q := scylla.Query(MessageTable.Insert()).BindStruct(&message)
	if err := q.ExecRelease(); err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed save message")
		return err
	}

	return nil

}

func (r *Repository) ListMessages(ctx context.Context, userID, channelID []byte) ([]Message, error) {

	span := tracer.SpanFromContext(ctx, "repo.SaveMessage")
	defer tracer.Finish(span)

	messages := []Message{}
	q := r.scylla.Query(MessageByChannelIDTable.Select()).BindMap(qb.M{"channel_id": channelID})
	if err := q.Select(&messages); err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed save message")
		return messages, err
	}

	return messages, nil

}
