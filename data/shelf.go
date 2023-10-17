package data

import (
	"time"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/pkg/themoviedb"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type ShelfData struct {
	DB  pg.DB
	Env config.Environments
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

func (s *ShelfData) GetAvailableMovies(shelfID uuid.UUID, searchTerm string, excludeExisting bool) ([]themoviedb.Movie, error) {
	movieDB := themoviedb.NewMovieDBOptions(s.Env.MovieDBAuthToken, "")

	movies, err := movieDB.SearchMovies(searchTerm)
	if err != nil {
		return nil, err
	}

	if excludeExisting {
		// Get existing movies for the specified shelf.
		var existingMovies []Movie
		if err := s.DB.Model(&existingMovies).Where("shelf_id = ?", shelfID).Select(); err != nil {
			return nil, err
		}

		// Create a map to efficiently check if a movie exists in existingMovies.
		existingMovieMap := make(map[uint]struct{})
		for _, movie := range existingMovies {
			existingMovieMap[movie.MovieID] = struct{}{}
		}

		// Filter out movies that already exist in the specified shelf.
		var availableMovies []themoviedb.Movie
		for _, movie := range movies {
			if _, exists := existingMovieMap[movie.ID]; !exists {
				availableMovies = append(availableMovies, movie)
			}
		}

		return availableMovies, nil
	}

	// get exluded movies, movies that does not exists inside existingMovies

	return movies, nil
}
