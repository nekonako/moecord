package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

type CallbackRequest struct {
	Provider          string `json:"provider"`
	AuthorizationCode string `json:"authorization_code"`
	State             string `json:"state"`
}

type oauthTokenExchange struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type githubTokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type googleTokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

func (u *UseCase) Callback(input CallbackRequest) (any, error) {

	switch input.Provider {
	case "github":
		return u.githubTokenExchange(input.AuthorizationCode)
	case "google":
		return u.googleTokenExchange(input.AuthorizationCode)
	default:
		return nil, errors.New("oauth provider not implemented")
	}

}

func (u *UseCase) githubTokenExchange(authCode string) (githubTokenExchangeResponse, error) {

	tokenExchange := oauthTokenExchange{
		ClientID:     u.config.Oauth.Github.ClientID,
		ClientSecret: u.config.Oauth.Github.ClientSecret,
		Code:         authCode,
		RedirectURI:  u.config.Oauth.RedirectURI + u.config.Oauth.Github.Name,
	}
	tokenExchangeURL := u.config.Oauth.Github.TokenExchangeURL
	r := githubTokenExchangeResponse{}

	byteTokenExchange, err := json.Marshal(tokenExchange)
	if err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}

	c := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, tokenExchangeURL, bytes.NewBuffer(byteTokenExchange))
	if err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}

	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}

	return r, nil

}

func (u *UseCase) googleTokenExchange(authCode string) (googleTokenExchangeResponse, error) {

	data := url.Values{}
	data.Set("client_id", u.config.Oauth.Google.ClientID)
	data.Set("client_secret", u.config.Oauth.Google.ClientSecret)
	data.Set("code", authCode)
	data.Set("redirect_uri", u.config.Oauth.RedirectURI+u.config.Oauth.Google.Name)
	data.Set("grant_type", "authorization_code")

	tokenExchangeURL := u.config.Oauth.Google.TokenExchangeURL
	r := googleTokenExchangeResponse{}
	c := http.DefaultClient

	req, err := http.NewRequest(http.MethodPost, tokenExchangeURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}

	return r, nil

}
