package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"log"
)

// 'X' | 'O' | ' '

type playerType rune

const (
	playerX    playerType = 'X'
	playerO    playerType = 'O'
	playerNone playerType = ' '
)

type client struct {
	connection *websocket.Conn
	role       playerType
}

type room struct {
	id      string
	client1 *client
	client2 *client
}

type roomMessage struct {
	room    string
	message string
}

type wsMessage struct {
	Cmd    string `json:"cmd"`
	Arg    string `json:"arg"`
	Player int    `json:"player"`
}

var rooms = make(map[string]*room)
var register = make(chan *websocket.Conn)
var roomMessages = make(chan roomMessage)
var unregister = make(chan string)

func (r room) sendMessage(message string) error { // p в данном случае называют receiver-ом.
	if err := r.client1.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}
	if err := r.client2.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}
	return nil
}

func runHub() {
	for {
		select {
		case connection := <-register:
			roomId := connection.Params("roomId")
			currentRoom, ok := rooms[roomId]
			if !ok {
				rooms[roomId] = &room{roomId, &client{connection, playerNone}, nil}
				log.Println("room was created #", roomId)
				continue
			}

			currentRoom.client2 = &client{connection, ' '}
			log.Println("room was fulfilled #", roomId)

		case curMessage := <-roomMessages:
			log.Println("message received:", curMessage)

			currentRoom := rooms[curMessage.room]

			var dat wsMessage
			if err := json.Unmarshal([]byte(curMessage.message), &dat); err != nil {
				log.Println("message wrong format:", curMessage)
				continue
			}

			switch dat.Cmd {
			case "setRole":
				if dat.Player == 1 {
					currentRoom.client1.role = playerType(dat.Arg[0])
				} else {
					currentRoom.client2.role = playerType(dat.Arg[0])
				}
				if err := currentRoom.sendMessage(curMessage.message); err != nil {
					log.Println("write error:", err)
					unregister <- currentRoom.id
				}
				fmt.Println("room now", currentRoom)
			case "makeTurn":
				if err := currentRoom.sendMessage(curMessage.message); err != nil {
					log.Println("write error:", err)
					unregister <- currentRoom.id
				}
			default:
				log.Println("message received:", currentRoom)
			}

		case roomId := <-unregister:
			currentRoom := rooms[roomId]
			currentRoom.client1.connection.WriteMessage(websocket.CloseMessage, []byte{})
			currentRoom.client2.connection.WriteMessage(websocket.CloseMessage, []byte{})
			currentRoom.client1.connection.Close()
			currentRoom.client2.connection.Close()

			// Remove the client from the hub
			delete(rooms, roomId)

			log.Println("room was closed #", roomId)
		}
	}
}

func main() {
	app := fiber.New()

	id := uuid.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	go runHub()

	app.Get("/ws/:roomId", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c.Params("roomId")
		}()

		// Register the client
		register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				roomMessages <- roomMessage{c.Params("roomId"), string(message)}
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}

	}))

	app.Get("/get-id", func(c *fiber.Ctx) error {
		return c.SendString(id.String())
	})

	log.Fatal(app.Listen(":3000"))
}
