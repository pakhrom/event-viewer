package routers

import (
	// _ "help/docs"
	Event "help/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// get all events with sort by event date
// @Summary Get Events
// @Schemes
// @Tags Events
// @Accept json
// @Param filter query string false "Filter events by date" Enums(past, present, future)
// @Produce json
// @Success 200 {array} Event.Event
// @Router /events [get]
func RequestEvents(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		searchBy := c.QueryParam("filter")

		filter := bson.M{}
		sortBy := options.Find()

		switch searchBy {
		case "past":
			filter = bson.M{
				"endDate": bson.M{"$lt": time.Now()},
			}
			sortBy.SetSort(bson.D{{Key: "endDate", Value: -1}})

		case "present":
			filter = bson.M{
				"startDate": bson.M{"$lte": time.Now()},
				"endDate":   bson.M{"$gt": time.Now()},
			}
			sortBy.SetSort(bson.D{{Key: "endDate", Value: -1}})

		case "future":
			filter = bson.M{
				"startDate": bson.M{"$gt": time.Now()},
			}
			sortBy.SetSort(bson.D{{Key: "startDate", Value: 1}})

		}

		events, err := Event.ReadMany(*collection, filter, *sortBy)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": err.Error(),
			})

		}
		return c.JSON(200, events)
	}
}

// func CreateEvent(collection *mongo.Collection) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var eventData Event.EventData
// 		if err := c.Bind(&eventData); err != nil {
// 		}
// 	}
// }
