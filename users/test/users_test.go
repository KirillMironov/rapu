//go:build integration

package test

import (
	"context"
	"fmt"
	"github.com/KirillMironov/rapu/domain"
	"github.com/KirillMironov/rapu/internal/repository/postgres"
	"github.com/KirillMironov/rapu/internal/service"
	"github.com/KirillMironov/rapu/pkg/auth"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

const jwtKey = "secretKey"

var ctx = context.Background()

func TestUsers_SignUp(t *testing.T) {
	db, container := postgresSetup(t)
	defer container.Terminate(ctx)
	svc := usersServiceSetup(t, db, time.Minute)

	var user = domain.User{
		Username: "Lisa",
		Email:    "lisa@gmail.com",
		Password: "qwerty",
	}

	token, err := svc.SignUp(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	token, err = svc.SignUp(user)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestUsers_SignIn(t *testing.T) {
	db, container := postgresSetup(t)
	defer container.Terminate(ctx)
	svc := usersServiceSetup(t, db, time.Minute)

	var user = domain.User{
		Username: "Lisa",
		Email:    "lisa@gmail.com",
		Password: "qwerty",
	}

	token, err := svc.SignIn(user)
	assert.Error(t, err)
	assert.Empty(t, token)

	token, err = svc.SignUp(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	token, err = svc.SignIn(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	token, err = svc.SignIn(domain.User{
		Email:    user.Email,
		Password: "a",
	})
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestUsers_Authenticate(t *testing.T) {
	db, container := postgresSetup(t)
	defer container.Terminate(ctx)
	svc := usersServiceSetup(t, db, time.Millisecond*500)

	var user = domain.User{
		Username: "Lisa",
		Email:    "lisa@gmail.com",
		Password: "qwerty",
	}

	token, err := svc.SignUp(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	userId, err := svc.Authenticate(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, userId)

	time.Sleep(time.Second * 2)

	userId, err = svc.Authenticate(token)
	assert.Error(t, err)
	assert.Empty(t, userId)
}

func usersServiceSetup(t *testing.T, db *sqlx.DB, tokenTTL time.Duration) domain.UsersService {
	manager, err := auth.NewManager(jwtKey, tokenTTL)
	if err != nil {
		t.Fatal(err)
	}
	repo := postgres.NewUsersRepository(db)
	return service.NewUsersService(repo, manager)
}

func postgresSetup(t *testing.T) (*sqlx.DB, testcontainers.Container) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13.0",
		ExposedPorts: []string{"5432"},
		Env:          map[string]string{"POSTGRES_PASSWORD": "postgres"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ip, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	conn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", ip, mappedPort.Port())
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.New("file://testdata", conn)
	if err != nil {
		t.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		t.Fatal(err)
	}

	return db, container
}
