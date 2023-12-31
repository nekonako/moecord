package infra

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/pkg/log"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/redis/go-redis/v9"
	"github.com/scylladb/gocqlx/v2"
)

type Infra struct {
	Postgres   *sqlx.DB
	Cloudinary *cloudinary.Cloudinary
	Redis      *redis.Client
	Nats       *nats.Conn
	Scylla     *gocqlx.Session
}

func New(c *config.Config) *Infra {
	if c.Apm.Enable {
		tracer.Init(c.Apm.Host, c.Apm.ServiceName)
	}
	log.Init(c.Apm.LogLevel)

	return &Infra{
		Redis:      newRedis(c),
		Postgres:   newPostgres(c),
		Cloudinary: newCloudinary(c),
		Nats:       newNats(c),
	}
}
