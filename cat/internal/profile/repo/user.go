package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID        ulid.ULID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Avatar    string    `db:"avatar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *Repository) SaveOrUpdateUser(ctx context.Context, tx *sqlx.Tx, user User) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveOrUpdateUser")
	defer tracer.Finish(span)

	query := `
		INSERT INTO users (
			id,
			username,
			email,
			created_at,
			updated_at
		) VALUES (:id, :username, :email, :created_at, :updated_at)
		ON CONFLICT (email) DO NOTHING
	`

	_, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert user")
		return err
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, user User) error {
	span := tracer.SpanFromContext(ctx, "repo.SaveOrUpdateUser")
	defer tracer.Finish(span)

	query := `
    UPDATE users SET
        id=:id,
        username=:username,
        email=:email,
        avatar=:avatar
    WHERE id = :id
	`

	_, err := r.postgres.NamedExecContext(ctx, query, user)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert user")
		return err
	}

	return nil

}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {

	span := tracer.SpanFromContext(ctx, "repo.GetUserByEmail")
	defer tracer.Finish(span)

	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE email = $1"

	result := User{}
	err := r.postgres.GetContext(ctx, &result, query, email)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed get user")
		return result, err
	}

	return result, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id ulid.ULID) (User, error) {

	span := tracer.SpanFromContext(ctx, "repo.GetUserByEmail")
	defer tracer.Finish(span)

	query := "SELECT id, username, email, avatar, created_at, updated_at FROM users WHERE id = $1"

	result := User{}
	err := r.postgres.GetContext(ctx, &result, query, id)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed get user")
		return result, err
	}

	return result, nil
}
