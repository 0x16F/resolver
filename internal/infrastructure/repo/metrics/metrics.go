package metrics

import (
	"context"

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

func (r *Repo) CreateMetric(ctx context.Context, req entity.MetricCreateReq) error {
	query := `
		INSERT INTO
			vp_usage_metrics (user_id, server_id, ip, agent)
		VALUES
			(@user_id, @server_id, @ip, @agent)
	`

	args := pgx.NamedArgs{
		"user_id":   req.UserID,
		"server_id": req.ServerID,
		"ip":        req.IP,
		"agent":     req.Agent,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
