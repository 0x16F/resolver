package servers

import (
	"context"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/0x16f/vpn-resolver/pkg/codes"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type serversSrv interface {
	CreateServer(ctx context.Context, req entity.CreateServerReq) (entity.Server, error)
	GetServer(ctx context.Context, id int) (entity.Server, error)
	GetServers(ctx context.Context) ([]entity.Server, error)
	DeleteServer(ctx context.Context, id int) error
}

type errorSrv interface {
	GetError(code int) error
}

type Handler struct {
	serversSrv serversSrv
	errorsSrv  errorSrv
}

func New(serversSrv serversSrv, errorsSrv errorSrv) *Handler {
	return &Handler{
		serversSrv: serversSrv,
		errorsSrv:  errorsSrv,
	}
}

func (h *Handler) CreateServer(c *fiber.Ctx) error {
	var req entity.CreateServerReq

	err := c.BodyParser(&req)
	if err != nil {
		logrus.Errorf("failed to parse body: %v", err)

		return h.errorsSrv.GetError(codes.IncorrectBody)
	}

	var server entity.Server

	server, err = h.serversSrv.CreateServer(c.Context(), req)
	if err != nil {
		logrus.Errorf("failed to create server: %v", err)

		return err
	}

	return c.Status(fiber.StatusCreated).JSON(server)
}

func (h *Handler) GetServer(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logrus.Errorf("failed to parse id: %v", err)

		return h.errorsSrv.GetError(codes.IncorrectID)
	}

	var server entity.Server

	server, err = h.serversSrv.GetServer(c.Context(), id)
	if err != nil {
		logrus.Errorf("failed to get server: %v", err)

		return err
	}

	return c.JSON(server)
}

func (h *Handler) GetServers(c *fiber.Ctx) error {
	servers, err := h.serversSrv.GetServers(c.Context())
	if err != nil {
		logrus.Errorf("failed to get servers: %v", err)

		return err
	}

	return c.JSON(servers)
}

func (h *Handler) DeleteServer(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		logrus.Errorf("failed to parse id: %v", err)

		return h.errorsSrv.GetError(codes.IncorrectID)
	}

	err = h.serversSrv.DeleteServer(c.Context(), id)
	if err != nil {
		logrus.Errorf("failed to delete server: %v", err)

		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
