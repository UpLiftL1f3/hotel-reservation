package api

import (
	"fmt"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
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

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		//-> Get the ID
		id = c.Params("id")
	)

	//* Fetch the user from the database
	user, err := h.UserStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)

}
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.UserStore.GetUsers(c.Context())
	fmt.Println("USERS: ", users)
	if err != nil {
		return err
	}

	fmt.Println("USERS: ", users)
	return c.JSON(users)
}
