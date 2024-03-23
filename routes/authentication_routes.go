package routes

import (
	"github.com/gofiber/fiber/v3"
	"matheuslanduci.com/api-fiber/controllers"
	"matheuslanduci.com/api-fiber/middlewares"
	"matheuslanduci.com/api-fiber/services"
)

func AuthenticationRoutes(app fiber.Router, services *services.Services, middlewares *middlewares.Middlewares) {
	controller := controllers.NewAuthenticationController(services)

	app.Use("/me", middlewares.AuthenticationMiddleware.ValidateSession())
	app.Get("/me", controller.GetCurrentUser)

	app.Use("/sessions/password", middlewares.AuthenticationMiddleware.GuestOnly())
	app.Post("/sessions/password", controller.CreateSessionWithPassword)

	app.Post("/sessions/refresh", controller.CreateSessionWithRememberToken)

	app.Use("/sessions", middlewares.AuthenticationMiddleware.ValidateSession())
	app.Delete("/sessions", controller.DeleteSession)
}
