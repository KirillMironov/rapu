//go:build integration

package integration

import (
	"context"
	"fmt"
	"github.com/KirillMironov/rapu/posts/internal/delivery"
	"github.com/KirillMironov/rapu/posts/internal/delivery/proto"
	_mongo "github.com/KirillMironov/rapu/posts/internal/repository/mongo"
	"github.com/KirillMironov/rapu/posts/internal/service"
	"github.com/KirillMironov/rapu/posts/test/mocks"
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
	mongoDB         = "rapu"
	mongoCollection = "posts"
	maxLimit        = 10
)

func newClient(t *testing.T) (proto.PostsClient, func()) {
	db, terminate := mongoSetup(t)
	handler := handlerSetup(db)

	var listener = bufconn.Listen(1024 * 1024)
	go func() {
		log.Fatal(handler.Serve(listener))
	}()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}))
	require.NoError(t, err)

	return proto.NewPostsClient(conn), func() {
		conn.Close()
		terminate()
	}
}

func handlerSetup(db *mongo.Collection) *grpc.Server {
	repo := _mongo.NewPostsRepository(db)
	svc := service.NewPostsService(repo, maxLimit)
	return delivery.NewHandler(svc, &mocks.LoggerMock{})
}

func mongoSetup(t *testing.T) (*mongo.Collection, func()) {
	request := testcontainers.ContainerRequest{
		Image:        mongoImage,
		ExposedPorts: []string{mongoPort},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	require.NoError(t, err)

	ip, err := container.Host(ctx)
	require.NoError(t, err)
	mappedPort, err := container.MappedPort(ctx, mongoPort)
	require.NoError(t, err)

	conn := fmt.Sprintf("mongodb://%s:%s", ip, mappedPort.Port())

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	require.NoError(t, err)

	err = client.Ping(ctx, readpref.Primary())
	require.NoError(t, err)

	return client.Database(mongoDB).Collection(mongoCollection), func() {
		container.Terminate(ctx)
	}
}
