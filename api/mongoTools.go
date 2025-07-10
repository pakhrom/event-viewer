package main

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	collection mongo.Collection
	client     mongo.Client
}

func (DB) initDB(uri string) (DB, error) {
	Client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		return DB{}, err
	}

	// defer func() {
	// 	if err := Client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	Collection := Client.Database("EducationalEventsTest").Collection("EducationalEvents")
	return DB{
		collection: *Collection,
		client:     *Client,
	}, nil
}

func (s *DB) closeConnection() error {
	if err := s.client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}
