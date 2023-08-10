package api

import (
	"errors"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type UserHandler struct {
	UserStore db.UserStore
}

// -> CREATE NEW USER
func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
	}
}

// -> DELETE USER
func (h *UserHandler) HandlerDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	if err := h.UserStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}

	return c.JSON(map[string]string{"Deleted": userID})
}

// -> UPDATE USER
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		// updates bson.M
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	// how to determine which document to update
	filter := bson.M{"_id": oid}
	if err := h.UserStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": userID})
}

// -> POST USER
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.UserStore.InsertUser(c.Context(), user)
	if err != nil {
		return nil
	}

	return c.JSON(insertedUser)
}

// -> GET USER
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		//-> Get the ID
		id = c.Params("id")
	)

	//* Fetch the user from the database
	user, err := h.UserStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"msg": "not found"})
		}
	}

	return c.JSON(user)

}

// -> GET USERS
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.UserStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	// fmt.Println("USERS: ", users)
	return c.JSON(users)
}
