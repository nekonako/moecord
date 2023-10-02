package usecase

import (
	"context"
	"fmt"
	"time"

	"errors"
	"io"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/nekonako/moecord/internal/server/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type UpdateServerRequest struct {
	ID     ulid.ULID `json:"id"`
	Name   string    `json:"name"`
	Avatar io.Reader `json:"avatar"`
}

type UpdateServerResponse struct {
	ID            ulid.ULID `json:"id"`
	OwnerID       ulid.ULID `json:"owner_id"`
	Name          string    `json:"name"`
	DirectMessage bool      `json:"direct_message"`
	Avatar        string    `json:"avatar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (u *UseCase) UpdateServer(ctx context.Context, p UpdateServerRequest) (UpdateServerResponse, error) {
	span := tracer.SpanFromContext(ctx, "usecase.UpdateServer")
	defer tracer.Finish(span)

	res := UpdateServerResponse{}
	server, err := u.repo.GetServerByID(ctx, p.ID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return res, errors.New("failed update server")
	}

	if p.Avatar != nil {
		response, err := u.infra.Cloudinary.Upload.Upload(ctx, p.Avatar, uploader.UploadParams{
			Folder:           "moecord",
			FilenameOverride: p.ID.String(),
		})
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return res, errors.New("failed update server")
		}
		server.Avatar = response.URL
	}

	s := repo.Server{
		ID:     p.ID,
		Name:   p.Name,
		Avatar: server.Avatar,
	}

	err = u.repo.UpdateServer(ctx, s)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return res, errors.New("failed update server")
	}

	res = UpdateServerResponse{
		ID:            server.ID,
		OwnerID:       server.OwnerID,
		Name:          p.Name,
		DirectMessage: server.DirectMessage,
		Avatar:        server.Avatar,
		CreatedAt:     server.CreatedAt,
		UpdatedAt:     time.Now().UTC(),
	}

	return res, nil
}
