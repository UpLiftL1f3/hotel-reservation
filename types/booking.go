package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID    primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	FromDate  time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate  time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	NumGuests int                `bson:"numGuests,omitempty" json:"numGuests,omitempty"`
	Canceled  bool               `bson:"canceled" json:"canceled"`
}

type UpdateBookingParams struct {
	RoomID    string    `json:"roomID,omitempty"`
	FromDate  time.Time `json:"fromDate,omitempty"`
	TillDate  time.Time `json:"tillDate,omitempty"`
	NumGuests int       `json:"numGuests,omitempty"`
	Canceled  bool      `json:"canceled,omitempty"`
}
