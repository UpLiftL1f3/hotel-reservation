package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > minFirstNameLen {
		m["firstName"] = p.FirstName
	}
	if len(p.FirstName) > minLastNameLen {
		m["lastName"] = p.LastName
	}

	return m
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first Name Length must be at least %d characters", minFirstNameLen)
	}

	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("last Name Length must be at least %d characters", minLastNameLen)
	}

	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password Length must be at least %d characters", minPasswordLen)
	}

	if !IsEmailValid(params.Email) {
		errors["Eamil"] = fmt.Sprintf("email %s is invalid", params.Email)
	}

	return errors
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e)
}

func IsValidPassword(encpw, password string) bool {
	// -> encpw == encrypted password
	//* compareHash returns an err so if result is nil and == nil then it will return TRUE else FALSE

	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(password)) == nil
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
