package outline

import (
	"context"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/0x16f/vpn-resolver/pkg/codes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type serversSrv interface {
	GetServer(ctx context.Context, id int) (entity.Server, error)
}

type usersSrv interface {
	GetUser(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type parserSrv interface {
	ParseConfig(cfgLine string) (uuid.UUID, error)
}

type errorsSrv interface {
	GetError(code int) error
}

type Handler struct {
	serversSrv serversSrv
	usersSrv   usersSrv
	parserSrv  parserSrv
	errorsSrv  errorsSrv
}

func New(serversSrv serversSrv, usersSrv usersSrv, parserSrv parserSrv, errorsSrv errorsSrv) *Handler {
	return &Handler{
		serversSrv: serversSrv,
		usersSrv:   usersSrv,
		parserSrv:  parserSrv,
		errorsSrv:  errorsSrv,
	}
}

func (h *Handler) GetConfig(c *fiber.Ctx) error {
	userID, err := h.parserSrv.ParseConfig(c.Params("id"))
	if err != nil {
		logrus.Errorf("failed to parse config: %v", err)

		return h.errorsSrv.GetError(codes.IncorrectCfgLine)
	}

	user, err := h.usersSrv.GetUser(c.Context(), userID)
	if err != nil {
		logrus.Errorf("failed to get user: %v", err)

		return h.errorsSrv.GetError(codes.ServerError)
	}

	server, err := h.serversSrv.GetServer(c.Context(), user.ServerID)
	if err != nil {
		logrus.Errorf("failed to get server: %v", err)

		return h.errorsSrv.GetError(codes.ServerError)
	}

	return c.Status(fiber.StatusOK).JSON(entity.SSConf{
		Server:   server.IP,
		Port:     server.UserPort,
		Password: user.Password,
		Method:   "chacha20-ietf-poly1305",
	})
}
