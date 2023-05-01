package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"log"
)

type HttpController struct {
	app *fiber.App
}

func (s HttpController) httpGetRoomId(c *fiber.Ctx) error {
	return c.SendString(uuid.New().String())
}

func NewHttpController(app *fiber.App) *HttpController {
	controller := HttpController{app}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	app.Get("/get-room-id", controller.httpGetRoomId)

	return &controller
}

func (s HttpController) start() {
	log.Fatal(s.app.Listen(":3000"))
}
