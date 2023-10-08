package data

import (
	"fmt"
	"time"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/db"
	"github.com/adamelfsborg-code/movie-nest/shared"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

type UserData struct {
	Env config.Environments
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"max=20,min=3"`
	Password  string    `json:"password" db:"password" validate:"max=50,min=10"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func NewRegisterUser(name, password string) (*User, error) {
	validate := validator.New()

	user := &User{
		Name:     name,
		Password: password,
	}

	errs := validate.Struct(user)
	if errs != nil {
		return nil, errs
	}

	hash, err := shared.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user.Password = hash
	return user, nil
}

func (u *UserData) Register(user User) error {
	_, err := db.Store.Model(&user).Insert()
	return err
}

func (u *UserData) List() []User {
	var users []User
	db.Store.Model(&users).Select()
	return users
}

func (u *UserData) Login(name, password string) (string, error) {
	user := getUserByName(name)

	valid := shared.CheckPasswordHash(password, user.Password)
	if valid == false {
		return "", fmt.Errorf("User does not exists")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.Env.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getUserByName(name string) User {
	var user User
	db.Store.Model(&user).Where("name = ?", name).Select()
	return user
}
