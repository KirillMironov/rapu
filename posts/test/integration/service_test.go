//go:build integration

package integration

import (
	"context"
	"fmt"
	"github.com/KirillMironov/rapu/posts/domain"
	_mongo "github.com/KirillMironov/rapu/posts/internal/repository/mongo"
	"github.com/KirillMironov/rapu/posts/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

const (
	mongoPort  = "27017"
	database   = "rapu"
	collection = "posts"
	userId     = "1"
)

var (
	ctx  = context.Background()
	post = domain.Post{
		UserId:    userId,
		Message:   "Hello",
		CreatedAt: time.Now(),
	}
)

func TestPosts_Create(t *testing.T) {
	db, terminate, err := mongoSetup()
	defer terminate()
	if err != nil {
		t.Fatal(err)
	}

	repo := _mongo.NewPostsRepository(db)
	svc := service.NewPostsService(repo, 10)

	err = svc.Create(post)
	assert.NoError(t, err)

	err = svc.Create(post)
	assert.NoError(t, err)

	err = svc.Create(domain.Post{
		UserId:    "2",
		Message:   "Empty",
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	posts, err := svc.GetByUserId(userId, "", 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts))
	assert.NotEmpty(t, posts[0])
	assert.NotEmpty(t, posts[1])
}

func TestPosts_GetByUserId(t *testing.T) {
	db, terminate, err := mongoSetup()
	defer terminate()
	if err != nil {
		t.Fatal(err)
	}

	repo := _mongo.NewPostsRepository(db)
	svc := service.NewPostsService(repo, 10)

	err = svc.Create(post)
	assert.NoError(t, err)

	err = svc.Create(domain.Post{
		UserId:    userId,
		Message:   "New post",
		CreatedAt: time.Now().Add(time.Hour),
	})
	assert.NoError(t, err)

	posts, err := svc.GetByUserId(userId, "", 0)
	assert.NoError(t, err)
	assert.True(t, posts[0].CreatedAt.After(posts[1].CreatedAt))
}

func mongoSetup() (*mongo.Collection, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.4-rc-focal",
		ExposedPorts: []string{mongoPort},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	mappedPort, err := container.MappedPort(ctx, mongoPort)
	if err != nil {
		return nil, nil, err
	}

	conn := fmt.Sprintf("mongodb://%s:%s", ip, mappedPort.Port())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	return client.Database(database).Collection(collection), func() {
		container.Terminate(ctx)
	}, nil
}
