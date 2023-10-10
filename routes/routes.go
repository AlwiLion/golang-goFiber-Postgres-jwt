package routes

import (
	"github.com/AlwiLion/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/user")

	api.Get("/welcome", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!!")
	})
	api.Get("/get-user", controllers.User)

	api.Post("/register", controllers.Register)

	api.Post("/login", controllers.Login)
	api.Put("/updateKyc", controllers.UpdateKyc)
	api.Post("/deposit-withdraw", controllers.DepositWithdrawMoney)
	api.Post("/transfer", controllers.TransferMoney)
	api.Get("/checkBalance", controllers.CheckBalance)
}
