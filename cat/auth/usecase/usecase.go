package usecase

import (
	"github.com/nekonako/moecord/auth/repo"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
)

type UseCase struct {
	config *config.Config
	infra  *infra.Infra
	repo   *repo.Repository
}

func New(
	c *config.Config,
	infra *infra.Infra,
	repo *repo.Repository,
) *UseCase {
	return &UseCase{
		config: c,
		infra:  infra,
		repo:   repo,
	}
}
