package usecase

import (
	"context"

	"time"

	"github.com/nekonako/moecord/internal/channel/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ChannelCategory struct {
	CategoryID   ulid.ULID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Channels     []Channel `json:"channels"`
}

type Channel struct {
	ID          ulid.ULID `json:"id"`
	ServerID    ulid.ULID `json:"server_id"`
	Name        string    `json:"name"`
	ChannelType string    `json:"channel_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *UseCase) ListChannel(ctx context.Context, userID ulid.ULID, serverId string) ([]ChannelCategory, error) {

	span := tracer.SpanFromContext(ctx, "usecase.ListChannel")
	defer tracer.Finish(span)

	res := make([]ChannelCategory, 0)
	channelCategory := []repo.ChannelCategory{}
	if serverId == "@me" {
		server, err := u.repo.GetFistServer(ctx)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return nil, err
		}
		cc, err := u.repo.ListChannelCategory(ctx, server.ID)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return nil, err
		}
		channelCategory = cc
	} else {
		sID, _ := ulid.Parse(serverId)
		cc, err := u.repo.ListChannelCategory(ctx, sID)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return nil, err
		}
		channelCategory = cc
	}

	for _, v := range channelCategory {
		channels, err := u.repo.ListChannel(ctx, userID, v.ID)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return nil, err
		}

		cn := make([]Channel, len(channels))
		for i, c := range channels {
			cn[i] = Channel{
				ID:          c.ID,
				ServerID:    c.ServerID,
				Name:        c.Name,
				ChannelType: c.ChannelType,
				CreatedAt:   c.CreatedAt,
				UpdatedAt:   c.UpdatedAt,
			}
		}

		res = append(res, ChannelCategory{
			CategoryID:   v.ID,
			CategoryName: v.Name,
			Channels:     cn,
		})
	}

	return res, nil

}
