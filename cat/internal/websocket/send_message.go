package websocket

import (
	"context"
	"encoding/json"

	"github.com/gobwas/ws/wsutil"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/rs/zerolog/log"

	"github.com/nekonako/moecord/pkg/tracer"
)

func SendMessage[T any](ctx context.Context, w *Websocket, m api.WebSocketMessage[T], channelID, senderID string) {
	ctx, span := tracer.Start(ctx, "websocket.SendMessage")
	defer tracer.Finish(span)

	if conn, ok := connMap.channels[channelID]; ok {
		mb, err := json.Marshal(m)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Ctx(ctx).Msg(err.Error())
			return
		}
		for _, c := range conn {
			if c.UserID.String() == senderID {
				continue
			}
			wsutil.WriteServerMessage(c.Conn, 0x1, mb)
		}
	}
}
