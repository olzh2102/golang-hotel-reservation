package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/olzh2102/golang-hotel-reservation/db"
	"github.com/olzh2102/golang-hotel-reservation/types"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "james@foo.com",
		FirstName: "James",
		LastName:  "Longstaff",
		Password:  "verystrongpassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	// * arrange database
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)

	// * arrange server
	app := fiber.New()

	// * arrange user controller
	authHandler := NewAuthHandler(tdb.UserStore)

	// * act
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "verystrongpassword",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", res.StatusCode)
	}

	var authRes AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authRes); err != nil {
		t.Error(err)
	}
	if authRes.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}

	// * set the encrypted password to an empty string,
	// * because we do NOT return that in any JSON response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authRes.User) {
		t.Fatal("expected the user to be the inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	// * arrange database
	tdb := setup(t)
	defer tdb.teardown(t)
	insertTestUser(t, tdb.UserStore)

	// * arrange server
	app := fiber.New()

	// * arrange user controller
	authHandler := NewAuthHandler(tdb.UserStore)

	// * act
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "verystrongpasswordNOTCORRECT",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", res.StatusCode)
	}

	var genRes genericRes
	if err := json.NewDecoder(res.Body).Decode(&genRes); err != nil {
		t.Fatal(err)
	}
	if genRes.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genRes.Type)
	}
	if genRes.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be <invalid credential>, but got %s", genRes.Msg)
	}
}
