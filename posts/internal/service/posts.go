package service

import (
	"context"
	"github.com/KirillMironov/rapu/posts/internal/domain"
	"time"
)

type Posts struct {
	postsRepository PostsRepository
	maxPostsPerPage int64
	logger          Logger
}

type PostsRepository interface {
	Create(context.Context, domain.Post) error
	GetByUserId(ctx context.Context, userId, offset string, limit int64) ([]domain.Post, error)
}

type Logger interface {
	Error(args ...interface{})
}

func NewPosts(postsRepository PostsRepository, maxPostsPerPage int64, logger Logger) *Posts {
	return &Posts{
		postsRepository: postsRepository,
		maxPostsPerPage: maxPostsPerPage,
		logger:          logger,
	}
}

func (p *Posts) Create(ctx context.Context, post domain.Post) error {
	if post.UserId == "" || post.Message == "" {
		return domain.ErrEmptyParameters
	}

	post.CreatedAt = time.Now()

	err := p.postsRepository.Create(ctx, post)
	if err != nil {
		p.logger.Error(err)
	}
	return err
}

func (p *Posts) GetByUserId(ctx context.Context, userId, offset string, limit int64) ([]domain.Post, error) {
	if userId == "" {
		return nil, domain.ErrEmptyParameters
	}

	if limit < 1 || limit > p.maxPostsPerPage {
		limit = p.maxPostsPerPage
	}

	posts, err := p.postsRepository.GetByUserId(ctx, userId, offset, limit)
	if err != nil {
		p.logger.Error(err)
		return nil, err
	}

	return posts, nil
}
