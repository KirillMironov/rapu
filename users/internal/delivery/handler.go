package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/users/domain"
	"github.com/KirillMironov/rapu/users/internal/delivery/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service domain.UsersService
	logger  Logger
	proto.UnimplementedUsersServer
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

func NewHandler(usersService domain.UsersService, logger Logger) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterUsersServer(server, &Handler{
		service: usersService,
		logger:  logger,
	})
	return server
}

func (h *Handler) SignUp(_ context.Context, request *proto.SignUpRequest) (*proto.Response, error) {
	var user = domain.User{
		Username: request.GetUsername(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	token, err := h.service.SignUp(user)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			h.logger.Info(err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserAlreadyExists:
			h.logger.Info(err)
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			h.logger.Error(err)
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h *Handler) SignIn(_ context.Context, request *proto.SignInRequest) (*proto.Response, error) {
	var user = domain.User{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	token, err := h.service.SignIn(user)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			h.logger.Info(err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserNotFound, domain.ErrWrongPassword:
			h.logger.Info(err)
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			h.logger.Error(err)
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h *Handler) Authenticate(_ context.Context, request *proto.AuthRequest) (*proto.AuthResponse, error) {
	userId, err := h.service.Authenticate(request.GetAccessToken())
	if err != nil {
		h.logger.Info(err)
		if err == domain.ErrEmptyParameters {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &proto.AuthResponse{UserId: userId}, nil
}
