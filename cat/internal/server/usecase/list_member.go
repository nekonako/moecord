package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ListServerMemberResponse struct {
	ID        ulid.ULID `json:"id"`
	UserID    ulid.ULID `json:"user_id"`
	ServerID  ulid.ULID `json:"server_ud"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *UseCase) ListServerMember(ctx context.Context, serverID string) ([]ListServerMemberResponse, error) {

	span := tracer.SpanFromContext(ctx, "usecase.ListServer")
	defer tracer.Finish(span)

	id, err := ulid.Parse(serverID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return nil, errors.New("invalid user id")
	}

	server, err := u.repo.ListServerMember(ctx, id)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return nil, errors.New("failed get list server")
	}

	fmt.Println("-=", server, serverID)

	res := make([]ListServerMemberResponse, len(server))
	for i, v := range server {
		res[i] = ListServerMemberResponse{
			ID:        v.ID,
			UserID:    v.UserID,
			ServerID:  v.ServerID,
			Username:  v.Username,
			CreatedAt: v.CreatedAt,
		}
	}

	return res, nil

}
