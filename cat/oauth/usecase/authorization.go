package usecase

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidOauthProvider = errors.New("invalid oauth provider")
	oauthMapProvider        = map[string]bool{
		"github": true,
		"google": true,
	}
)

type OauthRequest struct {
	Provider string `json:"provider"`
}

func (r OauthRequest) validate(c *config.Config) error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Provider, validation.Required, validation.By(func(value interface{}) error {
			if !oauthMapProvider[r.Provider] {
				return ErrInvalidOauthProvider
			}
			return nil
		})),
	)
}

func (u *UseCase) Authorization(input OauthRequest) (string, error) {

	log.Info().Msg("start oauth")
	if err := input.validate(u.config); err != nil {
		log.Error().Msg(err.Error())
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
	default:
		return "", errors.New("oauth provider not implemented")
	}

}
