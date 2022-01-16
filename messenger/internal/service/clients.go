package service

import "github.com/KirillMironov/rapu/messenger/domain"

type ClientsService struct {
	messagesService domain.MessagesService
}

func NewClientsService(messagesService domain.MessagesService) *ClientsService {
	return &ClientsService{messagesService}
}

func (c *ClientsService) Connect(client domain.Client) {
	defer client.Conn.Close()
	done := make(chan struct{})

	go c.messagesService.Writer(client, done)
	c.messagesService.Reader(client)
}

func (c *ClientsService) disconnect(client domain.Client, done chan<- struct{}) {
	client.Conn.Close()
	done <- struct{}{}
}
