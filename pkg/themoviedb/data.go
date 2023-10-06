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

func (m *MovieDBOptions) GetMovie(movieID uint) (*Movie, error) {
	byteMovie, err := m.Get(fmt.Sprintf("movie/%v", movieID))
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
