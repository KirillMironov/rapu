package mock

import (
	"context"
	"github.com/KirillMironov/rapu/gateway/internal/delivery/proto"
	"google.golang.org/grpc"
)

type PostsClient struct{}

func (PostsClient) Create(context.Context, *proto.CreateRequest,
	...grpc.CallOption) (*proto.CreateResponse, error) {
	return nil, nil
}

func (PostsClient) GetByUserId(context.Context, *proto.GetByUserIdRequest,
	...grpc.CallOption) (*proto.GetByUserIdResponse, error) {
	return nil, nil
}
