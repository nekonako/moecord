package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nekonako/moecord/internal/message/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (h *Handler) SaveMessage(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.SaveMessage")
	defer tracer.Finish(span)

	reqBody := usecase.SaveMessagRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	userIDStr := ctx.Value(middleware.Claim("user_id")).(string)
	userID, _ := ulid.Parse(userIDStr)
	reqBody.SenderID = userID

	err := h.usecase.SaveMessage(ctx, reqBody)
	if err != nil {
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
		WithData(reqBody).
		SendJSON(w)

}
