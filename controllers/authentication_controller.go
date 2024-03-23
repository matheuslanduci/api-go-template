package controllers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"matheuslanduci.com/api-fiber/database/models"
	"matheuslanduci.com/api-fiber/dto"
	"matheuslanduci.com/api-fiber/errors"
	"matheuslanduci.com/api-fiber/services"
)

type AuthenticationController struct {
	service   *services.AuthenticationService
	validator *services.XValidator
}

func NewAuthenticationController(services *services.Services) *AuthenticationController {
	return &AuthenticationController{
		service:   services.AuthenticationService,
		validator: services.Validator,
	}
}

func (controller *AuthenticationController) CreateSessionWithPassword(c fiber.Ctx) error {
	request := &dto.CreateSessionWithPasswordRequest{}

	if err := json.Unmarshal(c.Body(), request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors.InvalidRequest())
	}

	if err := controller.validator.Validate(request); len(err) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors.ValidationError(err))
	}

	session, err := controller.service.CreateSessionWithPassword(request)

	if err != nil {
		println(err.Error())

		switch err.Error() {
		case services.ErrUserNotFound:
			return c.Status(fiber.StatusNotFound).JSON(
				errors.HttpError{
					Code:     "user_not_found",
					Status:   fiber.StatusNotFound,
					Message:  "The user was not found.",
					Metadata: nil,
				},
			)
		case services.ErrInvalidPassword:
			return c.Status(fiber.StatusUnauthorized).JSON(
				errors.HttpError{
					Code:     "invalid_password",
					Status:   fiber.StatusUnauthorized,
					Message:  "The password is invalid.",
					Metadata: nil,
				},
			)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
		}
	}

	handler := csrf.HandlerFromContext(c)

	if handler != nil {
		if err := handler.DeleteToken(c); err != nil {
			println(err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth",
		Value:    session.SessionToken,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(time.Hour * 1),
	})

	if request.Remember {
		c.Cookie(&fiber.Cookie{
			Name:     "remember",
			Value:    *session.RememberToken,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Lax",
			Expires:  time.Now().Add(time.Hour * 24 * 7),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (controller *AuthenticationController) DeleteSession(c fiber.Ctx) error {
	authCookie := c.Cookies("auth")
	rememberCookie := c.Cookies("remember")

	err := controller.service.DeleteSession(authCookie, &rememberCookie)

	if err != nil {
		println(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
	}

	handler := csrf.HandlerFromContext(c)

	if handler != nil {
		if err := handler.DeleteToken(c); err != nil {
			println(err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
		}
	}

	c.ClearCookie("auth_")
	c.ClearCookie("remember")

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (controller *AuthenticationController) GetCurrentUser(c fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	return c.Status(fiber.StatusOK).JSON(user)
}

func (controller *AuthenticationController) CreateSessionWithRememberToken(c fiber.Ctx) error {
	rememberCookie := c.Cookies("remember")

	if rememberCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			errors.HttpError{
				Code:     "invalid_remember_token",
				Status:   fiber.StatusUnauthorized,
				Message:  "The remember token is invalid.",
				Metadata: nil,
			},
		)
	}

	session, err := controller.service.CreateSessionWithRememberToken(rememberCookie)

	if err != nil {
		println(err.Error())

		switch err.Error() {
		case services.ErrInvalidSession:
			return c.Status(fiber.StatusUnauthorized).JSON(errors.InvalidSession())
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
		}
	}

	handler := csrf.HandlerFromContext(c)

	if handler != nil {
		if err := handler.DeleteToken(c); err != nil {
			println(err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(errors.ServerError())
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth",
		Value:    session.SessionToken,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(time.Hour * 1),
	})

	c.ClearCookie("remember")

	c.Cookie(&fiber.Cookie{
		Name:     "remember",
		Value:    *session.RememberToken,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
		Expires:  time.Now().Add(time.Hour * 24 * 7),
	})

	return c.Status(fiber.StatusNoContent).Send(nil)
}
