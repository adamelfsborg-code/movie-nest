package data

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type ShelfData struct {
	DB pg.DB
}

type Shelf struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	RoomID    uuid.UUID `json:"room_id" db:"room_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type ShelfMovies struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Movies []*Movie  `json:"movies" db:"movies"`
}

func NewShelf(name string, roomID uuid.UUID) *Shelf {
	return &Shelf{
		Name:   name,
		RoomID: roomID,
	}
}

func (s *ShelfData) CreateShelf(shelf Shelf) error {
	_, err := s.DB.Model(&shelf).Insert()
	return err
}

func (s *ShelfData) GetShelvesByRoomID(roomID uuid.UUID) []Shelf {
	var shelf []Shelf
	s.DB.Model(&shelf).Where("room_id = ?", &roomID).Select()
	return shelf
}

func (s *ShelfData) GetShelfMoviesByID(shelfID uuid.UUID) []Movie {
	var movies []Movie
	s.DB.Model(&movies).Where("shelf_id = ?", &shelfID).Select()
	if len(movies) > 0 {
		return movies
	}
	return make([]Movie, 0)
}

func (s *ShelfData) GetShelfInfoByID(shelfID uuid.UUID) Shelf {
	var shelf Shelf
	s.DB.Model(&shelf).Where("id = ?", &shelfID).Select()
	return shelf
}
