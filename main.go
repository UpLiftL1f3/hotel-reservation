package main

import (
	"flag"

	"github.com/UpLiftL1f3/hotel-reservation/api"
	"github.com/UpLiftL1f3/hotel-reservation/api/middleware"
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
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(db.NewMongoUserStore(client))
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		app            = *fiber.New(config)
		auth           = app.Group("/api")
		apiV1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))

		// -> since you wrote apiV1 you inherit the "/api/v1" and add it to .../admin
		admin = apiV1.Group("/admin", middleware.AdminAuth)
	)

	// -> Auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// -> User Handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandlerDeleteUser)

	// -> hotel Handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// -> Rooms Handlers
	apiV1.Get("/room", roomHandler.HandleGetRooms)
	apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	//TODO: Cancel a booking
	// -> Booking Handlers
	apiV1.Get("/booking", bookingHandler.HandleGetBookings)
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Put("/booking/:id", bookingHandler.HandleUpdateBooking)

	//-> Admin Handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	app.Listen(*listenAddr)

}
