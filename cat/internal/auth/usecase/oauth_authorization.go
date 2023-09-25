package usecase

import (
	"context"

	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/nekonako/moecord/pkg/validation"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidOauthProvider = errors.New("invalid oauth provider")
)

type OauthRequest struct {
	Provider string `json:"provider" validate:"required,oneof=github google discord"`
}

func (r OauthRequest) validate() error {
	return validation.Validate.Struct(&r)
}

func (u *UseCase) Authorization(ctx context.Context, input OauthRequest) (string, error) {

	span := tracer.SpanFromContext(ctx, "usecase.authorization")
	defer tracer.Finish(span)

	if err := input.validate(); err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return "", err
	}

	oauthURL := ""
	state := uuid.NewString()
	switch input.Provider {
	case "github":
		oauthURL = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s", u.config.Oauth.Github.AuthURL, u.config.Oauth.Github.ClientID, u.config.Oauth.RedirectURI+input.Provider, u.config.Oauth.Google.Scope, state)
		return oauthURL, nil
	case "google":
		oauthURL = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s&response_type=code", u.config.Oauth.Google.AuthURL, u.config.Oauth.Google.ClientID, u.config.Oauth.RedirectURI+input.Provider, u.config.Oauth.Google.Scope, state)
		return oauthURL, nil
	case "discord":
		oauthURL = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&state=%s&response_type=code", u.config.Oauth.Discord.AuthURL, u.config.Oauth.Discord.ClientID, u.config.Oauth.RedirectURI+input.Provider, u.config.Oauth.Discord.Scope, state)
		return oauthURL, nil
	default:
		tracer.SpanError(span, ErrInvalidOauthProvider)
		return "", ErrInvalidOauthProvider
	}

}
