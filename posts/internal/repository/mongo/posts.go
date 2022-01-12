package mongo

import (
	"context"
	"github.com/KirillMironov/rapu/posts/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type PostsRepository struct {
	db *mongo.Collection
}

func NewPostsRepository(db *mongo.Collection) *PostsRepository {
	return &PostsRepository{db: db}
}

func (p *PostsRepository) Create(post domain.Post) error {
	_, err := p.db.InsertOne(ctx, post)
	return err
}

func (p *PostsRepository) GetByUserId(userId, offset string, limit int64) ([]domain.Post, error) {
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
	return posts, cur.All(ctx, &posts)
}
