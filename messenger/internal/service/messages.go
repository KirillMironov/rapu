package service

import (
	"encoding/json"
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/gorilla/websocket"
)

type MessagesService struct {
	repository domain.MessagesRepository
	logger     Logger
}

type Logger interface {
	Error(args ...interface{})
}

func NewMessagesService(repository domain.MessagesRepository, logger Logger) *MessagesService {
	return &MessagesService{repository, logger}
}

func (m *MessagesService) Reader(client domain.Client) {
	roomId := m.getRoomId(client.UserId, client.ToUserId)

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway) {
				m.logger.Error(err)
			}
			return
		}

		var message = domain.Message{
			From: client.UserId,
			Text: string(p),
		}

		err = m.repository.Publish(message, roomId)
		if err != nil {
			m.logger.Error(err)
			return
		}

		err = m.repository.Save(message, roomId)
		if err != nil {
			m.logger.Error(err)
			return
		}
	}
}

func (m *MessagesService) Writer(client domain.Client, done <-chan struct{}) {
	roomId := m.getRoomId(client.UserId, client.ToUserId)

	history, err := m.loadChatHistory(roomId)
	if err != nil {
		m.logger.Error(err)
		return
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, history)
	if err != nil {
		m.logger.Error(err)
		return
	}

	sub := m.repository.Subscribe(roomId)
	defer sub.Close()

	for {
		select {
		case msg := <-sub.Channel():
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				m.logger.Error(err)
				return
			}
		case <-done:
			return
		}
	}
}

func (m *MessagesService) loadChatHistory(roomId string) ([]byte, error) {
	messages, err := m.repository.Get(roomId)
	if err != nil {
		return nil, err
	}

	return json.Marshal(messages)
}

func (m *MessagesService) getRoomId(from, to string) string {
	if to < from {
		return to + ":" + from
	}
	return from + ":" + to
}
