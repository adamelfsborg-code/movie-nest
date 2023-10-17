package themoviedb

import (
	"encoding/json"
	"fmt"
)

type Movie struct {
	ID          uint   `json:"id"`
	Title       string `json:"original_title"`
	Poster      string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
}

type SearchMovieResp struct {
	Page   uint    `json:"page"`
	Movies []Movie `json:"results"`
}

func (m *MovieDBOptions) GetMovie(movieID uint) (*Movie, error) {
	byteMovie, err := m.get(fmt.Sprintf("movie/%v", movieID))
	if err != nil {
		return nil, fmt.Errorf("Failed to get movie: %w", err)
	}

	resp := &Movie{}
	json.Unmarshal(byteMovie, &resp)

	movie := &Movie{
		ID:          resp.ID,
		Title:       resp.Title,
		Poster:      resp.Poster,
		ReleaseDate: resp.ReleaseDate,
	}
	return movie, nil
}

func (m *MovieDBOptions) SearchMovies(searchTerm string) ([]Movie, error) {
	byteMovies, err := m.get(fmt.Sprintf("search/movie?query=%v", searchTerm))
	if err != nil {
		return nil, fmt.Errorf("Failed to get movie: %w", err)
	}

	resp := SearchMovieResp{}
	json.Unmarshal(byteMovies, &resp)

	return resp.Movies, nil
}
