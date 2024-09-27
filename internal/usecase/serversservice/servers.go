package serversservice

import (
	"context"
	"errors"
	"net"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/0x16f/vpn-resolver/pkg/codes"
	"github.com/sirupsen/logrus"
)

type serversRepo interface {
	CreateServer(ctx context.Context, req entity.CreateServerReq) (entity.Server, error)
	GetServer(ctx context.Context, id int) (entity.Server, error)
	GetServers(ctx context.Context) ([]entity.Server, error)
	DeleteServer(ctx context.Context, id int) error
}

type errorService interface {
	GetError(code int) error
}

type Service struct {
	serversRepo  serversRepo
	errorService errorService
}

func New(serverRepo serversRepo, errorService errorService) *Service {
	return &Service{
		serversRepo:  serverRepo,
		errorService: errorService,
	}
}

func (s *Service) CreateServer(ctx context.Context, req entity.CreateServerReq) (entity.Server, error) {
	ip := net.ParseIP(req.IP)
	if ip == nil {
		return entity.Server{}, s.errorService.GetError(codes.ServerInvalidIP)
	}

	if req.URL == "" {
		return entity.Server{}, s.errorService.GetError(codes.ServerEmptyURL)
	}

	if req.Port <= 0 || req.Port > 65535 {
		return entity.Server{}, s.errorService.GetError(codes.ServerEmptyPort)
	}

	if req.Secret == "" {
		return entity.Server{}, s.errorService.GetError(codes.ServerEmptySecret)
	}

	if req.UserPort <= 0 || req.UserPort > 65535 {
		return entity.Server{}, s.errorService.GetError(codes.IncorrectUserPort)

	}

	server, err := s.serversRepo.CreateServer(ctx, req)
	if err != nil {
		logrus.Errorf("failed to create server: %v", err)

		return entity.Server{}, s.errorService.GetError(codes.ServerError)
	}

	return server, nil
}

func (s *Service) GetServer(ctx context.Context, id int) (entity.Server, error) {
	server, err := s.serversRepo.GetServer(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrServerNotFound) {
			return entity.Server{}, s.errorService.GetError(codes.ServerNotFound)
		}

		logrus.Errorf("failed to get server: %v", err)

		return entity.Server{}, s.errorService.GetError(codes.ServerError)
	}

	return server, nil
}

func (s *Service) GetServers(ctx context.Context) ([]entity.Server, error) {
	servers, err := s.serversRepo.GetServers(ctx)
	if err != nil {
		logrus.Errorf("failed to get servers: %v", err)

		return nil, s.errorService.GetError(codes.ServerError)
	}

	return servers, nil
}

func (s *Service) DeleteServer(ctx context.Context, id int) error {
	_, err := s.GetServer(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrServerNotFound) {
			return s.errorService.GetError(codes.ServerNotFound)
		}

		logrus.Errorf("failed to get server: %v", err)

		return s.errorService.GetError(codes.ServerError)
	}

	err = s.serversRepo.DeleteServer(ctx, id)
	if err != nil {
		logrus.Errorf("failed to delete server: %v", err)

		return s.errorService.GetError(codes.ServerError)
	}

	return nil
}
