package db

import (
	"context"
	"fmt"

	"github.com/UpLiftL1f3/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	bookingCollection = "bookings"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, types.UpdateBookingParams) (*types.Booking, error)
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(bookingCollection),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	resp, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := resp.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}
	resp := s.collection.FindOne(ctx, filter)

	var booking *types.Booking
	if err := resp.Decode(&booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, updateParams types.UpdateBookingParams) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	m := bson.M{
		"$set": updateParams,
	}

	fmt.Printf("INFO FOR UPDATE %s, %#v", oid, updateParams)

	resp, err := s.collection.UpdateByID(ctx, oid, m)
	if err != nil {
		return nil, fmt.Errorf("%v (HELPPP)", err)
	}

	var booking *types.Booking

	if resp.MatchedCount == 1 && resp.ModifiedCount == 1 {
		err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Update operation did not match or modify exactly one document")
	}

	return booking, nil
}
