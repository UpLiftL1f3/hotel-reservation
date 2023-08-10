package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type RoomHandler struct {
	store *db.Store
}

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate`
	TillDate  time.Time `json:"tillDate`
	NumGuests int       `json:"numGuests"`
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

// -> HELPER FUNCTIONS
func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room on a date that has already passed")
	}
	return nil
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {

	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		return false, err
	}

	if len(bookings) > 0 {
		return false, nil
	}

	// fmt.Println("bookings: ", bookings)

	return true, nil
}

// -> CRUD HANDLER FUNCTIONS

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	// fmt.Println("Made it to handle BookRoom 2")

	user, err := GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	isAvailable, err := h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if !isAvailable || err != nil {
		if err != nil {
			return err
		}
		return c.Status(http.StatusInternalServerError).JSON(&GenericResponse{
			Type: "error",
			Msg:  "room is already booked",
		})
	}

	booking := types.Booking{
		RoomID:    roomID,
		UserID:    user.ID,
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumGuests: params.NumGuests,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})

	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
