package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olzh2102/golang-hotel-reservation/api"
	"github.com/olzh2102/golang-hotel-reservation/db"
	"github.com/olzh2102/golang-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURUI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}

func main() {
	seedHotel("Ci Senior", "Alicante", 3)
	seedHotel("Boq", "Shymkala", 5)
	hotel := seedHotel("Sleep as dead men", "Newcastle", 4)

	james := seedUser(false, "james", "foo", "james@foo.com", "verystrongpassword")
	seedUser(true, "admin", "admin", "admin@admin.com", "adminpassword123")

	seedRoom("small", true, 89.99, hotel.ID)
	seedRoom("medium", true, 189.99, hotel.ID)
	room := seedRoom("large", false, 289.99, hotel.ID)

	seedBooking(james.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return user
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	res, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking: ", res.ID)
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  "small",
			Price: 24.99,
		},
		{
			Size:  "normal",
			Price: 109.99,
		},
		{
			Size:  "kingsize",
			Price: 89.99,
		},
	}
	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
	return insertedHotel
}
