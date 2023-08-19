package infra

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gocql/gocql"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/nekonako/moecord/config"
	"github.com/redis/go-redis/v9"
)

type Infra struct {
	Postgres   *sqlx.DB
	Cloudinary *cloudinary.Cloudinary
	Redis      *redis.Client
	Nats       *nats.Conn
	Scylla     *gocql.Session
}

func New(c *config.Config) *Infra {
	if c.Apm.Enable {
		initTracer(c)
	}
	initLogger(c)

	return &Infra{
		Redis:      newRedis(c),
		Postgres:   newPostgres(c),
		Cloudinary: newCloudinary(c),
		Nats:       newNats(c),
		Scylla:     newScylladb(c),
	}
}
