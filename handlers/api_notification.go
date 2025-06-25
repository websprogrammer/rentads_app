package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"rentads_app/database"
	"rentads_app/schemas"
)

func AddNotification(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		c.Accepts("application/json")

		notification := new(schemas.Notification)

		if err := c.BodyParser(notification); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse notification in request",
			})
		}

		err := db.UpsertNotificationToDB(
			ctx,
			notification,
		)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to add notification",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{})
	}
}
