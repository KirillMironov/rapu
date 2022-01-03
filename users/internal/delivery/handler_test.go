package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/users/internal/delivery/proto"
	"github.com/KirillMironov/rapu/users/test/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

const (
	email    = "lisa@gmail.com"
	password = "qwerty"
	username = "Lisa"
)

var ctx = context.Background()

func newClient() (proto.UsersClient, func()) {
	var listener = bufconn.Listen(1024 * 1024)
	var server = grpc.NewServer()

	proto.RegisterUsersServer(server, &Handler{
		service: mocks.UsersServiceMock{},
		logger:  mocks.LoggerMock{},
	})

	go func() {
		log.Fatal(server.Serve(listener))
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
	}
}

func TestHandler_SignUp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		username           string
		email              string
		password           string
		expectedStatusCode codes.Code
	}{
		{username, email, password, codes.OK},
		{"", email, password, codes.InvalidArgument},
		{username, "", password, codes.InvalidArgument},
		{username, email, "", codes.InvalidArgument},
		{"", "", password, codes.InvalidArgument},
		{"", email, "", codes.InvalidArgument},
		{username, "", "", codes.InvalidArgument},
		{"", "", "", codes.InvalidArgument},
	}

	client, closeConn := newClient()
	defer closeConn()

	for _, tc := range testCases {
		tc := tc

		resp, err := client.SignUp(ctx, &proto.SignUpRequest{
			Username: tc.username,
			Email:    tc.email,
			Password: tc.password,
		})

		e, _ := status.FromError(err)
		assert.Equal(t, tc.expectedStatusCode, e.Code())

		if tc.expectedStatusCode == codes.OK {
			assert.NotEmpty(t, resp.GetAccessToken())
		} else {
			assert.Empty(t, resp.GetAccessToken())
		}
	}
}

func TestHandler_SignIn(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		email              string
		password           string
		expectedStatusCode codes.Code
	}{
		{email, password, codes.OK},
		{"", password, codes.InvalidArgument},
		{email, "", codes.InvalidArgument},
		{"", "", codes.InvalidArgument},
	}

	client, closeConn := newClient()
	defer closeConn()

	for _, tc := range testCases {
		tc := tc

		resp, err := client.SignIn(ctx, &proto.SignInRequest{
			Email:    tc.email,
			Password: tc.password,
		})

		e, _ := status.FromError(err)
		assert.Equal(t, tc.expectedStatusCode, e.Code())

		if tc.expectedStatusCode == codes.OK {
			assert.NotEmpty(t, resp.GetAccessToken())
		} else {
			assert.Empty(t, resp.GetAccessToken())
		}
	}
}

func TestHandler_Authenticate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		token              string
		expectedStatusCode codes.Code
	}{
		{"token", codes.OK},
		{"", codes.InvalidArgument},
	}

	client, closeConn := newClient()
	defer closeConn()

	for _, tc := range testCases {
		tc := tc

		resp, err := client.Authenticate(ctx, &proto.AuthRequest{AccessToken: tc.token})

		e, _ := status.FromError(err)
		assert.Equal(t, tc.expectedStatusCode, e.Code())

		if tc.expectedStatusCode == codes.OK {
			assert.NotEmpty(t, resp.GetUserId())
		} else {
			assert.Empty(t, resp.GetUserId())
		}
	}
}
