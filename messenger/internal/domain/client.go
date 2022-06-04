package domain

import "github.com/gorilla/websocket"

type Client struct {
	UserId   string
	ToUserId string
	Conn     *websocket.Conn
}
