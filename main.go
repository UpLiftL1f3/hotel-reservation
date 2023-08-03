package main

import (
	"flag"

	"github.com/UpLiftL1f3/hotel-reservation/api"
	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
)

// Create a new fiber instance with custom config
var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	//-> Setting up connection to the Database
	client, err := db.GenerateClient()
	if err != nil {
		panic(err)
	}

	//-> USER HANDLERS
	var (
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client))
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler = api.NewHotelHandler(db.NewMongoHotelStore(client), roomStore)
		app          = *fiber.New(config)
		apiV1        = app.Group("/api/v1")
	)

	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandlerDeleteUser)

	// hotel Handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	app.Listen(*listenAddr)
}
