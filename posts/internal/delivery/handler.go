package delivery

import (
	"context"
	"encoding/json"
	"github.com/KirillMironov/rapu/posts/domain"
	"github.com/KirillMironov/rapu/posts/internal/delivery/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service domain.PostsService
	logger  Logger
	proto.UnimplementedPostsServer
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

func NewHandler(service domain.PostsService, logger Logger) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterPostsServer(server, &Handler{
		service: service,
		logger:  logger,
	})
	return server
}

func (h *Handler) Create(_ context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	var post = domain.Post{
		UserId:  request.GetUserId(),
		Message: request.GetMessage(),
	}

	err := h.service.Create(post)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			h.logger.Info(err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			h.logger.Error(err)
			return nil, err
		}
	}

	return &proto.CreateResponse{}, nil
}

func (h *Handler) GetByUserId(_ context.Context, request *proto.GetByUserIdRequest) (*proto.GetByUserIdResponse, error) {
	posts, err := h.service.GetByUserId(request.GetUserId(), request.GetOffset(), request.GetLimit())
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			h.logger.Info(err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrEmptyResult:
			h.logger.Info(err)
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			h.logger.Error(err)
			return nil, err
		}
	}

	encoded, err := json.Marshal(posts)
	if err != nil {
		h.logger.Error(err)
		return nil, err
	}

	return &proto.GetByUserIdResponse{Posts: encoded}, nil
}
