package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/pkg/themoviedb"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type ShelfData struct {
	DB   pg.DB
	Env  config.Environments
	Nats *nats.Conn
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
	data, _ := json.Marshal(shelf)
	s.Nats.Publish(fmt.Sprintf("rooms.%v.shelves.create", &shelf.RoomID), []byte(data))
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

func (s *ShelfData) GetAvailableMovies(shelfID uuid.UUID, searchTerm string, excludeExisting bool) ([]themoviedb.Movie, error) {
	movieDB := themoviedb.NewMovieDBOptions(s.Env.MovieDBAuthToken, "")

	movies, err := movieDB.SearchMovies(searchTerm)
	if err != nil {
		return nil, err
	}

	if excludeExisting {
		var existingMovies []Movie
		if err := s.DB.Model(&existingMovies).Where("shelf_id = ?", shelfID).Select(); err != nil {
			return nil, err
		}

		existingMovieMap := make(map[uint]struct{})
		for _, movie := range existingMovies {
			existingMovieMap[movie.MovieID] = struct{}{}
		}

		var availableMovies []themoviedb.Movie
		for _, movie := range movies {
			if _, exists := existingMovieMap[movie.ID]; !exists {
				availableMovies = append(availableMovies, movie)
			}
		}

		return availableMovies, nil
	}

	return movies, nil
}

func (s *ShelfData) GetShelfAccess(shelfID, userID uuid.UUID) (bool, error) {
	var room Room

	query := s.DB.Model(&room).
		Join(`JOIN room_users ru ON "room".id = "ru".room_id`).
		Join(`JOIN shelves s ON "s".room_id = "ru".room_id`).
		Where(`"s".id = ? AND "ru".user_id = ?`, &shelfID, &userID).
		Select()

	if query == pg.ErrNoRows {
		return false, nil
	}

	if query != nil {
		return false, nil
	}

	return true, nil

}
