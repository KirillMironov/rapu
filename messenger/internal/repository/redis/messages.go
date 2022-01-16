package redis

import (
	"encoding/json"
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/go-redis/redis"
)

type MessagesRepository struct {
	client *redis.Client
}

func (m *MessagesRepository) Publish(message domain.Message) error {
	roomId := m.getRoomId(message.From, message.To)

	encoded, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return m.client.Publish(roomId, encoded).Err()
}

func (m *MessagesRepository) Subscribe(from, to string) *redis.PubSub {
	roomId := m.getRoomId(from, to)

	return m.client.Subscribe(roomId)
}

func (m *MessagesRepository) getRoomId(from, to string) string {
	if to < from {
		return to + ":" + from
	}
	return from + ":" + to
}
