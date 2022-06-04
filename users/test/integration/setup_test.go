//go:build integration

package integration

import (
	"context"
	"fmt"
	"github.com/KirillMironov/rapu/users/internal/delivery"
	"github.com/KirillMironov/rapu/users/internal/delivery/proto"
	"github.com/KirillMironov/rapu/users/internal/repository"
	"github.com/KirillMironov/rapu/users/internal/service"
	"github.com/KirillMironov/rapu/users/test/mock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	postgresImage       = "postgres:12.7-alpine3.14"
	postgresPort        = "5432"
	postgresPasswordEnv = "POSTGRES_PASSWORD"
	postgresPassword    = "postgres"
	postgresSchemaPath  = "../../config/schema.sql"
	jwtKey              = "secretKey"
	tokenTTL            = time.Minute
)

func newClient(t *testing.T) proto.UsersClient {
	t.Helper()

	var (
		db       = newPostgres(t)
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

	return proto.NewUsersClient(conn)
}

func newHandler(t *testing.T, db *sqlx.DB) *grpc.Server {
	t.Helper()

	jwtService, err := service.NewJWTManager(jwtKey, tokenTTL)
	require.NoError(t, err)
	usersRepository := repository.NewUsers(db)
	usersService := service.NewUsers(usersRepository, jwtService, mock.Logger{})

	return delivery.NewHandler(usersService)
}

func newPostgres(t *testing.T) *sqlx.DB {
	t.Helper()

	var ctx = context.Background()

	containerRequest := testcontainers.ContainerRequest{
		Image:        postgresImage,
		ExposedPorts: []string{postgresPort},
		Env:          map[string]string{postgresPasswordEnv: postgresPassword},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	mappedPort, err := container.MappedPort(ctx, postgresPort)
	require.NoError(t, err)

	connectionString := fmt.Sprintf("postgres://postgres:%s@%s:%s/postgres?sslmode=disable",
		postgresPassword, host, mappedPort.Port())

	db, err := sqlx.Connect("postgres", connectionString)
	require.NoError(t, err)

	query, err := ioutil.ReadFile(postgresSchemaPath)
	require.NoError(t, err)

	_, err = db.Exec(string(query))
	require.NoError(t, err)

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	return db
}
