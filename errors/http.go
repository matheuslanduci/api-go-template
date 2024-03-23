package errors

import "github.com/gofiber/fiber/v3"

type HttpError struct {
	Code     string      `json:"code"`
	Status   int         `json:"status"`
	Message  string      `json:"message"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func (e HttpError) Error() string {
	return e.Message
}

func InvalidSession() HttpError {
	return HttpError{
		Code:     "invalid_session",
		Status:   fiber.StatusUnauthorized,
		Message:  "The session is invalid.",
		Metadata: nil,
	}
}

func ServerError() HttpError {
	return HttpError{
		Code:     "internal_server_error",
		Status:   fiber.StatusInternalServerError,
		Message:  "There was an error processing the request.",
		Metadata: nil,
	}
}

func ValidationError(metadata interface{}) HttpError {
	return HttpError{
		Code:     "validation_error",
		Status:   fiber.StatusBadRequest,
		Message:  "The request contains invalid data.",
		Metadata: metadata,
	}
}

func GuestOnly() HttpError {
	return HttpError{
		Code:    "guest_only",
		Status:  fiber.StatusUnauthorized,
		Message: "The user must not be authenticated.",
	}
}

func InvalidRequest() HttpError {
	return HttpError{
		Code:     "invalid_request",
		Status:   fiber.StatusBadRequest,
		Message:  "The body of the request is invalid.",
		Metadata: nil,
	}
}

func NotFound() HttpError {
	return HttpError{
		Code:     "not_found",
		Status:   fiber.StatusNotFound,
		Message:  "The resource was not found.",
		Metadata: nil,
	}
}
