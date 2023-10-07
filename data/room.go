package data

import (
	"time"

	"github.com/adamelfsborg-code/movie-nest/db"
	"github.com/google/uuid"
)

type RoomData struct{}

type Room struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type RoomUser struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoomID    uuid.UUID `json:"room_id" db:"room_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type RoomWithUser struct {
	Room  *Room   `json:"room" db:"room"`
	Users []*User `json:"users" db:"users"`
}

type UserRooms struct {
	User  *User   `json:"user" db:"user"`
	Rooms []*Room `json:"rooms" db:"rooms"`
}

func NewRoom(name string) *Room {
	return &Room{
		Name: name,
	}
}

func NewRoomUser(roomID, userID uuid.UUID) *RoomUser {
	return &RoomUser{
		RoomID: roomID,
		UserID: userID,
	}
}

func (r *RoomData) CreateRoom(room Room) error {
	_, err := db.Store.Model(&room).Insert()
	return err
}

func (r *RoomData) ListRooms() []Room {
	var rooms []Room
	db.Store.Model(&rooms).Select()
	return rooms
}

func (r *RoomData) GetRoomByID(roomID uuid.UUID) (*Room, error) {
	var room Room
	err := db.Store.Model(&room).Where("id = ?", &roomID).Select()
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomData) AddUserToUser(roomUser RoomUser) error {
	_, err := db.Store.Model(&roomUser).Insert()
	return err
}

func (r *RoomData) ListRoomsWithUsers() []RoomWithUser {
	var roomsWithUsers []RoomWithUser

	_, _ = db.Store.Query(&roomsWithUsers, `
		SELECT 
			jsonb_build_object
			(
				'id', r.id, 'name', r."name", 'timestamp', r."timestamp"
			) AS room,
			jsonb_agg
			(
				jsonb_build_object
				(
					'id', u.id, 'name', u."name", 'timestamp', u."timestamp"
				)
			) AS users
		FROM rooms r
		LEFT JOIN room_users ru ON r.id = ru.room_id
		LEFT JOIN users u ON u.id = ru.user_id
		GROUP BY r.id
	`)
	return roomsWithUsers
}

func (r *RoomData) GetRoomWithUsersByID(roomID uuid.UUID) RoomWithUser {
	var roomWithUsers RoomWithUser

	_, _ = db.Store.Query(&roomWithUsers, `
		SELECT 
			jsonb_build_object
			(
				'id', r.id, 'name', r."name", 'timestamp', r."timestamp"
			) AS room,
			jsonb_agg
			(
				jsonb_build_object
				(
					'id', u.id, 'name', u."name", 'timestamp', u."timestamp"
				)
			) AS users
		FROM rooms r
		LEFT JOIN room_users ru ON r.id = ru.room_id
		LEFT JOIN users u ON u.id = ru.user_id
		WHERE r.id = ?
		GROUP BY r.id
	`, &roomID)

	return roomWithUsers
}

func (r *RoomData) GetUserRoomsByID(userID uuid.UUID) UserRooms {
	var userRooms UserRooms

	_, _ = db.Store.Query(&userRooms, `
		SELECT 
			jsonb_build_object
			(
				'id', u.id, 'name', u."name", 'timestamp', u."timestamp"
			) AS user,
			jsonb_agg
			(
				jsonb_build_object
				(
					'id', r.id, 'name', r."name", 'timestamp', r."timestamp"
				)
			) AS rooms
		FROM users u
		LEFT JOIN room_users ru ON u.id = ru.user_id
		LEFT JOIN rooms r ON r.id = ru.room_id
		WHERE u.id = ?
		GROUP BY u.id
	`, &userID)

	return userRooms
}
