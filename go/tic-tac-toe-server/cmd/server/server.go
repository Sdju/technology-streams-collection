package server

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	id        string
	client1   *client
	client2   *client
	gameField []playerType
	turn      playerType
}

type roomMessage struct {
	room    *room
	message string
	player  *client
}

type wsMessage struct {
	Cmd    string `json:"cmd"`
	Arg    string `json:"arg"`
	Player int    `json:"player"`
}

type wsGameStateMessage struct {
	Id        string       `json:"roomId"`
	You       playerType   `json:"you"`
	Enemy     playerType   `json:"enemy"`
	GameField []playerType `json:"gameField"`
	Turn      playerType   `json:"turn"`
}

type Server struct {
	rooms        map[string]*room
	register     chan *websocket.Conn
	roomMessages chan roomMessage
	unregister   chan string

	app *fiber.App
	id  uuid.UUID
}

func (r room) getGameDtoState(player *client) wsGameStateMessage {
	var you *client
	var enemy *client
	if player == r.client1 {
		you = r.client1
		enemy = r.client2
	} else {
		you = r.client2
		enemy = r.client1
	}
	enemyRole := playerNone
	if enemy != nil {
		enemyRole = enemy.role
	}
	return wsGameStateMessage{
		r.id,
		you.role,
		enemyRole,
		r.gameField,
		r.turn,
	}
}

func (r room) sendMessage(message string) error {
	if err := r.client1.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}

	if r.client2 == nil {
		return nil
	}

	if err := r.client2.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}
	return nil
}
func (r room) sendCurStateMessage(message roomMessage) error {
	msg1 := r.getGameDtoState(message.room.client1)
	msg2 := r.getGameDtoState(message.room.client2)

	fmt.Println("room now", r)
	if err := message.room.client1.connection.WriteJSON(msg1); err != nil {
		return err
	}
	if message.room.client2 == nil {
		return nil
	}

	if err := message.room.client2.connection.WriteJSON(msg2); err != nil {
		return err
	}
	return nil
}

func (s Server) httpWsHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (s Server) httpWsRoomId(c *websocket.Conn) {
	defer func() {
		s.unregister <- c.Params("roomId")
	}()

	s.register <- c

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return
		}

		if messageType == websocket.TextMessage {
			room := s.rooms[c.Params("roomId")]
			player := room.client1
			if player.connection != c {
				player = room.client2
			}
			s.roomMessages <- roomMessage{room, string(message), player}
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}

func (s Server) httpGetRoomId(c *fiber.Ctx) error {
	return c.SendString(uuid.New().String())
}

func New() *Server {
	srv := Server{
		make(map[string]*room),
		make(chan *websocket.Conn),
		make(chan roomMessage),
		make(chan string),
		fiber.New(),
		uuid.New(),
	}

	srv.app.Use(cors.New())

	srv.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	srv.app.Use("/ws", srv.httpWsHandler)

	go srv.runHub()

	srv.app.Get("/ws/:roomId", websocket.New(srv.httpWsRoomId))

	srv.app.Get("/get-room-id", srv.httpGetRoomId)

	log.Fatal(srv.app.Listen(":3000"))

	return &srv
}

func (s Server) registerWs(connection *websocket.Conn) {
	roomId := connection.Params("roomId")
	currentRoom, ok := s.rooms[roomId]
	if !ok {
		player := client{connection, playerNone}

		gameField := make([]playerType, 9)

		for i := range gameField {
			gameField[i] = playerNone
		}

		currentRoom = &room{
			roomId,
			&player,
			nil,
			gameField,
			playerX,
		}

		s.rooms[roomId] = currentRoom

		msg := currentRoom.getGameDtoState(&player)

		connection.WriteJSON(msg)
		log.Println("room was created #", roomId)
		return
	}

	currentRoom.client2 = &client{connection, ' '}

	msg := currentRoom.getGameDtoState(currentRoom.client2)

	connection.WriteJSON(msg)

	log.Println("room was fulfilled #", roomId)
}

func (s Server) applyMessage(curMessage roomMessage) {
	log.Println("message received:", curMessage)

	currentRoom := curMessage.room

	var dat wsMessage
	if err := json.Unmarshal([]byte(curMessage.message), &dat); err != nil {
		log.Println("message wrong format:", curMessage)
		return
	}

	switch dat.Cmd {
	case "setRole":
		if dat.Player == 1 {
			currentRoom.client1.role = playerType(dat.Arg[0])
		} else {
			currentRoom.client2.role = playerType(dat.Arg[0])
		}

		fmt.Println("room now", currentRoom)
		if err := currentRoom.sendCurStateMessage(curMessage); err != nil {
			log.Println("write error:", err)
			s.unregister <- currentRoom.id
		}
	case "makeTurn":
		if err := currentRoom.sendCurStateMessage(curMessage); err != nil {
			log.Println("write error:", err)
			s.unregister <- currentRoom.id
		}
	default:
		log.Println("message received:", currentRoom)
	}
}

func (s Server) unregisterWs(roomId string) {
	currentRoom := s.rooms[roomId]
	currentRoom.client1.connection.WriteMessage(websocket.CloseMessage, []byte{})
	currentRoom.client2.connection.WriteMessage(websocket.CloseMessage, []byte{})
	currentRoom.client1.connection.Close()
	currentRoom.client2.connection.Close()

	delete(s.rooms, roomId)

	log.Println("room was closed #", roomId)
}

func (s Server) runHub() {
	for {
		select {
		case connection := <-s.register:
			s.registerWs(connection)

		case curMessage := <-s.roomMessages:
			s.applyMessage(curMessage)

		case roomId := <-s.unregister:
			s.unregisterWs(roomId)
		}
	}
}
