package mocks

import "github.com/KirillMironov/rapu/posts/domain"

type PostsRepositoryMock struct{}

func (PostsRepositoryMock) Create(post domain.Post) error {
	return nil
}

func (PostsRepositoryMock) GetByUserId(userId, offset string, limit int64) ([]domain.Post, error) {
	return nil, nil
}
