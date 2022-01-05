package service

import (
	"errors"
	"github.com/KirillMironov/rapu/posts/domain"
	"time"
)

var errNotEnoughArgs = errors.New("not enough arguments")

type PostsService struct {
	repository domain.PostsRepository
	maxLimit   int64
}

func NewPostsService(repository domain.PostsRepository, maxLimit int64) *PostsService {
	return &PostsService{repository: repository, maxLimit: maxLimit}
}

func (p *PostsService) Create(post domain.Post) error {
	if post.UserId == "" || post.Message == "" {
		return errNotEnoughArgs
	}

	post.CreatedAt = time.Now()

	return p.repository.Create(post)
}

func (p *PostsService) GetByUserId(userId, offset string, limit int64) ([]domain.Post, error) {
	if userId == "" {
		return nil, errNotEnoughArgs
	}

	if limit < 1 || limit > p.maxLimit {
		limit = p.maxLimit
	}

	return p.repository.GetByUserId(userId, offset, limit)
}
