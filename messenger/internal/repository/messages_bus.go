package repository

import (
	"encoding/json"
	"github.com/KirillMironov/rapu/messenger/internal/domain"
	"github.com/go-redis/redis"
)

type MessagesBus struct {
	client *redis.Client
}

func NewMessagesBus(client *redis.Client) *MessagesBus {
	return &MessagesBus{client}
}

func (m *MessagesBus) Publish(message domain.Message, roomId string) error {
	encoded, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return m.client.Publish(roomId, encoded).Err()
}

func (m *MessagesBus) Subscribe(roomId string) *redis.PubSub {
	return m.client.Subscribe(roomId)
}
