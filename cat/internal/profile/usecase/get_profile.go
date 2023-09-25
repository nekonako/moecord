package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Profile struct {
	ID        ulid.ULID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UseCase) GetProfile(ctx context.Context, userID string) (Profile, error) {
	span := tracer.SpanFromContext(ctx, "usecase.GetProfile")
	defer tracer.Finish(span)

	res := Profile{}
	id, _ := ulid.Parse(userID)

	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return res, errors.New("failed get profile")
	}

	res = Profile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return res, nil

}
