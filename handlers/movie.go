package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adamelfsborg-code/movie-nest/data"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MovieHandler struct {
	Data data.MovieData
}

func (m *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var body struct {
		MovieID uint      `json:"movie_id"`
		ShelfID uuid.UUID `json:"shelf_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie := data.NewMovie(body.MovieID, body.ShelfID)
	err = m.Data.CreateMovie(*movie)
	if err != nil {
		fmt.Println("Failed to create movie: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (m *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "movie_id")

	movieID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		fmt.Println("Failed to convert param to uint: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie, err := m.Data.GetMovie(uint(movieID))
	if err != nil {
		fmt.Println("Failed to get movie: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(movie)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (m *MovieHandler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "movie_id")

	movieID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie, err := m.Data.GetMovieDetails(movieID)
	if err != nil {
		fmt.Println("Failed to get movie: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(movie)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (m *MovieHandler) RateMovie(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Rating float64 `json:"rating"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println("Failed to decode json: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "movie_id")

	movieID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam = r.Header.Get("X-UserID")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Failed to parse id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movieRating := &data.MovieRating{
		MovieID: movieID,
		UserID:  userID,
		Rating:  body.Rating,
	}

	err = m.Data.RateMovie(*movieRating)
	if err != nil {
		fmt.Println("Failed to search movies: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
