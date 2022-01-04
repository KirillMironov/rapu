package mongo

import (
	"github.com/KirillMironov/rapu/posts/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostsRepository struct {
	db *mongo.Collection
}

func NewPostsRepository(db *mongo.Collection) *PostsRepository {
	return &PostsRepository{db: db}
}

func (p *PostsRepository) Create(post domain.Post) (string, error) {
	res, err := p.db.InsertOne(nil, post)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), nil
}

func (p *PostsRepository) GetByUserId(userId string) ([]domain.Post, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cur, err := p.db.Find(nil, bson.M{"user_id": userId}, opts)
	if err != nil {
		return nil, err
	}

	var posts []domain.Post
	err = cur.All(nil, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
