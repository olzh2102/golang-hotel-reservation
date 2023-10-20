package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olzh2102/golang-hotel-reservation/db"
	"github.com/olzh2102/golang-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(
	store *db.Store,
	name string,
	loc string,
	rating int,
	rooms []primitive.ObjectID,
) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDs,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddRoom(
	store *db.Store,
	size string,
	ss bool,
	price float64,
	hid primitive.ObjectID,
) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hid,
	}
	insertedRoom, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(
	store *db.Store,
	uid,
	rid primitive.ObjectID,
	from,
	till time.Time,
) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
