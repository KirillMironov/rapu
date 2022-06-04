//go:build integration

package integration

import (
	"context"
	"fmt"
	"github.com/KirillMironov/rapu/posts/internal/delivery"
	"github.com/KirillMironov/rapu/posts/internal/delivery/proto"
	"github.com/KirillMironov/rapu/posts/internal/repository"
	"github.com/KirillMironov/rapu/posts/internal/service"
	"github.com/KirillMironov/rapu/posts/test/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const (
	mongoImage      = "mongo:4.4-rc-focal"
	mongoPort       = "27017"
	mongoDB         = "mongo"
	mongoCollection = "posts"
	postsPerPage    = 10
)

func newClient(t *testing.T) proto.PostsClient {
	t.Helper()

	var (
		db       = newMongo(t)
		handler  = newHandler(t, db)
		listener = bufconn.Listen(1024 * 1024)
		ctx      = context.Background()
	)

	go func() {
		err := handler.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close()
	})

	return proto.NewPostsClient(conn)
}

func newHandler(t *testing.T, db *mongo.Collection) *grpc.Server {
	t.Helper()

	postsRepository := repository.NewPosts(db)
	postsService := service.NewPosts(postsRepository, postsPerPage, mock.Logger{})

	return delivery.NewHandler(postsService)
}

func newMongo(t *testing.T) *mongo.Collection {
	t.Helper()

	var ctx = context.Background()

	containerRequest := testcontainers.ContainerRequest{
		Image:        mongoImage,
		ExposedPorts: []string{mongoPort},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	mappedPort, err := container.MappedPort(ctx, mongoPort)
	require.NoError(t, err)

	connectionString := fmt.Sprintf("mongodb://%s:%s", host, mappedPort.Port())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	require.NoError(t, err)

	err = client.Ping(ctx, readpref.Primary())
	require.NoError(t, err)

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	return client.Database(mongoDB).Collection(mongoCollection)
}
