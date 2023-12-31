package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nekonako/moecord/internal/channel/usecase"
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
	userID := ctx.Value(middleware.Claim("user_id")).(string)
	id, _ := ulid.Parse(userID)
	res, err := h.usecase.ListChannel(ctx, id, serverIDStr)
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

func (h *Handler) CreateChannelCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.CreateChannelCategory")
	defer tracer.Finish(span)

	reqBody := usecase.CreateChannelCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	err := h.usecase.CreateChannelCategory(ctx, reqBody)
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
		SendJSON(w)

}

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.CreateChannel")
	defer tracer.Finish(span)

	reqBody := usecase.CreateChannelRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusBadRequest).
			WitMessage("invalid request").
			SendJSON(w)
		return
	}

	err := h.usecase.CreateChannel(ctx, reqBody)
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
		SendJSON(w)

}

func (h *Handler) Typing(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handler.Typing")
	defer tracer.Finish(span)

	chanStr := mux.Vars(r)["channel_id"]
	channelID, _ := ulid.Parse(chanStr)

	userStr := ctx.Value(middleware.Claim("user_id")).(string)
	userID, _ := ulid.Parse(userStr)

	reqBody := usecase.TypingRequest{
		ChannelID: channelID,
		UserID:    userID,
	}
	err := h.usecase.Typing(ctx, reqBody)
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
