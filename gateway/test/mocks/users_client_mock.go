package mocks

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
)

type UsersClientMock struct{}

func (UsersClientMock) SignUp(ctx context.Context, in *proto.SignUpRequest,
	opts ...grpc.CallOption) (*proto.Response, error) {
	return nil, nil
}

func (UsersClientMock) SignIn(ctx context.Context, in *proto.SignInRequest,
	opts ...grpc.CallOption) (*proto.Response, error) {
	return nil, nil
}

func (UsersClientMock) Authenticate(ctx context.Context, in *proto.AuthRequest,
	opts ...grpc.CallOption) (*proto.AuthResponse, error) {
	return nil, nil
}
