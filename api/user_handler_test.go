package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/olzh2102/golang-hotel-reservation/types"
)

func TestPostUser(t *testing.T) {
	// * arrange database
	tdb := setup(t)
	defer tdb.teardown(t)

	// * arrange server
	app := fiber.New()

	// * arrange user controller
	userHandler := NewUserHandler(tdb.User)

	// * act
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "foo@bar.com",
		FirstName: "Smith",
		LastName:  "Anderson",
		Password:  "somerandompassword",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(res.Body).Decode(&user)

	// * assertions
	if len(user.ID) == 0 {
		t.Errorf("expected a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected EncryptedPassword not to be included on the json response")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected first name %s, got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected last name %s, got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s, got %s", params.Email, user.Email)
	}
}
