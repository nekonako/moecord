package handler

import (
	"net/http"

	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.GetProfile")
	defer tracer.Finish(span)

	userID := ctx.Value(middleware.Claim("user_id")).(string)

	res, err := h.usecase.GetProfile(ctx, userID)
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
