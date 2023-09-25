package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ListMessage(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.ListMessage")
	defer tracer.Finish(span)

	channelID := mux.Vars(r)["channel_id"]
	userID := ctx.Value(middleware.Claim("user_id")).(string)

	res, err := h.usecase.ListMessage(ctx, userID, channelID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusInternalServerError).
			WitMessage("internal server error").
			SendJSON(w)
		return
	}

	api.NewHttpResponse().
		WithCode(http.StatusOK).
		WithData(res).
		SendJSON(w)

}
