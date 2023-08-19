package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/UpLiftL1f3/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Println("token not present")
			return ErrUnAuthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		// fmt.Println(claims)

		// -> Check token expiration
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		userID := claims["userID"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnAuthorized()
		}

		// -> Set the current authenticated user to the context
		// fmt.Printf("Set user : %+v", user)
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid Signing method", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	// fmt.Println("TOKEN: ", token)

	if err != nil {
		fmt.Println("failed to parse jwt token", err)
		return nil, ErrUnAuthorized()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, ErrUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnAuthorized()
	}

	return claims, nil

}
