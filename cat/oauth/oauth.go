package oauth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/oauth/handler"
	"github.com/nekonako/moecord/oauth/usecase"
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

	u := usecase.New(o.Config)
	h := handler.New(o.Config, u)

	r.HandleFunc("/v1/login/oauth/authorization/{provider}", h.Authorization).Methods(http.MethodGet)
	r.HandleFunc("/v1/login/oauth/callback/{provider}", h.Callback).Methods(http.MethodPost)

}