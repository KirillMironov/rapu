package domain

import "github.com/go-redis/redis"

type Message struct {
	From string `json:"from"`
	To   string `json:"-"`
	Text string `json:"text"`
}

type MessagesService interface {
	Reader(client Client)
	Writer(client Client, done <-chan struct{})
}

type MessagesRepository interface {
	Publish(message Message) error
	Subscribe(from, to string) *redis.PubSub
}
