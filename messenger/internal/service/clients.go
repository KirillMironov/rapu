package service

import (
	"context"
	"errors"
	"github.com/KirillMironov/rapu/gateway/pkg/logger"
	"github.com/KirillMironov/rapu/messenger/internal/delivery/proto"
	"github.com/KirillMironov/rapu/messenger/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Clients struct {
	usersClient     proto.UsersClient
	messagesService MessagesService
	logger          logger.Logger
}

type MessagesService interface {
	Reader(client domain.Client, done chan<- struct{})
	Writer(client domain.Client, done <-chan struct{})
}

func NewClients(usersClient proto.UsersClient, messagesService MessagesService, logger logger.Logger) *Clients {
	return &Clients{
		usersClient:     usersClient,
		messagesService: messagesService,
		logger:          logger,
	}
}

func (c *Clients) Connect(ctx context.Context, accessToken string, client domain.Client) error {
	resp, err := c.usersClient.Authenticate(ctx, &proto.AuthRequest{AccessToken: accessToken})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.logger.Error(err)
			return err
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return domain.ErrEmptyParameters
		case codes.NotFound:
			return domain.ErrUserNotFound
		default:
			c.logger.Error(err)
			return err
		}
	}

	client.UserId = resp.GetUserId()
	if client.UserId == "" {
		return errors.New("user id was not found")
	}

	existsResp, err := c.usersClient.UserExists(ctx, &proto.UserExistsRequest{UserId: client.ToUserId})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.logger.Error(err)
			return err
		}
		switch st.Code() {
		case codes.NotFound:
			return domain.ErrUserNotFound
		default:
			c.logger.Error(err)
			return err
		}
	}

	if !existsResp.GetExists() {
		return domain.ErrUserNotFound
	}

	var done = make(chan struct{})
	go c.messagesService.Writer(client, done)
	go c.messagesService.Reader(client, done)

	return nil
}
