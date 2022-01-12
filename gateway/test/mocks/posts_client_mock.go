package mocks

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
)

type PostsClientMock struct{}

func (PostsClientMock) Create(context.Context, *proto.CreateRequest,
	...grpc.CallOption) (*proto.CreateResponse, error) {
	return nil, nil
}

func (PostsClientMock) GetByUserId(context.Context, *proto.GetByUserIdRequest,
	...grpc.CallOption) (*proto.GetByUserIdResponse, error) {
	return nil, nil
}
