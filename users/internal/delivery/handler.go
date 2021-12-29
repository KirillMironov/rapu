package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/domain"
	"github.com/KirillMironov/rapu/internal/delivery/proto"
	"github.com/KirillMironov/rapu/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service      domain.UsersService
	tokenManager auth.TokenManager
	proto.UnimplementedUsersServer
}

func NewHandler(usersService domain.UsersService, tokenManager auth.TokenManager) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterUsersServer(server, &Handler{
		service:      usersService,
		tokenManager: tokenManager,
	})
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
	userId, err := h.tokenManager.VerifyAuthToken(request.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &proto.AuthResponse{UserId: userId}, nil
}
