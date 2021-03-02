/*
The package implements the backend logic of Arendator application.

It requires the following parts:
- server logging
- getting adverts with requested params
- sending user feedback to server
- sending token with necessary application params
*/
package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/stew/slice"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func init() {
	// Set logging for the server
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	// Open or create the log file with required permissions
	file, err := os.OpenFile(
		filepath.Join(exPath, "/application.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666)

	if err != nil {
		log.Fatal(err)
	}

	// Set the output destination for the standard logger
	log.SetOutput(file)
}

// Get adverts with requested params from database
func (db *DB) getAdverts(c echo.Context) error {
	var query []bson.M

	lastDate, err := strconv.ParseUint(
		c.QueryParam("last_date"),
		10,
		32)

	if err == nil && lastDate != 0 {
		query = append(
			query,
			bson.M{
				"Date": bson.M{
					"$lt": lastDate,
				},
			},
		)
	}

	availableCities := []string{"nn", "msc", "spb"}
	city := c.QueryParam("city")
	cityMap := bson.M{"City": "nn"}
	if slice.Contains(availableCities, city) {
		cityMap = bson.M{"City": city}
	}
	query = append(query, cityMap)

	index := mgo.Index{
		Key:  []string{"$text:Description"},
		Name: "Description_text",
	}
	_ = db.adverts.EnsureIndex(index)

	keyWords := c.QueryParam("key_words")
	if keyWords != "" {
		query = append(
			query,
			bson.M{
				"$text": bson.M{
					"$search": strings.Replace(keyWords, "|", " ", -1),
				},
			})
	}

	availableTypes := []int{1, 2}

	rentType, err := strconv.Atoi(c.QueryParam("rent_type"))
	if err == nil && slice.Contains(availableTypes, rentType) {
		query = append(query, bson.M{"RentType": rentType})
	}

	roomType, err := strconv.Atoi(
		c.QueryParam("room_type"))

	if err == nil && slice.Contains(availableTypes, roomType) {
		query = append(query, bson.M{"RoomType": roomType})
	}

	districts := c.QueryParam("districts")
	if districts != "" {
		query = append(
			query,
			bson.M{
				"District": bson.M{
					"$in": strings.Split(districts, "|"),
				},
			})
	}

	subDistrict := c.QueryParam("sub_district")
	if subDistrict != "" {
		query = append(
			query,
			bson.M{
				"SubDistrict": bson.M{
					"$in": strings.Split(subDistrict, "|"),
				},
			})
	}

	metro := c.QueryParam("metro")
	if metro != "" {
		query = append(
			query,
			bson.M{
				"Metro": bson.M{
					"$in": strings.Split(metro, "|"),
				},
			})
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$and": query,
			},
		},
		{"$sort": bson.M{"Date": -1}},
		{"$limit": 20},
	}

	fmt.Println(pipeline)


	var adverts []Advert
	err = db.adverts.Pipe(pipeline).All(&adverts)
	if err != nil {
		results := struct {
			status string
		}{
			DbError,
		}
		log.Println(DbError)
		return c.JSON(http.StatusInternalServerError, results)
	}

	var lastDateValue uint64
	if len(adverts) > 0 && len(adverts) == 20 {
		lastDateValue = adverts[len(adverts)-1].Date
	}

	results := Results{
		Adverts:  adverts,
		LastDate: lastDateValue,
	}
	return c.JSON(http.StatusOK, results)
}

// Get Feedbacks
func (db *DB) getFeedbacks(c echo.Context) error  {

	pipeline := []bson.M{}

	var feedbacks []Feedback
	err:= db.feedbacks.Pipe(pipeline).All(&feedbacks)
	if err != nil {
		results := struct {
			status string
		}{
			DbError,
		}
		log.Println(DbError)
		return c.JSON(http.StatusInternalServerError, results)
	}

	return c.JSON(http.StatusOK, feedbacks)
}

