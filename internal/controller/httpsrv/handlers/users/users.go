package users

import (
	"context"
	"fmt"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/0x16f/vpn-resolver/pkg/codes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type errorsSrv interface {
	GetError(code int) error
}

type usersSrv interface {
	GetUser(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetUsers(ctx context.Context) ([]entity.User, error)
	CreateUser(ctx context.Context, req entity.CreateUserReq) ([]entity.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type parserSrv interface {
	CreateConfig(userID uuid.UUID) (string, error)
}

type Handler struct {
	url       string
	usersSrv  usersSrv
	errorsSrv errorsSrv
	parserSrv parserSrv
}

func New(uri string, usersSrv usersSrv, errorsSrv errorsSrv, parserSrv parserSrv) *Handler {
	return &Handler{
		url:       uri,
		usersSrv:  usersSrv,
		errorsSrv: errorsSrv,
		parserSrv: parserSrv,
	}
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req entity.CreateUserReq

	err := c.BodyParser(&req)
	if err != nil {
		logrus.Errorf("failed to parse body: %v", err)

		return h.errorsSrv.GetError(codes.IncorrectBody)
	}

	users, err := h.usersSrv.CreateUser(c.Context(), req)
	if err != nil {
		logrus.Errorf("failed to create user: %v", err)

		return err
	}

	// TODO: do it normal

	if len(users) == 0 {
		return c.Status(fiber.StatusNoContent).SendString("nil")
	}

	cfgline, err := h.parserSrv.CreateConfig(users[0].ID)
	if err != nil {
		logrus.Errorf("failed to create config: %v", err)

		return h.errorsSrv.GetError(codes.ServerError)
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("ssconf://%s/%s", h.url, cfgline))
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logrus.Errorf("failed to parse id: %v", err)

		return h.errorsSrv.GetError(codes.IncorrectID)
	}

	var user entity.User

	user, err = h.usersSrv.GetUser(c.Context(), id)
	if err != nil {
		logrus.Errorf("failed to get user: %v", err)

		return err
	}

	return c.JSON(user)
}

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	users, err := h.usersSrv.GetUsers(c.Context())
	if err != nil {
		logrus.Errorf("failed to get users: %v", err)

		return err
	}

	return c.JSON(users)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logrus.Errorf("failed to parse id: %v", err)

		return h.errorsSrv.GetError(codes.IncorrectID)
	}

	err = h.usersSrv.DeleteUser(c.Context(), id)
	if err != nil {
		logrus.Errorf("failed to delete user: %v", err)

		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
