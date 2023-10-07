package data

import (
	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/db"
	"github.com/adamelfsborg-code/movie-nest/pkg/themoviedb"
	"github.com/google/uuid"
)

type MovieData struct {
	Env config.Environments
}

type Movie struct {
	ID      uuid.UUID `json:"id" db:"id"`
	MovieID uint      `json:"movie_id" db:"movie_id"`
	ShelfID uuid.UUID `json:"shelf_id" db:"shelf_id"`
}

func NewMovie(movieID uint, shelfID uuid.UUID) *Movie {
	return &Movie{
		MovieID: movieID,
		ShelfID: shelfID,
	}
}

func (m *MovieData) CreateMovie(movie Movie) error {
	_, err := db.Store.Model(&movie).Insert()
	return err
}

func (m *MovieData) GetMovie(movieID uint) (*themoviedb.Movie, error) {
	movieDB := themoviedb.NewMovieDBOptions(m.Env.MovieDBAuthToken, "")

	movie, err := movieDB.GetMovie(movieID)
	if err != nil {
		return nil, err
	}

	return movie, nil
}
