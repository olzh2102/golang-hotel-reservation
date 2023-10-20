package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/olzh2102/golang-hotel-reservation/db/fixtures"
)

func TestAuthenticateSuccess(t *testing.T) {
	// * arrange database
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.Store, "james", "foo", false)

	// * arrange server
	app := fiber.New()

	// * arrange user controller
	authHandler := NewAuthHandler(tdb.User)

	// * act
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "james_foo",
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
	fixtures.AddUser(tdb.Store, "james", "foo", false)

	// * arrange server
	app := fiber.New()

	// * arrange user controller
	authHandler := NewAuthHandler(tdb.User)

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
