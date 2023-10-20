package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/olzh2102/golang-hotel-reservation/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(tdb.Store, "botqa", "a", 4, nil)
	room := fixtures.AddRoom(tdb.Store, "small", true, 4.4, hotel.ID)

	from := time.Now()
	till := time.Now().AddDate(0, 0, 5)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, from, till)
	fmt.Println(booking)
}
