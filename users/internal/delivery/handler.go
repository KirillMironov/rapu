package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/domain"
	"github.com/KirillMironov/rapu/internal/delivery/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service domain.UsersService
	proto.UnimplementedUsersServer
}

func NewHandler(usersService domain.UsersService) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterUsersServer(server, &Handler{service: usersService})
	return server
}

func (h *Handler) SignUp(ctx context.Context, request *proto.SignUpRequest) (*proto.Response, error) {
	var user = domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	token, err := h.service.SignUp(user)
	if err != nil {
		return nil, err
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h *Handler) SignIn(ctx context.Context, request *proto.SignInRequest) (*proto.Response, error) {
	var user = domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	token, err := h.service.SignIn(user)
	if err != nil {
		return nil, err
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h *Handler) Authenticate(ctx context.Context, request *proto.AuthRequest) (*proto.AuthResponse, error) {
	userId, err := h.service.Authenticate(request.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &proto.AuthResponse{UserId: userId}, nil
}
