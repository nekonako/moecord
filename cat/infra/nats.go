package infra

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
)

func newNats(c *config.Config) *nats.Conn {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s@%s:%d", c.Nats.Username, c.Nats.Password, c.Nats.Host, c.Nats.Port))
	if err != nil {
		log.Fatal().Msg(err.Error())
		return nil
	}
	return nc
}
