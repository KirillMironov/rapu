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
	logger       Logger
	proto.UnimplementedUsersServer
}

type UsersService interface {
	SignUp(context.Context, domain.User) (token string, err error)
	SignIn(context.Context, domain.User) (token string, err error)
	Authenticate(token string) (userId string, err error)
	UserExists(ctx context.Context, userId string) (bool, error)
}

type Logger interface {
	Error(args ...interface{})
}

func NewHandler(usersService UsersService, logger Logger) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterUsersServer(server, &Handler{
		usersService: usersService,
		logger:       logger,
	})
	return server
}

func (h Handler) SignUp(ctx context.Context, request *proto.SignUpRequest) (*proto.Response, error) {
	var user = domain.User{
		Username: request.GetUsername(),
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	token, err := h.usersService.SignUp(ctx, user)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			h.logger.Error(err)
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h Handler) SignIn(ctx context.Context, request *proto.SignInRequest) (*proto.Response, error) {
	var user = domain.User{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	token, err := h.usersService.SignIn(ctx, user)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserNotFound, domain.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			h.logger.Error(err)
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.Response{AccessToken: token}, nil
}

func (h Handler) Authenticate(_ context.Context, request *proto.AuthRequest) (*proto.AuthResponse, error) {
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

func (h Handler) UserExists(ctx context.Context, request *proto.UserExistsRequest) (*proto.UserExistsResponse, error) {
	exists, err := h.usersService.UserExists(ctx, request.GetUserId())
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			h.logger.Error(err)
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.UserExistsResponse{Exists: exists}, nil
}
