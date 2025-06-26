package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"rentads_app/database"
	"rentads_app/schemas"
)

func SendFeedback(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := context.Background()
		c.Accepts("application/json")

		feedback := new(schemas.Feedback)
		if err := c.BodyParser(feedback); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse feedback in request",
			})
		}

		err := db.AddFeedbackToDB(
			ctx,
			feedback,
		)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to add feedback",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{})
	}
}
