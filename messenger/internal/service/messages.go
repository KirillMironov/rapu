package service

import (
	"encoding/json"
	"github.com/KirillMironov/rapu/gateway/pkg/logger"
	"github.com/KirillMironov/rapu/messenger/internal/domain"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

type Messages struct {
	messagesBus        MessagesBus
	messagesRepository MessagesRepository
	logger             logger.Logger
}

type MessagesBus interface {
	Publish(message domain.Message, roomId string) error
	Subscribe(roomId string) *redis.PubSub
}

type MessagesRepository interface {
	Save(message domain.Message, roomId string) error
	Get(roomId string) ([]domain.Message, error)
}

func NewMessages(messagesBus MessagesBus, messagesRepository MessagesRepository, logger logger.Logger) *Messages {
	return &Messages{
		messagesBus:        messagesBus,
		messagesRepository: messagesRepository,
		logger:             logger,
	}
}

func (m *Messages) Reader(client domain.Client, done chan<- struct{}) {
	roomId := m.getRoomId(client.UserId, client.ToUserId)
	defer close(done)

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) &&
				!websocket.IsCloseError(err, websocket.CloseGoingAway) {
				m.logger.Error(err)
			}
			return
		}

		var message = domain.Message{
			From: client.UserId,
			Text: string(p),
		}

		err = m.messagesBus.Publish(message, roomId)
		if err != nil {
			m.logger.Error(err)
			return
		}

		err = m.messagesRepository.Save(message, roomId)
		if err != nil {
			m.logger.Error(err)
			return
		}
	}
}

func (m *Messages) Writer(client domain.Client, done <-chan struct{}) {
	roomId := m.getRoomId(client.UserId, client.ToUserId)
	defer client.Conn.Close()

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

	sub := m.messagesBus.Subscribe(roomId)
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

func (m *Messages) loadChatHistory(roomId string) ([]byte, error) {
	messages, err := m.messagesRepository.Get(roomId)
	if err != nil {
		return nil, err
	}

	return json.Marshal(messages)
}

func (m *Messages) getRoomId(from, to string) string {
	if to < from {
		return to + ":" + from
	}
	return from + ":" + to
}
