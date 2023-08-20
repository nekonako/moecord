package infra

import (
	"os"
	"time"

	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger(c *config.Config) {

	now := time.Now().String()
	file, err := os.OpenFile(
		"../../dev/log/"+now+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	level, err := zerolog.ParseLevel(c.Apm.LogLevel)
	if err != nil {
		panic(err)
	}

	log.Logger = zerolog.New(zerolog.MultiLevelWriter(file, os.Stdout)).
		With().
		Caller().
		Logger().
		Level(level)

}
