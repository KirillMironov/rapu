//go:build integration

package integration

import (
	"context"
	"fmt"
	"github.com/KirillMironov/rapu/users/domain"
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
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"time"
)

const (
	jwtKey   = "secretKey"
	tokenTTL = time.Minute
)

func newClient() (proto.UsersClient, func(), error) {
	var listener = bufconn.Listen(1024 * 1024)

	db, terminate, err := postgresSetup()
	if err != nil {
		return nil, nil, err
	}
	usersService, err := serviceSetup(db)
	if err != nil {
		return nil, nil, err
	}
	handler := delivery.NewHandler(usersService, mocks.LoggerMock{})

	go func() {
		log.Fatal(handler.Serve(listener))
	}()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}))
	if err != nil {
		log.Fatal(err)
	}

	return proto.NewUsersClient(conn), func() {
		conn.Close()
		terminate()
	}, nil
}

func serviceSetup(db *sqlx.DB) (domain.UsersService, error) {
	manager, err := auth.NewManager(jwtKey, tokenTTL)
	if err != nil {
		return nil, err
	}
	repo := postgres.NewUsersRepository(db)
	return service.NewUsersService(repo, manager), nil
}

func postgresSetup() (*sqlx.DB, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12.7-alpine3.14",
		ExposedPorts: []string{"5432"},
		Env:          map[string]string{"POSTGRES_PASSWORD": "postgres"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
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
	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	conn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", ip, mappedPort.Port())
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, nil, err
	}

	m, err := migrate.New("file://../testdata", conn)
	if err != nil {
		return nil, nil, err
	}
	err = m.Up()
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		container.Terminate(context.Background())
	}, nil
}
