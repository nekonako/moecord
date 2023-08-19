package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	errFailedTokenExchange = errors.New("failed token exchange")
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

type response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *UseCase) Callback(input CallbackRequest) (response, error) {

	var (
		err   error
		email string
		r     response
	)

	switch input.Provider {
	case "github":
		email, err = u.githubTokenExchange(input.AuthorizationCode)
		if err != nil {
			log.Error().Msg(err.Error())
			return r, err
		}
	case "google":
		email, err = u.googleTokenExchange(input.AuthorizationCode)
		if err != nil {
			log.Error().Msg(err.Error())
			return r, err
		}
	default:
		return r, errors.New("oauth provider not implemented")
	}

	id := ulid.Make()
	accessToken, refreshToken, err := u.generateToken(id, email)
	if err != nil {
		log.Error().Msg(err.Error())
		return r, err
	}

	r.AccessToken = accessToken
	r.RefreshToken = refreshToken
	return r, nil

}

func (u *UseCase) githubTokenExchange(authCode string) (string, error) {

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
		return "", err
	}

	c := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, tokenExchangeURL, bytes.NewBuffer(byteTokenExchange))
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed exchange token, with http status code %d", res.StatusCode))
		return "", errFailedTokenExchange
	}

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	ru, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}
	ru.Header.Set("Authorization", "Bearer "+r.AccessToken)

	res, err = c.Do(ru)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	var user struct {
		Email string `json:"email"`
	}

	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	return user.Email, nil

}

func (u *UseCase) googleTokenExchange(authCode string) (string, error) {

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
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed exchange token, with http status code %d", res.StatusCode))
		return "", errFailedTokenExchange
	}

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	ru, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + r.AccessToken)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	var user struct {
		Email string `json:"email"`
	}

	if err = json.NewDecoder(ru.Body).Decode(&user); err != nil {
		log.Error().Msg(err.Error())
		return "", err
	}

	return user.Email, nil

}

func (u *UseCase) generateToken(id ulid.ULID, email string) (string, string, error) {

	now := time.Now()
	secretKey := []byte(u.config.JWT.PrivateKey)
	claims := jwt.MapClaims{
		"iat":   now.Unix(),
		"exp":   now.Add(time.Minute * time.Duration(u.config.JWT.AccessTokenDuration)).Unix(),
		"email": email,
		"sub":   id,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"exp": now.Add(time.Minute * time.Duration(u.config.JWT.RefreshTokenDuration)).Unix(),
		"sub": id,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
