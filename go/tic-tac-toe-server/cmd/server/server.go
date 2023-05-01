package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

// 'X' | 'O' | ' '

type playerType rune

const (
	playerX    playerType = 'X'
	playerO    playerType = 'O'
	playerNone playerType = ' '
)

type wsMessage struct {
	Connection *websocket.Conn
	RoomId     string
	Cmd        string `json:"cmd"`
	Arg        string `json:"arg"`
}

type wsGameStateMessage struct {
	Id        string       `json:"roomId"`
	You       playerType   `json:"you"`
	Enemy     playerType   `json:"enemy"`
	GameField []playerType `json:"gameField"`
	Turn      playerType   `json:"turn"`
}

type Server struct {
	rooms map[string]*room

	wsHandler *WsEventHandler

	app            *fiber.App
	httpController *HttpController
	WsController   *WsController
}

func New() *Server {
	app := fiber.New()
	handler := &WsEventHandler{
		make(chan *websocket.Conn),
		make(chan wsMessage),
		make(chan *websocket.Conn),
	}

	srv := Server{
		make(map[string]*room),
		handler,
		app,
		NewHttpController(app),
		NewWsController(app, handler),
	}

	go srv.runHub()

	srv.httpController.start()

	return &srv
}

func (s Server) registerWs(connection *websocket.Conn) {
	roomId := connection.Params("roomId")
	currentRoom, ok := s.rooms[roomId]
	if !ok {
		gameField := make([]playerType, 9)

		for i := range gameField {
			gameField[i] = playerNone
		}

		currentRoom = &room{
			roomId,
			nil,
			nil,
			gameField,
			playerX,
		}

		s.rooms[roomId] = currentRoom

		log.Println("room was created #", roomId)
	}

	client := &client{connection, playerNone}

	if currentRoom.client1 == nil {
		currentRoom.client1 = client
		log.Println("client1 connect #", roomId)
	} else if currentRoom.client2 == nil {
		currentRoom.client2 = client
		log.Println("client2 connect #", roomId)
	} else {
		log.Println("room has been already fulfilled #", roomId)
		connection.Close()
		return
	}

	if err := currentRoom.sendCurStateMessage(); err != nil {

	}
}

func (s Server) applyMessage(curMessage wsMessage) {
	log.Println("message received:", curMessage)

	currentRoom := s.rooms[curMessage.RoomId]
	player := currentRoom.client1

	if curMessage.Connection != currentRoom.client1.connection {
		player = currentRoom.client2
	}

	switch curMessage.Cmd {
	case "setRole":
		currentRoom.handlerSetRole(player, curMessage.Arg)
	case "makeTurn":
		currentRoom.handlerMakeTurn(curMessage.Arg)
	default:
		log.Println("message received:", currentRoom)
	}
}

func (s Server) unregisterWs(conn *websocket.Conn) {
	roomId := conn.Params("roomId")
	currentRoom := s.rooms[roomId]

	if err := conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {

	}
	if err := conn.Close(); err != nil {

	}

	if currentRoom.client1 != nil && currentRoom.client1.connection == conn {
		currentRoom.client1 = nil
		log.Println("client1 disconnect #", roomId)
	} else if currentRoom.client2 != nil && currentRoom.client2.connection == conn {
		currentRoom.client2 = nil
		log.Println("client2 disconnect #", roomId)
	}

	if currentRoom.client1 == nil && currentRoom.client2 == nil {
		delete(s.rooms, roomId)
		log.Println("room was closed #", roomId)
	}
}

func (s Server) runHub() {
	for {
		select {
		case connection := <-s.wsHandler.register:
			s.registerWs(connection)

		case curMessage := <-s.wsHandler.message:
			s.applyMessage(curMessage)

		case connection := <-s.wsHandler.unregister:
			s.unregisterWs(connection)
		}
	}
}
