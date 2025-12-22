package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/stew/slice"
	"rentads_app/config"
	"rentads_app/database"
	"rentads_app/schemas"
	"strconv"
)

func DeleteFeedback(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := context.Background()
		c.Accepts("application/json")

		id, _ := strconv.ParseUint(
			c.Query("feedback_id"),
			10,
			32,
		)
		err := db.DeleteFeedbackFromDB(
			ctx,
			id,
		)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to delete feedback",
			})
		}

		return c.Redirect("/get_feedbacks")
	}
}

func GetFeedbacks(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := context.Background()
		c.Accepts("application/json")

		id, _ := strconv.ParseUint(
			c.Query("id"),
			10,
			32,
		)

		availableCities := config.GetSupportedCities()
		city := c.Query("city")
		if !slice.Contains(availableCities, city) {
			city = "nn"
		}

		feedbacks, err := db.GetFeedbacksFromDB(
			ctx,
			id,
			city,
		)
		if err != nil {
			println("Error getting feedbacks from DB: {}", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch feedbacks",
			})
		}

		return c.Render("index", fiber.Map{
			"Title":     "Feedbacks",
			"Feedbacks": feedbacks,
		})

	}
}

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
