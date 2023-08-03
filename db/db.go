package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//-> this file is going to be used for General DB things (aka helper files, etc)

const DB_URI = "mongodb://localhost:27017"
const DBNAME = "hotel-reservation"
const TestDBNAME = "hotel-reservation-test"

func GenerateClient() (*mongo.Client, error) {
	//-> Setting up connection to the Database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		return nil, err
	}

	return client, nil
}
