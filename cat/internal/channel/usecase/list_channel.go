package usecase

import (
	"context"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ListChannelResponse struct {
	ID          ulid.ULID `json:"id"`
	ServerID    ulid.ULID `json:"server_id"`
	Name        string    `json:"name"`
	ChannelType string    `json:"channel_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *UseCase) ListChannel(ctx context.Context, userID, serverId ulid.ULID) ([]ListChannelResponse, error) {

	span := tracer.SpanFromContext(ctx, "usecase.ListChannel")
	defer tracer.Finish(span)

	channels, err := u.repo.ListChannel(ctx, userID, serverId)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return nil, err
	}

	res := make([]ListChannelResponse, len(channels))
	for i, v := range channels {
		res[i] = ListChannelResponse(v)
	}

	return res, nil

}
