package data

import (
	"fmt"
	"time"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/shared"
	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/golang-jwt/jwt/v5"
)

type UserData struct {
	Env  config.Environments
	DB   pg.DB
	Nats *nats.Conn
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"max=20,min=3"`
	Password  string    `json:"password" db:"password" validate:"max=50,min=10"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type UserResp struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"max=20,min=3"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func NewRegisterUser(name, password string) (*User, error) {
	fmt.Println(password, name)
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
	_, err := u.DB.Model(&user).Insert()
	return err
}

func (u *UserData) List() []User {
	var users []User
	u.DB.Model(&users).Select()
	return users
}

func (u *UserData) Login(name, password string) (string, error) {
	user := u.getUserByName(name)

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

func (u *UserData) GetUserInfoByID(userID uuid.UUID) User {
	var user User
	u.DB.Model(&user).Where("id = ?", userID).Select()
	return user
}

func (u *UserData) CheckUserExistsByID(userID uuid.UUID) bool {
	var user []User
	u.DB.Model(&user).Where("id = ?", userID).Select()
	if len(user) > 0 {
		return true
	}
	return false
}

func (u *UserData) GetUsersInRoom(roomID uuid.UUID, userID uuid.UUID, excludeSelf bool) []User {
	var users []User

	query := u.DB.Model(&users).Join(fmt.Sprintf("JOIN room_users ru ON %q.id = ru.user_id", "user")).Where("ru.room_id = ?", &roomID)

	if excludeSelf {
		query.Where("ru.user_id <> ?", &userID)
	}

	query.Select()

	return users
}

func (u *UserData) getUserByName(name string) User {
	var user User

	u.DB.Model(&user).Where("name = ?", name).Select()
	return user
}
