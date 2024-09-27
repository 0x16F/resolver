package entity

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	OutlineID string    `json:"outline_id"`
	ServerID  int       `json:"server_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
}

type CreateUserReq struct {
	Username string `json:"username"`
}

type CreateRepoUserReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	OutlineID string `json:"outline_id"`
	ServerID  int    `json:"server_id"`
}

func NewUser(req CreateRepoUserReq) User {
	return User{
		ID:        uuid.New(),
		OutlineID: req.OutlineID,
		ServerID:  req.ServerID,
		Username:  req.Username,
		Password:  req.Password,
	}
}
