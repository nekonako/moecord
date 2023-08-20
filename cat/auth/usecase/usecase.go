package usecase

import (
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
)

type UseCase struct {
	config *config.Config
	infra  *infra.Infra
}

func New(c *config.Config) *UseCase {
	return &UseCase{
		config: c,
	}
}
