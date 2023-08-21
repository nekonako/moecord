package usecase

import (
	"context"
	"time"

	"github.com/nekonako/moecord/internal/server/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/nekonako/moecord/pkg/validation"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type CreateServerRequest struct {
	Name string `json:"name" validate:"required"`
}

func (r CreateServerRequest) validate() error {
	return validation.Validate.Struct(&r)
}

func (u *UseCase) CreateServer(ctx context.Context, userID ulid.ULID, input CreateServerRequest) error {
	span := tracer.SpanFromContext(ctx, "usecase.CreateServer")
	defer tracer.Finish(span)

	if err := input.validate(); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}

	serverID := ulid.Make()
	now := time.Now().UTC()

	tx, err := u.infra.Postgres.BeginTxx(ctx, nil)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}

	defer tx.Commit()

	server := repo.Server{
		ID:            serverID,
		OwnerID:       userID,
		Name:          input.Name,
		DirectMessage: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := u.repo.SaveServer(ctx, tx, server); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}

	member := repo.ServerMember{
		ID:        ulid.Make(),
		ServerID:  serverID,
		UserID:    userID,
		CreatedAt: now,
	}

	if err := u.repo.SaveServerMember(ctx, tx, member); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}

	if err = tx.Commit(); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return err
	}

	return nil
}
