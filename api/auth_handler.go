package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/UpLiftL1f3/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type GenericResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(GenericResponse{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

// -> LOGIN USER
func (a *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params = AuthParams{}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	fmt.Println(params)

	// -> FIND USER
	user, err := a.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNilDocument) {
			return invalidCredentials(c)
		}
		return err
	}
	fmt.Println(user)

	// -> Compare Password
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}

	authResp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}

	return c.JSON(authResp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"userID":  user.ID,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := os.Getenv("JWT_SECRET")
	fmt.Println("Secret from Getenv: ", secret)
	result := fmt.Sprintf("%T", secret)
	fmt.Println("type of secret: ", result)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err, token)
	}
	return tokenString
}
