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
	post.Id = primitive.NewObjectID()

	_, err := p.db.InsertOne(ctx, post)
	return err
}

func (p *PostsRepository) GetByUserId(userId, offset string, limit int64) ([]domain.Post, error) {
	var posts []domain.Post

	id, err := primitive.ObjectIDFromHex(offset)
	if err != nil && err != primitive.ErrInvalidHex {
		return nil, err
	}

	var query = bson.M{"user_id": userId, "_id": bson.M{"$gt": id}}
	var opts = options.Find().
		SetSort(bson.D{{"created_at", -1}}).
		SetLimit(limit)

	cur, err := p.db.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	return posts, cur.All(ctx, &posts)
}
