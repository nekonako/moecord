package usecase

import (
	"context"
	"errors"
	"io"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/nekonako/moecord/internal/profile/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type UpdateProfileRequest struct {
	ID       ulid.ULID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Avatar   io.Reader `json:"avatar"`
}

func (u *UseCase) UpdateProfile(ctx context.Context, p UpdateProfileRequest) error {
	span := tracer.SpanFromContext(ctx, "usecase.ListMessage")
	defer tracer.Finish(span)

	user, err := u.repo.GetUserByID(ctx, p.ID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
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
		user.Avatar = res.URL
	}

	err = u.repo.UpdateUser(ctx, repo.User{
		ID:       p.ID,
		Username: p.Username,
		Email:    p.Email,
		Avatar:   user.Avatar,
	})
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Ctx(ctx).Msg(err.Error())
		return errors.New("failed get list server")
	}

	return nil

}
