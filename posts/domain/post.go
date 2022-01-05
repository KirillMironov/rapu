package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    string             `bson:"user_id"`
	Message   string             `bson:"message"`
	CreatedAt time.Time          `bson:"created_at"`
}

type PostsRepository interface {
	Create(Post) (string, error)
	GetByUserId(userId string) ([]Post, error)
}
