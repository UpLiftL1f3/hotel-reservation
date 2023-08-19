package api

import (
	"fmt"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	user, err := GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if !user.IsAdmin && booking.UserID != user.ID {
		return ErrUnAuthorized()
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleUpdateBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	var updateParams types.UpdateBookingParams
	if err := c.BodyParser(&updateParams); err != nil {
		return ErrBadRequest()
	}
	fmt.Printf("Main 2 (update booking params): %#v", updateParams)

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	fmt.Println("Main 3")
	user, err := GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	fmt.Println("Main 4")
	if !user.IsAdmin && booking.UserID != user.ID {
		return ErrUnAuthorized()
	}

	fmt.Println("Main 5")
	updatedBooking, err := h.store.Booking.UpdateBooking(c.Context(), id, updateParams)
	if err != nil {
		return err
	}

	return c.JSON(updatedBooking)
}
