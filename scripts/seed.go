package main

import (
	"context"
	"fmt"
	"log"

	"github.com/UpLiftL1f3/hotel-reservation/api"
	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func main() {
	seedHotel("Bellucia", "France", 3)
	seedHotel("The cozy hotel", "The Netherlands", 4)
	seedHotel("High as a cloud", "Tennessee", 5)
	seedUser("James", "Foo", "jamesFoo@gmail.com", "superSecurePassword", false)
	seedUser("Admin", "Admin", "Admin@admin.com", "adminAdminAdmin", true)
}

func init() {
	var err error
	client, err = db.GenerateClient()
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

func seedUser(fName, lName, email, password string, isAdmin bool) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fName,
		LastName:  lName,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 199.9,
		},
		{
			Size:  "king",
			Price: 122.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}

	// fmt.Println(insertedHotel)

	fmt.Println("seeding the database")
}
