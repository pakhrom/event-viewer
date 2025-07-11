package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DBService struct {
	EducationalEvents mongo.Collection
	Client            mongo.Client
}

func InitDB(uri string) (DBService, error) {
	Client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		return DBService{}, err
	}

	// defer func() {
	// 	if err := Client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	DB := Client.Database("EducationalEventsTest")
	return DBService{
		Client:            *Client,
		EducationalEvents: *DB.Collection("EducationalEvents"),
	}, nil
}

func (s *DBService) CloseConnection() error {
	if err := s.Client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}
