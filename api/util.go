package api

import (
	"net/http"

	"github.com/UpLiftL1f3/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func GetAuthenticatedUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		// -> i CHOSE internal server error over unauthorized because if we didn't get a user it means it wasn't set correctly
		return nil, c.Status(http.StatusInternalServerError).JSON(GenericResponse{
			Type: "error",
			Msg:  "internal server error",
		})
	}

	return user, nil
}
