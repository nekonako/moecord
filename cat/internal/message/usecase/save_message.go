package usecase

import (
	"context"
	"time"

	"github.com/nekonako/moecord/internal/message/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type SaveMessagRequest struct {
	ChannelID ulid.ULID `json:"channel_id"`
	SenderID  ulid.ULID `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UseCase) SaveMessage(ctx context.Context, m SaveMessagRequest) error {

	ctx, span := tracer.Start(ctx, "usecase.SaveMessage")
	defer tracer.Finish(span)
	now := time.Now().UTC()

	message := repo.Message{
		ID:        ulid.Make().Bytes(),
		ChannelID: m.ChannelID.Bytes(),
		SenderID:  m.SenderID.Bytes(),
		Content:   m.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := u.repo.SaveMessage(ctx, u.infra.Scylla, message)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}

	return nil

}
