package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/users/domain"
	"github.com/KirillMironov/rapu/users/internal/delivery/proto"
	"github.com/KirillMironov/rapu/users/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const notEnoughArgs = "not enough arguments"

type Handler struct {
	service domain.UsersService
	logger  logger.Logger
	proto.UnimplementedUsersServer
}

func NewHandler(usersService domain.UsersService, logger logger.Logger) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterUsersServer(server, &Handler{
		service: usersService,
		logger:  logger,
	})
	return server
}

func (h *Handler) SignUp(ctx context.Context, request *proto.SignUpRequest) (*proto.Response, error) {
	if request.GetUsername() == "" || request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, notEnoughArgs)
	}

	var user = domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	token, err := h.service.SignUp(user)
	if err != nil {
		h.logger.Info(err)
		return nil, err
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h *Handler) SignIn(ctx context.Context, request *proto.SignInRequest) (*proto.Response, error) {
	if request.GetEmail() == "" || request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, notEnoughArgs)
	}

	var user = domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	token, err := h.service.SignIn(user)
	if err != nil {
		h.logger.Info(err)
		return nil, err
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h *Handler) Authenticate(ctx context.Context, request *proto.AuthRequest) (*proto.AuthResponse, error) {
	if request.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, notEnoughArgs)
	}

	userId, err := h.service.Authenticate(request.AccessToken)
	if err != nil {
		h.logger.Info(err)
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &proto.AuthResponse{UserId: userId}, nil
}
