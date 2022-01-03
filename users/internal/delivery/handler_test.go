package delivery

import (
	"context"
	"github.com/KirillMironov/rapu/users/domain"
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
	bufSize  = 1024 * 1024
	email    = "lisa@gmail.com"
	password = "qwerty"
	username = "Lisa"
)

var (
	lis *bufconn.Listener
	ctx = context.Background()
)

func init() {
	lis = bufconn.Listen(bufSize)

	s := grpc.NewServer()
	proto.RegisterUsersServer(s, &Handler{
		service: mocks.UsersServiceMock{},
		logger:  mocks.LoggerMock{},
	})

	go func() {
		log.Fatal(s.Serve(lis))
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestHandler_SignUp(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		user               domain.User
		expectedStatusCode codes.Code
	}{
		{domain.User{Username: username, Email: email, Password: password}, codes.OK},
		{domain.User{}, codes.InvalidArgument},
	}

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewUsersClient(conn)

	for _, tc := range testCases {
		tc := tc

		resp, err := client.SignUp(ctx, &proto.SignUpRequest{
			Username: tc.user.Username,
			Email:    tc.user.Email,
			Password: tc.user.Password,
		})

		e, _ := status.FromError(err)
		assert.Equal(t, tc.expectedStatusCode, e.Code())

		if e.Code() == codes.OK {
			assert.NotEmpty(t, resp.GetAccessToken())
		} else {
			assert.Empty(t, resp.GetAccessToken())
		}
	}
}
