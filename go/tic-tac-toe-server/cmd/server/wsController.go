package server

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

type WsEventHandler struct {
	register   chan *websocket.Conn
	message    chan wsMessage
	unregister chan *websocket.Conn
}

type WsController struct {
	app     *fiber.App
	handler *WsEventHandler
}

func (controller WsController) onWsRequest(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (controller WsController) onRoomRequest(c *websocket.Conn) {
	controller.handler.register <- c
	defer func() {
		controller.handler.unregister <- c
	}()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return
		}

		if messageType == websocket.TextMessage {
			var dat wsMessage
			if err := json.Unmarshal(message, &dat); err != nil {
				log.Println("message wrong format:", message)
				return
			}

			dat.Connection = c
			dat.RoomId = c.Params("roomId")

			controller.handler.message <- dat
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}

func (controller WsController) messageRender() {

}

func NewWsController(app *fiber.App, handler *WsEventHandler) *WsController {
	controller := WsController{app, handler}

	app.Use("/ws", controller.onWsRequest)
	app.Use("/ws/:roomId", websocket.New(controller.onRoomRequest))

	return &controller
}
