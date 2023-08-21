package channel

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/internal/channel/handler"
	"github.com/nekonako/moecord/internal/channel/repo"
	"github.com/nekonako/moecord/internal/channel/usecase"
	"github.com/nekonako/moecord/pkg/middleware"
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

	sub := r.PathPrefix("/v1/channels").Subrouter()
	sub.Use(middleware.Authentication(o.Config))

	repo := repo.New(o.Infra.Postgres)
	u := usecase.New(o.Config, o.Infra, repo)
	h := handler.New(o.Config, u)

	sub.HandleFunc("", h.ListChannel).Methods(http.MethodGet)

}
