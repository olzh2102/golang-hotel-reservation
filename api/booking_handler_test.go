package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/olzh2102/golang-hotel-reservation/db/fixtures"
	"github.com/olzh2102/golang-hotel-reservation/types"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		nonAuthUser = fixtures.AddUser(tdb.Store, "beer", "drinkwater", false)
		user        = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel       = fixtures.AddHotel(tdb.Store, "botqa", "a", 4, nil)
		room        = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)

		from    = time.Now()
		till    = time.Now().AddDate(0, 0, 5)
		booking = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)

		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(tdb.User))
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var bookingRes *types.Booking
	if err := json.NewDecoder(res.Body).Decode(&bookingRes); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("non 200 code got %d", res.StatusCode)
	}
	if bookingRes.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingRes.ID)
	}
	if bookingRes.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingRes.UserID)
	}

	// * non auth user
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 code got %d", res.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user      = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(tdb.Store, "botqa", "a", 4, nil)
		room      = fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)

		from    = time.Now()
		till    = time.Now().AddDate(0, 0, 5)
		booking = fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)

		app = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		admin          = app.Group("/", JWTAuthentication(tdb.User), AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", res.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(res.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	// * test non-admin cannot access the bookings

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	res, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized but got %d", res.StatusCode)
	}
}