// Send user feedback to server
func (db *DB) sendFeedback(c echo.Context) error {
	var results struct {
		status string
	}

	city := c.QueryParam("city")

	postId, err := strconv.ParseUint(
		c.QueryParam("post_id"),
		0,
		64)

	if err != nil {
		results.status = wrongPostId
		return c.JSON(http.StatusInternalServerError, results)
	}

	availableTypes := []uint64{1, 2, 3, 4}
	feedbackType, err := strconv.ParseUint(
		c.QueryParam("type"),
		10,
		32)
	if err != nil || !slice.Contains(availableTypes, feedbackType) {
		results.status = wrongFeedbackType
		return c.JSON(http.StatusInternalServerError, results)
	}

	var message string
	if feedbackType == 4 {
		message = c.QueryParam("message")
	}

	feedback := &Feedback{
		ID:      bson.NewObjectId(),
		PostId:  postId,
		City:    city,
		Type:    uint(feedbackType),
		Message: message,
	}

	err = db.feedbacks.Insert(feedback)
	if err != nil {
		log.Println(DbError)

		results.status = DbError
		return c.JSON(http.StatusInternalServerError, results)
	}

	results.status = "Feedback added."
	return c.JSON(http.StatusOK, results)
}

// Send token with necessary application params
func (db *DB) sendToken(c echo.Context) error {
	var results struct {
		status string
	}

	availableCities := []string{"nn", "msc", "spb"}
	city := c.QueryParam("city")
	if !slice.Contains(availableCities, city) {
		city = "nn"
	}

	tokenId := c.QueryParam("token")
	keyWords := c.QueryParam("key_words")
	districts := c.QueryParam("districts")

	rentType, err := strconv.ParseUint(
		c.QueryParam("rent_type"),
		10,
		32,
	)
	if err != nil {
		log.Println(wrongRentType)

		results.status = wrongRentType
		return c.JSON(http.StatusInternalServerError, results)
	}

	roomType, err := strconv.ParseUint(
		c.QueryParam("room_type"),
		10,
		32,
	)
	if err != nil {
		log.Println(wrongRoomType)

		results.status = wrongRoomType
		return c.JSON(http.StatusInternalServerError, results)
	}

	notifications, err := strconv.ParseUint(
		c.QueryParam("notifications"),
		10,
		32,
	)
	if err != nil {
		log.Println(wrongNotificationType)

		results.status = wrongNotificationType
		return c.JSON(http.StatusInternalServerError, results)
	}

	selector := bson.M{"Token": tokenId}
	update := bson.M{
		"$set": bson.M{
			"City":          city,
			"KeyWords":      keyWords,
			"Districts":     districts,
			"RentType":      uint(rentType),
			"RoomType":      uint(roomType),
			"Notifications": uint(notifications),
			"Updated":       time.Now().Unix(),
		},
	}

	info, err := db.tokens.Upsert(selector, update)
	key := "added"
	if info != nil && info.Updated > 0 {
		key = "updated"
	}

	if err != nil {
		log.Println(DbError)

		results.status = DbError
		return c.JSON(http.StatusInternalServerError, results)
	}

	msg := fmt.Sprintf("Token %s %s.", tokenId, key)
	results.status = msg
	log.Println(msg)
	return c.JSON(http.StatusOK, results)
}

func main() {
	// Database initialization
	session, err := mgo.Dial(mongoURI)
	if err != nil {
		panic(err)
	}
	dbInstance := session.DB("rentads")
	db := &DB{
		session:   session,
		adverts:   dbInstance.C("adverts"),
		tokens:    dbInstance.C("tokens"),
		feedbacks: dbInstance.C("feedbacks"),
	}

	defer session.Close()

	// Create Echo instance
	e := echo.New()

	// Routes
	e.GET("/", db.getAdverts)
	e.GET("/send_token", db.sendToken)
	e.GET("/send_feedback", db.sendFeedback)
	e.GET("/feedbacks", db.getFeedbacks)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
