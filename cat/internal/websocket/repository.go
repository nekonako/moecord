package websocket

import (
	"context"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (w *Websocket) getUserServer(ctx context.Context, userID ulid.ULID) ([]server, error) {
	ctx, span := tracer.Start(ctx, "websocket.GetUserChannel")
	defer tracer.Finish(span)

	result := []server{}
	query := `
       SELECT
        server_id as id
       FROM server_member WHERE user_id = $1
       `
	err := w.infra.Postgres.SelectContext(ctx, &result, query, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return result, err
	}

	return result, nil

}

func (w *Websocket) getUserChannel(ctx context.Context, userID ulid.ULID) ([]channel, error) {
	ctx, span := tracer.Start(ctx, "websocket.GetUserChannel")
	defer tracer.Finish(span)

	result := []channel{}
	query := `
       SELECT
         distinct(c.id)
       FROM channel AS c
       LEFT JOIN channel_member AS cm ON cm.channel_id = c.id
       WHERE (c.is_private = false OR cm.user_id = $1)
    `
	err := w.infra.Postgres.SelectContext(ctx, &result, query, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return result, err
	}

	return result, nil

}
