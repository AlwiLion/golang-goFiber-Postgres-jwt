package main

import (
	"github.com/AlwiLion/database"
	"github.com/AlwiLion/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.DBconn()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true, //Very important while using a HTTPonly Cookie, frontend can easily get and return back the cookie.
	}))
	routes.Setup(app)

	app.Listen(":8000")
}
