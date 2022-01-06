package delivery

import (
	"context"
	"encoding/json"
	"github.com/KirillMironov/rapu/posts/domain"
	"github.com/KirillMironov/rapu/posts/internal/delivery/proto"
	"github.com/KirillMironov/rapu/posts/pkg/logger"
	"google.golang.org/grpc"
)

type Handler struct {
	service domain.PostsService
	logger  logger.Logger
	proto.UnimplementedPostsServer
}

func NewHandler(service domain.PostsService, logger logger.Logger) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterPostsServer(server, &Handler{
		service: service,
		logger:  logger,
	})
	return server
}

func (h *Handler) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	var post = domain.Post{
		UserId:  request.GetUserId(),
		Message: request.GetMessage(),
	}

	err := h.service.Create(post)
	if err != nil {
		h.logger.Info(err)
	}
	return nil, err
}

func (h *Handler) GetByUserId(ctx context.Context, request *proto.GetByUserIdRequest) (*proto.GetByUserIdResponse, error) {
	posts, err := h.service.GetByUserId(request.GetUserId(), request.GetOffset(), request.GetLimit())
	if err != nil {
		h.logger.Info(err)
		return nil, err
	}

	encoded, err := json.Marshal(posts)
	if err != nil {
		h.logger.Error(err)
		return nil, err
	}

	return &proto.GetByUserIdResponse{Posts: encoded}, nil
}
