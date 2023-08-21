package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/pkg/api"
)

type Claim string

func Authentication(c *config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractTokenFromHeader(r)
			if tokenString == "" {
				api.NewHttpResponse().
					WithCode(http.StatusUnauthorized).
					WitMessage("Unauthorized").
					SendJSON(w)
				return
			}

			claim := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
				return []byte(c.JWT.PrivateKey), nil
			})
			if err != nil || !token.Valid {
				api.NewHttpResponse().
					WithCode(http.StatusUnauthorized).
					WitMessage("Unauthorized").
					SendJSON(w)
				return
			}
			ctx := context.WithValue(r.Context(), Claim("user_id"), claim["sub"])
			req := r.WithContext(ctx)
			*r = *req
			next.ServeHTTP(w, r)
		})
	}

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
