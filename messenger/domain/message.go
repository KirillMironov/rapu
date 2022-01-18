package domain

import "github.com/go-redis/redis"

type Message struct {
	From string `json:"from"`
	Text string `json:"text"`
}

type MessagesService interface {
	Reader(client Client)
	Writer(client Client, done <-chan struct{})
}

type MessagesRepository interface {
	Publish(message Message, roomId string) error
	Subscribe(roomId string) *redis.PubSub
}
