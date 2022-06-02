package repository

import (
	"context"
	"github.com/KirillMironov/rapu/posts/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Posts struct {
	db *mongo.Collection
}

func NewPosts(db *mongo.Collection) *Posts {
	return &Posts{db: db}
}

func (p *Posts) Create(ctx context.Context, post domain.Post) error {
	_, err := p.db.InsertOne(ctx, post)
	return err
}

func (p *Posts) GetByUserId(ctx context.Context, userId, offset string, limit int64) ([]domain.Post, error) {
	id, err := primitive.ObjectIDFromHex(offset)
	if err != nil && err != primitive.ErrInvalidHex {
		return nil, err
	}

	var query = bson.M{"user_id": userId}
	if id != primitive.NilObjectID {
		query = bson.M{"user_id": userId, "_id": bson.M{"$lt": id}}
	}

	var opts = options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(limit)

	cur, err := p.db.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var posts []domain.Post
	err = cur.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, domain.ErrEmptyResult
	}

	return posts, nil
}
