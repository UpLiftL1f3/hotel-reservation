package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/UpLiftL1f3/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func RunCrudTest(t *testing.T) {}

func TestPostUser(t *testing.T) {
	// testing database
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	app.Post("/", UserHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "James",
		LastName:  "Doe",
		Password:  "randomPassword",
	}
	b, _ := json.Marshal(params)

	// creating a request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	// running the test
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	// fmt.Println(user)

	if len(user.ID) == 0 {
		t.Errorf("expected a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Make sure not to send the encrypted password in the resp body")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected userName %s, but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected userName %s, but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected userName %s, but got %s", params.Email, user.Email)
	}
}

func TestGetUser(t *testing.T) {
	tbd := setup(t)
	defer tbd.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tbd.User)

	// First Create User to get ID
	app.Post("/", UserHandler.HandlePostUser)
	app.Get("/", UserHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "GETsome@GETfoo.com",
		FirstName: "James",
		LastName:  "Goodie",
		Password:  "randomPassword",
	}
	b, _ := json.Marshal(params)

	// creating a request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	// running the test
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	// fmt.Println(user)

	if len(user.ID) == 0 {
		t.Errorf("expected a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Make sure not to send the encrypted password in the resp body")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected userName %s, but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected userName %s, but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected userName %s, but got %s", params.Email, user.Email)
	}

	getParams := types.User{
		ID: user.ID,
	}
	getBody, _ := json.Marshal(getParams)

	// Now Get User
	request := httptest.NewRequest("GET", "/", bytes.NewReader(getBody))
	request.Header.Add("Content-Type", "application/json")

	var getUser types.User
	response, err := app.Test(request)
	if err != nil {
		t.Errorf(err.Error())
	}

	json.NewDecoder(response.Body).Decode(&getUser)
	// fmt.Println(user)

}
