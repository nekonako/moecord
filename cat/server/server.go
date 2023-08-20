package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/middleware"
	"github.com/nekonako/moecord/server/handler"
	"github.com/nekonako/moecord/server/repo"
	"github.com/nekonako/moecord/server/usecase"
)

type Oauth struct {
	Config *config.Config
	Infra  *infra.Infra
}

func New(
	c *config.Config,
	infra *infra.Infra,
) *Oauth {
	return &Oauth{
		Config: c,
		Infra:  infra,
	}
}

func (o *Oauth) InitRouter(r *mux.Router) {

	sub := r.PathPrefix("/v1/servers").Subrouter()
	sub.Use(middleware.Authentication(o.Config))

	repo := repo.New(o.Infra.Postgres)
	u := usecase.New(o.Config, o.Infra, repo)
	h := handler.New(o.Config, u)

	sub.HandleFunc("", h.ListServer).Methods(http.MethodGet)

}
