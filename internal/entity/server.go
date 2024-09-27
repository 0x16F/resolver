package entity

import (
	"errors"
)

var (
	ErrServerNotFound = errors.New("server not found")
	ErrUserNotFound   = errors.New("user not found")
)

type Server struct {
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	URL      string `json:"url"`
	Port     uint16 `json:"port"`
	UserPort uint16 `json:"user_port"`
	Secret   string `json:"secret"`
}

type CreateServerReq struct {
	IP       string `json:"ip"`
	URL      string `json:"url"`
	UserPort uint16 `json:"user_port"`
	Port     uint16 `json:"port"`
	Secret   string `json:"secret"`
}
