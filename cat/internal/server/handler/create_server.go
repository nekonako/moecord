package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nekonako/moecord/internal/server/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateServer(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "server", "handler.CreateServer")
	defer tracer.Finish(span)

	reqBody := usecase.CreateServerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	userID := ctx.Value(middleware.Claim("user_id")).(string)
	id, err := ulid.Parse(userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
	}

	if err := h.usecase.CreateServer(ctx, id, reqBody); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusInternalServerError).
			WitMessage("internal server error").
			SendJSON(w)
		return
	}

	api.NewHttpResponse().
		WithCode(http.StatusOK).
		WitMessage("Success").
		SendJSON(w)

}
