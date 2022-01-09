package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Message   string             `bson:"message" json:"message"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type PostsService interface {
	Create(Post) error
	GetByUserId(userId, offset string, limit int64) ([]Post, error)
}

type PostsRepository interface {
	Create(Post) error
	GetByUserId(userId, offset string, limit int64) ([]Post, error)
}
