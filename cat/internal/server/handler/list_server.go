package handler

import (
	"net/http"

	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
)

func (h *Handler) ListServer(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "oauth", "handler.ListServer")
	defer tracer.Finish(span)

	userID := ctx.Value(middleware.Claim("user_id")).(string)
	res, err := h.usecase.ListServer(ctx, userID)
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
