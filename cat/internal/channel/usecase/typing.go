package usecase

import (
	"context"
	"encoding/json"

	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type TypingRequest struct {
	ChannelID ulid.ULID `json:"channel_id"`
	UserID    ulid.ULID `json:"user_id"`
}

type TypingMessage struct {
	ChannelID ulid.ULID `json:"channel_id"`
	UserID    ulid.ULID `json:"user_id"`
	ServerID  ulid.ULID `json:"server_id"`
}

func (u *UseCase) Typing(ctx context.Context, p TypingRequest) error {
	span := tracer.SpanFromContext(ctx, "usecase.Typing")
	defer tracer.Finish(span)

	channel, err := u.repo.GetChannelByID(ctx, p.ChannelID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	m := api.WebSocketMessage[TypingMessage]{
		EventID: "TYPING",
		Data: TypingMessage{
			ChannelID: channel.ID,
			UserID:    p.UserID,
			ServerID:  channel.ServerID,
		},
	}

	bm, _ := json.Marshal(m)
	go u.infra.Nats.Publish("TYPING", bm)

	return nil
}
