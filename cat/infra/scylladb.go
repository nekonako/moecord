package infra

import (
	"github.com/gocql/gocql"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
)

func newScylladb(c *config.Config) *gocql.Session {

	cluster := gocql.NewCluster(c.Scylla.Host)
	cluster.Keyspace = c.Scylla.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal().Err(err).Msg("failed create sycladb session")
		panic(err)
	}

	return session

}
