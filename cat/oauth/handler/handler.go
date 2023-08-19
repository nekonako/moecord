package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/oauth/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	config  *config.Config
	usecase *usecase.UseCase
}

func New(c *config.Config, u *usecase.UseCase) *Handler {
	return &Handler{
		config:  c,
		usecase: u,
	}
}

func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {

	_, span := tracer.Start(r.Context(), "oauth", "authorization")
	defer tracer.Finish(span)

	reqBody := usecase.OauthRequest{
		Provider: mux.Vars(r)["provider"],
	}

	span.AddEvent(fmt.Sprintf("Authorization_%s", reqBody.Provider))
	log.Info().Msg("start oauth")
	url, err := h.usecase.Authorization(reqBody)
	if err != nil {
		log.Error().Msg(err.Error())
		api.NewResponse().
			WithCode(http.StatusInternalServerError).
			WitMessage("internal server error").
			SendJSON(w)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
	w.Write(nil)

}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {

	reqBody := usecase.CallbackRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Error().Msg(err.Error())
		api.NewResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	defer r.Body.Close()

	b, err := h.usecase.Callback(reqBody)
	if err != nil {
		log.Error().Msg(err.Error())
		api.NewResponse().
			WithCode(http.StatusInternalServerError).
			WitMessage("internal server error").
			SendJSON(w)
		return
	}

	api.NewResponse().
		WithCode(http.StatusOK).
		WithData(b).
		SendJSON(w)

}
