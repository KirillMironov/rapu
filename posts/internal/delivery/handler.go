package delivery

import (
	"context"
	"encoding/json"
	"github.com/KirillMironov/rapu/posts/internal/delivery/proto"
	"github.com/KirillMironov/rapu/posts/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	postsService PostsService
	proto.UnimplementedPostsServer
}

type PostsService interface {
	Create(context.Context, domain.Post) error
	GetByUserId(ctx context.Context, userId, offset string, limit int64) ([]domain.Post, error)
}

func NewHandler(postsService PostsService) *grpc.Server {
	var server = grpc.NewServer()
	proto.RegisterPostsServer(server, &Handler{
		postsService: postsService,
	})
	return server
}

func (h *Handler) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	var post = domain.Post{
		UserId:  request.GetUserId(),
		Message: request.GetMessage(),
	}

	err := h.postsService.Create(ctx, post)
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.CreateResponse{}, nil
}

func (h *Handler) GetByUserId(ctx context.Context, request *proto.GetByUserIdRequest) (*proto.GetByUserIdResponse, error) {
	posts, err := h.postsService.GetByUserId(ctx, request.GetUserId(), request.GetOffset(), request.GetLimit())
	if err != nil {
		switch err {
		case domain.ErrEmptyParameters:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrEmptyResult:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	encoded, err := json.Marshal(posts)
	if err != nil {
		return nil, err
	}

	return &proto.GetByUserIdResponse{Posts: encoded}, nil
}
