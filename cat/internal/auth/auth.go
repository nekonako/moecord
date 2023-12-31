package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/internal/auth/handler"
	"github.com/nekonako/moecord/internal/auth/repo"
	"github.com/nekonako/moecord/internal/auth/usecase"
)

type Auth struct {
	Config *config.Config
	Infra  *infra.Infra
}

func New(
	c *config.Config,
	infra *infra.Infra,
) *Auth {
	return &Auth{
		Config: c,
		Infra:  infra,
	}
}

func (o *Auth) InitRouter(r *mux.Router) {

	repo := repo.New(o.Infra.Postgres)
	u := usecase.New(o.Config, o.Infra, repo)
	h := handler.New(o.Config, u)

	r.HandleFunc("/v1/login/oauth/authorization/{provider}", h.Authorization).Methods(http.MethodGet)
	r.HandleFunc("/v1/login/oauth/callback/{provider}", h.Callback).Methods(http.MethodPost)

}
