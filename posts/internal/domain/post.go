package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    string             `json:"user_id" bson:"user_id"`
	Message   string             `json:"message" bson:"message"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
