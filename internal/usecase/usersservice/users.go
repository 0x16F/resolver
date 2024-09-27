package usersservice

import (
	"context"
	"errors"

	"github.com/0x16f/vpn-resolver/internal/entity"
	"github.com/0x16f/vpn-resolver/pkg/codes"
	"github.com/0x16f/vpn-resolver/pkg/generator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type errorSrv interface {
	GetError(code int) error
}

type usersRepo interface {
	GetUser(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetUsers(ctx context.Context) ([]entity.User, error)
	CreateUser(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type outlineSrv interface {
	CreateUser(ctx context.Context, req entity.OutlineCreateUserReq) (entity.OutlineCreateUserResp, error)
	DeleteUser(ctx context.Context, req entity.OutlineDeleteUserReq) error
}

type serversSrv interface {
	GetServers(ctx context.Context) ([]entity.Server, error)
	GetServer(ctx context.Context, id int) (entity.Server, error)
}

type Service struct {
	usersRepo  usersRepo
	errorSrv   errorSrv
	serversSrv serversSrv
	outlineSrv outlineSrv
}

func New(usersRepo usersRepo, errorSrv errorSrv, serversSrv serversSrv, outlineSrv outlineSrv) *Service {
	return &Service{
		usersRepo:  usersRepo,
		errorSrv:   errorSrv,
		serversSrv: serversSrv,
		outlineSrv: outlineSrv,
	}
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (entity.User, error) {
	user, err := s.usersRepo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return entity.User{}, s.errorSrv.GetError(codes.UserNotFound)
		}

		logrus.Errorf("failed to get user: %v", err)

		return entity.User{}, s.errorSrv.GetError(codes.ServerError)
	}

	return user, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := s.usersRepo.GetUsers(ctx)
	if err != nil {
		logrus.Errorf("failed to get users: %v", err)

		return nil, s.errorSrv.GetError(codes.ServerError)
	}

	return users, nil
}

func (s *Service) CreateUser(ctx context.Context, req entity.CreateUserReq) ([]entity.User, error) {
	if req.Username == "" {
		return nil, s.errorSrv.GetError(codes.UserEmptyUsername)
	}

	servers, err := s.serversSrv.GetServers(ctx)
	if err != nil {
		logrus.Errorf("failed to get servers: %v", err)

		return nil, s.errorSrv.GetError(codes.ServerError)
	}

	users := make([]entity.User, 0, len(servers))

	password, err := generator.GeneratePassword(32)
	if err != nil {
		logrus.Errorf("failed to generate password: %v", err)

		return nil, s.errorSrv.GetError(codes.ServerError)
	}

	for _, server := range servers {
		resp, err := s.outlineSrv.CreateUser(ctx, entity.OutlineCreateUserReq{
			Name:     req.Username,
			Password: password,
			OutlineInfo: entity.OutlineInfo{
				OutlineURL:    server.IP,
				OutlinePort:   server.Port,
				OutlineSecret: server.Secret,
			},
		})

		if err != nil {
			logrus.Errorf("failed to create user: %v", err)

			return nil, s.errorSrv.GetError(codes.ServerError)
		}

		user := entity.NewUser(entity.CreateRepoUserReq{
			Username:  req.Username,
			Password:  password,
			OutlineID: resp.Id,
			ServerID:  server.ID,
		})

		err = s.usersRepo.CreateUser(ctx, user)
		if err != nil {
			logrus.Errorf("failed to create user: %v", err)

			return nil, s.errorSrv.GetError(codes.ServerError)
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.usersRepo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return s.errorSrv.GetError(codes.UserNotFound)
		}

		logrus.Errorf("failed to get user: %v", err)

		return s.errorSrv.GetError(codes.ServerError)
	}

	server, err := s.serversSrv.GetServer(ctx, user.ServerID)
	if err != nil {
		if errors.Is(err, entity.ErrServerNotFound) {
			return s.errorSrv.GetError(codes.ServerNotFound)
		}

		logrus.Errorf("failed to get server: %v", err)

		return s.errorSrv.GetError(codes.ServerError)
	}

	err = s.outlineSrv.DeleteUser(ctx, entity.OutlineDeleteUserReq{
		UserID: user.OutlineID,
		OutlineInfo: entity.OutlineInfo{
			OutlineURL:    server.IP,
			OutlinePort:   server.Port,
			OutlineSecret: server.Secret,
		},
	})

	if err != nil {
		logrus.Errorf("failed to delete user: %v", err)

		return s.errorSrv.GetError(codes.ServerError)
	}

	if err = s.usersRepo.DeleteUser(ctx, id); err != nil {
		logrus.Errorf("failed to delete user: %v", err)

		return s.errorSrv.GetError(codes.ServerError)
	}

	return nil
}
