package usecase

import (
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/internal/message/repo"
	"github.com/nekonako/moecord/internal/websocket"
)

type UseCase struct {
	config *config.Config
	infra  *infra.Infra
	repo   *repo.Repository
	ws     *websocket.Websocket
}

func New(
	c *config.Config,
	infra *infra.Infra,
	repo *repo.Repository,
	ws *websocket.Websocket,
) *UseCase {
	return &UseCase{
		config: c,
		infra:  infra,
		repo:   repo,
		ws:     ws,
	}
}
