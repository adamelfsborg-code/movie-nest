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

type MovieData struct {
	Env  config.Environments
	DB   pg.DB
	Nats *nats.Conn
}

type Movie struct {
	ID      uuid.UUID `json:"id" db:"id"`
	MovieID uint      `json:"movie_id" db:"movie_id"`
	ShelfID uuid.UUID `json:"shelf_id" db:"shelf_id"`
}

type MovieAvgRating struct {
	MovieID uuid.UUID `json:"movie_id" db:"movie_id"`
	Rating  float64   `json:"rating" db:"rating"`
}

type MovieDetails struct {
	Movie          Movie             `json:"movie" db:"movie"`
	MovieDetails   themoviedb.Movie  `json:"details" db:"details"`
	MovieAvgRating MovieAvgRating    `json:"avg_rating" db:"avg_rating"`
	MovieRatings   []MovieRatingResp `json:"ratings" db:"ratings"`
}

type MovieRating struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	MovieID   uuid.UUID `json:"movie_id" db:"movie_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Rating    float64   `json:"rating" db:"rating"`
}

type MovieRatingResp struct {
	User      UserResp  `json:"user" db:"user"`
	Rating    float64   `json:"rating" db:"rating"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func NewMovie(movieID uint, shelfID uuid.UUID) *Movie {
	return &Movie{
		MovieID: movieID,
		ShelfID: shelfID,
	}
}

func (m *MovieData) CreateMovie(movie Movie) error {
	_, err := m.DB.Model(&movie).Insert()
	data, _ := json.Marshal(&movie)
	m.Nats.Publish(fmt.Sprintf("shelves.%v.movies.new", &movie.ShelfID), []byte(data))
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

func (m *MovieData) GetMovieDetails(movieID uuid.UUID) (*MovieDetails, error) {
	movieDB := themoviedb.NewMovieDBOptions(m.Env.MovieDBAuthToken, "")

	movie := &Movie{}

	err := m.DB.Model(movie).Where("id = ?", &movieID).Select()
	if err != nil {
		return nil, err
	}

	details, err := movieDB.GetMovie(movie.MovieID)
	if err != nil {
		return nil, err
	}

	var movieRatingResp []MovieRatingResp

	m.DB.Query(&movieRatingResp, `
		SELECT 
			jsonb_build_object
			(
				'id', u.id, 'name', u."name", 'timestamp', u."timestamp"
			) AS user,
			mr."timestamp",
			mr.rating
		FROM movie_ratings mr
		JOIN users u ON mr.user_id = u.id
		WHERE mr.movie_id = ?
	`, &movieID)

	var avgRating MovieAvgRating

	m.DB.Query(&avgRating, `
		SELECT mr.movie_id, avg(rating) as rating
		FROM movie_ratings mr
		WHERE mr.movie_id = ?
		GROUP BY mr.movie_id
	`, &movieID)

	movieDetails := &MovieDetails{
		Movie:          *movie,
		MovieDetails:   *details,
		MovieAvgRating: avgRating,
		MovieRatings:   movieRatingResp,
	}

	return movieDetails, nil
}

func (m *MovieData) RateMovie(rating MovieRating) error {
	_, err := m.DB.Model(&rating).Insert()
	data, _ := json.Marshal(&rating)
	m.Nats.Publish(fmt.Sprintf("movies.%v.rated", &rating.MovieID), []byte(data))
	return err
}
