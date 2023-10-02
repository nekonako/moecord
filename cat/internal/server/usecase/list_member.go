package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nekonako/moecord/internal/websocket"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ListServerMemberResponse struct {
	ID        ulid.ULID `json:"id"`
	UserID    ulid.ULID `json:"user_id"`
	Avatar    string    `json:"avatar"`
	Username  string    `json:"username"`
	Online    bool      `json:"online"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *UseCase) ListServerMember(ctx context.Context, serverID string) ([]ListServerMemberResponse, error) {

	span := tracer.SpanFromContext(ctx, "usecase.ListServerMember")
	defer tracer.Finish(span)

	id, err := ulid.Parse(serverID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return nil, errors.New("invalid user id")
	}

	server, err := u.repo.ListServerMember(ctx, id)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return nil, errors.New("failed get list server")
	}

	onlineUsers := websocket.GetOnlineUser(u.ws, id)
	fmt.Println("-===", onlineUsers)

	res := make([]ListServerMemberResponse, len(server))
	for i, v := range server {
		res[i] = ListServerMemberResponse{
			ID:        v.ID,
			UserID:    v.UserID,
			Username:  v.Username,
			Avatar:    v.Avatar,
			Online:    onlineUsers[v.UserID.String()],
			CreatedAt: v.CreatedAt,
		}
	}

	return res, nil

}
