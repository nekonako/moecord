package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
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
		ID:        ulid.Make(),
		ChannelID: m.ChannelID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		CreatedAt: now,
		UpdatedAt: sql.NullTime{},
	}

	err := u.repo.SaveMessage(ctx, message)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}
	go u.publishMessage("NEW_CHANNEL_MESSAGE", message)
	return nil

}

func (u *UseCase) publishMessage(topic string, m any) {

	ctx, span := tracer.Start(context.Background(), "usecase.SaveMessage")
	defer tracer.Finish(span)
	tick := time.NewTicker(1)
	maxRetry := 3
	b, err := json.Marshal(m)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return
	}

	retry := 0
	for retry < maxRetry {
		<-tick.C
		err := u.infra.Nats.Publish(topic, b)
		if err == nil {
			log.Info().Ctx(ctx).Msg("success publish message with topic : " + topic)
			return
		}
		retry++
		tick = time.NewTicker(time.Duration(retry) * time.Second)
	}
}
