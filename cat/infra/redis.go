package infra

import (
	"context"
	"fmt"

	"github.com/nekonako/moecord/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func newRedis(c *config.Config) *redis.Client {

	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		Username: c.Redis.Username,
		DB:       c.Redis.Database,
	})

	if err := r.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Msg(err.Error())
		panic(err)
	}

	return r

}
