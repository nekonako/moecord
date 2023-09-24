package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ListMessages struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"sender_id"`
	ChannelID string    `json:"channel_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
}

func (u *UseCase) ListMessage(ctx context.Context, uID, cID string) ([]ListMessages, error) {

	span := tracer.SpanFromContext(ctx, "usecase.ListMessage")
	defer tracer.Finish(span)

	channelID, _ := ulid.Parse(cID)

	messages, err := u.repo.ListMessages(ctx, channelID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Ctx(ctx).Msg(err.Error())
		return nil, errors.New("failed get list server")
	}

	res := make([]ListMessages, len(messages))
	for i, v := range messages {
		res[i] = ListMessages{
			ID:        v.ID.String(),
			SenderID:  v.SenderID.String(),
			ChannelID: v.ChannelID.String(),
			Content:   v.Content,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt.Time,
			Username:  v.Username,
			Avatar:    v.Avatar,
		}
	}

	return res, nil

}
