package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fName, lName string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fName, lName),
		FirstName: fName,
		LastName:  lName,
		Password:  fmt.Sprintf("%s_%s", fName, lName),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if roomIDS == nil {
		roomIDS = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store *db.Store, uid, roomID primitive.ObjectID, from, till time.Time, numGuests int, canceled bool) *types.Booking {
	booking := &types.Booking{
		RoomID:    roomID,
		FromDate:  from,
		TillDate:  till,
		UserID:    uid,
		NumGuests: numGuests,
		Canceled:  canceled,
	}

	booking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("booking: ", booking.ID)
	return booking
}
