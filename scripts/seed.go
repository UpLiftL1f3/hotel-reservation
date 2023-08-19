package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/UpLiftL1f3/hotel-reservation/api"
	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/db/fixtures"
)

var ()

func main() {
	ctx := context.Background()

	var err error
	client, err := db.GenerateClient()
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   hotelStore,
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	fmt.Println("Print user ->", user)

	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(admin))
	fmt.Println("Print admin ->", admin)

	hotel := fixtures.AddHotel(store, "hotel1", "bermuda", 5, nil)
	fmt.Println("Print hotel ->", hotel)

	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	fmt.Println("Print room ->", room)

	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2), 2, false)
	fmt.Println("Print booking ->", booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random_hotel_name %d", i)
		location := fmt.Sprintf("random_hotel_location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
