package data

import (
	"time"

	"github.com/adamelfsborg-code/movie-nest/db"
	"github.com/google/uuid"
)

type ShelfData struct{}

type Shelf struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	RoomID    uuid.UUID `json:"room_id" db:"room_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func NewShelf(name string, roomID uuid.UUID) *Shelf {
	return &Shelf{
		Name:   name,
		RoomID: roomID,
	}
}

func (s *ShelfData) CreateShelf(shelf Shelf) error {
	_, err := db.Store.Model(&shelf).Insert()
	return err
}
