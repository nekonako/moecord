package message

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/internal/message/handler"
	"github.com/nekonako/moecord/internal/message/repo"
	"github.com/nekonako/moecord/internal/message/usecase"
	"github.com/nekonako/moecord/pkg/middleware"
)

type Message struct {
	Config *config.Config
	Infra  *infra.Infra
}

func New(
	c *config.Config,
	infra *infra.Infra,
) *Message {
	return &Message{
		Config: c,
		Infra:  infra,
	}
}

func (o *Message) InitRouter(r *mux.Router) {

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.Use(middleware.Authentication(o.Config))

	repo := repo.New(o.Infra.Postgres)
	u := usecase.New(o.Config, o.Infra, repo)
	h := handler.New(o.Config, u)

	v1.HandleFunc("/messages", h.SaveMessage).Methods(http.MethodPost)
	v1.HandleFunc("/messages/channels/{channel_id}", h.ListMessage).Methods(http.MethodGet)

}
