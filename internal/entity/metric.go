package entity

import "github.com/google/uuid"

type Metric struct {
	ID       int       `json:"id"`
	ServerID int       `json:"server_id"`
	UserID   uuid.UUID `json:"user_id"`
}

type MetricCreateReq struct {
	ServerID int       `json:"server_id"`
	UserID   uuid.UUID `json:"user_id"`
}
