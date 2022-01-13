package mocks

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
)

type UsersClientMock struct{}

func (UsersClientMock) SignUp(context.Context, *proto.SignUpRequest,
	...grpc.CallOption) (*proto.Response, error) {
	return nil, nil
}

func (UsersClientMock) SignIn(context.Context, *proto.SignInRequest,
	...grpc.CallOption) (*proto.Response, error) {
	return nil, nil
}

func (UsersClientMock) Authenticate(context.Context, *proto.AuthRequest,
	...grpc.CallOption) (*proto.AuthResponse, error) {
	return &proto.AuthResponse{UserId: "1"}, nil
}
