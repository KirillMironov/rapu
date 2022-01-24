package domain

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type Message struct {
	From string `json:"from"`
	Text string `json:"text"`
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

type MessagesService interface {
	Reader(client Client, done chan<- struct{})
	Writer(client Client, done <-chan struct{})
}

type MessagesBus interface {
	Publish(message Message, roomId string) error
	Subscribe(roomId string) *redis.PubSub
}

type MessagesRepository interface {
	Save(message Message, roomId string) error
	Get(roomId string) ([]Message, error)
}
