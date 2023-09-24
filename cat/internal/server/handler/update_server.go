package handler

import (
	"net/http"

	"github.com/nekonako/moecord/internal/server/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (h *Handler) UpdateServer(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.CreateServer")
	defer tracer.Finish(span)

	f, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}
	if f != nil {
		defer f.Close()
	}

	id, err := ulid.Parse(r.FormValue("id"))
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	reqBody := usecase.UpdateServerRequest{
		ID:     id,
		Name:   r.FormValue("name"),
		Avatar: f,
	}

	if err := h.usecase.UpdateServer(ctx, reqBody); err != nil {
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
