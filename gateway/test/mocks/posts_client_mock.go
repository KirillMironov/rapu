package mocks

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
)

type PostsClientMock struct{}

func (PostsClientMock) Create(ctx context.Context, in *proto.CreateRequest,
	opts ...grpc.CallOption) (*proto.CreateResponse, error) {
	return nil, nil
}

func (PostsClientMock) GetByUserId(ctx context.Context, in *proto.GetByUserIdRequest,
	opts ...grpc.CallOption) (*proto.GetByUserIdResponse, error) {
	return nil, nil
}
