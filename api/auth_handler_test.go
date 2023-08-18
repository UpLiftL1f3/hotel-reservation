package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/UpLiftL1f3/hotel-reservation/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

// func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
// 	user, err := types.NewUserFromParams(types.CreateUserParams{
// 		Email:     "gil@foo.com",
// 		FirstName: "gil",
// 		LastName:  "bill",
// 		Password:  "gil_bill",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	resp, err := userStore.InsertUser(context.TODO(), user)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	return resp
// }

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	// insertedUser := insertTestUser(t, tdb.User)
	insertedUser := fixtures.AddUser(tdb.Store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "james_foo",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http.Status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}

	// fmt.Println(authResp)

	if authResp.Token == "" {
		t.Fatalf("expected the JWT token to be present in the auth response")
	}

	insertedUser.EncryptedPassword = ""
	fmt.Println("Inserted User: ", insertedUser)
	fmt.Println("authResp User: ", authResp.User)
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be the inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	// insertTestUser(t, tdb.User)
	fixtures.AddUser(tdb.Store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "superSecurePass",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http.Status of 200 but got %d", resp.StatusCode)
	}

	var genericResp GenericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}

	if genericResp.Type != "error" {
		t.Fatalf("expected genericResponse type to be error but got %s", genericResp.Type)
	}

	if genericResp.Msg != "invalid credentials" {
		t.Fatalf("expected genericResponse Message to be 'invalid credentials' but got %s", genericResp.Msg)
	}

}
