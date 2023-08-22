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

type Channel struct {
	Config *config.Config
	Infra  *infra.Infra
}

func New(
	c *config.Config,
	infra *infra.Infra,
) *Channel {
	return &Channel{
		Config: c,
		Infra:  infra,
	}
}

func (o *Channel) InitRouter(r *mux.Router) {

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.Use(middleware.Authentication(o.Config))

	repo := repo.New(o.Infra.Postgres)
	u := usecase.New(o.Config, o.Infra, repo)
	h := handler.New(o.Config, u)

	v1.HandleFunc("/channels", h.ListChannel).Methods(http.MethodGet)

}
