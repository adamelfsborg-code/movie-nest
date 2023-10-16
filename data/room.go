package data

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type RoomData struct {
	DB pg.DB
}

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

type RoomInfo struct {
	Room    *Room          `json:"room" db:"room"`
	Users   []*User        `json:"users" db:"users"`
	Shelves []*ShelfMovies `json:"shelves" db:"shelves"`
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
	_, err := r.DB.Model(&room).Insert()
	return err
}

func (r *RoomData) ListRooms() []Room {
	var rooms []Room
	r.DB.Model(&rooms).Select()
	return rooms
}

func (r *RoomData) GetRoomByID(roomID uuid.UUID) (*Room, error) {
	var room Room
	err := r.DB.Model(&room).Where("id = ?", &roomID).Select()
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomData) GetRoomInfoByID(roomID uuid.UUID) (*RoomInfo, error) {
	var roomInfo RoomInfo

	_, _ = r.DB.Query(&roomInfo, `
		WITH ShelfMovies AS (
			SELECT s.id AS shelf_id, 
				jsonb_build_object(
					'id', m.id,
					'movie_id', m.movie_id,
					'timestamp', m."timestamp"
				) AS movie
			FROM shelves s
			JOIN movies m ON s.id = m.shelf_id
			WHERE s.room_id = '26e126f7-84d9-41dd-843f-44931badece5'
		)
		SELECT 
			jsonb_build_object(
				'id', r.id, 'name', r."name", 'timestamp', r."timestamp"
			) AS room,
			(
				SELECT jsonb_agg(
					jsonb_build_object(
						'id', u.id, 'name', u."name", 'timestamp', u."timestamp"
					)
				)
				FROM room_users ru
				JOIN users u ON u.id = ru.user_id
				WHERE ru.room_id = r.id
			) AS users,
			(
				SELECT jsonb_agg(
					jsonb_build_object(
						'id', s.id, 'name', s."name", 'timestamp', s."timestamp", 'movies', 
						COALESCE((SELECT jsonb_agg(movie) FROM ShelfMovies WHERE shelf_id = s.id), '[]'::jsonb)
					)
				)
				FROM shelves s
				WHERE s.room_id = r.id
			) AS shelves
		FROM rooms r
		where r.id = ?
	`, &roomID)
	return &roomInfo, nil
}

func (r *RoomData) AddUserToUser(roomUser RoomUser) error {
	_, err := r.DB.Model(&roomUser).Insert()
	return err
}

func (r *RoomData) ListRoomsWithUsers() []RoomWithUser {
	var roomsWithUsers []RoomWithUser

	_, _ = r.DB.Query(&roomsWithUsers, `
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

	_, _ = r.DB.Query(&roomWithUsers, `
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

	_, _ = r.DB.Query(&userRooms, `
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

func (r *RoomData) GetAvailableUsers(roomID uuid.UUID, userID uuid.UUID, searchTerm string, excludeSelf bool, excludeExisting bool) []User {
	var users []User

	query := r.DB.Model(&users).Where(`LOWER("user"."name") LIKE ?`, fmt.Sprintf("%%%s%%", searchTerm))

	if excludeSelf {
		query = query.Where(`"user".id <> ?`, userID)
	}

	if excludeExisting {
		query = query.Where(`"user".id NOT IN (SELECT user_id FROM room_users WHERE room_id = ?)`, roomID)
	}

	query.Select()

	return users
}
