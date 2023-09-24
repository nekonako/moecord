package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type ListServerResponse struct {
	ID            ulid.ULID `json:"id"`
	OwnerID       ulid.ULID `json:"owner_id"`
	Name          string    `json:"name"`
	DirectMessage bool      `json:"direct_message"`
	Avatar        string    `json:"avatar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (u *UseCase) ListServer(ctx context.Context, userID string) ([]ListServerResponse, error) {

	span := tracer.SpanFromContext(ctx, "usecase.ListServer")
	defer tracer.Finish(span)

	id, err := ulid.Parse(userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return nil, errors.New("invalid user id")
	}

	server, err := u.repo.ListServerUser(ctx, id)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return nil, errors.New("failed get list server")
	}

	res := make([]ListServerResponse, len(server))
	for i, v := range server {
		res[i] = ListServerResponse(v)
	}

	return res, nil

}
