package usecase

import (
	"context"
	"time"

	"github.com/nekonako/moecord/internal/channel/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type CreateChannelRequest struct {
	Name       string    `json:"name"`
	ServerID   ulid.ULID `json:"server_id"`
	CategoryID ulid.ULID `json:"category_id"`
	Type       string    `json:"type"`
}

func (u *UseCase) CreateChannel(ctx context.Context, p CreateChannelRequest) error {
	span := tracer.SpanFromContext(ctx, "usecase.CreateChannel")
	defer tracer.Finish(span)

	now := time.Now().UTC()
	e := repo.Channel{
		ID:          ulid.Make(),
		ServerID:    p.ServerID,
		Name:        p.Name,
		CategoryID:  p.CategoryID,
		ChannelType: p.Type,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := u.repo.CreateChannel(ctx, e)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	return nil
}
