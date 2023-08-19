package infra

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresMigration "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
)

func newPostgres(c *config.Config) *sqlx.DB {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.Postgres.Username, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Database, "disable")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed create database connection")
		panic(err)
	}

	db.SetMaxOpenConns(c.Postgres.MaxOpenConn)
	db.SetMaxIdleConns(c.Postgres.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(c.Postgres.ConnMaxLifeTime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(c.Postgres.ConnMaxIdleTime) * time.Minute)

	driver, err := postgresMigration.WithInstance(db.DB, &postgresMigration.Config{DatabaseName: c.Postgres.Database})
	if err != nil {
		log.Fatal().Err(err).Msg("failed create database migration instance")
		return db
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration", c.Postgres.Database, driver)
	if err != nil {
		log.Fatal().Err(err).Msg("failed create migration instance")
		return db
	}

	if err = m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Info().Msg("no migration")
			return db
		}

		log.Error().Err(err).Msg("failed up migration")
		log.Info().Msg("try down migration to clean dirty migration")
		if err = m.Down(); err != nil {
			log.Fatal().Err(err).Msg("failed down migration")
			return db
		}

		return db
	}

	log.Info().Msg("successfully run migration")

	return db
}
