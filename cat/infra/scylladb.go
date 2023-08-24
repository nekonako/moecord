package infra

import (
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	scyllaMigration "github.com/golang-migrate/migrate/v4/database/cassandra"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx/v2"
)

func newScylladb(c *config.Config) *gocqlx.Session {

	cluster := gocql.NewCluster(c.Scylla.Host)
	cluster.Keyspace = c.Scylla.Keyspace
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal().Err(err).Msg("failed create sycladb session")
		panic(err)
	}

	runScyllaMigraton(&session, c.Scylla.Keyspace)

	return &session

}

func runScyllaMigraton(session *gocqlx.Session, keySpace string) {
	driver, err := scyllaMigration.WithInstance(session.Session, &scyllaMigration.Config{
		KeyspaceName: keySpace,
	})

	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migration/scylla", keySpace, driver)
	if err != nil {
		log.Fatal().Err(err).Msg("failed create migration instance")
		return
	}

	err = m.Up()
	if err != nil && err.Error() == "no change" {
		log.Info().Msg("no migration")
		return
	}

	if err != nil {
		log.Error().Err(err).Msg("failed up migration")
		log.Info().Msg("try down migration to clean dirty migration")
		if err = m.Down(); err != nil {
			log.Fatal().Err(err).Msg("failed down migration")
			return
		}
	}

	log.Info().Msg("successfully run migration")

}
