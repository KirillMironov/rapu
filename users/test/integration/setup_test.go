//go:build integration

package integration

import (
	"context"
	"fmt"
	"github.com/KirillMironov/rapu/users/internal/delivery"
	"github.com/KirillMironov/rapu/users/internal/delivery/proto"
	"github.com/KirillMironov/rapu/users/internal/repository/postgres"
	"github.com/KirillMironov/rapu/users/internal/service"
	"github.com/KirillMironov/rapu/users/pkg/auth"
	"github.com/KirillMironov/rapu/users/test/mocks"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"
)

const (
	postgresImage      = "postgres:12.7-alpine3.14"
	postgresSchemaPath = "../../config/schema.sql"
	jwtKey             = "secretKey"
	tokenTTL           = time.Minute
)

func newClient(t *testing.T) proto.UsersClient {
	db := postgresSetup(t)
	handler := handlerSetup(t, db)

	var listener = bufconn.Listen(1024 * 1024)
	go func() {
		log.Fatal(handler.Serve(listener))
	}()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}))
	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close()
	})

	return proto.NewUsersClient(conn)
}

func handlerSetup(t *testing.T, db *sqlx.DB) *grpc.Server {
	manager, err := auth.NewTokenManager(jwtKey, tokenTTL)
	require.NoError(t, err)
	repo := postgres.NewUsersRepository(db)
	svc := service.NewUsersService(repo, manager)
	return delivery.NewHandler(svc, mocks.LoggerMock{})
}

func postgresSetup(t *testing.T) *sqlx.DB {
	request := testcontainers.ContainerRequest{
		Image:        postgresImage,
		ExposedPorts: []string{"5432"},
		Env:          map[string]string{"POSTGRES_PASSWORD": "postgres"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	conn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", host, mappedPort.Port())
	db, err := sqlx.Connect("postgres", conn)
	require.NoError(t, err)

	query, err := ioutil.ReadFile(postgresSchemaPath)
	require.NoError(t, err)

	_, err = db.Exec(string(query))
	require.NoError(t, err)

	t.Cleanup(func() {
		container.Terminate(context.Background())
	})

	return db
}
