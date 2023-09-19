package usecase

import (
	"context"
	"time"

	"github.com/nekonako/moecord/internal/channel/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type CreateChannelCategoryRequest struct {
	Name     string    `json:"name"`
	ServerID ulid.ULID `json:"server_id"`
}

func (u *UseCase) CreateChannelCategory(ctx context.Context, p CreateChannelCategoryRequest) error {
	span := tracer.SpanFromContext(ctx, "usecase.CreateChannelCategory")
	defer tracer.Finish(span)

	e := repo.ChannelCategory{
		ID:        ulid.Make(),
		ServerID:  p.ServerID,
		Name:      p.Name,
		CreatedAt: time.Now().UTC(),
	}

	err := u.repo.SaveChannelCategory(ctx, e)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}

	return nil
}
