package handler

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/internal/auth/usecase"
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

	ctx, span := tracer.Start(r.Context(), "oauth", "handler.authorization")
	defer tracer.Finish(span)

	reqBody := usecase.OauthRequest{
		Provider: mux.Vars(r)["provider"],
	}

	url, err := h.usecase.Authorization(ctx, reqBody)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		if ok, ve := validation.IsValidationError(err); ok {
			api.NewHttpResponse().
				WithCode(http.StatusBadRequest).
				WithErrors(ve).
				WitMessage("bad request").
				SendJSON(w)
			return
		}

		if errors.Is(err, usecase.ErrInvalidOauthProvider) {
			api.NewHttpResponse().
				WithCode(http.StatusBadRequest).
				WithError(err).
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
