package service

import "github.com/KirillMironov/rapu/messenger/internal/domain"

type Clients struct {
	messagesService MessagesService
}

type MessagesService interface {
	Reader(client domain.Client, done chan<- struct{})
	Writer(client domain.Client, done <-chan struct{})
}

func NewClients(messagesService MessagesService) *Clients {
	return &Clients{messagesService}
}

func (c *Clients) Connect(client domain.Client) {
	done := make(chan struct{})
	go c.messagesService.Writer(client, done)
	go c.messagesService.Reader(client, done)
}
