package service

import "github.com/KirillMironov/rapu/messenger/domain"

type ClientsService struct {
	messagesService domain.MessagesService
}

func NewClientsService(messagesService domain.MessagesService) *ClientsService {
	return &ClientsService{messagesService}
}

func (c *ClientsService) Connect(client domain.Client) {
	done := make(chan struct{})
	go c.messagesService.Writer(client, done)
	go c.messagesService.Reader(client, done)
}
