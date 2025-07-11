package Event

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EventData struct {
	Name        string    `bson:"name,omnitempty"`
	Description string    `bson:"description,omnitempty"`
	StartDate   time.Time `bson:"startDate,omnitempty"`
	EndDate     time.Time `bson:"startDate,omnitempty"`
}
type Event struct {
	ID primitive.ObjectID `bson:"_id,omitembty"`
	EventData
}

func (s Event) write(collection *mongo.Collection) (primitive.ObjectID, error) {
	insertedResult, err := collection.InsertOne(context.TODO(), s)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertedResult.InsertedID.(primitive.ObjectID), nil
}

func (Event) read(collection *mongo.Collection, ID primitive.ObjectID) (Event, error) {
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
		err := cursor.Decode(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
