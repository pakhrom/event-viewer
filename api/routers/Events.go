package routers

import (
	// _ "help/docs"

	"context"
	Event "help/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// @Summary Get Event
// @Schemes
// @Tags Events
// @Accept json
// @Param id path string true "ID"
// @Produce json
// @Success 200 {object} Event.Event
// @Router /events/{id} [get]
func RequestEvent(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			idRaw string = c.Param("id")
		)
		if id, err := bson.ObjectIDFromHex(idRaw); err != nil {
			return err
		} else if event, err := Event.Read(collection, id); err != nil {
			return err
		} else {
			return c.JSON(http.StatusOK, event)
		}

	}
}

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
			return err

		}
		return c.JSON(200, events)
	}
}

// Create new event
// @Summary Create Event
// @Tags Events
// @Accept json
// @Param event_data body Event.EventData true "Event"
// @Produce json
// @Success 200 {object} bson.ObjectID
// @Router /events [post]
func CreateEvent(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		var eventData Event.EventData
		if err := c.Bind(&eventData); err != nil {
			return err
		}
		id, err := eventData.Write(collection)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, id)
	}
}

// Update Existing event
// @Summary Update Event
// @Tags Events
// @Accept json
// @Param event_data body Event.EventData true "EventData"
// @Param id path string true "ID"
// @Produce json
// @Success 200 {object} bson.ObjectID
// @Failure 404 {object} error
// @Router /events/{id} [put]
func UpdateEvent(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			eventData Event.EventData
			idRaw     string = c.Param("id")
		)
		if err := c.Bind(&eventData); err != nil {
			return err
		}
		if id, err := bson.ObjectIDFromHex(idRaw); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})
		} else if collection.FindOne(context.TODO(), bson.M{"_id": id}) == nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Event not found",
			})
		} else if err := eventData.Update(collection, id); err != nil {
			return err
		} else {
			return c.JSON(http.StatusOK, id)
		}
	}
}

// Delete event
// @Summary Delete Event
// @Tags Events
// @Accept json
// @Param id path string true "Event ID"      // исправлен тип на string
// @Produce json
// @Success 200 {object} bool
// @Failure 404 {object} error       // стандартная ошибка Echo
// @Router /events/{id} [delete]
func DeleteEvent(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			idRaw string = c.Param("id")
		)
		if id, err := bson.ObjectIDFromHex(idRaw); err != nil {
			return err
		} else if collection.FindOne(context.TODO(), bson.M{"_id": id}) == nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Event not found",
			})
		} else if err := Event.Delete(collection, id); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, true)
	}
}
