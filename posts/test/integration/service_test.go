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
	"log"
	"testing"
	"time"
)

const (
	mongoImage = "mongo:4.4-rc-focal"
	mongoPort  = "27017"
	dbname     = "rapu"
	collection = "posts"
	userId     = "1"
	maxLimit   = 10
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
	assert.NoError(t, err)

	repo := _mongo.NewPostsRepository(db)
	svc := service.NewPostsService(repo, maxLimit)

	err = svc.Create(post)
	assert.NoError(t, err)

	err = svc.Create(post)
	assert.NoError(t, err)

	// Another UserId
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

	log.Println(posts)

	// Another UserId
	posts, err = svc.GetByUserId("2", "", 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.NotEmpty(t, posts[0])
}

func TestPosts_GetByUserId(t *testing.T) {
	db, terminate, err := mongoSetup()
	defer terminate()
	assert.NoError(t, err)

	repo := _mongo.NewPostsRepository(db)
	svc := service.NewPostsService(repo, maxLimit)

	for i := 0; i < maxLimit; i++ {
		err = svc.Create(post)
		assert.NoError(t, err)
	}

	posts, err := svc.GetByUserId(userId, "", 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, posts)

	// Checking that posts are sorted in descending order by CreatedAt field
	for i := 0; i < len(posts)-1; i++ {
		assert.True(t, posts[i].CreatedAt.After(posts[i+1].CreatedAt) || posts[i].CreatedAt.Equal(posts[i+1].CreatedAt))
	}
}

func TestPosts_GetByUserId_Pagination_offset(t *testing.T) {
	db, terminate, err := mongoSetup()
	defer terminate()
	assert.NoError(t, err)

	repo := _mongo.NewPostsRepository(db)
	svc := service.NewPostsService(repo, maxLimit)

	for i := 0; i < maxLimit; i++ {
		err = svc.Create(post)
		assert.NoError(t, err)
	}

	// Default offset
	posts, err := svc.GetByUserId(userId, "", 0)
	assert.NoError(t, err)
	assert.Len(t, posts, maxLimit)

	// Create another post
	var message = "Next page"
	err = svc.Create(domain.Post{
		UserId:    userId,
		Message:   message,
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	// Manual offset
	postsWithOffset, err := svc.GetByUserId(userId, posts[0].Id.Hex(), 0)
	assert.NoError(t, err)
	assert.Len(t, postsWithOffset, 1)
	assert.Equal(t, postsWithOffset[0].Message, message)
}

func TestPosts_GetByUserId_Pagination_limit(t *testing.T) {
	db, terminate, err := mongoSetup()
	defer terminate()
	assert.NoError(t, err)

	repo := _mongo.NewPostsRepository(db)
	svc := service.NewPostsService(repo, maxLimit)

	for i := 0; i < maxLimit*2; i++ {
		err = svc.Create(post)
		assert.NoError(t, err)
	}

	testCases := []struct {
		limit          int64
		expectedLength int
	}{
		{0, maxLimit},
		{-1, maxLimit},
		{maxLimit * 2, maxLimit},
		{maxLimit / 2, maxLimit / 2},
	}

	for _, tc := range testCases {
		posts, err := svc.GetByUserId(userId, "", tc.limit)
		assert.NoError(t, err)
		assert.Len(t, posts, tc.expectedLength)
	}
}

func mongoSetup() (*mongo.Collection, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        mongoImage,
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

	return client.Database(dbname).Collection(collection), func() {
		container.Terminate(ctx)
	}, nil
}
