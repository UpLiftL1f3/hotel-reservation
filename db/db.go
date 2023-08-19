package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Pagination struct {
	Limit int64
	Page  int64
}

// -> this file is going to be used for General DB things (aka helper files, etc)
const (
	DB_URI     = "mongodb://localhost:27017"
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
)

func GenerateClient() (*mongo.Client, error) {
	//-> Setting up connection to the Database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		return nil, err
	}

	return client, nil
}
