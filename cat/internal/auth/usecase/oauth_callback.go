package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nekonako/moecord/internal/auth/repo"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/nekonako/moecord/pkg/util"
	"github.com/nekonako/moecord/pkg/validation"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrFailedTokenExchange = errors.New("failed token exchange")
	ErrFailedGetUserInfo   = errors.New("failed get user info")
)

const servervatar = "https://res.cloudinary.com/da9bihi2v/image/upload/v1695292886/moecord/tzqgovnkhdaqcjnyj7t3.jpg"

type CallbackRequest struct {
	Provider          string `json:"provider" validate:"required,oneof=github google discord"`
	AuthorizationCode string `json:"authorization_code" validate:"required"`
	State             string `json:"state" validate:"required"`
}

func (r CallbackRequest) validate() error {
	return validation.Validate.Struct(&r)
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

type discordTokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
}

type response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *UseCase) Callback(ctx context.Context, input CallbackRequest) (response, error) {

	var (
		err   error
		email string
		r     response
		now   = time.Now().UTC()
	)

	span := tracer.SpanFromContext(ctx, "usecase.callback")
	defer tracer.Finish(span)

	if err := input.validate(); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return r, err
	}

	switch input.Provider {
	case "github":
		email, err = u.githubTokenExchange(input.AuthorizationCode)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return r, err
		}
	case "google":
		email, err = u.googleTokenExchange(input.AuthorizationCode)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return r, err
		}
	case "discord":
		email, err = u.discordTokenExchange(input.AuthorizationCode)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return r, err
		}
	default:
		return r, errors.New("oauth provider not implemented")
	}

	userID := ulid.Make()
	username, _ := util.RandomHex(8)
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return r, errors.New("failed get user")
	}

	newUser := err == sql.ErrNoRows
	if newUser {
		user = repo.User{
			ID:        userID,
			Username:  username,
			Email:     email,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	tx, err := u.infra.Postgres.BeginTxx(ctx, nil)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return r, err
	}

	defer tx.Rollback()

	err = u.repo.SaveOrUpdateUser(ctx, tx, user)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return r, errors.New("failed create user")
	}

	if newUser {
		serverID := ulid.Make()
		server := repo.Server{
			ID:            serverID,
			OwnerID:       user.ID,
			Name:          "@me",
			DirectMessage: true,
			Avatar:        servervatar,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err = u.repo.SaveServer(ctx, tx, server); err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return r, errors.New("failed create server")
		}
		member := repo.ServerMember{
			ID:        ulid.Make(),
			ServerID:  serverID,
			UserID:    user.ID,
			CreatedAt: now,
		}
		if err = u.repo.SaveServerMember(ctx, tx, member); err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return r, errors.New("failed create server")
		}

		directMessageChannelCategoryID := ulid.Make()
		channelCategory := []repo.ChannelCategory{
			{
				ID:        directMessageChannelCategoryID,
				ServerID:  serverID,
				Name:      "Direct Messages",
				CreatedAt: now,
			},
		}

		err = u.repo.CreateChannelCategory(ctx, tx, channelCategory)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Msg(err.Error())
			return r, errors.New("failed create channel category")
		}

	}

	err = tx.Commit()
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return r, errors.New("failed commit transaction")
	}

	accessToken, refreshToken, err := u.generateToken(user.ID, user.Username)
	if err != nil {
		tracer.SpanError(span, err)
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
		return "", ErrFailedTokenExchange
	}

	c := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, tokenExchangeURL, bytes.NewBuffer(byteTokenExchange))
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed exchange token, with http status code %d", res.StatusCode))
		return "", ErrFailedTokenExchange
	}

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}

	ru, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
	}
	ru.Header.Set("Authorization", "Bearer "+r.AccessToken)

	res, err = c.Do(ru)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
	}

	if res.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed get user information, with http status code %d", res.StatusCode))
		return "", ErrFailedGetUserInfo
	}

	var user struct {
		Email string `json:"email"`
	}

	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
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
		return "", ErrFailedTokenExchange
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed exchange token, with http status code %d", res.StatusCode))
		return "", ErrFailedTokenExchange
	}

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}

	ru, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + r.AccessToken)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
	}

	if ru.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed get user information, with http status code %d", ru.StatusCode))
		return "", ErrFailedGetUserInfo
	}

	var user struct {
		Email string `json:"email"`
	}

	if err = json.NewDecoder(ru.Body).Decode(&user); err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
	}

	return user.Email, nil

}

func (u *UseCase) discordTokenExchange(authCode string) (string, error) {

	data := url.Values{}
	data.Set("client_id", u.config.Oauth.Discord.ClientID)
	data.Set("client_secret", u.config.Oauth.Discord.ClientSecret)
	data.Set("code", authCode)
	data.Set("redirect_uri", u.config.Oauth.RedirectURI+u.config.Oauth.Discord.Name)
	data.Set("grant_type", "authorization_code")

	tokenExchangeURL := u.config.Oauth.Discord.TokenExchangeURL
	r := discordTokenExchangeResponse{}
	c := http.DefaultClient

	req, err := http.NewRequest(http.MethodPost, tokenExchangeURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed exchange token, with http status code %d", res.StatusCode))
		return "", ErrFailedTokenExchange
	}

	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedTokenExchange
	}

	ru, err := http.NewRequest(http.MethodGet, "https://discord.com/api/users/@me", nil)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
	}
	ru.Header.Set("Authorization", "Bearer "+r.AccessToken)

	res, err = c.Do(ru)
	if err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
	}

	if res.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("fialed get user information, with http status code %d", res.StatusCode))
		return "", ErrFailedGetUserInfo
	}

	var user struct {
		Email string `json:"email"`
	}

	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		log.Error().Msg(err.Error())
		return "", ErrFailedGetUserInfo
	}

	return user.Email, nil

}

func (u *UseCase) generateToken(id ulid.ULID, username string) (string, string, error) {

	now := time.Now().UTC()
	secretKey := []byte(u.config.JWT.PrivateKey)
	claims := jwt.MapClaims{
		"iat":      now.Unix(),
		"exp":      now.Add(time.Minute * time.Duration(u.config.JWT.AccessTokenDuration)).Unix(),
		"sub":      id,
		"username": username,
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
