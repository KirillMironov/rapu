package repository

import (
	"github.com/KirillMironov/rapu/messenger/internal/domain"
	"github.com/go-redis/redis"
)

type Messages struct {
	client *redis.Client
}

func NewMessages(client *redis.Client) *Messages {
	return &Messages{client}
}

func (m *Messages) Save(message domain.Message, roomId string) error {
	return m.client.RPush(roomId, message).Err()
}

func (m *Messages) Get(roomId string) (messages []domain.Message, _ error) {
	return messages, m.client.LRange(roomId, 0, -1).ScanSlice(&messages)
}
