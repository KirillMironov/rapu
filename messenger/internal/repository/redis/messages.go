package redis

import (
	"encoding/json"
	"github.com/KirillMironov/rapu/messenger/domain"
	"github.com/go-redis/redis"
)

type MessagesRepository struct {
	client *redis.Client
}

func NewMessagesRepository(client *redis.Client) *MessagesRepository {
	return &MessagesRepository{client}
}

func (m *MessagesRepository) Publish(message domain.Message, roomId string) error {
	encoded, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return m.client.Publish(roomId, encoded).Err()
}

func (m *MessagesRepository) Subscribe(roomId string) *redis.PubSub {
	return m.client.Subscribe(roomId)
}

func (m *MessagesRepository) Save(message domain.Message, roomId string) error {
	return m.client.RPush(roomId, message).Err()
}

func (m *MessagesRepository) Get(roomId string) (messages []domain.Message, _ error) {
	return messages, m.client.LRange(roomId, 0, -1).ScanSlice(&messages)
}
