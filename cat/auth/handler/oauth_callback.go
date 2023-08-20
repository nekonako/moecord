package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nekonako/moecord/auth/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/validation"
	"github.com/rs/zerolog/log"
)

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {

	reqBody := usecase.CallbackRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	defer r.Body.Close()

	res, err := h.usecase.Callback(reqBody)
	if err != nil {
		log.Error().Msg(err.Error())
		ive, ve := validation.IsValidationError(err)
		switch {
		case ive:
			api.NewHttpResponse().
				WithCode(http.StatusBadRequest).
				WithErrors(ve).
				SendJSON(w)
		case errors.Is(err, usecase.ErrFailedTokenExchange):
			api.NewHttpResponse().
				WithCode(http.StatusBadRequest).
				WithError(err).
				SendJSON(w)
		case errors.Is(err, usecase.ErrFailedGetUserInfo):
			api.NewHttpResponse().
				WithCode(http.StatusUnauthorized).
				WithError(err).
				SendJSON(w)
		default:
			api.NewHttpResponse().
				WithCode(http.StatusInternalServerError).
				WitMessage("internal server error").
				SendJSON(w)
		}
		return
	}

	api.NewHttpResponse().
		WithCode(http.StatusOK).
		WithData(res).
		SendJSON(w)

}
