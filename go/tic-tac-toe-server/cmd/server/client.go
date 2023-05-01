package server

import "github.com/gofiber/websocket/v2"

type client struct {
	connection *websocket.Conn
	role       playerType
}
