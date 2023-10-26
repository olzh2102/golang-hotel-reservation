package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/olzh2102/golang-hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store      *db.Store
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}
	filter := db.Map{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrNotResourceFound("rooms")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return ErrNotResourceFound("hotels")
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return ErrNotResourceFound("hotel")
	}
	return c.JSON(hotel)
}
