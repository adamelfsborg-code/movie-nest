package data

import (
	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/pkg/themoviedb"
)

type MovieData struct {
	Env config.Environments
}

func (u *MovieData) GetMovie(movieID uint) (*themoviedb.Movie, error) {
	movieDB := themoviedb.NewMovieDBOptions(u.Env.MovieDBAuthToken, "")

	movie, err := movieDB.GetMovie(movieID)
	if err != nil {
		return nil, err
	}

	return movie, nil
}
