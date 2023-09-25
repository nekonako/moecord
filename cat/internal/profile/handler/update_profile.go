package handler

import (
	"net/http"

	"github.com/nekonako/moecord/internal/profile/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.updateProfile")
	defer tracer.Finish(span)

	f, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
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
		log.Error().Ctx(ctx).Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	reqBody := usecase.UpdateProfileRequest{
		ID:       id,
		Avatar:   f,
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
	}

	if err := h.usecase.UpdateProfile(ctx, reqBody); err != nil {
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
		WitMessage("Success").
		SendJSON(w)

}
