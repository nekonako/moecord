package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/auth/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/nekonako/moecord/pkg/validation"
	"github.com/rs/zerolog/log"
)

type ApiError struct {
	Field string
	Msg   string
}

func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {

	_, span := tracer.Start(r.Context(), "oauth", "redirect")
	defer tracer.Finish(span)

	reqBody := usecase.OauthRequest{
		Provider: mux.Vars(r)["provider"],
	}

	span.AddEvent(fmt.Sprintf("aouth %s", reqBody.Provider))
	log.Info().Msg("start oauth")
	url, err := h.usecase.Authorization(reqBody)
	if err != nil {
		log.Error().Msg(err.Error())
		if ok, ve := validation.IsValidationError(err); ok {
			api.NewHttpResponse().
				WithCode(http.StatusBadRequest).
				WithErrors(ve).
				WitMessage("bad request").
				SendJSON(w)
			return
		}

		api.NewHttpResponse().
			WithCode(http.StatusInternalServerError).
			WitMessage("internal server error").
			SendJSON(w)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
	w.Write(nil)

}
