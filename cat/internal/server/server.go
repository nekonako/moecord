package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/internal/server/handler"
	"github.com/nekonako/moecord/internal/server/repo"
	"github.com/nekonako/moecord/internal/server/usecase"
	"github.com/nekonako/moecord/pkg/middleware"
)

type Server struct {
	Config *config.Config
	Infra  *infra.Infra
}

func New(
	c *config.Config,
	infra *infra.Infra,
) *Server {
	return &Server{
		Config: c,
		Infra:  infra,
	}
}

func (o *Server) InitRouter(r *mux.Router) {

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.Use(middleware.Authentication(o.Config))

	repo := repo.New(o.Infra.Postgres)
	u := usecase.New(o.Config, o.Infra, repo)
	h := handler.New(o.Config, u)

	v1.HandleFunc("/servers", h.ListServer).Methods(http.MethodGet)
	v1.HandleFunc("/servers", h.CreateServer).Methods(http.MethodPost)
	v1.HandleFunc("/servers", h.UpdateServer).Methods(http.MethodPut)
	v1.HandleFunc("/servers/{server_id}/member", h.ListServerMember).Methods(http.MethodGet)
	v1.HandleFunc("/servers/{server_id}/member/invite", h.InviteServerMember).Methods(http.MethodGet)

}
