package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
)

type ListMessages struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	ChannelID string    `json:"channel_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UseCase) ListMessage(ctx context.Context, userID, channelID []byte) ([]ListMessages, error) {

	span := tracer.SpanFromContext(ctx, "usecase.ListMessage")
	defer tracer.Finish(span)

	messages, err := u.repo.ListMessages(ctx, userID, channelID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return nil, errors.New("failed get list server")
	}

	res := make([]ListMessages, len(messages))
	for i, v := range messages {
		res[i] = ListMessages{
			ID:        string(v.ID),
			SenderID:  string(v.SenderID),
			ChannelID: string(v.ChannelID),
			Content:   v.Content,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
	}

	return res, nil

}
