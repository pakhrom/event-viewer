package Event

import (
	"context"
	"time"

	"github.com/billowdev/clog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EventData struct {
	Name        string    `bson:"name" json:"name,omitempty" example:"Концерт"`
	Description string    `bson:"description" json:"description,omitempty" example:"Описание мероприятия"`
	StartDate   time.Time `bson:"startDate" json:"startDate" example:"2025-07-11T18:00:00Z"`
	EndDate     time.Time `bson:"endDate" json:"endDate" example:"2025-07-11T21:00:00Z"`
}

type Event struct {
	ID        bson.ObjectID `bson:"_id" json:"id" example:"60b8d295f1d2c916c4d5f1a3"`
	EventData EventData
}

func (s EventData) Write(collection *mongo.Collection) (bson.ObjectID, error) {
	insertedResult, err := collection.InsertOne(context.TODO(), map[string]EventData{
		"EventData": s,
	})
	if err != nil {
		return bson.NilObjectID, err
	}
	return insertedResult.InsertedID.(bson.ObjectID), nil
}

func Read(collection *mongo.Collection, ID bson.ObjectID) (Event, error) {
	var loaded Event
	err := collection.FindOne(context.TODO(), bson.M{"_id": ID}).Decode(&loaded)
	if err != nil {
		return loaded, err
	}
	return loaded, nil
}

func ReadMany(collection mongo.Collection, filter bson.M, findOptions options.FindOptionsBuilder) ([]Event, error) {
	var events []Event
	cursor, err := collection.Find(context.TODO(), filter, &findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var event Event
		clog.Debug("Event struct: %+v", event)
		err := cursor.Decode(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s EventData) Update(collection *mongo.Collection, ID bson.ObjectID) error {
	_, err := collection.UpdateByID(context.TODO(), ID, s)
	return err
}

func Delete(collection *mongo.Collection, ID bson.ObjectID) error {
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": ID})
	return err
}
