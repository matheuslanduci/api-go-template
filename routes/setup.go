package routes

import (
	"github.com/gofiber/fiber/v3"
	"matheuslanduci.com/api-fiber/middlewares"
	"matheuslanduci.com/api-fiber/services"
)

func Setup(app *fiber.App, services *services.Services, middlewares *middlewares.Middlewares) {
	AuthenticationRoutes(app, services, middlewares)
}
