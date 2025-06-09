package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/stew/slice"
	"rentads_app/database"
	"strconv"
	"strings"
)

func GetAdverts(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		c.Accepts("application/json")

		olderThan, _ := strconv.ParseUint(
			c.Query("older_than"),
			10,
			32,
		)

		availableCities := []string{"nn", "msc", "spb", "soc", "crm", "kzn"}
		city := c.Query("city")
		if !slice.Contains(availableCities, city) {
			city = "nn"
		}

		availableTypes := []int{1, 2}
		rentType, _ := strconv.Atoi(c.Query("rent_type"))
		if !slice.Contains(availableTypes, rentType) {
			rentType = 0
		}
		roomType, _ := strconv.Atoi(c.Query("room_type"))
		if !slice.Contains(availableTypes, roomType) {
			roomType = 0
		}

		var districts []string
		districtsQuery := c.Query("districts")
		if districtsQuery == "" {
			districts = []string{}
		} else {
			districts = strings.Split(districtsQuery, "|")
		}

		var subDistricts []string
		subDistrictQuery := c.Query("sub_district")
		if subDistrictQuery == "" {
			subDistricts = []string{}
		} else {
			subDistricts = strings.Split(subDistrictQuery, "|")
		}

		var metro []string
		metroQuery := c.Query("metro")
		if metroQuery == "" {
			metro = []string{}
		} else {
			metro = strings.Split(metroQuery, "|")
		}

		keyWords := strings.ToLower(c.Query("key_words", ""))

		adverts, err := db.GetAdvertsFromDB(
			ctx,
			olderThan,
			city,
			rentType,
			roomType,
			districts,
			subDistricts,
			metro,
			keyWords,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch adverts",
			})
		}

		return c.Status(fiber.StatusOK).JSON(adverts)
	}
}
