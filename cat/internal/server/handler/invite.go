package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
)

func (h *Handler) InviteServerMember(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.ListServer")
	defer tracer.Finish(span)

	serverID := mux.Vars(r)["server_id"]
	userID := ctx.Value(middleware.Claim("user_id")).(string)
	err := h.usecase.InviteServerMember(ctx, serverID, userID)
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
		SendJSON(w)

}
