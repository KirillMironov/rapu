package mock

import (
	"context"
	"github.com/KirillMironov/rapu/posts/internal/domain"
)

type PostsRepository struct{}

func (PostsRepository) Create(context.Context, domain.Post) error {
	return nil
}

func (PostsRepository) GetByUserId(context.Context, string, string, int64) ([]domain.Post, error) {
	return nil, nil
}
