package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"matheuslanduci.com/api-fiber/errors"
	"matheuslanduci.com/api-fiber/services"
)

type AuthenticationMiddleware struct {
	service *services.AuthenticationService
}

func NewAuthenticationMiddleware(services *services.Services) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		service: services.AuthenticationService,
	}
}

func (middleware *AuthenticationMiddleware) ValidateSession() fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionCookie := c.Cookies("auth")

		if sessionCookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.InvalidSession())
		}

		user, err := middleware.service.GetSession(sessionCookie)

		if err != nil {
			switch err.Error() {
			case services.ErrInvalidSession:
				c.ClearCookie("auth")

				return c.Status(fiber.StatusUnauthorized).JSON(errors.InvalidSession())
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
			}
		}

		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.InvalidSession())
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func (middleware *AuthenticationMiddleware) GuestOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionCookie := c.Cookies("auth")

		if sessionCookie != "" {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.GuestOnly())
		}

		return c.Next()
	}
}
