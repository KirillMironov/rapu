package service

import (
	"context"
	"github.com/KirillMironov/rapu/posts/internal/domain"
	"time"
)

type Posts struct {
	postsRepository PostsRepository
	maxPostsPerPage int64
}

type PostsRepository interface {
	Create(context.Context, domain.Post) error
	GetByUserId(ctx context.Context, userId, offset string, limit int64) ([]domain.Post, error)
}

func NewPosts(postsRepository PostsRepository, maxPostsPerPage int64) *Posts {
	return &Posts{
		postsRepository: postsRepository,
		maxPostsPerPage: maxPostsPerPage,
	}
}

func (p Posts) Create(ctx context.Context, post domain.Post) error {
	if post.UserId == "" || post.Message == "" {
		return domain.ErrEmptyParameters
	}

	post.CreatedAt = time.Now()

	return p.postsRepository.Create(ctx, post)
}

func (p Posts) GetByUserId(ctx context.Context, userId, offset string, limit int64) ([]domain.Post, error) {
	if userId == "" {
		return nil, domain.ErrEmptyParameters
	}

	if limit < 1 || limit > p.maxPostsPerPage {
		limit = p.maxPostsPerPage
	}

	return p.postsRepository.GetByUserId(ctx, userId, offset, limit)
}
