package usecase

import (
	"context"

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

func (u *UseCase) UpdateServer(ctx context.Context, p UpdateServerRequest) error {
	span := tracer.SpanFromContext(ctx, "usecase.UpdateServer")
	defer tracer.Finish(span)

	server, err := u.repo.GetServerByID(ctx, p.ID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return errors.New("failed update server")
	}

	if p.Avatar != nil {
		res, err := u.infra.Cloudinary.Upload.Upload(ctx, p.Avatar, uploader.UploadParams{
			Folder:           "moecord",
			FilenameOverride: p.ID.String(),
		})
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return errors.New("failed update server")
		}
		server.Avatar = res.URL
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
		return errors.New("failed update server")
	}

	return nil
}
