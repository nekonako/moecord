package handler

import (
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/server/usecase"
)

type Handler struct {
	config  *config.Config
	usecase *usecase.UseCase
}

func New(c *config.Config, u *usecase.UseCase) *Handler {
	return &Handler{
		config:  c,
		usecase: u,
	}
}
