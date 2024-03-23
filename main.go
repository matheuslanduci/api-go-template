package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/utils/v2"
	"matheuslanduci.com/api-fiber/config"
	"matheuslanduci.com/api-fiber/database"
	"matheuslanduci.com/api-fiber/errors"
	"matheuslanduci.com/api-fiber/middlewares"
	"matheuslanduci.com/api-fiber/routes"
	"matheuslanduci.com/api-fiber/services"
)

var mode string

func setupFlags() {
	flag.StringVar(&mode, "mode", "development", "The mode the application should run")

	flag.Parse()
}

func main() {
	setupFlags()

	config := config.New(mode)

	db := database.New(config)

	redis := redis.New()

	services := &services.Services{
		Validator:             services.NewValidator(),
		AuthenticationService: services.NewAuthenticationService(db, redis),
	}

	middlewares := &middlewares.Middlewares{
		AuthenticationMiddleware: middlewares.NewAuthenticationMiddleware(services),
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			println(err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.Cors.AllowedOrigins,
		AllowMethods: config.Cors.AllowedMethods,
		AllowHeaders: config.Cors.AllowedHeaders,
	}))
	app.Use(csrf.New(csrf.Config{
		CookieName:     "csrf",
		CookieSameSite: "Lax",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUIDv4,
		Storage:        redis,
		SessionKey:     "app.csrf.token",
		CookieDomain:   config.Server.Host,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			println(err.Error())

			return c.Status(fiber.StatusForbidden).JSON(
				errors.HttpError{
					Code:     "invalid-csrf-token",
					Status:   fiber.StatusForbidden,
					Message:  "The CSRF token is invalid.",
					Metadata: nil,
				},
			)
		},
	}))

	app.Get("/csrf", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNoContent).Send(nil)
	})

	routes.Setup(app, services, middlewares)

	app.Listen(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
}
