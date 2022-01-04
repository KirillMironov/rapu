package domain

import "time"

type Post struct {
	Id        string    `bson:"_id"`
	UserId    string    `bson:"user_id"`
	Message   string    `bson:"message"`
	CreatedAt time.Time `bson:"created_at"`
}

type PostsRepository interface {
	Create(Post) (string, error)
	GetByUserId(userId string) ([]Post, error)
}
