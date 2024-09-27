package users

import (
	"context"
	"errors"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetUser(ctx context.Context, id uuid.UUID) (entity.User, error) {
	query := `
		SELECT
			id, username, password, outline_id, server_id
		FROM
			vp_users
		WHERE
			id = @user_id
	`

	args := pgx.NamedArgs{
		"user_id": id,
	}

	u := entity.User{}

	err := r.db.QueryRow(ctx, query, args).Scan(&u.ID, &u.Username, &u.Password, &u.OutlineID, &u.ServerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return u, err
}

func (r *Repo) GetUsers(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT
			id, username, password, outline_id, server_id
		FROM
			vp_users
	`

	var users []entity.User

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	u := entity.User{}

	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Username, &u.Password, &u.OutlineID, &u.ServerID)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *Repo) CreateUser(ctx context.Context, req entity.User) error {
	query := `
		INSERT INTO
			vp_users (id, username, password, outline_id, server_id)
		VALUES
			(@id, @username, @password, @outline_id, @server_id)
	`

	args := pgx.NamedArgs{
		"id":         req.ID,
		"username":   req.Username,
		"password":   req.Password,
		"outline_id": req.OutlineID,
		"server_id":  req.ServerID,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM vp_users WHERE id = @user_id
	`

	args := pgx.NamedArgs{
		"user_id": id,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
