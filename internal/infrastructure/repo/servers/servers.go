package servers

import (
	"context"
	"errors"

	"github.com/0x16f/vpn-resolver/internal/entity"
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

func (r *Repo) GetServer(ctx context.Context, id int) (entity.Server, error) {
	query := `
		SELECT
			id, ip, url, port, secret
		FROM
			vp_servers
		WHERE
			id = @server_id
	`

	args := pgx.NamedArgs{
		"server_id": id,
	}

	server := entity.Server{}

	err := r.db.QueryRow(ctx, query, args).Scan(&server.ID, &server.IP, &server.URL, &server.Port, &server.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Server{}, entity.ErrServerNotFound
		}

		return entity.Server{}, err
	}

	return server, err
}

func (r *Repo) GetServers(ctx context.Context) ([]entity.Server, error) {
	query := `
		SELECT
			id, ip, url, port, secret
		FROM
			vp_servers
		WHERE blocked = @blocked
	`

	args := pgx.NamedArgs{
		"blocked": false,
	}

	var servers []entity.Server

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	server := entity.Server{}

	for rows.Next() {
		err = rows.Scan(&server.ID, &server.IP, &server.URL, &server.Port, &server.Secret)
		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	return servers, err
}

func (r *Repo) CreateServer(ctx context.Context, req entity.CreateServerReq) (entity.Server, error) {
	query := `
		INSERT INTO vp_servers (ip, url, port, secret) VALUES (@ip, @url, @port, @secret) RETURNING id
	`

	args := pgx.NamedArgs{
		"ip":     req.IP,
		"url":    req.URL,
		"port":   req.Port,
		"secret": req.Secret,
	}

	server := entity.Server{}

	err := r.db.QueryRow(ctx, query, args).Scan(&server.ID)
	if err != nil {
		return entity.Server{}, err
	}

	server.IP = req.IP
	server.URL = req.URL
	server.Port = req.Port
	server.Secret = req.Secret

	return server, err
}

func (r *Repo) DeleteServer(ctx context.Context, id int) error {
	query := `
		DELETE FROM vp_servers WHERE id = @server_id
	`

	args := pgx.NamedArgs{
		"server_id": id,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
