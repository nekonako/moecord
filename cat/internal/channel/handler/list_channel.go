package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ListChannel(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.ListChannel")
	defer tracer.Finish(span)

	serverIDStr := mux.Vars(r)["server_id"]
	serverID, _ := ulid.Parse(serverIDStr)
	userID := ctx.Value(middleware.Claim("user_id")).(string)
	id, _ := ulid.Parse(userID)
	res, err := h.usecase.ListChannel(ctx, id, serverID)
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
		WithData(res).
		SendJSON(w)

}
