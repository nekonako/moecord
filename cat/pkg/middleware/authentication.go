package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
)

type Claim string

func Authentication(c *config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(r.Context(), "middleware.authentication")
			defer tracer.Finish(span)
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				err := errors.New("empty token")
				tracer.SpanError(span, err)
				log.Error().Ctx(ctx).Msg(err.Error())
				api.NewHttpResponse().
					WithCode(http.StatusUnauthorized).
					WitMessage("Unauthorized").
					SendJSON(w)
				return
			}

			claim, err := ValidateToken(tokenString, c.JWT.PrivateKey)
			if err != nil {
				err := errors.New("invalid token")
				tracer.SpanError(span, err)
				log.Error().Ctx(ctx).Msg(err.Error())
				api.NewHttpResponse().
					WithCode(http.StatusUnauthorized).
					WitMessage("Unauthorized").
					SendJSON(w)
				return
			}
			xCtx := context.WithValue(ctx, Claim("user_id"), claim["sub"])
			xCtx = context.WithValue(xCtx, Claim("username"), claim["username"])
			req := r.WithContext(xCtx)
			*r = *req
			next.ServeHTTP(w, r)
		})
	}

}

func ValidateToken(tokenString, privateKey string) (jwt.MapClaims, error) {
	claim := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	})

	if err != nil {
		log.Error().Msg(err.Error())
		return claim, errors.New("invalid token")
	}

	if !token.Valid {
		return claim, errors.New("invalid token")
	}

	return claim, nil
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		authParts := strings.Split(authHeader, " ")
		if len(authParts) == 2 && authParts[0] == "Bearer" {
			return authParts[1]
		}
	}
	return ""
}
