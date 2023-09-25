package usecase

import (
	"context"
	"time"

	"github.com/nekonako/moecord/internal/server/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) InviteServerMember(ctx context.Context, serverID, userID string) error {

	span := tracer.SpanFromContext(ctx, "usecase.InviteServerMember")
	defer tracer.Finish(span)

	sid, err := ulid.Parse(serverID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	uid, err := ulid.Parse(userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	tx, err := u.infra.Postgres.BeginTxx(ctx, nil)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	defer tx.Rollback()

	publicChannel, err := u.repo.GetPublicChannel(ctx, sid)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	newMebers := []repo.ChannelMember{}
	for _, v := range publicChannel {
		newMebers = append(newMebers, repo.ChannelMember{
			ID:        ulid.Make(),
			ChannelID: v.ID,
			UserID:    uid,
			CreatedAt: time.Now().UTC(),
		})
	}

	err = u.repo.CreateChannelMember(ctx, tx, newMebers)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	err = u.repo.CreateServerMember(ctx, tx, repo.ServerMember{
		ID:        ulid.Make(),
		ServerID:  sid,
		UserID:    uid,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	if err = tx.Commit(); err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return err
	}

	return nil

}
