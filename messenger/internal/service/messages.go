package service

import (
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/gorilla/websocket"
	"log"
)

type MessagesService struct {
	repository domain.MessagesRepository
}

func NewMessagesService(repository domain.MessagesRepository) *MessagesService {
	return &MessagesService{repository}
}

func (m *MessagesService) Reader(client domain.Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var message = domain.Message{
			From: client.UserId,
			To:   client.ToUserId,
			Text: string(p),
		}

		err = m.repository.Publish(message)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (m *MessagesService) Writer(client domain.Client, done <-chan struct{}) {
	sub := m.repository.Subscribe(client.UserId, client.ToUserId)
	defer sub.Close()

	for {
		select {
		case msg := <-sub.Channel():
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				log.Println(err)
				return
			}
		case <-done:
			return
		}
	}
}
