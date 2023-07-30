package main

import (
	"context"
	"flag"

	"github.com/UpLiftL1f3/hotel-reservation/api"
	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		panic(err)
	}

	//-> USER HANDLERS
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)
}
