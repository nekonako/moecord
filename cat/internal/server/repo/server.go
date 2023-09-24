package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type Server struct {
	ID            ulid.ULID `db:"id"`
	OwnerID       ulid.ULID `db:"owner_id"`
	Name          string    `db:"name"`
	DirectMessage bool      `db:"direct_message"`
	Avatar        string    `db:"avatar"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type ServerMember struct {
	ID        ulid.ULID `db:"id"`
	ServerID  ulid.ULID `db:"server_id"`
	Username  string    `db:"username"`
	Avatar    string    `db:"avatar"`
	UserID    ulid.ULID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *Repository) SaveServer(ctx context.Context, tx *sqlx.Tx, server Server) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveServer")
	defer tracer.Finish(span)

	query := `
		INSERT INTO servers (
			id,
			owner_id,
			name,
			direct_message,
            avatar,
			created_at,
			updated_at
    ) VALUES (:id, :owner_id, :name, :direct_message, :avatar, :created_at, :updated_at)
	`

	_, err := tx.NamedExecContext(ctx, query, server)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert user")
		return err
	}

	return nil

}

func (r *Repository) SaveServerMember(ctx context.Context, tx *sqlx.Tx, member ServerMember) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveServer")
	defer tracer.Finish(span)

	query := `
		INSERT INTO server_member (
			id,
			server_id,
			user_id,
			created_at
		) VALUES (:id, :server_id, :user_id, :created_at)
	`

	_, err := tx.NamedExecContext(ctx, query, member)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert server member")
		return err
	}

	return nil
}

func (r *Repository) InviteServerMember(ctx context.Context, tx *sqlx.Tx, member ServerMember) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveServer")
	defer tracer.Finish(span)

	query := `
		INSERT INTO server_member (
			id,
			server_id,
			user_id,
			created_at
		) VALUES (:id, :server_id, :user_id, :created_at)
	`

	_, err := tx.NamedExecContext(ctx, query, member)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert server member")
		return err
	}

	return nil
}

func (r *Repository) ListServerUser(ctx context.Context, userID ulid.ULID) ([]Server, error) {

	span := tracer.SpanFromContext(ctx, "repo.ListServer")
	defer tracer.Finish(span)

	query := `
	SELECT
		s.id,
		s.owner_id,
		s.name,
		s.direct_message,
        s.avatar,
		s.created_at,
		s.updated_at
	FROM server_member AS sm
	INNER JOIN servers AS s ON s.id = sm.server_id
	WHERE sm.user_id = $1 ORDER BY sm.id
	`

	result := []Server{}
	err := r.postgres.SelectContext(ctx, &result, query, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed get server")
		return result, err
	}

	return result, nil

}

func (r *Repository) ListServerMember(ctx context.Context, serverID ulid.ULID) ([]ServerMember, error) {
	span := tracer.SpanFromContext(ctx, "repo.ListServer")
	defer tracer.Finish(span)

	query := `
	SELECT
		sm.id,
		sm.user_id,
		sm.server_id,
        u.username,
        u.avatar,
		sm.created_at
	FROM server_member AS sm
    INNER JOIN users As u ON u.id = sm.user_id
    WHERE sm.server_id = $1
    `

	result := []ServerMember{}
	err := r.postgres.SelectContext(ctx, &result, query, serverID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed get server")
		return result, err
	}

	return result, nil

}

func (r *Repository) UpdateServer(ctx context.Context, s Server) error {

	span := tracer.SpanFromContext(ctx, "repo.SaveServer")
	defer tracer.Finish(span)

	query := `UPDATE servers SET name=:name, avatar=:avatar WHERE id = :id`
	_, err := r.postgres.NamedExecContext(ctx, query, s)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert server member")
		return err
	}

	return nil
}

func (r *Repository) GetServerByID(ctx context.Context, serverID ulid.ULID) (Server, error) {

	span := tracer.SpanFromContext(ctx, "repo.SaveServer")
	defer tracer.Finish(span)

	result := Server{}
	query := `SELECT id, name, owner_id, direct_message, avatar, created_at, updated_at FROM servers WHERE id = $1`
	err := r.postgres.GetContext(ctx, &result, query, serverID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Err(err).Msg("failed insert server member")
		return result, err
	}

	return result, nil
}
