package mocks

import "github.com/KirillMironov/rapu/posts/domain"

type PostsRepositoryMock struct{}

func (PostsRepositoryMock) Create(domain.Post) error {
	return nil
}

func (PostsRepositoryMock) GetByUserId(string, string, int64) ([]domain.Post, error) {
	return nil, nil
}
