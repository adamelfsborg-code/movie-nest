package data

import (
	"time"

	"github.com/adamelfsborg-code/movie-nest/db"
	"github.com/google/uuid"
)

type UserData struct{}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func NewUser(name string) *User {
	return &User{
		Name: name,
	}
}

func (u *UserData) List() []User {
	var users []User
	db.Store.Model(&users).Select()
	return users
}

func (u *UserData) Create(user User) error {
	_, err := db.Store.Model(&user).Insert()
	return err
}
