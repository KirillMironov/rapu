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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

const (
	postgresImage = "postgres:12.7-alpine3.14"
	jwtKey        = "secretKey"
	tokenTTL      = time.Minute
)

func newClient(t *testing.T) (proto.UsersClient, func()) {
	db, terminate := postgresSetup(t)
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

	return proto.NewUsersClient(conn), func() {
		conn.Close()
		terminate()
	}
}

func handlerSetup(t *testing.T, db *sqlx.DB) *grpc.Server {
	manager, err := auth.NewTokenManager(jwtKey, tokenTTL)
	require.NoError(t, err)
	repo := postgres.NewUsersRepository(db)
	svc := service.NewUsersService(repo, manager)
	return delivery.NewHandler(svc, mocks.LoggerMock{})
}

func postgresSetup(t *testing.T) (*sqlx.DB, func()) {
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

	m, err := migrate.New("file://../testdata", conn)
	require.NoError(t, err)
	err = m.Up()
	require.NoError(t, err)

	return db, func() {
		container.Terminate(context.Background())
	}
}
