package api

import (
	"context"
	"log"
	"testing"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName       = "hotel-reservation-test"
	testMongoURI = "mongodb://localhost:27017"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	//-> Setting up connection to the Database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)

	return &testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
			Hotel:   hotelStore,
		},
	}
}
