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

func (p *PostsRepository) Create(post domain.Post) (string, error) {
	post.Id = primitive.NewObjectID()

	res, err := p.db.InsertOne(ctx, post)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (p *PostsRepository) GetByUserId(userId string) ([]domain.Post, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cur, err := p.db.Find(ctx, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}

	var posts []domain.Post
	err = cur.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
