package api

import (
	"context"
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
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	//-> Setting up connection to the Database
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoURI))
	if err != nil {
		panic(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}
