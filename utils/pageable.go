package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func GetPageParams(c fiber.Ctx) (int, int) {
	page, err := strconv.Atoi(c.Query("page", "1"))

	if err != nil {
		page = 1
	}

	size, err := strconv.Atoi(c.Query("size", "10"))

	if err != nil {
		size = 10
	}

	return page, size
}
