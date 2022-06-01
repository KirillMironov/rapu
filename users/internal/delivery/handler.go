package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/users/internal/delivery/proto"
	"github.com/KirillMironov/rapu/users/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	usersService UsersService
	proto.UnimplementedUsersServer
}

type UsersService interface {
	SignUp(domain.User) (token string, err error)
	SignIn(domain.User) (token string, err error)
	Authenticate(token string) (userId string, err error)
	UserExists(userId string) (bool, error)
}

func NewHandler(usersService UsersService) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterUsersServer(server, &Handler{
		usersService: usersService,
	})
	return server
}

func (h *Handler) SignUp(_ context.Context, request *proto.SignUpRequest) (*proto.Response, error) {
	var user = domain.User{
		Username: request.GetUsername(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	token, err := h.usersService.SignUp(user)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
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

	token, err := h.usersService.SignIn(user)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserNotFound, domain.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h *Handler) Authenticate(_ context.Context, request *proto.AuthRequest) (*proto.AuthResponse, error) {
	userId, err := h.usersService.Authenticate(request.GetAccessToken())
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
	}

	return &proto.AuthResponse{UserId: userId}, nil
}

func (h *Handler) UserExists(_ context.Context, request *proto.UserExistsRequest) (*proto.UserExistsResponse, error) {
	exists, err := h.usersService.UserExists(request.GetUserId())
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.UserExistsResponse{Exists: exists}, nil
}
