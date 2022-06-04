package mock

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
)

type UsersClient struct{}

func (UsersClient) SignUp(context.Context, *proto.SignUpRequest,
	...grpc.CallOption) (*proto.Response, error) {
	return nil, nil
}

func (UsersClient) SignIn(context.Context, *proto.SignInRequest,
	...grpc.CallOption) (*proto.Response, error) {
	return nil, nil
}

func (UsersClient) Authenticate(context.Context, *proto.AuthRequest,
	...grpc.CallOption) (*proto.AuthResponse, error) {
	return &proto.AuthResponse{UserId: "1"}, nil
}

func (UsersClient) UserExists(context.Context, *proto.UserExistsRequest,
	...grpc.CallOption) (*proto.UserExistsResponse, error) {
	return &proto.UserExistsResponse{Exists: true}, nil
}
