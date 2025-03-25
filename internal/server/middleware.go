package server

import (
	"github.com/gofiber/fiber/v2"
)

func UnauthMiddleware(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	c.Set("Content_Encodig", "gzip")

	err := c.Next()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	return nil
}
